package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission đại diện cho quyền trong hệ thống,
// Các quyền được kết cấu theo các quyền gọi các API trong router.
// Các quyèn này được tạo ra khi khởi tạo hệ thống và không thể thay đổi.
type FbConversation struct {
	ID               primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`                        // ID của quyền
	PageId           string                 `json:"pageId" bson:"pageId" index:"text"`                        // ID của trang
	PageUsername     string                 `json:"pageUsername" bson:"pageUsername" index:"text"`            // Tên người dùng của trang
	ConversationId   string                 `json:"conversationId" bson:"conversationId" index:"unique;text"` // ID của trang
	CustomerId       string                 `json:"customerId" bson:"customerId" index:"text"`                // ID của khách hàng
	PanCakeData      map[string]interface{} `json:"panCakeData" bson:"panCakeData"`                           // Dữ liệu API
	PanCakeUpdatedAt int64                  `json:"panCakeUpdatedAt" bson:"panCakeUpdatedAt"`                 // Thời gian cập nhật dữ liệu API
	CreatedAt        int64                  `json:"createdAt" bson:"createdAt"`                               // Thời gian tạo quyền
	UpdatedAt        int64                  `json:"updatedAt" bson:"updatedAt"`                               // Thời gian cập nhật quyền
}
