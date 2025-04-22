package main

import (
	"context"
	"meta_commerce/config"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/database"
	"meta_commerce/core/global"

	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// Hàm khởi tạo các biến toàn cục
func InitGlobal() {
	initColNames()         // Khởi tạo tên các collection trong database
	initValidator()        // Khởi tạo validator
	initConfig()           // Khởi tạo cấu hình server
	initDatabase_MongoDB() // Khởi tạo kết nối database
}

// Hàm khởi tạo tên các collection trong database
func initColNames() {
	global.MongoDB_ColNames.Users = "auth_users"
	global.MongoDB_ColNames.Permissions = "auth_permissions"
	global.MongoDB_ColNames.Roles = "auth_roles"
	global.MongoDB_ColNames.RolePermissions = "auth_role_permissions"
	global.MongoDB_ColNames.UserRoles = "auth_user_roles"
	global.MongoDB_ColNames.Agents = "agents"
	global.MongoDB_ColNames.AccessTokens = "access_tokens"
	global.MongoDB_ColNames.FbPages = "fb_pages"
	global.MongoDB_ColNames.FbConvesations = "fb_conversations"
	global.MongoDB_ColNames.FbMessages = "fb_messages"
	global.MongoDB_ColNames.FbPosts = "fb_posts"
	global.MongoDB_ColNames.PcOrders = "pc_orders"

	logrus.Info("Initialized collection names") // Ghi log thông báo đã khởi tạo tên các collection
}

// Hàm khởi tạo validator
func initValidator() {
	global.Validate = validator.New()
	logrus.Info("Initialized validator") // Ghi log thông báo đã khởi tạo validator
}

// Hàm khởi tạo cấu hình server
func initConfig() {
	var err error
	global.MongoDB_ServerConfig = config.NewConfig()
	if err != nil {
		logrus.Fatalf("Failed to initialize config: %v", err) // Ghi log lỗi nếu khởi tạo cấu hình thất bại
	}
	logrus.Info("Initialized server config") // Ghi log thông báo đã khởi tạo cấu hình server
}

// Hàm khởi tạo kết nối database
func initDatabase_MongoDB() {
	var err error
	global.MongoDB_Session, err = database.GetInstance(global.MongoDB_ServerConfig)
	if err != nil {
		logrus.Fatalf("Failed to get database instance: %v", err) // Ghi log lỗi nếu kết nối database thất bại
	}
	logrus.Info("Connected to MongoDB") // Ghi log thông báo đã kết nối database thành công

	// Khởi tạo các db và collections nếu chưa có
	database.EnsureDatabaseAndCollections(global.MongoDB_Session)
	logrus.Info("Ensured database and collections") // Ghi log thông báo đã đảm bảo database và các collection

	// Khơi tạo các index cho các collection
	dbName := global.MongoDB_ServerConfig.MongoDB_DBName_Auth
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.Users), models.User{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.Permissions), models.Permission{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.Roles), models.Role{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.UserRoles), models.UserRole{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.RolePermissions), models.RolePermission{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.Agents), models.Agent{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.AccessTokens), models.AccessToken{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.FbPages), models.FbPage{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.FbConvesations), models.FbConversation{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.FbMessages), models.FbMessage{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.FbPosts), models.FbPost{})
	database.CreateIndexes(context.TODO(), global.MongoDB_Session.Database(dbName).Collection(global.MongoDB_ColNames.PcOrders), models.PcOrder{})
}
