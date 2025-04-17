package middleware

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"

	"meta_commerce/core/api/handler"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/utility"
)

// AuthManager quản lý xác thực và phân quyền người dùng
type AuthManager struct {
	UserCRUD           *services.UserService
	RoleCRUD           *services.RoleService
	PermissionCRUD     *services.PermissionService
	RolePermissionCRUD *services.RolePermissionService
	UserRoleCRUD       *services.UserRoleService
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
		var err error
		authManagerInstance, err = newAuthManager()
		if err != nil {
			panic(err)
		}
	})
	return authManagerInstance
}

// newAuthManager khởi tạo một instance mới của AuthManager (private constructor)
func newAuthManager() (*AuthManager, error) {
	newManager := new(AuthManager)

	// Khởi tạo các service với BaseService để thực hiện các thao tác CRUD
	userService, err := services.NewUserService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %v", err)
	}
	newManager.UserCRUD = userService

	roleService, err := services.NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}
	newManager.RoleCRUD = roleService

	permissionService, err := services.NewPermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create permission service: %v", err)
	}
	newManager.PermissionCRUD = permissionService

	rolePermissionService, err := services.NewRolePermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role permission service: %v", err)
	}
	newManager.RolePermissionCRUD = rolePermissionService

	userRoleService, err := services.NewUserRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user role service: %v", err)
	}
	newManager.UserRoleCRUD = userRoleService

	// Khởi tạo cache với thời gian sống 5 phút và thời gian dọn dẹp 10 phút
	newManager.Cache = utility.NewCache(5*time.Minute, 10*time.Minute)

	// Khởi tạo response handler một lần duy nhất
	newManager.responseHandler = &handler.BaseHandler[models.User, models.UserCreateInput, models.UserChangeInfoInput]{
		BaseService: newManager.UserCRUD,
	}

	return newManager, nil
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
		return nil, common.ConvertMongoError(err)
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
			authManager.responseHandler.HandleResponse(c, nil, common.ErrTokenMissing)
			return nil
		}

		// Kiểm tra định dạng token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			authManager.responseHandler.HandleResponse(c, nil, common.ErrTokenInvalid)
			return nil
		}

		token := parts[1]

		// Tìm user có token
		user, err := authManager.UserCRUD.FindOne(context.Background(), bson.M{"tokens.jwtToken": token}, nil)
		if err != nil {
			authManager.responseHandler.HandleResponse(c, nil, common.ErrTokenInvalid)
			return nil
		}

		// Kiểm tra user có bị block không
		if user.IsBlock {
			authManager.responseHandler.HandleResponse(c, nil, common.NewError(
				common.ErrCodeAuthCredentials,
				"Tài khoản đã bị khóa: "+user.BlockNote,
				common.StatusForbidden,
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
			authManager.responseHandler.HandleResponse(c, nil, common.NewError(
				common.ErrCodeAuthRole,
				"Không thể lấy thông tin quyền",
				common.StatusForbidden,
				nil,
			))
			return nil
		}

		// Kiểm tra user có permission cần thiết không
		scope, hasPermission := permissions[requirePermission]
		if !hasPermission {
			authManager.responseHandler.HandleResponse(c, nil, common.NewError(
				common.ErrCodeAuthRole,
				"Không có quyền truy cập",
				common.StatusForbidden,
				nil,
			))
			return nil
		}

		// Lưu scope tối thiểu vào context để sử dụng trong handler
		c.Locals("minScope", scope)
		return c.Next()
	}
}
