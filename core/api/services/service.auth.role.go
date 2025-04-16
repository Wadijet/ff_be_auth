package services

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/global"
	"meta_commerce/core/utility"
	"meta_commerce/pkg/registry"
)

// RoleService là cấu trúc chứa các phương thức liên quan đến vai trò
type RoleService struct {
	*BaseServiceMongoImpl[models.Role]
}

// NewRoleService tạo mới RoleService
func NewRoleService() (*RoleService, error) {
	roleCollection, exist := registry.Collections.Get(global.MongoDB_ColNames.Roles)
	if !exist {
		return nil, fmt.Errorf("failed to get roles collection: %v", utility.ErrNotFound)
	}

	return &RoleService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Role](roleCollection),
	}, nil
}
