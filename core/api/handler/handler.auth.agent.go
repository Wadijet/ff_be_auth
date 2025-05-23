package handler

import (
	"context"
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/utility"

	"github.com/gofiber/fiber/v3"
)

// AgentHandler xử lý các route liên quan đến đại lý
// Kế thừa từ BaseHandler để có các chức năng CRUD cơ bản
type AgentHandler struct {
	*BaseHandler[models.Agent, models.AgentCreateInput, models.AgentUpdateInput]
	agentService *services.AgentService
}

// NewAgentHandler tạo một instance mới của AgentHandler
// Returns:
//   - *AgentHandler: Instance mới của AgentHandler
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewAgentHandler() (*AgentHandler, error) {
	handler := &AgentHandler{}

	// Khởi tạo base handler
	baseHandler := &BaseHandler[models.Agent, models.AgentCreateInput, models.AgentUpdateInput]{}
	handler.BaseHandler = baseHandler

	// Khởi tạo agent service
	agentService, err := services.NewAgentService()
	if err != nil {
		return nil, fmt.Errorf("failed to create agent service: %v", err)
	}
	handler.agentService = agentService
	handler.BaseService = agentService // Gán service cho BaseHandler

	return handler, nil
}

// HandleCheckIn xử lý check-in cho Agent
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
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
	// Lấy id từ param thông qua hàm tiện ích của BaseHandler
	idParam := h.GetIDFromContext(c)
	if idParam == "" {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Thiếu id agent", common.StatusBadRequest, nil))
		return nil
	}

	objID := utility.String2ObjectID(idParam)
	if objID.IsZero() {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "ID agent không hợp lệ", common.StatusBadRequest, nil))
		return nil
	}

	result, err := h.agentService.CheckIn(context.Background(), objID)
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
	userID := c.Locals("user_id")
	if userID == nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeAuth, "User not authenticated", common.StatusUnauthorized, nil))
		return nil
	}

	objID := utility.String2ObjectID(userID.(string))
	result, err := h.agentService.CheckOut(context.Background(), objID)
	h.HandleResponse(c, result, err)
	return nil
}
