package logger

import (
	"os"
	"strings"
)

// LogConfig chứa cấu hình cho hệ thống logging
type LogConfig struct {
	// Log Level: trace, debug, info, warn, error, fatal
	Level string `env:"LOG_LEVEL" envDefault:"info"`

	// Log Format: json, text
	Format string `env:"LOG_FORMAT" envDefault:"text"`

	// Log Output: file, stdout, both
	Output string `env:"LOG_OUTPUT" envDefault:"both"`

	// Log Rotation
	MaxSize    int  `env:"LOG_MAX_SIZE" envDefault:"100"`    // MB
	MaxBackups int  `env:"LOG_MAX_BACKUPS" envDefault:"7"`    // Số file cũ giữ lại
	MaxAge     int  `env:"LOG_MAX_AGE" envDefault:"7"`       // Số ngày giữ lại
	Compress   bool `env:"LOG_COMPRESS" envDefault:"true"`    // Nén file cũ

	// Log Paths
	LogPath         string `env:"LOG_PATH" envDefault:"./logs"`
	AppFile         string `env:"LOG_APP_FILE" envDefault:"app.log"`
	AuditFile       string `env:"LOG_AUDIT_FILE" envDefault:"audit.log"`
	PerformanceFile string `env:"LOG_PERF_FILE" envDefault:"performance.log"`
	ErrorFile       string `env:"LOG_ERROR_FILE" envDefault:"error.log"`
}

// DefaultConfig trả về cấu hình mặc định
func DefaultConfig() *LogConfig {
	// Lấy environment
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	config := &LogConfig{
		Level:           "info",
		Format:          "text",
		Output:          "both",
		MaxSize:         100,
		MaxBackups:      7,
		MaxAge:          7,
		Compress:        true,
		LogPath:         "./logs",
		AppFile:         "app.log",
		AuditFile:       "audit.log",
		PerformanceFile: "performance.log",
		ErrorFile:       "error.log",
	}

	// Điều chỉnh theo môi trường
	if env == "development" {
		config.Level = "debug"
		config.Format = "text"
	} else {
		config.Level = "info"
		config.Format = "json"
	}

	// Override từ environment variables
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Level = strings.ToLower(level)
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Format = strings.ToLower(format)
	}
	if output := os.Getenv("LOG_OUTPUT"); output != "" {
		config.Output = strings.ToLower(output)
	}

	return config
}
