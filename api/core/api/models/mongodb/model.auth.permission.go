package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission đại diện cho quyền trong hệ thống,
// Các quyền được kết cấu theo các quyền gọi các API trong router.
// Các quyèn này được tạo ra khi khởi tạo hệ thống và không thể thay đổi.
type Permission struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // ID của quyền
	Name      string             `json:"name" bson:"name" index:"unique"`   // Tên của quyền
	Describe  string             `json:"describe" bson:"describe"`          // Mô tả quyền
	Category  string             `json:"category" bson:"category"`          // Danh mục của quyền
	Group     string             `json:"group" bson:"group"`                // Nhóm của quyền
	IsSystem  bool               `json:"-" bson:"isSystem" index:"single:1"` // true = dữ liệu hệ thống, không thể xóa (chỉ dùng nội bộ, không expose ra API)
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`        // Thời gian tạo quyền
	UpdatedAt int64              `json:"updatedAt" bson:"updatedAt"`        // Thời gian cập nhật quyền
}
