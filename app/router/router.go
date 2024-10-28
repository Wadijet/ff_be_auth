package router

import (
	"atk-go-server/app/handler"
	"atk-go-server/app/middleware"
	"atk-go-server/config"

	"github.com/fasthttp/router"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	preBase = "/api"
	preV1   = preBase + "/v1"
)

// InitRounters khởi tạo các route cho ứng dụng
func InitRounters(r *router.Router, c *config.Configuration, db *mongo.Client) {
	middle := middleware.NewJwtToken(c, db)

	// ====================================  INIT API ===============================================
	// Các API khởi tạo hệ thống
	if c.InitMode == true {
		ApiInit := handler.NewInitHandler(c, db)
		r.GET(preV1+"/init/permissions", ApiInit.InitPermission) // Khởi tạo quyền
		r.GET(preV1+"/init/roles", ApiInit.InitRole)             // Khởi tạo vai trò
		r.POST(preV1+"/init/setadmin", ApiInit.SetAdmin)         // Thiết lập admin
	}

	// ====================================  STATIC API ===============================================
	// Các API tĩnh
	ApiStatic := handler.NewStaticHandler()
	r.GET(preV1+"/static/test", ApiStatic.TestApi)                                      // API kiểm tra
	r.GET(preV1+"/static/system", middle.CheckUserAuth(nil, ApiStatic.GetSystemStatic)) // Lấy thông tin hệ thống
	r.GET(preV1+"/static/api", middle.CheckUserAuth(nil, ApiStatic.GetApiStatic))       // Lấy thông tin API

	// ====================================  PERMISSIONS API ========================================
	// Các API liên quan đến quyền
	ApiPermission := handler.NewPermissionHandler(c, db)
	r.POST(preV1+"/permissions", middle.CheckUserAuth([]string{"permission.create"}, ApiPermission.Create))               // Tạo quyền
	r.GET(preV1+"/permissions/{id}", middle.CheckUserAuth([]string{"permission.readall"}, ApiPermission.FindOneById))     // Lấy quyền theo ID
	r.GET(preV1+"/permissions", middle.CheckUserAuth([]string{"permission.readall"}, ApiPermission.FindAll))              // Lấy tất cả quyền
	r.PUT(preV1+"/permissions/{id}", middle.CheckUserAuth([]string{"permission.update"}, ApiPermission.UpdateOneById))    // Cập nhật quyền theo ID
	r.DELETE(preV1+"/permissions/{id}", middle.CheckUserAuth([]string{"permission.delete"}, ApiPermission.DeleteOneById)) // Xóa quyền theo ID

	// ====================================  ROLES API =============================================
	// Các API liên quan đến vai trò
	ApiRole := handler.NewRoleHandler(c, db)
	//r.POST(preV1+"/roles", middle.CheckAuth(ApiRole.Create))
	r.GET(preV1+"/roles/{id}", middle.CheckUserAuth([]string{"role.readall"}, ApiRole.FindOneById)) // Lấy vai trò theo ID
	r.GET(preV1+"/roles", middle.CheckUserAuth([]string{"role.readall"}, ApiRole.FindAll))          // Lấy tất cả vai trò
	//r.PUT(preV1+"/roles/{id}", middle.CheckAuth(ApiRole.UpdateOneById))
	//r.DELETE(preV1+"/roles/{id}", middle.CheckAuth(ApiRole.DeleteOneById))

	// ====================================  ADMIN API =============================================
	// Các API dành cho admin
	ApiAdmin := handler.NewAdminHandler(c, db)
	r.POST(preV1+"/admin/set_role", middle.CheckUserAuth([]string{"admin.set_role"}, ApiAdmin.SetRole))             // Thiết lập vai trò cho người dùng
	r.POST(preV1+"/admin/block_user", middle.CheckUserAuth([]string{"admin.block_user"}, ApiAdmin.BlockUser))       // Khóa người dùng
	r.POST(preV1+"/admin/unblock_user", middle.CheckUserAuth([]string{"admin.unblock_user"}, ApiAdmin.UnBlockUser)) // Mở khóa người dùng

	// ====================================  USERS API =============================================
	// Các API liên quan đến người dùng
	ApiUser := handler.NewUserHandler(c, db)
	r.POST(preV1+"/users/register", ApiUser.Registry)                                                 // Đăng ký người dùng
	r.POST(preV1+"/users/login", ApiUser.Login)                                                       // Đăng nhập người dùng
	r.POST(preV1+"/users/logout", middle.CheckUserAuth(nil, ApiUser.Logout))                          // Đăng xuất người dùng
	r.GET(preV1+"/users/me", middle.CheckUserAuth(nil, ApiUser.GetMyInfo))                            // Lấy thông tin cá nhân
	r.GET(preV1+"/users", middle.CheckUserAuth([]string{"user.read_all"}, ApiUser.FindAllWithFilter)) // Lấy tất cả người dùng với bộ lọc
	r.POST(preV1+"/users/change_password", middle.CheckUserAuth(nil, ApiUser.ChangePassword))         // Đổi mật khẩu
	r.POST(preV1+"/users/change_info", middle.CheckUserAuth(nil, ApiUser.ChangeInfo))                 // Đổi thông tin cá nhân

	r.POST(preV1+"/users/check_token", middle.CheckMtServiceAuth(nil, ApiUser.CheckToken)) // Kiểm tra token

	// ====================================  MUTI SERVICES API =============================================
	// Các API liên quan đến dịch vụ MT
	ApiMtService := handler.NewMtServiceHandler(c, db)
	r.POST(preV1+"/mtservices/register", ApiMtService.Registry)                                                           // Đăng ký dịch vụ MT
	r.POST(preV1+"/mtservices/login", ApiMtService.Login)                                                                 // Đăng nhập dịch vụ MT
	r.GET(preV1+"/mtservices/me", middle.CheckMtServiceAuth(nil, ApiMtService.GetMyInfo))                                 // Lấy thông tin cá nhân dịch vụ MT
	r.GET(preV1+"/mtservices", middle.CheckMtServiceAuth([]string{"mtservice.read_all"}, ApiMtService.FindAllWithFilter)) // Lấy tất cả dịch vụ MT với bộ lọc
	r.POST(preV1+"/mtservices/change_password", middle.CheckMtServiceAuth(nil, ApiMtService.ChangePassword))              // Đổi mật khẩu dịch vụ MT
	r.POST(preV1+"/mtservices/change_info", middle.CheckMtServiceAuth(nil, ApiMtService.ChangeInfo))                      // Đổi thông tin dịch vụ MT
	r.POST(preV1+"/mtservices/check_token", middle.CheckMtServiceAuth(nil, ApiMtService.CheckToken))                      // Kiểm tra token dịch vụ MT
}
