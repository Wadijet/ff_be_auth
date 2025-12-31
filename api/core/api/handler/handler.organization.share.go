package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/utility"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationShareHandler xử lý các request liên quan đến Organization Share
type OrganizationShareHandler struct {
	BaseHandler[models.OrganizationShare, dto.OrganizationShareCreateInput, dto.OrganizationShareUpdateInput]
	OrganizationShareService *services.OrganizationShareService
}

// NewOrganizationShareHandler tạo mới OrganizationShareHandler
func NewOrganizationShareHandler() (*OrganizationShareHandler, error) {
	shareService, err := services.NewOrganizationShareService()
	if err != nil {
		return nil, fmt.Errorf("failed to create organization share service: %v", err)
	}

	baseHandler := NewBaseHandler[models.OrganizationShare, dto.OrganizationShareCreateInput, dto.OrganizationShareUpdateInput](shareService)
	handler := &OrganizationShareHandler{
		BaseHandler:              *baseHandler,
		OrganizationShareService: shareService,
	}

	return handler, nil
}

// CreateShare tạo sharing giữa 2 organizations
// POST /api/v1/organization-shares
func (h *OrganizationShareHandler) CreateShare(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		var input dto.OrganizationShareCreateInput
		if err := h.ParseRequestBody(c, &input); err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Dữ liệu gửi lên không đúng định dạng: %v", err),
				common.StatusBadRequest,
				err,
			))
			return nil
		}

		// Validate ObjectIDs
		ownerOrgID, err := primitive.ObjectIDFromHex(input.OwnerOrganizationID)
		if err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("ownerOrganizationId không hợp lệ: %v", err),
				common.StatusBadRequest,
				err,
			))
			return nil
		}

		toOrgID, err := primitive.ObjectIDFromHex(input.ToOrgID)
		if err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("toOrgId không hợp lệ: %v", err),
				common.StatusBadRequest,
				err,
			))
			return nil
		}

		// Validate: ownerOrgID và toOrgID không được giống nhau
		if ownerOrgID == toOrgID {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationInput,
				"ownerOrganizationId và toOrgId không được giống nhau",
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		// Validate: user có quyền share data của ownerOrg
		userIDStr, ok := c.Locals("user_id").(string)
		if !ok {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeAuth,
				"Không tìm thấy user ID",
				common.StatusUnauthorized,
				nil,
			))
			return nil
		}
		userID, _ := primitive.ObjectIDFromHex(userIDStr)

		// Kiểm tra user có quyền truy cập ownerOrg không
		allowedOrgIDs, err := services.GetUserAllowedOrganizationIDs(c.Context(), userID, "")
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		hasAccess := false
		for _, orgID := range allowedOrgIDs {
			if orgID == ownerOrgID {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeAuth,
				"Bạn không có quyền share data của organization này",
				common.StatusForbidden,
				nil,
			))
			return nil
		}

		// Kiểm tra share đã tồn tại chưa
		_, err = h.OrganizationShareService.FindOne(c.Context(), bson.M{
			"ownerOrganizationId": ownerOrgID,
			"toOrgId":             toOrgID,
		}, nil)

		if err == nil {
			// Share đã tồn tại
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeBusinessOperation,
				"Share đã tồn tại giữa 2 organizations này",
				common.StatusConflict,
				nil,
			))
			return nil
		}
		// Nếu err != nil và không phải ErrNotFound, có thể là lỗi khác
		if err != common.ErrNotFound {
			h.HandleResponse(c, nil, err)
			return nil
		}
		// err == ErrNotFound, tiếp tục tạo mới

		// Tạo share record
		share := models.OrganizationShare{
			OwnerOrganizationID: ownerOrgID,
			ToOrgID:             toOrgID,
			PermissionNames:     input.PermissionNames,
			CreatedAt:           utility.CurrentTimeInMilli(),
			CreatedBy:           userID,
		}

		data, err := h.BaseService.InsertOne(c.Context(), share)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// DeleteShare xóa sharing
// DELETE /api/v1/organization-shares/:id
func (h *OrganizationShareHandler) DeleteShare(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		id := c.Params("id")
		if !primitive.IsValidObjectID(id) {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("ID không hợp lệ: %s", id),
				common.StatusBadRequest,
				nil,
			))
			return nil
		}
		shareID := utility.String2ObjectID(id)

		// Lấy share để kiểm tra quyền
		share, err := h.OrganizationShareService.FindOneById(c.Context(), shareID)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		// Validate: user có quyền xóa share này (phải là người tạo hoặc có quyền với fromOrg)
		userIDStr, ok := c.Locals("user_id").(string)
		if !ok {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeAuth,
				"Không tìm thấy user ID",
				common.StatusUnauthorized,
				nil,
			))
			return nil
		}
		userID, _ := primitive.ObjectIDFromHex(userIDStr)

		// Kiểm tra user có phải người tạo không
		if share.CreatedBy != userID {
			// Kiểm tra user có quyền với ownerOrg không
			allowedOrgIDs, err := services.GetUserAllowedOrganizationIDs(c.Context(), userID, "")
			if err != nil {
				h.HandleResponse(c, nil, err)
				return nil
			}

			hasAccess := false
			for _, orgID := range allowedOrgIDs {
				if orgID == share.OwnerOrganizationID {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				h.HandleResponse(c, nil, common.NewError(
					common.ErrCodeAuth,
					"Bạn không có quyền xóa share này",
					common.StatusForbidden,
					nil,
				))
				return nil
			}
		}

		// Xóa share
		err = h.BaseService.DeleteById(c.Context(), shareID)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		h.HandleResponse(c, map[string]interface{}{
			"message": "Xóa share thành công",
		}, nil)
		return nil
	})
}

// ListShares liệt kê các shares của organization
// GET /api/v1/organization-shares?ownerOrganizationId=xxx hoặc ?toOrgId=xxx
func (h *OrganizationShareHandler) ListShares(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		ownerOrgIDStr := c.Query("ownerOrganizationId")
		toOrgIDStr := c.Query("toOrgId")

		filter := bson.M{}

		// Filter theo ownerOrganizationID
		if ownerOrgIDStr != "" {
			ownerOrgID, err := primitive.ObjectIDFromHex(ownerOrgIDStr)
			if err != nil {
				h.HandleResponse(c, nil, common.NewError(
					common.ErrCodeValidationFormat,
					fmt.Sprintf("ownerOrganizationId không hợp lệ: %v", err),
					common.StatusBadRequest,
					err,
				))
				return nil
			}

			// Validate: user có quyền xem shares của ownerOrg này
			userIDStr, ok := c.Locals("user_id").(string)
			if ok {
				userID, _ := primitive.ObjectIDFromHex(userIDStr)
				allowedOrgIDs, err := services.GetUserAllowedOrganizationIDs(c.Context(), userID, "")
				if err == nil {
					hasAccess := false
					for _, orgID := range allowedOrgIDs {
						if orgID == ownerOrgID {
							hasAccess = true
							break
						}
					}
					if !hasAccess {
						h.HandleResponse(c, nil, common.NewError(
							common.ErrCodeAuth,
							"Bạn không có quyền xem shares của organization này",
							common.StatusForbidden,
							nil,
						))
						return nil
					}
				}
			}

			filter["ownerOrganizationId"] = ownerOrgID
		}

		// Filter theo toOrgId
		if toOrgIDStr != "" {
			toOrgID, err := primitive.ObjectIDFromHex(toOrgIDStr)
			if err != nil {
				h.HandleResponse(c, nil, common.NewError(
					common.ErrCodeValidationFormat,
					fmt.Sprintf("toOrgId không hợp lệ: %v", err),
					common.StatusBadRequest,
					err,
				))
				return nil
			}
			filter["toOrgId"] = toOrgID
		}

		// Nếu không có filter nào, trả về lỗi
		if len(filter) == 0 {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationInput,
				"Cần cung cấp ít nhất một trong các tham số: ownerOrganizationId hoặc toOrgId",
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		// Query shares
		data, err := h.BaseService.Find(c.Context(), filter, nil)
		h.HandleResponse(c, data, err)
		return nil
	})
}
