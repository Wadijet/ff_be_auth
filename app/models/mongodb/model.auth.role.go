package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Vai trò
type Role struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`              // ID của vai trò
	Name      string             `json:"name" bson:"name"`           // Tên của vai trò
	Describe  string             `json:"describe" bson:"describe"`   // Mô tả vai trò
	CreatedAt int64              `json:"createdAt" bson:"createdAt"` // Thời gian tạo
	UpdatedAt int64              `json:"updatedAt" bson:"updatedAt"` // Thời gian cập nhật
}

// ==========================================================================================

// Dữ liệu đầu vào tạo vai trò
type RoleCreateInput struct {
	Name     string `json:"name" bson:"name" validate:"required"`         // Tên của vai trò
	Describe string `json:"describe" bson:"describe" validate:"required"` // Mô tả vai trò
}

// Dữ liệu đầu vào cập nhật vai trò
type RoleUpdateInput struct {
	Name     string `json:"name" bson:"name"`         // Tên của vai trò
	Describe string `json:"describe" bson:"describe"` // Mô tả vai trò
}
