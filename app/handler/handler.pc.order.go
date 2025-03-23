package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// PcOrderHandler là cấu trúc xử lý các yêu cầu liên quan đến đơn hàng
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type PcOrderHandler struct {
	BaseHandler
	PcOrderService services.PcOrderService
}

// NewPcOrderHandler khởi tạo một PcOrderHandler mới
func NewPcOrderHandler(c *config.Configuration, db *mongo.Client) *PcOrderHandler {
	newHandler := new(PcOrderHandler)
	newHandler.PcOrderService = *services.NewPcOrderService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới PcOrder
func (h *PcOrderHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.PcOrderCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		inputStruct := input.(*models.PcOrderCreateInput)
		return h.PcOrderService.ReviceData(ctx, inputStruct)
	})
}

// FindOneById tìm PcOrder theo ID
func (h *PcOrderHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.PcOrderService.FindOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các PcOrder với phân trang
func (h *PcOrderHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	data, err := h.PcOrderService.FindAll(ctx, page, limit)
	h.HandleResponse(ctx, data, err)
}
