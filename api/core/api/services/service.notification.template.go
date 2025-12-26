package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// NotificationTemplateService là cấu trúc chứa các phương thức liên quan đến Notification Template
type NotificationTemplateService struct {
	*BaseServiceMongoImpl[models.NotificationTemplate]
}

// NewNotificationTemplateService tạo mới NotificationTemplateService
func NewNotificationTemplateService() (*NotificationTemplateService, error) {
	collection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.NotificationTemplates)
	if !exist {
		return nil, fmt.Errorf("failed to get notification_templates collection: %v", common.ErrNotFound)
	}

	return &NotificationTemplateService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.NotificationTemplate](collection),
	}, nil
}

