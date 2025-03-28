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
		r.POST(preV1+"/init/setadmin/{id}", middle.CheckUserAuth("User.SetRole", ApiInit.SetAdministrator)) // Thiết lập admin
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
	r.GET(preV1+"/permissions/{id}", middle.CheckUserAuth("Permission.Read", PermissionHandler.FindOneById))    // Lấy quyền theo ID
	r.GET(preV1+"/permissions", middle.CheckUserAuth("Permission.Read", PermissionHandler.Find))                // Lấy tất cả quyền
	r.POST(preV1+"/permissions", middle.CheckUserAuth("Permission.Insert", PermissionHandler.InsertOne))        // Tạo quyền
	r.PUT(preV1+"/permissions/{id}", middle.CheckUserAuth("Permission.Update", PermissionHandler.UpdateOne))    // Cập nhật quyền
	r.DELETE(preV1+"/permissions/{id}", middle.CheckUserAuth("Permission.Delete", PermissionHandler.DeleteOne)) // Xóa quyền

	// ====================================  ROLES API =============================================
	// Các API liên quan đến vai trò
	RoleHandler := handler.NewRoleHandler(c, db)
	r.POST(preV1+"/roles", middle.CheckUserAuth("Role.Insert", RoleHandler.InsertOne))        // Tạo vai trò
	r.GET(preV1+"/roles/{id}", middle.CheckUserAuth("Role.Read", RoleHandler.FindOneById))    // Lấy vai trò theo ID
	r.GET(preV1+"/roles", middle.CheckUserAuth("Role.Read", RoleHandler.Find))                // Lấy tất cả vai trò
	r.PUT(preV1+"/roles/{id}", middle.CheckUserAuth("Role.Update", RoleHandler.UpdateOne))    // Cập nhật vai trò
	r.DELETE(preV1+"/roles/{id}", middle.CheckUserAuth("Role.Delete", RoleHandler.DeleteOne)) // Xóa vai trò

	// ====================================  ROLE PERMISSIONS API ====================================
	// Các API liên quan đến quyền của vai trò
	RolePermissionHandler := handler.NewRolePermissionHandler(c, db)
	r.POST(preV1+"/role_permissions", middle.CheckUserAuth("RolePermission.Insert", RolePermissionHandler.InsertOne))        // Tạo quyền cho vai trò
	r.GET(preV1+"/role_permissions", middle.CheckUserAuth("RolePermission.Read", RolePermissionHandler.Find))                // Lấy danh sách quyền của vai trò
	r.PUT(preV1+"/role_permissions/{id}", middle.CheckUserAuth("RolePermission.Update", RolePermissionHandler.UpdateOne))    // Cập nhật quyền của vai trò
	r.DELETE(preV1+"/role_permissions/{id}", middle.CheckUserAuth("RolePermission.Delete", RolePermissionHandler.DeleteOne)) // Xóa quyền của vai trò

	// ====================================  USER ROLES API ========================================
	// Các API liên quan đến vai trò của người dùng
	UserRoleHanlder := handler.NewUserRoleHandler(c, db)
	r.POST(preV1+"/user_roles", middle.CheckUserAuth("UserRole.Insert", UserRoleHanlder.InsertOne))        // Tạo vai trò cho người dùng
	r.GET(preV1+"/user_roles", middle.CheckUserAuth("UserRole.Read", UserRoleHanlder.Find))                // Lấy danh sách vai trò của người dùng
	r.PUT(preV1+"/user_roles/{id}", middle.CheckUserAuth("UserRole.Update", UserRoleHanlder.UpdateOne))    // Cập nhật vai trò của người dùng
	r.DELETE(preV1+"/user_roles/{id}", middle.CheckUserAuth("UserRole.Delete", UserRoleHanlder.DeleteOne)) // Xóa vai trò của người dùng

	// ====================================  ADMIN API =============================================
	// Các API dành cho admin
	AdminHandler := handler.NewAdminHandler(c, db)
	r.POST(preV1+"/admin/set_role", middle.CheckUserAuth("User.SetRole", AdminHandler.SetRole))       // Thiết lập vai trò cho người dùng
	r.POST(preV1+"/admin/block_user", middle.CheckUserAuth("User.Block", AdminHandler.BlockUser))     // Khóa người dùng
	r.POST(preV1+"/admin/unblock_user", middle.CheckUserAuth("User.Block", AdminHandler.UnBlockUser)) // Mở khóa người dùng

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
	r.GET(preV1+"/users/{id}", middle.CheckUserAuth("User.Read", UserHandler.FindOneById))       // Lấy thông tin người dùng theo ID
	r.GET(preV1+"/users", middle.CheckUserAuth("User.Read", UserHandler.Find))                   // Lấy danh sách người dùng
	r.POST(preV1+"/users", middle.CheckUserAuth("User.Insert", UserHandler.InsertOne))           // Tạo người dùng
	r.PUT(preV1+"/users/{id}", middle.CheckUserAuth("User.Update", UserHandler.UpdateOne))       // Cập nhật thông tin người dùng
	r.DELETE(preV1+"/users/{id}", middle.CheckUserAuth("User.Delete", UserHandler.DeleteOne))    // Xóa người dùng

	// ====================================  AGENTS API =============================================
	// Các API liên quan đến đại lý
	AgentHandler := handler.NewAgentHandler(c, db)
	r.POST(preV1+"/agents", middle.CheckUserAuth("Agent.Insert", AgentHandler.InsertOne))             // Tạo đại lý
	r.GET(preV1+"/agents/{id}", middle.CheckUserAuth("Agent.Read", AgentHandler.FindOneById))         // Lấy đại lý theo ID
	r.GET(preV1+"/agents", middle.CheckUserAuth("Agent.Read", AgentHandler.Find))                     // Lấy tất cả đại lý
	r.PUT(preV1+"/agents/{id}", middle.CheckUserAuth("Agent.Update", AgentHandler.UpdateOne))         // Cập nhật đại lý
	r.DELETE(preV1+"/agents/{id}", middle.CheckUserAuth("Agent.Delete", AgentHandler.DeleteOne))      // Xóa đại lý
	r.POST(preV1+"/agents/checkin/{id}", middle.CheckUserAuth("Agent.CheckIn", AgentHandler.CheckIn)) // Kiểm tra trạng thái online của đại lý

	// ====================================  ACCESSTOKEN API ========================================
	// Các API liên quan đến token
	AccessTokenHandler := handler.NewAccessTokenHandler(c, db)
	r.POST(preV1+"/access_tokens", middle.CheckUserAuth("AccessToken.Insert", AccessTokenHandler.InsertOne))        // Tạo token
	r.GET(preV1+"/access_tokens/{id}", middle.CheckUserAuth("AccessToken.Read", AccessTokenHandler.FindOne))        // Lấy token theo ID
	r.GET(preV1+"/access_tokens", middle.CheckUserAuth("AccessToken.Read", AccessTokenHandler.Find))                // Lấy tất cả token
	r.PUT(preV1+"/access_tokens/{id}", middle.CheckUserAuth("AccessToken.Update", AccessTokenHandler.UpdateOne))    // Cập nhật token
	r.DELETE(preV1+"/access_tokens/{id}", middle.CheckUserAuth("AccessToken.Delete", AccessTokenHandler.DeleteOne)) // Xóa token

	// ====================================  FBPAGE API =============================================
	// Các API liên quan đến trang Facebook
	FbPageHandler := handler.NewFbPageHandler(c, db)
	r.POST(preV1+"/fb_pages", middle.CheckUserAuth("FbPage.Insert", FbPageHandler.InsertOne))                     // Tạo trang Facebook
	r.GET(preV1+"/fb_pages/{id}", middle.CheckUserAuth("FbPage.Read", FbPageHandler.FindOne))                     // Lấy trang Facebook theo ID
	r.GET(preV1+"/fb_pages", middle.CheckUserAuth("FbPage.Read", FbPageHandler.Find))                             // Lấy tất cả trang Facebook
	r.PUT(preV1+"/fb_pages/{id}", middle.CheckUserAuth("FbPage.Update", FbPageHandler.UpdateOne))                 // Cập nhật trang Facebook
	r.DELETE(preV1+"/fb_pages/{id}", middle.CheckUserAuth("FbPage.Delete", FbPageHandler.DeleteOne))              // Xóa trang Facebook
	r.POST(preV1+"/fb_pages/update_token", middle.CheckUserAuth("FbPage.UpdateToken", FbPageHandler.UpdateToken)) // Cập nhật token trang Facebook
	r.GET(preV1+"/fb_pages/pageId/{id}", middle.CheckUserAuth("FbPage.Read", FbPageHandler.FindOneByPageID))      // Lấy trang Facebook theo PageID

	// ====================================  FBCONVERSATION API =====================================
	// Các API liên quan đến cuộc trò chuyện trên Facebook
	FbConversationHandler := handler.NewFbConversationHandler(c, db)
	r.POST(preV1+"/fb_conversations", middle.CheckUserAuth("FbConversation.Insert", FbConversationHandler.InsertOne))                  // Tạo cuộc trò chuyện
	r.GET(preV1+"/fb_conversations/{id}", middle.CheckUserAuth("FbConversation.Read", FbConversationHandler.FindOne))                  // Lấy cuộc trò chuyện theo ID
	r.GET(preV1+"/fb_conversations", middle.CheckUserAuth("FbConversation.Read", FbConversationHandler.Find))                          // Lấy tất cả cuộc trò chuyện
	r.PUT(preV1+"/fb_conversations/{id}", middle.CheckUserAuth("FbConversation.Update", FbConversationHandler.UpdateOne))              // Cập nhật cuộc trò chuyện
	r.DELETE(preV1+"/fb_conversations/{id}", middle.CheckUserAuth("FbConversation.Delete", FbConversationHandler.DeleteOne))           // Xóa cuộc trò chuyện
	r.GET(preV1+"/fb_conversations/newest", middle.CheckUserAuth("FbConversation.Read", FbConversationHandler.FindAllSortByApiUpdate)) // Lấy tất cả cuộc trò chuyện mới nhất

	// ====================================  FBMESSAGE API ==========================================
	// Các API liên quan đến tin nhắn trên Facebook
	FbMessageHandler := handler.NewFbMessageHandler(c, db)
	r.POST(preV1+"/fb_messages", middle.CheckUserAuth("FbMessage.Insert", FbMessageHandler.InsertOne))        // Tạo tin nhắn
	r.GET(preV1+"/fb_messages/{id}", middle.CheckUserAuth("FbMessage.Read", FbMessageHandler.FindOne))        // Lấy tin nhắn theo ID
	r.GET(preV1+"/fb_messages", middle.CheckUserAuth("FbMessage.Read", FbMessageHandler.Find))                // Lấy tất cả tin nhắn
	r.PUT(preV1+"/fb_messages/{id}", middle.CheckUserAuth("FbMessage.Update", FbMessageHandler.UpdateOne))    // Cập nhật tin nhắn
	r.DELETE(preV1+"/fb_messages/{id}", middle.CheckUserAuth("FbMessage.Delete", FbMessageHandler.DeleteOne)) // Xóa tin nhắn

	// ====================================  FBPOST API =============================================
	// Các API liên quan đến bài viết trên Facebook
	FbPostHandler := handler.NewFbPostHandler(c, db)
	r.POST(preV1+"/fb_posts", middle.CheckUserAuth("FbPost.Insert", FbPostHandler.InsertOne))        // Tạo bài viết
	r.GET(preV1+"/fb_posts/{id}", middle.CheckUserAuth("FbPost.Read", FbPostHandler.FindOne))        // Lấy bài viết theo ID
	r.GET(preV1+"/fb_posts", middle.CheckUserAuth("FbPost.Read", FbPostHandler.Find))                // Lấy tất cả bài viết
	r.PUT(preV1+"/fb_posts/{id}", middle.CheckUserAuth("FbPost.Update", FbPostHandler.UpdateOne))    // Cập nhật bài viết
	r.DELETE(preV1+"/fb_posts/{id}", middle.CheckUserAuth("FbPost.Delete", FbPostHandler.DeleteOne)) // Xóa bài viết

	// ====================================  PCORDER API ============================================
	// Các API liên quan đến đơn hàng trên Pancake
	PcOrderHandler := handler.NewPcOrderHandler(c, db)
	r.POST(preV1+"/pc_orders", middle.CheckUserAuth("PcOrder.Insert", PcOrderHandler.InsertOne))        // Tạo đơn hàng
	r.GET(preV1+"/pc_orders/{id}", middle.CheckUserAuth("PcOrder.Read", PcOrderHandler.FindOne))        // Lấy đơn hàng theo ID
	r.GET(preV1+"/pc_orders", middle.CheckUserAuth("PcOrder.Read", PcOrderHandler.Find))                // Lấy tất cả đơn hàng
	r.PUT(preV1+"/pc_orders/{id}", middle.CheckUserAuth("PcOrder.Update", PcOrderHandler.UpdateOne))    // Cập nhật đơn hàng
	r.DELETE(preV1+"/pc_orders/{id}", middle.CheckUserAuth("PcOrder.Delete", PcOrderHandler.DeleteOne)) // Xóa đơn hàng
}
