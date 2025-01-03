package services

import (
	"atk-go-server/config"
	"atk-go-server/global"

	"go.mongodb.org/mongo-driver/mongo"
)

// AccessTokenService là cấu trúc chứa các phương thức liên quan đến người dùng
type AccessTokenService struct {
	crudUser     RepositoryService
	crudUserRole RepositoryService
}

// Khởi tạo UserService với cấu hình và kết nối cơ sở dữ liệu
func NewAccessTokenService(c *config.Configuration, db *mongo.Client) *AccessTokenService {
	newService := new(AccessTokenService)
	newService.crudUser = *NewRepository(c, db, global.MongoDB_ColNames.Users)
	newService.crudUserRole = *NewRepository(c, db, global.MongoDB_ColNames.UserRoles)
	return newService
}
