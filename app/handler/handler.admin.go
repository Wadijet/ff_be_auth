package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/config"
	"meta_commerce/global"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdminHandler là cấu trúc chứa các dịch vụ cần thiết để xử lý các yêu cầu liên quan đến quản trị viên
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type AdminHandler struct {
	BaseHandler
	UserCRUD       services.BaseService[models.User]
	PermissionCRUD services.BaseService[models.Permission]
	RoleCRUD       services.BaseService[models.Role]
	InitService    services.InitService
	AdminService   services.AdminService
}

// NewAdminHandler khởi tạo một AdminHandler mới với cấu hình và kết nối cơ sở dữ liệu
func NewAdminHandler(c *config.Configuration, db *mongo.Client) *AdminHandler {
	newHandler := new(AdminHandler)

	// Khởi tạo các collection
	userCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Users)).Collection(global.MongoDB_ColNames.Users)
	permissionCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Permissions)).Collection(global.MongoDB_ColNames.Permissions)
	roleCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Roles)).Collection(global.MongoDB_ColNames.Roles)

	// Khởi tạo các service với BaseService
	newHandler.UserCRUD = services.NewBaseService[models.User](userCol)
	newHandler.PermissionCRUD = services.NewBaseService[models.Permission](permissionCol)
	newHandler.RoleCRUD = services.NewBaseService[models.Role](roleCol)
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
