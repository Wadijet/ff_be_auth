package dto

// RoleCreateInput dùng cho tạo vai trò (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi tạo role
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type RoleCreateInput struct {
	Name                string `json:"name" validate:"required"`      // Tên của vai trò - BẮT BUỘC
	Describe            string `json:"describe" validate:"required"`  // Mô tả vai trò - BẮT BUỘC
	OwnerOrganizationID string `json:"ownerOrganizationId,omitempty"` // Tổ chức sở hữu dữ liệu (phân quyền) - Optional, nếu không có → dùng context
	// Lưu ý: Nếu có ownerOrganizationId trong request, backend sẽ validate quyền với organization đó
}

// RoleUpdateInput dùng cho cập nhật vai trò (tầng transport)
// Đây là contract/interface cho Frontend - định nghĩa cấu trúc dữ liệu cần gửi khi cập nhật role
// Lưu ý: Backend parse trực tiếp vào Model, nhưng DTO này dùng để Frontend biết cấu trúc cần gửi
type RoleUpdateInput struct {
	Name                string `json:"name"`                          // Tên của vai trò - Optional
	Describe            string `json:"describe"`                      // Mô tả vai trò - Optional
	OwnerOrganizationID string `json:"ownerOrganizationId,omitempty"` // Tổ chức sở hữu dữ liệu (phân quyền) - Optional, có thể update với validation quyền
	// Lưu ý: Nếu update ownerOrganizationId, backend sẽ validate quyền với organization mới và document hiện tại
}
