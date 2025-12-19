package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// PcPosShopHandler là cấu trúc xử lý các yêu cầu liên quan đến Pancake POS Shop cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.PcPosShop
// - CreateInput: dto.PcPosShopCreateInput
// - UpdateInput: dto.PcPosShopCreateInput
type PcPosShopHandler struct {
	BaseHandler[models.PcPosShop, dto.PcPosShopCreateInput, dto.PcPosShopCreateInput]
	PcPosShopService *services.PcPosShopService
}

// NewPcPosShopHandler khởi tạo một PcPosShopHandler mới
// Returns:
//   - *PcPosShopHandler: Instance mới của PcPosShopHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewPcPosShopHandler() (*PcPosShopHandler, error) {
	handler := &PcPosShopHandler{}

	// Khởi tạo PcPosShopService và xử lý error
	service, err := services.NewPcPosShopService()
	if err != nil {
		return nil, fmt.Errorf("failed to create pc pos shop service: %v", err)
	}
	handler.PcPosShopService = service
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
}
