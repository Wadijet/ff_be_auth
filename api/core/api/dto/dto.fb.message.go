package dto

// ============================================
// DTO CHO CRUD ROUTES (Logic chung)
// ============================================
// FbMessageCreateInput dữ liệu đầu vào cho CRUD operations (insert-one, update-one)
// - Dùng cho các endpoint CRUD chuẩn
// - KHÔNG có logic tách messages[] ra khỏi panCakeData
// - PanCakeData có thể chứa messages[] (tương thích ngược)
type FbMessageCreateInput struct {
	PageId         string                 `json:"pageId" validate:"required"`
	PageUsername   string                 `json:"pageUsername" validate:"required"`
	ConversationId string                 `json:"conversationId" validate:"required"`
	CustomerId     string                 `json:"customerId" validate:"required"`
	PanCakeData    map[string]interface{} `json:"panCakeData" validate:"required"` // PanCakeData (có thể có messages[])
}

// ============================================
// DTO CHO ENDPOINT ĐẶC BIỆT (Logic tách messages)
// ============================================
// FbMessageUpsertMessagesInput dữ liệu đầu vào cho endpoint đặc biệt upsert-messages
// - Dùng cho endpoint: POST /api/v1/facebook/message/upsert-messages
// - CÓ logic tự động tách messages[] ra khỏi panCakeData và lưu vào 2 collections
// - API bên ngoài vẫn gửi panCakeData đầy đủ (bao gồm messages[]), server sẽ tự động tách
type FbMessageUpsertMessagesInput struct {
	PageId         string                 `json:"pageId" validate:"required"`
	PageUsername   string                 `json:"pageUsername" validate:"required"`
	ConversationId string                 `json:"conversationId" validate:"required"`
	CustomerId     string                 `json:"customerId" validate:"required"`
	PanCakeData    map[string]interface{} `json:"panCakeData" validate:"required"` // PanCakeData đầy đủ (bao gồm messages[])
	HasMore        bool                   `json:"hasMore"`                         // Còn messages để sync không
}
