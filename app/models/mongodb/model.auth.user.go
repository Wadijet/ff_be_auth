package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User , định nghĩa mô hình người dùng
// Token chứa token xác thực mới nhất của người dùng
// Tokens chứa danh sách các token, mỗi thiết bị khác nhau sẽ có một token riêng để xác thực (bằng hwid)
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // ID của người dùng
	Name      string             `json:"name" bson:"name"`                  // Tên của người dùng
	Email     string             `json:"email" bson:"email" index:"unique"` // Email của người dùng
	Password  string             `json:"-" bson:"password"`                 // Mật khẩu của người dùng
	Salt      string             `json:"-" bson:"salt"`                     // Muối để mã hóa mật khẩu
	Tokens    []Token            `json:"tokens" bson:"tokens"`              // Danh sách các token đang hiệụ lực (mỗi hwid sẽ có một token)
	IsBlock   bool               `json:"isBlock" bson:"isBlock"`            // Trạng thái bị khóa
	BlockNote string             `json:"blockNote" bson:"blockNote"`        // Ghi chú về việc bị khóa
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`        // Thời gian tạo
	UpdatedAt int64              `json:"updatedAt" bson:"updatedAt"`        // Thời gian cập nhật
}

// ComparePassword so sánh mật khẩu
func (u *User) ComparePassword(password string) error {
	existing := []byte(u.Password)
	incoming := []byte(password + u.Salt)
	err := bcrypt.CompareHashAndPassword(existing, incoming)
	return err
}

// API INPUT STRUCT =======================================================================================

// UserCreateInput , đầu vào tạo người dùng
type UserCreateInput struct {
	Name     string `json:"name" bson:"name" validate:"required"`         // Tên của người dùng
	Email    string `json:"email" bson:"email" validate:"required"`       // Email của người dùng
	Password string `json:"password" bson:"password" validate:"required"` // Mật khẩu của người dùng
}

// UserLoginInput , đầu vào đăng nhập người dùng
type UserLoginInput struct {
	Email    string `json:"email" bson:"email" validate:"required"`       // Email của người dùng
	Password string `json:"password" bson:"password" validate:"required"` // Mật khẩu của người dùng
	Hwid     string `json:"hwid" bson:"hwid" validate:"required"`         // ID phần cứng
}

// UserLoginInput , đầu vào đăng nhập người dùng
type UserSetWorkingRoleInput struct {
	RoleID string `json:"roleId" bson:"roleId" validate:"required"` // ID của vai trò
}

// UserLogoutInput , đầu vào đăng xuất người dùng
type UserLogoutInput struct {
	Hwid string `json:"hwid" bson:"hwid" validate:"required"` // ID phần cứng
}

// UserChangePasswordInput , đầu vào thay đổi mật khẩu người dùng
type UserChangePasswordInput struct {
	OldPassword string `json:"oldPassword" bson:"oldPassword" validate:"required"` // Mật khẩu cũ
	NewPassword string `json:"newPassword" bson:"newPassword" validate:"required"` // Mật khẩu mới
}

// UserChangeInfoInput , đầu vào thay đổi thông tin người dùng
type UserChangeInfoInput struct {
	Name string `json:"name" bson:"name"` // Tên của người dùng
}
