package handler

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// PcOrderHandler là cấu trúc xử lý các yêu cầu liên quan đến đơn hàng cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.PcOrder
// - CreateInput: models.PcOrderCreateInput
// - UpdateInput: models.PcOrderCreateInput
type PcOrderHandler struct {
	BaseHandler[models.PcOrder, models.PcOrderCreateInput, models.PcOrderCreateInput]
	PcOrderService *services.PcOrderService
}

// NewPcOrderHandler khởi tạo một FiberPcOrderHandler mới
// Returns:
//   - *FiberPcOrderHandler: Instance mới của FiberPcOrderHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewPcOrderHandler() (*PcOrderHandler, error) {
	handler := &PcOrderHandler{}

	// Khởi tạo PcOrderService và xử lý error
	service, err := services.NewPcOrderService()
	if err != nil {
		return nil, fmt.Errorf("failed to create order service: %v", err)
	}
	handler.PcOrderService = service

	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng PcOrderService trực tiếp
	return handler, nil
}

// Các phương thức được kế thừa từ FiberBaseHandler:
// - HandleCreate: Tạo mới một PcOrder
// - HandleUpdate: Cập nhật một PcOrder
// - HandleDelete: Xóa một PcOrder
// - HandleFindById: Tìm PcOrder theo ID
// - HandleFindAll: Lấy danh sách PcOrder có phân trang

// Các hàm đặc thù của PcOrder (nếu có) sẽ được thêm vào đây
