package global

import (
	"atk-go-server/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	validator "gopkg.in/go-playground/validator.v9"
)

// MongoDB_Auth_CollectionName chứa tên các collection trong MongoDB
type MongoDB_Auth_CollectionName struct {
	Users           string // Tên collection cho người dùng
	Permissions     string // Tên collection cho quyền
	Roles           string // Tên collection cho vai trò
	RolePermissions string // Tên collection cho vai trò và quyền
	UserRoles       string // Tên collection cho người dùng và vai trò
	Agents          string // Tên collection cho bot
	AccessTokens    string // Tên collection cho token
	FbPages         string // Tên collection cho trang Facebook
	FbConvesations  string // Tên collection cho cuộc trò chuyện trên Facebook
	FbMessages      string // Tên collection cho tin nhắn trên Facebook
}

// Các biến toàn cục
var Validate *validator.Validate                                                     // Biến để xác thực dữ liệu
var MongoDB_Session *mongo.Client                                                    // Phiên kết nối tới MongoDB
var MongoDB_ServerConfig *config.Configuration                                       // Cấu hình của server
var MongoDB_ColNames MongoDB_Auth_CollectionName = *new(MongoDB_Auth_CollectionName) // Tên các collection
var MySQL_Session *sql.DB                                                            // Add this line to define MySQLDB

// Tên các collection trong MongoDB
