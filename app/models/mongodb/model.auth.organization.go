package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Organization đại diện cho một tổ chức trong hệ thống
type Organization struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`              // ID của tổ chức
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`           // Tên của tổ chức
	Describe  string             `json:"describe,omitempty" bson:"describe,omitempty"`   // Mô tả tổ chức
	ParentID  primitive.ObjectID `json:"parentId,omitempty" bson:"parentId,omitempty"`   // ID tổ chức 
	Level     int                `json:"level,omitempty" bson:"level,omitempty"`         // Cấp độ tổ chức
	CreatedAt int64              `json:"createdAt,omitempty" bson:"createdAt,omitempty"` // Thời gian tạo
	UpdatedAt int64              `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"` // Thời gian cập nhật
}

// OrganizationCreateInput đại diện cho dữ liệu đầu vào khi tạo tổ chức
type OrganizationCreateInput struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty" validate:"required"`         // Tên của tổ chức (bắt buộc)
	Describe string `json:"describe,omitempty" bson:"describe,omitempty" validate:"required"` // Mô tả tổ chức (bắt buộc)
	ParentID string `json:"parentId,omitempty" bson:"parentId,omitempty"`                     // ID tổ chức cha
}

// OrganizationUpdateInput đại diện cho dữ liệu đầu vào khi cập nhật tổ chức
type OrganizationUpdateInput struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`         // Tên của tổ chức
	Describe string `json:"describe,omitempty" bson:"describe,omitempty"` // Mô tả tổ chức
	ParentID string `json:"parentId,omitempty" bson:"parentId,omitempty"` // ID tổ chức cha
}


