package services

import (
	"context"
	"fmt"
	"time"

	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/registry"
	"meta_commerce/app/utility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RolePermissionService là cấu trúc chứa các phương thức liên quan đến quyền của vai trò
type RolePermissionService struct {
	*BaseServiceMongoImpl[models.RolePermission]
	roleService       *RoleService
	permissionService *PermissionService
}

// NewRolePermissionService tạo mới RolePermissionService
func NewRolePermissionService() (*RolePermissionService, error) {
	rolePermissionCollection, err := registry.Collections.MustGet(global.MongoDB_ColNames.RolePermissions)
	if err != nil {
		return nil, fmt.Errorf("failed to get role_permissions collection: %v", err)
	}

	roleService, err := NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}

	permissionService, err := NewPermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create permission service: %v", err)
	}

	return &RolePermissionService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.RolePermission](rolePermissionCollection),
		roleService:          roleService,
		permissionService:    permissionService,
	}, nil
}

// Create tạo mới một quyền cho vai trò
func (s *RolePermissionService) Create(ctx context.Context, input *models.RolePermissionCreateInput) (*models.RolePermission, error) {
	// Kiểm tra Role có tồn tại không
	if _, err := s.roleService.FindOneById(ctx, input.RoleID); err != nil {
		return nil, utility.ErrNotFound
	}

	// Kiểm tra Permission có tồn tại không
	if _, err := s.permissionService.FindOneById(ctx, input.PermissionID); err != nil {
		return nil, utility.ErrNotFound
	}

	// Kiểm tra RolePermission đã tồn tại chưa
	exists, err := s.IsExist(ctx, input.RoleID, input.PermissionID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utility.ErrInvalidInput
	}

	// Tạo rolePermission mới
	rolePermission := &models.RolePermission{
		ID:           primitive.NewObjectID(),
		RoleID:       input.RoleID,
		PermissionID: input.PermissionID,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	// Lưu rolePermission
	createdRolePermission, err := s.BaseServiceMongoImpl.InsertOne(ctx, *rolePermission)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return &createdRolePermission, nil
}

// IsExist kiểm tra xem một RolePermission đã tồn tại chưa
func (s *RolePermissionService) IsExist(ctx context.Context, roleID, permissionID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"roleId":       roleID,
		"permissionId": permissionID,
	}
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, utility.ConvertMongoError(err)
	}
	return count > 0, nil
}
