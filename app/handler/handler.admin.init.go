package handler

import (
	"fmt"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
)

// InitHandler xử lý các route liên quan đến khởi tạo hệ thống
// Kế thừa từ BaseHandler để có các chức năng CRUD cơ bản
type InitHandler struct {
	*BaseHandler[interface{}, interface{}, interface{}]
	userCRUD       *services.UserService
	permissionCRUD *services.PermissionService
	roleCRUD       *services.RoleService
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

	// Khởi tạo các service
	userService, err := services.NewUserService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %v", err)
	}
	handler.userCRUD = userService

	permissionService, err := services.NewPermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create permission service: %v", err)
	}
	handler.permissionCRUD = permissionService

	roleService, err := services.NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}
	handler.roleCRUD = roleService

	// Khởi tạo InitService
	initService, err := services.NewInitService()
	if err != nil {
		return nil, fmt.Errorf("failed to create init service: %v", err)
	}
	handler.initService = initService

	// Gán UserCRUD cho BaseHandler
	handler.BaseService = nil

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
