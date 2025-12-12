package dto

// UserCreateInput , đầu vào tạo người dùng
// Lưu ý: User được tạo tự động từ Firebase, không cần tạo thủ công
// DTO này chỉ dùng cho CRUD operations, không dùng cho authentication
type UserCreateInput struct {
	Name  string `json:"name" validate:"required"`  // Tên của người dùng
	Email string `json:"email" validate:"required"` // Email của người dùng
	// Password đã bị deprecated - User được tạo từ Firebase, không cần password
}

// UserSetWorkingRoleInput , đầu vào đăng nhập người dùng
type UserSetWorkingRoleInput struct {
	RoleID string `json:"roleId" validate:"required"` // ID của vai trò
}

// UserLogoutInput , đầu vào đăng xuất người dùng
type UserLogoutInput struct {
	Hwid string `json:"hwid" validate:"required"` // ID phần cứng
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
