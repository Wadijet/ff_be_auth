package handler

import (
	"fmt"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/utility"

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
// CHỈ DÙNG KHI CHƯA CÓ ADMIN - Không check quyền
// @Summary Thiết lập administrator (chỉ khi chưa có admin)
// @Description Gán quyền Administrator cho một người dùng. Chỉ hoạt động khi hệ thống chưa có admin nào.
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse "Hệ thống đã có admin, vui lòng sử dụng endpoint /admin/user/administrator/:id"
// @Router /init/set-administrator/:id [post]
func (h *InitHandler) HandleSetAdministrator(c fiber.Ctx) error {
	// Kiểm tra xem đã có admin chưa
	hasAdmin, err := h.initService.HasAnyAdministrator()
	if err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeInternalServer, "Không thể kiểm tra trạng thái admin", common.StatusInternalServerError, err))
		return nil
	}

	if hasAdmin {
		h.HandleResponse(c, nil, common.NewError(
			common.ErrCodeBusinessState,
			"Hệ thống đã có admin. Vui lòng sử dụng endpoint /admin/user/set-administrator/:id với quyền Init.SetAdmin.",
			common.StatusForbidden,
			nil,
		))
		return nil
	}

	id := h.GetIDFromContext(c)
	if id == "" {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "ID không hợp lệ", common.StatusBadRequest, nil))
		return nil
	}

	result, err := h.initService.SetAdministrator(utility.String2ObjectID(id))
	h.HandleResponse(c, result, err)
	return nil
}

// HandleInitOrganization khởi tạo Organization Root
// @Summary Khởi tạo Organization Root
// @Description Tạo Organization Root (Group - Level 0) nếu chưa tồn tại
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /init/organization [post]
func (h *InitHandler) HandleInitOrganization(c fiber.Ctx) error {
	err := h.initService.InitRootOrganization()
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}
	h.HandleResponse(c, map[string]string{"message": "Organization Root đã được khởi tạo thành công"}, nil)
	return nil
}

// HandleInitPermissions khởi tạo Permissions
// @Summary Khởi tạo Permissions
// @Description Tạo tất cả các quyền mặc định của hệ thống
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /init/permissions [post]
func (h *InitHandler) HandleInitPermissions(c fiber.Ctx) error {
	err := h.initService.InitPermission()
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}
	h.HandleResponse(c, map[string]string{"message": "Permissions đã được khởi tạo thành công"}, nil)
	return nil
}

// HandleInitRoles khởi tạo Roles
// @Summary Khởi tạo Roles
// @Description Tạo Role Administrator và gán tất cả quyền
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /init/roles [post]
func (h *InitHandler) HandleInitRoles(c fiber.Ctx) error {
	err := h.initService.InitRole()
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}
	// Đảm bảo Administrator có đầy đủ quyền
	err = h.initService.CheckPermissionForAdministrator()
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}
	h.HandleResponse(c, map[string]string{"message": "Roles đã được khởi tạo thành công"}, nil)
	return nil
}

// HandleInitAdminUser khởi tạo Admin User từ Firebase UID
// @Summary Khởi tạo Admin User
// @Description Tạo user admin từ Firebase UID (user phải đã tồn tại trong Firebase)
// @Accept json
// @Produce json
// @Param input body dto.InitAdminUserInput true "Firebase UID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /init/admin-user [post]
func (h *InitHandler) HandleInitAdminUser(c fiber.Ctx) error {
	type InitAdminUserInput struct {
		FirebaseUID string `json:"firebaseUid" validate:"required"`
	}

	var input InitAdminUserInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	err := h.initService.InitAdminUser(input.FirebaseUID)
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}
	h.HandleResponse(c, map[string]string{"message": "Admin user đã được khởi tạo thành công"}, nil)
	return nil
}

// HandleInitAll khởi tạo tất cả các đơn vị cơ bản
// @Summary Khởi tạo tất cả
// @Description Khởi tạo Organization, Permissions, Roles (one-click setup)
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /init/all [post]
func (h *InitHandler) HandleInitAll(c fiber.Ctx) error {
	results := make(map[string]interface{})

	// 1. Init Organization Root
	if err := h.initService.InitRootOrganization(); err != nil {
		results["organization"] = map[string]string{"status": "failed", "error": err.Error()}
	} else {
		results["organization"] = map[string]string{"status": "success"}
	}

	// 2. Init Permissions
	if err := h.initService.InitPermission(); err != nil {
		results["permissions"] = map[string]string{"status": "failed", "error": err.Error()}
	} else {
		results["permissions"] = map[string]string{"status": "success"}
	}

	// 3. Init Roles
	if err := h.initService.InitRole(); err != nil {
		results["roles"] = map[string]string{"status": "failed", "error": err.Error()}
	} else {
		results["roles"] = map[string]string{"status": "success"}
		// Đảm bảo Administrator có đầy đủ quyền
		_ = h.initService.CheckPermissionForAdministrator()
	}

	h.HandleResponse(c, results, nil)
	return nil
}

// HandleInitStatus kiểm tra trạng thái khởi tạo hệ thống
// @Summary Kiểm tra trạng thái init
// @Description Kiểm tra các đơn vị cơ bản đã được khởi tạo chưa
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Router /init/status [get]
func (h *InitHandler) HandleInitStatus(c fiber.Ctx) error {
	status, err := h.initService.GetInitStatus()
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}
	h.HandleResponse(c, status, nil)
	return nil
}
