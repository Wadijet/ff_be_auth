package services

import (
	"context"
	"time"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"
	"errors"

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

// IsExist kiểm tra vai trò người dùng có tồn tại hay không
func (s *UserRoleService) IsExist(ctx context.Context, userID, roleID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"userId": userID,
		"roleId": roleID,
	}
	var userRole models.UserRole
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&userRole)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Create tạo mới một vai trò người dùng
func (s *UserRoleService) Create(ctx context.Context, input *models.UserRoleCreateInput) (*models.UserRole, error) {
	// Kiểm tra User có tồn tại không
	if _, err := s.userService.FindOne(ctx, input.UserID.Hex()); err != nil {
		return nil, errors.New("User not found")
	}

	// Kiểm tra Role có tồn tại không
	if _, err := s.roleService.FindOne(ctx, input.RoleID.Hex()); err != nil {
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
	createdUserRole, err := s.BaseServiceImpl.Create(ctx, *userRole)
	if err != nil {
		return nil, err
	}

	return &createdUserRole, nil
}

// Delete xóa vai trò người dùng
func (s *UserRoleService) Delete(ctx context.Context, id string) error {
	// Kiểm tra UserRole có tồn tại không
	if _, err := s.BaseServiceImpl.FindOne(ctx, id); err != nil {
		return errors.New("UserRole not found")
	}

	return s.BaseServiceImpl.Delete(ctx, id)
}
