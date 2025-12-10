package dto

// FbConversationCreateInput dữ liệu đầu vào khi tạo conversation
type FbConversationCreateInput struct {
	PageId       string                 `json:"pageId" validate:"required"`
	PageUsername string                 `json:"pageUsername" validate:"required"`
	PanCakeData  map[string]interface{} `json:"panCakeData"`
}

