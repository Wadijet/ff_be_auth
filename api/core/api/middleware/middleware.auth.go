package middleware

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/logger"
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

	return newManager, nil
}

// getUserPermissions lấy danh sách permissions của user từ cache hoặc database
// Nếu activeRoleID được cung cấp, chỉ lấy permissions từ role đó (role context)
// Nếu activeRoleID là nil, lấy permissions từ tất cả roles của user (backward compatibility)
func (am *AuthManager) getUserPermissions(userID string, activeRoleID *primitive.ObjectID) (map[string]byte, error) {
	// Tạo cache key dựa trên userID và activeRoleID (nếu có)
	var cacheKey string
	if activeRoleID != nil {
		cacheKey = fmt.Sprintf("user_permissions:%s:role:%s", userID, activeRoleID.Hex())
	} else {
		cacheKey = "user_permissions:" + userID
	}

	// Kiểm tra cache trước để tối ưu hiệu suất
	if cached, found := am.Cache.Get(cacheKey); found {
		return cached.(map[string]byte), nil
	}

	// Nếu không có trong cache, lấy từ database
	permissions := make(map[string]byte)

	// Nếu có activeRoleID, chỉ lấy permissions từ role đó
	if activeRoleID != nil {
		// Validate user có role này không
		_, err := am.UserRoleCRUD.FindOne(context.TODO(), bson.M{
			"userId": utility.String2ObjectID(userID),
			"roleId": *activeRoleID,
		}, nil)
		if err != nil {
			// User không có role này, trả về map rỗng
			am.Cache.Set(cacheKey, permissions)
			return permissions, nil
		}

		// Lấy danh sách permissions của role
		findRolePermissions, err := am.RolePermissionCRUD.Find(context.TODO(), bson.M{"roleId": *activeRoleID}, nil)
		if err != nil {
			am.Cache.Set(cacheKey, permissions)
			return permissions, nil
		}

		// Lấy thông tin chi tiết của từng permission
		for _, rolePermission := range findRolePermissions {
			permission, err := am.PermissionCRUD.FindOneById(context.TODO(), rolePermission.PermissionID)
			if err != nil {
				continue
			}
			permissions[permission.Name] = rolePermission.Scope
		}
	} else {
		// Lấy permissions từ tất cả roles của user (backward compatibility)
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
	}

	// Lưu vào cache để sử dụng cho các lần sau
	am.Cache.Set(cacheKey, permissions)
	return permissions, nil
}

// AuthMiddleware middleware xác thực cho Fiber
func AuthMiddleware(requirePermission string) fiber.Handler {
	// Log khi tạo middleware instance
	fmt.Printf("[AUTH] ⚙️ Creating AuthMiddleware with permission: %s\n", requirePermission)

	// Sử dụng singleton instance của AuthManager
	authManager := GetAuthManager()

	return func(c fiber.Ctx) error {
		// Log ngay đầu hàm để xác nhận middleware được gọi - dùng GetAppLogger để ghi vào file
		logger.GetAppLogger().WithFields(logrus.Fields{
			"path":       c.Path(),
			"method":     c.Method(),
			"permission": requirePermission,
		}).Error("🔒 [AUTH] AuthMiddleware EXECUTING - FORCE LOG - FIRST LINE")

		// Lấy token từ header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Ghi log vào file để debug
			logrus.WithFields(logrus.Fields{
				"path": c.Path(),
				"method": c.Method(),
			}).Error("❌ Missing Authorization header")
			HandleErrorResponse(c, common.ErrTokenMissing)
			return nil
		}
		
		// Log để đảm bảo middleware được gọi - dùng GetAppLogger để ghi vào file
		logger.GetAppLogger().WithFields(logrus.Fields{
			"path": c.Path(),
			"method": c.Method(),
			"has_auth_header": authHeader != "",
		}).Error("🔍 [AUTH] AuthMiddleware processing request - FORCE LOG")

		// Kiểm tra định dạng token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			HandleErrorResponse(c, common.ErrTokenInvalid)
			return nil
		}

		token := parts[1]

		// Log token để debug - dùng Info level để đảm bảo hiển thị
		tokenPreview := token
		if len(token) > 50 {
			tokenPreview = token[:50] + "..."
		}
		// Log token để debug - dùng GetAppLogger để ghi vào file
		logger.GetAppLogger().WithFields(logrus.Fields{
			"path":         c.Path(),
			"token":        tokenPreview,
			"token_length": len(token),
		}).Error("🔍 [AUTH] Searching for user with token - FORCE LOG")

		// Tìm user có token
		// Ưu tiên query field "token" (token mới nhất) trước vì nó được cập nhật mỗi lần login
		// Nếu không tìm thấy, query trong array "tokens" (tokens theo hwid)
		var user models.User
		var err error
		var query bson.M

		// Cách 1: Query field "token" (token mới nhất) - ĐÂY LÀ CÁCH CHÍNH
		query = bson.M{"token": token}
		logger.GetAppLogger().WithFields(logrus.Fields{
			"path":         c.Path(),
			"query":        query,
			"token_length": len(token),
			"token_preview": tokenPreview,
		}).Error("🔍 [AUTH] Executing Query 1: token field - FORCE LOG")
		user, err = authManager.UserCRUD.FindOne(context.Background(), query, nil)
		if err != nil {
			logger.GetAppLogger().WithFields(logrus.Fields{
				"path":    c.Path(),
				"query":   query,
				"error":   err.Error(),
				"has_user": false,
			}).Error("❌ [AUTH] Query 1 FAILED - FORCE LOG")
		} else {
			logger.GetAppLogger().WithFields(logrus.Fields{
				"path":     c.Path(),
				"query":    query,
				"user_id":  user.ID.Hex(),
				"has_user": true,
			}).Error("✅ [AUTH] Query 1 SUCCESS - FORCE LOG")
		}

		if err != nil {
			logger.GetAppLogger().WithFields(logrus.Fields{
				"path":  c.Path(),
				"error": err.Error(),
				"query": query,
			}).Error("⚠️ [AUTH] Token not found in 'token' field, trying 'tokens' array - FORCE LOG")

			// Cách 2: Query trong array "tokens" với dot notation
			query = bson.M{"tokens.jwtToken": token}
			user, err = authManager.UserCRUD.FindOne(context.Background(), query, nil)

			if err != nil {
				logger.GetAppLogger().WithFields(logrus.Fields{
					"path":  c.Path(),
					"error": err.Error(),
					"query": query,
				}).Error("⚠️ [AUTH] Query 2 failed, trying $elemMatch - FORCE LOG")

				// Cách 3: Query với $elemMatch
				query = bson.M{
					"tokens": bson.M{
						"$elemMatch": bson.M{
							"jwtToken": token,
						},
					},
				}
				user, err = authManager.UserCRUD.FindOne(context.Background(), query, nil)
				if err != nil {
					logger.GetAppLogger().WithFields(logrus.Fields{
						"path":  c.Path(),
						"error": err.Error(),
						"query": query,
					}).Error("⚠️ [AUTH] Query 3 ($elemMatch) also failed - FORCE LOG")
				}
			}
		}

		if err != nil {
			// Log chi tiết lỗi - dùng GetAppLogger để ghi vào file
			logger.GetAppLogger().WithFields(logrus.Fields{
				"path":  c.Path(),
				"error": err.Error(),
				"token": token[:20] + "...",
				"query": query,
			}).Error("❌ [AUTH] Token not found in database - FORCE LOG")
			// Log thêm thông tin query để debug
			logger.GetAppLogger().WithFields(logrus.Fields{
				"path":          c.Path(),
				"query":         query,
				"token_preview": token[:20] + "...",
			}).Error("❌ [AUTH] Token query details - FORCE LOG")
			HandleErrorResponse(c, common.ErrTokenInvalid)
			return nil
		}

		// Log khi tìm thấy token - dùng Info level để đảm bảo hiển thị
		fmt.Printf("[AUTH] ✅ Token found, user authenticated - Path: %s, UserID: %s\n",
			c.Path(), user.ID.Hex())
		logrus.WithFields(logrus.Fields{
			"path":    c.Path(),
			"user_id": user.ID.Hex(),
		}).Info("✅ Token found, user authenticated")

		// Kiểm tra user có bị block không
		if user.IsBlock {
			HandleErrorResponse(c, common.NewError(
				common.ErrCodeAuthCredentials,
				"Tài khoản đã bị khóa: "+user.BlockNote,
				common.StatusForbidden,
				nil,
			))
			return nil
		}

		// Lưu thông tin user vào context
		c.Locals("user_id", user.ID.Hex())
		c.Locals("user", user)

		// Nếu không yêu cầu permission cụ thể, cho phép truy cập NGAY
		// Đây là endpoint đặc biệt như /auth/roles - chỉ cần xác thực, không cần permission
		if requirePermission == "" {
			fmt.Printf("[AUTH] ✅ No permission required - Path: %s, UserID: %s - ALLOWING ACCESS\n",
				c.Path(), user.ID.Hex())
			logrus.WithFields(logrus.Fields{
				"path":    c.Path(),
				"user_id": user.ID.Hex(),
			}).Info("✅ No permission required - allowing access")
			return c.Next()
		}

		// Lấy active role ID từ header (role context)
		// Logic: Nếu route có require permission, PHẢI có header X-Active-Role-ID để chỉ định role context
		activeRoleIDStr := c.Get("X-Active-Role-ID")

		// Log tất cả headers để debug
		allHeaders := make(map[string]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			allHeaders[string(key)] = string(value)
		})
		fmt.Printf("[AUTH] 🔍 Checking headers - Path: %s, X-Active-Role-ID: %s, Permission: %s\n",
			c.Path(), activeRoleIDStr, requirePermission)
		logrus.WithFields(logrus.Fields{
			"path":               c.Path(),
			"x_active_role_id":   activeRoleIDStr,
			"require_permission": requirePermission,
			"all_headers":        allHeaders,
		}).Info("🔍 Checking headers and permission")

		// Header X-Active-Role-ID là BẮT BUỘC khi route yêu cầu permission
		if activeRoleIDStr == "" {
			fmt.Printf("[AUTH] ❌ BLOCKING: Missing X-Active-Role-ID header - User: %s, Path: %s\n",
				user.Email, c.Path())
			logrus.WithFields(logrus.Fields{
				"user_id":    user.ID.Hex(),
				"user_email": user.Email,
				"path":       c.Path(),
				"permission": requirePermission,
			}).Error("❌ Missing X-Active-Role-ID header - BLOCKING REQUEST")
			HandleErrorResponse(c, common.NewError(
				common.ErrCodeAuthRole,
				"Thiếu header X-Active-Role-ID. Vui lòng chọn role để làm việc.",
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		// Parse và validate role ID
		roleID, err := primitive.ObjectIDFromHex(activeRoleIDStr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id":        user.ID.Hex(),
				"active_role_id": activeRoleIDStr,
				"path":           c.Path(),
				"error":          err.Error(),
			}).Error("❌ Invalid X-Active-Role-ID format")
			HandleErrorResponse(c, common.NewError(
				common.ErrCodeValidationFormat,
				"X-Active-Role-ID không đúng định dạng",
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		// Lấy danh sách roles của user để kiểm tra
		userRoles, err := authManager.UserRoleCRUD.Find(context.Background(), bson.M{"userId": utility.String2ObjectID(user.ID.Hex())}, nil)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id": user.ID.Hex(),
				"error":   err.Error(),
				"path":    c.Path(),
			}).Error("Failed to get user roles")
			HandleErrorResponse(c, common.NewError(
				common.ErrCodeAuthRole,
				"Không thể kiểm tra quyền truy cập",
				common.StatusForbidden,
				nil,
			))
			return nil
		}

		// Log để debug - dùng Info level để đảm bảo hiển thị
		fmt.Printf("[AUTH] 🔐 Checking permissions - User: %s, Roles: %d, Path: %s, Permission: %s, ActiveRole: %s\n",
			user.Email, len(userRoles), c.Path(), requirePermission, roleID.Hex())
		logrus.WithFields(logrus.Fields{
			"user_id":        user.ID.Hex(),
			"user_email":     user.Email,
			"roles_count":    len(userRoles),
			"path":           c.Path(),
			"permission":     requirePermission,
			"active_role_id": roleID.Hex(),
		}).Info("🔐 Checking user permissions")

		// Nếu user không có role nào, từ chối truy cập ngay
		if len(userRoles) == 0 {
			fmt.Printf("[AUTH] ❌ BLOCKING: User has no roles - User: %s, Path: %s\n",
				user.Email, c.Path())
			logrus.WithFields(logrus.Fields{
				"user_id":    user.ID.Hex(),
				"user_email": user.Email,
				"path":       c.Path(),
				"permission": requirePermission,
			}).Error("❌ User has no roles, denying access")
			HandleErrorResponse(c, common.NewError(
				common.ErrCodeAuthRole,
				"Người dùng chưa được gán vai trò. Vui lòng liên hệ quản trị viên để được cấp quyền truy cập.",
				common.StatusForbidden,
				nil,
			))
			return nil
		}

		// Validate user có role này không
		hasRole := false
		for _, userRole := range userRoles {
			if userRole.RoleID == roleID {
				hasRole = true
				break
			}
		}

		// Nếu user không có role này, từ chối truy cập
		if !hasRole {
			logrus.WithFields(logrus.Fields{
				"user_id":        user.ID.Hex(),
				"active_role_id": roleID.Hex(),
				"path":           c.Path(),
			}).Error("❌ User does not have this role")
			HandleErrorResponse(c, common.NewError(
				common.ErrCodeAuthRole,
				"Người dùng không có quyền sử dụng role này. Vui lòng chọn role khác hoặc liên hệ quản trị viên.",
				common.StatusForbidden,
				nil,
			))
			return nil
		}

		activeRoleID := &roleID

		// Kiểm tra permission của user trong role context (active role)
		permissions, err := authManager.getUserPermissions(user.ID.Hex(), activeRoleID)
		if err != nil {
			HandleErrorResponse(c, common.NewError(
				common.ErrCodeAuthRole,
				"Không thể lấy thông tin quyền",
				common.StatusForbidden,
				nil,
			))
			return nil
		}

		// Kiểm tra user có permission cần thiết trong role context không
		scope, hasPermission := permissions[requirePermission]
		if !hasPermission {
			logrus.WithFields(logrus.Fields{
				"user_id":             user.ID.Hex(),
				"user_email":          user.Email,
				"active_role_id":      activeRoleID.Hex(),
				"required_permission": requirePermission,
				"path":                c.Path(),
				"permissions":         permissions,
			}).Error("❌ User does not have required permission")
			HandleErrorResponse(c, common.NewError(
				common.ErrCodeAuthRole,
				"Không có quyền truy cập. Vui lòng kiểm tra lại role context hoặc liên hệ quản trị viên.",
				common.StatusForbidden,
				nil,
			))
			return nil
		}

		logrus.WithFields(logrus.Fields{
			"user_id":        user.ID.Hex(),
			"active_role_id": activeRoleID.Hex(),
			"permission":     requirePermission,
			"scope":          scope,
			"path":           c.Path(),
		}).Info("✅ Permission check passed")

		// Lưu scope tối thiểu vào context để sử dụng trong handler
		c.Locals("minScope", scope)
		return c.Next()
	}
}
