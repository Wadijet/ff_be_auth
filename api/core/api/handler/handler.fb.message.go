package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// FbMessageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Message cho Fiber
// Kế thừa từ FiberBaseHandler để sử dụng các phương thức xử lý chung
type FbMessageHandler struct {
	BaseHandler[models.FbMessage, dto.FbMessageCreateInput, dto.FbMessageCreateInput]
	FbMessageService *services.FbMessageService
}

// NewFbMessageHandler khởi tạo một FiberFbMessageHandler mới
// Returns:
//   - *FiberFbMessageHandler: Instance mới của FiberFbMessageHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewFbMessageHandler() (*FbMessageHandler, error) {
	// Khởi tạo FbMessageService và xử lý error
	service, err := services.NewFbMessageService()
	if err != nil {
		return nil, fmt.Errorf("failed to create message service: %v", err)
	}

	handler := &FbMessageHandler{
		BaseHandler:      *NewBaseHandler[models.FbMessage, dto.FbMessageCreateInput, dto.FbMessageCreateInput](service.BaseServiceMongoImpl),
		FbMessageService: service,
	}

	return handler, nil
}

// Các phương thức được kế thừa từ FiberBaseHandler:
// - HandleCreate: Tạo mới một FbMessage
// - HandleUpdate: Cập nhật một FbMessage
// - HandleDelete: Xóa một FbMessage
// - HandleFindById: Tìm FbMessage theo ID
// - HandleFindAll: Lấy danh sách FbMessage có phân trang

// Các hàm đặc thù của FbMessage (nếu có) sẽ được thêm vào đây
