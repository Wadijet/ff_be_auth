package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/valyala/fasthttp"
)

// AgentHandler là cấu trúc xử lý các yêu cầu liên quan đến đại lý
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type AgentHandler struct {
	BaseHandler[models.Agent, models.AgentCreateInput, models.AgentUpdateInput]
	AgentService *services.AgentService
}

// NewAgentHandler khởi tạo một AgentHandler mới
func NewAgentHandler() *AgentHandler {
	newHandler := new(AgentHandler)
	newHandler.AgentService = services.NewAgentService()
	newHandler.BaseHandler.Service = newHandler.AgentService // Gán service cho BaseHandler
	return newHandler
}

// Các hàm đặc thù của Agent (nếu có) sẽ được thêm vào đây

// CheckIn xử lý check-in cho Agent
func (h *AgentHandler) CheckIn(ctx *fasthttp.RequestCtx) {

	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	context := context.Background()
	data, err := h.AgentService.CheckIn(context, utility.String2ObjectID(strMyID))
	h.HandleResponse(ctx, data, err)
}

// CheckOut xử lý check-out cho Agent
func (h *AgentHandler) CheckOut(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("userId") == nil {
		h.HandleError(ctx, nil)
		return
	}

	strMyID := ctx.UserValue("userId").(string)
	context := context.Background()
	data, err := h.AgentService.CheckOut(context, utility.String2ObjectID(strMyID))
	h.HandleResponse(ctx, data, err)
}
