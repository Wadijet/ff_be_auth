package services

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/config"
	"meta_commerce/global"

	"go.mongodb.org/mongo-driver/mongo"
)

// PermissionService là cấu trúc chứa các phương thức liên quan đến quyền
type PermissionService struct {
	*BaseServiceImpl[models.Permission]
}

// NewPermissionService tạo mới PermissionService
func NewPermissionService(c *config.Configuration, db *mongo.Client) *PermissionService {
	permissionCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.Permissions)).Collection(global.MongoDB_ColNames.Permissions)
	return &PermissionService{
		BaseServiceImpl: NewBaseService[models.Permission](permissionCollection),
	}
}
