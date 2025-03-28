package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/config"

	"go.mongodb.org/mongo-driver/mongo"
)

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type RoleHandler struct {
	BaseHandler[models.Role, models.RoleCreateInput, models.RoleUpdateInput]
	RoleService *services.RoleService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewRoleHandler(c *config.Configuration, db *mongo.Client) *RoleHandler {
	newHandler := new(RoleHandler)
	newHandler.RoleService = services.NewRoleService(c, db)
	newHandler.BaseHandler.Service = newHandler.RoleService // Gán service cho BaseHandler
	return newHandler
}

// Các hàm đặc thù của Role (nếu có) sẽ được thêm vào đây
