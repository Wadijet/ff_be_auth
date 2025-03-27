package services

import (
	"context"
	"errors"
	"time"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FbPageService là cấu trúc chứa các phương thức liên quan đến trang Facebook
type FbPageService struct {
	*BaseServiceImpl[models.FbPage]
}

// NewFbPageService tạo mới FbPageService
func NewFbPageService(c *config.Configuration, db *mongo.Client) *FbPageService {
	fbPageCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.FbPages)).Collection(global.MongoDB_ColNames.FbPages)
	return &FbPageService{
		BaseServiceImpl: NewBaseService[models.FbPage](fbPageCollection),
	}
}

// IsPageExist kiểm tra trang Facebook có tồn tại hay không
func (s *FbPageService) IsPageExist(ctx context.Context, pageId string) (bool, error) {
	filter := bson.M{"pageId": pageId}
	var page models.FbPage
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ReviceData nhận data từ Facebook và lưu vào cơ sở dữ liệu
func (s *FbPageService) ReviceData(ctx context.Context, input *models.FbPageCreateInput) (*models.FbPage, error) {
	if input.PanCakeData == nil {
		return nil, errors.New("ApiData is required")
	}

	// Lấy thông tin PageID từ ApiData đưa vào biến
	pageId := input.PanCakeData["id"].(string)

	// Kiểm tra FbPage đã tồn tại chưa
	exists, err := s.IsPageExist(ctx, pageId)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Tạo một FbPage mới
		page := &models.FbPage{
			ID:           primitive.NewObjectID(),
			AccessToken:  input.AccessToken,
			PanCakeData:  input.PanCakeData,
			PageName:     input.PanCakeData["name"].(string),
			PageUsername: input.PanCakeData["username"].(string),
			PageId:       input.PanCakeData["id"].(string),
			IsSync:       false,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}

		// Lưu FbPage
		createdPage, err := s.BaseServiceImpl.InsertOne(ctx, *page)
		if err != nil {
			return nil, err
		}

		return &createdPage, nil
	} else {
		// Lấy FbPage hiện tại
		filter := bson.M{"pageId": pageId}
		page, err := s.BaseServiceImpl.FindOne(ctx, filter, nil)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin mới
		page.PanCakeData = input.PanCakeData
		page.AccessToken = input.AccessToken
		page.PageName = input.PanCakeData["name"].(string)
		page.PageUsername = input.PanCakeData["username"].(string)
		page.UpdatedAt = time.Now().Unix()

		// Cập nhật FbPage
		updatedPage, err := s.BaseServiceImpl.UpdateById(ctx, page.ID, page)
		if err != nil {
			return nil, err
		}

		return &updatedPage, nil
	}
}

// FindOneByPageID tìm một FbPage theo PageID
func (s *FbPageService) FindOneByPageID(ctx context.Context, pageID string) (models.FbPage, error) {
	filter := bson.M{"pageId": pageID}
	var page models.FbPage
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return page, errors.New("page not found")
		}
		return page, err
	}
	return page, nil
}

// FindAll tìm tất cả các FbPage với phân trang
func (s *FbPageService) FindAll(ctx context.Context, page int64, limit int64) ([]models.FbPage, error) {
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{"updatedAt", 1}})

	cursor, err := s.BaseServiceImpl.collection.Find(ctx, nil, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.FbPage
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// UpdateToken cập nhật access token của một FbPage theo ID
func (s *FbPageService) UpdateToken(ctx context.Context, input *models.FbPageUpdateTokenInput) (*models.FbPage, error) {
	// Tìm FbPage theo page
	page, err := s.FindOneByPageID(ctx, input.PageId)
	if err != nil {
		return nil, err
	}

	// Cập nhật thông tin FbPage
	page.PageAccessToken = input.PageAccessToken
	page.UpdatedAt = time.Now().Unix()

	// Cập nhật FbPage
	updatedPage, err := s.BaseServiceImpl.UpdateById(ctx, page.ID, page)
	if err != nil {
		return nil, err
	}

	return &updatedPage, nil
}
