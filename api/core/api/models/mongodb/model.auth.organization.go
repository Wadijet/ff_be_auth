package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationType định nghĩa các loại tổ chức
const (
	OrganizationTypeSystem     = "system"     // Hệ thống (Level -1) - Cấp cao nhất, chứa Administrator, không thể xóa
	OrganizationTypeGroup      = "group"      // Tập đoàn (Level 0)
	OrganizationTypeCompany    = "company"    // Công ty (Level 1)
	OrganizationTypeDepartment = "department" // Phòng ban (Level 2)
	OrganizationTypeDivision   = "division"   // Bộ phận (Level 3)
	OrganizationTypeTeam       = "team"      // Team (Level 4+)
)

// Organization đại diện cho cấu trúc tổ chức hình cây (Hệ thống, Tập đoàn, Công ty, Phòng ban, Bộ phận, Team)
type Organization struct {
	ID        primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`                        // ID của tổ chức
	Name      string              `json:"name" bson:"name" index:"single:1"`                        // Tên tổ chức
	Code      string              `json:"code" bson:"code" index:"unique"`                          // Mã tổ chức (unique)
	Type      string              `json:"type" bson:"type" index:"single:1"`                        // Loại tổ chức (system, group, company, department, division, team)
	ParentID  *primitive.ObjectID `json:"parentId,omitempty" bson:"parentId,omitempty" index:"single:1"` // ID tổ chức cha (null nếu là root system)
	Path      string              `json:"path" bson:"path" index:"single:1"`                        // Đường dẫn cây (ví dụ: "/system/root_group/company1/dept1")
	Level     int                 `json:"level" bson:"level" index:"single:1"`                      // Cấp độ (-1 = system root, 0 = group, 1 = company, 2 = department, ...)
	IsActive  bool                `json:"isActive" bson:"isActive" index:"single:1"`                // Trạng thái hoạt động
	CreatedAt int64               `json:"createdAt" bson:"createdAt"`                               // Thời gian tạo
	UpdatedAt int64               `json:"updatedAt" bson:"updatedAt"`                               // Thời gian cập nhật
}

