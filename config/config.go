package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Configuration chứa thông tin tĩnh cần thiết để chạy ứng dụng
// Nó chứa thông tin cơ sở dữ liệu
type Configuration struct {
	InitMode               bool   `env:"INITMODE" envDefault:"false"`     // Chế độ khởi tạo
	Address                string `env:"ADDRESS" envDefault:":8080"`      // Địa chỉ server
	JwtSecret              string `env:"JWT_SECRET,required"`             // Bí mật JWT
	MongoDB_ConnectionURI  string `env:"MONGODB_CONNECTION_URI,required"` // URL kết nối cơ sở dữ liệu
	MongoDB_DBName_Auth    string `env:"MONGODB_DBNAME_AUTH,required"`    // Tên cơ sở dữ liệu xác thực
	MongoDB_DBName_Staging string `env:"MONGODB_DBNAME_STAGING,required"` // Tên cơ sở dữ liệu staging
	MongoDB_DBName_Data    string `env:"MONGODB_DBNAME_DATA,required"`    // Tên cơ sở dữ liệu data
}

// getEnvPath trả về đường dẫn đến file env dựa trên môi trường
func getEnvPath() string {
	// Mặc định sử dụng môi trường development
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Tìm thư mục config
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Không thể lấy được thư mục hiện tại: %v\n", err)
		return ""
	}

	// Tìm thư mục config/env
	for {
		envDir := filepath.Join(currentDir, "config", "env")
		if _, err := os.Stat(envDir); err == nil {
			// Tìm thấy thư mục config/env
			envPath := filepath.Join(envDir, fmt.Sprintf("%s.env", env))
			return envPath
		}

		// Đi lên thư mục cha
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return ""
		}
		currentDir = parentDir
	}
}

// NewConfig sẽ đọc dữ liệu cấu hình từ file env được cung cấp
func NewConfig(files ...string) *Configuration {
	envPath := getEnvPath()
	if envPath == "" {
		log.Printf("Không tìm thấy thư mục config/env")
		return nil
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Không thể load file env tại %s: %v\n", envPath, err)
		return nil
	}

	cfg := Configuration{}
	err = env.Parse(&cfg)
	if err != nil {
		fmt.Printf("Lỗi khi parse config: %+v\n", err)
		return nil
	}

	return &cfg
}
