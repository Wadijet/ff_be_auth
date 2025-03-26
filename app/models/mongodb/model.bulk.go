package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BulkOperation là model chung cho các thao tác bulk
type BulkOperation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      string             `bson:"type" json:"type"`       // Loại thao tác (create, update, delete)
	Status    string             `bson:"status" json:"status"`   // Trạng thái (pending, processing, completed, failed)
	Total     int                `bson:"total" json:"total"`     // Tổng số record
	Success   int                `bson:"success" json:"success"` // Số record thành công
	Failed    int                `bson:"failed" json:"failed"`   // Số record thất bại
	Errors    []string           `bson:"errors" json:"errors"`   // Danh sách lỗi
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
