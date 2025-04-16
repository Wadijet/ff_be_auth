package services

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/global"
	"meta_commerce/core/utility"
	"meta_commerce/pkg/registry"
)

// PermissionService là cấu trúc chứa các phương thức liên quan đến quyền
type PermissionService struct {
	*BaseServiceMongoImpl[models.Permission]
}

// NewPermissionService tạo mới PermissionService
func NewPermissionService() (*PermissionService, error) {
	permissionCollection, exist := registry.Collections.Get(global.MongoDB_ColNames.Permissions)
	if !exist {
		return nil, fmt.Errorf("failed to get permissions collection: %v", utility.ErrNotFound)
	}

	return &PermissionService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Permission](permissionCollection),
	}, nil
}
