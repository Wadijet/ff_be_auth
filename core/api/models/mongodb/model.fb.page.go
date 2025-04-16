package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission đại diện cho quyền trong hệ thống,
// Các quyền được kết cấu theo các quyền gọi các API trong router.
// Các quyèn này được tạo ra khi khởi tạo hệ thống và không thể thay đổi.
type FbPage struct {
	ID              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`        // ID của quyền
	PageName        string                 `json:"pageName" bson:"pageName"`                 // Tên của trang
	PageUsername    string                 `json:"pageUsername" bson:"pageUsername"`         // Tên người dùng của trang
	PageId          string                 `json:"pageId" bson:"pageId" index:"unique;text"` // ID của trang
	IsSync          bool                   `json:"isSync" bson:"isSync"`                     // Trạng thái đồng bộ
	AccessToken     string                 `json:"accessToken" bson:"accessToken"`
	PageAccessToken string                 `json:"pageAccessToken" bson:"pageAccessToken"` // Mã truy cập của trang
	PanCakeData     map[string]interface{} `json:"panCakeData" bson:"panCakeData"`         // Dữ liệu API
	CreatedAt       int64                  `json:"createdAt" bson:"createdAt"`             // Thời gian tạo quyền
	UpdatedAt       int64                  `json:"updatedAt" bson:"updatedAt"`             // Thời gian cập nhật quyền
}

// API INPUT STRUCT =======================================================================================

// FbPageCreateInput dữ liệu đầu vào khi tạo page
type FbPageCreateInput struct {
	AccessToken string                 `json:"accessToken" bson:"accessToken" validate:"required"` // Mã truy cập của trang
	PanCakeData map[string]interface{} `json:"panCakeData" bson:"panCakeData" validate:"required"` // Dữ liệu API
}

type FbPageUpdateTokenInput struct {
	PageId          string `json:"pageId" bson:"pageId" validate:"required" validate:"required"` // ID của trang
	PageAccessToken string `json:"pageAccessToken" bson:"pageAccessToken" validate:"required"`   // Mã truy cập của trang
}
