package router

import (
	"fmt"
	"meta_commerce/core/api/handler"
	"meta_commerce/core/api/middleware"

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

// registerCRUDRoutes đăng ký các route CRUD cơ bản cho một handler
func registerCRUDRoutes(router fiber.Router, prefix string, handler interface{}, permission string) {
	// Tạo group với prefix
	group := router.Group(prefix)

	// Create operations
	if h, ok := handler.(interface{ InsertOne(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Insert")).Post("/", h.InsertOne)
	}
	if h, ok := handler.(interface{ InsertMany(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Insert")).Post("/batch", h.InsertMany)
	}

	// Read operations
	if h, ok := handler.(interface{ FindOne(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Read")).Get("/find-one", h.FindOne)
	}
	if h, ok := handler.(interface{ FindOneById(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Read")).Get("/:id", h.FindOneById)
	}
	if h, ok := handler.(interface{ FindManyByIds(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Read")).Get("/by-ids", h.FindManyByIds)
	}
	if h, ok := handler.(interface{ FindWithPagination(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Read")).Get("/paginate", h.FindWithPagination)
	}
	if h, ok := handler.(interface{ Find(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Read")).Get("/", h.Find)
	}

	// Update operations
	if h, ok := handler.(interface{ UpdateOne(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Update")).Put("/update-one", h.UpdateOne)
	}
	if h, ok := handler.(interface{ UpdateMany(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Update")).Put("/batch", h.UpdateMany)
	}
	if h, ok := handler.(interface{ UpdateById(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Update")).Put("/:id", h.UpdateById)
	}
	if h, ok := handler.(interface{ FindOneAndUpdate(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Update")).Put("/find-and-update", h.FindOneAndUpdate)
	}

	// Delete operations
	if h, ok := handler.(interface{ DeleteOne(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Delete")).Delete("/delete-one", h.DeleteOne)
	}
	if h, ok := handler.(interface{ DeleteMany(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Delete")).Delete("/batch", h.DeleteMany)
	}
	if h, ok := handler.(interface{ DeleteById(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Delete")).Delete("/:id", h.DeleteById)
	}
	if h, ok := handler.(interface{ FindOneAndDelete(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Delete")).Delete("/find-and-delete", h.FindOneAndDelete)
	}

	// Other operations
	if h, ok := handler.(interface{ CountDocuments(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Read")).Get("/count", h.CountDocuments)
	}
	if h, ok := handler.(interface{ Distinct(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Read")).Get("/distinct/:field", h.Distinct)
	}
	if h, ok := handler.(interface{ Upsert(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Insert")).Post("/upsert", h.Upsert)
	}
	if h, ok := handler.(interface{ UpsertMany(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Insert")).Post("/upsert/batch", h.UpsertMany)
	}
	if h, ok := handler.(interface{ DocumentExists(c fiber.Ctx) error }); ok {
		group.Use(middleware.AuthMiddleware(permission+".Read")).Get("/exists", h.DocumentExists)
	}
}

// registerAuthRoutes đăng ký các route cho authentication
func registerAuthRoutes(router fiber.Router) error {
	// Khởi tạo các handler
	permissionHandler, err := handler.NewPermissionHandler()
	if err != nil {
		return fmt.Errorf("failed to create permission handler: %v", err)
	}

	roleHandler, err := handler.NewRoleHandler()
	if err != nil {
		return fmt.Errorf("failed to create role handler: %v", err)
	}

	rolePermissionHandler, err := handler.NewRolePermissionHandler()
	if err != nil {
		return fmt.Errorf("failed to create role permission handler: %v", err)
	}

	userRoleHandler, err := handler.NewUserRoleHandler()
	if err != nil {
		return fmt.Errorf("failed to create user role handler: %v", err)
	}

	userHandler, err := handler.NewUserHandler()
	if err != nil {
		return fmt.Errorf("failed to create user handler: %v", err)
	}

	// Route đặc biệt cho user
	userGroup := router.Group("/users")
	// Route không cần xác thực
	userGroup.Post("/login", userHandler.HandleLogin)
	userGroup.Post("/register", userHandler.HandleRegister)
	// Route cần xác thực
	userGroup.Use(middleware.AuthMiddleware("")).Post("/logout", userHandler.HandleLogout)
	userGroup.Use(middleware.AuthMiddleware("")).Get("/profile", userHandler.HandleGetProfile)
	userGroup.Use(middleware.AuthMiddleware("")).Put("/profile", userHandler.HandleUpdateProfile)
	userGroup.Use(middleware.AuthMiddleware("")).Put("/change-password", userHandler.HandleChangePassword)

	// Tạo group riêng cho các route CRUD
	crudGroup := router.Group("/crud")
	// Đăng ký route CRUD cho từng handler
	registerCRUDRoutes(crudGroup, "/permissions", permissionHandler, "Permission")
	registerCRUDRoutes(crudGroup, "/roles", roleHandler, "Role")
	registerCRUDRoutes(crudGroup, "/role-permissions", rolePermissionHandler, "RolePermission")
	registerCRUDRoutes(crudGroup, "/user-roles", userRoleHandler, "UserRole")
	registerCRUDRoutes(crudGroup, "/users", userHandler, "User")

	return nil
}

// SetupRoutes thiết lập tất cả các route cho ứng dụng
func SetupRoutes(app *fiber.App) error {
	// Khởi tạo route prefix
	prefix := NewRoutePrefix()

	// Khởi tạo các handler
	staticHandler := handler.NewStaticHandler()

	fbConversationHandler, err := handler.NewFbConversationHandler()
	if err != nil {
		return fmt.Errorf("failed to create conversation handler: %v", err)
	}

	fbPostHandler, err := handler.NewFbPostHandler()
	if err != nil {
		return fmt.Errorf("failed to create post handler: %v", err)
	}

	fbPageHandler, err := handler.NewFbPageHandler()
	if err != nil {
		return fmt.Errorf("failed to create page handler: %v", err)
	}

	adminInitHandler, err := handler.NewInitHandler()
	if err != nil {
		return fmt.Errorf("failed to create init handler: %v", err)
	}

	// API v1 routes
	v1 := app.Group(prefix.V1)

	// Admin routes
	v1.Post("/admin/init", adminInitHandler.HandleSetAdministrator)

	// Auth routes
	if err := registerAuthRoutes(v1); err != nil {
		return fmt.Errorf("failed to register auth routes: %v", err)
	}

	// Static routes
	v1.Get("/static/test", staticHandler.HandleTestApi)

	// Facebook routes
	registerCRUDRoutes(v1, "/conversations", fbConversationHandler, "Facebook")
	registerCRUDRoutes(v1, "/posts", fbPostHandler, "Facebook")
	registerCRUDRoutes(v1, "/pages", fbPageHandler, "Facebook")

	// Route đặc thù cho Facebook
	v1.Use(middleware.AuthMiddleware("Facebook.Read")).Get("/conversations/newest", fbConversationHandler.HandleFindAllSortByApiUpdate)
	v1.Use(middleware.AuthMiddleware("Facebook.Read")).Get("/posts/:postId", fbPostHandler.HandleFindOneByPostID)
	v1.Use(middleware.AuthMiddleware("Facebook.Update")).Put("/posts/token", fbPostHandler.HandleUpdateToken)
	v1.Use(middleware.AuthMiddleware("Facebook.Read")).Get("/pages/:id", fbPageHandler.HandleFindOneByPageID)
	v1.Use(middleware.AuthMiddleware("Facebook.Update")).Put("/pages/token", fbPageHandler.HandleUpdateToken)

	return nil
}
