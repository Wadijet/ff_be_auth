package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission đại diện cho quyền trong hệ thống,
// Các quyền được kết cấu theo các quyền gọi các API trong router.
// Các quyèn này được tạo ra khi khởi tạo hệ thống và không thể thay đổi.
type FbPost struct {
	ID          primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`        // ID của quyền
	PageId      string                 `json:"pageId" bson:"pageId" index:"unique;text"` // ID của trang
	PostId      string                 `json:"postId" bson:"postId" index:"unique;text"` // ID của bài viết
	PanCakeData map[string]interface{} `json:"panCakeData" bson:"panCakeData"`           // Dữ liệu API
	CreatedAt   int64                  `json:"createdAt" bson:"createdAt"`               // Thời gian tạo quyền
	UpdatedAt   int64                  `json:"updatedAt" bson:"updatedAt"`               // Thời gian cập nhật quyền
}
