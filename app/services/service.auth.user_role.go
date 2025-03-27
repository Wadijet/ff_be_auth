package services

import (
	"context"
	"time"

	"errors"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/config"
	"meta_commerce/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRoleService là cấu trúc chứa các phương thức liên quan đến vai trò người dùng
type UserRoleService struct {
	*BaseServiceImpl[models.UserRole]
	userService *UserService
	roleService *RoleService
}

// NewUserRoleService tạo mới UserRoleService
func NewUserRoleService(c *config.Configuration, db *mongo.Client) *UserRoleService {
	userRoleCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.UserRoles)).Collection(global.MongoDB_ColNames.UserRoles)
	return &UserRoleService{
		BaseServiceImpl: NewBaseService[models.UserRole](userRoleCollection),
		userService:     NewUserService(c, db),
		roleService:     NewRoleService(c, db),
	}
}

// Create tạo mới một vai trò người dùng
func (s *UserRoleService) Create(ctx context.Context, input *models.UserRoleCreateInput) (*models.UserRole, error) {
	// Kiểm tra User có tồn tại không
	if _, err := s.userService.FindOneById(ctx, input.UserID); err != nil {
		return nil, errors.New("User not found")
	}

	// Kiểm tra Role có tồn tại không
	if _, err := s.roleService.FindOneById(ctx, input.RoleID); err != nil {
		return nil, errors.New("Role not found")
	}

	// Kiểm tra UserRole đã tồn tại chưa
	exists, err := s.IsExist(ctx, input.UserID, input.RoleID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("UserRole already exists")
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
	createdUserRole, err := s.BaseServiceImpl.InsertOne(ctx, *userRole)
	if err != nil {
		return nil, err
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
		return false, err
	}
	return count > 0, nil
}
