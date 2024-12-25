package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission đại diện cho quyền trong hệ thống,
// Các quyền được kết cấu theo các quyền gọi các API trong router.
// Các quyèn này được tạo ra khi khởi tạo hệ thống và không thể thay đổi.
type FbPage struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`      // ID của quyền
	Name            string             `json:"name" bson:"name" index:"unique"`        // Tên của quyền
	Describe        string             `json:"describe" bson:"describe"`               // Mô tả quyền
	PageName        string             `json:"pageName" bson:"pageName"`               // Tên của trang
	PageID          string             `json:"pageID" bson:"pageID"`                   // ID của trang
	PageAccessToken string             `json:"pageAccessToken" bson:"pageAccessToken"` // Mã truy cập của trang
	ApiData         []byte             `json:"apiData" bson:"apiData"`                 // Dữ liệu API
	CreatedAt       int64              `json:"createdAt" bson:"createdAt"`             // Thời gian tạo quyền
	UpdatedAt       int64              `json:"updatedAt" bson:"updatedAt"`             // Thời gian cập nhật quyền
}

// API INPUT STRUCT =======================================================================================
