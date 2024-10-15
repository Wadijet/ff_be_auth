package config

import (
	"fmt"
	"log"

	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Configuration chứa thông tin tĩnh cần thiết để chạy ứng dụng
// Nó chứa thông tin cơ sở dữ liệu
type Configuration struct {
	InitMode              bool   `env:"INITMODE" envDefault:false`   // Chế độ khởi tạo
	Address               string `env:"ADDRESS" envDefault:":8080"`  // Địa chỉ server
	JwtSecret             string `env:"JWT_SECRET,required"`         // Bí mật JWT
	DataBaseConnectionURL string `env:"CONNECTION_URL,required"`     // URL kết nối cơ sở dữ liệu
	DataBaseNameAuth      string `env:"DATABASE_NAME_AUTH,required"` // Tên cơ sở dữ liệu xác thực
}

// NewConfig sẽ đọc dữ liệu cấu hình từ file .env được cung cấp
func NewConfig(files ...string) *Configuration {
	err := godotenv.Load(filepath.Join(".env")) // Tải cấu hình từ file .env
	if err != nil {
		log.Printf("Không tìm thấy file .env %q\n", files)
	}

	cfg := Configuration{}

	// Phân tích env vào cấu hình
	err = env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}
