// Package handler chứa các handler xử lý request HTTP cho phần xác thực và quản lý người dùng
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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler xử lý các request liên quan đến xác thực và quản lý thông tin người dùng
type UserHandler struct {
	*BaseHandler[models.User, dto.UserCreateInput, dto.UserChangeInfoInput]
	userService     *services.UserService
	roleService     *services.RoleService
	userRoleService *services.UserRoleService
}

// NewUserHandler tạo một instance mới của UserHandler
func NewUserHandler() (*UserHandler, error) {
	// Khởi tạo các service
	userService, err := services.NewUserService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %v", err)
	}

	roleService, err := services.NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}

	userRoleService, err := services.NewUserRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user role service: %v", err)
	}

	baseHandler := NewBaseHandler[models.User, dto.UserCreateInput, dto.UserChangeInfoInput](userService)
	handler := &UserHandler{
		BaseHandler:     baseHandler,
		userService:     userService,
		roleService:     roleService,
		userRoleService: userRoleService,
	}

	return handler, nil
}

// HandleLogout xử lý đăng xuất người dùng
func (h *UserHandler) HandleLogout(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeAuth, "User not authenticated", common.StatusUnauthorized, nil))
		return nil
	}

	var input dto.UserLogoutInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Invalid user ID", common.StatusBadRequest, err))
		return nil
	}

	err = h.userService.Logout(context.Background(), objID, &input)
	h.HandleResponse(c, nil, err)
	return nil
}

// --------------------------------
// User Profile Methods
// --------------------------------

// HandleGetProfile lấy thông tin profile của người dùng
func (h *UserHandler) HandleGetProfile(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeAuth, "User not authenticated", common.StatusUnauthorized, nil))
		return nil
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Invalid user ID", common.StatusBadRequest, err))
		return nil
	}

	user, err := h.userService.FindOneById(context.Background(), objID)
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	// Loại bỏ thông tin nhạy cảm
	user.Password = ""
	user.Salt = ""
	user.Tokens = nil

	h.HandleResponse(c, user, nil)
	return nil
}

// HandleUpdateProfile cập nhật thông tin profile của người dùng
func (h *UserHandler) HandleUpdateProfile(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeAuth, "User not authenticated", common.StatusUnauthorized, nil))
		return nil
	}

	var input dto.UserChangeInfoInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Invalid user ID", common.StatusBadRequest, err))
		return nil
	}

	// Tạo update data với các trường cần update
	update := &services.UpdateData{
		Set: map[string]interface{}{
			"name": input.Name,
			// Thêm các trường khác nếu cần
		},
	}

	updatedUser, err := h.userService.UpdateById(context.Background(), objID, update)
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	// Loại bỏ thông tin nhạy cảm
	updatedUser.Password = ""
	updatedUser.Salt = ""
	updatedUser.Tokens = nil

	h.HandleResponse(c, updatedUser, nil)
	return nil
}

// HandleGetUserRoles lấy danh sách tất cả các role của người dùng
// @Summary Lấy danh sách role của người dùng
// @Description Trả về danh sách các role mà người dùng hiện có
// @Accept json
// @Produce json
// @Success 200 {array} models.Role
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/roles [get]
func (h *UserHandler) HandleGetUserRoles(c fiber.Ctx) error {
	// Lấy user ID từ context
	userID := c.Locals("user_id")
	if userID == nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeAuth, "User not authenticated", common.StatusUnauthorized, nil))
		return nil
	}

	// Chuyển đổi string ID thành ObjectID
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Invalid user ID", common.StatusBadRequest, err))
		return nil
	}

	// Lấy danh sách user role
	filter := bson.M{"userId": objID}
	userRoles, err := h.userRoleService.Find(context.Background(), filter, nil)
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	// Lấy thông tin chi tiết của từng role
	var roles []models.Role
	for _, userRole := range userRoles {
		role, err := h.roleService.FindOneById(context.Background(), userRole.RoleID)
		if err != nil {
			continue // Bỏ qua role không tìm thấy
		}
		roles = append(roles, role)
	}

	h.HandleResponse(c, roles, nil)
	return nil
}

// HandleLoginWithFirebase xử lý đăng nhập bằng Firebase ID token
// @Summary Đăng nhập bằng Firebase
// @Description Xác thực Firebase ID token và trả về JWT token nếu thành công
// @Accept json
// @Produce json
// @Param input body dto.FirebaseLoginInput true "Firebase ID token và hwid"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login/firebase [post]
func (h *UserHandler) HandleLoginWithFirebase(c fiber.Ctx) error {
	var input dto.FirebaseLoginInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	user, err := h.userService.LoginWithFirebase(context.Background(), &input)
	if err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	// Loại bỏ thông tin nhạy cảm trước khi trả về
	user.Password = ""
	user.Salt = ""
	user.Tokens = nil

	h.HandleResponse(c, user, nil)
	return nil
}