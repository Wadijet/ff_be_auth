package services

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// PermissionService là cấu trúc chứa các phương thức liên quan đến quyền
type PermissionService struct {
	*BaseServiceMongoImpl[models.Permission]
}

// NewPermissionService tạo mới PermissionService
func NewPermissionService() (*PermissionService, error) {
	permissionCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.Permissions)
	if !exist {
		return nil, fmt.Errorf("failed to get permissions collection: %v", common.ErrNotFound)
	}

	return &PermissionService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Permission](permissionCollection),
	}, nil
}
