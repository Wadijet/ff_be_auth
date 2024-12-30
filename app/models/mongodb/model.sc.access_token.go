package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AccessToken lưu các access tokens để truy cập vào các hệ thống khác nhau
type AccessToken struct {
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`  // ID của access token
	Name          string               `json:"name" bson:"name" index:"unique"`    // Tên của access token
	Describe      string               `json:"describe" bson:"describe"`           // Mô tả access token
	System        string               `json:"system" bson:"system"`               // Hệ thống của access token
	Value         string               `json:"value" bson:"value"`                 // Giá trị của access token
	AssignedUsers []primitive.ObjectID `json:"assignedUsers" bson:"assignedUsers"` // Danh sách người dùng được gán access token
	CreatedAt     int64                `json:"createdAt" bson:"createdAt"`         // Thời gian tạo access token
	UpdatedAt     int64                `json:"updatedAt" bson:"updatedAt"`         // Thời gian cập nhật access token
}

// AccessTokenCreateInput dữ liệu đầu vào khi tạo access token
type AccessTokenCreateInput struct {
	Name          string   `json:"name" bson:"name" validate:"required"`         // Tên của access token
	Describe      string   `json:"describe" bson:"describe" validate:"required"` // Mô tả access token
	System        string   `json:"system" bson:"system" validate:"required"`     // Hệ thống của access token
	Value         string   `json:"value" bson:"value" validate:"required"`       // Giá trị của access token
	AssignedUsers []string `json:"assignedUsers" bson:"assignedUsers"`           // Danh sách người dùng được gán access token
}

// AccessTokenUpdateInput dữ liệu đầu vào khi cập nhật access token
type AccessTokenUpdateInput struct {
	Name          string   `json:"name" bson:"name"`                   // Tên của access token
	Describe      string   `json:"describe" bson:"describe"`           // Mô tả access token
	System        string   `json:"system" bson:"system"`               // Hệ thống của access token
	Value         string   `json:"value" bson:"value"`                 // Giá trị của access token
	AssignedUsers []string `json:"assignedUsers" bson:"assignedUsers"` // Danh sách người dùng được gán access token
}

// ==========================================================================================================
