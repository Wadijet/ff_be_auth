package services

import (
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/registry"
)

// PermissionService là cấu trúc chứa các phương thức liên quan đến quyền
type PermissionService struct {
	*BaseServiceMongoImpl[models.Permission]
}

// NewPermissionService tạo mới PermissionService
func NewPermissionService() *PermissionService {
	permissionCollection := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Permissions)
	return &PermissionService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Permission](permissionCollection),
	}
}
