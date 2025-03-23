package handler

import (
	"atk-go-server/app/services"
	"atk-go-server/config"
	"atk-go-server/global"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdminHandler là cấu trúc chứa các dịch vụ cần thiết để xử lý các yêu cầu liên quan đến quản trị viên
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type AdminHandler struct {
	BaseHandler
	UserCRUD       services.RepositoryService
	PermissionCRUD services.RepositoryService
	RoleCRUD       services.RepositoryService
	InitService    services.InitService
	AdminService   services.AdminService
}

// NewAdminHandler khởi tạo một AdminHandler mới với cấu hình và kết nối cơ sở dữ liệu
func NewAdminHandler(c *config.Configuration, db *mongo.Client) *AdminHandler {
	newHandler := new(AdminHandler)
	newHandler.UserCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Users)
	newHandler.PermissionCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newHandler.RoleCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Roles)
	newHandler.InitService = *services.NewInitService(c, db)
	newHandler.AdminService = *services.NewAdminService(c, db)
	return newHandler
}

//=============================================================================

// SetRoleStruct là cấu trúc dữ liệu đầu vào cho việc thiết lập vai trò người dùng
type SetRoleStruct struct {
	Email  string             `json:"email" bson:"email" validate:"required"`
	RoleID primitive.ObjectID `json:"roleID" bson:"roleID" validate:"required"`
}

// SetRole xử lý yêu cầu thiết lập vai trò cho người dùng
func (h *AdminHandler) SetRole(ctx *fasthttp.RequestCtx) {
	input := new(SetRoleStruct)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		inputStruct := input.(*SetRoleStruct)
		return h.AdminService.SetRole(ctx, inputStruct.Email, inputStruct.RoleID)
	})
}

// =================================================================================

// BlockUserInput là cấu trúc dữ liệu đầu vào cho việc khóa người dùng
type BlockUserInput struct {
	Email string `json:"email" bson:"email" validate:"required"`
	Note  string `json:"note" bson:"note" validate:"required"`
}

// BlockUser xử lý yêu cầu khóa người dùng
func (h *AdminHandler) BlockUser(ctx *fasthttp.RequestCtx) {
	input := new(BlockUserInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		inputStruct := input.(*BlockUserInput)
		return h.AdminService.BlockUser(ctx, inputStruct.Email, true, inputStruct.Note)
	})
}

// UnBlockUser xử lý yêu cầu mở khóa người dùng
func (h *AdminHandler) UnBlockUser(ctx *fasthttp.RequestCtx) {
	input := new(BlockUserInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		inputStruct := input.(*BlockUserInput)
		return h.AdminService.BlockUser(ctx, inputStruct.Email, false, inputStruct.Note)
	})
}
