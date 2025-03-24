package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"
	"context"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PcOrderHandler là cấu trúc xử lý các yêu cầu liên quan đến đơn hàng
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type PcOrderHandler struct {
	BaseHandler
	PcOrderService *services.PcOrderService
}

// NewPcOrderHandler khởi tạo một PcOrderHandler mới
func NewPcOrderHandler(c *config.Configuration, db *mongo.Client) *PcOrderHandler {
	newHandler := new(PcOrderHandler)
	newHandler.PcOrderService = services.NewPcOrderService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới PcOrder
func (h *PcOrderHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.PcOrderCreateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.PcOrderService.ReviceData(context, input)
	h.HandleResponse(ctx, data, err)
}

// FindOne tìm một PcOrder theo ID
func (h *PcOrderHandler) FindOne(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)

	// Chuyển đổi string ID thành ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	context := context.Background()
	data, err := h.PcOrderService.FindOne(context, objectID)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các PcOrder với phân trang
func (h *PcOrderHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	data, err := h.PcOrderService.FindAll(context, page, limit)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một PcOrder
func (h *PcOrderHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)

	// Chuyển đổi string ID thành ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	input := new(models.PcOrderCreateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	pcOrder := models.PcOrder{
		ID:          objectID,
		PanCakeData: input.PanCakeData,
		UpdatedAt:   time.Now().Unix(),
	}

	data, err := h.PcOrderService.Update(context, objectID, pcOrder)
	h.HandleResponse(ctx, data, err)
}

// Delete xóa một PcOrder
func (h *PcOrderHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)

	// Chuyển đổi string ID thành ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	context := context.Background()
	err = h.PcOrderService.Delete(context, objectID)
	h.HandleResponse(ctx, nil, err)
}
