package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
	"meta_commerce/config"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// FbMessageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Message
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbMessageHandler struct {
	BaseHandler
	FbMessageService *services.FbMessageService
}

// NewFbMessageHandler khởi tạo một FbMessageHandler mới
func NewFbMessageHandler(c *config.Configuration, db *mongo.Client) *FbMessageHandler {
	newHandler := new(FbMessageHandler)
	newHandler.FbMessageService = services.NewFbMessageService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới FbMessage
func (h *FbMessageHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.FbMessageCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		fbMessageInput := input.(*models.FbMessageCreateInput)
		return h.FbMessageService.Upsert(context.Background(), fbMessageInput)
	})
}

// FindOne tìm một FbMessage theo ID
func (h *FbMessageHandler) FindOne(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}
	data, err := h.FbMessageService.FindOneById(context.Background(), objectID)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các FbMessage với phân trang
func (h *FbMessageHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	data, err := h.FbMessageService.FindAll(context.Background(), page, limit)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một FbMessage
func (h *FbMessageHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	input := new(models.FbMessageCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		fbMessageInput := input.(*models.FbMessageCreateInput)
		message := models.FbMessage{
			ID:           objectID,
			PageId:       fbMessageInput.PageId,
			PageUsername: fbMessageInput.PageUsername,
			PanCakeData:  fbMessageInput.PanCakeData,
			UpdatedAt:    time.Now().Unix(),
		}
		return h.FbMessageService.UpdateById(context.Background(), objectID, message)
	})
}

// Delete xóa một FbMessage
func (h *FbMessageHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	err := h.FbMessageService.DeleteById(context.Background(), utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}
