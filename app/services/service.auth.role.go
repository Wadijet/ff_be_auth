package services

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/config"
	"meta_commerce/global"

	"go.mongodb.org/mongo-driver/mongo"
)

// RoleService là cấu trúc chứa các phương thức liên quan đến vai trò
type RoleService struct {
	*BaseServiceImpl[models.Role]
}

// NewRoleService tạo mới RoleService
func NewRoleService(c *config.Configuration, db *mongo.Client) *RoleService {
	roleCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.Roles)).Collection(global.MongoDB_ColNames.Roles)
	return &RoleService{
		BaseServiceImpl: NewBaseService[models.Role](roleCollection),
	}
}
