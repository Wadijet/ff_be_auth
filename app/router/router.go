package router

import (
	"meta_commerce/app/handler"
	"meta_commerce/app/middleware"
	"meta_commerce/config"

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
	r.POST(preV1+"/agents/checkin/{id}", middle.CheckUserAuth("Agent.Update", AgentHandler.CheckIn)) // Kiểm tra trạng thái online của tất cả đại lý

	// ====================================  ACCESSTOKEN API ========================================
	// Các API liên quan đến token
	AccessTokenHandler := handler.NewAccessTokenHandler(c, db)
	r.POST(preV1+"/access_tokens", middle.CheckUserAuth("AccessToken.Create", AccessTokenHandler.Create))        // Tạo token
	r.GET(preV1+"/access_tokens/{id}", middle.CheckUserAuth("AccessToken.Read", AccessTokenHandler.FindOne))     // Lấy token theo ID
	r.GET(preV1+"/access_tokens", middle.CheckUserAuth("AccessToken.Read", AccessTokenHandler.FindAll))          // Lấy tất cả token
	r.PUT(preV1+"/access_tokens/{id}", middle.CheckUserAuth("AccessToken.Update", AccessTokenHandler.Update))    // Cập nhật token theo ID
	r.DELETE(preV1+"/access_tokens/{id}", middle.CheckUserAuth("AccessToken.Delete", AccessTokenHandler.Delete)) // Xóa token theo ID

	// ====================================  FBPAGE API =============================================
	// Các API liên quan đến trang Facebook
	FbPageHandler := handler.NewFbPageHandler(c, db)
	r.POST(preV1+"/fb_pages", middle.CheckUserAuth("FbPage.Create", FbPageHandler.Create))                   // Tạo trang Facebook
	r.GET(preV1+"/fb_pages/{id}", middle.CheckUserAuth("FbPage.Read", FbPageHandler.FindOne))                // Lấy trang Facebook theo ID
	r.GET(preV1+"/fb_pages", middle.CheckUserAuth("FbPage.Read", FbPageHandler.FindAll))                     // Lấy tất cả trang Facebook
	r.POST(preV1+"/fb_pages/update_token", middle.CheckUserAuth("FbPage.Update", FbPageHandler.UpdateToken)) // Cập nhật token trang Facebook
	r.GET(preV1+"/fb_pages/pageId/{id}", middle.CheckUserAuth("FbPage.Read", FbPageHandler.FindOneByPageID)) // Lấy trang Facebook theo ID

	// ====================================  FBCONVERSATION API =====================================
	// Các API liên quan đến cuộc trò chuyện trên Facebook
	FbConversationHandler := handler.NewFbConversationHandler(c, db)
	r.POST(preV1+"/fb_conversations", middle.CheckUserAuth("FbConversation.Create", FbConversationHandler.Create))                     // Tạo cuộc trò chuyện
	r.GET(preV1+"/fb_conversations/{id}", middle.CheckUserAuth("FbConversation.Read", FbConversationHandler.FindOne))                  // Lấy cuộc trò chuyện theo ID
	r.GET(preV1+"/fb_conversations", middle.CheckUserAuth("FbConversation.Read", FbConversationHandler.FindAll))                       // Lấy tất cả cuộc trò chuyện
	r.GET(preV1+"/fb_conversations/newest", middle.CheckUserAuth("FbConversation.Read", FbConversationHandler.FindAllSortByApiUpdate)) // Lấy tất cả cuộc trò chuyện

	// ====================================  FBMESSAGE API ==========================================
	// Các API liên quan đến tin nhắn trên Facebook
	FbMessageHandler := handler.NewFbMessageHandler(c, db)
	r.POST(preV1+"/fb_messages", middle.CheckUserAuth("FbMessage.Create", FbMessageHandler.Create))    // Tạo tin nhắn
	r.GET(preV1+"/fb_messages/{id}", middle.CheckUserAuth("FbMessage.Read", FbMessageHandler.FindOne)) // Lấy tin nhắn theo ID
	r.GET(preV1+"/fb_messages", middle.CheckUserAuth("FbMessage.Read", FbMessageHandler.FindAll))      // Lấy tất cả tin nhắn

	// ====================================  FBPOST API =============================================
	// Các API liên quan đến bài viết trên Facebook
	FbPostHandler := handler.NewFbPostHandler(c, db)
	r.POST(preV1+"/fb_posts", middle.CheckUserAuth("FbPost.Create", FbPostHandler.Create))    // Tạo bài viết
	r.GET(preV1+"/fb_posts/{id}", middle.CheckUserAuth("FbPost.Read", FbPostHandler.FindOne)) // Lấy bài viết theo ID
	r.GET(preV1+"/fb_posts", middle.CheckUserAuth("FbPost.Read", FbPostHandler.FindAll))      // Lấy tất cả bài viết

	// ====================================  PCORDER API ============================================
	// Các API liên quan đến đơn hàng trên Pancake
	PcOrderHandler := handler.NewPcOrderHandler(c, db)
	r.POST(preV1+"/pc_orders", middle.CheckUserAuth("PcOrder.Create", PcOrderHandler.Create))    // Tạo đơn hàng
	r.GET(preV1+"/pc_orders/{id}", middle.CheckUserAuth("PcOrder.Read", PcOrderHandler.FindOne)) // Lấy đơn hàng theo ID
	r.GET(preV1+"/pc_orders", middle.CheckUserAuth("PcOrder.Read", PcOrderHandler.FindAll))      // Lấy tất cả đơn hàng

}
