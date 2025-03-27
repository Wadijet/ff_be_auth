package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
	"meta_commerce/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò của người dùng
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type UserRoleHandler struct {
	BaseHandler
	UserRoleService *services.UserRoleService
}

// NewUserRoleHandler khởi tạo một UserRoleHandler mới
func NewUserRoleHandler(c *config.Configuration, db *mongo.Client) *UserRoleHandler {
	newHandler := new(UserRoleHandler)
	newHandler.UserRoleService = services.NewUserRoleService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới UserRole
func (h *UserRoleHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.UserRoleCreateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.UserRoleService.Create(context, input)
	h.HandleResponse(ctx, data, err)
}

// FindOne tìm một UserRole theo ID
func (h *UserRoleHandler) FindOne(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.UserRoleService.FindOneById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các UserRole với phân trang
func (h *UserRoleHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{} // Có thể thêm filter từ query params nếu cần
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit)
	data, err := h.UserRoleService.Find(context, filter, opts)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một UserRole
func (h *UserRoleHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	input := new(models.UserRoleCreateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	userRole := models.UserRole{
		ID:     objectID,
		UserID: input.UserID,
		RoleID: input.RoleID,
	}
	data, err := h.UserRoleService.UpdateById(context, objectID, userRole)
	h.HandleResponse(ctx, data, err)
}

// Delete xóa một UserRole
func (h *UserRoleHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.UserRoleService.DeleteById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}
