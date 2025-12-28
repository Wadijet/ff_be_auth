package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// NotificationSenderService là cấu trúc chứa các phương thức liên quan đến Notification Sender
type NotificationSenderService struct {
	*BaseServiceMongoImpl[models.NotificationChannelSender]
}

// NewNotificationSenderService tạo mới NotificationSenderService
func NewNotificationSenderService() (*NotificationSenderService, error) {
	collection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.NotificationSenders)
	if !exist {
		return nil, fmt.Errorf("failed to get notification_senders collection: %v", common.ErrNotFound)
	}

	return &NotificationSenderService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.NotificationChannelSender](collection),
	}, nil
}

// ✅ Các method InsertOne, DeleteById, UpdateById đã được xử lý bởi BaseServiceMongoImpl
// với cơ chế bảo vệ dữ liệu hệ thống chung (IsSystem)

