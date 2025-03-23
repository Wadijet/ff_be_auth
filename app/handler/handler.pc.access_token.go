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

// AccessTokenHandler là cấu trúc xử lý các yêu cầu liên quan đến Access Token
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type AccessTokenHandler struct {
	BaseHandler
	AccessTokenService *services.AccessTokenService
}

// NewAccessTokenHandler khởi tạo một AccessTokenHandler mới
func NewAccessTokenHandler(c *config.Configuration, db *mongo.Client) *AccessTokenHandler {
	newHandler := new(AccessTokenHandler)
	newHandler.AccessTokenService = services.NewAccessTokenService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create tạo mới một Access Token
func (h *AccessTokenHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.AccessTokenCreateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.AccessTokenService.Create(context, input)
	h.HandleResponse(ctx, data, err)
}

// FindOne tìm một Access Token theo ID
func (h *AccessTokenHandler) FindOne(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.AccessTokenService.FindOne(context, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các Access Token với phân trang
func (h *AccessTokenHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{} // Có thể thêm filter từ query params nếu cần
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit)
	data, err := h.AccessTokenService.FindAll(context, filter, opts)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một Access Token theo ID
func (h *AccessTokenHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	input := new(models.AccessTokenUpdateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.AccessTokenService.Update(context, id, input)
	h.HandleResponse(ctx, data, err)
}

// Delete xóa một Access Token theo ID
func (h *AccessTokenHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.AccessTokenService.Delete(context, id)
	h.HandleResponse(ctx, nil, err)
}
