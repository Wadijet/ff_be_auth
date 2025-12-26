package services

import (
	"context"
	"fmt"
	"time"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationQueueService là cấu trúc chứa các phương thức liên quan đến Notification Queue
type NotificationQueueService struct {
	*BaseServiceMongoImpl[models.NotificationQueueItem]
}

// NewNotificationQueueService tạo mới NotificationQueueService
func NewNotificationQueueService() (*NotificationQueueService, error) {
	collection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.NotificationQueue)
	if !exist {
		return nil, fmt.Errorf("failed to get notification_queue collection: %v", common.ErrNotFound)
	}

	return &NotificationQueueService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.NotificationQueueItem](collection),
	}, nil
}

// FindPending tìm các items có status="pending" và nextRetryAt <= now (hoặc null)
func (s *NotificationQueueService) FindPending(ctx context.Context, limit int) ([]models.NotificationQueueItem, error) {
	now := time.Now().Unix()
	filter := bson.M{
		"status": bson.M{"$in": []string{"pending"}},
		"$or": []bson.M{
			{"nextRetryAt": bson.M{"$lte": now}},
			{"nextRetryAt": nil},
		},
	}

	opts := options.Find().
		SetSort(bson.M{"createdAt": 1}).
		SetLimit(int64(limit))

	cursor, err := s.BaseServiceMongoImpl.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []models.NotificationQueueItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateStatus cập nhật status cho nhiều items
func (s *NotificationQueueService) UpdateStatus(ctx context.Context, ids []interface{}, status string) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now().Unix()}}

	_, err := s.BaseServiceMongoImpl.collection.UpdateMany(ctx, filter, update)
	return err
}

