package handler

import (
	"context"
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"

	"github.com/gofiber/fiber/v3"
)

// FbPostHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Post cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.FbPost
// - CreateInput: dto.FbPostCreateInput
// - UpdateInput: dto.FbPostCreateInput
type FbPostHandler struct {
	BaseHandler[models.FbPost, dto.FbPostCreateInput, dto.FbPostCreateInput]
	FbPostService *services.FbPostService
}

// NewFbPostHandler khởi tạo một FiberFbPostHandler mới
// Returns:
//   - *FiberFbPostHandler: Instance mới của FiberFbPostHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewFbPostHandler() (*FbPostHandler, error) {
	handler := &FbPostHandler{}

	// Khởi tạo FbPostService và xử lý error
	service, err := services.NewFbPostService()
	if err != nil {
		return nil, fmt.Errorf("failed to create post service: %v", err)
	}
	handler.FbPostService = service
	handler.BaseService = handler.FbPostService

	return handler, nil
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
func (h *FbPostHandler) HandleFindOneByPostID(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	data, err := h.FbPostService.FindOneByPostID(context.Background(), id)
	h.HandleResponse(c, data, err)
	return nil
}
