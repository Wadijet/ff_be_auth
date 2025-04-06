package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
)

// FiberFbPageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Page cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.FbPage
// - CreateInput: models.FbPageCreateInput
// - UpdateInput: models.FbPageCreateInput
type FiberFbPageHandler struct {
	FiberBaseHandler[models.FbPage, models.FbPageCreateInput, models.FbPageCreateInput]
	FbPageService *services.FbPageService
}

// NewFiberFbPageHandler khởi tạo một FiberFbPageHandler mới
// Returns:
//   - *FiberFbPageHandler: Instance mới của FiberFbPageHandler đã được khởi tạo với các service cần thiết
func NewFiberFbPageHandler() *FiberFbPageHandler {
	handler := &FiberFbPageHandler{}
	handler.FbPageService = services.NewFbPageService()
	handler.Service = handler.FbPageService
	return handler
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
func (h *FiberFbPageHandler) HandleFindOneByPageID(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	data, err := h.FbPageService.FindOneByPageID(context.Background(), id)
	if err != nil {
		if customErr, ok := err.(*utility.Error); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"code":    customErr.Code,
				"message": customErr.Message,
				"details": customErr.Details,
			})
		}
		return c.Status(utility.StatusInternalServerError).JSON(fiber.Map{
			"code":    utility.ErrCodeDatabase,
			"message": err.Error(),
		})
	}

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
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
func (h *FiberFbPageHandler) HandleUpdateToken(c fiber.Ctx) error {
	input := new(models.FbPageUpdateTokenInput)
	if response := h.ParseRequestBody(c, input); response != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationInput,
			"message": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	data, err := h.FbPageService.UpdateToken(context.Background(), input)
	if err != nil {
		if customErr, ok := err.(*utility.Error); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"code":    customErr.Code,
				"message": customErr.Message,
				"details": customErr.Details,
			})
		}
		return c.Status(utility.StatusInternalServerError).JSON(fiber.Map{
			"code":    utility.ErrCodeDatabase,
			"message": err.Error(),
		})
	}

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
}
