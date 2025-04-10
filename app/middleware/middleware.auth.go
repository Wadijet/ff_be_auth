package middleware

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"

	"meta_commerce/app/global"
	"meta_commerce/app/handler"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
)

// AuthManager quản lý xác thực và phân quyền người dùng
type AuthManager struct {
	UserCRUD           services.BaseServiceMongo[models.User]
	RoleCRUD           services.BaseServiceMongo[models.Role]
	PermissionCRUD     services.BaseServiceMongo[models.Permission]
	RolePermissionCRUD services.BaseServiceMongo[models.RolePermission]
	UserRoleCRUD       services.BaseServiceMongo[models.UserRole]
	Cache              *utility.Cache
	responseHandler    *handler.BaseHandler[models.User, models.UserCreateInput, models.UserChangeInfoInput]
}

var (
	authManagerInstance *AuthManager
	authManagerOnce     sync.Once
)

// GetAuthManager trả về instance duy nhất của AuthManager (singleton pattern)
func GetAuthManager() *AuthManager {
	authManagerOnce.Do(func() {
		authManagerInstance = newAuthManager()
	})
	return authManagerInstance
}

// newAuthManager khởi tạo một instance mới của AuthManager (private constructor)
func newAuthManager() *AuthManager {
	newManager := new(AuthManager)

	// Khởi tạo các collection từ registry
	userCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.Users)
	roleCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.Roles)
	permissionCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.Permissions)
	rolePermissionCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.RolePermissions)
	userRoleCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.UserRoles)

	// Khởi tạo các service với BaseService để thực hiện các thao tác CRUD
	newManager.UserCRUD = services.NewBaseServiceMongo[models.User](userCol)
	newManager.RoleCRUD = services.NewBaseServiceMongo[models.Role](roleCol)
	newManager.PermissionCRUD = services.NewBaseServiceMongo[models.Permission](permissionCol)
	newManager.RolePermissionCRUD = services.NewBaseServiceMongo[models.RolePermission](rolePermissionCol)
	newManager.UserRoleCRUD = services.NewBaseServiceMongo[models.UserRole](userRoleCol)

	// Khởi tạo cache với thời gian sống 5 phút và thời gian dọn dẹp 10 phút
	newManager.Cache = utility.NewCache(5*time.Minute, 10*time.Minute)

	// Khởi tạo response handler một lần duy nhất
	newManager.responseHandler = &handler.BaseHandler[models.User, models.UserCreateInput, models.UserChangeInfoInput]{
		Service: newManager.UserCRUD,
	}

	return newManager
}

// getUserPermissions lấy danh sách permissions của user từ cache hoặc database
func (am *AuthManager) getUserPermissions(userID string) (map[string]byte, error) {
	// Kiểm tra cache trước để tối ưu hiệu suất
	cacheKey := "user_permissions:" + userID
	if cached, found := am.Cache.Get(cacheKey); found {
		return cached.(map[string]byte), nil
	}

	// Nếu không có trong cache, lấy từ database
	permissions := make(map[string]byte)

	// Lấy danh sách vai trò của user
	findRoles, err := am.UserRoleCRUD.Find(context.TODO(), bson.M{"userId": utility.String2ObjectID(userID)}, nil)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	// Duyệt qua từng vai trò để lấy permissions
	for _, userRole := range findRoles {
		// Lấy danh sách permissions của vai trò
		findRolePermissions, err := am.RolePermissionCRUD.Find(context.TODO(), bson.M{"roleId": userRole.RoleID}, nil)
		if err != nil {
			continue
		}

		// Lấy thông tin chi tiết của từng permission
		for _, rolePermission := range findRolePermissions {
			permission, err := am.PermissionCRUD.FindOneById(context.TODO(), rolePermission.PermissionID)
			if err != nil {
				continue
			}
			permissions[permission.Name] = rolePermission.Scope
		}
	}

	// Lưu vào cache để sử dụng cho các lần sau
	am.Cache.Set(cacheKey, permissions)
	return permissions, nil
}

// AuthMiddleware middleware xác thực cho Fiber
func AuthMiddleware(requirePermission string) fiber.Handler {
	// Sử dụng singleton instance của AuthManager
	authManager := GetAuthManager()

	return func(c fiber.Ctx) error {

		// Lấy token từ header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			authManager.responseHandler.HandleResponse(c, nil, utility.ErrTokenMissing)
			return nil
		}

		// Kiểm tra định dạng token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			authManager.responseHandler.HandleResponse(c, nil, utility.ErrTokenInvalid)
			return nil
		}

		token := parts[1]

		// Tìm user có token
		user, err := authManager.UserCRUD.FindOne(context.Background(), bson.M{"tokens.jwtToken": token}, nil)
		if err != nil {
			authManager.responseHandler.HandleResponse(c, nil, utility.ErrTokenInvalid)
			return nil
		}

		// Kiểm tra user có bị block không
		if user.IsBlock {
			authManager.responseHandler.HandleResponse(c, nil, utility.NewError(
				utility.ErrCodeAuthCredentials,
				"Tài khoản đã bị khóa: "+user.BlockNote,
				utility.StatusForbidden,
				nil,
			))
			return nil
		}

		// Lưu thông tin user vào context
		c.Locals("userId", user.ID.Hex())
		c.Locals("user", user)

		// Nếu không yêu cầu permission cụ thể, cho phép truy cập
		if requirePermission == "" {
			return c.Next()
		}

		// Kiểm tra permission của user
		permissions, err := authManager.getUserPermissions(user.ID.Hex())
		if err != nil {
			authManager.responseHandler.HandleResponse(c, nil, utility.NewError(
				utility.ErrCodeAuthRole,
				"Không thể lấy thông tin quyền",
				utility.StatusForbidden,
				nil,
			))
			return nil
		}

		// Kiểm tra user có permission cần thiết không
		scope, hasPermission := permissions[requirePermission]
		if !hasPermission {
			authManager.responseHandler.HandleResponse(c, nil, utility.NewError(
				utility.ErrCodeAuthRole,
				"Không có quyền truy cập",
				utility.StatusForbidden,
				nil,
			))
			return nil
		}

		// Lưu scope tối thiểu vào context để sử dụng trong handler
		c.Locals("minScope", scope)
		return c.Next()
	}
}
