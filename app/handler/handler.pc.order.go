package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/config"

	"go.mongodb.org/mongo-driver/mongo"
)

// PcOrderHandler là cấu trúc xử lý các yêu cầu liên quan đến đơn hàng
// Kế thừa từ BaseHandler với các type parameter:
// - Model: models.PcOrder
// - CreateInput: models.PcOrderCreateInput
// - UpdateInput: models.PcOrderCreateInput
type PcOrderHandler struct {
	BaseHandler[models.PcOrder, models.PcOrderCreateInput, models.PcOrderCreateInput]
	PcOrderService *services.PcOrderService
}

// NewPcOrderHandler khởi tạo một PcOrderHandler mới
func NewPcOrderHandler(c *config.Configuration, db *mongo.Client) *PcOrderHandler {
	handler := &PcOrderHandler{}
	handler.PcOrderService = services.NewPcOrderService(c, db)
	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng PcOrderService trực tiếp
	return handler
}

// Các hàm đặc thù của PcOrder (nếu có) sẽ được thêm vào đây
