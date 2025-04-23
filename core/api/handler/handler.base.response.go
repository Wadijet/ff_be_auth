package handler

import (
	"errors"
	"fmt"
	"meta_commerce/core/common"
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
)

// SafeHandler bọc các handler với recover để bắt panic và xử lý lỗi an toàn.
// Hàm này đảm bảo rằng server luôn trả về response cho client, kể cả khi có panic xảy ra.
//
// Parameters:
// - c: Fiber context
// - handler: Function xử lý chính của handler
func (h *BaseHandler[T, CreateInput, UpdateInput]) SafeHandler(c fiber.Ctx, handler func() error) error {
	defer func() {
		if r := recover(); r != nil {
			// Log stack trace để debug
			debug.PrintStack()

			// Trả về lỗi cho client
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeInternalServer,
				fmt.Sprintf("Lỗi hệ thống không mong muốn: %v", r),
				common.StatusInternalServerError,
				nil,
			))
		}
	}()
	return handler()
}

// HandleResponse xử lý và chuẩn hóa response trả về cho client.
// Phương thức này đảm bảo format response thống nhất trong toàn bộ ứng dụng.
//
// Parameters:
// - c: Fiber context
// - data: Dữ liệu trả về cho client (có thể là nil nếu chỉ trả về lỗi)
// - err: Lỗi nếu có (nil nếu không có lỗi)
func (h *BaseHandler[T, CreateInput, UpdateInput]) HandleResponse(c fiber.Ctx, data interface{}, err error) {
	if err != nil {
		var customErr *common.Error
		if errors.As(err, &customErr) {
			c.Status(customErr.StatusCode).JSON(fiber.Map{
				"code":    customErr.Code.Code,
				"message": customErr.Message,
				"details": customErr.Details,
				"status":  "error",
			})
			return
		}
		// Nếu không phải custom error, trả về internal server error
		c.Status(common.StatusInternalServerError).JSON(fiber.Map{
			"code":    common.ErrCodeDatabase.Code,
			"message": err.Error(),
			"status":  "error",
		})
		return
	}

	// Trường hợp thành công
	c.Status(common.StatusOK).JSON(fiber.Map{
		"code":    common.StatusOK,
		"message": common.MsgSuccess,
		"data":    data,
		"status":  "success",
	})
}
