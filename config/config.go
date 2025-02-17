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
	InitMode              bool   `env:"INITMODE" envDefault:false`       // Chế độ khởi tạo
	Address               string `env:"ADDRESS" envDefault:":8080"`      // Địa chỉ server
	JwtSecret             string `env:"JWT_SECRET,required"`             // Bí mật JWT
	MongoDB_ConnectionURL string `env:"MONGODB_CONNECTION_URL,required"` // URL kết nối cơ sở dữ liệu
	MongoDB_DBNameAuth    string `env:"MONGODB_DBNAME_AUTH,required"`    // Tên cơ sở dữ liệu xác thực
	MySQLConnectionURL    string `env:"MYSQL_CONNECTION_URL,required"`   // URL kết nối MySQL

	Metadata_Path    string `env:"METADATA_PATH,required"`    // Đường dẫn đến file metadata
	MongoDb_Uri_Auth string `env:"MONGODB_URI_AUTH,required"` // URI kết nối MongoDB
	MongoDb_Uri_Data string `env:"MONGODB_URI_DATA,required"` // URI kết nối MongoDB
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
