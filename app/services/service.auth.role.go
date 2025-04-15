package services

import (
	"fmt"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/registry"
)

// RoleService là cấu trúc chứa các phương thức liên quan đến vai trò
type RoleService struct {
	*BaseServiceMongoImpl[models.Role]
}

// NewRoleService tạo mới RoleService
func NewRoleService() (*RoleService, error) {
	roleCollection, err := registry.Collections.MustGet(global.MongoDB_ColNames.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles collection: %v", err)
	}

	return &RoleService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Role](roleCollection),
	}, nil
}
