package handler

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// NotificationSenderHandler xử lý các request liên quan đến Notification Sender
type NotificationSenderHandler struct {
	BaseHandler[models.NotificationChannelSender, models.NotificationChannelSender, models.NotificationChannelSender]
}

// NewNotificationSenderHandler tạo mới NotificationSenderHandler
func NewNotificationSenderHandler() (*NotificationSenderHandler, error) {
	senderService, err := services.NewNotificationSenderService()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification sender service: %v", err)
	}

	handler := &NotificationSenderHandler{}
	handler.BaseService = senderService

	// Khởi tạo filterOptions với giá trị mặc định
	handler.filterOptions = FilterOptions{
		DeniedFields: []string{
			"smtpPassword",
			"botToken",
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

