package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/registry"
	"meta_commerce/app/router"
	"meta_commerce/app/services"
	"meta_commerce/config"
	"meta_commerce/database"
)

// Hàm khởi tạo các biến toàn cục
func initGlobal() {
	initColNames()         // Khởi tạo tên các collection trong database
	initValidator()        // Khởi tạo validator
	initConfig()           // Khởi tạo cấu hình server
	initDatabase_MongoDB() // Khởi tạo kết nối database
}

// Hàm khởi tạo tên các collection trong database
func initColNames() {
	global.MongoDB_ColNames.Users = "users"
	global.MongoDB_ColNames.Permissions = "permissions"
	global.MongoDB_ColNames.Roles = "roles"
	global.MongoDB_ColNames.RolePermissions = "role_permissions"
	global.MongoDB_ColNames.UserRoles = "user_roles"
	global.MongoDB_ColNames.Agents = "agents"
	global.MongoDB_ColNames.AccessTokens = "access_tokens"
	global.MongoDB_ColNames.FbPages = "fb_pages"
	global.MongoDB_ColNames.FbConvesations = "fb_conversations"
	global.MongoDB_ColNames.FbMessages = "fb_messages"
	global.MongoDB_ColNames.FbPosts = "fb_posts"
	global.MongoDB_ColNames.PcOrders = "pc_orders"

	logrus.Info("Initialized collection names") // Ghi log thông báo đã khởi tạo tên các collection
}

// Hàm khởi tạo validator
func initValidator() {
	global.Validate = validator.New()
	logrus.Info("Initialized validator") // Ghi log thông báo đã khởi tạo validator
}

// Hàm khởi tạo cấu hình server
func initConfig() {
	var err error
	global.MongoDB_ServerConfig = config.NewConfig()
	if err != nil {
		logrus.Fatalf("Failed to initialize config: %v", err) // Ghi log lỗi nếu khởi tạo cấu hình thất bại
	}
	logrus.Info("Initialized server config") // Ghi log thông báo đã khởi tạo cấu hình server
}

// Hàm khởi tạo kết nối database
func initDatabase_MongoDB() {
	var err error
	global.MongoDB_Session, err = database.GetInstance(global.MongoDB_ServerConfig)
	if err != nil {
		logrus.Fatalf("Failed to get database instance: %v", err) // Ghi log lỗi nếu kết nối database thất bại
	}
	logrus.Info("Connected to MongoDB") // Ghi log thông báo đã kết nối database thành công

	// Khởi tạo các db và collections nếu chưa có
	database.EnsureDatabaseAndCollections(global.MongoDB_Session)
	logrus.Info("Ensured database and collections") // Ghi log thông báo đã đảm bảo database và các collection

	// Khởi tạo registry và đăng ký các collections
	err = registry.InitCollections(global.MongoDB_Session, global.MongoDB_ServerConfig)
	if err != nil {
		logrus.Fatalf("Failed to initialize collections: %v", err)
	}
	logrus.Info("Initialized collection registry")

	// Khơi tạo các index cho các collection
	dbName := global.MongoDB_ServerConfig.MongoDB_DBNameAuth
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.Users), models.User{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.Permissions), models.Permission{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.Roles), models.Role{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.UserRoles), models.UserRole{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.RolePermissions), models.RolePermission{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.Agents), models.Agent{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.AccessTokens), models.AccessToken{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.FbPages), models.FbPage{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.FbConvesations), models.FbConversation{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.FbMessages), models.FbMessage{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.FbPosts), models.FbPost{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.PcOrders), models.PcOrder{})

	// Khởi tạo InitService và xử lý error
	initService, err := services.NewInitService()
	if err != nil {
		logrus.Fatalf("Failed to create init service: %v", err)
	}

	// Gọi hàm khởi tạo các quyền mặc định
	if err := initService.InitPermission(); err != nil {
		logrus.Fatalf("Failed to initialize permissions: %v", err)
	}

	if err := initService.CheckPermissionForAdministrator(); err != nil {
		logrus.Fatalf("Failed to check administrator permissions: %v", err)
	}
}

// initFiberApp khởi tạo ứng dụng Fiber với các middleware cần thiết
func initFiberApp() *fiber.App {
	// Khởi tạo app với cấu hình nâng cao
	app := fiber.New(fiber.Config{
		// =========================================
		// 1. CẤU HÌNH CƠ BẢN
		// =========================================
		AppName:       "Meta Commerce API", // Tên ứng dụng hiển thị
		ServerHeader:  "Meta Commerce API", // Header server trong response
		StrictRouting: true,                // /foo và /foo/ là khác nhau
		CaseSensitive: true,                // /Foo và /foo là khác nhau
		UnescapePath:  true,                // Tự động decode URL-encoded paths
		Immutable:     false,               // Tính năng immutable cho ctx

		// =========================================
		// 2. CẤU HÌNH PERFORMANCE
		// =========================================
		BodyLimit:       10 * 1024 * 1024, // Max size của request body (10MB)
		Concurrency:     256 * 1024,       // Số lượng goroutines tối đa
		ReadBufferSize:  4096,             // Buffer size cho request reading
		WriteBufferSize: 4096,             // Buffer size cho response writing

		// =========================================
		// 3. CẤU HÌNH TIMEOUT
		// =========================================
		ReadTimeout:  5 * time.Second,   // Timeout đọc request
		WriteTimeout: 10 * time.Second,  // Timeout ghi response
		IdleTimeout:  120 * time.Second, // Timeout cho idle connections

		// =========================================
		// 4. CẤU HÌNH ERROR HANDLING
		// =========================================
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			// Log error
			logrus.WithFields(logrus.Fields{
				"code":    code,
				"message": message,
				"path":    c.Path(),
				"method":  c.Method(),
				"ip":      c.IP(),
			}).Error("Request error")

			// Return JSON error
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"code":    code,
				"message": message,
				"time":    time.Now().Format(time.RFC3339),
			})
		},
	})

	// =========================================
	// MIDDLEWARE STACK
	// =========================================

	// 1. Recover Middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		Next: func(c fiber.Ctx) bool {
			return false
		},
	}))

	// 2. Logger Middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${ip} | ${status} | ${latency} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Ho_Chi_Minh",
		Output:     os.Stdout,
		Next: func(c fiber.Ctx) bool {
			return c.Path() == "/health"
		},
	}))

	// 3. CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split("*", ","),                                                       // Cho phép tất cả origins, trong thực tế nên giới hạn domain cụ thể
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"},                     // Các HTTP methods được phép
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"}, // Các headers được phép
		//AllowCredentials: true,                                                                          // Cho phép gửi credentials (cookies, auth headers)
		ExposeHeaders: []string{"Content-Length", "Content-Range"}, // Headers cho phép client đọc từ response
		MaxAge:        24 * 60 * 60,                                // Thời gian cache preflight requests (24 giờ)
	}))

	// Khởi tạo routes trước khi đăng ký response middleware
	router.SetupRoutes(app)

	return app
}

// main_thread khởi tạo và chạy Fiber server
func main_thread() {
	// Khởi tạo app với cấu hình
	app := initFiberApp()

	// Khởi động server với cấu hình listen
	logrus.Info("Starting Fiber server...")
	listenConfig := fiber.ListenConfig{}
	if err := app.Listen(":"+global.MongoDB_ServerConfig.Address, listenConfig); err != nil {
		logrus.Fatalf("Error in Fiber ListenAndServe: %v", err)
	}
}

// Hàm main
func main() {
	// Khởi tạo các biến toàn cục
	initGlobal()

	// Chạy Fiber server
	main_thread()
}
