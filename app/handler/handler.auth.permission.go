package handler

import (
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// PermissionHandler là struct chứa các phương thức xử lý quyền
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type PermissionHandler struct {
	BaseHandler
	PermissionService *services.PermissionService
}

// NewPermissionHandler khởi tạo một PermissionHandler mới
func NewPermissionHandler(c *config.Configuration, db *mongo.Client) *PermissionHandler {
	newHandler := new(PermissionHandler)
	newHandler.PermissionService = services.NewPermissionService(c, db)
	return newHandler
}

// CRUD functions =========================================================================

// FindOneById tìm kiếm một quyền theo ID
func (h *PermissionHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.PermissionService.FindOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm kiếm tất cả các quyền với phân trang
func (h *PermissionHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	data, err := h.PermissionService.FindAll(ctx, page, limit)
	h.HandleResponse(ctx, data, err)
}

// Other functions =========================================================================
