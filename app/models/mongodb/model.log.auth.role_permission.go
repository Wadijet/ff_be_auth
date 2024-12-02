package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Vai trò
type LogRolePermission struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`              // ID của vai trò
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`           // Tên của vai trò
	Describe  string             `json:"describe,omitempty" bson:"describe,omitempty"`   // Mô tả vai trò
	CreatedAt int64              `json:"createdAt,omitempty" bson:"createdAt,omitempty"` // Thời gian tạo
	UpdatedAt int64              `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"` // Thời gian cập nhật
}