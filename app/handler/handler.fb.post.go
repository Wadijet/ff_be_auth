package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"
	"context"
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
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.FbPostService.ReviceData(context, input)
	h.HandleResponse(ctx, data, err)
}

// FindOne tìm một FbPost theo ID
func (h *FbPostHandler) FindOne(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.FbPostService.FindOne(context, id)
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
	input := new(models.FbPostCreateInput)
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

	fbPost := models.FbPost{
		ID:          objectID,
		PanCakeData: input.PanCakeData,
		UpdatedAt:   time.Now().Unix(),
	}
	data, err := h.FbPostService.Update(context, id, fbPost)
	h.HandleResponse(ctx, data, err)
}

// Delete xóa một FbPost
func (h *FbPostHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.FbPostService.Delete(context, id)
	h.HandleResponse(ctx, nil, err)
}

// UpdateToken cập nhật access token của một FbPost
func (h *FbPostHandler) UpdateToken(ctx *fasthttp.RequestCtx) {
	input := new(models.FbPostUpdateTokenInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.FbPostService.UpdateToken(context, input)
	h.HandleResponse(ctx, data, err)
}
