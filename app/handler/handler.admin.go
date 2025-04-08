package handler

import (
	"fmt"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/registry"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminHandler xử lý các route liên quan đến quản trị viên cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
type AdminHandler struct {
	BaseHandler[models.User, models.UserCreateInput, models.UserChangeInfoInput]
	UserCRUD       services.BaseServiceMongo[models.User]
	PermissionCRUD services.BaseServiceMongo[models.Permission]
	RoleCRUD       services.BaseServiceMongo[models.Role]
	AdminService   services.AdminService
}

// NewAdminHandler tạo một instance mới của FiberAdminHandler
// Returns:
//   - *FiberAdminHandler: Instance mới của FiberAdminHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewAdminHandler() (*AdminHandler, error) {
	handler := &AdminHandler{}

	// Khởi tạo các collection từ registry
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

	// Khởi tạo các service với BaseService
	handler.UserCRUD = services.NewBaseServiceMongo[models.User](userCol)
	handler.PermissionCRUD = services.NewBaseServiceMongo[models.Permission](permissionCol)
	handler.RoleCRUD = services.NewBaseServiceMongo[models.Role](roleCol)

	// Khởi tạo AdminService và xử lý error
	adminService, err := services.NewAdminService()
	if err != nil {
		return nil, fmt.Errorf("failed to create admin service: %v", err)
	}
	handler.AdminService = *adminService

	// Gán UserCRUD cho BaseHandler
	handler.Service = handler.UserCRUD
	return handler, nil
}

// SetRoleInput là cấu trúc dữ liệu đầu vào cho việc thiết lập vai trò người dùng
type SetRoleInput struct {
	Email  string             `json:"email" validate:"required"`
	RoleID primitive.ObjectID `json:"roleID" validate:"required"`
}

// HandleSetRole xử lý thiết lập vai trò cho người dùng
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - email: Email của người dùng cần set role
//   - roleID: ID của role cần gán
//
// Response:
//   - 200: Thiết lập role thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "email": "...",
//     "name": "...",
//     "roles": ["..."],
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 400: Dữ liệu đầu vào không hợp lệ
//   - 404: Không tìm thấy người dùng hoặc role
//   - 500: Lỗi server
func (h *AdminHandler) HandleSetRole(c fiber.Ctx) error {
	var input SetRoleInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, err.Error(), utility.StatusBadRequest, nil))
		return nil
	}

	result, err := h.AdminService.SetRole(c.Context(), input.Email, input.RoleID)
	h.HandleResponse(c, result, err)
	return nil
}

// HandleBlockUser xử lý khóa người dùng
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - email: Email của người dùng cần khóa
//   - note: Ghi chú lý do khóa
//
// Response:
//   - 200: Khóa người dùng thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "email": "...",
//     "name": "...",
//     "isBlock": true,
//     "blockNote": "...",
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 400: Dữ liệu đầu vào không hợp lệ
//   - 404: Không tìm thấy người dùng
//   - 500: Lỗi server
func (h *AdminHandler) HandleBlockUser(c fiber.Ctx) error {
	var input models.BlockUserInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, err.Error(), utility.StatusBadRequest, nil))
		return nil
	}

	result, err := h.AdminService.BlockUser(c.Context(), input.Email, true, input.Note)
	h.HandleResponse(c, result, err)
	return nil
}

// HandleUnBlockUser xử lý mở khóa người dùng
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - email: Email của người dùng cần mở khóa
//
// Response:
//   - 200: Mở khóa người dùng thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "email": "...",
//     "name": "...",
//     "isBlock": false,
//     "blockNote": "",
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 400: Dữ liệu đầu vào không hợp lệ
//   - 404: Không tìm thấy người dùng
//   - 500: Lỗi server
func (h *AdminHandler) HandleUnBlockUser(c fiber.Ctx) error {
	var input models.UnBlockUserInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, err.Error(), utility.StatusBadRequest, nil))
		return nil
	}

	result, err := h.AdminService.BlockUser(c.Context(), input.Email, false, "")
	h.HandleResponse(c, result, err)
	return nil
}
