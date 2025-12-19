package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// PcPosCustomerHandler xử lý các route liên quan đến Pancake POS Customer
type PcPosCustomerHandler struct {
	BaseHandler[models.PcPosCustomer, dto.PcPosCustomerCreateInput, dto.PcPosCustomerCreateInput]
	PcPosCustomerService *services.PcPosCustomerService
}

// NewPcPosCustomerHandler tạo một instance mới của PcPosCustomerHandler
func NewPcPosCustomerHandler() (*PcPosCustomerHandler, error) {
	handler := &PcPosCustomerHandler{}

	// Khởi tạo PcPosCustomerService
	service, err := services.NewPcPosCustomerService()
	if err != nil {
		return nil, fmt.Errorf("failed to create pc pos customer service: %v", err)
	}
	handler.PcPosCustomerService = service

	// Gán BaseServiceMongoImpl cho BaseHandler để các method CRUD cơ bản hoạt động
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
}
