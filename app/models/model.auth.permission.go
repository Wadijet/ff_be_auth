package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission đại diện cho quyền trong hệ thống
type Permission struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`              // ID của quyền
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`           // Tên của quyền
	Describe  string             `json:"describe,omitempty" bson:"describe,omitempty"`   // Mô tả quyền
	Category  string             `json:"category,omitempty" bson:"category,omitempty"`   // Danh mục của quyền
	CreatedAt int64              `json:"createdAt,omitempty" bson:"createdAt,omitempty"` // Thời gian tạo quyền
	UpdatedAt int64              `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"` // Thời gian cập nhật quyền
}

// =======================================================================================

// PermissionCreateInput đại diện cho dữ liệu đầu vào khi tạo quyền
type PermissionCreateInput struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty" validate:"required"`         // Tên của quyền (bắt buộc)
	Describe string `json:"describe,omitempty" bson:"describe,omitempty" validate:"required"` // Mô tả quyền (bắt buộc)
	Category string `json:"category,omitempty" bson:"category,omitempty"`                     // Danh mục của quyền
}

// PermissionUpdateInput đại diện cho dữ liệu đầu vào khi cập nhật quyền
type PermissionUpdateInput struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`         // Tên của quyền
	Describe string `json:"describe,omitempty" bson:"describe,omitempty"` // Mô tả quyền
	Category string `json:"category,omitempty" bson:"category,omitempty"` // Danh mục của quyền
}
