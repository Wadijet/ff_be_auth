package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PcOrder lưu thông tin đơn hàng
type PcOrder struct {
	ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`                        // ID của đơn hàng
	PancakeOrderId string                 `json:"pancakeOrderId" bson:"pancakeOrderId" index:"unique;text"` // ID của đơn hàng trên hệ thống Pancake
	Status         byte                   `json:"status" bson:"status"`                                     // Trạng thái của đơn hàng (0 = active, 1 = inactive)
	PanCakeData    map[string]interface{} `json:"panCakeData" bson:"panCakeData"`                           // Dữ liệu API
	CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`                               // Thời gian tạo order
	UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`                               // Thời gian cập nhật order
}
