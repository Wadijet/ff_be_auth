package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// NotificationHistoryService là cấu trúc chứa các phương thức liên quan đến Notification History
type NotificationHistoryService struct {
	*BaseServiceMongoImpl[models.NotificationHistory]
}

// NewNotificationHistoryService tạo mới NotificationHistoryService
func NewNotificationHistoryService() (*NotificationHistoryService, error) {
	collection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.NotificationHistory)
	if !exist {
		return nil, fmt.Errorf("failed to get notification_history collection: %v", common.ErrNotFound)
	}

	return &NotificationHistoryService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.NotificationHistory](collection),
	}, nil
}

