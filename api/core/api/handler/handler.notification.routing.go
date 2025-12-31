package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// NotificationRoutingHandler xử lý các request liên quan đến Notification Routing Rule
type NotificationRoutingHandler struct {
	BaseHandler[models.NotificationRoutingRule, dto.NotificationRoutingRuleCreateInput, dto.NotificationRoutingRuleUpdateInput]
}

// NewNotificationRoutingHandler tạo mới NotificationRoutingHandler
func NewNotificationRoutingHandler() (*NotificationRoutingHandler, error) {
	routingService, err := services.NewNotificationRoutingService()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification routing service: %v", err)
	}

	baseHandler := NewBaseHandler[models.NotificationRoutingRule, dto.NotificationRoutingRuleCreateInput, dto.NotificationRoutingRuleUpdateInput](routingService)
	handler := &NotificationRoutingHandler{
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

