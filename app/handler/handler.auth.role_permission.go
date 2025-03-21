package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// RolePermissionHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
type RolePermissionHandler struct {
	RolePermissionService *services.RolePermissionService
}

// NewRolePermissionHandler khởi tạo một RolePermissionHandler mới
func NewRolePermissionHandler(c *config.Configuration, db *mongo.Client) *RolePermissionHandler {
	newHandler := new(RolePermissionHandler)
	newHandler.RolePermissionService = services.NewRolePermissionService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới RolePermission
func (h *RolePermissionHandler) Create(ctx *fasthttp.RequestCtx) {
	utility.GenericHandler[models.RolePermissionCreateInput](ctx, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		inputStruct := input.(*models.RolePermissionCreateInput)
		return h.RolePermissionService.Create(ctx, inputStruct)
	})
}

// Xóa một RolePermission
func (h *RolePermissionHandler) Delete(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.RolePermissionService.Delete(ctx, id))
	if response["error"] != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK) // 200 OK
	}
	utility.JSON(ctx, response)
}
