package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// PcPosProductHandler là cấu trúc xử lý các yêu cầu liên quan đến Pancake POS Product cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.PcPosProduct
// - CreateInput: dto.PcPosProductCreateInput
// - UpdateInput: dto.PcPosProductCreateInput
type PcPosProductHandler struct {
	BaseHandler[models.PcPosProduct, dto.PcPosProductCreateInput, dto.PcPosProductCreateInput]
	PcPosProductService *services.PcPosProductService
}

// NewPcPosProductHandler khởi tạo một PcPosProductHandler mới
// Returns:
//   - *PcPosProductHandler: Instance mới của PcPosProductHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewPcPosProductHandler() (*PcPosProductHandler, error) {
	handler := &PcPosProductHandler{}

	// Khởi tạo PcPosProductService và xử lý error
	service, err := services.NewPcPosProductService()
	if err != nil {
		return nil, fmt.Errorf("failed to create pc pos product service: %v", err)
	}
	handler.PcPosProductService = service
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
}
