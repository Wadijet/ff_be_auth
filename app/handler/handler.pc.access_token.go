package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"context"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		accessTokenInput := input.(*models.AccessTokenCreateInput)
		return h.AccessTokenService.Create(context.Background(), accessTokenInput)
	})
}

// FindOne tìm một Access Token theo ID
func (h *AccessTokenHandler) FindOne(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.AccessTokenService.FindOneById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các Access Token với phân trang
func (h *AccessTokenHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{} // Có thể thêm filter từ query params nếu cần
	data, err := h.AccessTokenService.FindWithPagination(context, filter, page, limit)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một Access Token theo ID
func (h *AccessTokenHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	input := new(models.AccessTokenUpdateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		accessTokenInput := input.(*models.AccessTokenUpdateInput)
		return h.AccessTokenService.Update(context.Background(), utility.String2ObjectID(id), accessTokenInput)
	})
}

// Delete xóa một Access Token theo ID
func (h *AccessTokenHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.AccessTokenService.DeleteById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}
