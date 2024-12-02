package models

import (
	"go.mongodb.org/mongo-driver/bson"
)

// PaginateResult đại diện cho kết quả phân trang
type PaginateResult struct {
	// Trang hiện tại
	Page int64 `json:"page" bson:"page"`
	// Số lượng mục trên mỗi trang
	Limit int64 `json:"limit" bson:"limit"`
	// Số lượng mục trong trang hiện tại
	ItemCount int64 `json:"itemCount" bson:"itemCount"`
	// Danh sách các mục
	Items []bson.M `json:"items" bson:"items"`
}

// CountResult đại diện cho kết quả đếm
type CountResult struct {
	// Tổng số lượng mục
	TotalCount int64 `json:"totalCount" bson:"totalCount"`
	// Số lượng mục trên mỗi trang
	Limit int64 `json:"limit" bson:"limit"`
	// Tổng số trang
	TotalPage int64 `json:"totalPage" bson:"totalPage"`
}
