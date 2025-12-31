package dto

// NotificationChannelSenderCreateInput dùng cho tạo notification sender (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi tạo sender
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type NotificationChannelSenderCreateInput struct {
	ChannelType string `json:"channelType" validate:"required"` // Loại kênh: email, telegram, webhook - BẮT BUỘC
	Name        string `json:"name" validate:"required"`        // Tên sender - BẮT BUỘC
	IsActive    bool   `json:"isActive"`                         // Sender có đang hoạt động không - Optional (default: true)

	// Email sender config (chỉ dùng khi channelType = "email")
	SMTPHost     string `json:"smtpHost,omitempty"`     // SMTP host - Optional
	SMTPPort     int    `json:"smtpPort,omitempty"`     // SMTP port - Optional
	SMTPUsername string `json:"smtpUsername,omitempty"` // SMTP username - Optional
	SMTPPassword string `json:"smtpPassword,omitempty"` // SMTP password - Optional
	FromEmail    string `json:"fromEmail,omitempty"`     // Email gửi từ - Optional
	FromName     string `json:"fromName,omitempty"`     // Tên người gửi - Optional

	// Telegram sender config (chỉ dùng khi channelType = "telegram")
	BotToken    string `json:"botToken,omitempty"`    // Telegram bot token - Optional
	BotUsername string `json:"botUsername,omitempty"` // Telegram bot username - Optional

	// Lưu ý: KHÔNG cần gửi ownerOrganizationId - Backend tự động gán từ context (phân quyền dữ liệu, nullable: null = System Organization)
}

// NotificationChannelSenderUpdateInput dùng cho cập nhật notification sender (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi cập nhật sender
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type NotificationChannelSenderUpdateInput struct {
	ChannelType string `json:"channelType"` // Loại kênh: email, telegram, webhook - Optional
	Name        string `json:"name"`        // Tên sender - Optional
	IsActive    *bool  `json:"isActive"`    // Sender có đang hoạt động không - Optional

	// Email sender config
	SMTPHost     string `json:"smtpHost,omitempty"`     // SMTP host - Optional
	SMTPPort     *int   `json:"smtpPort,omitempty"`    // SMTP port - Optional
	SMTPUsername string `json:"smtpUsername,omitempty"` // SMTP username - Optional
	SMTPPassword string `json:"smtpPassword,omitempty"` // SMTP password - Optional
	FromEmail    string `json:"fromEmail,omitempty"`     // Email gửi từ - Optional
	FromName     string `json:"fromName,omitempty"`     // Tên người gửi - Optional

	// Telegram sender config
	BotToken    string `json:"botToken,omitempty"`    // Telegram bot token - Optional
	BotUsername string `json:"botUsername,omitempty"` // Telegram bot username - Optional

	// Lưu ý: KHÔNG thể update ownerOrganizationId - Backend sẽ tự động xóa field này nếu có trong request (bảo mật)
}
