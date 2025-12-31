package services

import (
	"context"
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationChannelService là cấu trúc chứa các phương thức liên quan đến Notification Channel
type NotificationChannelService struct {
	*BaseServiceMongoImpl[models.NotificationChannel]
}

// NewNotificationChannelService tạo mới NotificationChannelService
func NewNotificationChannelService() (*NotificationChannelService, error) {
	collection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.NotificationChannels)
	if !exist {
		return nil, fmt.Errorf("failed to get notification_channels collection: %v", common.ErrNotFound)
	}

	return &NotificationChannelService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.NotificationChannel](collection),
	}, nil
}

// FindByOrganizationID tìm tất cả channels của một organization, có thể filter theo channelTypes
func (s *NotificationChannelService) FindByOrganizationID(ctx context.Context, orgID primitive.ObjectID, channelTypes []string) ([]models.NotificationChannel, error) {
	filter := bson.M{
		"ownerOrganizationId": orgID, // Phân quyền dữ liệu
		"isActive":           true,
	}

	// Filter theo channelTypes nếu có
	if len(channelTypes) > 0 {
		filter["channelType"] = bson.M{"$in": channelTypes}
	}

	opts := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := s.BaseServiceMongoImpl.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var channels []models.NotificationChannel
	if err := cursor.All(ctx, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

// ✅ Các method InsertOne, DeleteById, UpdateById đã được xử lý bởi BaseServiceMongoImpl
// với cơ chế bảo vệ dữ liệu hệ thống chung (IsSystem)

