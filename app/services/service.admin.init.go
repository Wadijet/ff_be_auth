package services

import (
	"go.mongodb.org/mongo-driver/mongo"

	"atk-go-server/config"
	"atk-go-server/global"
)

// InitService định nghĩa các CRUD repository cho User, Permission và Role
type InitService struct {
	UserCRUD       Repository
	PermissionCRUD Repository
	RoleCRUD       Repository
}

// NewInitService khởi tạo các repository và trả về một đối tượng InitService
func NewInitService(c *config.Configuration, db *mongo.Client) *InitService {
	newService := new(InitService)
	newService.UserCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Users)
	newService.PermissionCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newService.RoleCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Roles)
	return newService
}
