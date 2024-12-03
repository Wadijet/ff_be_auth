package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User , định nghĩa mô hình người dùng
// Token chứa token xác thực mới nhất của người dùng
// Tokens chứa danh sách các token, mỗi thiết bị khác nhau sẽ có một token riêng để xác thực (bằng hwid)
type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`      // ID của người dùng
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`   // Tên của người dùng
	Email    string             `json:"email,omitempty" bson:"email,omitempty"` // Email của người dùng
	Password string             `json:"-" bson:"password,omitempty"`            // Mật khẩu của người dùng
	Salt     string             `json:"-" bson:"salt,omitempty"`                // Muối để mã hóa mật khẩu
	Token    string             `json:"token" bson:"token,omitempty"`           // Token xác thực
	Tokens   []Token            `json:"tokens" bson:"tokens,omitempty"`         // Danh sách các token đang hiệụ lực (mỗi hwid sẽ có một token)
	IsBlock   bool   `json:"isBlock,omitempty" bson:"isBlock,omitempty"`     // Trạng thái bị khóa
	BlockNote string `json:"blockNote,omitempty" bson:"blockNote,omitempty"` // Ghi chú về việc bị khóa
	CreatedAt int64  `json:"createdAt,omitempty" bson:"createdAt,omitempty"` // Thời gian tạo
	UpdatedAt int64  `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"` // Thời gian cập nhật
}

// ComparePassword so sánh mật khẩu
func (u *User) ComparePassword(password string) error {
	existing := []byte(u.Password)
	incoming := []byte(password + u.Salt)
	err := bcrypt.CompareHashAndPassword(existing, incoming)
	return err
}

// =======================================================================================

// UserCreateInput , đầu vào tạo người dùng
type UserCreateInput struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty" validate:"required"`   // Tên của người dùng
	Email    string `json:"email,omitempty" bson:"email,omitempty" validate:"required"` // Email của người dùng
	Password string `json:"password" bson:"password" validate:"required"`               // Mật khẩu của người dùng
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
	Name string `json:"name,omitempty" bson:"name,omitempty"` // Tên của người dùng
}

// UserFilterInput , đầu vào lọc người dùng
type UserFilterInput struct {
	Emails []string `json:"emails" bson:"emails"` // Danh sách email
}
