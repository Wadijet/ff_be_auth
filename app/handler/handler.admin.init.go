package handler

import (
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/registry"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
)

// InitHandler xử lý các route liên quan đến khởi tạo hệ thống cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
type InitHandler struct {
	BaseHandler[interface{}, interface{}, interface{}]
	UserCRUD       services.BaseServiceMongo[models.User]
	PermissionCRUD services.BaseServiceMongo[models.Permission]
	RoleCRUD       services.BaseServiceMongo[models.Role]
	InitService    services.InitService
}

// NewInitHandler tạo một instance mới của FiberInitHandler
// Returns:
//   - *FiberInitHandler: Instance mới của FiberInitHandler đã được khởi tạo với các service cần thiết
func NewInitHandler() *InitHandler {
	handler := &InitHandler{}

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
func (h *InitHandler) HandleSetAdministrator(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	if id == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	result, err := h.InitService.SetAdministrator(utility.String2ObjectID(id))
	h.HandleResponse(c, result, err)
	return nil
}
