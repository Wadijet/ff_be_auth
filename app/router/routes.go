package router

import (
	"meta_commerce/app/handler"
	"meta_commerce/app/middleware"

	"github.com/gofiber/fiber/v3"
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

// registerFiberCRUDRoutes đăng ký các route CRUD cơ bản cho một handler
func registerFiberCRUDRoutes(router fiber.Router, prefix string, handler interface{}, permission string) {
	// Tạo group với prefix
	group := router.Group(prefix)

	// Create operations
	if h, ok := handler.(interface{ InsertOne(c fiber.Ctx) error }); ok {
		group.Post("/", middleware.FiberAuthMiddleware(permission+".Create"), h.InsertOne)
	}
	if h, ok := handler.(interface{ InsertMany(c fiber.Ctx) error }); ok {
		group.Post("/batch", middleware.FiberAuthMiddleware(permission+".Create"), h.InsertMany)
	}

	// Read operations
	if h, ok := handler.(interface{ FindOne(c fiber.Ctx) error }); ok {
		group.Get("/find-one", middleware.FiberAuthMiddleware(permission+".Read"), h.FindOne)
	}
	if h, ok := handler.(interface{ FindOneById(c fiber.Ctx) error }); ok {
		group.Get("/:id", middleware.FiberAuthMiddleware(permission+".Read"), h.FindOneById)
	}
	if h, ok := handler.(interface{ FindManyByIds(c fiber.Ctx) error }); ok {
		group.Get("/by-ids", middleware.FiberAuthMiddleware(permission+".Read"), h.FindManyByIds)
	}
	if h, ok := handler.(interface{ FindWithPagination(c fiber.Ctx) error }); ok {
		group.Get("/paginate", middleware.FiberAuthMiddleware(permission+".Read"), h.FindWithPagination)
	}
	if h, ok := handler.(interface{ Find(c fiber.Ctx) error }); ok {
		group.Get("/", middleware.FiberAuthMiddleware(permission+".Read"), h.Find)
	}

	// Update operations
	if h, ok := handler.(interface{ UpdateOne(c fiber.Ctx) error }); ok {
		group.Put("/update-one", middleware.FiberAuthMiddleware(permission+".Update"), h.UpdateOne)
	}
	if h, ok := handler.(interface{ UpdateMany(c fiber.Ctx) error }); ok {
		group.Put("/batch", middleware.FiberAuthMiddleware(permission+".Update"), h.UpdateMany)
	}
	if h, ok := handler.(interface{ UpdateById(c fiber.Ctx) error }); ok {
		group.Put("/:id", middleware.FiberAuthMiddleware(permission+".Update"), h.UpdateById)
	}
	if h, ok := handler.(interface{ FindOneAndUpdate(c fiber.Ctx) error }); ok {
		group.Put("/find-and-update", middleware.FiberAuthMiddleware(permission+".Update"), h.FindOneAndUpdate)
	}

	// Delete operations
	if h, ok := handler.(interface{ DeleteOne(c fiber.Ctx) error }); ok {
		group.Delete("/delete-one", middleware.FiberAuthMiddleware(permission+".Delete"), h.DeleteOne)
	}
	if h, ok := handler.(interface{ DeleteMany(c fiber.Ctx) error }); ok {
		group.Delete("/batch", middleware.FiberAuthMiddleware(permission+".Delete"), h.DeleteMany)
	}
	if h, ok := handler.(interface{ DeleteById(c fiber.Ctx) error }); ok {
		group.Delete("/:id", middleware.FiberAuthMiddleware(permission+".Delete"), h.DeleteById)
	}
	if h, ok := handler.(interface{ FindOneAndDelete(c fiber.Ctx) error }); ok {
		group.Delete("/find-and-delete", middleware.FiberAuthMiddleware(permission+".Delete"), h.FindOneAndDelete)
	}

	// Other operations
	if h, ok := handler.(interface{ CountDocuments(c fiber.Ctx) error }); ok {
		group.Get("/count", middleware.FiberAuthMiddleware(permission+".Read"), h.CountDocuments)
	}
	if h, ok := handler.(interface{ Distinct(c fiber.Ctx) error }); ok {
		group.Get("/distinct/:field", middleware.FiberAuthMiddleware(permission+".Read"), h.Distinct)
	}
	if h, ok := handler.(interface{ Upsert(c fiber.Ctx) error }); ok {
		group.Post("/upsert", middleware.FiberAuthMiddleware(permission+".Create"), h.Upsert)
	}
	if h, ok := handler.(interface{ UpsertMany(c fiber.Ctx) error }); ok {
		group.Post("/upsert/batch", middleware.FiberAuthMiddleware(permission+".Create"), h.UpsertMany)
	}
	if h, ok := handler.(interface{ DocumentExists(c fiber.Ctx) error }); ok {
		group.Get("/exists", middleware.FiberAuthMiddleware(permission+".Read"), h.DocumentExists)
	}
}

// registerFiberAuthRoutes đăng ký các route cho authentication
func registerFiberAuthRoutes(router fiber.Router) {
	// Khởi tạo các handler
	permissionHandler := handler.NewFiberPermissionHandler()
	roleHandler := handler.NewFiberRoleHandler()
	rolePermissionHandler := handler.NewFiberRolePermissionHandler()
	userRoleHandler := handler.NewFiberUserRoleHandler()
	userHandler := handler.NewFiberAuthUserHandler()

	// Đăng ký route cho từng handler
	registerFiberCRUDRoutes(router, "/permission", permissionHandler, "Permission")
	registerFiberCRUDRoutes(router, "/role", roleHandler, "Role")
	registerFiberCRUDRoutes(router, "/role-permission", rolePermissionHandler, "RolePermission")
	registerFiberCRUDRoutes(router, "/user-role", userRoleHandler, "UserRole")
	registerFiberCRUDRoutes(router, "/user", userHandler, "User")

	// Route đặc biệt cho user
	userGroup := router.Group("/user")
	userGroup.Post("/login", userHandler.HandleLogin)
	userGroup.Post("/register", userHandler.HandleRegister)
	userGroup.Post("/logout", middleware.FiberAuthMiddleware("User.Logout"), userHandler.HandleLogout)
	userGroup.Get("/profile", middleware.FiberAuthMiddleware("User.Profile"), userHandler.HandleGetMyInfo)
	userGroup.Put("/profile", middleware.FiberAuthMiddleware("User.Profile"), userHandler.HandleChangeInfo)
	userGroup.Put("/change-password", middleware.FiberAuthMiddleware("User.ChangePassword"), userHandler.HandleChangePassword)
}

// SetupFiberRoutes thiết lập tất cả các route cho ứng dụng
func SetupFiberRoutes(app *fiber.App) {
	// Khởi tạo route prefix
	prefix := NewRoutePrefix()

	// Khởi tạo các handler
	staticHandler := handler.NewFiberStaticHandler()
	fbConversationHandler := handler.NewFiberFbConversationHandler()
	fbPostHandler := handler.NewFiberFbPostHandler()
	fbPageHandler := handler.NewFiberFbPageHandler()
	adminInitHandler := handler.NewFiberInitHandler()

	// API v1 routes
	v1 := app.Group(prefix.V1)

	// Static routes
	v1.Get("/static/test", staticHandler.HandleTestApi)
	v1.Get("/static/system", staticHandler.HandleGetSystemStatic)
	v1.Get("/static/api", staticHandler.HandleGetApiStatic)

	// Facebook routes
	registerFiberCRUDRoutes(v1, "/conversation", fbConversationHandler, "Facebook")
	registerFiberCRUDRoutes(v1, "/post", fbPostHandler, "Facebook")
	registerFiberCRUDRoutes(v1, "/page", fbPageHandler, "Facebook")

	// Route đặc thù cho Facebook
	v1.Get("/conversation/newest", middleware.FiberAuthMiddleware("Facebook.Read"), fbConversationHandler.HandleFindAllSortByApiUpdate)
	v1.Get("/post/:postId", middleware.FiberAuthMiddleware("Facebook.Read"), fbPostHandler.HandleFindOneByPostID)
	v1.Put("/post/token", middleware.FiberAuthMiddleware("Facebook.Update"), fbPostHandler.HandleUpdateToken)
	v1.Get("/page/:id", middleware.FiberAuthMiddleware("Facebook.Read"), fbPageHandler.HandleFindOneByPageID)
	v1.Put("/page/token", middleware.FiberAuthMiddleware("Facebook.Update"), fbPageHandler.HandleUpdateToken)

	// Admin routes
	v1.Post("/admin/init", adminInitHandler.HandleSetAdministrator)

	// Auth routes
	registerFiberAuthRoutes(v1)
}
