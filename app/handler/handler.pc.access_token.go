package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// AccessTokenHandler là cấu trúc xử lý các yêu cầu liên quan đến Access Token
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type AccessTokenHandler struct {
	BaseHandler
	AccessTokenService services.AccessTokenService
}

// NewAccessTokenHandler khởi tạo một AccessTokenHandler mới
func NewAccessTokenHandler(c *config.Configuration, db *mongo.Client) *AccessTokenHandler {
	newHandler := new(AccessTokenHandler)
	newHandler.AccessTokenService = *services.NewAccessTokenService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create tạo mới một Access Token
func (h *AccessTokenHandler) Create(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.AccessTokenCreateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.AccessTokenService.Create(ctx, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// FindOneById tìm một Access Token theo ID
func (h *AccessTokenHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.AccessTokenService.FindOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các Access Token với phân trang
func (h *AccessTokenHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	data, err := h.AccessTokenService.FindAll(ctx, page, limit)
	h.HandleResponse(ctx, data, err)
}

// UpdateOneById cập nhật một Access Token theo ID
func (h *AccessTokenHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	inputStruct := new(models.AccessTokenUpdateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.AccessTokenService.UpdateOneById(ctx, id, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// DeleteOneById xóa một Access Token theo ID
func (h *AccessTokenHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.AccessTokenService.DeleteOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}
