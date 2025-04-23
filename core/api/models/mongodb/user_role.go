package models

// UserRoleUpdateInput là struct chứa dữ liệu đầu vào cho việc cập nhật user role
type UserRoleUpdateInput struct {
	UserID  string   `json:"userId" validate:"required"`
	RoleIDs []string `json:"roleIds" validate:"required,min=1"`
}
