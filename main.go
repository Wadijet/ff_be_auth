package main

import (
	"context"
	"time"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	validator "gopkg.in/go-playground/validator.v9"

	"meta_commerce/app/database/registry"
	"meta_commerce/app/global"
	"meta_commerce/app/middleware"
	models "meta_commerce/app/models/mongodb"
	api "meta_commerce/app/router"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
	"meta_commerce/config"
	"meta_commerce/database"
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

	// Khởi tạo registry và đăng ký các collections
	err = registry.InitCollections(global.MongoDB_Session, global.MongoDB_ServerConfig)
	if err != nil {
		logrus.Fatalf("Failed to initialize collections: %v", err)
	}
	logrus.Info("Initialized collection registry")

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
	InitService, err := services.NewInitService(global.MongoDB_ServerConfig, global.MongoDB_Session)
	if err != nil {
		logrus.Fatalf("Failed to initialize init service: %v", err)
	}

	InitService.InitPermission()
	InitService.CheckPermissionForAdministrator()
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
	r := router.New()                                                        // Khởi tạo router mới
	api.InitRounters(r, global.MongoDB_ServerConfig, global.MongoDB_Session) // Khởi tạo các route cho router
	r.PanicHandler = panicHandler                                            // Xử lý panic

	// Sử dụng các middleware theo thứ tự
	handler := r.Handler
	handler = middleware.Recovery(handler)                                   // Recovery middleware để bắt panic
	handler = middleware.Timeout(time.Second * 30)(handler)                  // Timeout middleware
	handler = middleware.RateLimit(100, time.Second)(handler)                // Rate limiting
	handler = middleware.Measure(middleware.DefaultMeasureConfig())(handler) // Measure middleware với cấu hình mặc định
	handler = middleware.CORS(middleware.DefaultCorsConfig())(handler)       // CORS middleware với cấu hình mặc định

	// Cấu hình server với các timeout
	server := &fasthttp.Server{
		Handler:               handler,          // Thêm handler vào server
		ReadTimeout:           time.Second * 10, // Thời gian đọc request tối đa là 10 giây
		WriteTimeout:          time.Second * 10, // Thời gian ghi response tối đa là 10 giây
		MaxRequestsPerConn:    1000,             // Số lượng request tối đa trên mỗi connection là 1000
		MaxConnsPerIP:         100,              // Số lượng connection tối đa trên mỗi IP là 100
		MaxKeepaliveDuration:  time.Second * 5,  // Thời gian keepalive tối đa là 5 giây
		IdleTimeout:           time.Second * 5,  // Thời gian idle tối đa là 5 giây
		ReadBufferSize:        4096,             // Kích thước buffer đọc là 4096 bytes
		WriteBufferSize:       4096,             // Kích thước buffer ghi là 4096 bytes
		MaxRequestBodySize:    10 * 1024 * 1024, // Kích thước body request tối đa là 10MB
		NoDefaultServerHeader: true,             // Không thêm header server mặc định
		DisableKeepalive:      false,            // Không tắt keepalive
	}

	// Chạy server với cấu hình đã thiết lập
	logrus.Info("Starting server...")
	if err := server.ListenAndServe(":8080"); err != nil {
		logrus.Fatalf("Error in ListenAndServe: %v", err)
	}
}

// Hàm main
func main() {
	initGlobal()  // Khởi tạo các biến toàn cục
	main_thread() // Chạy server
}
