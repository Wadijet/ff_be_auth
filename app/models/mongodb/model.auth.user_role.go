package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRole đại diện cho vai trò người dùng trong hệ thống.
// ID: ID của vai trò người dùng, được lưu trữ dưới dạng ObjectID của MongoDB.
// UserID: ID của người dùng, được lưu trữ dưới dạng ObjectID của MongoDB.
// RoleID: ID của vai trò, được lưu trữ dưới dạng ObjectID của MongoDB.
// CreatedAt: Thời gian tạo vai trò người dùng, được lưu trữ dưới dạng timestamp.
// UpdatedAt: Thời gian cập nhật vai trò người dùng, được lưu trữ dưới dạng timestamp.
type UserRole struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`     // ID của vai trò người dùng
	UserID    primitive.ObjectID `json:"userId" bson:"userId" index:"single:1"` // ID của người dùng
	RoleID    primitive.ObjectID `json:"roleId" bson:"roleId" index:"single:1"` // ID của vai trò
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`            // Thời gian tạo
	UpdatedAt int64              `json:"updatedAt" bson:"updatedAt"`            // Thời gian cập nhật
}
