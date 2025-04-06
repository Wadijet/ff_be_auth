package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
)

// FiberPcOrderHandler là cấu trúc xử lý các yêu cầu liên quan đến đơn hàng cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.PcOrder
// - CreateInput: models.PcOrderCreateInput
// - UpdateInput: models.PcOrderCreateInput
type FiberPcOrderHandler struct {
	FiberBaseHandler[models.PcOrder, models.PcOrderCreateInput, models.PcOrderCreateInput]
	PcOrderService *services.PcOrderService
}

// NewFiberPcOrderHandler khởi tạo một FiberPcOrderHandler mới
// Returns:
//   - *FiberPcOrderHandler: Instance mới của FiberPcOrderHandler đã được khởi tạo với các service cần thiết
func NewFiberPcOrderHandler() *FiberPcOrderHandler {
	handler := &FiberPcOrderHandler{}
	handler.PcOrderService = services.NewPcOrderService()
	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng PcOrderService trực tiếp
	return handler
}

// Các phương thức được kế thừa từ FiberBaseHandler:
// - HandleCreate: Tạo mới một PcOrder
// - HandleUpdate: Cập nhật một PcOrder
// - HandleDelete: Xóa một PcOrder
// - HandleFindById: Tìm PcOrder theo ID
// - HandleFindAll: Lấy danh sách PcOrder có phân trang

// Các hàm đặc thù của PcOrder (nếu có) sẽ được thêm vào đây
