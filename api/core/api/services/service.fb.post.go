package services

import (
	"context"
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

// FbPostService là cấu trúc chứa các phương thức liên quan đến bài viết Facebook
type FbPostService struct {
	*BaseServiceMongoImpl[models.FbPost]
	fbPageService *FbPageService
}

// NewFbPostService tạo mới FbPostService
func NewFbPostService() (*FbPostService, error) {
	fbPostCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.FbPosts)
	if !exist {
		return nil, fmt.Errorf("failed to get fb_posts collection: %v", common.ErrNotFound)
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
		return false, common.ConvertMongoError(err)
	}
	return true, nil
}

// FindOneByPostID tìm một FbPost theo PostID
func (s *FbPostService) FindOneByPostID(ctx context.Context, postID string) (models.FbPost, error) {
	filter := bson.M{"postId": postID}
	var post models.FbPost
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return post, common.ErrNotFound
		}
		return post, common.ConvertMongoError(err)
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
		return nil, common.ConvertMongoError(err)
	}
	defer cursor.Close(ctx)

	var results []models.FbPost
	if err = cursor.All(ctx, &results); err != nil {
		return nil, common.ConvertMongoError(err)
	}

	return results, nil
}

// UpdateToken cập nhật access token của một FbPost theo ID
func (s *FbPostService) UpdateToken(ctx context.Context, input *dto.FbPostUpdateTokenInput) (*models.FbPost, error) {
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
		return nil, common.ConvertMongoError(err)
	}

	return &updatedPost, nil
}
