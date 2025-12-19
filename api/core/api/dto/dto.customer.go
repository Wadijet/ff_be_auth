package dto

// CustomerCreateInput dữ liệu đầu vào khi tạo customer
type CustomerCreateInput struct {
	PanCakeData map[string]interface{} `json:"panCakeData" validate:"required"`
}
