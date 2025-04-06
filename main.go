package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gofiber/contrib/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"meta_commerce/app/database/registry"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
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

	// gọi hàm khởi tạo các quyền mặc định
	InitService := services.NewInitService()
	InitService.InitPermission()
	InitService.CheckPermissionForAdministrator()
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
	// Mục đích: Phục hồi từ panic, tránh crash server
	// Nên bật: LUÔN LUÔN, đây là middleware quan trọng nhất
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true, // Hiển thị stack trace khi panic
		Next: func(c fiber.Ctx) bool {
			return false // Luôn xử lý panic
		},
	}))

	// 2. Logger Middleware
	// Mục đích: Ghi log mọi request để debug và monitor
	// Nên bật: Môi trường development và staging
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${ip} | ${status} | ${latency} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Ho_Chi_Minh",
		Output:     os.Stdout,
		Next: func(c fiber.Ctx) bool {
			return c.Path() == "/health" // Bỏ qua health check
		},
	}))

	// 3. CORS Middleware
	// Mục đích: Cho phép request từ domain khác
	// Nên bật: Khi cần cho phép frontend từ domain khác gọi API
	app.Use(cors.New(cors.Config{
		AllowOrigins:  strings.Split("*", ","),                                                         // Domain được phép gọi API
		AllowMethods:  strings.Split("GET,POST,PUT,DELETE,OPTIONS,PATCH", ","),                         // Method được phép
		AllowHeaders:  strings.Split("Origin,Content-Type,Accept,Authorization,X-Requested-With", ","), // Header được phép
		ExposeHeaders: strings.Split("Content-Length,Content-Range", ","),                              // Header được phép client đọc
		//AllowCredentials: true,                                                                            // Cho phép gửi cookie
		MaxAge: 24 * 60 * 60, // Cache preflight request
	}))

	// 4. Compress Middleware
	// Mục đích: Nén response để giảm bandwidth
	// Nên bật: Khi cần tối ưu performance
	/*app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Mức độ nén (1-9)
		Next: func(c fiber.Ctx) bool {
			return strings.Contains(c.Path(), "/ws") // Không nén WebSocket
		},
	}))*/

	// 5. Cache Middleware
	// Mục đích: Cache response để tăng tốc độ
	// Nên bật: Với các endpoint ít thay đổi data
	/*app.Use(cache.New(cache.Config{
		Next: func(c fiber.Ctx) bool {
			return c.Query("refresh") == "true" // Bỏ qua cache nếu có query refresh
		},
		KeyGenerator: func(c fiber.Ctx) string {
			return c.Path() // Key là path
		},
		CacheControl:         true,             // Thêm header Cache-Control
		StoreResponseHeaders: true,             // Cache cả response header
		MaxBytes:             10 * 1024 * 1024, // Max size cache (10MB)
	}))*/

	// 6. BasicAuth Middleware
	// Mục đích: Bảo vệ API bằng username/password
	// Nên bật: Khi cần xác thực đơn giản
	/*app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "123456", // Username/password
		},
		Realm: "Restricted Area",
		Unauthorized: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	}))*/

	// 7. CSRF Middleware
	// Mục đích: Bảo vệ khỏi tấn công CSRF
	// Nên bật: Khi có form submission
	/*app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token", // Vị trí lưu token
		CookieName:     "csrf_",               // Tên cookie
		CookieSameSite: "Strict",              // SameSite policy
		CookieSecure:   true,                  // Chỉ gửi qua HTTPS
		CookieHTTPOnly: true,                  // Không cho JS đọc cookie
		KeyGenerator: func() string {
			return fmt.Sprintf("%d", time.Now().UnixNano()) // Tạo token ngẫu nhiên
		},
	}))*/

	// 8. Limiter Middleware
	// Mục đích: Giới hạn số request
	// Nên bật: Chống tấn công DDoS
	/*app.Use(limiter.New(limiter.Config{
		Max:        100,             // Số request tối đa
		Expiration: 1 * time.Minute, // Thời gian reset
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP() // Dùng IP làm key
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
		SkipFailedRequests:     false, // Tính cả request lỗi
		SkipSuccessfulRequests: false, // Tính cả request thành công
	}))*/

	// 9. RequestID Middleware
	// Mục đích: Gắn ID cho mỗi request
	// Nên bật: Khi cần track request
	/*app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		},
		Next: func(c fiber.Ctx) bool {
			return false // Luôn tạo request ID
		},
	}))*/

	// 10. Helmet Middleware
	// Mục đích: Thêm security headers
	// Nên bật: Trên môi trường production
	/*app.Use(helmet.New(helmet.Config{
		ContentSecurityPolicy: "default-src 'self';",
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            31536000,
		ReferrerPolicy:        "same-origin",
	}))*/

	// 11. EncryptCookie Middleware
	// Mục đích: Mã hóa cookie
	// Nên bật: Khi lưu thông tin nhạy cảm
	/*app.Use(encryptcookie.New(encryptcookie.Config{
		Key:    "secret-thirty-2-character-string", // Key mã hóa 32 ký tự
		Except: []string{"non_encrypted_cookie"},   // Cookie không cần mã hóa
	}))*/

	// 12. ETag Middleware
	// Mục đích: Tối ưu cache bằng ETag
	// Nên bật: Khi muốn tiết kiệm bandwidth
	/*app.Use(etag.New(etag.Config{
		Weak: false, // Dùng strong ETag
		Next: func(c fiber.Ctx) bool {
			return c.Method() != "GET" // Chỉ tạo ETag cho GET
		},
	}))*/

	// 13. Favicon Middleware
	// Mục đích: Xử lý request favicon
	// Nên bật: Khi muốn tối ưu log
	//app.Use(favicon.New(favicon.Config{
	//	File:         "./favicon.ico", // Path tới file favicon
	//	CacheControl: "public, max-age=31536000",
	//}))

	// 14. Adaptor Middleware
	// Mục đích: Chuyển đổi net/http handlers sang/từ Fiber handlers
	// Nên bật: Khi cần tích hợp với các handler net/http có sẵn
	/*app.Use("/legacy", adaptor.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from net/http handler")
	})))*/

	// 15. EarlyData Middleware
	// Mục đích: Hỗ trợ tính năng early data ("0-RTT") của TLS 1.3
	// Nên bật: Khi cần tối ưu hiệu suất với TLS 1.3
	/*app.Use(earlydata.New())*/

	// 16. EnvVar Middleware
	// Mục đích: Expose biến môi trường với cấu hình tùy chọn
	// Nên bật: Khi cần expose biến môi trường an toàn
	/*app.Use(envvar.New(envvar.Config{
		ExportVars: map[string]string{
			"APP_NAME":    os.Getenv("APP_NAME"),
			"APP_VERSION": os.Getenv("APP_VERSION"),
		},
	}))*/

	// 17. ExpVar Middleware
	// Mục đích: Phục vụ các biến runtime qua HTTP server dưới dạng JSON
	// Nên bật: Khi cần monitor runtime variables
	/*app.Use(expvar.New())*/

	// 18. Healthcheck Middleware
	// Mục đích: Kiểm tra health của server
	// Nên bật: Khi triển khai trong môi trường container/kubernetes
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	// 19. Idempotency Middleware
	// Mục đích: Cho phép API chịu lỗi khi có request trùng lặp
	// Nên bật: Khi cần đảm bảo request trùng lặp không gây lỗi
	/*app.Use(idempotency.New(idempotency.Config{
		Lifetime:  24 * time.Hour, // Thời gian lưu key
		KeyHeader: "X-Idempotency-Key",
	}))*/

	// 20. Keyauth Middleware
	// Mục đích: Hỗ trợ xác thực dựa trên key
	// Nên bật: Khi cần xác thực API key
	/*app.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:X-API-Key",
		Validator: func(c fiber.Ctx, key string) (bool, error) {
			return key == "valid-api-key", nil
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API Key",
			})
		},
	}))*/

	// 21. Pprof Middleware
	// Mục đích: Phục vụ dữ liệu profiling runtime dưới dạng pprof
	// Nên bật: Khi cần debug performance
	/*app.Use(pprof.New())*/

	// 22. Proxy Middleware
	// Mục đích: Cho phép proxy requests tới nhiều server
	// Nên bật: Khi cần load balancing
	/*app.Get("/proxy/:server", func(c fiber.Ctx) error {
		server := c.Params("server")
		url := fmt.Sprintf("http://localhost:300%s", server)

		// Forward request
		if resp, err := http.Get(url); err == nil {
			return c.Status(resp.StatusCode).SendString("Proxied")
		}
		return c.SendStatus(fiber.StatusBadGateway)
	})*/

	// 23. Redirect Middleware
	// Mục đích: Middleware chuyển hướng
	// Nên bật: Khi cần chuyển hướng URL
	/*app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/old":       "/new",
			"/old-api/*": "/api/$1",
		},
		StatusCode: 301,
	}))*/

	// 24. Rewrite Middleware
	// Mục đích: Viết lại đường dẫn URL dựa trên rules
	// Nên bật: Khi cần viết lại URL cho backward compatibility
	/*app.Use(rewrite.New(rewrite.Config{
		Rules: map[string]string{
			"/api/v1/*":   "/api/v2/$1",
			"/blog/:post": "/blog/:post/comments",
		},
	}))*/

	// 25. Monitor Middleware
	// Mục đích: Hiển thị metrics của server như CPU, Memory, số request, v.v.
	// Nên bật: Khi cần monitor hiệu năng server
	app.Get("/metrics", monitor.New(monitor.Config{
		Title: "Meta Commerce Metrics",
	}))

	return app
}

// main_thread khởi tạo và chạy Fiber server
func main_thread() {
	// Khởi tạo app với cấu hình
	app := initFiberApp()

	// Thiết lập routes
	router.SetupFiberRoutes(app)

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
