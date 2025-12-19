package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// PcPosVariationHandler là cấu trúc xử lý các yêu cầu liên quan đến Pancake POS Variation cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.PcPosVariation
// - CreateInput: dto.PcPosVariationCreateInput
// - UpdateInput: dto.PcPosVariationCreateInput
type PcPosVariationHandler struct {
	BaseHandler[models.PcPosVariation, dto.PcPosVariationCreateInput, dto.PcPosVariationCreateInput]
	PcPosVariationService *services.PcPosVariationService
}

// NewPcPosVariationHandler khởi tạo một PcPosVariationHandler mới
// Returns:
//   - *PcPosVariationHandler: Instance mới của PcPosVariationHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewPcPosVariationHandler() (*PcPosVariationHandler, error) {
	handler := &PcPosVariationHandler{}

	// Khởi tạo PcPosVariationService và xử lý error
	service, err := services.NewPcPosVariationService()
	if err != nil {
		return nil, fmt.Errorf("failed to create pc pos variation service: %v", err)
	}
	handler.PcPosVariationService = service
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
}
