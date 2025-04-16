package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Lưu lại log các hành động trong nhóm chức năng AUTH
type AuthLog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`              // ID của vai trò
	UserID    primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`       // ID của người dùng
	RoleID   primitive.ObjectID `json:"roleId,omitempty" bson:"roleId,omitempty"`       // ID của vai trò
	OrganizationID primitive.ObjectID `json:"organizationId,omitempty" bson:"organizationId,omitempty"` // ID tổ chức
	Collection string             `json:"collection,omitempty" bson:"collection,omitempty"` // Tên bảng	
	Action    string             `json:"action,omitempty" bson:"action,omitempty"`       // Hành động
	Describe  string             `json:"describe,omitempty" bson:"describe,omitempty"`   // Mô tả hành động
	OldData   string             `json:"oldData,omitempty" bson:"oldData,omitempty"`     // Dữ liệu cũ
	NewData   string             `json:"newData,omitempty" bson:"newData,omitempty"`     // Dữ liệu mới
	CreatedAt int64              `json:"createdAt,omitempty" bson:"createdAt,omitempty"` // Thời gian tạo
	UpdatedAt int64              `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"` // Thời gian cập nhật
}