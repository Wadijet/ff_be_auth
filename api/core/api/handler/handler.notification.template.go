package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// NotificationTemplateHandler xử lý các request liên quan đến Notification Template
type NotificationTemplateHandler struct {
	BaseHandler[models.NotificationTemplate, dto.NotificationTemplateCreateInput, dto.NotificationTemplateUpdateInput]
}

// NewNotificationTemplateHandler tạo mới NotificationTemplateHandler
func NewNotificationTemplateHandler() (*NotificationTemplateHandler, error) {
	templateService, err := services.NewNotificationTemplateService()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification template service: %v", err)
	}

	baseHandler := NewBaseHandler[models.NotificationTemplate, dto.NotificationTemplateCreateInput, dto.NotificationTemplateUpdateInput](templateService)
	handler := &NotificationTemplateHandler{
		BaseHandler: *baseHandler,
	}

	// Khởi tạo filterOptions với giá trị mặc định
	handler.filterOptions = FilterOptions{
		DeniedFields: []string{},
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

