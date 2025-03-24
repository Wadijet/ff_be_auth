package services

import (
	"context"
	"time"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PermissionService là cấu trúc chứa các phương thức liên quan đến quyền
type PermissionService struct {
	*BaseServiceImpl[models.Permission]
}

// NewPermissionService tạo mới PermissionService
func NewPermissionService(c *config.Configuration, db *mongo.Client) *PermissionService {
	permissionCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.Permissions)).Collection(global.MongoDB_ColNames.Permissions)
	return &PermissionService{
		BaseServiceImpl: NewBaseService[models.Permission](permissionCollection),
	}
}

// IsNameExist kiểm tra tên quyền có tồn tại hay không
func (s *PermissionService) IsNameExist(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}
	var permission models.Permission
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&permission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Create tạo mới một quyền
func (s *PermissionService) Create(ctx context.Context, input *models.PermissionCreateInput) (*models.Permission, error) {
	// Kiểm tra tên tồn tại
	exists, err := s.IsNameExist(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicate
	}

	// Tạo permission mới
	permission := &models.Permission{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		Describe:  input.Describe,
		Category:  input.Category,
		Group:     input.Group,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Lưu permission
	createdPermission, err := s.BaseServiceImpl.Create(ctx, *permission)
	if err != nil {
		return nil, err
	}

	return &createdPermission, nil
}

// Update cập nhật thông tin quyền
func (s *PermissionService) Update(ctx context.Context, id primitive.ObjectID, input *models.PermissionUpdateInput) (*models.Permission, error) {
	// Kiểm tra permission tồn tại
	permission, err := s.BaseServiceImpl.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Nếu có thay đổi tên, kiểm tra tên mới
	if input.Name != "" && input.Name != permission.Name {
		exists, err := s.IsNameExist(ctx, input.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrDuplicate
		}
		permission.Name = input.Name
	}

	// Cập nhật thông tin khác
	if input.Describe != "" {
		permission.Describe = input.Describe
	}
	if input.Category != "" {
		permission.Category = input.Category
	}
	if input.Group != "" {
		permission.Group = input.Group
	}
	permission.UpdatedAt = time.Now().Unix()

	// Cập nhật permission
	updatedPermission, err := s.BaseServiceImpl.Update(ctx, id, permission)
	if err != nil {
		return nil, err
	}

	return &updatedPermission, nil
}

// Delete xóa quyền
func (s *PermissionService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.BaseServiceImpl.Delete(ctx, id)
}
