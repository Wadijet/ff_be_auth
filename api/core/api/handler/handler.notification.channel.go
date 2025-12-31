package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// NotificationChannelHandler xử lý các request liên quan đến Notification Channel
type NotificationChannelHandler struct {
	BaseHandler[models.NotificationChannel, dto.NotificationChannelCreateInput, dto.NotificationChannelUpdateInput]
}

// NewNotificationChannelHandler tạo mới NotificationChannelHandler
func NewNotificationChannelHandler() (*NotificationChannelHandler, error) {
	channelService, err := services.NewNotificationChannelService()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification channel service: %v", err)
	}

	baseHandler := NewBaseHandler[models.NotificationChannel, dto.NotificationChannelCreateInput, dto.NotificationChannelUpdateInput](channelService)
	handler := &NotificationChannelHandler{
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

