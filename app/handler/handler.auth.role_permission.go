package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"context"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RolePermissionHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò và quyền
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type RolePermissionHandler struct {
	BaseHandler
	RolePermissionService *services.RolePermissionService
}

// NewRolePermissionHandler khởi tạo một RolePermissionHandler mới
func NewRolePermissionHandler(c *config.Configuration, db *mongo.Client) *RolePermissionHandler {
	newHandler := new(RolePermissionHandler)
	newHandler.RolePermissionService = services.NewRolePermissionService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới RolePermission
func (h *RolePermissionHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.RolePermissionCreateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.RolePermissionService.Create(context, input)
	h.HandleResponse(ctx, data, err)
}

// FindOneById tìm một RolePermission theo ID
func (h *RolePermissionHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.RolePermissionService.FindOneById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các RolePermission với phân trang
func (h *RolePermissionHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{} // Có thể thêm filter từ query params nếu cần

	// Tạo options cho phân trang
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

	data, err := h.RolePermissionService.Find(context, filter, findOptions)
	h.HandleResponse(ctx, data, err)
}

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
		_, err = h.RolePermissionService.Create(context, &models.RolePermissionCreateInput{
			RoleID:       rolePermission.RoleID,
			PermissionID: rolePermission.PermissionID,
			Scope:        rolePermission.Scope,
		})
		if err != nil {
			h.HandleError(ctx, err)
			return
		}
	}

	h.HandleResponse(ctx, rolePermissions, nil)
}

// Delete xóa một RolePermission
func (h *RolePermissionHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.RolePermissionService.DeleteById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}
