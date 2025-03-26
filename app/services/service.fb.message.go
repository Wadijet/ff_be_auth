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

// ReviceData nhận data từ Facebook và lưu vào cơ sở dữ liệu
func (s *FbMessageService) ReviceData(ctx context.Context, input *models.FbMessageCreateInput) (*models.FbMessage, error) {
	if input.PanCakeData == nil {
		return nil, errors.New("ApiData is required")
	}

	// Lấy thông tin MessageId từ ApiData đưa vào biến
	conversationId := input.PanCakeData["conversation_id"].(string)

	// Kiểm tra FbMessage đã tồn tại chưa
	exists, err := s.IsMessageExist(ctx, conversationId, input.CustomerId)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Tạo một FbMessage mới
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

		// Lưu FbMessage
		createdMessage, err := s.BaseServiceImpl.Create(ctx, *message)
		if err != nil {
			return nil, err
		}

		return &createdMessage, nil
	} else {
		// Lấy FbMessage hiện tại
		filter := bson.M{"conversationId": conversationId}
		message, err := s.BaseServiceImpl.FindOneByFilter(ctx, filter, nil)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin mới
		message.PanCakeData = input.PanCakeData
		message.PageId = input.PageId
		message.PageUsername = input.PageUsername
		message.ConversationId = conversationId
		message.CustomerId = input.CustomerId
		message.UpdatedAt = time.Now().Unix()

		// Cập nhật FbMessage
		updatedMessage, err := s.BaseServiceImpl.Update(ctx, message.ID, message)
		if err != nil {
			return nil, err
		}

		return &updatedMessage, nil
	}
}

// FindOneById tìm một FbMessage theo ID
func (s *FbMessageService) FindOneById(ctx context.Context, id primitive.ObjectID) (models.FbMessage, error) {
	return s.BaseServiceImpl.FindOne(ctx, id)
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
