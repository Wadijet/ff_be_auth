package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User , định nghĩa mô hình người dùng
// Token chứa token xác thực mới nhất của người dùng
// Tokens chứa danh sách các token, mỗi thiết bị khác nhau sẽ có một token riêng để xác thực (bằng hwid)
type User struct {
	_Relationships struct{}          `relationship:"collection:user_roles,field:userId,message:Không thể xóa user vì có %d role đang được gán cho user này. Vui lòng gỡ các role trước."` // Relationship definitions - không export, chỉ dùng cho tag parsing
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                                                                                                                      // ID của người dùng
	Name          string             `json:"name" bson:"name"`                                             // Tên của người dùng
	Email         string             `json:"email,omitempty" bson:"email,omitempty" index:"unique,sparse"` // Email của người dùng (sparse để cho phép null) - Optional vì dùng Firebase
	Password      string             `json:"-" bson:"password,omitempty"`                                  // DEPRECATED: Không còn sử dụng - Firebase quản lý authentication
	Salt          string             `json:"-" bson:"salt,omitempty"`                                      // DEPRECATED: Không còn sử dụng - Firebase quản lý authentication
	Phone         string             `json:"phone,omitempty" bson:"phone,omitempty" index:"unique,sparse"` // Số điện thoại (sparse để cho phép null) - Optional vì dùng Firebase
	FirebaseUID   string             `json:"firebaseUid" bson:"firebaseUid" index:"unique"`                // Firebase User ID (primary key để link với Firebase)
	EmailVerified bool               `json:"emailVerified" bson:"emailVerified"`                           // Email đã được xác thực
	PhoneVerified bool               `json:"phoneVerified" bson:"phoneVerified"`                           // Số điện thoại đã được xác thực
	AvatarURL     string             `json:"avatarUrl" bson:"avatarUrl"`                                   // URL avatar
	Token         string             `json:"token" bson:"token"`                                           // Token xác thực mới nhất của người dùng
	Tokens        []Token            `json:"-" bson:"tokens"`                                              // Danh sách các token đang hiệụ lực (mỗi hwid sẽ có một token)
	IsBlock       bool               `json:"-" bson:"isBlock"`                                             // Trạng thái bị khóa
	BlockNote     string             `json:"-" bson:"blockNote"`                                           // Ghi chú về việc bị khóa
	CreatedAt     int64              `json:"createdAt" bson:"createdAt"`                                   // Thời gian tạo
	UpdatedAt     int64              `json:"updatedAt" bson:"updatedAt"`                                   // Thời gian cập nhật
}

// PaginateResult đại diện cho kết quả phân trang
type UserPaginateResult struct {
	// Trang hiện tại
	Page int64 `json:"page" bson:"page"`
	// Số lượng mục trên mỗi trang
	Limit int64 `json:"limit" bson:"limit"`
	// Số lượng mục trong trang hiện tại
	ItemCount int64 `json:"itemCount" bson:"itemCount"`
	// Danh sách các mục
	Items []User `json:"items" bson:"items"`
}
