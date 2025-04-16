// Package handler chứa các handler xử lý request HTTP cho phần xác thực và quản lý người dùng
package handler

import (
	"context"
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/utility"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler xử lý các request liên quan đến xác thực và quản lý thông tin người dùng
type UserHandler struct {
	*BaseHandler[models.User, models.UserCreateInput, models.UserChangeInfoInput]
	userService *services.UserService
	roleService *services.RoleService
}

// NewUserHandler tạo một instance mới của UserHandler
func NewUserHandler() (*UserHandler, error) {
	handler := &UserHandler{}

	// Khởi tạo base handler
	baseHandler := &BaseHandler[models.User, models.UserCreateInput, models.UserChangeInfoInput]{}
	handler.BaseHandler = baseHandler

	// Khởi tạo các service
	userService, err := services.NewUserService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %v", err)
	}
	handler.userService = userService

	roleService, err := services.NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}
	handler.roleService = roleService

	return handler, nil
}

// HandleLogin xử lý đăng nhập người dùng
// @Summary Đăng nhập người dùng
// @Description Xác thực thông tin đăng nhập và trả về token nếu thành công
// @Accept json
// @Produce json
// @Param input body models.UserLoginInput true "Thông tin đăng nhập"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) HandleLogin(c fiber.Ctx) error {
	var input models.UserLoginInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	user, err := h.userService.Login(context.Background(), &input)
	h.HandleResponse(c, user, err)
	return nil
}

// HandleRegister xử lý đăng ký tài khoản mới
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - email: Email đăng ký
//   - password: Mật khẩu
//   - name: Tên người dùng
//
// Response:
//   - 200: Đăng ký thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "email": "...",
//     "name": "...",
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 400: Dữ liệu không hợp lệ
//   - 409: Email đã tồn tại
//   - 500: Lỗi server
func (h *UserHandler) HandleRegister(c fiber.Ctx) error {
	var input models.UserCreateInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	user, err := h.userService.Registry(context.Background(), &input)
	h.HandleResponse(c, user, err)
	return nil
}

// HandleLogout xử lý đăng xuất người dùng
func (h *UserHandler) HandleLogout(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeAuth, "User not authenticated", utility.StatusUnauthorized, nil))
		return nil
	}

	var input models.UserLogoutInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Invalid user ID", utility.StatusBadRequest, err))
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
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeAuth, "User not authenticated", utility.StatusUnauthorized, nil))
		return nil
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Invalid user ID", utility.StatusBadRequest, err))
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
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeAuth, "User not authenticated", utility.StatusUnauthorized, nil))
		return nil
	}

	var input models.UserChangeInfoInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Invalid user ID", utility.StatusBadRequest, err))
		return nil
	}

	// Tạo user object từ input
	updateUser := &models.User{
		Name: input.Name,
	}

	updatedUser, err := h.userService.UpdateById(context.Background(), objID, *updateUser)
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

// HandleChangePassword xử lý thay đổi mật khẩu người dùng
func (h *UserHandler) HandleChangePassword(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeAuth, "User not authenticated", utility.StatusUnauthorized, nil))
		return nil
	}

	var input models.UserChangePasswordInput
	if err := h.ParseRequestBody(c, &input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Invalid user ID", utility.StatusBadRequest, err))
		return nil
	}

	err = h.userService.ChangePassword(context.Background(), objID, &input)
	h.HandleResponse(c, nil, err)
	return nil
}
