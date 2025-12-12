package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
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
	cfg := global.MongoDB_ServerConfig
	address := ":" + cfg.Address
	
	logrus.Info("Starting Fiber server...")
	
	// Helper function để resolve đường dẫn từ thư mục api
	resolvePath := func(path string) string {
		if filepath.IsAbs(path) {
			return path
		}
		// Tìm thư mục api
		currentDir, err := os.Getwd()
		if err != nil {
			return path
		}
		for {
			envDir := filepath.Join(currentDir, "config", "env")
			if _, err := os.Stat(envDir); err == nil {
				return filepath.Join(currentDir, path)
			}
			parentDir := filepath.Dir(currentDir)
			if parentDir == currentDir {
				return path
			}
			currentDir = parentDir
		}
	}

	// Kiểm tra xem có bật TLS không
	if cfg.EnableTLS && cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
		// Resolve đường dẫn certificate và key từ thư mục api
		certPath := resolvePath(cfg.TLSCertFile)
		keyPath := resolvePath(cfg.TLSKeyFile)
		
		// Kiểm tra file certificate và key tồn tại
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			logrus.Fatalf("TLS certificate file not found: %s (resolved from: %s)", certPath, cfg.TLSCertFile)
		}
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			logrus.Fatalf("TLS key file not found: %s (resolved from: %s)", keyPath, cfg.TLSKeyFile)
		}
		
		// Load certificate và key
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			logrus.Fatalf("Error loading TLS certificate: %v", err)
		}
		
		// Tạo listener với TLS
		ln, err := net.Listen("tcp", address)
		if err != nil {
			logrus.Fatalf("Error creating listener: %v", err)
		}
		
		// Cấu hình TLS
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		
		// Wrap listener với TLS
		tlsListener := tls.NewListener(ln, tlsConfig)
		
		logrus.WithFields(logrus.Fields{
			"address": address,
			"cert":    certPath,
			"key":     keyPath,
		}).Info("Starting server with HTTPS/TLS")
		
		// Khởi động server với TLS listener
		if err := app.Listener(tlsListener); err != nil {
			logrus.Fatalf("Error in Fiber Listener with TLS: %v", err)
		}
	} else {
		// Khởi động server HTTP thông thường
		logrus.WithFields(logrus.Fields{
			"address":  address,
			"protocol": "HTTP",
		}).Info("Starting server with HTTP")
		
		listenConfig := fiber.ListenConfig{}
		if err := app.Listen(address, listenConfig); err != nil {
			logrus.Fatalf("Error in Fiber Listen: %v", err)
		}
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

	// Chạy Fiber server trên main thread
	main_thread()
}
