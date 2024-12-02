package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Organization đại diện cho một tổ chức trong hệ thống
type Organization struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`              // ID của tổ chức
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`           // Tên của tổ chức
	Describe  string             `json:"describe,omitempty" bson:"describe,omitempty"`   // Mô tả tổ chức
	ParentID  primitive.ObjectID `json:"parentId,omitempty" bson:"parentId,omitempty"`   // ID tổ chức cha
	CreatedAt int64              `json:"createdAt,omitempty" bson:"createdAt,omitempty"` // Thời gian tạo
	UpdatedAt int64              `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"` // Thời gian cập nhật
}
