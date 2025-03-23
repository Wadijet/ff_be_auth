package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// FbMessageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Message
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbMessageHandler struct {
	BaseHandler
	FbMessageService services.FbMessageService
}

// NewFbMessageHandler khởi tạo một FbMessageHandler mới
func NewFbMessageHandler(c *config.Configuration, db *mongo.Client) *FbMessageHandler {
	newHandler := new(FbMessageHandler)
	newHandler.FbMessageService = *services.NewFbMessageService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới FbMessage
func (h *FbMessageHandler) Create(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.FbMessageCreateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	data, err := h.FbMessageService.ReviceData(ctx, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// FindOneById tìm một FbMessage theo ID
func (h *FbMessageHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.FbMessageService.FindOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các FbMessage với phân trang
func (h *FbMessageHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	data, err := h.FbMessageService.FindAll(ctx, page, limit)
	h.HandleResponse(ctx, data, err)
}
