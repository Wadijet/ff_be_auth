package router

import (
	"meta_commerce/app/handler"
	"meta_commerce/app/middleware"
	"meta_commerce/config"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	preBase = "/api"
	preV1   = preBase + "/v1"
)

// registerCRUDRoutes đăng ký các route CRUD cơ bản
func registerCRUDRoutes(r *router.Router, prefix string, handler interface{}, middle *middleware.JwtToken) {
	// Định nghĩa các route CRUD cơ bản
	r.POST(prefix, middle.CheckUserAuth(prefix+".Insert", handler.(interface {
		InsertOne(ctx *fasthttp.RequestCtx)
	}).InsertOne))
	r.POST(prefix+"/many", middle.CheckUserAuth(prefix+".Insert", handler.(interface {
		InsertMany(ctx *fasthttp.RequestCtx)
	}).InsertMany))

	r.GET(prefix+"/{id}", middle.CheckUserAuth(prefix+".Read", handler.(interface {
		FindOne(ctx *fasthttp.RequestCtx)
	}).FindOne))
	r.GET(prefix+"/by_id/{id}", middle.CheckUserAuth(prefix+".Read", handler.(interface {
		FindOneById(ctx *fasthttp.RequestCtx)
	}).FindOneById))
	r.GET(prefix+"/by_ids", middle.CheckUserAuth(prefix+".Read", handler.(interface {
		FindManyByIds(ctx *fasthttp.RequestCtx)
	}).FindManyByIds))
	r.GET(prefix+"/pagination", middle.CheckUserAuth(prefix+".Read", handler.(interface {
		FindWithPagination(ctx *fasthttp.RequestCtx)
	}).FindWithPagination))
	r.GET(prefix, middle.CheckUserAuth(prefix+".Read", handler.(interface {
		Find(ctx *fasthttp.RequestCtx)
	}).Find))

	r.PUT(prefix+"/{id}", middle.CheckUserAuth(prefix+".Update", handler.(interface {
		UpdateOne(ctx *fasthttp.RequestCtx)
	}).UpdateOne))
	r.PUT(prefix+"/many", middle.CheckUserAuth(prefix+".Update", handler.(interface {
		UpdateMany(ctx *fasthttp.RequestCtx)
	}).UpdateMany))
	r.PUT(prefix+"/by_id/{id}", middle.CheckUserAuth(prefix+".Update", handler.(interface {
		UpdateById(ctx *fasthttp.RequestCtx)
	}).UpdateById))

	r.DELETE(prefix+"/{id}", middle.CheckUserAuth(prefix+".Delete", handler.(interface {
		DeleteOne(ctx *fasthttp.RequestCtx)
	}).DeleteOne))
	r.DELETE(prefix+"/many", middle.CheckUserAuth(prefix+".Delete", handler.(interface {
		DeleteMany(ctx *fasthttp.RequestCtx)
	}).DeleteMany))
	r.DELETE(prefix+"/by_id/{id}", middle.CheckUserAuth(prefix+".Delete", handler.(interface {
		DeleteById(ctx *fasthttp.RequestCtx)
	}).DeleteById))

	// Các route atomic và tiện ích
	r.PUT(prefix+"/find_and_update", middle.CheckUserAuth(prefix+".Update", handler.(interface {
		FindOneAndUpdate(ctx *fasthttp.RequestCtx)
	}).FindOneAndUpdate))
	r.DELETE(prefix+"/find_and_delete", middle.CheckUserAuth(prefix+".Delete", handler.(interface {
		FindOneAndDelete(ctx *fasthttp.RequestCtx)
	}).FindOneAndDelete))

	r.GET(prefix+"/count", middle.CheckUserAuth(prefix+".Read", handler.(interface {
		CountDocuments(ctx *fasthttp.RequestCtx)
	}).CountDocuments))
	r.GET(prefix+"/distinct/{field}", middle.CheckUserAuth(prefix+".Read", handler.(interface {
		Distinct(ctx *fasthttp.RequestCtx)
	}).Distinct))

	r.POST(prefix+"/upsert", middle.CheckUserAuth(prefix+".Upsert", handler.(interface {
		Upsert(ctx *fasthttp.RequestCtx)
	}).Upsert))
	r.POST(prefix+"/upsert/many", middle.CheckUserAuth(prefix+".Upsert", handler.(interface {
		UpsertMany(ctx *fasthttp.RequestCtx)
	}).UpsertMany))

	r.GET(prefix+"/exists", middle.CheckUserAuth(prefix+".Read", handler.(interface {
		DocumentExists(ctx *fasthttp.RequestCtx)
	}).DocumentExists))
}

// registerAuthRoutes đăng ký các route liên quan đến xác thực
func registerAuthRoutes(r *router.Router, c *config.Configuration, db *mongo.Client, middle *middleware.JwtToken) {
	// ====================================  INIT API ===============================================
	if c.InitMode == true {
		ApiInit := handler.NewInitHandler()
		r.POST(preV1+"/init/setadmin/{id}", ApiInit.SetAdministrator)
	}

	// ====================================  PERMISSIONS API ========================================
	PermissionHandler := handler.NewPermissionHandler()
	registerCRUDRoutes(r, preV1+"/permissions", PermissionHandler, middle)

	// ====================================  ROLES API =============================================
	RoleHandler := handler.NewRoleHandler()
	registerCRUDRoutes(r, preV1+"/roles", RoleHandler, middle)

	// ====================================  ROLE PERMISSIONS API ====================================
	RolePermissionHandler := handler.NewRolePermissionHandler()
	registerCRUDRoutes(r, preV1+"/role_permissions", RolePermissionHandler, middle)

	// ====================================  USER ROLES API ========================================
	UserRoleHanlder := handler.NewUserRoleHandler()
	registerCRUDRoutes(r, preV1+"/user_roles", UserRoleHanlder, middle)

	// ====================================  ADMIN API =============================================
	AdminHandler := handler.NewAdminHandler()
	r.POST(preV1+"/admin/set_role", middle.CheckUserAuth("User.SetRole", AdminHandler.SetRole))
	r.POST(preV1+"/admin/block_user", middle.CheckUserAuth("User.Block", AdminHandler.BlockUser))
	r.POST(preV1+"/admin/unblock_user", middle.CheckUserAuth("User.Block", AdminHandler.UnBlockUser))

	// ====================================  USERS API =============================================
	UserHandler := handler.NewUserHandler()
	r.POST(preV1+"/users/register", UserHandler.Registry)
	r.POST(preV1+"/users/login", UserHandler.Login)
	r.POST(preV1+"/users/logout", middle.CheckUserAuth("", UserHandler.Logout))
	r.GET(preV1+"/users/me", middle.CheckUserAuth("", UserHandler.GetMyInfo))
	r.GET(preV1+"/users/roles", middle.CheckUserAuth("", UserHandler.GetMyRoles))
	r.POST(preV1+"/users/change_password", middle.CheckUserAuth("", UserHandler.ChangePassword))
	r.POST(preV1+"/users/change_info", middle.CheckUserAuth("", UserHandler.ChangeInfo))
	registerCRUDRoutes(r, preV1+"/users", UserHandler, middle)
}

// InitRounters khởi tạo các route cho ứng dụng
func InitRounters(r *router.Router, c *config.Configuration, db *mongo.Client) {
	middle := middleware.NewJwtToken(c, db)

	// ====================================  STATIC API ===============================================
	StaticHandler := handler.NewStaticHandler()
	r.GET(preV1+"/static/test", StaticHandler.TestApi)
	r.GET(preV1+"/static/system", middle.CheckUserAuth("", StaticHandler.GetSystemStatic))
	r.GET(preV1+"/static/api", middle.CheckUserAuth("", StaticHandler.GetApiStatic))

	// Đăng ký các route xác thực
	registerAuthRoutes(r, c, db, middle)

	// ====================================  AGENTS API =============================================
	AgentHandler := handler.NewAgentHandler()
	registerCRUDRoutes(r, preV1+"/agents", AgentHandler, middle)
	r.POST(preV1+"/agents/checkin/{id}", middle.CheckUserAuth("Agent.CheckIn", AgentHandler.CheckIn))
	r.POST(preV1+"/agents/checkout/{id}", middle.CheckUserAuth("Agent.CheckOut", AgentHandler.CheckOut))

	// ====================================  ACCESSTOKEN API ========================================
	AccessTokenHandler := handler.NewAccessTokenHandler()
	registerCRUDRoutes(r, preV1+"/access_tokens", AccessTokenHandler, middle)

	// ====================================  FBPAGE API =============================================
	FbPageHandler := handler.NewFbPageHandler()
	registerCRUDRoutes(r, preV1+"/fb_pages", FbPageHandler, middle)
	r.POST(preV1+"/fb_pages/update_token", middle.CheckUserAuth("FbPage.UpdateToken", FbPageHandler.UpdateToken))
	r.GET(preV1+"/fb_pages/pageId/{id}", middle.CheckUserAuth("FbPage.Read", FbPageHandler.FindOneByPageID))

	// ====================================  FBCONVERSATION API =====================================
	FbConversationHandler := handler.NewFbConversationHandler()
	registerCRUDRoutes(r, preV1+"/fb_conversations", FbConversationHandler, middle)
	r.GET(preV1+"/fb_conversations/newest", middle.CheckUserAuth("FbConversation.Read", FbConversationHandler.FindAllSortByApiUpdate))

	// ====================================  FBMESSAGE API ==========================================
	FbMessageHandler := handler.NewFbMessageHandler()
	registerCRUDRoutes(r, preV1+"/fb_messages", FbMessageHandler, middle)

	// ====================================  FBPOST API =============================================
	FbPostHandler := handler.NewFbPostHandler()
	registerCRUDRoutes(r, preV1+"/fb_posts", FbPostHandler, middle)

	// ====================================  PCORDER API ============================================
	PcOrderHandler := handler.NewPcOrderHandler()
	registerCRUDRoutes(r, preV1+"/pc_orders", PcOrderHandler, middle)
}
