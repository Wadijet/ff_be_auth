package router

import (
	"context"
	"fmt"
	"meta_commerce/core/api/handler"
	"meta_commerce/core/api/middleware"
	"meta_commerce/core/api/services"
	"meta_commerce/core/global"
	"time"

	"github.com/gofiber/fiber/v3"
)

// CONFIGS

// CRUDHandler định nghĩa interface cho các handler CRUD
type CRUDHandler interface {
	// Create
	InsertOne(c fiber.Ctx) error
	InsertMany(c fiber.Ctx) error

	// Read
	Find(c fiber.Ctx) error
	FindOne(c fiber.Ctx) error
	FindOneById(c fiber.Ctx) error
	FindManyByIds(c fiber.Ctx) error
	FindWithPagination(c fiber.Ctx) error

	// Update
	UpdateOne(c fiber.Ctx) error
	UpdateMany(c fiber.Ctx) error
	UpdateById(c fiber.Ctx) error
	FindOneAndUpdate(c fiber.Ctx) error

	// Delete
	DeleteOne(c fiber.Ctx) error
	DeleteMany(c fiber.Ctx) error
	DeleteById(c fiber.Ctx) error
	FindOneAndDelete(c fiber.Ctx) error

	// Other
	CountDocuments(c fiber.Ctx) error
	Distinct(c fiber.Ctx) error
	Upsert(c fiber.Ctx) error
	UpsertMany(c fiber.Ctx) error
	DocumentExists(c fiber.Ctx) error
}

// Router quản lý việc định tuyến cho API
type Router struct {
	app *fiber.App
}

// CRUDConfig cấu hình các operation được phép cho mỗi collection
type CRUDConfig struct {
	// Create
	InsOne  bool // Insert One
	InsMany bool // Insert Many

	// Read
	Find     bool // Find All
	FindOne  bool // Find One
	FindById bool // Find By Id
	FindIds  bool // Find Many By Ids
	Paginate bool // Find With Pagination

	// Update
	UpdOne  bool // Update One
	UpdMany bool // Update Many
	UpdById bool // Update By Id
	FindUpd bool // Find One And Update

	// Delete
	DelOne  bool // Delete One
	DelMany bool // Delete Many
	DelById bool // Delete By Id
	FindDel bool // Find One And Delete

	// Other
	Count    bool // Count Documents
	Distinct bool // Distinct
	Upsert   bool // Upsert One
	UpsMany  bool // Upsert Many
	Exists   bool // Document Exists
}

// Config cho từng collection
var (
	readOnlyConfig = CRUDConfig{
		InsOne: false, InsMany: false,
		Find: true, FindOne: true, FindById: true,
		FindIds: true, Paginate: true,
		UpdOne: false, UpdMany: false, UpdById: false,
		FindUpd: false,
		DelOne:  false, DelMany: false, DelById: false,
		FindDel: false,
		Count:   true, Distinct: true,
		Upsert: false, UpsMany: false, Exists: true,
	}

	readWriteConfig = CRUDConfig{
		InsOne: true, InsMany: true,
		Find: true, FindOne: true, FindById: true,
		FindIds: true, Paginate: true,
		UpdOne: true, UpdMany: true, UpdById: true,
		FindUpd: true,
		DelOne:  true, DelMany: true, DelById: true,
		FindDel: true,
		Count:   true, Distinct: true,
		Upsert: true, UpsMany: true, Exists: true,
	}

	// Auth Module Collections
	userConfig     = readOnlyConfig
	permConfig     = readOnlyConfig
	roleConfig     = readWriteConfig
	rolePermConfig = readWriteConfig
	userRoleConfig = readWriteConfig
	agentConfig    = readWriteConfig

	// Pancake Module Collections
	accessTokenConfig = readWriteConfig
	fbPageConfig      = readWriteConfig
	fbPostConfig      = readWriteConfig
	fbConvConfig      = readWriteConfig
	fbMessageConfig   = readWriteConfig
	pcOrderConfig     = readWriteConfig
)

// RoutePrefix chứa các prefix cơ bản cho API
type RoutePrefix struct {
	Base string // Prefix cơ bản (/api)
	V1   string // Prefix cho API version 1 (/api/v1)
}

// NewRoutePrefix tạo mới một instance của RoutePrefix với các giá trị mặc định
func NewRoutePrefix() RoutePrefix {
	base := "/api"
	return RoutePrefix{
		Base: base,
		V1:   base + "/v1",
	}
}

// NewRouter tạo mới một instance của Router
func NewRouter(app *fiber.App) *Router {
	return &Router{
		app: app,
	}
}

// registerCRUDRoutes đăng ký các route CRUD cho một collection
func (r *Router) registerCRUDRoutes(router fiber.Router, prefix string, h CRUDHandler, config CRUDConfig, permissionPrefix string) {
	// Create operations
	if config.InsOne {
		router.Post(fmt.Sprintf("%s/insert-one", prefix), h.InsertOne, middleware.AuthMiddleware(permissionPrefix+".Insert"))
	}
	if config.InsMany {
		router.Post(fmt.Sprintf("%s/insert-many", prefix), h.InsertMany, middleware.AuthMiddleware(permissionPrefix+".Insert"))
	}

	// Read operations
	if config.Find {
		router.Get(fmt.Sprintf("%s/find", prefix), h.Find, middleware.AuthMiddleware(permissionPrefix+".Read"))
	}
	if config.FindOne {
		router.Get(fmt.Sprintf("%s/find-one", prefix), h.FindOne, middleware.AuthMiddleware(permissionPrefix+".Read"))
	}
	if config.FindById {
		router.Get(fmt.Sprintf("%s/find-by-id/:id", prefix), h.FindOneById, middleware.AuthMiddleware(permissionPrefix+".Read"))
	}
	if config.FindIds {
		router.Post(fmt.Sprintf("%s/find-by-ids", prefix), h.FindManyByIds, middleware.AuthMiddleware(permissionPrefix+".Read"))
	}
	if config.Paginate {
		router.Get(fmt.Sprintf("%s/find-with-pagination", prefix), h.FindWithPagination, middleware.AuthMiddleware(permissionPrefix+".Read"))
	}

	// Update operations
	if config.UpdOne {
		router.Put(fmt.Sprintf("%s/update-one", prefix), h.UpdateOne, middleware.AuthMiddleware(permissionPrefix+".Update"))
	}
	if config.UpdMany {
		router.Put(fmt.Sprintf("%s/update-many", prefix), h.UpdateMany, middleware.AuthMiddleware(permissionPrefix+".Update"))
	}
	if config.UpdById {
		router.Put(fmt.Sprintf("%s/update-by-id/:id", prefix), h.UpdateById, middleware.AuthMiddleware(permissionPrefix+".Update"))
	}
	if config.FindUpd {
		router.Put(fmt.Sprintf("%s/find-one-and-update", prefix), h.FindOneAndUpdate, middleware.AuthMiddleware(permissionPrefix+".Update"))
	}

	// Delete operations
	if config.DelOne {
		router.Delete(fmt.Sprintf("%s/delete-one", prefix), h.DeleteOne, middleware.AuthMiddleware(permissionPrefix+".Delete"))
	}
	if config.DelMany {
		router.Delete(fmt.Sprintf("%s/delete-many", prefix), h.DeleteMany, middleware.AuthMiddleware(permissionPrefix+".Delete"))
	}
	if config.DelById {
		router.Delete(fmt.Sprintf("%s/delete-by-id/:id", prefix), h.DeleteById, middleware.AuthMiddleware(permissionPrefix+".Delete"))
	}
	if config.FindDel {
		router.Delete(fmt.Sprintf("%s/find-one-and-delete", prefix), h.FindOneAndDelete, middleware.AuthMiddleware(permissionPrefix+".Delete"))
	}

	// Other operations
	if config.Count {
		fmt.Printf("Registering COUNT route: %s/count\n", prefix)
		router.Get(fmt.Sprintf("%s/count", prefix), h.CountDocuments, middleware.AuthMiddleware(permissionPrefix+".Read"))
	}
	if config.Distinct {
		router.Get(fmt.Sprintf("%s/distinct", prefix), h.Distinct, middleware.AuthMiddleware(permissionPrefix+".Read"))
	}
	if config.Upsert {
		router.Post(fmt.Sprintf("%s/upsert-one", prefix), h.Upsert, middleware.AuthMiddleware(permissionPrefix+".Update"))
	}
	if config.UpsMany {
		router.Post(fmt.Sprintf("%s/upsert-many", prefix), h.UpsertMany, middleware.AuthMiddleware(permissionPrefix+".Update"))
	}
	if config.Exists {
		router.Get(fmt.Sprintf("%s/exists", prefix), h.DocumentExists, middleware.AuthMiddleware(permissionPrefix+".Read"))
	}
}

// CÁC HÀM ĐĂNG KÝ ROUTES

// registerAdminRoutes đăng ký các route cho admin operations
func registerAdminRoutes(router fiber.Router) error {
	// Admin routes
	adminHandler, err := handler.NewAdminHandler()
	if err != nil {
		return fmt.Errorf("failed to create admin handler: %v", err)
	}

	// Các route đặc biệt cho quản trị viên
	router.Post("/admin/user/block", middleware.AuthMiddleware("User.Block"), adminHandler.HandleBlockUser)
	router.Post("/admin/user/unblock", middleware.AuthMiddleware("User.Block"), adminHandler.HandleUnBlockUser)
	router.Post("/admin/user/role", middleware.AuthMiddleware("User.SetRole"), adminHandler.HandleSetRole)
	// Thiết lập administrator (yêu cầu quyền Init.SetAdmin)
	router.Post("/admin/user/set-administrator/:id", middleware.AuthMiddleware("Init.SetAdmin"), adminHandler.HandleAddAdministrator)

	return nil
}

// registerSystemRoutes đăng ký các route cho system operations
func registerSystemRoutes(router fiber.Router) error {
	// System routes
	router.Get("/system/health", func(c fiber.Ctx) error {
		// Kiểm tra database connection
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		healthStatus := fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"services": fiber.Map{
				"api": "ok",
			},
		}

		// Kiểm tra MongoDB connection
		if global.MongoDB_Session != nil {
			err := global.MongoDB_Session.Ping(ctx, nil)
			if err != nil {
				healthStatus["status"] = "degraded"
				healthStatus["services"].(fiber.Map)["database"] = "error"
				healthStatus["database_error"] = err.Error()
				return c.Status(fiber.StatusServiceUnavailable).JSON(healthStatus)
			}
			healthStatus["services"].(fiber.Map)["database"] = "ok"
		} else {
			healthStatus["status"] = "degraded"
			healthStatus["services"].(fiber.Map)["database"] = "not_initialized"
		}

		return c.JSON(healthStatus)
	})

	return nil
}

// registerAuthRoutes đăng ký các route cho authentication cá nhân
func (r *Router) registerAuthRoutes(router fiber.Router) error {
	// User routes
	userHandler, err := handler.NewUserHandler()
	if err != nil {
		return fmt.Errorf("failed to create user handler: %v", err)
	}

	// Các route xác thực cá nhân
	// Firebase Authentication - Nhận Firebase ID token và tạo JWT
	router.Post("/auth/login/firebase", userHandler.HandleLoginWithFirebase)

	// Logout - Xóa JWT token
	router.Post("/auth/logout", userHandler.HandleLogout, middleware.AuthMiddleware(""))

	// Profile - Lấy và cập nhật thông tin user
	router.Get("/auth/profile", userHandler.HandleGetProfile, middleware.AuthMiddleware(""))
	router.Put("/auth/profile", userHandler.HandleUpdateProfile, middleware.AuthMiddleware(""))

	// Roles - Lấy danh sách roles của user
	router.Get("/auth/roles", userHandler.HandleGetUserRoles, middleware.AuthMiddleware(""))

	return nil
}

// registerRBACRoutes đăng ký các route cho Role-Based Access Control
func (r *Router) registerRBACRoutes(router fiber.Router) error {
	// User routes (Quản lý người dùng)
	userHandler, err := handler.NewUserHandler()
	if err != nil {
		return fmt.Errorf("failed to create user handler: %v", err)
	}
	r.registerCRUDRoutes(router, "/user", userHandler, userConfig, "User")

	// Permission routes
	permHandler, err := handler.NewPermissionHandler()
	if err != nil {
		return fmt.Errorf("failed to create permission handler: %v", err)
	}
	fmt.Printf("Registering permission routes with prefix: /permission\n")
	r.registerCRUDRoutes(router, "/permission", permHandler, permConfig, "Permission")

	// Role routes
	roleHandler, err := handler.NewRoleHandler()
	if err != nil {
		return fmt.Errorf("failed to create role handler: %v", err)
	}
	r.registerCRUDRoutes(router, "/role", roleHandler, roleConfig, "Role")

	// RolePermission routes
	rolePermHandler, err := handler.NewRolePermissionHandler()
	if err != nil {
		return fmt.Errorf("failed to create role permission handler: %v", err)
	}
	// Route đặc biệt cho cập nhật quyền của vai trò
	router.Put("/role-permission/update-role", middleware.AuthMiddleware("RolePermission.Update"), rolePermHandler.HandleUpdateRolePermissions)
	// CRUD routes
	r.registerCRUDRoutes(router, "/role-permission", rolePermHandler, rolePermConfig, "RolePermission")

	// UserRole routes
	userRoleHandler, err := handler.NewUserRoleHandler()
	if err != nil {
		return fmt.Errorf("failed to create user role handler: %v", err)
	}
	r.registerCRUDRoutes(router, "/user-role", userRoleHandler, userRoleConfig, "UserRole")

	// Organization routes
	organizationHandler, err := handler.NewOrganizationHandler()
	if err != nil {
		return fmt.Errorf("failed to create organization handler: %v", err)
	}
	fmt.Printf("Registering organization routes with prefix: /organization\n")
	r.registerCRUDRoutes(router, "/organization", organizationHandler, readWriteConfig, "Organization")
	fmt.Printf("Organization routes registered successfully\n")

	// Agent routes
	agentHandler, err := handler.NewAgentHandler()
	if err != nil {
		return fmt.Errorf("failed to create agent handler: %v", err)
	}
	// Đăng ký các route đặc biệt cho agent: check-in/check-out
	router.Post("/agent/check-in/:id", middleware.AuthMiddleware("Agent.CheckIn"), agentHandler.HandleCheckIn)    // Route check-in cho agent
	router.Post("/agent/check-out/:id", middleware.AuthMiddleware("Agent.CheckOut"), agentHandler.HandleCheckOut) // Route check-out cho agent
	r.registerCRUDRoutes(router, "/agent", agentHandler, agentConfig, "Agent")

	return nil
}

// registerFacebookRoutes đăng ký các route cho Facebook integration
func (r *Router) registerFacebookRoutes(router fiber.Router) error {
	// Access Token routes
	accessTokenHandler, err := handler.NewAccessTokenHandler()
	if err != nil {
		return fmt.Errorf("failed to create access token handler: %v", err)
	}
	r.registerCRUDRoutes(router, "/access-token", accessTokenHandler, accessTokenConfig, "AccessToken")

	// Facebook Page routes
	fbPageHandler, err := handler.NewFbPageHandler()
	if err != nil {
		return fmt.Errorf("failed to create facebook page handler: %v", err)
	}
	r.registerCRUDRoutes(router, "/facebook/page", fbPageHandler, fbPageConfig, "FbPage")

	// Facebook Post routes
	fbPostHandler, err := handler.NewFbPostHandler()
	if err != nil {
		return fmt.Errorf("failed to create facebook post handler: %v", err)
	}
	r.registerCRUDRoutes(router, "/facebook/post", fbPostHandler, fbPostConfig, "FbPost")

	// Facebook Conversation routes
	fbConvHandler, err := handler.NewFbConversationHandler()
	if err != nil {
		return fmt.Errorf("failed to create facebook conversation handler: %v", err)
	}
	// Route đặc biệt cho lấy cuộc trò chuyện sắp xếp theo thời gian cập nhật API
	router.Get("/facebook/conversation/sort-by-api-update", middleware.AuthMiddleware("FbConversation.Read"), fbConvHandler.HandleFindAllSortByApiUpdate)
	// CRUD routes
	r.registerCRUDRoutes(router, "/facebook/conversation", fbConvHandler, fbConvConfig, "FbConversation")

	// Facebook Message routes
	fbMessageHandler, err := handler.NewFbMessageHandler()
	if err != nil {
		return fmt.Errorf("failed to create facebook message handler: %v", err)
	}
	r.registerCRUDRoutes(router, "/facebook/message", fbMessageHandler, fbMessageConfig, "FbMessage")

	// Pancake Order routes
	pcOrderHandler, err := handler.NewPcOrderHandler()
	if err != nil {
		return fmt.Errorf("failed to create pancake order handler: %v", err)
	}
	r.registerCRUDRoutes(router, "/pancake/order", pcOrderHandler, pcOrderConfig, "PcOrder")

	return nil
}

// registerInitRoutes đăng ký các route cho khởi tạo hệ thống
func (r *Router) registerInitRoutes(router fiber.Router) error {
	// Kiểm tra xem đã có admin chưa
	// Nếu đã có admin, không đăng ký bất kỳ init endpoint nào (tối ưu hiệu suất và bảo mật)
	initService, err := services.NewInitService()
	if err == nil {
		hasAdmin, err := initService.HasAnyAdministrator()
		if err == nil && hasAdmin {
			// Đã có admin, không đăng ký bất kỳ init endpoint nào
			// Endpoint thêm admin sẽ ở /admin/user/set-administrator/:id
			return nil
		}
	}

	// Chưa có admin, đăng ký tất cả init endpoints
	initHandler, err := handler.NewInitHandler()
	if err != nil {
		return fmt.Errorf("failed to create init handler: %v", err)
	}

	// Route kiểm tra trạng thái init (chỉ khi chưa có admin)
	router.Get("/init/status", initHandler.HandleInitStatus)

	// Các route khởi tạo các đơn vị cơ bản
	router.Post("/init/organization", initHandler.HandleInitOrganization)
	router.Post("/init/permissions", initHandler.HandleInitPermissions)
	router.Post("/init/roles", initHandler.HandleInitRoles)
	router.Post("/init/admin-user", initHandler.HandleInitAdminUser)
	router.Post("/init/all", initHandler.HandleInitAll) // One-click setup

	// Route thiết lập administrator lần đầu (chưa có admin, không cần quyền)
	// Handler sẽ tự check xem đã có admin chưa
	router.Post("/init/set-administrator/:id", initHandler.HandleSetAdministrator)

	return nil
}

// SetupRoutes thiết lập tất cả các route cho ứng dụng
func SetupRoutes(app *fiber.App) error {
	// Khởi tạo route prefix
	prefix := NewRoutePrefix()
	v1 := app.Group(prefix.V1)

	// Khởi tạo router
	router := NewRouter(app)

	// 1. Init Routes
	if err := router.registerInitRoutes(v1); err != nil {
		return fmt.Errorf("failed to register init routes: %v", err)
	}

	// 2. Admin Routes
	if err := registerAdminRoutes(v1); err != nil {
		return fmt.Errorf("failed to register admin routes: %v", err)
	}

	// 3. System Routes
	if err := registerSystemRoutes(v1); err != nil {
		return fmt.Errorf("failed to register system routes: %v", err)
	}

	// 4. Auth Routes (Xác thực cá nhân)
	if err := router.registerAuthRoutes(v1); err != nil {
		return fmt.Errorf("failed to register auth routes: %v", err)
	}

	// 5. RBAC Routes (Bao gồm User Management)
	if err := router.registerRBACRoutes(v1); err != nil {
		return fmt.Errorf("failed to register RBAC routes: %v", err)
	}

	// 6. Facebook Routes
	if err := router.registerFacebookRoutes(v1); err != nil {
		return fmt.Errorf("failed to register Facebook routes: %v", err)
	}

	return nil
}
