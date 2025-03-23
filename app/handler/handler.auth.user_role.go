package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò của người dùng
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type UserRoleHandler struct {
	BaseHandler
	UserRoleService *services.UserRoleService
}

// NewUserRoleHandler khởi tạo một UserRoleHandler mới
func NewUserRoleHandler(c *config.Configuration, db *mongo.Client) *UserRoleHandler {
	newHandler := new(UserRoleHandler)
	newHandler.UserRoleService = services.NewUserRoleService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới UserRole
func (h *UserRoleHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.UserRoleCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		return h.UserRoleService.Create(ctx, input.(*models.UserRoleCreateInput))
	})
}

// Delete xóa một UserRole
func (h *UserRoleHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.UserRoleService.Delete(ctx, id)
	h.HandleResponse(ctx, data, err)
}
