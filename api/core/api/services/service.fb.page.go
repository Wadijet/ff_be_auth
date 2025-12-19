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

// FbPageService là cấu trúc chứa các phương thức liên quan đến Facebook page
type FbPageService struct {
	*BaseServiceMongoImpl[models.FbPage]
}

// NewFbPageService tạo mới FbPageService
func NewFbPageService() (*FbPageService, error) {
	fbPageCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.FbPages)
	if !exist {
		return nil, fmt.Errorf("failed to get fb_pages collection: %v", common.ErrNotFound)
	}

	return &FbPageService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.FbPage](fbPageCollection),
	}, nil
}

// IsPageExist kiểm tra trang Facebook có tồn tại hay không
func (s *FbPageService) IsPageExist(ctx context.Context, pageId string) (bool, error) {
	filter := bson.M{"pageId": pageId}
	var page models.FbPage
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, common.ConvertMongoError(err)
	}
	return true, nil
}

// FindOneByPageID tìm một FbPage theo PageID
func (s *FbPageService) FindOneByPageID(ctx context.Context, pageID string) (models.FbPage, error) {
	filter := bson.M{"pageId": pageID}
	var page models.FbPage
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return page, common.ErrNotFound
		}
		return page, common.ConvertMongoError(err)
	}
	return page, nil
}

// FindAll tìm tất cả các FbPage với phân trang
func (s *FbPageService) FindAll(ctx context.Context, page int64, limit int64) ([]models.FbPage, error) {
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{"updatedAt", 1}})

	cursor, err := s.BaseServiceMongoImpl.collection.Find(ctx, nil, opts)
	if err != nil {
		return nil, common.ConvertMongoError(err)
	}
	defer cursor.Close(ctx)

	var results []models.FbPage
	if err = cursor.All(ctx, &results); err != nil {
		return nil, common.ConvertMongoError(err)
	}

	return results, nil
}

// UpdateToken cập nhật access token của một FbPage theo ID
func (s *FbPageService) UpdateToken(ctx context.Context, input *dto.FbPageUpdateTokenInput) (*models.FbPage, error) {
	// Tìm FbPage theo page
	page, err := s.FindOneByPageID(ctx, input.PageId)
	if err != nil {
		return nil, err
	}

	// Cập nhật thông tin FbPage
	page.PageAccessToken = input.PageAccessToken
	page.UpdatedAt = time.Now().Unix()

	// Cập nhật FbPage
	updatedPage, err := s.BaseServiceMongoImpl.UpdateById(ctx, page.ID, page)
	if err != nil {
		return nil, common.ConvertMongoError(err)
	}

	return &updatedPage, nil
}
