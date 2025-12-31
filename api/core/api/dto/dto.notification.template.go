package dto

// NotificationTemplateCreateInput dùng cho tạo notification template (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi tạo template
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type NotificationTemplateCreateInput struct {
	EventType   string   `json:"eventType" validate:"required"`   // Loại event: conversation_unreplied, order_created, ... - BẮT BUỘC
	ChannelType string   `json:"channelType" validate:"required"` // Loại kênh: email, telegram, webhook - BẮT BUỘC
	Subject     string   `json:"subject,omitempty"`                // Subject (cho email) - Optional
	Content     string   `json:"content" validate:"required"`      // Nội dung template (có thể chứa {{variable}}) - BẮT BUỘC
	Variables   []string `json:"variables,omitempty"`               // Danh sách variables: ["conversationId", "minutes"] - Optional

	// CTA buttons (optional)
	CTAs []NotificationCTACreateInput `json:"ctas,omitempty"` // CTA buttons - Optional

	IsActive bool `json:"isActive"` // Template có đang hoạt động không - Optional (default: true)

	// Lưu ý: KHÔNG cần gửi ownerOrganizationId - Backend tự động gán từ context (phân quyền dữ liệu, nullable: null = System Organization)
}

// NotificationTemplateUpdateInput dùng cho cập nhật notification template (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi cập nhật template
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type NotificationTemplateUpdateInput struct {
	EventType   string   `json:"eventType"`   // Loại event - Optional
	ChannelType string   `json:"channelType"` // Loại kênh - Optional
	Subject     string   `json:"subject,omitempty"` // Subject (cho email) - Optional
	Content     string   `json:"content"`      // Nội dung template - Optional
	Variables   []string `json:"variables,omitempty"` // Danh sách variables - Optional

	// CTA buttons (optional)
	CTAs []NotificationCTACreateInput `json:"ctas,omitempty"` // CTA buttons - Optional

	IsActive *bool `json:"isActive"` // Template có đang hoạt động không - Optional

	// Lưu ý: KHÔNG thể update ownerOrganizationId - Backend sẽ tự động xóa field này nếu có trong request (bảo mật)
}

// NotificationCTACreateInput dùng cho CTA button trong template
type NotificationCTACreateInput struct {
	Label  string `json:"label" validate:"required"`  // Label của CTA: "Xem chi tiết", "Phản hồi", "Đã xem" - BẮT BUỘC
	Action string `json:"action" validate:"required"` // URL (có thể chứa {{variable}}) - BẮT BUỘC
	Style  string `json:"style,omitempty"`            // Style: "primary", "success", "secondary" (chỉ để styling) - Optional
}
