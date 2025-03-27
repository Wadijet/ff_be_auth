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

// FbPageHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Page
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbPageHandler struct {
	BaseHandler
	FbPageService *services.FbPageService
}

// NewFbPageHandler khởi tạo một FbPageHandler mới
func NewFbPageHandler(c *config.Configuration, db *mongo.Client) *FbPageHandler {
	newHandler := new(FbPageHandler)
	newHandler.FbPageService = services.NewFbPageService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create tạo mới một FbPage
func (h *FbPageHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.FbPageCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		fbPageInput := input.(*models.FbPageCreateInput)
		return h.FbPageService.ReviceData(context.Background(), fbPageInput)
	})
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.FbPageService.ReviceData(context, input)
	h.HandleResponse(ctx, data, err)
}

// FindOne tìm một FbPage theo ID
func (h *FbPageHandler) FindOne(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.FbPageService.FindOneById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, data, err)
}

// FindOneByPageID tìm một FbPage theo PageID
func (h *FbPageHandler) FindOneByPageID(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.FbPageService.FindOneByPageID(context, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các FbPage với phân trang
func (h *FbPageHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()

	data, err := h.FbPageService.FindAll(context, page, limit)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một FbPage
func (h *FbPageHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	input := new(models.FbPageCreateInput)
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

	fbPage := models.FbPage{
		ID:          objectID,
		PanCakeData: input.PanCakeData,
		UpdatedAt:   time.Now().Unix(),
	}
	data, err := h.FbPageService.UpdateById(context, utility.String2ObjectID(id), fbPage)
	h.HandleResponse(ctx, data, err)
}

// Delete xóa một FbPage
func (h *FbPageHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.FbPageService.DeleteById(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}

// UpdateToken cập nhật access token của một FbPage
func (h *FbPageHandler) UpdateToken(ctx *fasthttp.RequestCtx) {
	input := new(models.FbPageUpdateTokenInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.FbPageService.UpdateToken(context, input)
	h.HandleResponse(ctx, data, err)
}
