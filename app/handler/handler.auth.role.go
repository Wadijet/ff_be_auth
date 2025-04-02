package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
)

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type RoleHandler struct {
	BaseHandler[models.Role, models.RoleCreateInput, models.RoleUpdateInput]
	RoleService *services.RoleService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewRoleHandler() *RoleHandler {
	newHandler := new(RoleHandler)
	newHandler.RoleService = services.NewRoleService()
	newHandler.BaseHandler.Service = newHandler.RoleService // Gán service cho BaseHandler
	return newHandler
}

// Các hàm đặc thù của Role (nếu có) sẽ được thêm vào đây
