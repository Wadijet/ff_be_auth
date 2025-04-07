package services

import (
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/registry"
)

// RoleService là cấu trúc chứa các phương thức liên quan đến vai trò
type RoleService struct {
	*BaseServiceMongoImpl[models.Role]
}

// NewRoleService tạo mới RoleService
func NewRoleService() *RoleService {
	roleCollection := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Roles)
	return &RoleService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Role](roleCollection),
	}
}
