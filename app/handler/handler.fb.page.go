package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"

	"github.com/valyala/fasthttp"
)

// FbPageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Page
// Kế thừa từ BaseHandler với các type parameter:
// - Model: models.FbPage
// - CreateInput: models.FbPageCreateInput
// - UpdateInput: models.FbPageCreateInput
type FbPageHandler struct {
	BaseHandler[models.FbPage, models.FbPageCreateInput, models.FbPageCreateInput]
	FbPageService *services.FbPageService
}

// NewFbPageHandler khởi tạo một FbPageHandler mới
func NewFbPageHandler() *FbPageHandler {
	handler := &FbPageHandler{}
	handler.FbPageService = services.NewFbPageService()
	handler.Service = handler.FbPageService
	return handler
}

// OTHER FUNCTIONS ==========================================================================

// FindOneByPageID tìm một FbPage theo PageID
func (h *FbPageHandler) FindOneByPageID(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.FbPageService.FindOneByPageID(context.Background(), id)
	h.HandleResponse(ctx, data, err)
}

// UpdateToken cập nhật access token của một FbPage
func (h *FbPageHandler) UpdateToken(ctx *fasthttp.RequestCtx) {
	input := new(models.FbPageUpdateTokenInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.FbPageService.UpdateToken(context.Background(), input)
	h.HandleResponse(ctx, data, err)
}
