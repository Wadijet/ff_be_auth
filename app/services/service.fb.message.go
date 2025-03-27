package services

import (
	"context"
	"errors"
	"time"

	models "meta_commerce/app/models/mongodb"
	"meta_commerce/config"
	"meta_commerce/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FbMessageService là cấu trúc chứa các phương thức liên quan đến tin nhắn Facebook
type FbMessageService struct {
	*BaseServiceImpl[models.FbMessage]
}

// NewFbMessageService tạo mới FbMessageService
func NewFbMessageService(c *config.Configuration, db *mongo.Client) *FbMessageService {
	fbMessageCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.FbMessages)).Collection(global.MongoDB_ColNames.FbMessages)
	return &FbMessageService{
		BaseServiceImpl: NewBaseService[models.FbMessage](fbMessageCollection),
	}
}

// IsMessageExist kiểm tra tin nhắn có tồn tại hay không
func (s *FbMessageService) IsMessageExist(ctx context.Context, conversationId string, customerId string) (bool, error) {
	filter := bson.M{"conversationId": conversationId, "customerId": customerId}
	var message models.FbMessage
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Upsert nhận data từ Facebook và lưu vào cơ sở dữ liệu
func (s *FbMessageService) Upsert(ctx context.Context, input *models.FbMessageCreateInput) (*models.FbMessage, error) {
	if input.PanCakeData == nil {
		return nil, errors.New("ApiData is required")
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
	upsertedMessage, err := s.BaseServiceImpl.Upsert(ctx, filter, *message)
	if err != nil {
		return nil, err
	}

	return &upsertedMessage, nil
}

// FindOneByConversationID tìm một FbMessage theo ConversationID
func (s *FbMessageService) FindOneByConversationID(ctx context.Context, conversationID string) (models.FbMessage, error) {
	filter := bson.M{"conversationId": conversationID}
	var message models.FbMessage
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return message, errors.New("message not found")
		}
		return message, err
	}
	return message, nil
}

// FindAll tìm tất cả các FbMessage với phân trang
func (s *FbMessageService) FindAll(ctx context.Context, page int64, limit int64) ([]models.FbMessage, error) {
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{"updatedAt", 1}})

	cursor, err := s.BaseServiceImpl.collection.Find(ctx, nil, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.FbMessage
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
