package handler

import (
	"context"
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// FbConversationHandler xử lý các route liên quan đến Facebook Conversation cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
type FbConversationHandler struct {
	BaseHandler[models.FbConversation, dto.FbConversationCreateInput, dto.FbConversationCreateInput]
	FbConversationService *services.FbConversationService
}

// NewFbConversationHandler tạo một instance mới của FiberFbConversationHandler
// Returns:
//   - *FiberFbConversationHandler: Instance mới của FiberFbConversationHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewFbConversationHandler() (*FbConversationHandler, error) {
	handler := &FbConversationHandler{}

	// Khởi tạo FbConversationService và xử lý error
	service, err := services.NewFbConversationService()
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation service: %v", err)
	}
	handler.FbConversationService = service

	// Gán BaseServiceMongoImpl cho BaseHandler để các method CRUD cơ bản hoạt động
	// FbConversationService có method Upsert với signature khác nên không implement BaseServiceMongo
	// Nhưng BaseServiceMongoImpl bên trong nó vẫn implement interface này
	handler.BaseService = service.BaseServiceMongoImpl

	return handler, nil
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
func (h *FbConversationHandler) HandleFindAllSortByApiUpdate(c fiber.Ctx) error {
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
		if customErr, ok := err.(*common.Error); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"code":    customErr.Code,
				"message": customErr.Message,
				"details": customErr.Details,
			})
		}
		return c.Status(common.StatusInternalServerError).JSON(fiber.Map{
			"code":    common.ErrCodeDatabase,
			"message": err.Error(),
		})
	}

	return c.Status(common.StatusOK).JSON(fiber.Map{
		"message": common.MsgSuccess,
		"data":    result,
	})
}
