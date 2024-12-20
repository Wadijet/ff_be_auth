package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RolePermission đại diện cho quyền vai trò trong hệ thống.
// ID: ID của quyền vai trò, được lưu trữ dưới dạng ObjectID của MongoDB.
// RoleID: ID của vai trò, được lưu trữ dưới dạng ObjectID của MongoDB.
// PermissionID: ID của quyền, được lưu trữ dưới dạng ObjectID của MongoDB.
// CreatedAt: Thời gian tạo quyền vai trò, được lưu trữ dưới dạng timestamp.
// UpdatedAt: Thời gian cập nhật quyền vai trò, được lưu trữ dưới dạng timestamp.
type RolePermission struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // ID của quyền vai trò
	RoleID       primitive.ObjectID `json:"roleId" bson:"roleId"`              // ID của vai trò
	PermissionID primitive.ObjectID `json:"permissionId" bson:"permissionId"`  // ID của quyền
	Scope        byte               `json:"scope" bson:"scope"`                // Phạm vi của quyền (0: All, 1: Assign)
	CreatedAt    int64              `json:"createdAt" bson:"createdAt"`        // Thời gian tạo
	UpdatedAt    int64              `json:"updatedAt" bson:"updatedAt"`        // Thời gian cập nhật
}

// ==========================================================================================
