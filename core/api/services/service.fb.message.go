package services

import (
	"context"
	"fmt"
	"time"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FbMessageService là cấu trúc chứa các phương thức liên quan đến tin nhắn Facebook
type FbMessageService struct {
	*BaseServiceMongoImpl[models.FbMessage]
	fbPageService *FbPageService
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

	return &FbMessageService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.FbMessage](fbMessageCollection),
		fbPageService:        fbPageService,
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

// Upsert nhận data từ Facebook và lưu vào cơ sở dữ liệu
func (s *FbMessageService) Upsert(ctx context.Context, input *models.FbMessageCreateInput) (*models.FbMessage, error) {
	if input.PanCakeData == nil {
		return nil, common.ErrInvalidInput
	}

	// Lấy thông tin MessageId từ ApiData đưa vào biến
	conversationId := input.PanCakeData["conversation_id"].(string)

	// Tạo filter để tìm kiếm document
	filter := bson.M{
		"conversationId": conversationId,
		"customerId":     input.CustomerId,
	}

	// Tạo message mới
	message := &models.FbMessage{
		ID:             primitive.NewObjectID(),
		PageId:         input.PageId,
		PageUsername:   input.PageUsername,
		PanCakeData:    input.PanCakeData,
		CustomerId:     input.CustomerId,
		ConversationId: conversationId,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}

	// Sử dụng Upsert từ base.service
	upsertedMessage, err := s.BaseServiceMongoImpl.Upsert(ctx, filter, *message)
	if err != nil {
		return nil, common.ConvertMongoError(err)
	}

	return &upsertedMessage, nil
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
		SetSort(bson.D{{"updatedAt", 1}})

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
