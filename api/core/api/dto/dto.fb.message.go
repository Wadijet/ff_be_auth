package dto

// FbMessageCreateInput dữ liệu đầu vào khi tạo message
type FbMessageCreateInput struct {
	PageId         string                 `json:"pageId" validate:"required"`
	PageUsername   string                 `json:"pageUsername" validate:"required"`
	ConversationId string                 `json:"conversationId" validate:"required"`
	CustomerId     string                 `json:"customerId" validate:"required"`
	PanCakeData    map[string]interface{} `json:"panCakeData" validate:"required"`
}

