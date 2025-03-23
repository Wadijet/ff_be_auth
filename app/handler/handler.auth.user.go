package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserHandler là struct chứa các dịch vụ và repository cần thiết để xử lý người dùng
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type UserHandler struct {
	BaseHandler
	RoleService services.RoleService
	UserService services.UserService
}

// NewUserHandler khởi tạo một UserHandler mới
func NewUserHandler(c *config.Configuration, db *mongo.Client) *UserHandler {
	newHandler := new(UserHandler)
	newHandler.UserService = *services.NewUserService(c, db)
	newHandler.RoleService = *services.NewRoleService(c, db)

	return newHandler
}

// CRUD functions ======================================================

// FindOneById tìm một người dùng theo ID
func (h *UserHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.UserService.FindOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindAllWithFilter tìm tất cả người dùng với bộ lọc
func (h *UserHandler) FindAllWithFilter(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	data, err := h.UserService.FindAll(ctx, page, limit)
	h.HandleResponse(ctx, data, err)
}

// OTHER functions =======================================================

// Registry đăng ký người dùng mới
func (h *UserHandler) Registry(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.UserCreateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.UserService.Create(ctx, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// Login đăng nhập người dùng
func (h *UserHandler) Login(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.UserLoginInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.UserService.Login(ctx, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// Logout đăng xuất người dùng
func (h *UserHandler) Logout(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	inputStruct := new(models.UserLogoutInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.UserService.Logout(ctx, strMyID, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// GetMyInfo lấy thông tin của người dùng hiện tại
func (h *UserHandler) GetMyInfo(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	data, err := h.UserService.FindOneById(ctx, strMyID)
	h.HandleResponse(ctx, data, err)
}

// GetMyRoles lấy danh sách vai trò của người dùng hiện tại
func (h *UserHandler) GetMyRoles(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	data, err := h.UserService.GetRoles(ctx, strMyID)
	h.HandleResponse(ctx, data, err)
}

// ChangePassword thay đổi mật khẩu người dùng
func (h *UserHandler) ChangePassword(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	inputStruct := new(models.UserChangePasswordInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.UserService.ChangePassword(ctx, strMyID, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// ChangeInfo thay đổi thông tin người dùng
func (h *UserHandler) ChangeInfo(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	inputStruct := new(models.UserChangeInfoInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.UserService.ChangeInfo(ctx, strMyID, inputStruct)
	h.HandleResponse(ctx, data, err)
}

//
