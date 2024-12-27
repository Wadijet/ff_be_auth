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

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
type RolePermissionHandler struct {
	crud                  services.RepositoryService
	RolePermissionService services.RolePermissionService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewRolePermissionHandler(c *config.Configuration, db *mongo.Client) *RolePermissionHandler {
	newHandler := new(RolePermissionHandler)
	newHandler.crud = *services.NewRepository(c, db, global.MongoDB_ColNames.RolePermissions)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo mới một RolePermission
func (h *RolePermissionHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.RolePermissionCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.RolePermissionService.Create(ctx, inputStruct))
		}
	}
	utility.JSON(ctx, response)
}

// Xóa một RolePermission
func (h *RolePermissionHandler) Delete(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crud.DeleteOneById(ctx, id))
	utility.JSON(ctx, response)
}