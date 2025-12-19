package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// PcPosCategoryHandler là cấu trúc xử lý các yêu cầu liên quan đến Pancake POS Category cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.PcPosCategory
// - CreateInput: dto.PcPosCategoryCreateInput
// - UpdateInput: dto.PcPosCategoryCreateInput
type PcPosCategoryHandler struct {
	BaseHandler[models.PcPosCategory, dto.PcPosCategoryCreateInput, dto.PcPosCategoryCreateInput]
	PcPosCategoryService *services.PcPosCategoryService
}

// NewPcPosCategoryHandler khởi tạo một PcPosCategoryHandler mới
// Returns:
//   - *PcPosCategoryHandler: Instance mới của PcPosCategoryHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewPcPosCategoryHandler() (*PcPosCategoryHandler, error) {
	handler := &PcPosCategoryHandler{}

	// Khởi tạo PcPosCategoryService và xử lý error
	service, err := services.NewPcPosCategoryService()
	if err != nil {
		return nil, fmt.Errorf("failed to create pc pos category service: %v", err)
	}
	handler.PcPosCategoryService = service
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
}
