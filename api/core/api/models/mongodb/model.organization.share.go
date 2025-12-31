package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationShare đại diện cho việc share dữ liệu giữa các organizations
// Organization A có thể share tất cả data của mình với Organization B
type OrganizationShare struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerOrganizationID primitive.ObjectID `json:"ownerOrganizationId" bson:"ownerOrganizationId" index:"single:1"` // Tổ chức sở hữu dữ liệu (phân quyền) - Organization share data với ToOrgID
	ToOrgID             primitive.ObjectID `json:"toOrgId" bson:"toOrgId" index:"single:1"`                         // Organization nhận data
	PermissionNames     []string           `json:"permissionNames,omitempty" bson:"permissionNames,omitempty"`      // [] hoặc nil = tất cả permissions, ["Order.Read", "Order.Create"] = chỉ share với permissions cụ thể
	CreatedAt           int64              `json:"createdAt" bson:"createdAt"`
	CreatedBy           primitive.ObjectID `json:"createdBy" bson:"createdBy"`
}
