package main

import (
	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	validator "gopkg.in/go-playground/validator.v9"

	api "atk-go-server/app/router"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/database"
	"atk-go-server/global"
)

// Hàm khởi tạo các biến toàn cục
func initGlobal() {
	initColNames()  // Khởi tạo tên các collection trong database
	initValidator() // Khởi tạo validator
	initConfig()    // Khởi tạo cấu hình server
	initDatabase()  // Khởi tạo kết nối database
}

// Hàm khởi tạo tên các collection trong database
func initColNames() {
	global.MongoDB_ColNames.Users = "users"
	global.MongoDB_ColNames.Permissions = "permissions"
	global.MongoDB_ColNames.Roles = "roles"
	global.MongoDB_ColNames.RolePermissions = "role_permissions"
	global.MongoDB_ColNames.UserRoles = "user_roles"

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
func initDatabase() {
	var err error
	global.MongoDB_Session, err = database.GetInstance(global.MongoDB_ServerConfig)
	if err != nil {
		logrus.Fatalf("Failed to get database instance: %v", err) // Ghi log lỗi nếu kết nối database thất bại
	}
	logrus.Info("Connected to MongoDB") // Ghi log thông báo đã kết nối database thành công
	database.EnsureDatabaseAndCollections(global.MongoDB_Session)

	global.MySQL_Session, err = database.GetMySQLInstance(global.MongoDB_ServerConfig)
	if err != nil {
		logrus.Fatalf("Failed to get MySQL instance: %v", err)
	}
	logrus.Info("Connected to MySQL")
}

// Hàm xử lý panic
func panicHandler(ctx *fasthttp.RequestCtx, data interface{}) {
	logrus.Errorf("Panic occurred: %v", data)                     // Ghi log lỗi khi xảy ra panic
	utility.JSON(ctx, utility.Payload(false, data, "Lỗi panic!")) // Trả về JSON thông báo lỗi panic
}

// Hàm chính để chạy server
func main_thread() {
	// Khởi tạo router
	r := router.New()
	api.InitRounters(r, global.MongoDB_ServerConfig, global.MongoDB_Session) // Khởi tạo các route cho API
	r.PanicHandler = panicHandler                                            // Đặt hàm xử lý panic

	// Chạy server
	logrus.Info("Starting server...") // Ghi log thông báo bắt đầu chạy server
	if err := fasthttp.ListenAndServe(":8080", r.Handler); err != nil {
		logrus.Fatalf("Error in ListenAndServe: %v", err) // Ghi log lỗi nếu server không thể chạy
	}
}

// Hàm main
func main() {
	initGlobal()  // Khởi tạo các biến toàn cục
	main_thread() // Chạy server
}

// test
