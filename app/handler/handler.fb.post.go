package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"

	"github.com/gofiber/fiber/v3"
)

// FbPostHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Post cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.FbPost
// - CreateInput: models.FbPostCreateInput
// - UpdateInput: models.FbPostCreateInput
type FbPostHandler struct {
	BaseHandler[models.FbPost, models.FbPostCreateInput, models.FbPostCreateInput]
	FbPostService *services.FbPostService
}

// NewFbPostHandler khởi tạo một FiberFbPostHandler mới
// Returns:
//   - *FiberFbPostHandler: Instance mới của FiberFbPostHandler đã được khởi tạo với các service cần thiết
func NewFbPostHandler() *FbPostHandler {
	handler := &FbPostHandler{}
	handler.FbPostService = services.NewFbPostService()
	handler.Service = handler.FbPostService
	return handler
}

// OTHER FUNCTIONS ==========================================================================

// HandleFindOneByPostID tìm một FbPost theo PostID
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Tìm thấy FbPost
//   - 404: Không tìm thấy FbPost
//   - 500: Lỗi server
func (h *FbPostHandler) HandleFindOneByPostID(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	data, err := h.FbPostService.FindOneByPostID(context.Background(), id)
	h.HandleResponse(c, data, err)
	return nil
}

// HandleUpdateToken cập nhật access token của một FbPost
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Request Body:
//   - models.FbPostUpdateTokenInput: Thông tin token cần cập nhật
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Cập nhật token thành công
//   - 400: Dữ liệu đầu vào không hợp lệ
//   - 404: Không tìm thấy FbPost
//   - 500: Lỗi server
func (h *FbPostHandler) HandleUpdateToken(c fiber.Ctx) error {
	input := new(models.FbPostUpdateTokenInput)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	data, err := h.FbPostService.UpdateToken(context.Background(), input)
	h.HandleResponse(c, data, err)
	return nil
}
