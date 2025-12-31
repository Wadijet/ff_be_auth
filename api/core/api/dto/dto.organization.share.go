package dto

// OrganizationShareCreateInput dùng cho tạo organization share (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi tạo share
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type OrganizationShareCreateInput struct {
	OwnerOrganizationID string   `json:"ownerOrganizationId" validate:"required"` // Tổ chức sở hữu dữ liệu (phân quyền) - Organization share data - BẮT BUỘC
	ToOrgID             string   `json:"toOrgId" validate:"required"`             // Organization nhận data - BẮT BUỘC
	PermissionNames     []string `json:"permissionNames,omitempty"`              // [] hoặc null = tất cả permissions, ["Order.Read", "Order.Create"] = chỉ share với permissions cụ thể - Optional

	// Lưu ý: KHÔNG cần gửi createdAt, createdBy - Backend tự động set
}

// OrganizationShareUpdateInput dùng cho cập nhật organization share (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi cập nhật share
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
// Lưu ý: OrganizationShare thường không cần update, nhưng nếu có thì chỉ update PermissionNames
type OrganizationShareUpdateInput struct {
	PermissionNames []string `json:"permissionNames,omitempty"` // [] hoặc null = tất cả permissions, ["Order.Read", "Order.Create"] = chỉ share với permissions cụ thể - Optional

	// Lưu ý: KHÔNG thể update ownerOrganizationId, toOrgId - Backend sẽ tự động xóa các fields này nếu có trong request (bảo mật)
	// Lưu ý: KHÔNG thể update createdAt, createdBy - Backend sẽ tự động xóa các fields này nếu có trong request
}
