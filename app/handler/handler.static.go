package handler

import (
	"meta_commerce/app/utility"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// FiberStaticHandler là cấu trúc xử lý các yêu cầu liên quan đến thông tin tĩnh cho Fiber
type FiberStaticHandler struct {
	FiberBaseHandler[interface{}, interface{}, interface{}]
}

// NewFiberStaticHandler khởi tạo một FiberStaticHandler mới
// Returns:
//   - *FiberStaticHandler: Instance mới của FiberStaticHandler
func NewFiberStaticHandler() *FiberStaticHandler {
	return new(FiberStaticHandler)
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
func (h *FiberStaticHandler) HandleTestApi(c fiber.Ctx) error {
	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
	})
}

// FiberSystemStaticResponse là cấu trúc chứa thông tin về tài nguyên hệ thống
type FiberSystemStaticResponse struct {
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
func (h *FiberStaticHandler) HandleGetSystemStatic(c fiber.Ctx) error {
	result := new(FiberSystemStaticResponse)
	result.Cpu = utility.GetCpuStatic()
	result.Memory = utility.GetMemoryStatic()

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    result,
	})
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
func (h *FiberStaticHandler) HandleGetApiStatic(c fiber.Ctx) error {
	inseconds := c.Query("inseconds", "30")
	insesonds, err := strconv.ParseInt(inseconds, 10, 64)
	if err != nil {
		insesonds = 30
	}

	data := utility.GetApiStatic(insesonds)
	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
}
