package router

import (
	"fmt"
	"meta_commerce/core/api/handler"
	"meta_commerce/core/api/middleware"

	"github.com/gofiber/fiber/v3"
)

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
func (r *Router) registerCRUDRoutes(prefix string, h CRUDHandler, config CRUDConfig, permissionPrefix string) {
	// Create operations
	if config.InsOne {
		r.app.Post(fmt.Sprintf("%s/insertOne", prefix), middleware.AuthMiddleware(permissionPrefix+".Insert"), h.InsertOne)
	}
	if config.InsMany {
		r.app.Post(fmt.Sprintf("%s/insertMany", prefix), middleware.AuthMiddleware(permissionPrefix+".Insert"), h.InsertMany)
	}

	// Read operations
	if config.Find {
		r.app.Get(fmt.Sprintf("%s/find", prefix), middleware.AuthMiddleware(permissionPrefix+".Read"), h.Find)
	}
	if config.FindOne {
		r.app.Get(fmt.Sprintf("%s/findOne", prefix), middleware.AuthMiddleware(permissionPrefix+".Read"), h.FindOne)
	}
	if config.FindById {
		r.app.Get(fmt.Sprintf("%s/findById/:id", prefix), middleware.AuthMiddleware(permissionPrefix+".Read"), h.FindOneById)
	}
	if config.FindIds {
		r.app.Post(fmt.Sprintf("%s/findByIds", prefix), middleware.AuthMiddleware(permissionPrefix+".Read"), h.FindManyByIds)
	}
	if config.Paginate {
		r.app.Get(fmt.Sprintf("%s/findWithPagination", prefix), middleware.AuthMiddleware(permissionPrefix+".Read"), h.FindWithPagination)
	}

	// Update operations
	if config.UpdOne {
		r.app.Put(fmt.Sprintf("%s/updateOne", prefix), middleware.AuthMiddleware(permissionPrefix+".Update"), h.UpdateOne)
	}
	if config.UpdMany {
		r.app.Put(fmt.Sprintf("%s/updateMany", prefix), middleware.AuthMiddleware(permissionPrefix+".Update"), h.UpdateMany)
	}
	if config.UpdById {
		r.app.Put(fmt.Sprintf("%s/updateById/:id", prefix), middleware.AuthMiddleware(permissionPrefix+".Update"), h.UpdateById)
	}
	if config.FindUpd {
		r.app.Put(fmt.Sprintf("%s/findOneAndUpdate", prefix), middleware.AuthMiddleware(permissionPrefix+".Update"), h.FindOneAndUpdate)
	}

	// Delete operations
	if config.DelOne {
		r.app.Delete(fmt.Sprintf("%s/deleteOne", prefix), middleware.AuthMiddleware(permissionPrefix+".Delete"), h.DeleteOne)
	}
	if config.DelMany {
		r.app.Delete(fmt.Sprintf("%s/deleteMany", prefix), middleware.AuthMiddleware(permissionPrefix+".Delete"), h.DeleteMany)
	}
	if config.DelById {
		r.app.Delete(fmt.Sprintf("%s/deleteById/:id", prefix), middleware.AuthMiddleware(permissionPrefix+".Delete"), h.DeleteById)
	}
	if config.FindDel {
		r.app.Delete(fmt.Sprintf("%s/findOneAndDelete", prefix), middleware.AuthMiddleware(permissionPrefix+".Delete"), h.FindOneAndDelete)
	}

	// Other operations
	if config.Count {
		r.app.Get(fmt.Sprintf("%s/count", prefix), middleware.AuthMiddleware(permissionPrefix+".Read"), h.CountDocuments)
	}
	if config.Distinct {
		r.app.Get(fmt.Sprintf("%s/distinct", prefix), middleware.AuthMiddleware(permissionPrefix+".Read"), h.Distinct)
	}
	if config.Upsert {
		r.app.Post(fmt.Sprintf("%s/upsertOne", prefix), middleware.AuthMiddleware(permissionPrefix+".Update"), h.Upsert)
	}
	if config.UpsMany {
		r.app.Post(fmt.Sprintf("%s/upsertMany", prefix), middleware.AuthMiddleware(permissionPrefix+".Update"), h.UpsertMany)
	}
	if config.Exists {
		r.app.Get(fmt.Sprintf("%s/exists", prefix), middleware.AuthMiddleware(permissionPrefix+".Read"), h.DocumentExists)
	}
}

// registerAuthRoutes đăng ký các route cho authentication
func (r *Router) registerAuthRoutes(router fiber.Router) error {
	// User routes
	userHandler, err := handler.NewUserHandler()
	if err != nil {
		return fmt.Errorf("failed to create user handler: %v", err)
	}

	// Các route xác thực đặc biệt
	router.Post("/auth/register", userHandler.HandleRegister)
	router.Post("/auth/login", userHandler.HandleLogin)
	router.Post("/auth/logout", middleware.AuthMiddleware(""), userHandler.HandleLogout)
	router.Get("/auth/profile", middleware.AuthMiddleware(""), userHandler.HandleGetProfile)
	router.Put("/auth/profile", middleware.AuthMiddleware(""), userHandler.HandleUpdateProfile)
	router.Put("/auth/password", middleware.AuthMiddleware(""), userHandler.HandleChangePassword)

	// CRUD routes cho User với quyền từ InitialPermissions
	r.registerCRUDRoutes("/users", userHandler, userConfig, "User")

	return nil
}

// registerRBACRoutes đăng ký các route cho Role-Based Access Control
func (r *Router) registerRBACRoutes(router fiber.Router) error {
	// Permission routes
	permHandler, err := handler.NewPermissionHandler()
	if err != nil {
		return fmt.Errorf("failed to create permission handler: %v", err)
	}
	r.registerCRUDRoutes("/permissions", permHandler, permConfig, "Permission")

	// Role routes
	roleHandler, err := handler.NewRoleHandler()
	if err != nil {
		return fmt.Errorf("failed to create role handler: %v", err)
	}
	r.registerCRUDRoutes("/roles", roleHandler, roleConfig, "Role")

	// RolePermission routes
	rolePermHandler, err := handler.NewRolePermissionHandler()
	if err != nil {
		return fmt.Errorf("failed to create role permission handler: %v", err)
	}
	// Route đặc biệt cho cập nhật quyền của vai trò
	router.Put("/role-permissions/update-role", middleware.AuthMiddleware("RolePermission.Update"), rolePermHandler.HandleUpdateRolePermissions)
	// CRUD routes
	r.registerCRUDRoutes("/role-permissions", rolePermHandler, rolePermConfig, "RolePermission")

	// UserRole routes
	userRoleHandler, err := handler.NewUserRoleHandler()
	if err != nil {
		return fmt.Errorf("failed to create user role handler: %v", err)
	}
	r.registerCRUDRoutes("/user-roles", userRoleHandler, userRoleConfig, "UserRole")

	// Agent routes
	agentHandler, err := handler.NewAgentHandler()
	if err != nil {
		return fmt.Errorf("failed to create agent handler: %v", err)
	}
	r.registerCRUDRoutes("/agents", agentHandler, agentConfig, "Agent")

	return nil
}

// registerFacebookRoutes đăng ký các route cho Facebook integration
func (r *Router) registerFacebookRoutes(router fiber.Router) error {
	// Access Token routes
	accessTokenHandler, err := handler.NewAccessTokenHandler()
	if err != nil {
		return fmt.Errorf("failed to create access token handler: %v", err)
	}
	r.registerCRUDRoutes("/access-tokens", accessTokenHandler, accessTokenConfig, "AccessToken")

	// Facebook Page routes
	fbPageHandler, err := handler.NewFbPageHandler()
	if err != nil {
		return fmt.Errorf("failed to create facebook page handler: %v", err)
	}
	r.registerCRUDRoutes("/facebook/pages", fbPageHandler, fbPageConfig, "FbPage")

	// Facebook Post routes
	fbPostHandler, err := handler.NewFbPostHandler()
	if err != nil {
		return fmt.Errorf("failed to create facebook post handler: %v", err)
	}
	r.registerCRUDRoutes("/facebook/posts", fbPostHandler, fbPostConfig, "FbPost")

	// Facebook Conversation routes
	fbConvHandler, err := handler.NewFbConversationHandler()
	if err != nil {
		return fmt.Errorf("failed to create facebook conversation handler: %v", err)
	}
	// Route đặc biệt cho lấy cuộc trò chuyện sắp xếp theo thời gian cập nhật API
	router.Get("/facebook/conversations/sort-by-api-update", middleware.AuthMiddleware("FbConversation.Read"), fbConvHandler.HandleFindAllSortByApiUpdate)
	// CRUD routes
	r.registerCRUDRoutes("/facebook/conversations", fbConvHandler, fbConvConfig, "FbConversation")

	// Facebook Message routes
	fbMessageHandler, err := handler.NewFbMessageHandler()
	if err != nil {
		return fmt.Errorf("failed to create facebook message handler: %v", err)
	}
	r.registerCRUDRoutes("/facebook/messages", fbMessageHandler, fbMessageConfig, "FbMessage")

	// Pancake Order routes
	pcOrderHandler, err := handler.NewPcOrderHandler()
	if err != nil {
		return fmt.Errorf("failed to create pancake order handler: %v", err)
	}
	r.registerCRUDRoutes("/pancake/orders", pcOrderHandler, pcOrderConfig, "PcOrder")

	return nil
}

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

	return nil
}

// registerSystemRoutes đăng ký các route cho system operations
func registerSystemRoutes(router fiber.Router) error {
	// System routes
	router.Get("/system/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})
	router.Get("/system/version", func(c fiber.Ctx) error {
		return c.SendString("v1.0.0")
	})

	return nil
}

// SetupRoutes thiết lập tất cả các route cho ứng dụng
func SetupRoutes(app *fiber.App) error {
	// Khởi tạo route prefix
	prefix := NewRoutePrefix()
	v1 := app.Group(prefix.V1)

	// Khởi tạo router
	router := NewRouter(app)

	// 1. Admin Routes
	if err := registerAdminRoutes(v1); err != nil {
		return fmt.Errorf("failed to register admin routes: %v", err)
	}

	// 2. System Routes
	if err := registerSystemRoutes(v1); err != nil {
		return fmt.Errorf("failed to register system routes: %v", err)
	}

	// 3. Auth Routes
	if err := router.registerAuthRoutes(v1); err != nil {
		return fmt.Errorf("failed to register auth routes: %v", err)
	}

	// 4. RBAC Routes
	if err := router.registerRBACRoutes(v1); err != nil {
		return fmt.Errorf("failed to register RBAC routes: %v", err)
	}

	// 5. Facebook Routes
	if err := router.registerFacebookRoutes(v1); err != nil {
		return fmt.Errorf("failed to register Facebook routes: %v", err)
	}

	return nil
}
