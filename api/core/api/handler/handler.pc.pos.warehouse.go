package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// PcPosWarehouseHandler là cấu trúc xử lý các yêu cầu liên quan đến Pancake POS Warehouse cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.PcPosWarehouse
// - CreateInput: dto.PcPosWarehouseCreateInput
// - UpdateInput: dto.PcPosWarehouseCreateInput
type PcPosWarehouseHandler struct {
	BaseHandler[models.PcPosWarehouse, dto.PcPosWarehouseCreateInput, dto.PcPosWarehouseCreateInput]
	PcPosWarehouseService *services.PcPosWarehouseService
}

// NewPcPosWarehouseHandler khởi tạo một PcPosWarehouseHandler mới
// Returns:
//   - *PcPosWarehouseHandler: Instance mới của PcPosWarehouseHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewPcPosWarehouseHandler() (*PcPosWarehouseHandler, error) {
	handler := &PcPosWarehouseHandler{}

	// Khởi tạo PcPosWarehouseService và xử lý error
	service, err := services.NewPcPosWarehouseService()
	if err != nil {
		return nil, fmt.Errorf("failed to create pc pos warehouse service: %v", err)
	}
	handler.PcPosWarehouseService = service
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
}
