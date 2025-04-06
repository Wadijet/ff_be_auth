package handler

import (
	"meta_commerce/app/database/registry"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FiberAdminHandler xử lý các route liên quan đến quản trị viên cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
type FiberAdminHandler struct {
	FiberBaseHandler[models.User, models.UserCreateInput, models.UserChangeInfoInput]
	UserCRUD       services.BaseServiceMongo[models.User]
	PermissionCRUD services.BaseServiceMongo[models.Permission]
	RoleCRUD       services.BaseServiceMongo[models.Role]
	AdminService   services.AdminService
}

// NewFiberAdminHandler tạo một instance mới của FiberAdminHandler
// Returns:
//   - *FiberAdminHandler: Instance mới của FiberAdminHandler đã được khởi tạo với các service cần thiết
func NewFiberAdminHandler() *FiberAdminHandler {
	handler := &FiberAdminHandler{}

	// Khởi tạo các collection từ registry
	userCol := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Users)
	permissionCol := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Permissions)
	roleCol := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Roles)

	// Khởi tạo các service với BaseService
	handler.UserCRUD = services.NewBaseServiceMongo[models.User](userCol)
	handler.PermissionCRUD = services.NewBaseServiceMongo[models.Permission](permissionCol)
	handler.RoleCRUD = services.NewBaseServiceMongo[models.Role](roleCol)
	handler.AdminService = *services.NewAdminService()

	// Gán UserCRUD cho BaseHandler
	handler.Service = handler.UserCRUD
	return handler
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
func (h *FiberAdminHandler) HandleSetRole(c fiber.Ctx) error {
	var input SetRoleInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": err.Error(),
		})
	}

	result, err := h.AdminService.SetRole(c.Context(), input.Email, input.RoleID)
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
func (h *FiberAdminHandler) HandleBlockUser(c fiber.Ctx) error {
	var input models.BlockUserInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": err.Error(),
		})
	}

	result, err := h.AdminService.BlockUser(c.Context(), input.Email, true, input.Note)
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
func (h *FiberAdminHandler) HandleUnBlockUser(c fiber.Ctx) error {
	var input models.UnBlockUserInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": err.Error(),
		})
	}

	result, err := h.AdminService.BlockUser(c.Context(), input.Email, false, "")
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
