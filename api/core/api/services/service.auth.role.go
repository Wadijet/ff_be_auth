package services

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// RoleService là cấu trúc chứa các phương thức liên quan đến vai trò
type RoleService struct {
	*BaseServiceMongoImpl[models.Role]
}

// NewRoleService tạo mới RoleService
func NewRoleService() (*RoleService, error) {
	roleCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.Roles)
	if !exist {
		return nil, fmt.Errorf("failed to get roles collection: %v", common.ErrNotFound)
	}

	return &RoleService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Role](roleCollection),
	}, nil
}
