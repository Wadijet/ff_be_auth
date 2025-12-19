package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission đại diện cho quyền trong hệ thống,
// Các quyền được kết cấu theo các quyền gọi các API trong router.
// Các quyèn này được tạo ra khi khởi tạo hệ thống và không thể thay đổi.
type FbPage struct {
	ID              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`                                   // ID của quyền
	PageName        string                 `json:"pageName" bson:"pageName" extract:"PanCakeData\\.name"`               // Tên của trang (extract từ PanCakeData["name"])
	PageUsername    string                 `json:"pageUsername" bson:"pageUsername" extract:"PanCakeData\\.username"`   // Tên người dùng của trang (extract từ PanCakeData["username"])
	PageId          string                 `json:"pageId" bson:"pageId" index:"unique;text" extract:"PanCakeData\\.id"` // ID của trang (extract từ PanCakeData["id"])
	IsSync          bool                   `json:"isSync" bson:"isSync"`                                                // Trạng thái đồng bộ
	AccessToken     string                 `json:"accessToken" bson:"accessToken"`
	PageAccessToken string                 `json:"pageAccessToken" bson:"pageAccessToken"` // Mã truy cập của trang
	PanCakeData     map[string]interface{} `json:"panCakeData" bson:"panCakeData"`         // Dữ liệu API
	CreatedAt       int64                  `json:"createdAt" bson:"createdAt"`             // Thời gian tạo quyền
	UpdatedAt       int64                  `json:"updatedAt" bson:"updatedAt"`             // Thời gian cập nhật quyền
}
