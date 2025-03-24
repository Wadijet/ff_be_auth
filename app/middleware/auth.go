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

// JwtToken , basic jwt model
type JwtToken struct {
	C                  *config.Configuration
	UserCRUD           services.BaseService[models.User]
	RoleCRUD           services.BaseService[models.Role]
	PermissionCRUD     services.BaseService[models.Permission]
	RolePermissionCRUD services.BaseService[models.RolePermission]
	UserRoleCRUD       services.BaseService[models.UserRole]
	Cache              *utility.Cache // Cache cho permissions và roles
}

// NewJwtToken , khởi tạo một JwtToken mới
func NewJwtToken(c *config.Configuration, db *mongo.Client) *JwtToken {
	newHandler := new(JwtToken)
	newHandler.C = c

	// Khởi tạo các collection
	userCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Users)).Collection(global.MongoDB_ColNames.Users)
	roleCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Roles)).Collection(global.MongoDB_ColNames.Roles)
	permissionCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Permissions)).Collection(global.MongoDB_ColNames.Permissions)
	rolePermissionCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.RolePermissions)).Collection(global.MongoDB_ColNames.RolePermissions)
	userRoleCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.UserRoles)).Collection(global.MongoDB_ColNames.UserRoles)

	// Khởi tạo các service với BaseService
	newHandler.UserCRUD = services.NewBaseService[models.User](userCol)
	newHandler.RoleCRUD = services.NewBaseService[models.Role](roleCol)
	newHandler.PermissionCRUD = services.NewBaseService[models.Permission](permissionCol)
	newHandler.RolePermissionCRUD = services.NewBaseService[models.RolePermission](rolePermissionCol)
	newHandler.UserRoleCRUD = services.NewBaseService[models.UserRole](userRoleCol)

	// Khởi tạo cache
	newHandler.Cache = utility.NewCache(5*time.Minute, 10*time.Minute)

	return newHandler
}

// validateToken kiểm tra và xác thực JWT token
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
func (jt *JwtToken) getUserPermissions(userID string) (map[string]byte, error) {
	// Kiểm tra cache trước
	cacheKey := "user_permissions:" + userID
	if cached, found := jt.Cache.Get(cacheKey); found {
		return cached.(map[string]byte), nil
	}

	// Nếu không có trong cache, lấy từ database
	permissions := make(map[string]byte)
	findRoles, err := jt.UserRoleCRUD.FindAll(context.TODO(), bson.M{"userId": utility.String2ObjectID(userID)}, nil)
	if err != nil {
		return nil, err
	}

	for _, userRole := range findRoles {
		findRolePermissions, err := jt.RolePermissionCRUD.FindAll(context.TODO(), bson.M{"roleId": userRole.RoleID}, nil)
		if err != nil {
			continue
		}

		for _, rolePermission := range findRolePermissions {
			permission, err := jt.PermissionCRUD.FindOne(context.TODO(), rolePermission.PermissionID)
			if err != nil {
				continue
			}
			permissions[permission.Name] = rolePermission.Scope
		}
	}

	// Lưu vào cache
	jt.Cache.Set(cacheKey, permissions)
	return permissions, nil
}

// CheckUserAuth , kiểm tra xác thực người dùng
func (jt *JwtToken) CheckUserAuth(requirePermission string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		notAuthError := "An unauthorized access!"
		notPermissionError := "You do not have permission to perform the action!"

		jwtTokenString := string(ctx.Request.Header.Peek("Authorization"))
		if jwtTokenString == "" {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			return
		}

		splitToken := strings.Split(jwtTokenString, "Bearer ")
		if len(splitToken) != 2 {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			return
		}

		// Xác thực token
		t, err := jt.validateToken(splitToken[1])
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, err, notAuthError))
			return
		}

		// Kiểm tra user
		findUser, err := jt.UserCRUD.FindOne(context.TODO(), utility.String2ObjectID(t.UserID))
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, err, notAuthError))
			return
		}

		if findUser.IsBlock {
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			return
		}

		// Kiểm tra token có hợp lệ không
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
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			return
		}

		// Nếu không yêu cầu permission cụ thể
		if requirePermission == "" {
			next(ctx)
			return
		}

		// Kiểm tra permission
		permissions, err := jt.getUserPermissions(t.UserID)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
			return
		}

		scope, hasPermission := permissions[requirePermission]
		if !hasPermission {
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			utility.JSON(ctx, utility.Payload(false, nil, notPermissionError))
			return
		}

		ctx.SetUserValue("minScope", scope)
		next(ctx)
	}
}
