package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"

	"github.com/valyala/fasthttp"
)

// FbPostHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Post
// Kế thừa từ BaseHandler với các type parameter:
// - Model: models.FbPost
// - CreateInput: models.FbPostCreateInput
// - UpdateInput: models.FbPostCreateInput
type FbPostHandler struct {
	BaseHandler[models.FbPost, models.FbPostCreateInput, models.FbPostCreateInput]
	FbPostService *services.FbPostService
}

// NewFbPostHandler khởi tạo một FbPostHandler mới
func NewFbPostHandler() *FbPostHandler {
	handler := &FbPostHandler{}
	handler.FbPostService = services.NewFbPostService()
	handler.Service = handler.FbPostService
	return handler
}

// OTHER FUNCTIONS ==========================================================================

// FindOneByPostID tìm một FbPost theo PostID
func (h *FbPostHandler) FindOneByPostID(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.FbPostService.FindOneByPostID(context.Background(), id)
	h.HandleResponse(ctx, data, err)
}

// UpdateToken cập nhật access token của một FbPost
func (h *FbPostHandler) UpdateToken(ctx *fasthttp.RequestCtx) {
	input := new(models.FbPostUpdateTokenInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		fbPostInput := input.(*models.FbPostUpdateTokenInput)
		return h.FbPostService.UpdateToken(context.Background(), fbPostInput)
	})
}
