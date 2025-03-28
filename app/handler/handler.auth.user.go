package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
	"meta_commerce/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserHandler là struct chứa các dịch vụ và repository cần thiết để xử lý người dùng
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type UserHandler struct {
	BaseHandler[models.User, models.UserCreateInput, models.UserChangeInfoInput]
	RoleService *services.RoleService
	UserService *services.UserService
}

// NewUserHandler khởi tạo một UserHandler mới
func NewUserHandler(c *config.Configuration, db *mongo.Client) *UserHandler {
	newHandler := new(UserHandler)
	newHandler.UserService = services.NewUserService(c, db)
	newHandler.RoleService = services.NewRoleService(c, db)
	newHandler.BaseHandler.Service = newHandler.UserService // Gán service cho BaseHandler

	return newHandler
}

// CRUD functions ======================================================

// FindOneById tìm một người dùng theo ID
func (h *UserHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.UserService.FindOneById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, data, err)
}

// FindAllWithFilter tìm tất cả người dùng với bộ lọc
func (h *UserHandler) FindAllWithFilter(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{} // Có thể thêm filter từ query params nếu cần
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit)
	data, err := h.UserService.Find(context, filter, opts)
	h.HandleResponse(ctx, data, err)
}

// OTHER functions =======================================================

// Registry đăng ký người dùng mới
func (h *UserHandler) Registry(ctx *fasthttp.RequestCtx) {
	input := new(models.UserCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		userInput := input.(*models.UserCreateInput)
		return h.UserService.Registry(context.Background(), userInput)
	})
}

// Login đăng nhập người dùng
func (h *UserHandler) Login(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.UserLoginInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, utility.NewError(utility.ErrCodeValidationFormat, "Invalid request body", utility.StatusBadRequest, nil))
		return
	}

	context := context.Background()
	data, err := h.UserService.Login(context, inputStruct)
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

	context := context.Background()
	err := h.UserService.Logout(context, utility.String2ObjectID(strMyID), inputStruct)
	h.HandleResponse(ctx, nil, err)
}

// GetMyInfo lấy thông tin của người dùng hiện tại
func (h *UserHandler) GetMyInfo(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	context := context.Background()
	data, err := h.UserService.FindOneById(context, utility.String2ObjectID(strMyID))
	h.HandleResponse(ctx, data, err)
}

// GetMyRoles lấy danh sách vai trò của người dùng hiện tại
func (h *UserHandler) GetMyRoles(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	context := context.Background()

	// Lấy thông tin user
	user, err := h.UserService.FindOneById(context, utility.String2ObjectID(strMyID))
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	// Lấy thông tin role từ token
	if user.Token == "" {
		h.HandleResponse(ctx, nil, nil)
		return
	}

	role, err := h.RoleService.FindOneById(context, utility.String2ObjectID(user.Token))
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	h.HandleResponse(ctx, role, nil)
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

	context := context.Background()
	err := h.UserService.ChangePassword(context, utility.String2ObjectID(strMyID), inputStruct)
	h.HandleResponse(ctx, nil, err)
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

	context := context.Background()
	user := models.User{
		Name: inputStruct.Name,
	}
	data, err := h.UserService.UpdateById(context, utility.String2ObjectID(strMyID), user)
	h.HandleResponse(ctx, data, err)
}

//
