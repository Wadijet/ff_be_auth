package handler

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// OrganizationHandler xử lý các request liên quan đến Organization
type OrganizationHandler struct {
	BaseHandler[models.Organization, models.Organization, models.Organization]
	OrganizationService *services.OrganizationService
}

// NewOrganizationHandler tạo mới OrganizationHandler
func NewOrganizationHandler() (*OrganizationHandler, error) {
	organizationService, err := services.NewOrganizationService()
	if err != nil {
		return nil, fmt.Errorf("failed to create organization service: %v", err)
	}

	handler := &OrganizationHandler{
		OrganizationService: organizationService,
	}
	handler.BaseService = handler.OrganizationService
	
	// Khởi tạo filterOptions với giá trị mặc định
	handler.filterOptions = FilterOptions{
		DeniedFields: []string{
			"password",
			"token",
			"secret",
			"key",
			"hash",
		},
		AllowedOperators: []string{
			"$eq",
			"$gt",
			"$gte",
			"$lt",
			"$lte",
			"$in",
			"$nin",
			"$exists",
		},
		MaxFields: 10,
	}
	
	return handler, nil
}

// Các method CRUD đã được kế thừa từ BaseHandler

