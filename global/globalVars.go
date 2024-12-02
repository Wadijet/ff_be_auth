package global

import (
	"atk-go-server/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	validator "gopkg.in/go-playground/validator.v9"
)

// MongoDB_CollectionName chứa tên các collection trong MongoDB
type MongoDB_CollectionName struct {
	Users          string // Tên collection cho người dùng
	Permissions    string // Tên collection cho quyền
	Roles          string // Tên collection cho vai trò
	RolePermission string // Tên collection cho vai trò và quyền
	UserRole       string // Tên collection cho người dùng và vai trò
}

// Các biến toàn cục
var Validate *validator.Validate                                           // Biến để xác thực dữ liệu
var MongoDB_Session *mongo.Client                                          // Phiên kết nối tới MongoDB
var MongoDB_ServerConfig *config.Configuration                             // Cấu hình của server
var MongoDB_ColNames MongoDB_CollectionName = *new(MongoDB_CollectionName) // Tên các collection
var MySQL_Session *sql.DB                                                  // Add this line to define MySQLDB
