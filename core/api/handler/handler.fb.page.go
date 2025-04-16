package handler

import (
	"context"
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"

	"github.com/gofiber/fiber/v3"
)

// FbPageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Page cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.FbPage
// - CreateInput: models.FbPageCreateInput
// - UpdateInput: models.FbPageCreateInput
type FbPageHandler struct {
	BaseHandler[models.FbPage, models.FbPageCreateInput, models.FbPageCreateInput]
	FbPageService *services.FbPageService
}

// NewFbPageHandler khởi tạo một FiberFbPageHandler mới
// Returns:
//   - *FiberFbPageHandler: Instance mới của FiberFbPageHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewFbPageHandler() (*FbPageHandler, error) {
	handler := &FbPageHandler{}

	// Khởi tạo FbPageService và xử lý error
	service, err := services.NewFbPageService()
	if err != nil {
		return nil, fmt.Errorf("failed to create page service: %v", err)
	}
	handler.FbPageService = service
	handler.BaseService = handler.FbPageService

	return handler, nil
}

// OTHER FUNCTIONS ==========================================================================

// HandleFindOneByPageID tìm một FbPage theo PageID
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Tìm thấy FbPage
//   - 404: Không tìm thấy FbPage
//   - 500: Lỗi server
func (h *FbPageHandler) HandleFindOneByPageID(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	data, err := h.FbPageService.FindOneByPageID(context.Background(), id)
	h.HandleResponse(c, data, err)
	return nil
}

// HandleUpdateToken cập nhật access token của một FbPage
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Request Body:
//   - models.FbPageUpdateTokenInput: Thông tin token cần cập nhật
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Cập nhật token thành công
//   - 400: Dữ liệu đầu vào không hợp lệ
//   - 404: Không tìm thấy FbPage
//   - 500: Lỗi server
func (h *FbPageHandler) HandleUpdateToken(c fiber.Ctx) error {
	input := new(models.FbPageUpdateTokenInput)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	data, err := h.FbPageService.UpdateToken(context.Background(), input)
	h.HandleResponse(c, data, err)
	return nil
}
