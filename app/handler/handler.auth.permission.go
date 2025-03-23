package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"
	"context"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Create tạo mới một quyền
func (h *PermissionHandler) Create(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.PermissionCreateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.PermissionService.Create(context, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// FindOneById tìm kiếm một quyền theo ID
func (h *PermissionHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.PermissionService.FindOne(context, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm kiếm tất cả các quyền với phân trang
func (h *PermissionHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{} // Có thể thêm filter từ query params nếu cần

	// Tạo options cho phân trang
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

	data, err := h.PermissionService.FindAll(context, filter, findOptions)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một quyền
func (h *PermissionHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	inputStruct := new(models.PermissionUpdateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.PermissionService.Update(context, id, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// Delete xóa một quyền
func (h *PermissionHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.PermissionService.Delete(context, id)
	h.HandleResponse(ctx, nil, err)
}

// Other functions =========================================================================
