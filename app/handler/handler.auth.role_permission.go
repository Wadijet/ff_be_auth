package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// RolePermissionHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò và quyền
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type RolePermissionHandler struct {
	BaseHandler
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
	input := new(models.RolePermissionCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		return h.RolePermissionService.Create(ctx, input.(*models.RolePermissionCreateInput))
	})
}

// Delete xóa một RolePermission
func (h *RolePermissionHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.RolePermissionService.Delete(ctx, id)
	h.HandleResponse(ctx, data, err)
}
