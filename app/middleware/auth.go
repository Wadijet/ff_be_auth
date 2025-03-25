package middleware

import (
	"context"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
)

// JwtToken là middleware xử lý xác thực và phân quyền người dùng thông qua JWT
// Các thành phần chính:
// - C: Cấu hình hệ thống
// - UserCRUD: Service quản lý thông tin người dùng
// - RoleCRUD: Service quản lý vai trò
// - PermissionCRUD: Service quản lý quyền hạn
// - RolePermissionCRUD: Service quản lý phân quyền cho vai trò
// - UserRoleCRUD: Service quản lý phân vai trò cho người dùng
// - Cache: Cache để lưu trữ tạm thời permissions và roles
type JwtToken struct {
	C                  *config.Configuration
	UserCRUD           services.BaseService[models.User]
	RoleCRUD           services.BaseService[models.Role]
	PermissionCRUD     services.BaseService[models.Permission]
	RolePermissionCRUD services.BaseService[models.RolePermission]
	UserRoleCRUD       services.BaseService[models.UserRole]
	Cache              *utility.Cache
}

// NewJwtToken khởi tạo một instance mới của JwtToken middleware
// Parameters:
//   - c: Cấu hình hệ thống
//   - db: Kết nối MongoDB
//
// Returns:
//   - *JwtToken: Instance mới của JwtToken middleware
func NewJwtToken(c *config.Configuration, db *mongo.Client) *JwtToken {
	newHandler := new(JwtToken)
	newHandler.C = c

	// Khởi tạo các collection từ database
	userCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Users)).Collection(global.MongoDB_ColNames.Users)
	roleCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Roles)).Collection(global.MongoDB_ColNames.Roles)
	permissionCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Permissions)).Collection(global.MongoDB_ColNames.Permissions)
	rolePermissionCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.RolePermissions)).Collection(global.MongoDB_ColNames.RolePermissions)
	userRoleCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.UserRoles)).Collection(global.MongoDB_ColNames.UserRoles)

	// Khởi tạo các service với BaseService để thực hiện các thao tác CRUD
	newHandler.UserCRUD = services.NewBaseService[models.User](userCol)
	newHandler.RoleCRUD = services.NewBaseService[models.Role](roleCol)
	newHandler.PermissionCRUD = services.NewBaseService[models.Permission](permissionCol)
	newHandler.RolePermissionCRUD = services.NewBaseService[models.RolePermission](rolePermissionCol)
	newHandler.UserRoleCRUD = services.NewBaseService[models.UserRole](userRoleCol)

	// Khởi tạo cache với thời gian sống 5 phút và thời gian dọn dẹp 10 phút
	newHandler.Cache = utility.NewCache(5*time.Minute, 10*time.Minute)

	return newHandler
}

// validateToken kiểm tra và xác thực JWT token
// Parameters:
//   - tokenString: Chuỗi JWT token cần xác thực
//
// Returns:
//   - *models.JwtToken: Thông tin token đã được xác thực
//   - error: Lỗi nếu có
func (jt *JwtToken) validateToken(tokenString string) (*models.JwtToken, error) {
	t := models.JwtToken{}
	token, err := jwt.ParseWithClaims(tokenString, &t, func(token *jwt.Token) (interface{}, error) {
		return []byte(jt.C.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return &t, nil
}

// getUserPermissions lấy danh sách permissions của user từ cache hoặc database
// Parameters:
//   - userID: ID của user cần lấy permissions
//
// Returns:
//   - map[string]byte: Map chứa tên permission và scope tương ứng
//   - error: Lỗi nếu có
func (jt *JwtToken) getUserPermissions(userID string) (map[string]byte, error) {
	// Kiểm tra cache trước để tối ưu hiệu suất
	cacheKey := "user_permissions:" + userID
	if cached, found := jt.Cache.Get(cacheKey); found {
		return cached.(map[string]byte), nil
	}

	// Nếu không có trong cache, lấy từ database
	permissions := make(map[string]byte)

	// Lấy danh sách vai trò của user
	findRoles, err := jt.UserRoleCRUD.FindAll(context.TODO(), bson.M{"userId": utility.String2ObjectID(userID)}, nil)
	if err != nil {
		return nil, err
	}

	// Duyệt qua từng vai trò để lấy permissions
	for _, userRole := range findRoles {
		// Lấy danh sách permissions của vai trò
		findRolePermissions, err := jt.RolePermissionCRUD.FindAll(context.TODO(), bson.M{"roleId": userRole.RoleID}, nil)
		if err != nil {
			continue
		}

		// Lấy thông tin chi tiết của từng permission
		for _, rolePermission := range findRolePermissions {
			permission, err := jt.PermissionCRUD.FindOne(context.TODO(), rolePermission.PermissionID)
			if err != nil {
				continue
			}
			permissions[permission.Name] = rolePermission.Scope
		}
	}

	// Lưu vào cache để sử dụng cho các lần sau
	jt.Cache.Set(cacheKey, permissions)
	return permissions, nil
}

// CheckUserAuth là middleware kiểm tra xác thực và phân quyền người dùng
// Parameters:
//   - requirePermission: Tên permission cần thiết để truy cập endpoint
//   - next: Handler tiếp theo trong chuỗi middleware
//
// Returns:
//   - fasthttp.RequestHandler: Handler đã được bọc với logic xác thực và phân quyền
func (jt *JwtToken) CheckUserAuth(requirePermission string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Lấy token từ header Authorization
		jwtTokenString := string(ctx.Request.Header.Peek("Authorization"))
		if jwtTokenString == "" {
			ctx.SetStatusCode(utility.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, nil, utility.MsgTokenMissing))
			return
		}

		// Tách token từ chuỗi "Bearer <token>"
		splitToken := strings.Split(jwtTokenString, "Bearer ")
		if len(splitToken) != 2 {
			ctx.SetStatusCode(utility.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, nil, utility.MsgTokenInvalid))
			return
		}

		// Xác thực token
		t, err := jt.validateToken(splitToken[1])
		if err != nil {
			ctx.SetStatusCode(utility.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, err, utility.MsgTokenInvalid))
			return
		}

		// Kiểm tra user có tồn tại không
		findUser, err := jt.UserCRUD.FindOne(context.TODO(), utility.String2ObjectID(t.UserID))
		if err != nil {
			ctx.SetStatusCode(utility.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, err, utility.MsgUnauthorized))
			return
		}

		// Kiểm tra user có bị block không
		if findUser.IsBlock {
			ctx.SetStatusCode(utility.StatusForbidden)
			utility.JSON(ctx, utility.Payload(false, nil, utility.MsgForbidden))
			return
		}

		// Kiểm tra token có hợp lệ với user không
		isRightToken := false
		for _, _token := range findUser.Tokens {
			if _token.JwtToken == splitToken[1] {
				ctx.SetUserValue("userId", t.UserID)
				ctx.SetUserValue("userToken", _token.JwtToken)
				isRightToken = true
				break
			}
		}

		if !isRightToken {
			ctx.SetStatusCode(utility.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, nil, utility.MsgTokenInvalid))
			return
		}

		// Nếu không yêu cầu permission cụ thể, cho phép truy cập
		if requirePermission == "" {
			next(ctx)
			return
		}

		// Kiểm tra permission của user
		permissions, err := jt.getUserPermissions(t.UserID)
		if err != nil {
			ctx.SetStatusCode(utility.StatusForbidden)
			utility.JSON(ctx, utility.Payload(false, err, utility.MsgForbidden))
			return
		}

		// Kiểm tra user có permission cần thiết không
		scope, hasPermission := permissions[requirePermission]
		if !hasPermission {
			ctx.SetStatusCode(utility.StatusForbidden)
			utility.JSON(ctx, utility.Payload(false, nil, utility.MsgForbidden))
			return
		}

		// Lưu scope tối thiểu vào context để sử dụng trong handler
		ctx.SetUserValue("minScope", scope)
		next(ctx)
	}
}
