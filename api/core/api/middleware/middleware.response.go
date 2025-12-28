package middleware

import (
	"errors"
	"meta_commerce/core/common"

	"github.com/gofiber/fiber/v3"
)

// HandleErrorResponse xử lý và trả về error response cho client
// Tách riêng để tránh import cycle với handler package
func HandleErrorResponse(c fiber.Ctx, err error) {
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
}

