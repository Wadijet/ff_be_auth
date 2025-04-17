package main

import (
	"meta_commerce/core/api/router"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/sirupsen/logrus"
)

// InitFiberApp khởi tạo ứng dụng Fiber với các middleware cần thiết
func InitFiberApp() *fiber.App {
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
