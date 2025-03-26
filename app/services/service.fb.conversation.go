package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "meta_commerce/app/models/mongodb"
	"meta_commerce/config"
	"meta_commerce/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FbConversationService là cấu trúc chứa các phương thức liên quan đến cuộc trò chuyện Facebook
type FbConversationService struct {
	*BaseServiceImpl[models.FbConversation]
}

// NewFbConversationService tạo mới FbConversationService
func NewFbConversationService(c *config.Configuration, db *mongo.Client) *FbConversationService {
	fbConversationCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.FbConvesations)).Collection(global.MongoDB_ColNames.FbConvesations)
	return &FbConversationService{
		BaseServiceImpl: NewBaseService[models.FbConversation](fbConversationCollection),
	}
}

// IsConversationIdExist kiểm tra ID cuộc trò chuyện có tồn tại hay không
func (s *FbConversationService) IsConversationIdExist(ctx context.Context, conversationId string) (bool, error) {
	filter := bson.M{"conversationId": conversationId}
	var conversation models.FbConversation
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ReviceData nhận data từ Facebook và lưu vào cơ sở dữ liệu
func (s *FbConversationService) ReviceData(ctx context.Context, input *models.FbConversationCreateInput) (*models.FbConversation, error) {
	if input.PanCakeData == nil {
		return nil, errors.New("ApiData is required")
	}

	// Lấy thông tin ConversationID từ ApiData đưa vào biến
	conversationId := input.PanCakeData["id"].(string)
	customerId := input.PanCakeData["customer_id"].(string)
	pancakeUpdatedAtStr := input.PanCakeData["updated_at"].(string)

	// Chuyển đổi thời gian từ string sang time.Time
	parsedTime, err := time.Parse("2006-01-02T15:04:05", pancakeUpdatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("lỗi phân tích thời gian: %v", err)
	}

	// Chuyển sang kiểu float64 (Unix timestamp dạng float64)
	pancakeUpdatedAt := int64(parsedTime.Unix())

	// Kiểm tra FbConversation đã tồn tại chưa
	exists, err := s.IsConversationIdExist(ctx, conversationId)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Tạo một FbConversation mới
		conversation := &models.FbConversation{
			PageId:           input.PageId,
			PageUsername:     input.PageUsername,
			PanCakeData:      input.PanCakeData,
			ConversationId:   conversationId,
			CustomerId:       customerId,
			PanCakeUpdatedAt: pancakeUpdatedAt,
			CreatedAt:        time.Now().Unix(),
			UpdatedAt:        time.Now().Unix(),
		}

		// Lưu FbConversation
		createdConversation, err := s.BaseServiceImpl.Create(ctx, *conversation)
		if err != nil {
			return nil, err
		}

		return &createdConversation, nil
	} else {
		// Lấy FbConversation hiện tại
		filter := bson.M{"conversationId": conversationId}
		conversation, err := s.BaseServiceImpl.FindOneByFilter(ctx, filter, nil)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin mới
		conversation.PanCakeData = input.PanCakeData
		conversation.PageId = input.PageId
		conversation.PageUsername = input.PageUsername
		conversation.ConversationId = conversationId
		conversation.CustomerId = customerId
		conversation.PanCakeUpdatedAt = pancakeUpdatedAt
		conversation.UpdatedAt = time.Now().Unix()

		// Cập nhật FbConversation
		updatedConversation, err := s.BaseServiceImpl.Update(ctx, conversation.ID, conversation)
		if err != nil {
			return nil, err
		}

		return &updatedConversation, nil
	}
}

// FindAllSortByApiUpdate tìm tất cả các FbConversation với phân trang sắp xếp theo thời gian cập nhật của dữ liệu API
func (s *FbConversationService) FindAllSortByApiUpdate(ctx context.Context, page int64, limit int64, filter bson.M) ([]models.FbConversation, error) {
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{"panCakeUpdatedAt", -1}})

	cursor, err := s.BaseServiceImpl.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.FbConversation
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
