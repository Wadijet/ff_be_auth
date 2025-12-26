package handler

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// NotificationHistoryHandler xử lý các request liên quan đến Notification History
type NotificationHistoryHandler struct {
	BaseHandler[models.NotificationHistory, models.NotificationHistory, models.NotificationHistory]
}

// NewNotificationHistoryHandler tạo mới NotificationHistoryHandler
func NewNotificationHistoryHandler() (*NotificationHistoryHandler, error) {
	historyService, err := services.NewNotificationHistoryService()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification history service: %v", err)
	}

	handler := &NotificationHistoryHandler{}
	handler.BaseService = historyService

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

