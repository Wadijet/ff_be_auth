package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// MtService , định nghĩa mô hình dịch vụ
type MtService struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`              // ID của dịch vụ
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`           // Tên dịch vụ
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`         // Email của dịch vụ
	Password  string             `json:"-" bson:"password,omitempty"`                    // Mật khẩu của dịch vụ
	Salt      string             `json:"-" bson:"salt,omitempty"`                        // Muối để mã hóa mật khẩu
	Tokens    []Token            `json:"tokens" bson:"tokens,omitempty"`                 // Danh sách token của dịch vụ
	Role      primitive.ObjectID `json:"-" bson:"role,omitempty"`                        // Vai trò của dịch vụ
	IsBlock   bool               `json:"isBlock,omitempty" bson:"isBlock,omitempty"`     // Trạng thái bị khóa của dịch vụ
	BlockNote string             `json:"blockNote,omitempty" bson:"blockNote,omitempty"` // Ghi chú về việc bị khóa
	CreatedAt int64              `json:"createdAt,omitempty" bson:"createdAt,omitempty"` // Thời gian tạo dịch vụ
	UpdatedAt int64              `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"` // Thời gian cập nhật dịch vụ
}

// ComparePassword , so sánh mật khẩu
func (u *MtService) ComparePassword(password string) error {
	existing := []byte(u.Password)
	incoming := []byte(password + u.Salt)
	err := bcrypt.CompareHashAndPassword(existing, incoming)
	return err
}

// =======================================================================================

// MtServiceCreateInput , đầu vào để tạo dịch vụ
type MtServiceCreateInput struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty" validate:"required"`   // Tên dịch vụ
	Email    string `json:"email,omitempty" bson:"email,omitempty" validate:"required"` // Email của dịch vụ
	Password string `json:"password" bson:"password" validate:"required"`               // Mật khẩu của dịch vụ
}

// MtServiceLoginInput , đầu vào để đăng nhập dịch vụ
type MtServiceLoginInput struct {
	Email    string `json:"email" bson:"email" validate:"required"`       // Email của dịch vụ
	Password string `json:"password" bson:"password" validate:"required"` // Mật khẩu của dịch vụ
	Hwid     string `json:"hwid" bson:"hwid" validate:"required"`         // ID phần cứng
}

// MtServiceChangePasswordInput , đầu vào để thay đổi mật khẩu dịch vụ
type MtServiceChangePasswordInput struct {
	OldPassword string `json:"oldPassword" bson:"oldPassword" validate:"required"` // Mật khẩu cũ
	NewPassword string `json:"newPassword" bson:"newPassword" validate:"required"` // Mật khẩu mới
}

// MtServiceChangeInfoInput , đầu vào để thay đổi thông tin dịch vụ
type MtServiceChangeInfoInput struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"` // Tên dịch vụ
}

// MtServiceFilterInput , đầu vào để lọc dịch vụ
type MtServiceFilterInput struct {
	Emails  []string             `json:"emails" bson:"emails"`   // Danh sách email
	RoleIDs []primitive.ObjectID `json:"roleIDs" bson:"roleIDs"` // Danh sách ID vai trò
}

// MtServiceCheckTokenInput , đầu vào để kiểm tra token dịch vụ
type MtServiceCheckTokenInput struct {
	Token       string   `json:"token" bson:"token,omitempty"`             // Token của dịch vụ
	Permissions []string `json:"permissions" bson:"permissions,omitempty"` // Danh sách quyền
}
