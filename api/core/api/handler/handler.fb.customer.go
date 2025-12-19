package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// FbCustomerHandler xử lý các route liên quan đến Facebook Customer
type FbCustomerHandler struct {
	BaseHandler[models.FbCustomer, dto.FbCustomerCreateInput, dto.FbCustomerCreateInput]
	FbCustomerService *services.FbCustomerService
}

// NewFbCustomerHandler tạo một instance mới của FbCustomerHandler
func NewFbCustomerHandler() (*FbCustomerHandler, error) {
	handler := &FbCustomerHandler{}

	// Khởi tạo FbCustomerService
	service, err := services.NewFbCustomerService()
	if err != nil {
		return nil, fmt.Errorf("failed to create fb customer service: %v", err)
	}
	handler.FbCustomerService = service

	// Gán BaseServiceMongoImpl cho BaseHandler để các method CRUD cơ bản hoạt động
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
}
