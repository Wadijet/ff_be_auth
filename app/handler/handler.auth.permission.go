package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
)

// PermissionHandler là struct chứa các phương thức xử lý quyền
// Kế thừa từ BaseHandler với các type parameter:
// - Model: models.Permission
// - CreateInput: models.PermissionCreateInput
// - UpdateInput: models.PermissionUpdateInput
type PermissionHandler struct {
	BaseHandler[models.Permission, models.PermissionCreateInput, models.PermissionUpdateInput]
}

// NewPermissionHandler khởi tạo một PermissionHandler mới
func NewPermissionHandler() *PermissionHandler {
	handler := &PermissionHandler{}
	handler.Service = services.NewPermissionService()
	return handler
}

// Các hàm đặc thù của Permission (nếu có) sẽ được thêm vào đây
