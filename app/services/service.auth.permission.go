package services

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/config"
	"meta_commerce/global"

	"go.mongodb.org/mongo-driver/mongo"
)

// PermissionService là cấu trúc chứa các phương thức liên quan đến quyền
type PermissionService struct {
	*BaseServiceMongoImpl[models.Permission]
}

// NewPermissionService tạo mới PermissionService
func NewPermissionService(c *config.Configuration, db *mongo.Client) *PermissionService {
	permissionCollection := GetCollectionFromName(db, GetDBNameFromCollectionName(c, global.MongoDB_ColNames.Permissions), global.MongoDB_ColNames.Permissions)
	return &PermissionService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Permission](permissionCollection),
	}
}
