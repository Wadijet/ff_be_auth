package dto

// AgentCreateInput dữ liệu đầu vào tạo trợ lý
type AgentCreateInput struct {
	Name          string                 `json:"name" validate:"required"`     // Tên của vai trò
	Describe      string                 `json:"describe" validate:"required"` // Mô tả vai trò
	AssignedUsers []string               `json:"assignedUsers"`                // Danh sách người dùng được gán access token
	ConfigData    map[string]interface{} `json:"configData"`                   // Dữ liệu cấu hình
}

// AgentUpdateInput dữ liệu đầu vào cập nhật trợ lý
type AgentUpdateInput struct {
	Name          string                 `json:"name"`          // Tên của vai trò
	Describe      string                 `json:"describe"`      // Mô tả vai trò
	Status        byte                   `json:"status"`        // Trạng thái của trợ lý (0: offline, 1: online)
	Command       byte                   `json:"command"`       // Lệnh điều khiển trợ lý (0: stop, 1: play)
	AssignedUsers []string               `json:"assignedUsers"` // Danh sách người dùng được gán access token
	ConfigData    map[string]interface{} `json:"configData"`    // Dữ liệu cấu hình
}

