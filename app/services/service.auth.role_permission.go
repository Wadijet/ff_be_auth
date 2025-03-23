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

// RolePermissionService là cấu trúc chứa các phương thức liên quan đến quyền vai trò
type RolePermissionService struct {
	*BaseServiceImpl[models.RolePermission]
	roleService       *RoleService
	permissionService *PermissionService
}

// NewRolePermissionService tạo mới RolePermissionService
func NewRolePermissionService(c *config.Configuration, db *mongo.Client) *RolePermissionService {
	rolePermissionCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.RolePermissions)).Collection(global.MongoDB_ColNames.RolePermissions)
	return &RolePermissionService{
		BaseServiceImpl:   NewBaseService[models.RolePermission](rolePermissionCollection),
		roleService:       NewRoleService(c, db),
		permissionService: NewPermissionService(c, db),
	}
}

// IsExist kiểm tra quyền vai trò có tồn tại hay không
func (s *RolePermissionService) IsExist(ctx context.Context, roleID, permissionID primitive.ObjectID, scope byte) (bool, error) {
	filter := bson.M{
		"roleId":       roleID,
		"permissionId": permissionID,
		"scope":        scope,
	}
	var rolePermission models.RolePermission
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&rolePermission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Create tạo mới một quyền vai trò
func (s *RolePermissionService) Create(ctx context.Context, input *models.RolePermissionCreateInput) (*models.RolePermission, error) {
	// Kiểm tra Role có tồn tại không
	if _, err := s.roleService.FindOne(ctx, input.RoleID.Hex()); err != nil {
		return nil, errors.New("Role not found")
	}

	// Kiểm tra Permission có tồn tại không
	if _, err := s.permissionService.FindOne(ctx, input.PermissionID.Hex()); err != nil {
		return nil, errors.New("Permission not found")
	}

	// Kiểm tra RolePermission đã tồn tại chưa
	exists, err := s.IsExist(ctx, input.RoleID, input.PermissionID, input.Scope)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("RolePermission already exists")
	}

	// Tạo rolePermission mới
	rolePermission := &models.RolePermission{
		ID:           primitive.NewObjectID(),
		RoleID:       input.RoleID,
		PermissionID: input.PermissionID,
		Scope:        input.Scope,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	// Lưu rolePermission
	createdRolePermission, err := s.BaseServiceImpl.Create(ctx, *rolePermission)
	if err != nil {
		return nil, err
	}

	return &createdRolePermission, nil
}

// Delete xóa quyền vai trò
func (s *RolePermissionService) Delete(ctx context.Context, id string) error {
	return s.BaseServiceImpl.Delete(ctx, id)
}
