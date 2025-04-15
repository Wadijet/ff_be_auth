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
)

// UserRoleService là cấu trúc chứa các phương thức liên quan đến vai trò của người dùng
type UserRoleService struct {
	*BaseServiceMongoImpl[models.UserRole]
	userService *UserService
	roleService *RoleService
}

// NewUserRoleService tạo mới UserRoleService
func NewUserRoleService() (*UserRoleService, error) {
	userRoleCollection, err := registry.Collections.MustGet(global.MongoDB_ColNames.UserRoles)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_roles collection: %v", err)
	}

	userService, err := NewUserService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %v", err)
	}

	roleService, err := NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}

	return &UserRoleService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.UserRole](userRoleCollection),
		userService:          userService,
		roleService:          roleService,
	}, nil
}

// Create tạo mới một vai trò người dùng
func (s *UserRoleService) Create(ctx context.Context, input *models.UserRoleCreateInput) (*models.UserRole, error) {
	// Kiểm tra User có tồn tại không
	if _, err := s.userService.FindOneById(ctx, input.UserID); err != nil {
		return nil, utility.ErrNotFound
	}

	// Kiểm tra Role có tồn tại không
	if _, err := s.roleService.FindOneById(ctx, input.RoleID); err != nil {
		return nil, utility.ErrNotFound
	}

	// Kiểm tra UserRole đã tồn tại chưa
	exists, err := s.IsExist(ctx, input.UserID, input.RoleID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utility.ErrInvalidInput
	}

	// Tạo userRole mới
	userRole := &models.UserRole{
		ID:        primitive.NewObjectID(),
		UserID:    input.UserID,
		RoleID:    input.RoleID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Lưu userRole
	createdUserRole, err := s.BaseServiceMongoImpl.InsertOne(ctx, *userRole)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return &createdUserRole, nil
}

// IsExist kiểm tra xem một UserRole đã tồn tại chưa
func (s *UserRoleService) IsExist(ctx context.Context, userID, roleID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"userId": userID,
		"roleId": roleID,
	}
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, utility.ConvertMongoError(err)
	}
	return count > 0, nil
}
