package dto

// AccessTokenCreateInput dữ liệu đầu vào khi tạo access token
type AccessTokenCreateInput struct {
	Name          string   `json:"name" validate:"required"`
	Describe      string   `json:"describe" validate:"required"`
	System        string   `json:"system" validate:"required"`
	Value         string   `json:"value" validate:"required"`
	AssignedUsers []string `json:"assignedUsers"`
}

// AccessTokenUpdateInput dữ liệu đầu vào khi cập nhật access token
type AccessTokenUpdateInput struct {
	Name          string   `json:"name"`
	Describe      string   `json:"describe"`
	System        string   `json:"system"`
	Value         string   `json:"value"`
	AssignedUsers []string `json:"assignedUsers"`
}

