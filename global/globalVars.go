package global

import (
	"atk-go-server/config"

	"go.mongodb.org/mongo-driver/mongo"
	validator "gopkg.in/go-playground/validator.v9"
)

// CollectionName chứa tên các collection trong MongoDB
type CollectionName struct {
	Permissions string // Tên collection cho quyền
	Roles       string // Tên collection cho vai trò
	Users       string // Tên collection cho người dùng
	MtServices  string // Tên collection cho dịch vụ MT
}

// Các biến toàn cục
var Validate *validator.Validate                   // Biến để xác thực dữ liệu
var DbSession *mongo.Client                        // Phiên kết nối tới MongoDB
var ServerConfig *config.Configuration             // Cấu hình của server
var ColNames CollectionName = *new(CollectionName) // Tên các collection
