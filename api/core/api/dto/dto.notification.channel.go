package dto

// NotificationChannelCreateInput dùng cho tạo notification channel (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi tạo channel
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type NotificationChannelCreateInput struct {
	ChannelType string   `json:"channelType" validate:"required"` // Loại kênh: email, telegram, webhook - BẮT BUỘC
	Name        string   `json:"name" validate:"required"`        // Tên channel - BẮT BUỘC
	IsActive    bool     `json:"isActive"`                        // Channel có đang hoạt động không - Optional (default: true)
	SenderIDs   []string `json:"senderIds,omitempty"`             // Mảng sender IDs (thứ tự ưu tiên) - Optional, null/empty = dùng inheritance

	// Email recipients
	Recipients []string `json:"recipients,omitempty"` // Email addresses - Optional (chỉ dùng khi channelType = "email")

	// Telegram recipients
	ChatIDs []string `json:"chatIds,omitempty"` // Telegram chat IDs - Optional (chỉ dùng khi channelType = "telegram")

	// Webhook recipients
	WebhookURL     string            `json:"webhookUrl,omitempty"`     // Webhook URL (chỉ 1 URL) - Optional (chỉ dùng khi channelType = "webhook")
	WebhookHeaders map[string]string `json:"webhookHeaders,omitempty"` // Webhook headers - Optional (chỉ dùng khi channelType = "webhook")

	// Lưu ý: KHÔNG cần gửi ownerOrganizationId - Backend tự động gán từ context (phân quyền dữ liệu)
}

// NotificationChannelUpdateInput dùng cho cập nhật notification channel (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi cập nhật channel
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type NotificationChannelUpdateInput struct {
	ChannelType string   `json:"channelType"` // Loại kênh: email, telegram, webhook - Optional
	Name        string   `json:"name"`        // Tên channel - Optional
	IsActive    *bool    `json:"isActive"`   // Channel có đang hoạt động không - Optional
	SenderIDs   []string `json:"senderIds,omitempty"` // Mảng sender IDs (thứ tự ưu tiên) - Optional

	// Email recipients
	Recipients []string `json:"recipients,omitempty"` // Email addresses - Optional

	// Telegram recipients
	ChatIDs []string `json:"chatIds,omitempty"` // Telegram chat IDs - Optional

	// Webhook recipients
	WebhookURL     string            `json:"webhookUrl,omitempty"`     // Webhook URL - Optional
	WebhookHeaders map[string]string `json:"webhookHeaders,omitempty"` // Webhook headers - Optional

	// Lưu ý: KHÔNG thể update ownerOrganizationId - Backend sẽ tự động xóa field này nếu có trong request (bảo mật)
}
