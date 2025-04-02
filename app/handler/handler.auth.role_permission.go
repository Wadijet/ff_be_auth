package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RolePermissionHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò và quyền
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type RolePermissionHandler struct {
	BaseHandler[models.RolePermission, models.RolePermissionCreateInput, models.RolePermissionUpdateInput]
	RolePermissionService *services.RolePermissionService
}

// NewRolePermissionHandler khởi tạo một RolePermissionHandler mới
func NewRolePermissionHandler() *RolePermissionHandler {
	newHandler := new(RolePermissionHandler)
	newHandler.RolePermissionService = services.NewRolePermissionService()
	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng RolePermissionService trực tiếp
	return newHandler
}

// Các hàm đặc thù của RolePermission (nếu có) sẽ được thêm vào đây

// Update cập nhật một RolePermission
func (h *RolePermissionHandler) Update(ctx *fasthttp.RequestCtx) {
	input := new(models.RolePermissionUpdateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	roleId, err := primitive.ObjectIDFromHex(input.RoleId)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	// Xóa các role permission cũ
	filter := bson.M{"roleId": roleId}
	_, err = h.RolePermissionService.DeleteMany(context, filter)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	// Tạo các role permission mới
	var rolePermissions []models.RolePermission
	for _, permissionId := range input.PermissionIds {
		permissionIdObj, err := primitive.ObjectIDFromHex(permissionId)
		if err != nil {
			continue
		}
		rolePermission := models.RolePermission{
			ID:           primitive.NewObjectID(),
			RoleID:       roleId,
			PermissionID: permissionIdObj,
			Scope:        0,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		rolePermissions = append(rolePermissions, rolePermission)
	}

	// Lưu các role permission mới
	for _, rolePermission := range rolePermissions {
		_, err = h.RolePermissionService.InsertOne(context, rolePermission)
		if err != nil {
			h.HandleError(ctx, err)
			return
		}
	}

	h.HandleResponse(ctx, rolePermissions, nil)
}
