package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"meta_commerce/app/global"
	"meta_commerce/app/handler"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
)

// FiberJwtToken là middleware xử lý xác thực và phân quyền người dùng thông qua JWT cho Fiber
type FiberJwtToken struct {
	UserCRUD           services.BaseServiceMongo[models.User]
	RoleCRUD           services.BaseServiceMongo[models.Role]
	PermissionCRUD     services.BaseServiceMongo[models.Permission]
	RolePermissionCRUD services.BaseServiceMongo[models.RolePermission]
	UserRoleCRUD       services.BaseServiceMongo[models.UserRole]
	Cache              *utility.Cache
}

// NewFiberJwtToken khởi tạo một instance mới của FiberJwtToken middleware
func NewFiberJwtToken() *FiberJwtToken {
	newHandler := new(FiberJwtToken)

	// Khởi tạo các collection từ registry
	userCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.Users)
	roleCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.Roles)
	permissionCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.Permissions)
	rolePermissionCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.RolePermissions)
	userRoleCol := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.UserRoles)

	// Khởi tạo các service với BaseService để thực hiện các thao tác CRUD
	newHandler.UserCRUD = services.NewBaseServiceMongo[models.User](userCol)
	newHandler.RoleCRUD = services.NewBaseServiceMongo[models.Role](roleCol)
	newHandler.PermissionCRUD = services.NewBaseServiceMongo[models.Permission](permissionCol)
	newHandler.RolePermissionCRUD = services.NewBaseServiceMongo[models.RolePermission](rolePermissionCol)
	newHandler.UserRoleCRUD = services.NewBaseServiceMongo[models.UserRole](userRoleCol)

	// Khởi tạo cache với thời gian sống 5 phút và thời gian dọn dẹp 10 phút
	newHandler.Cache = utility.NewCache(5*time.Minute, 10*time.Minute)

	return newHandler
}

// getUserPermissions lấy danh sách permissions của user từ cache hoặc database
func (jt *FiberJwtToken) getUserPermissions(userID string) (map[string]byte, error) {
	// Kiểm tra cache trước để tối ưu hiệu suất
	cacheKey := "user_permissions:" + userID
	if cached, found := jt.Cache.Get(cacheKey); found {
		return cached.(map[string]byte), nil
	}

	// Nếu không có trong cache, lấy từ database
	permissions := make(map[string]byte)

	// Lấy danh sách vai trò của user
	findRoles, err := jt.UserRoleCRUD.Find(context.TODO(), bson.M{"userId": utility.String2ObjectID(userID)}, nil)
	if err != nil {
		return nil, err
	}

	// Duyệt qua từng vai trò để lấy permissions
	for _, userRole := range findRoles {
		// Lấy danh sách permissions của vai trò
		findRolePermissions, err := jt.RolePermissionCRUD.Find(context.TODO(), bson.M{"roleId": userRole.RoleID}, nil)
		if err != nil {
			continue
		}

		// Lấy thông tin chi tiết của từng permission
		for _, rolePermission := range findRolePermissions {
			permission, err := jt.PermissionCRUD.FindOneById(context.TODO(), rolePermission.PermissionID)
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

// FiberAuthMiddleware middleware xác thực cho Fiber
func FiberAuthMiddleware(requirePermission string) fiber.Handler {
	h := &handler.FiberBaseHandler[interface{}, interface{}, interface{}]{}
	jt := NewFiberJwtToken()

	return func(c fiber.Ctx) error {
		// Lấy token từ header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Không tìm thấy token xác thực", utility.StatusUnauthorized, nil))
			return nil
		}

		// Kiểm tra định dạng token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Token không hợp lệ", utility.StatusUnauthorized, nil))
			return nil
		}

		token := parts[1]

		// Tìm user có token
		user, err := jt.UserCRUD.FindOne(context.Background(), bson.M{"tokens.jwtToken": token}, nil)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"token": token,
				"error": err.Error(),
			}).Error("Không tìm thấy user với token")
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Token không hợp lệ hoặc đã hết hạn", utility.StatusUnauthorized, nil))
			return nil
		}

		// Kiểm tra user có bị block không
		if user.IsBlock {
			logrus.WithFields(logrus.Fields{
				"userId": user.ID.Hex(),
				"reason": user.BlockNote,
			}).Warn("User bị block truy cập hệ thống")
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Tài khoản đã bị khóa: "+user.BlockNote, utility.StatusForbidden, nil))
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
		permissions, err := jt.getUserPermissions(user.ID.Hex())
		if err != nil {
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Không thể lấy thông tin quyền", utility.StatusForbidden, nil))
			return nil
		}

		// Kiểm tra user có permission cần thiết không
		scope, hasPermission := permissions[requirePermission]
		if !hasPermission {
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Không có quyền truy cập", utility.StatusForbidden, nil))
			return nil
		}

		// Lưu scope tối thiểu vào context để sử dụng trong handler
		c.Locals("minScope", scope)
		return c.Next()
	}
}

// FiberRoleMiddleware middleware kiểm tra quyền cho Fiber
func FiberRoleMiddleware(requiredRoles ...string) fiber.Handler {
	h := &handler.FiberBaseHandler[interface{}, interface{}, interface{}]{}
	return func(c fiber.Ctx) error {
		// Lấy user từ context
		user, ok := c.Locals("user").(models.User)
		if !ok {
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Không tìm thấy thông tin người dùng", utility.StatusUnauthorized, nil))
			return nil
		}

		// Kiểm tra token
		if user.Token == "" {
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Không có quyền truy cập", utility.StatusUnauthorized, nil))
			return nil
		}

		// Lấy role từ token
		collection := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.Roles)
		var role models.Role
		err := collection.FindOne(context.Background(), bson.M{"_id": utility.String2ObjectID(user.Token)}).Decode(&role)
		if err != nil {
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Không tìm thấy thông tin quyền", utility.StatusUnauthorized, nil))
			return nil
		}

		// Kiểm tra role có trong danh sách yêu cầu không
		hasRole := false
		for _, requiredRole := range requiredRoles {
			if role.Name == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Không có quyền truy cập", utility.StatusForbidden, nil))
			return nil
		}

		return c.Next()
	}
}
