package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
type UserRoleHandler struct {
	crudUserRole    services.RepositoryService
	UserRoleService services.UserRoleService
}

// NewUserRoleHandler khởi tạo một UserRoleHandler mới
func NewUserRoleHandler(c *config.Configuration, db *mongo.Client) *UserRoleHandler {
	newHandler := new(UserRoleHandler)
	newHandler.crudUserRole = *services.NewRepository(c, db, global.MongoDB_ColNames.UserRoles)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo mới một UserRole
func (h *UserRoleHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.UserRoleCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.UserRoleService.Create(ctx, inputStruct))
		}
	}
	utility.JSON(ctx, response)
}

// Xóa một UserRole
func (h *UserRoleHandler) Delete(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crudUserRole.DeleteOneById(ctx, id))
	utility.JSON(ctx, response)
}
