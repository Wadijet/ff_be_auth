package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"atk-go-server/global"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitDatabases khởi tạo các cơ sở dữ liệu dựa trên metadata
func InitDatabases() {
	// Tải các biến môi trường từ file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Đọc thông tin về các cơ sở dữ liệu từ biến global.Metadata
	databases, ok := global.Metadata["databases"].([]interface{})
	if !ok {
		log.Fatal("Định dạng databases trong metadata không hợp lệ")
	}

	// Duyệt qua từng cơ sở dữ liệu và kết nối đến nó
	for _, db := range databases {
		dbMap, ok := db.(map[interface{}]interface{})
		if !ok {
			log.Fatal("Định dạng entry của database trong metadata không hợp lệ")
		}

		name, ok := dbMap["name"].(string)
		if !ok {
			log.Fatal("Định dạng tên database trong metadata không hợp lệ")
		}

		connectionURI, ok := dbMap["connectionURI"].(string)
		if !ok {
			log.Fatal("Định dạng connection URI trong metadata không hợp lệ")
		}

		// Kiểm tra và thay thế biến môi trường trong connectionURI nếu có
		if strings.Contains(connectionURI, "${") {
			connectionURI = replaceEnvVariables(connectionURI)
		}

		dbType, ok := dbMap["type"].(string)
		if !ok {
			log.Fatal("Định dạng loại database trong metadata không hợp lệ")
		}

		// Kết nối đến cơ sở dữ liệu dựa trên loại của nó
		switch dbType {
		case "MongoDB":
			// Kết nối đến MongoDB
			client, err := ConnectToMongoDB(connectionURI)
			if err != nil {
				log.Fatalf("Kết nối đến cơ sở dữ liệu MongoDB %s thất bại: %v", name, err)
			}

			// Lưu client vào map dbMap
			dbMap["currentClient"] = client
			fmt.Printf("Đã kết nối đến cơ sở dữ liệu MongoDB %s\n", name)

		case "MySQL":
			db, err := ConnectToMySQL(connectionURI)
			if err != nil {
				log.Fatalf("Kết nối đến cơ sở dữ liệu MySQL %s thất bại: %v", name, err)
			}
			dbMap["currentClient"] = db
			fmt.Printf("Đã kết nối đến cơ sở dữ liệu MySQL %s\n", name)
		default:
			log.Fatalf("Loại cơ sở dữ liệu %s không được hỗ trợ cho cơ sở dữ liệu %s", dbType, name)
		}
	}
}

// replaceEnvVariables thay thế các biến môi trường trong chuỗi
func replaceEnvVariables(input string) string {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			key := fmt.Sprintf("${%s}", pair[0])
			value := pair[1]
			input = strings.ReplaceAll(input, key, value)
		}
	}
	return input
}

// DisconnectDatabases ngắt kết nối tất cả các client cơ sở dữ liệu
func DisconnectDatabases() {
	databases, ok := global.Metadata["databases"].([]interface{})
	if !ok {
		log.Fatal("Định dạng databases trong metadata không hợp lệ")
	}

	for _, db := range databases {
		dbMap, ok := db.(map[string]interface{})
		if !ok {
			log.Fatal("Định dạng entry của database trong metadata không hợp lệ")
		}

		name, ok := dbMap["name"].(string)
		if !ok {
			log.Fatal("Định dạng tên database trong metadata không hợp lệ")
		}

		client, ok := dbMap["currentClient"]
		if !ok {
			log.Printf("Không tìm thấy client cho cơ sở dữ liệu %s", name)
			continue
		}

		switch c := client.(type) {
		case *mongo.Client:
			err := c.Disconnect(context.Background())
			if err != nil {
				log.Printf("Ngắt kết nối client MongoDB %s thất bại: %v", name, err)
			} else {
				fmt.Printf("Đã ngắt kết nối client MongoDB %s\n", name)
			}
		case *sql.DB:
			err := c.Close()
			if err != nil {
				log.Printf("Đóng client MySQL %s thất bại: %v", name, err)
			} else {
				fmt.Printf("Đã đóng client MySQL %s\n", name)
			}
		}
	}
}
