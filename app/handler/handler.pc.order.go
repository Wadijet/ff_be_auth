package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/config"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
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
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		pcOrderInput := input.(*models.PcOrderCreateInput)
		return h.PcOrderService.ReviceData(context.Background(), pcOrderInput)
	})
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
	filter := bson.M{}
	data, err := h.PcOrderService.FindWithPagination(context, filter, page, limit)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một PcOrder
func (h *PcOrderHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	input := new(models.PcOrderCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		pcOrderInput := input.(*models.PcOrderCreateInput)
		pcOrder := models.PcOrder{
			ID:          objectID,
			PanCakeData: pcOrderInput.PanCakeData,
			UpdatedAt:   time.Now().Unix(),
		}
		return h.PcOrderService.UpdateById(context.Background(), objectID, pcOrder)
	})
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
	err = h.PcOrderService.DeleteById(context, objectID)
	h.HandleResponse(ctx, nil, err)
}
