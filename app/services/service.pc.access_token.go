package services

import (
	"context"
	"time"

	"meta_commerce/app/database/registry"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/utility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AccessTokenService là cấu trúc chứa các phương thức liên quan đến access token
type AccessTokenService struct {
	*BaseServiceMongoImpl[models.AccessToken]
}

// NewAccessTokenService tạo mới AccessTokenService
func NewAccessTokenService() *AccessTokenService {
	accessTokenCollection := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.AccessTokens)
	return &AccessTokenService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.AccessToken](accessTokenCollection),
	}
}

// IsNameExist kiểm tra tên access token có tồn tại hay không
func (s *AccessTokenService) IsNameExist(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}
	var accessToken models.AccessToken
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&accessToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, utility.ConvertMongoError(err)
	}
	return true, nil
}

// Create tạo mới một access token
func (s *AccessTokenService) Create(ctx context.Context, input *models.AccessTokenCreateInput) (*models.AccessToken, error) {
	// Kiểm tra tên tồn tại
	exists, err := s.IsNameExist(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utility.ErrInvalidInput
	}

	// Chuyển đổi input.AssignedUsers từ mảng []string sang mảng []ObjectID
	assignedUsers := make([]primitive.ObjectID, 0)
	for _, userID := range input.AssignedUsers {
		assignedUsers = append(assignedUsers, utility.String2ObjectID(userID))
	}

	// Tạo access token mới
	accessToken := &models.AccessToken{
		ID:            primitive.NewObjectID(),
		Name:          input.Name,
		Describe:      input.Describe,
		System:        input.System,
		Value:         input.Value,
		AssignedUsers: assignedUsers,
		Status:        0,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	// Lưu access token
	createdAccessToken, err := s.BaseServiceMongoImpl.InsertOne(ctx, *accessToken)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return &createdAccessToken, nil
}

// Update cập nhật thông tin access token
func (s *AccessTokenService) Update(ctx context.Context, id primitive.ObjectID, input *models.AccessTokenUpdateInput) (*models.AccessToken, error) {
	// Kiểm tra access token tồn tại
	accessToken, err := s.BaseServiceMongoImpl.FindOneById(ctx, id)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	// Nếu có thay đổi tên, kiểm tra tên mới
	if input.Name != "" && input.Name != accessToken.Name {
		exists, err := s.IsNameExist(ctx, input.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, utility.ErrInvalidInput
		}
		accessToken.Name = input.Name
	}

	// Cập nhật thông tin khác
	if input.Describe != "" {
		accessToken.Describe = input.Describe
	}
	if input.System != "" {
		accessToken.System = input.System
	}
	if input.Value != "" {
		accessToken.Value = input.Value
	}
	if len(input.AssignedUsers) > 0 {
		assignedUsers := make([]primitive.ObjectID, 0)
		for _, userID := range input.AssignedUsers {
			assignedUsers = append(assignedUsers, utility.String2ObjectID(userID))
		}
		accessToken.AssignedUsers = assignedUsers
	}
	accessToken.UpdatedAt = time.Now().Unix()

	// Cập nhật access token
	updatedAccessToken, err := s.BaseServiceMongoImpl.UpdateById(ctx, id, accessToken)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return &updatedAccessToken, nil
}
