package handler

import (
	"meta_commerce/app/database/registry"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
)

// FiberInitHandler xử lý các route liên quan đến khởi tạo hệ thống cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
type FiberInitHandler struct {
	FiberBaseHandler[interface{}, interface{}, interface{}]
	UserCRUD       services.BaseServiceMongo[models.User]
	PermissionCRUD services.BaseServiceMongo[models.Permission]
	RoleCRUD       services.BaseServiceMongo[models.Role]
	InitService    services.InitService
}

// NewFiberInitHandler tạo một instance mới của FiberInitHandler
// Returns:
//   - *FiberInitHandler: Instance mới của FiberInitHandler đã được khởi tạo với các service cần thiết
func NewFiberInitHandler() *FiberInitHandler {
	handler := &FiberInitHandler{}

	// Khởi tạo các collection từ registry
	userCol := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Users)
	permissionCol := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Permissions)
	roleCol := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Roles)

	// Khởi tạo các service với BaseService
	handler.UserCRUD = services.NewBaseServiceMongo[models.User](userCol)
	handler.PermissionCRUD = services.NewBaseServiceMongo[models.Permission](permissionCol)
	handler.RoleCRUD = services.NewBaseServiceMongo[models.Role](roleCol)
	handler.InitService = *services.NewInitService()
	return handler
}

// HandleSetAdministrator xử lý tạo người dùng quản trị hệ thống
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Path Params:
//   - id: ID của người dùng cần set làm admin
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
func (h *FiberInitHandler) HandleSetAdministrator(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": "ID không hợp lệ",
		})
	}

	result, err := h.InitService.SetAdministrator(utility.String2ObjectID(id))
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
