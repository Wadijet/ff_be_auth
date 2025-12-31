// Package handler chứa các handler xử lý request HTTP cho phần xác thực và quản lý người dùng
package handler

import (
	"context"
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/logger"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
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

// HandleGetUserRoles lấy danh sách tất cả các role của người dùng với thông tin organization
// @Summary Lấy danh sách role của người dùng
// @Description Trả về danh sách các role mà người dùng hiện có kèm thông tin organization.
// @Description QUAN TRỌNG: Context làm việc là ROLE, không phải organization.
// @Description CHỈ trả về các role trực tiếp của user, KHÔNG bao gồm children/parents organizations.
// @Description Đây là danh sách "context làm việc" - user sẽ chọn một ROLE trong danh sách này để làm việc.
// @Description Frontend sẽ gửi ROLE ID trong header X-Active-Role-ID, không phải organization ID.
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/roles [get]
func (h *UserHandler) HandleGetUserRoles(c fiber.Ctx) error {
	// Log để debug - kiểm tra handler có được gọi không
	logger.GetAppLogger().WithFields(logrus.Fields{
		"path":   c.Path(),
		"method": c.Method(),
	}).Error("🔵 [HANDLER] HandleGetUserRoles called - FORCE LOG")

	// Lấy user ID từ context
	userID := c.Locals("user_id")
	logger.GetAppLogger().WithFields(logrus.Fields{
		"path":        c.Path(),
		"user_id":     userID,
		"has_user_id": userID != nil,
	}).Error("🔵 [HANDLER] Checking user_id in context - FORCE LOG")

	if userID == nil {
		logger.GetAppLogger().WithFields(logrus.Fields{
			"path": c.Path(),
		}).Error("❌ [HANDLER] User not authenticated - returning 401 - FORCE LOG")
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeAuth, "User not authenticated", common.StatusUnauthorized, nil))
		return nil
	}

	// Chuyển đổi string ID thành ObjectID
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Invalid user ID", common.StatusBadRequest, err))
		return nil
	}

	// Lấy danh sách user role - CHỈ lấy các role trực tiếp của user
	// KHÔNG lấy children/parents organizations
	filter := bson.M{"userId": objID}
	userRoles, err := h.userRoleService.Find(context.Background(), filter, nil)
	if err != nil {
		logger.GetAppLogger().WithFields(logrus.Fields{
			"user_id": objID.Hex(),
			"error":   err.Error(),
		}).Error("❌ Failed to get user roles")
		h.HandleResponse(c, nil, err)
		return nil
	}
	
	logger.GetAppLogger().WithFields(logrus.Fields{
		"user_id":    objID.Hex(),
		"roles_count": len(userRoles),
	}).Info("📋 Found user roles")

	// Lấy thông tin chi tiết của từng role với organization
	// Mỗi role tương ứng với một organization - đây là "context làm việc"
	result := make([]map[string]interface{}, 0, len(userRoles))
	for _, userRole := range userRoles {
		// Lấy role
		role, err := h.roleService.FindOneById(context.Background(), userRole.RoleID)
		if err != nil {
			logger.GetAppLogger().WithFields(logrus.Fields{
				"role_id": userRole.RoleID.Hex(),
				"error":   err.Error(),
			}).Warn("⚠️ Failed to get role, skipping")
			continue
		}

		// Validate OwnerOrganizationID không được zero
		if role.OwnerOrganizationID.IsZero() {
			logger.GetAppLogger().WithFields(logrus.Fields{
				"role_id": role.ID.Hex(),
				"role_name": role.Name,
			}).Warn("⚠️ Role has zero OwnerOrganizationID, skipping")
			continue
		}

		// Lấy organization - CHỈ lấy organization trực tiếp của role (logic business)
		// KHÔNG lấy children/parents organizations
		organizationService, err := services.NewOrganizationService()
		if err != nil {
			logger.GetAppLogger().WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("⚠️ Failed to create organization service, skipping")
			continue
		}
		// Dùng OwnerOrganizationID trực tiếp (đã bỏ OrganizationID)
		orgID := role.OwnerOrganizationID
		org, err := organizationService.FindOneById(context.Background(), orgID)
		if err != nil {
			logger.GetAppLogger().WithFields(logrus.Fields{
				"role_id": role.ID.Hex(),
				"organization_id": orgID.Hex(),
				"error": err.Error(),
			}).Warn("⚠️ Failed to get organization, skipping")
			continue
		}

		// Trả về thông tin role và organization trực tiếp
		// Frontend sẽ dùng danh sách này để user chọn "context làm việc"
		// QUAN TRỌNG: Context làm việc là ROLE, không phải organization
		// Mỗi role = một context làm việc
		// Organization được tự động suy ra từ role khi user chọn role
		result = append(result, map[string]interface{}{
			"roleId":             role.ID.Hex(),
			"roleName":           role.Name,
			"ownerOrganizationId": org.ID.Hex(), // Nhất quán với model Role (OwnerOrganizationID)
			"organizationName":   org.Name,
			"organizationCode":   org.Code,
			"organizationType":   org.Type,
			"organizationLevel":  org.Level,
		})
	}

	logger.GetAppLogger().WithFields(logrus.Fields{
		"user_id":      objID.Hex(),
		"result_count": len(result),
		"user_roles_count": len(userRoles),
	}).Info("✅ Returning roles with organizations")

	h.HandleResponse(c, result, nil)
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
