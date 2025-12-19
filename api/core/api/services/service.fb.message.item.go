package services

import (
	"context"
	"fmt"
	"time"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FbMessageItemService là cấu trúc chứa các phương thức liên quan đến message items Facebook
type FbMessageItemService struct {
	*BaseServiceMongoImpl[models.FbMessageItem]
}

// NewFbMessageItemService tạo mới FbMessageItemService
func NewFbMessageItemService() (*FbMessageItemService, error) {
	fbMessageItemCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.FbMessageItems)
	if !exist {
		return nil, fmt.Errorf("failed to get fb_message_items collection: %v", common.ErrNotFound)
	}

	return &FbMessageItemService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.FbMessageItem](fbMessageItemCollection),
	}, nil
}

// UpsertMessages upsert nhiều messages vào collection (mỗi message là 1 document)
// Tự động tránh duplicate theo messageId
func (s *FbMessageItemService) UpsertMessages(
	ctx context.Context,
	conversationId string,
	messages []interface{},
) (int, error) {
	if len(messages) == 0 {
		return 0, nil
	}

	var operations []mongo.WriteModel
	now := time.Now().UnixMilli()

	for _, msg := range messages {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract messageId
		messageId, ok := msgMap["id"].(string)
		if !ok || messageId == "" {
			continue
		}

		// Extract inserted_at và convert sang Unix timestamp
		var insertedAt int64 = 0
		if insertedAtStr, ok := msgMap["inserted_at"].(string); ok {
			if t, err := time.Parse("2006-01-02T15:04:05.000000", insertedAtStr); err == nil {
				insertedAt = t.Unix()
			}
		}

		// Tạo document map trực tiếp (không dùng extract tag vì đã có dữ liệu sẵn)
		docMap := bson.M{
			"conversationId": conversationId,
			"messageId":      messageId,
			"messageData":    msgMap,
			"insertedAt":     insertedAt,
			"updatedAt":      now,
		}

		// Tạo upsert operation
		filter := bson.M{"messageId": messageId}
		update := bson.M{
			"$set": docMap,
			"$setOnInsert": bson.M{
				"createdAt": now,
			},
		}

		operation := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, operation)
	}

	if len(operations) == 0 {
		return 0, nil
	}

	// Bulk write
	opts := options.BulkWrite().SetOrdered(false)
	result, err := s.collection.BulkWrite(ctx, operations, opts)
	if err != nil {
		return 0, common.ConvertMongoError(err)
	}

	return int(result.UpsertedCount + result.ModifiedCount), nil
}

// FindByConversationId tìm tất cả messages của một conversation với phân trang
func (s *FbMessageItemService) FindByConversationId(
	ctx context.Context,
	conversationId string,
	page int64,
	limit int64,
) ([]models.FbMessageItem, int64, error) {
	filter := bson.M{"conversationId": conversationId}

	// Count total
	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, common.ConvertMongoError(err)
	}

	// Find with pagination
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{Key: "insertedAt", Value: -1}}) // Sort từ mới đến cũ

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, common.ConvertMongoError(err)
	}
	defer cursor.Close(ctx)

	var results []models.FbMessageItem
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, common.ConvertMongoError(err)
	}

	return results, total, nil
}

// CountByConversationId đếm số lượng messages của một conversation
func (s *FbMessageItemService) CountByConversationId(ctx context.Context, conversationId string) (int64, error) {
	filter := bson.M{"conversationId": conversationId}
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, common.ConvertMongoError(err)
	}
	return count, nil
}
