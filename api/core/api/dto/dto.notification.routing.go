package dto

// NotificationRoutingRuleCreateInput dùng cho tạo notification routing rule (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi tạo routing rule
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type NotificationRoutingRuleCreateInput struct {
	EventType       string   `json:"eventType" validate:"required"`       // Loại event: conversation_unreplied, order_created, ... - BẮT BUỘC
	OrganizationIDs []string `json:"organizationIds" validate:"required"` // Teams nào nhận (có thể nhiều) - BẮT BUỘC - Array of organization IDs
	ChannelTypes    []string `json:"channelTypes,omitempty"`               // Filter channels theo type (optional): ["email", "telegram", "webhook"] - Optional
	IsActive        bool     `json:"isActive"`                             // Rule có đang hoạt động không - Optional (default: true)

	// Lưu ý: KHÔNG cần gửi isSystem - Backend tự động set (chỉ dùng nội bộ)
}

// NotificationRoutingRuleUpdateInput dùng cho cập nhật notification routing rule (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi cập nhật routing rule
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type NotificationRoutingRuleUpdateInput struct {
	EventType       string   `json:"eventType"`       // Loại event - Optional
	OrganizationIDs []string `json:"organizationIds"` // Teams nào nhận - Optional
	ChannelTypes    []string `json:"channelTypes,omitempty"` // Filter channels theo type - Optional
	IsActive        *bool    `json:"isActive"`       // Rule có đang hoạt động không - Optional

	// Lưu ý: KHÔNG thể update isSystem - Backend sẽ tự động xóa field này nếu có trong request (bảo mật)
}
