package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

// StaticHandler là cấu trúc xử lý các yêu cầu liên quan đến thông tin tĩnh cho Fiber
type StaticHandler struct {
	BaseHandler[interface{}, interface{}, interface{}]
}

// NewStaticHandler khởi tạo một FiberStaticHandler mới
// Returns:
//   - *FiberStaticHandler: Instance mới của FiberStaticHandler
func NewStaticHandler() *StaticHandler {
	handler := new(StaticHandler)
	return handler
}

// ==========================================================================================

// HandleTestApi kiểm tra API có hoạt động không
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: API hoạt động bình thường
func (h *StaticHandler) HandleTestApi(c fiber.Ctx) error {
	data := map[string]interface{}{
		"status": "ok",
		"time":   time.Now().Unix(),
	}
	h.HandleResponse(c, data, nil)
	return nil
}
