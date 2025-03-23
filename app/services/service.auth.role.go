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

// RoleService là cấu trúc chứa các phương thức liên quan đến vai trò
type RoleService struct {
	*BaseServiceImpl[models.Role]
}

// NewRoleService tạo mới RoleService
func NewRoleService(c *config.Configuration, db *mongo.Client) *RoleService {
	roleCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.Roles)).Collection(global.MongoDB_ColNames.Roles)
	return &RoleService{
		BaseServiceImpl: NewBaseService[models.Role](roleCollection),
	}
}

// IsNameExist kiểm tra tên role có tồn tại hay không
func (s *RoleService) IsNameExist(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}
	var role models.Role
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Create tạo mới một vai trò
func (s *RoleService) Create(ctx context.Context, input *models.RoleCreateInput) (*models.Role, error) {
	// Kiểm tra tên tồn tại
	exists, err := s.IsNameExist(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicate
	}

	// Tạo role mới
	role := &models.Role{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		Describe:  input.Describe,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Lưu role
	createdRole, err := s.BaseServiceImpl.Create(ctx, *role)
	if err != nil {
		return nil, err
	}

	return &createdRole, nil
}

// Update cập nhật thông tin vai trò
func (s *RoleService) Update(ctx context.Context, id string, input *models.RoleUpdateInput) (*models.Role, error) {
	// Kiểm tra role tồn tại
	role, err := s.BaseServiceImpl.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Nếu có thay đổi tên, kiểm tra tên mới
	if input.Name != "" && input.Name != role.Name {
		exists, err := s.IsNameExist(ctx, input.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrDuplicate
		}
		role.Name = input.Name
	}

	// Cập nhật thông tin khác
	if input.Describe != "" {
		role.Describe = input.Describe
	}
	role.UpdatedAt = time.Now().Unix()

	// Cập nhật role
	updatedRole, err := s.BaseServiceImpl.Update(ctx, id, role)
	if err != nil {
		return nil, err
	}

	return &updatedRole, nil
}

// Delete xóa vai trò
func (s *RoleService) Delete(ctx context.Context, id string) error {
	return s.BaseServiceImpl.Delete(ctx, id)
}
