package dto

// FbMessageItemCreateInput dữ liệu đầu vào khi tạo message item
type FbMessageItemCreateInput struct {
	ConversationId string                 `json:"conversationId" validate:"required"` // ID của cuộc hội thoại
	MessageId      string                 `json:"messageId" validate:"required"`      // ID của message từ Pancake (unique)
	MessageData    map[string]interface{} `json:"messageData" validate:"required"`    // Toàn bộ dữ liệu của message
	InsertedAt     int64                  `json:"insertedAt"`                         // Thời gian insert message (optional, có thể extract từ MessageData)
}

// FbMessageItemUpdateInput dữ liệu đầu vào khi cập nhật message item
type FbMessageItemUpdateInput struct {
	ConversationId string                 `json:"conversationId" validate:"required"` // ID của cuộc hội thoại
	MessageId      string                 `json:"messageId" validate:"required"`      // ID của message từ Pancake (unique)
	MessageData    map[string]interface{} `json:"messageData" validate:"required"`    // Toàn bộ dữ liệu của message
	InsertedAt     int64                  `json:"insertedAt"`                         // Thời gian insert message (optional, có thể extract từ MessageData)
}
