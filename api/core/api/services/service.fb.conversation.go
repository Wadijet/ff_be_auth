package services

import (
	"context"
	"fmt"

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

// FindAllSortByApiUpdate tìm tất cả các FbConversation với phân trang sắp xếp theo thời gian cập nhật của dữ liệu API
// Hàm này cần để lấy dữ liệu cũ nhất để đồng bộ lại conversation mới
// Parameters:
//   - ctx: Context cho việc hủy bỏ hoặc timeout
//   - page: Trang hiện tại (bắt đầu từ 1)
//   - limit: Số lượng item trên mỗi trang
//   - filter: Điều kiện lọc (có thể là nil hoặc bson.M rỗng)
//
// Returns:
//   - *models.PaginateResult[models.FbConversation]: Kết quả phân trang với đầy đủ thông tin (page, limit, itemCount, items, total, totalPage)
//   - error: Lỗi nếu có
func (s *FbConversationService) FindAllSortByApiUpdate(ctx context.Context, page int64, limit int64, filter bson.M) (*models.PaginateResult[models.FbConversation], error) {
	// Tạo options với sort theo panCakeUpdatedAt giảm dần (cũ nhất trước)
	opts := options.Find().SetSort(bson.D{{Key: "panCakeUpdatedAt", Value: -1}})

	// Sử dụng FindWithPagination từ BaseServiceMongoImpl để có đầy đủ thông tin phân trang
	return s.BaseServiceMongoImpl.FindWithPagination(ctx, filter, page, limit, opts)
}
