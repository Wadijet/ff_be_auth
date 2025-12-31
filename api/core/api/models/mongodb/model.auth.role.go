package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Vai trò
type Role struct {
	_Relationships struct{}          `relationship:"collection:user_roles,field:roleId,message:Không thể xóa role vì có %d user đang sử dụng role này. Vui lòng gỡ role khỏi các user trước.|collection:role_permissions,field:roleId,message:Không thể xóa role vì có %d permission đang được gán cho role này. Vui lòng gỡ các permission trước."` // Relationship definitions - không export, chỉ dùng cho tag parsing
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                                                                                                                      // ID của vai trò
	Name               string             `json:"name" bson:"name" index:"compound:role_org_name_unique"`                                                                                                 // Tên vai trò (unique trong mỗi Organization)
	Describe           string             `json:"describe" bson:"describe"`                                                                                                                                 // Mô tả vai trò
	OwnerOrganizationID primitive.ObjectID `json:"ownerOrganizationId" bson:"ownerOrganizationId" index:"single:1,compound:role_org_name_unique"`                                                        // Tổ chức sở hữu dữ liệu (phân quyền) + Logic business - Có thể chỉ định khi create, có thể update với validation quyền
	IsSystem       bool               `json:"-" bson:"isSystem" index:"single:1"`                                                                                                                   // true = dữ liệu hệ thống, không thể xóa (chỉ dùng nội bộ, không expose ra API)
	CreatedAt      int64              `json:"createdAt" bson:"createdAt"`                                                                                                                              // Thời gian tạo
	UpdatedAt      int64              `json:"updatedAt" bson:"updatedAt"`                                                                                                                              // Thời gian cập nhật
}
