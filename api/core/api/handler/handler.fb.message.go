package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"github.com/gofiber/fiber/v3"
)

// FbMessageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Message cho Fiber
// Kế thừa từ FiberBaseHandler để sử dụng các phương thức xử lý chung
type FbMessageHandler struct {
	BaseHandler[models.FbMessage, dto.FbMessageCreateInput, dto.FbMessageCreateInput]
	FbMessageService *services.FbMessageService
}

// NewFbMessageHandler khởi tạo một FiberFbMessageHandler mới
// Returns:
//   - *FiberFbMessageHandler: Instance mới của FiberFbMessageHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewFbMessageHandler() (*FbMessageHandler, error) {
	// Khởi tạo FbMessageService và xử lý error
	service, err := services.NewFbMessageService()
	if err != nil {
		return nil, fmt.Errorf("failed to create message service: %v", err)
	}

	handler := &FbMessageHandler{
		BaseHandler:      *NewBaseHandler[models.FbMessage, dto.FbMessageCreateInput, dto.FbMessageCreateInput](service.BaseServiceMongoImpl),
		FbMessageService: service,
	}

	return handler, nil
}

// ============================================
// CRUD METHODS (Kế thừa từ BaseHandler)
// ============================================
// Các phương thức CRUD hoạt động bình thường, không có logic tách messages:
// - HandleCreate: Tạo mới một FbMessage (dùng FbMessageCreateInput)
// - HandleUpdate: Cập nhật một FbMessage (dùng FbMessageCreateInput)
// - HandleDelete: Xóa một FbMessage
// - HandleFindById: Tìm FbMessage theo ID
// - HandleFindAll: Lấy danh sách FbMessage có phân trang
//
// Lưu ý: CRUD routes KHÔNG tự động tách messages[] ra khỏi panCakeData
// Nếu cần tách messages, sử dụng endpoint đặc biệt HandleUpsertMessages

// ============================================
// ENDPOINT ĐẶC BIỆT: Upsert Messages
// ============================================
// HandleUpsertMessages xử lý upsert messages từ Pancake API
// Endpoint: POST /api/v1/facebook/message/upsert-messages
// DTO: FbMessageUpsertMessagesInput (có field HasMore)
//
// Logic nội bộ (TÁCH BIỆT với CRUD):
// 1. Tự động tách messages[] ra khỏi panCakeData
// 2. Lưu metadata (panCakeData không có messages[]) vào fb_messages
// 3. Lưu từng message riêng lẻ vào fb_message_items
// 4. Cập nhật totalMessages trong fb_messages
//
// Khác biệt với CRUD:
// - CRUD: Lưu nguyên panCakeData (có thể có messages[])
// - Endpoint này: Tách messages[] và lưu vào 2 collections riêng
func (h *FbMessageHandler) HandleUpsertMessages(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		// Parse request body
		var input dto.FbMessageUpsertMessagesInput
		if err := h.ParseRequestBody(c, &input); err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Dữ liệu gửi lên không đúng định dạng JSON. Chi tiết: %v", err),
				common.StatusBadRequest,
				err,
			))
			return nil
		}

		// Validate
		if err := global.Validate.Struct(input); err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Dữ liệu không hợp lệ: %v", err),
				common.StatusBadRequest,
				err,
			))
			return nil
		}

		// Gọi service để upsert (logic tách messages được xử lý trong service)
		result, err := h.FbMessageService.UpsertMessages(
			c.Context(),
			input.ConversationId,
			input.PageId,
			input.PageUsername,
			input.CustomerId,
			input.PanCakeData, // PanCakeData đầy đủ (bao gồm messages[])
			input.HasMore,
		)

		h.HandleResponse(c, result, err)
		return nil
	})
}
