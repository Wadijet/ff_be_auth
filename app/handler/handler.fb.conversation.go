package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// FiberFbConversationHandler xử lý các route liên quan đến Facebook Conversation cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
type FiberFbConversationHandler struct {
	FiberBaseHandler[models.FbConversation, models.FbConversationCreateInput, models.FbConversationCreateInput]
	FbConversationService *services.FbConversationService
}

// NewFiberFbConversationHandler tạo một instance mới của FiberFbConversationHandler
// Returns:
//   - *FiberFbConversationHandler: Instance mới của FiberFbConversationHandler đã được khởi tạo với các service cần thiết
func NewFiberFbConversationHandler() *FiberFbConversationHandler {
	handler := &FiberFbConversationHandler{}
	handler.FbConversationService = services.NewFbConversationService()
	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng FbConversationService trực tiếp
	return handler
}

// HandleFindAllSortByApiUpdate tìm tất cả các FbConversation với phân trang sắp xếp theo thời gian cập nhật của dữ liệu API
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Query Params:
//   - page: Trang hiện tại (mặc định: 1)
//   - limit: Số lượng item trên một trang (mặc định: 10)
//   - pageId: ID của page Facebook cần lọc (tùy chọn)
//
// Response:
//   - 200: Lấy dữ liệu thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "page": 1,
//     "limit": 10,
//     "itemCount": 5,
//     "items": [{
//     "id": "...",
//     "pageId": "...",
//     "messages": [...],
//     "updatedAt": 123,
//     "createdAt": 123
//     }]
//     }
//     }
//   - 400: Tham số không hợp lệ
//   - 500: Lỗi server
func (h *FiberFbConversationHandler) HandleFindAllSortByApiUpdate(c fiber.Ctx) error {
	// Parse page và limit từ query params
	pageInt, limitInt := h.ParsePagination(c)
	page := int64(pageInt)
	limit := int64(limitInt)

	// Tạo filter dựa trên pageId
	filter := bson.M{}
	if pageId := c.Query("pageId"); pageId != "" {
		filter = bson.M{"pageId": pageId}
	}

	// Gọi service để lấy dữ liệu
	result, err := h.FbConversationService.FindAllSortByApiUpdate(context.Background(), page, limit, filter)
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
		"data":    result,
	})
}
