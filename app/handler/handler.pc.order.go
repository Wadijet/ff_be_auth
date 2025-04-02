package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
)

// PcOrderHandler là cấu trúc xử lý các yêu cầu liên quan đến đơn hàng
// Kế thừa từ BaseHandler với các type parameter:
// - Model: models.PcOrder
// - CreateInput: models.PcOrderCreateInput
// - UpdateInput: models.PcOrderCreateInput
type PcOrderHandler struct {
	BaseHandler[models.PcOrder, models.PcOrderCreateInput, models.PcOrderCreateInput]
	PcOrderService *services.PcOrderService
}

// NewPcOrderHandler khởi tạo một PcOrderHandler mới
func NewPcOrderHandler() *PcOrderHandler {
	handler := &PcOrderHandler{}
	handler.PcOrderService = services.NewPcOrderService()
	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng PcOrderService trực tiếp
	return handler
}

// Các hàm đặc thù của PcOrder (nếu có) sẽ được thêm vào đây
