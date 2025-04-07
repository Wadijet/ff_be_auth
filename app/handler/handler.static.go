package handler

import (
	"meta_commerce/app/utility"
	"strconv"
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

// SystemStaticResponse là cấu trúc chứa thông tin về tài nguyên hệ thống
type SystemStaticResponse struct {
	Cpu    interface{} `json:"cpu" bson:"cpu"`       // Thông tin về CPU
	Memory interface{} `json:"memory" bson:"memory"` // Thông tin về bộ nhớ
}

// HandleGetSystemStatic lấy thông tin về tài nguyên hệ thống
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Lấy thông tin thành công
//   - 500: Lỗi server
func (h *StaticHandler) HandleGetSystemStatic(c fiber.Ctx) error {
	result := new(SystemStaticResponse)
	result.Cpu = utility.GetCpuStatic()
	result.Memory = utility.GetMemoryStatic()

	h.HandleResponse(c, result, nil)
	return nil
}

// HandleGetApiStatic lấy thông tin thống kê về API
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Query Parameters:
//   - inseconds: Số giây cần lấy thống kê (mặc định: 30)
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Lấy thông tin thành công
//   - 500: Lỗi server
func (h *StaticHandler) HandleGetApiStatic(c fiber.Ctx) error {
	inseconds := c.Query("inseconds", "30")
	insesonds, err := strconv.ParseInt(inseconds, 10, 64)
	if err != nil {
		insesonds = 30
	}

	data := utility.GetApiStatic(insesonds)
	h.HandleResponse(c, data, nil)
	return nil
}
