package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FbConversationService là cấu trúc chứa các phương thức liên quan đến Facebook conversation
type FbConversationService struct {
	*BaseServiceMongoImpl[models.FbConversation]
	fbPageService    *FbPageService
	fbMessageService *FbMessageService
}

// NewFbConversationService tạo mới FbConversationService
func NewFbConversationService() (*FbConversationService, error) {
	fbConversationCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.FbConvesations)
	if !exist {
		return nil, fmt.Errorf("failed to get fb_conversations collection: %v", common.ErrNotFound)
	}

	fbPageService, err := NewFbPageService()
	if err != nil {
		return nil, fmt.Errorf("failed to create fb_page service: %v", err)
	}

	fbMessageService, err := NewFbMessageService()
	if err != nil {
		return nil, fmt.Errorf("failed to create fb_message service: %v", err)
	}

	return &FbConversationService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.FbConversation](fbConversationCollection),
		fbPageService:        fbPageService,
		fbMessageService:     fbMessageService,
	}, nil
}

// IsConversationIdExist kiểm tra ID cuộc trò chuyện có tồn tại hay không
func (s *FbConversationService) IsConversationIdExist(ctx context.Context, conversationId string) (bool, error) {
	filter := bson.M{"conversationId": conversationId}
	var conversation models.FbConversation
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Upsert thực hiện thao tác upsert dữ liệu conversation từ Facebook vào cơ sở dữ liệu
func (s *FbConversationService) Upsert(ctx context.Context, input *dto.FbConversationCreateInput) (*models.FbConversation, error) {
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

	// Tạo filter để tìm conversation theo conversationId
	filter := bson.M{"conversationId": conversationId}

	// Tạo conversation mới với dữ liệu từ input
	conversation := &models.FbConversation{
		PageId:           input.PageId,
		PageUsername:     input.PageUsername,
		PanCakeData:      input.PanCakeData,
		ConversationId:   conversationId,
		CustomerId:       customerId,
		PanCakeUpdatedAt: pancakeUpdatedAt,
	}

	// Sử dụng Upsert để tạo mới hoặc cập nhật conversation
	upsertedConversation, err := s.BaseServiceMongoImpl.Upsert(ctx, filter, *conversation)
	if err != nil {
		return nil, err
	}

	return &upsertedConversation, nil
}

// FindAllSortByApiUpdate tìm tất cả các FbConversation với phân trang sắp xếp theo thời gian cập nhật của dữ liệu API
// Hàm này cần để lấy dữ liệu cũ nhất để đồng bộ lại conversation mới
func (s *FbConversationService) FindAllSortByApiUpdate(ctx context.Context, page int64, limit int64, filter bson.M) ([]models.FbConversation, error) {
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{"panCakeUpdatedAt", -1}})

	cursor, err := s.BaseServiceMongoImpl.collection.Find(ctx, filter, opts)
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
