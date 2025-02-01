package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Trợ lý
// Một trợ lý sẽ được cấp thông tin đăng nhâp của user để thực hiện các hành động như một user
// Trong trạng thái hoạt động, các Trợ lý sẽ thường xuyên điểm danh để cập nhật trạng thái hoạt động
type Agent struct {
	ID            primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`        // ID của vai trò
	Name          string                 `json:"name" bson:"name" index:"unique"`          // Tên của vai trò
	Describe      string                 `json:"describe" bson:"describe"`                 // Mô tả vai trò
	Status        byte                   `json:"status" bson:"status" index:"single:1;"`   // Trạng thái của trợ lý (0: offline, 1: online)
	Command       byte                   `json:"command" bson:"command" index:"single:1;"` // Lệnh điều khiển trợ lý (0: stop, 1: play)
	AssignedUsers []primitive.ObjectID   `json:"assignedUsers" bson:"assignedUsers"`       // Danh sách người dùng được gán access token
	CreatedAt     int64                  `json:"createdAt" bson:"createdAt"`               // Thời gian tạo
	UpdatedAt     int64                  `json:"updatedAt" bson:"updatedAt"`               // Thời gian cập nhật
	ConfigData    map[string]interface{} `json:"configData" bson:"configData"`             // Dữ liệu cấu hình
}

// API INPUT STRUCT ==========================================================================================
// Dữ liệu đầu vào tạo trợ lý
type AgentCreateInput struct {
	Name          string                 `json:"name" bson:"name" validate:"required"`         // Tên của vai trò
	Describe      string                 `json:"describe" bson:"describe" validate:"required"` // Mô tả vai trò
	AssignedUsers []string               `json:"assignedUsers" bson:"assignedUsers"`           // Danh sách người dùng được gán access token
	ConfigData    map[string]interface{} `json:"configData" bson:"configData"`                 // Dữ liệu cấu hình
}

// Dữ liệu đầu vào cập nhật trợ lý
type AgentUpdateInput struct {
	Name          string                 `json:"name" bson:"name"`                   // Tên của vai trò
	Describe      string                 `json:"describe" bson:"describe"`           // Mô tả vai trò
	Status        byte                   `json:"status" bson:"status"`               // Trạng thái của trợ lý (0: offline, 1: online)
	Command       byte                   `json:"command" bson:"command"`             // Lệnh điều khiển trợ lý (0: stop, 1: play)
	AssignedUsers []string               `json:"assignedUsers" bson:"assignedUsers"` // Danh sách người dùng được gán access token
	ConfigData    map[string]interface{} `json:"configData" bson:"configData"`       // Dữ liệu cấu hình
}

// ==========================================================================================================
