package services

import (
	"context"
	"fmt"
	"time"

	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/utility"
	"meta_commerce/registry"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FbPostService là cấu trúc chứa các phương thức liên quan đến bài viết Facebook
type FbPostService struct {
	*BaseServiceMongoImpl[models.FbPost]
	fbPageService *FbPageService
}

// NewFbPostService tạo mới FbPostService
func NewFbPostService() (*FbPostService, error) {
	fbPostCollection, err := registry.Collections.MustGet(global.MongoDB_ColNames.FbPosts)
	if err != nil {
		return nil, fmt.Errorf("failed to get fb_posts collection: %v", err)
	}

	fbPageService, err := NewFbPageService()
	if err != nil {
		return nil, fmt.Errorf("failed to create fb_page service: %v", err)
	}

	return &FbPostService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.FbPost](fbPostCollection),
		fbPageService:        fbPageService,
	}, nil
}

// IsPostExist kiểm tra bài viết có tồn tại hay không
func (s *FbPostService) IsPostExist(ctx context.Context, postId string) (bool, error) {
	filter := bson.M{"postId": postId}
	var post models.FbPost
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, utility.ConvertMongoError(err)
	}
	return true, nil
}

// ReviceData nhận data từ Facebook và lưu vào cơ sở dữ liệu
func (s *FbPostService) ReviceData(ctx context.Context, input *models.FbPostCreateInput) (*models.FbPost, error) {
	if input.PanCakeData == nil {
		return nil, utility.ErrInvalidInput
	}

	// Lấy thông tin PostId từ ApiData đưa vào biến
	pageId := input.PanCakeData["page_id"].(string)
	postId := input.PanCakeData["id"].(string)

	// Kiểm tra FbPost đã tồn tại chưa
	exists, err := s.IsPostExist(ctx, postId)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Tạo một FbPost mới
		post := &models.FbPost{
			ID:          primitive.NewObjectID(),
			PageId:      pageId,
			PostId:      postId,
			PanCakeData: input.PanCakeData,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		}

		// Lưu FbPost
		createdPost, err := s.BaseServiceMongoImpl.InsertOne(ctx, *post)
		if err != nil {
			return nil, utility.ConvertMongoError(err)
		}

		return &createdPost, nil
	} else {
		filter := bson.M{"postId": postId}
		// Lấy FbPost hiện tại
		post, err := s.BaseServiceMongoImpl.FindOne(ctx, filter, nil)
		if err != nil {
			return nil, utility.ConvertMongoError(err)
		}

		// Cập nhật thông tin mới
		post.PanCakeData = input.PanCakeData
		post.UpdatedAt = time.Now().Unix()

		// Cập nhật FbPost
		updatedPost, err := s.BaseServiceMongoImpl.UpdateById(ctx, post.ID, post)
		if err != nil {
			return nil, utility.ConvertMongoError(err)
		}

		return &updatedPost, nil
	}
}

// FindOneByPostID tìm một FbPost theo PostID
func (s *FbPostService) FindOneByPostID(ctx context.Context, postID string) (models.FbPost, error) {
	filter := bson.M{"postId": postID}
	var post models.FbPost
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return post, utility.ErrNotFound
		}
		return post, utility.ConvertMongoError(err)
	}
	return post, nil
}

// FindAll tìm kiếm tất cả bài viết
func (s *FbPostService) FindAll(ctx context.Context, page int64, limit int64) ([]models.FbPost, error) {
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{"updatedAt", 1}})

	cursor, err := s.BaseServiceMongoImpl.collection.Find(ctx, nil, opts)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}
	defer cursor.Close(ctx)

	var results []models.FbPost
	if err = cursor.All(ctx, &results); err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return results, nil
}

// UpdateToken cập nhật access token của một FbPost theo ID
func (s *FbPostService) UpdateToken(ctx context.Context, input *models.FbPostUpdateTokenInput) (*models.FbPost, error) {
	// Tìm FbPost theo post
	post, err := s.FindOneByPostID(ctx, input.PostId)
	if err != nil {
		return nil, err
	}

	// Cập nhật thông tin FbPost
	post.PanCakeData = input.PanCakeData
	post.UpdatedAt = time.Now().Unix()

	// Cập nhật FbPost
	updatedPost, err := s.BaseServiceMongoImpl.UpdateById(ctx, post.ID, post)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return &updatedPost, nil
}
