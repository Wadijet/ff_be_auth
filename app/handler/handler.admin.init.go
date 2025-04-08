package handler

import (
	"fmt"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/registry"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
)

// InitHandler xử lý các route liên quan đến khởi tạo hệ thống
// Kế thừa từ BaseHandler để có các chức năng CRUD cơ bản
type InitHandler struct {
	*BaseHandler[interface{}, interface{}, interface{}]
	userCRUD       services.BaseServiceMongo[models.User]
	permissionCRUD services.BaseServiceMongo[models.Permission]
	roleCRUD       services.BaseServiceMongo[models.Role]
	initService    *services.InitService
}

// NewInitHandler tạo một instance mới của InitHandler
// Returns:
//   - *InitHandler: Instance mới của InitHandler
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewInitHandler() (*InitHandler, error) {
	handler := &InitHandler{}

	// Khởi tạo base handler
	baseHandler := &BaseHandler[interface{}, interface{}, interface{}]{}
	handler.BaseHandler = baseHandler

	// Lấy các collection từ registry
	userCol, err := registry.Collections.MustGet(global.MongoDB_ColNames.Users)
	if err != nil {
		return nil, fmt.Errorf("failed to get users collection: %v", err)
	}

	permissionCol, err := registry.Collections.MustGet(global.MongoDB_ColNames.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions collection: %v", err)
	}

	roleCol, err := registry.Collections.MustGet(global.MongoDB_ColNames.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles collection: %v", err)
	}

	// Khởi tạo các service
	handler.userCRUD = services.NewBaseServiceMongo[models.User](userCol)
	handler.permissionCRUD = services.NewBaseServiceMongo[models.Permission](permissionCol)
	handler.roleCRUD = services.NewBaseServiceMongo[models.Role](roleCol)

	// Khởi tạo InitService
	initService, err := services.NewInitService()
	if err != nil {
		return nil, fmt.Errorf("failed to create init service: %v", err)
	}
	handler.initService = initService

	return handler, nil
}

// HandleSetAdministrator xử lý tạo người dùng quản trị hệ thống
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Thiết lập admin thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "email": "...",
//     "name": "...",
//     "isAdmin": true,
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 400: ID không hợp lệ
//   - 404: Không tìm thấy người dùng
//   - 500: Lỗi server
func (h *InitHandler) HandleSetAdministrator(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	if id == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	result, err := h.initService.SetAdministrator(utility.String2ObjectID(id))
	h.HandleResponse(c, result, err)
	return nil
}
