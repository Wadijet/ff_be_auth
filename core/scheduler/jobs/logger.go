package jobs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var jobLogger *logrus.Logger

func init() {
	// Lấy đường dẫn thực thi hiện tại
	executable, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Could not get executable path: %v", err))
	}

	// Lấy đường dẫn gốc của project (2 cấp trên thư mục cmd)
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(executable)))
	logPath := filepath.Join(rootDir, "logs")
	jobLogFile := filepath.Join(logPath, "jobs.log")

	// Log đường dẫn để debug
	fmt.Printf("Job Log File: %s\n", jobLogFile)

	// Tạo thư mục logs nếu chưa tồn tại
	if err := os.MkdirAll(logPath, 0755); err != nil {
		panic(fmt.Sprintf("Could not create logs directory at %s: %v", logPath, err))
	}

	// Tạo logger mới cho jobs
	jobLogger = logrus.New()

	// Cấu hình format với thông tin file, line number và goroutine ID
	jobLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, filepath.Base(f.File)
		},
	})

	// Mở file log với full path
	logFile, err := os.OpenFile(jobLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Could not open job log file %s: %v", jobLogFile, err))
	}

	// Kiểm tra xem file có thể ghi được không
	if _, err := logFile.Write([]byte("Job log file initialized\n")); err != nil {
		panic(fmt.Sprintf("Could not write to job log file %s: %v", jobLogFile, err))
	}

	// Ghi log ra cả stdout và file
	mw := io.MultiWriter(os.Stdout, logFile)
	jobLogger.SetOutput(mw)

	// Bật caller logging để hiển thị thông tin file và line number
	jobLogger.SetReportCaller(true)

	// Set log level (có thể điều chỉnh theo environment)
	jobLogger.SetLevel(logrus.DebugLevel)

	// Log thông tin khởi tạo
	jobLogger.WithFields(logrus.Fields{
		"log_file": jobLogFile,
		"level":    jobLogger.GetLevel().String(),
	}).Info("Job logger initialized successfully")
}

// GetJobLogger trả về instance của jobLogger để sử dụng bên ngoài package
func GetJobLogger() *logrus.Logger {
	return jobLogger
}
