package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"
	"context"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AgentHandler là cấu trúc xử lý các yêu cầu liên quan đến đại lý
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type AgentHandler struct {
	BaseHandler
	AgentService *services.AgentService
}

// NewAgentHandler khởi tạo một AgentHandler mới
func NewAgentHandler(c *config.Configuration, db *mongo.Client) *AgentHandler {
	newHandler := new(AgentHandler)
	newHandler.AgentService = services.NewAgentService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create tạo mới một Agent
func (h *AgentHandler) Create(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.AgentCreateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.AgentService.Create(context, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// FindOneById tìm một Agent theo ID
func (h *AgentHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.AgentService.FindOne(context, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các Agent với phân trang
func (h *AgentHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{} // Có thể thêm filter từ query params nếu cần

	// Tạo options cho phân trang
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

	data, err := h.AgentService.FindAll(context, filter, findOptions)
	h.HandleResponse(ctx, data, err)
}

// UpdateOneById cập nhật một Agent theo ID
func (h *AgentHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	inputStruct := new(models.AgentUpdateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.AgentService.Update(context, id, inputStruct)
	h.HandleResponse(ctx, data, err)
}

// DeleteOneById xóa một Agent theo ID
func (h *AgentHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.AgentService.Delete(context, id)
	h.HandleResponse(ctx, nil, err)
}

// CheckIn xử lý check-in cho Agent
func (h *AgentHandler) CheckIn(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	context := context.Background()
	data, err := h.AgentService.CheckIn(context, strMyID)
	h.HandleResponse(ctx, data, err)
}
