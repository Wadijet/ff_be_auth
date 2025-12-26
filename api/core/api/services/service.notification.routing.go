package services

import (
	"context"
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationRoutingService là cấu trúc chứa các phương thức liên quan đến Notification Routing Rule
type NotificationRoutingService struct {
	*BaseServiceMongoImpl[models.NotificationRoutingRule]
}

// NewNotificationRoutingService tạo mới NotificationRoutingService
func NewNotificationRoutingService() (*NotificationRoutingService, error) {
	collection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.NotificationRoutingRules)
	if !exist {
		return nil, fmt.Errorf("failed to get notification_routing_rules collection: %v", common.ErrNotFound)
	}

	return &NotificationRoutingService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.NotificationRoutingRule](collection),
	}, nil
}

// FindByEventType tìm tất cả rules theo eventType và isActive = true
func (s *NotificationRoutingService) FindByEventType(ctx context.Context, eventType string) ([]models.NotificationRoutingRule, error) {
	filter := bson.M{
		"eventType": eventType,
		"isActive":  true,
	}

	opts := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := s.BaseServiceMongoImpl.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rules []models.NotificationRoutingRule
	if err := cursor.All(ctx, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

