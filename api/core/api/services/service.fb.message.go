package services

import (
	"context"
	"fmt"
	"time"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FbMessageService là cấu trúc chứa các phương thức liên quan đến tin nhắn Facebook
type FbMessageService struct {
	*BaseServiceMongoImpl[models.FbMessage]
	fbPageService        *FbPageService
	fbMessageItemService *FbMessageItemService
}

// NewFbMessageService tạo mới FbMessageService
func NewFbMessageService() (*FbMessageService, error) {
	fbMessageCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.FbMessages)
	if !exist {
		return nil, fmt.Errorf("failed to get fb_messages collection: %v", common.ErrNotFound)
	}

	fbPageService, err := NewFbPageService()
	if err != nil {
		return nil, fmt.Errorf("failed to create fb_page service: %v", err)
	}

	fbMessageItemService, err := NewFbMessageItemService()
	if err != nil {
		return nil, fmt.Errorf("failed to create fb_message_item service: %v", err)
	}

	return &FbMessageService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.FbMessage](fbMessageCollection),
		fbPageService:        fbPageService,
		fbMessageItemService: fbMessageItemService,
	}, nil
}

// IsMessageExist kiểm tra tin nhắn có tồn tại hay không
func (s *FbMessageService) IsMessageExist(ctx context.Context, conversationId string, customerId string) (bool, error) {
	filter := bson.M{"conversationId": conversationId, "customerId": customerId}
	var message models.FbMessage
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, common.ConvertMongoError(err)
	}
	return true, nil
}

// FindOneByConversationID tìm một FbMessage theo ConversationID
func (s *FbMessageService) FindOneByConversationID(ctx context.Context, conversationID string) (models.FbMessage, error) {
	filter := bson.M{"conversationId": conversationID}
	var message models.FbMessage
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return message, common.ErrNotFound
		}
		return message, common.ConvertMongoError(err)
	}
	return message, nil
}

// FindAll tìm tất cả các FbMessage với phân trang
func (s *FbMessageService) FindAll(ctx context.Context, page int64, limit int64) ([]models.FbMessage, error) {
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{Key: "updatedAt", Value: 1}})

	cursor, err := s.BaseServiceMongoImpl.collection.Find(ctx, nil, opts)
	if err != nil {
		return nil, common.ConvertMongoError(err)
	}
	defer cursor.Close(ctx)

	var results []models.FbMessage
	if err = cursor.All(ctx, &results); err != nil {
		return nil, common.ConvertMongoError(err)
	}

	return results, nil
}

// UpsertMessages xử lý upsert messages từ panCakeData
// Logic nội bộ: Tách messages[] ra khỏi panCakeData và lưu vào 2 collections
// - Metadata (panCakeData không có messages[]) → fb_messages
// - Messages (từng message riêng lẻ) → fb_message_items
func (s *FbMessageService) UpsertMessages(
	ctx context.Context,
	conversationId string,
	pageId string,
	pageUsername string,
	customerId string,
	panCakeData map[string]interface{}, // PanCakeData đầy đủ (bao gồm messages[])
	hasMore bool,
) (models.FbMessage, error) {
	now := time.Now().UnixMilli()

	// 1. Tách messages[] ra khỏi panCakeData
	messages, _ := panCakeData["messages"].([]interface{})

	// Tạo metadata panCakeData: GIỮ LẠI TẤT CẢ các field khác ngoài "messages"
	// Ví dụ: conv_from, read_watermarks, activities, ad_clicks, is_banned, notes, etc.
	// Lưu ý: Range trên nil map là safe trong Go, không cần check nil
	metadataPanCakeData := make(map[string]interface{})
	for k, v := range panCakeData {
		// Chỉ bỏ qua field "messages", giữ lại tất cả các field khác
		if k != "messages" {
			metadataPanCakeData[k] = v
		}
	}

	// Đảm bảo metadataPanCakeData có conversation_id để extract tag hoạt động
	// (extract tag sẽ tìm conversationId từ panCakeData.conversation_id)
	// Nếu panCakeData không có conversation_id, thêm vào để extract có thể hoạt động
	// Lưu ý: Chỉ thêm nếu chưa có, không ghi đè giá trị hiện có
	if _, exists := metadataPanCakeData["conversation_id"]; !exists {
		metadataPanCakeData["conversation_id"] = conversationId
	}

	// 2. Upsert metadata vào fb_messages với merge panCakeData
	filter := bson.M{"conversationId": conversationId}

	// Kiểm tra xem document đã tồn tại chưa để merge panCakeData
	var existingDoc models.FbMessage
	err := s.collection.FindOne(ctx, filter).Decode(&existingDoc)
	exists := err == nil

	// Merge panCakeData: Giữ lại các field cũ, update các field mới
	mergedPanCakeData := make(map[string]interface{})

	if exists && existingDoc.PanCakeData != nil {
		// Copy tất cả field từ panCakeData cũ (GIỮ LẠI TẤT CẢ dữ liệu cũ)
		for k, v := range existingDoc.PanCakeData {
			mergedPanCakeData[k] = v
		}
		logrus.WithFields(logrus.Fields{
			"conversationId": conversationId,
			"existingFields": len(existingDoc.PanCakeData),
		}).Debug("UpsertMessages: Đã copy panCakeData cũ, số field:", len(existingDoc.PanCakeData))
	}

	// Update/Thêm các field từ metadataPanCakeData mới (ưu tiên dữ liệu mới)
	// Lưu ý: Chỉ update các field có trong request mới, giữ nguyên các field cũ không có trong request
	for k, v := range metadataPanCakeData {
		// Nếu là nested map, merge sâu hơn để giữ lại các field con
		if existingMap, ok := mergedPanCakeData[k].(map[string]interface{}); ok {
			if newMap, ok := v.(map[string]interface{}); ok {
				// Merge nested map: Giữ lại field con cũ, update field con mới
				for nk, nv := range newMap {
					existingMap[nk] = nv
				}
				mergedPanCakeData[k] = existingMap
			} else {
				// Không phải map, ghi đè (update field)
				mergedPanCakeData[k] = v
			}
		} else {
			// Field mới hoặc không phải nested map, update/ghi đè
			mergedPanCakeData[k] = v
		}
	}

	// Đảm bảo không có messages[] trong merged data
	delete(mergedPanCakeData, "messages")

	logrus.WithFields(logrus.Fields{
		"conversationId":   conversationId,
		"mergedFields":     len(mergedPanCakeData),
		"newFields":        len(metadataPanCakeData),
		"mergedFieldNames": getMapKeys(mergedPanCakeData),
	}).Debug("UpsertMessages: Sau khi merge panCakeData")

	// Tạo update operation với merge panCakeData
	update := bson.M{
		"$set": bson.M{
			"pageId":       pageId,
			"pageUsername": pageUsername,
			"customerId":   customerId,
			"panCakeData":  mergedPanCakeData, // Merge panCakeData, không ghi đè
			"lastSyncedAt": now,
			"hasMore":      hasMore,
			"updatedAt":    now,
		},
		"$setOnInsert": bson.M{
			"createdAt": now,
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var metadataResult models.FbMessage
	err = s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&metadataResult)
	if err != nil {
		return metadataResult, common.ConvertMongoError(err)
	}

	// 3. Upsert messages vào fb_message_items (nếu có messages)
	if len(messages) > 0 {
		_, err = s.fbMessageItemService.UpsertMessages(ctx, conversationId, messages)
		if err != nil {
			return metadataResult, fmt.Errorf("failed to upsert messages: %v", err)
		}
	}

	// 4. Cập nhật totalMessages
	totalMessages, err := s.fbMessageItemService.CountByConversationId(ctx, conversationId)
	if err != nil {
		return metadataResult, fmt.Errorf("failed to count messages: %v", err)
	}

	// Update totalMessages (dùng lại biến update và opts đã khai báo)
	update = bson.M{
		"$set": bson.M{
			"totalMessages": totalMessages,
			"updatedAt":     now,
		},
	}

	opts = options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var updated models.FbMessage
	err = s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		return metadataResult, common.ConvertMongoError(err)
	}

	return updated, nil
}
