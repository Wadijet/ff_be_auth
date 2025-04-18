package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"

	"meta_commerce/core/global"
)

// main_thread khởi tạo và chạy Fiber server
func main_thread() {
	// Khởi tạo app với cấu hình
	app := InitFiberApp()

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
	InitGlobal()

	// Khởi tạo registry
	InitRegistry()

	// Chạy Fiber server
	main_thread()
}
