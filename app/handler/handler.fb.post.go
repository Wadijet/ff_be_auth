package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
	"meta_commerce/config"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// FbPostHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Post
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbPostHandler struct {
	BaseHandler
	FbPostService *services.FbPostService
}

// NewFbPostHandler khởi tạo một FbPostHandler mới
func NewFbPostHandler(c *config.Configuration, db *mongo.Client) *FbPostHandler {
	newHandler := new(FbPostHandler)
	newHandler.FbPostService = services.NewFbPostService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create tạo mới một FbPost
func (h *FbPostHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.FbPostCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		fbPostInput := input.(*models.FbPostCreateInput)
		return h.FbPostService.ReviceData(context.Background(), fbPostInput)
	})
}

// FindOne tìm một FbPost theo ID
func (h *FbPostHandler) FindOne(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	context := context.Background()
	data, err := h.FbPostService.FindOneById(context, objectID)
	h.HandleResponse(ctx, data, err)
}

// FindOneByPostID tìm một FbPost theo PostID
func (h *FbPostHandler) FindOneByPostID(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.FbPostService.FindOneByPostID(context, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các FbPost với phân trang
func (h *FbPostHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()

	data, err := h.FbPostService.FindAll(context, page, limit)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một FbPost
func (h *FbPostHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	input := new(models.FbPostCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		fbPostInput := input.(*models.FbPostCreateInput)
		post := models.FbPost{
			ID:          objectID,
			PanCakeData: fbPostInput.PanCakeData,
			UpdatedAt:   time.Now().Unix(),
		}
		return h.FbPostService.UpdateById(context.Background(), objectID, post)
	})
}

// Delete xóa một FbPost
func (h *FbPostHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.FbPostService.DeleteById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}

// UpdateToken cập nhật access token của một FbPost
func (h *FbPostHandler) UpdateToken(ctx *fasthttp.RequestCtx) {
	input := new(models.FbPostUpdateTokenInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		fbPostInput := input.(*models.FbPostUpdateTokenInput)
		return h.FbPostService.UpdateToken(context.Background(), fbPostInput)
	})
}
