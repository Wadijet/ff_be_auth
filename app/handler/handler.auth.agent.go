package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
)

// AgentHandler xử lý các route liên quan đến đại lý cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
type AgentHandler struct {
	BaseHandler[models.Agent, models.AgentCreateInput, models.AgentUpdateInput]
	AgentService *services.AgentService
}

// NewAgentHandler tạo một instance mới của FiberAgentHandler
// Returns:
//   - *FiberAgentHandler: Instance mới của FiberAgentHandler đã được khởi tạo với các service cần thiết
func NewAgentHandler() *AgentHandler {
	handler := &AgentHandler{}
	handler.AgentService = services.NewAgentService()
	handler.Service = handler.AgentService // Gán service cho BaseHandler
	return handler
}

// HandleCheckIn xử lý check-in cho Agent
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Path Params:
//   - userId: ID của agent cần check-in
//
// Response:
//   - 200: Check-in thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "name": "...",
//     "checkInTime": 123,
//     "status": "CHECKED_IN"
//     }
//     }
//   - 400: ID không hợp lệ
//   - 404: Không tìm thấy agent
//   - 500: Lỗi server
func (h *AgentHandler) HandleCheckIn(c fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId == nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	strMyID := userId.(string)
	result, err := h.AgentService.CheckIn(context.Background(), utility.String2ObjectID(strMyID))
	h.HandleResponse(c, result, err)
	return nil
}

// HandleCheckOut xử lý check-out cho Agent
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Path Params:
//   - userId: ID của agent cần check-out
//
// Response:
//   - 200: Check-out thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "name": "...",
//     "checkOutTime": 123,
//     "status": "CHECKED_OUT"
//     }
//     }
//   - 400: ID không hợp lệ
//   - 404: Không tìm thấy agent
//   - 500: Lỗi server
func (h *AgentHandler) HandleCheckOut(c fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId == nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	strMyID := userId.(string)
	result, err := h.AgentService.CheckOut(context.Background(), utility.String2ObjectID(strMyID))
	h.HandleResponse(c, result, err)
	return nil
}
