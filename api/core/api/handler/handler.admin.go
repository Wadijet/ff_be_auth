package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminHandler xử lý các route liên quan đến quản trị viên cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
type AdminHandler struct {
	BaseHandler[models.User, dto.UserCreateInput, dto.UserChangeInfoInput]
	UserCRUD       *services.UserService
	PermissionCRUD *services.PermissionService
	RoleCRUD       *services.RoleService
	AdminService   *services.AdminService
}

// NewAdminHandler tạo một instance mới của FiberAdminHandler
// Returns:
//   - *FiberAdminHandler: Instance mới của FiberAdminHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewAdminHandler() (*AdminHandler, error) {
	handler := &AdminHandler{}

	// Khởi tạo các service với BaseService
	userService, err := services.NewUserService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %v", err)
	}
	handler.UserCRUD = userService

	permissionService, err := services.NewPermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create permission service: %v", err)
	}
	handler.PermissionCRUD = permissionService

	roleService, err := services.NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}
	handler.RoleCRUD = roleService

	// Khởi tạo AdminService và xử lý error
	adminService, err := services.NewAdminService()
	if err != nil {
		return nil, fmt.Errorf("failed to create admin service: %v", err)
	}
	handler.AdminService = adminService

	// Gán UserCRUD cho BaseHandler
	handler.BaseService = nil
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
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, err.Error(), common.StatusBadRequest, nil))
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
	var input dto.BlockUserInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, err.Error(), common.StatusBadRequest, nil))
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
	var input dto.UnBlockUserInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, err.Error(), common.StatusBadRequest, nil))
		return nil
	}

	result, err := h.AdminService.BlockUser(c.Context(), input.Email, false, "")
	h.HandleResponse(c, result, err)
	return nil
}

// HandleAddAdministrator xử lý thêm administrator khi đã có admin
// YÊU CẦU QUYỀN Init.SetAdmin
// @Summary Thiết lập administrator (khi đã có admin)
// @Description Gán quyền Administrator cho một người dùng. Yêu cầu quyền Init.SetAdmin.
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse "Không có quyền"
// @Router /admin/user/set-administrator/:id [post]
func (h *AdminHandler) HandleAddAdministrator(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	if id == "" {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "ID không hợp lệ", common.StatusBadRequest, nil))
		return nil
	}

	// Parse ID từ string sang ObjectID
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "ID không hợp lệ", common.StatusBadRequest, err))
		return nil
	}

	// Sử dụng InitService để set administrator
	initService, err := services.NewInitService()
	if err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeInternalServer, "Không thể khởi tạo InitService", common.StatusInternalServerError, err))
		return nil
	}

	result, err := initService.SetAdministrator(userID)
	h.HandleResponse(c, result, err)
	return nil
}

// HandleSyncAdministratorPermissions đồng bộ quyền cho Administrator
// Endpoint này sẽ tạo các quyền mới (nếu chưa có) và đảm bảo Administrator có đầy đủ tất cả quyền
// @Summary Đồng bộ quyền cho Administrator
// @Description Tạo các quyền mới (nếu chưa có) và đảm bảo Administrator có đầy đủ tất cả quyền trong hệ thống
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/sync-administrator-permissions [post]
func (h *AdminHandler) HandleSyncAdministratorPermissions(c fiber.Ctx) error {
	// Sử dụng InitService để sync permissions
	initService, err := services.NewInitService()
	if err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeInternalServer, "Không thể khởi tạo InitService", common.StatusInternalServerError, err))
		return nil
	}

	// 1. Tạo các quyền mới (nếu chưa có) - InitPermission chỉ tạo các quyền chưa tồn tại
	if err := initService.InitPermission(); err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeInternalServer, "Không thể khởi tạo permissions", common.StatusInternalServerError, err))
		return nil
	}

	// 2. Đảm bảo Administrator có đầy đủ tất cả quyền
	if err := initService.CheckPermissionForAdministrator(); err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeInternalServer, "Không thể đồng bộ quyền cho Administrator", common.StatusInternalServerError, err))
		return nil
	}

	h.HandleResponse(c, map[string]string{"message": "Đã đồng bộ quyền cho Administrator thành công"}, nil)
	return nil
}
