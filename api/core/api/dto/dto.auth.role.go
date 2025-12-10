package dto

// RoleCreateInput dùng cho tạo vai trò (tầng transport)
type RoleCreateInput struct {
	Name     string `json:"name" validate:"required"`     // Tên của vai trò
	Describe string `json:"describe" validate:"required"` // Mô tả vai trò
}

// RoleUpdateInput dùng cho cập nhật vai trò (tầng transport)
type RoleUpdateInput struct {
	Name     string `json:"name"`     // Tên của vai trò
	Describe string `json:"describe"` // Mô tả vai trò
}

