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
	StaticHandler := handler.NewStaticHandler()
	r.GET(preV1+"/static/test", StaticHandler.TestApi)                                     // API kiểm tra
	r.GET(preV1+"/static/system", middle.CheckUserAuth("", StaticHandler.GetSystemStatic)) // Lấy thông tin hệ thống
	r.GET(preV1+"/static/api", middle.CheckUserAuth("", StaticHandler.GetApiStatic))       // Lấy thông tin API

	// ====================================  PERMISSIONS API ========================================
	// Các API liên quan đến quyền
	PermissionHandler := handler.NewPermissionHandler(c, db)
	r.GET(preV1+"/permissions/{id}", middle.CheckUserAuth("Permission.Read", PermissionHandler.FindOneById)) // Lấy quyền theo ID
	r.GET(preV1+"/permissions", middle.CheckUserAuth("Permission.Read", PermissionHandler.FindAll))          // Lấy tất cả quyền

	// ====================================  ROLES API =============================================
	// Các API liên quan đến vai trò
	RoleHandler := handler.NewRoleHandler(c, db)
	r.POST(preV1+"/roles", middle.CheckUserAuth("Role.Create", RoleHandler.Create))               // Tạo vai trò
	r.GET(preV1+"/roles/{id}", middle.CheckUserAuth("Role.Read", RoleHandler.FindOneById))        // Lấy vai trò theo ID
	r.GET(preV1+"/roles", middle.CheckUserAuth("Role.Read", RoleHandler.FindAll))                 // Lấy tất cả vai trò
	r.PUT(preV1+"/roles/{id}", middle.CheckUserAuth("Role.Update", RoleHandler.UpdateOneById))    // Cập nhật vai trò theo ID
	r.DELETE(preV1+"/roles/{id}", middle.CheckUserAuth("Role.Delete", RoleHandler.DeleteOneById)) // Xóa vai trò theo ID

	// ====================================  ROLE PERMISSIONS API ====================================
	// Các API liên quan đến quyền của vai trò
	RolePermissionHandler := handler.NewRolePermissionHandler(c, db)
	r.POST(preV1+"/role_permissions", middle.CheckUserAuth("RolePermission.Create", RolePermissionHandler.Create))        // Tạo quyền cho vai trò
	r.DELETE(preV1+"/role_permissions/{id}", middle.CheckUserAuth("RolePermission.Delete", RolePermissionHandler.Delete)) // Xóa quyền của vai trò

	// ====================================  USER ROLES API ========================================
	// Các API liên quan đến vai trò của người dùng
	UserRoleHanlder := handler.NewUserRoleHandler(c, db)
	r.POST(preV1+"/user_roles", middle.CheckUserAuth("UserRole.Create", UserRoleHanlder.Create))        // Tạo vai trò cho người dùng
	r.DELETE(preV1+"/user_roles/{id}", middle.CheckUserAuth("UserRole.Delete", UserRoleHanlder.Delete)) // Xóa vai trò của người dùng

	// ====================================  ADMIN API =============================================
	// Các API dành cho admin
	AdminHandler := handler.NewAdminHandler(c, db)
	r.POST(preV1+"/admin/set_role", middle.CheckUserAuth("Admin.Set_role", AdminHandler.SetRole))           // Thiết lập vai trò cho người dùng
	r.POST(preV1+"/admin/block_user", middle.CheckUserAuth("Admin.Block_user", AdminHandler.BlockUser))     // Khóa người dùng
	r.POST(preV1+"/admin/unblock_user", middle.CheckUserAuth("Admin.Block_user", AdminHandler.UnBlockUser)) // Mở khóa người dùng

	// ====================================  USERS API =============================================
	// Các API liên quan đến người dùng
	UserHandler := handler.NewUserHandler(c, db)
	r.POST(preV1+"/users/register", UserHandler.Registry)                                        // Đăng ký người dùng
	r.POST(preV1+"/users/login", UserHandler.Login)                                              // Đăng nhập người dùng
	r.POST(preV1+"/users/logout", middle.CheckUserAuth("", UserHandler.Logout))                  // Đăng xuất người dùng
	r.GET(preV1+"/users/me", middle.CheckUserAuth("", UserHandler.GetMyInfo))                    // Lấy thông tin cá nhân
	r.GET(preV1+"/users/roles", middle.CheckUserAuth("", UserHandler.GetMyRoles))                // Lấy vai trò của người dùng
	r.POST(preV1+"/users/change_password", middle.CheckUserAuth("", UserHandler.ChangePassword)) // Đổi mật khẩu
	r.POST(preV1+"/users/change_info", middle.CheckUserAuth("", UserHandler.ChangeInfo))         // Đổi thông tin cá nhân
	// TODO: Bổ sung check quyền khi chạy thật
	r.GET(preV1+"/users/{id}", middle.CheckUserAuth("User.Read", UserHandler.FindOneById))  // Lấy tất cả người dùng với bộ lọc
	r.GET(preV1+"/users", middle.CheckUserAuth("User.Read", UserHandler.FindAllWithFilter)) // Lấy tất cả người dùng với bộ lọc

	// ====================================  AGENTS API =============================================
	// Các API liên quan đến đại lý
	AgentHandler := handler.NewAgentHandler(c, db)
	r.POST(preV1+"/agents", middle.CheckUserAuth("Agent.Create", AgentHandler.Create))               // Tạo đại lý
	r.GET(preV1+"/agents/{id}", middle.CheckUserAuth("Agent.Read", AgentHandler.FindOneById))        // Lấy đại lý theo ID
	r.GET(preV1+"/agents", middle.CheckUserAuth("Agent.Read", AgentHandler.FindAll))                 // Lấy tất cả đại lý
	r.PUT(preV1+"/agents/{id}", middle.CheckUserAuth("Agent.Update", AgentHandler.UpdateOneById))    // Cập nhật đại lý theo ID
	r.DELETE(preV1+"/agents/{id}", middle.CheckUserAuth("Agent.Delete", AgentHandler.DeleteOneById)) // Xóa đại lý theo ID
}
