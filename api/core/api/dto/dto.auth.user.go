package dto

// UserCreateInput , đầu vào tạo người dùng
type UserCreateInput struct {
	Name     string `json:"name" validate:"required"`     // Tên của người dùng
	Email    string `json:"email" validate:"required"`    // Email của người dùng
	Password string `json:"password" validate:"required"` // Mật khẩu của người dùng
}

// UserLoginInput , đầu vào đăng nhập người dùng
type UserLoginInput struct {
	Email    string `json:"email" validate:"required,email,min=6,max=100"` // Email của người dùng
	Password string `json:"password" validate:"required,min=8,max=32"`     // Mật khẩu của người dùng
	Hwid     string `json:"hwid" validate:"required"`                      // ID phần cứng
}

// UserSetWorkingRoleInput , đầu vào đăng nhập người dùng
type UserSetWorkingRoleInput struct {
	RoleID string `json:"roleId" validate:"required"` // ID của vai trò
}

// UserLogoutInput , đầu vào đăng xuất người dùng
type UserLogoutInput struct {
	Hwid string `json:"hwid" validate:"required"` // ID phần cứng
}

// UserChangePasswordInput , đầu vào thay đổi mật khẩu người dùng
type UserChangePasswordInput struct {
	OldPassword string `json:"oldPassword" validate:"required"` // Mật khẩu cũ
	NewPassword string `json:"newPassword" validate:"required"` // Mật khẩu mới
}

// UserChangeInfoInput , đầu vào thay đổi thông tin người dùng
type UserChangeInfoInput struct {
	Name string `json:"name"` // Tên của người dùng
}

// BlockUserInput là cấu trúc dữ liệu đầu vào cho việc khóa người dùng
type BlockUserInput struct {
	Email string `json:"email" validate:"required"`
	Note  string `json:"note" validate:"required"`
}

// UnBlockUserInput là cấu trúc dữ liệu đầu vào cho việc mở khóa người dùng
type UnBlockUserInput struct {
	Email string `json:"email" validate:"required"`
}

// FirebaseLoginInput đầu vào đăng nhập bằng Firebase ID token
type FirebaseLoginInput struct {
	IDToken string `json:"idToken" validate:"required"` // Firebase ID token
	Hwid    string `json:"hwid" validate:"required"`    // Device hardware ID
}

