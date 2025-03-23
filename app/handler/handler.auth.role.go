package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type RoleHandler struct {
	BaseHandler
	RoleService *services.RoleService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewRoleHandler(c *config.Configuration, db *mongo.Client) *RoleHandler {
	newHandler := new(RoleHandler)
	newHandler.RoleService = services.NewRoleService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới vai trò
func (h *RoleHandler) Create(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.RoleCreateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.RoleService.Create(ctx, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// FindOneById tìm vai trò theo ID
func (h *RoleHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.RoleService.FindOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các vai trò với phân trang
func (h *RoleHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	data, err := h.RoleService.FindAll(ctx, page, limit)
	h.HandleResponse(ctx, data, err)
}

// UpdateOneById cập nhật một vai trò theo ID
func (h *RoleHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	inputStruct := new(models.RoleUpdateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.RoleService.Update(ctx, id, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// DeleteOneById xóa một vai trò theo ID
func (h *RoleHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.RoleService.Delete(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// Other functions =========================================================================
