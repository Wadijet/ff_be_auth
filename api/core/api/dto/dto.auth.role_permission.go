package dto

// RolePermissionCreateInput đại diện cho dữ liệu đầu vào khi tạo quyền vai trò
type RolePermissionCreateInput struct {
	RoleID       string `json:"roleId" validate:"required"`       // ID của vai trò
	PermissionID string `json:"permissionId" validate:"required"` // ID của quyền
	Scope        byte   `json:"scope"`                            // Phạm vi của quyền (0: All, 1: Assign)
}

// RolePermissionUpdateInput dữ liệu đầu vào khi cập nhật quyền của vai trò
type RolePermissionUpdateInput struct {
	RoleId        string   `json:"roleId" validate:"required"`               // ID của vai trò
	PermissionIds []string `json:"permissionIds" validate:"required"`       // Danh sách ID của quyền
}

