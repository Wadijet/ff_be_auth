package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	loggers   = make(map[string]*logrus.Logger)
	loggersMu sync.Mutex
	rootDir   string
)

// Lấy rootDir của project (2 cấp trên thư mục cmd)
func getRootDir() string {
	if rootDir != "" {
		return rootDir
	}
	executable, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Could not get executable path: %v", err))
	}
	rootDir = filepath.Dir(filepath.Dir(filepath.Dir(executable)))
	return rootDir
}

// GetLogger trả về logger theo tên (app, jobs, ...)
func GetLogger(name string) *logrus.Logger {
	loggersMu.Lock()
	defer loggersMu.Unlock()

	if logger, ok := loggers[name]; ok {
		return logger
	}

	// Tạo logger mới nếu chưa có
	logPath := filepath.Join(getRootDir(), "logs")
	logFile := filepath.Join(logPath, fmt.Sprintf("%s.log", name))

	if err := os.MkdirAll(logPath, 0755); err != nil {
		panic(fmt.Sprintf("Could not create logs directory at %s: %v", logPath, err))
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Could not open log file %s: %v", logFile, err))
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, filepath.Base(f.File)
		},
	})
	mw := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(mw)
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.DebugLevel)

	logger.WithFields(logrus.Fields{
		"log_file": logFile,
		"level":    logger.GetLevel().String(),
	}).Info("Logger initialized successfully")

	loggers[name] = logger
	return logger
}
