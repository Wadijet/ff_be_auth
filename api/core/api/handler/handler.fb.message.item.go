package handler

import (
	"context"
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// FbMessageItemHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Message Item cho Fiber
// Kế thừa từ BaseHandler với các type parameter:
// - Model: models.FbMessageItem
// - CreateInput: dto.FbMessageItemCreateInput
// - UpdateInput: dto.FbMessageItemUpdateInput
type FbMessageItemHandler struct {
	BaseHandler[models.FbMessageItem, dto.FbMessageItemCreateInput, dto.FbMessageItemUpdateInput]
	FbMessageItemService *services.FbMessageItemService
}

// NewFbMessageItemHandler khởi tạo một FbMessageItemHandler mới
// Returns:
//   - *FbMessageItemHandler: Instance mới của FbMessageItemHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewFbMessageItemHandler() (*FbMessageItemHandler, error) {
	// Khởi tạo FbMessageItemService và xử lý error
	service, err := services.NewFbMessageItemService()
	if err != nil {
		return nil, fmt.Errorf("failed to create message item service: %v", err)
	}

	handler := &FbMessageItemHandler{
		BaseHandler:          *NewBaseHandler[models.FbMessageItem, dto.FbMessageItemCreateInput, dto.FbMessageItemUpdateInput](service.BaseServiceMongoImpl),
		FbMessageItemService: service,
	}

	return handler, nil
}

// OTHER FUNCTIONS ==========================================================================

// HandleFindByConversationId tìm tất cả message items của một conversation với phân trang
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Danh sách message items với phân trang
//   - 400: Lỗi validation (conversationId không hợp lệ)
//   - 500: Lỗi server
func (h *FbMessageItemHandler) HandleFindByConversationId(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		conversationId := c.Params("conversationId")
		if conversationId == "" {
			h.HandleResponse(c, nil, fmt.Errorf("conversationId không được để trống"))
			return nil
		}

		// Lấy query parameters cho phân trang
		page, err := strconv.ParseInt(c.Query("page", "1"), 10, 64)
		if err != nil || page < 1 {
			page = 1
		}
		limit, err := strconv.ParseInt(c.Query("limit", "50"), 10, 64)
		if err != nil || limit < 1 || limit > 100 {
			limit = 50
		}

		// Gọi service để lấy messages
		messages, total, err := h.FbMessageItemService.FindByConversationId(
			context.Background(),
			conversationId,
			int64(page),
			int64(limit),
		)

		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		// Trả về kết quả với phân trang
		result := map[string]interface{}{
			"data": messages,
			"pagination": map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		}

		h.HandleResponse(c, result, nil)
		return nil
	})
}

// HandleFindOneByMessageId tìm một FbMessageItem theo MessageId
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Tìm thấy FbMessageItem
//   - 404: Không tìm thấy FbMessageItem
//   - 500: Lỗi server
func (h *FbMessageItemHandler) HandleFindOneByMessageId(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		messageId := c.Params("messageId")
		if messageId == "" {
			h.HandleResponse(c, nil, fmt.Errorf("messageId không được để trống"))
			return nil
		}

		// Tìm message item theo messageId
		filter := map[string]interface{}{
			"messageId": messageId,
		}

		data, err := h.FbMessageItemService.FindOne(context.Background(), filter, nil)
		h.HandleResponse(c, data, err)
		return nil
	})
}
