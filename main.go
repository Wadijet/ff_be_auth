package main

import (
	"context"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	validator "gopkg.in/go-playground/validator.v9"

	"atk-go-server/app/middleware"
	models "atk-go-server/app/models/mongodb"
	api "atk-go-server/app/router"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/database"
	"atk-go-server/global"
	"atk-go-server/metadata"
)

// Hàm khởi tạo các biến toàn cục
func initGlobal() {

	initColNames()         // Khởi tạo tên các collection trong database
	initValidator()        // Khởi tạo validator
	initConfig()           // Khởi tạo cấu hình server
	initDatabase_MongoDB() // Khởi tạo kết nối database

}

// Hàm khởi tạo tên các collection trong database
func initColNames() {
	global.MongoDB_ColNames.Users = "users"
	global.MongoDB_ColNames.Permissions = "permissions"
	global.MongoDB_ColNames.Roles = "roles"
	global.MongoDB_ColNames.RolePermissions = "role_permissions"
	global.MongoDB_ColNames.UserRoles = "user_roles"
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
	dbName := global.MongoDB_ServerConfig.MongoDB_DBNameAuth
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

	// gọi hàm khởi tạo các quyền mặc định
	InitService := services.NewInitService(global.MongoDB_ServerConfig, global.MongoDB_Session)
	InitService.InitPermission()
	InitService.CheckPermissionForAdministrator()

}

func initMetadata() {

	readFileMetadata()       // Khởi tạo metadata
	database.InitDatabases() // Khởi tạo các cơ sở dữ liệu dựa trên metadata
}

// Hàm khởi tạo metadata
func readFileMetadata() {
	metadataFilePath := global.MongoDB_ServerConfig.Metadata_Path
	result, err := metadata.ReadMetadata(metadataFilePath)
	if err != nil {
		logrus.Fatalf("Failed to read metadata: %v", err)
	}
	global.Metadata = result
	logrus.Info("Initialized metadata")
}

// Hàm khởi tạo kết nối database
func initDatabase_MySql() {
	var err error

	global.MySQL_Session, err = database.GetMySQLInstance(global.MongoDB_ServerConfig)
	if err != nil {
		logrus.Fatalf("Failed to get MySQL instance: %v", err)
	}
	logrus.Info("Connected to MySQL")
}

// Hàm xử lý panic
func panicHandler(ctx *fasthttp.RequestCtx, data interface{}) {
	logrus.Errorf("Panic occurred: %v", data)
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)         // Ghi log lỗi khi xảy ra panic
	utility.JSON(ctx, utility.Payload(false, data, "Lỗi panic!")) // Trả về JSON thông báo lỗi panic
}

// Hàm chính để chạy server
func main_thread() {

	// Khởi tạo router
	r := router.New()
	api.InitRounters(r, global.MongoDB_ServerConfig, global.MongoDB_Session) // Khởi tạo các route cho API
	r.PanicHandler = panicHandler                                            // Đặt hàm xử lý panic

	// Sử dụng middleware Measure và COSR cho tất cả các route
	measuredHandler := middleware.CORS(middleware.Measure(r.Handler))

	// Chạy server
	logrus.Info("Starting server...") // Ghi log thông báo bắt đầu chạy server
	if err := fasthttp.ListenAndServe(":8080", measuredHandler); err != nil {
		logrus.Fatalf("Error in ListenAndServe: %v", err) // Ghi log lỗi nếu server không thể chạy
	}
}

// Hàm main
func main() {
	initGlobal()   // Khởi tạo các biến toàn cục
	initMetadata() // Khởi tạo metadata
	main_thread()  // Chạy server
}
