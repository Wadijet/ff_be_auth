package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// PcPosOrderHandler là cấu trúc xử lý các yêu cầu liên quan đến Pancake POS Order cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.PcPosOrder
// - CreateInput: dto.PcPosOrderCreateInput
// - UpdateInput: dto.PcPosOrderCreateInput
type PcPosOrderHandler struct {
	BaseHandler[models.PcPosOrder, dto.PcPosOrderCreateInput, dto.PcPosOrderCreateInput]
	PcPosOrderService *services.PcPosOrderService
}

// NewPcPosOrderHandler khởi tạo một PcPosOrderHandler mới
// Returns:
//   - *PcPosOrderHandler: Instance mới của PcPosOrderHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewPcPosOrderHandler() (*PcPosOrderHandler, error) {
	handler := &PcPosOrderHandler{}

	// Khởi tạo PcPosOrderService và xử lý error
	service, err := services.NewPcPosOrderService()
	if err != nil {
		return nil, fmt.Errorf("failed to create pc pos order service: %v", err)
	}
	handler.PcPosOrderService = service
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
}
