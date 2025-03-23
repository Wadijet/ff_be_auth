package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// FbPageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Page
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbPageHandler struct {
	BaseHandler
	FbPageHandlerService services.FbPageService
}

// NewFbPageHandler khởi tạo một FbPageHandler mới
func NewFbPageHandler(c *config.Configuration, db *mongo.Client) *FbPageHandler {
	newHandler := new(FbPageHandler)
	newHandler.FbPageHandlerService = *services.NewFbPageService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create tạo mới một FbPage
func (h *FbPageHandler) Create(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.FbPageCreateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.FbPageHandlerService.ReviceData(ctx, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// FindOneById tìm một FbPage theo ID
func (h *FbPageHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.FbPageHandlerService.FindOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindOneByPageID tìm một FbPage theo PageID
func (h *FbPageHandler) FindOneByPageID(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.FbPageHandlerService.FindOneByPageID(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các FbPage với phân trang
func (h *FbPageHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	data, err := h.FbPageHandlerService.FindAll(ctx, page, limit)
	h.HandleResponse(ctx, data, err)
}

// UpdateToken cập nhật access token của một FbPage
func (h *FbPageHandler) UpdateToken(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.FbPageUpdateTokenInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.FbPageHandlerService.UpdateToken(ctx, inputStruct)
	h.HandleResponse(ctx, data, err)
}
