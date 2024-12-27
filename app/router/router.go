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
		//r.GET(preV1+"/init/permissions", ApiInit.InitPermission) // Khởi tạo quyền
		//r.GET(preV1+"/init/roles", ApiInit.InitRole)             // Khởi tạo vai trò
		r.POST(preV1+"/init/setadmin/{id}", middle.CheckUserAuth("", ApiInit.SetAdministrator)) // Thiết lập admin
	}

	// ====================================  STATIC API ===============================================
	// Các API tĩnh
	ApiStatic := handler.NewStaticHandler()
	r.GET(preV1+"/static/test", ApiStatic.TestApi)                                     // API kiểm tra
	r.GET(preV1+"/static/system", middle.CheckUserAuth("", ApiStatic.GetSystemStatic)) // Lấy thông tin hệ thống
	r.GET(preV1+"/static/api", middle.CheckUserAuth("", ApiStatic.GetApiStatic))       // Lấy thông tin API

	// ====================================  PERMISSIONS API ========================================
	// Các API liên quan đến quyền
	ApiPermission := handler.NewPermissionHandler(c, db)
	r.GET(preV1+"/permissions/{id}", middle.CheckUserAuth("Permission.Read", ApiPermission.FindOneById)) // Lấy quyền theo ID
	r.GET(preV1+"/permissions", middle.CheckUserAuth("Permission.Read", ApiPermission.FindAll))          // Lấy tất cả quyền

	// ====================================  ROLES API =============================================
	// Các API liên quan đến vai trò
	ApiRole := handler.NewRoleHandler(c, db)
	r.POST(preV1+"/roles", middle.CheckUserAuth("Role.Create", ApiRole.Create))               // Tạo vai trò
	r.GET(preV1+"/roles/{id}", middle.CheckUserAuth("Role.Read", ApiRole.FindOneById))        // Lấy vai trò theo ID
	r.GET(preV1+"/roles", middle.CheckUserAuth("Role.Read", ApiRole.FindAll))                 // Lấy tất cả vai trò
	r.PUT(preV1+"/roles/{id}", middle.CheckUserAuth("Role.Update", ApiRole.UpdateOneById))    // Cập nhật vai trò theo ID
	r.DELETE(preV1+"/roles/{id}", middle.CheckUserAuth("Role.Delete", ApiRole.DeleteOneById)) // Xóa vai trò theo ID

	// ====================================  ROLE PERMISSIONS API ====================================
	// Các API liên quan đến quyền của vai trò
	ApiRolePermission := handler.NewRolePermissionHandler(c, db)
	r.POST(preV1+"/role_permissions", middle.CheckUserAuth("RolePermission.Create", ApiRolePermission.Create))        // Tạo quyền cho vai trò
	r.DELETE(preV1+"/role_permissions/{id}", middle.CheckUserAuth("RolePermission.Delete", ApiRolePermission.Delete)) // Xóa quyền của vai trò

	// ====================================  ADMIN API =============================================
	// Các API dành cho admin
	ApiAdmin := handler.NewAdminHandler(c, db)
	r.POST(preV1+"/admin/set_role", middle.CheckUserAuth("Admin.Set_role", ApiAdmin.SetRole))           // Thiết lập vai trò cho người dùng
	r.POST(preV1+"/admin/block_user", middle.CheckUserAuth("Admin.Block_user", ApiAdmin.BlockUser))     // Khóa người dùng
	r.POST(preV1+"/admin/unblock_user", middle.CheckUserAuth("Admin.Block_user", ApiAdmin.UnBlockUser)) // Mở khóa người dùng

	// ====================================  USERS API =============================================
	// Các API liên quan đến người dùng
	ApiUser := handler.NewUserHandler(c, db)
	r.POST(preV1+"/users/register", ApiUser.Registry)                                        // Đăng ký người dùng
	r.POST(preV1+"/users/login", ApiUser.Login)                                              // Đăng nhập người dùng
	r.POST(preV1+"/users/logout", middle.CheckUserAuth("", ApiUser.Logout))                  // Đăng xuất người dùng
	r.GET(preV1+"/users/me", middle.CheckUserAuth("", ApiUser.GetMyInfo))                    // Lấy thông tin cá nhân
	r.GET(preV1+"/users/roles", middle.CheckUserAuth("", ApiUser.GetMyRoles))                // Lấy vai trò của người dùng
	r.POST(preV1+"/users/change_password", middle.CheckUserAuth("", ApiUser.ChangePassword)) // Đổi mật khẩu
	r.POST(preV1+"/users/change_info", middle.CheckUserAuth("", ApiUser.ChangeInfo))         // Đổi thông tin cá nhân
	// TODO: Bổ sung check quyền khi chạy thật
	r.GET(preV1+"/users/{id}", middle.CheckUserAuth("User.Read", ApiUser.FindOneById))  // Lấy tất cả người dùng với bộ lọc
	r.GET(preV1+"/users/count", middle.CheckUserAuth("User.Read", ApiUser.Count))       // Lấy tất cả người dùng với bộ lọc
	r.GET(preV1+"/users", middle.CheckUserAuth("User.Read", ApiUser.FindAllWithFilter)) // Lấy tất cả người dùng với bộ lọc
}
