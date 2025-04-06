package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
)

// FiberFbPostHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Post cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.FbPost
// - CreateInput: models.FbPostCreateInput
// - UpdateInput: models.FbPostCreateInput
type FiberFbPostHandler struct {
	FiberBaseHandler[models.FbPost, models.FbPostCreateInput, models.FbPostCreateInput]
	FbPostService *services.FbPostService
}

// NewFiberFbPostHandler khởi tạo một FiberFbPostHandler mới
// Returns:
//   - *FiberFbPostHandler: Instance mới của FiberFbPostHandler đã được khởi tạo với các service cần thiết
func NewFiberFbPostHandler() *FiberFbPostHandler {
	handler := &FiberFbPostHandler{}
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
func (h *FiberFbPostHandler) HandleFindOneByPostID(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	data, err := h.FbPostService.FindOneByPostID(context.Background(), id)
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
func (h *FiberFbPostHandler) HandleUpdateToken(c fiber.Ctx) error {
	input := new(models.FbPostUpdateTokenInput)
	if response := h.ParseRequestBody(c, input); response != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationInput,
			"message": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	data, err := h.FbPostService.UpdateToken(context.Background(), input)
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
