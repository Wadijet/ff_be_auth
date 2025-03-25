package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"context"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type RoleHandler struct {
	BaseHandler
	RoleService *services.RoleService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewRoleHandler(c *config.Configuration, db *mongo.Client) *RoleHandler {
	newHandler := new(RoleHandler)
	newHandler.RoleService = services.NewRoleService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới vai trò
func (h *RoleHandler) Create(ctx *fasthttp.RequestCtx) {
	inputStruct := new(models.RoleCreateInput)
	if response := h.ParseRequestBody(ctx, inputStruct); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	role := models.Role{
		Name:     inputStruct.Name,
		Describe: inputStruct.Describe,
	}
	data, err := h.RoleService.Create(context, role)
	h.HandleResponse(ctx, data, err)
}

// FindOneById tìm vai trò theo ID
func (h *RoleHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	context := context.Background()
	data, err := h.RoleService.FindOne(context, objectID)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các vai trò với phân trang
func (h *RoleHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{} // Có thể thêm filter từ query params nếu cần
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit)
	data, err := h.RoleService.FindAll(context, filter, opts)
	h.HandleResponse(ctx, data, err)
}

// UpdateOneById cập nhật một vai trò theo ID
func (h *RoleHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {

	// Lấy ID từ context
	id := h.GetIDFromContext(ctx)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	input := new(models.RoleUpdateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	role := models.Role{
		Name:     input.Name,
		Describe: input.Describe,
	}
	data, err := h.RoleService.Update(context, objectID, role)
	h.HandleResponse(ctx, data, err)
}

// DeleteOneById xóa một vai trò theo ID
func (h *RoleHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.RoleService.Delete(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}

// Other functions =========================================================================
