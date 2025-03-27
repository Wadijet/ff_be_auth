package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"context"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PermissionHandler là struct chứa các phương thức xử lý quyền
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type PermissionHandler struct {
	BaseHandler
	PermissionService *services.PermissionService
}

// NewPermissionHandler khởi tạo một PermissionHandler mới
func NewPermissionHandler(c *config.Configuration, db *mongo.Client) *PermissionHandler {
	newHandler := new(PermissionHandler)
	newHandler.PermissionService = services.NewPermissionService(c, db)
	return newHandler
}

// CRUD functions =========================================================================

// Create tạo mới một quyền
func (h *PermissionHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.PermissionCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		permissionInput := input.(*models.PermissionCreateInput)
		permission := models.Permission{
			Name:     permissionInput.Name,
			Describe: permissionInput.Describe,
			Category: permissionInput.Category,
			Group:    permissionInput.Group,
		}
		return h.PermissionService.InsertOne(context.Background(), permission)
	})
}

// FindOneById tìm kiếm một quyền theo ID
func (h *PermissionHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.PermissionService.FindOneById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm kiếm tất cả các quyền với phân trang
func (h *PermissionHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{} // Có thể thêm filter từ query params nếu cần

	// Tạo options cho phân trang
	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

	data, err := h.PermissionService.Find(context, filter, findOptions)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một quyền
func (h *PermissionHandler) Update(ctx *fasthttp.RequestCtx) {
	input := new(models.PermissionUpdateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		permissionInput := input.(*models.PermissionUpdateInput)
		id := h.GetIDFromContext(ctx)
		permission := models.Permission{
			Name:     permissionInput.Name,
			Describe: permissionInput.Describe,
			Category: permissionInput.Category,
			Group:    permissionInput.Group,
		}
		return h.PermissionService.UpdateById(context.Background(), utility.String2ObjectID(id), permission)
	})
}

// Delete xóa một quyền
func (h *PermissionHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.PermissionService.DeleteById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}

// Other functions =========================================================================
