package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// CustomerHandler xử lý các route liên quan đến Customer
type CustomerHandler struct {
	BaseHandler[models.Customer, dto.CustomerCreateInput, dto.CustomerCreateInput]
	CustomerService *services.CustomerService
}

// NewCustomerHandler tạo một instance mới của CustomerHandler
func NewCustomerHandler() (*CustomerHandler, error) {
	handler := &CustomerHandler{}

	// Khởi tạo CustomerService
	service, err := services.NewCustomerService()
	if err != nil {
		return nil, fmt.Errorf("failed to create customer service: %v", err)
	}
	handler.CustomerService = service

	// Gán BaseServiceMongoImpl cho BaseHandler để các method CRUD cơ bản hoạt động
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
}
