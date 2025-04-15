package services

import (
	"fmt"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/registry"
)

// PermissionService là cấu trúc chứa các phương thức liên quan đến quyền
type PermissionService struct {
	*BaseServiceMongoImpl[models.Permission]
}

// NewPermissionService tạo mới PermissionService
func NewPermissionService() (*PermissionService, error) {
	permissionCollection, err := registry.Collections.MustGet(global.MongoDB_ColNames.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions collection: %v", err)
	}

	return &PermissionService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Permission](permissionCollection),
	}, nil
}
