package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission đại diện cho quyền trong hệ thống,
// Các quyền được kết cấu theo các quyền gọi các API trong router.
// Các quyèn này được tạo ra khi khởi tạo hệ thống và không thể thay đổi.
type FbConversation struct {
	ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`                        // ID của quyền
	ConversationId string                 `json:"conversationId" bson:"conversationId" index:"unique;text"` // ID của trang
	ApiData        map[string]interface{} `json:"apiData" bson:"apiData"`                                   // Dữ liệu API
	CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`                               // Thời gian tạo quyền
	UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`                               // Thời gian cập nhật quyền
}

// API INPUT STRUCT =======================================================================================

// FbPageCreateInput dữ liệu đầu vào khi tạo page
type FbConversationCreateInput struct {
	ApiData map[string]interface{} `json:"apiData" bson:"apiData" validate:"required"` // Dữ liệu API
}
