package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Vai trò
type Role struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                                                   // ID của vai trò
	Name           string             `json:"name" bson:"name" index:"compound:role_org_name_unique"`                              // Tên vai trò (unique trong mỗi Organization)
	Describe       string             `json:"describe" bson:"describe"`                                                            // Mô tả vai trò
	OrganizationID primitive.ObjectID `json:"organizationId" bson:"organizationId" index:"single:1,compound:role_org_name_unique"` // **BẮT BUỘC**: Role thuộc Organization nào
	IsSystem       bool               `json:"-" bson:"isSystem" index:"single:1"`                                                  // true = dữ liệu hệ thống, không thể xóa (chỉ dùng nội bộ, không expose ra API)
	CreatedAt      int64              `json:"createdAt" bson:"createdAt"`                                                          // Thời gian tạo
	UpdatedAt      int64              `json:"updatedAt" bson:"updatedAt"`                                                          // Thời gian cập nhật
}
