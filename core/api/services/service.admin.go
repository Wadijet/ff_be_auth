package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/utility"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminService là cấu trúc chứa các phương thức liên quan đến admin
type AdminService struct {
	userService           *UserService
	roleService           *RoleService
	permissionService     *PermissionService
	userRoleService       *UserRoleService
	rolePermissionService *RolePermissionService
}

// NewAdminService tạo mới AdminService
func NewAdminService() (*AdminService, error) {
	userService, err := NewUserService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %v", err)
	}

	roleService, err := NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}

	permissionService, err := NewPermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create permission service: %v", err)
	}

	userRoleService, err := NewUserRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user_role service: %v", err)
	}

	rolePermissionService, err := NewRolePermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role_permission service: %v", err)
	}

	return &AdminService{
		userService:           userService,
		roleService:           roleService,
		permissionService:     permissionService,
		userRoleService:       userRoleService,
		rolePermissionService: rolePermissionService,
	}, nil
}

// SetRole gán Role cho User dựa trên Email và RoleID
func (s *AdminService) SetRole(ctx context.Context, email string, roleID primitive.ObjectID) (*models.User, error) {
	// Kiểm tra Role có tồn tại không
	_, err := s.roleService.FindOneById(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Tìm User theo Email
	filter := bson.M{"email": email}
	var user models.User
	err = s.userService.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utility.ErrNotFound
		}
		return nil, utility.ConvertMongoError(err)
	}

	// Cập nhật Role cho User
	user.Token = roleID.Hex() // Sử dụng Token để lưu RoleID
	user.UpdatedAt = time.Now().Unix()

	// Cập nhật User
	updatedUser, err := s.userService.UpdateById(ctx, user.ID, user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// BlockUser chặn hoặc bỏ chặn User dựa trên Email và trạng thái Block
func (s *AdminService) BlockUser(ctx context.Context, email string, block bool, note string) (*models.User, error) {
	// Tìm User theo Email
	filter := bson.M{"email": email}
	var user models.User
	err := s.userService.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utility.ErrNotFound
		}
		return nil, utility.ConvertMongoError(err)
	}

	// Cập nhật trạng thái Block và ghi chú
	user.IsBlock = block
	user.BlockNote = note
	user.UpdatedAt = time.Now().Unix()

	// Cập nhật User
	updatedUser, err := s.userService.UpdateById(ctx, user.ID, user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// UnBlockUser mở khóa người dùng
func (s *AdminService) UnBlockUser(ctx context.Context, email string) (*models.User, error) {
	// Tìm User theo Email
	filter := bson.M{"email": email}
	var user models.User
	err := s.userService.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utility.ErrNotFound
		}
		return nil, utility.ConvertMongoError(err)
	}

	// Cập nhật trạng thái Block và ghi chú
	user.IsBlock = false
	user.BlockNote = ""
	user.UpdatedAt = time.Now().Unix()

	// Cập nhật User
	updatedUser, err := s.userService.UpdateById(ctx, user.ID, user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
