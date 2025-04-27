package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"

	"meta_commerce/core/global"
)

// initLogger khởi tạo và cấu hình logger cho toàn bộ ứng dụng
func initLogger() {
	// Lấy đường dẫn thực thi hiện tại
	executable, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Could not get executable path: %v", err))
	}

	// Lấy đường dẫn gốc của project (2 cấp trên thư mục cmd)
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(executable)))
	logPath := filepath.Join(rootDir, "logs")
	appLogFile := filepath.Join(logPath, "app.log")

	// Log đường dẫn để debug
	fmt.Printf("Executable: %s\n", executable)
	fmt.Printf("Root Dir: %s\n", rootDir)
	fmt.Printf("Log Path: %s\n", logPath)
	fmt.Printf("Log File: %s\n", appLogFile)

	// Tạo thư mục logs nếu chưa tồn tại
	if err := os.MkdirAll(logPath, 0755); err != nil {
		panic(fmt.Sprintf("Could not create logs directory at %s: %v", logPath, err))
	}

	// Kiểm tra quyền ghi vào thư mục logs
	tmpFile := filepath.Join(logPath, "test.tmp")
	if err := os.WriteFile(tmpFile, []byte("test"), 0666); err != nil {
		panic(fmt.Sprintf("No write permission in logs directory %s: %v", logPath, err))
	}
	os.Remove(tmpFile)

	// Mở file log với full path
	logFile, err := os.OpenFile(appLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Could not open log file %s: %v", appLogFile, err))
	}

	// Kiểm tra xem file có thể ghi được không
	if _, err := logFile.Write([]byte("Log file initialized\n")); err != nil {
		panic(fmt.Sprintf("Could not write to log file %s: %v", appLogFile, err))
	}

	// Cấu hình format với thông tin file, line number và goroutine ID
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, filepath.Base(f.File)
		},
	})

	// Ghi log ra cả stdout và file
	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)

	// Bật caller logging để hiển thị thông tin file và line number
	logrus.SetReportCaller(true)

	// Set log level (có thể điều chỉnh theo environment)
	logrus.SetLevel(logrus.DebugLevel)

	// Log thông tin khởi tạo
	logrus.WithFields(logrus.Fields{
		"log_file": appLogFile,
		"level":    logrus.GetLevel().String(),
	}).Info("Logger initialized successfully")
}

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
	// Khởi tạo logger
	initLogger()

	// Khởi tạo các biến toàn cục
	InitGlobal()

	// Khởi tạo registry
	InitRegistry()

	// Khởi tạo dữ liệu mặc định
	InitDefaultData()

	// Khởi tạo và chạy background jobs trong goroutines
	go func() {
		logrus.Info("Initializing background jobs...")
		if err := InitJobs(global.Scheduler); err != nil {
			logrus.WithError(err).Error("Failed to initialize jobs")
			return
		}
		logrus.Info("Background jobs initialized successfully")

		// Start scheduler sau khi khởi tạo jobs thành công
		global.Scheduler.Start()
		logrus.Info("Scheduler started successfully")
	}()

	// Chạy Fiber server trên main thread
	main_thread()
}
