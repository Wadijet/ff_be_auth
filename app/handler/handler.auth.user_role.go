package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/config"

	"go.mongodb.org/mongo-driver/mongo"
)

// UserRoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò của người dùng
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type UserRoleHandler struct {
	BaseHandler[models.UserRole, models.UserRoleCreateInput, models.UserRoleCreateInput]
	UserRoleService *services.UserRoleService
}

// NewUserRoleHandler khởi tạo một UserRoleHandler mới
func NewUserRoleHandler(c *config.Configuration, db *mongo.Client) *UserRoleHandler {
	newHandler := new(UserRoleHandler)
	newHandler.UserRoleService = services.NewUserRoleService(c, db)
	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng UserRoleService trực tiếp
	return newHandler
}

// Các hàm đặc thù của UserRole (nếu có) sẽ được thêm vào đây
