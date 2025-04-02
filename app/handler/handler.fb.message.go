package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
)

// FbMessageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Message
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbMessageHandler struct {
	BaseHandler[models.FbMessage, models.FbMessageCreateInput, models.FbMessageCreateInput]
	FbMessageService *services.FbMessageService
}

// NewFbMessageHandler khởi tạo một FbMessageHandler mới
func NewFbMessageHandler() *FbMessageHandler {
	newHandler := new(FbMessageHandler)
	newHandler.FbMessageService = services.NewFbMessageService()
	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng FbMessageService trực tiếp
	return newHandler
}

// Các hàm đặc thù của FbMessage (nếu có) sẽ được thêm vào đây
