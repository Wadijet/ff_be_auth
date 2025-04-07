package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
)

// FbMessageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Message cho Fiber
// Kế thừa từ FiberBaseHandler để sử dụng các phương thức xử lý chung
type FbMessageHandler struct {
	BaseHandler[models.FbMessage, models.FbMessageCreateInput, models.FbMessageCreateInput]
	FbMessageService *services.FbMessageService
}

// NewFbMessageHandler khởi tạo một FiberFbMessageHandler mới
// Returns:
//   - *FiberFbMessageHandler: Instance mới của FiberFbMessageHandler đã được khởi tạo với các service cần thiết
func NewFbMessageHandler() *FbMessageHandler {
	handler := &FbMessageHandler{}
	handler.FbMessageService = services.NewFbMessageService()
	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng FbMessageService trực tiếp
	return handler
}

// Các phương thức được kế thừa từ FiberBaseHandler:
// - HandleCreate: Tạo mới một FbMessage
// - HandleUpdate: Cập nhật một FbMessage
// - HandleDelete: Xóa một FbMessage
// - HandleFindById: Tìm FbMessage theo ID
// - HandleFindAll: Lấy danh sách FbMessage có phân trang

// Các hàm đặc thù của FbMessage (nếu có) sẽ được thêm vào đây
