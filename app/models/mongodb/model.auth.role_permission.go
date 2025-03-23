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
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                 // ID của quyền vai trò
	RoleID          primitive.ObjectID `json:"roleId" bson:"roleId" index:"single:1"`             // ID của vai trò
	PermissionID    primitive.ObjectID `json:"permissionId" bson:"permissionId" index:"single:1"` // ID của quyền
	Scope           byte               `json:"scope" bson:"scope"`                                // Phạm vi của quyền (0: All, 1: Assign)
	CreatedByRoleID primitive.ObjectID `json:"createdByRoleId" bson:"createdByRoleId"`            // ID của vai trò tạo quyền này
	CreatedByUserID primitive.ObjectID `json:"createdByUserId" bson:"createdByUserId"`            // ID của người dùng tạo quyền này
	CreatedAt       int64              `json:"createdAt" bson:"createdAt"`                        // Thời gian tạo
	UpdatedAt       int64              `json:"updatedAt" bson:"updatedAt"`                        // Thời gian cập nhật
}

// API INPUT STRUCT ==========================================================================================

// RolePermissionCreateInput đại diện cho dữ liệu đầu vào khi tạo quyền vai trò
type RolePermissionCreateInput struct {
	RoleID       primitive.ObjectID `json:"roleId" bson:"roleId" validate:"required"`             // ID của vai trò (bắt buộc)
	PermissionID primitive.ObjectID `json:"permissionId" bson:"permissionId" validate:"required"` // ID của quyền (bắt buộc)
	Scope        byte               `json:"scope" bson:"scope"`                                   // Phạm vi của quyền (0: All, 1: Assign)
}

// RolePermissionUpdateInput dữ liệu đầu vào khi cập nhật quyền của vai trò
type RolePermissionUpdateInput struct {
	RoleId        string   `json:"roleId" bson:"roleId" validate:"required"`               // ID của vai trò
	PermissionIds []string `json:"permissionIds" bson:"permissionIds" validate:"required"` // Danh sách ID của quyền
}
