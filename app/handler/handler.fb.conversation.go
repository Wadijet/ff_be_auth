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
)

// FbConversationHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Conversation
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbConversationHandler struct {
	BaseHandler
	FbConversationService *services.FbConversationService
}

// NewFbConversationHandler khởi tạo một FbConversationHandler mới
func NewFbConversationHandler(c *config.Configuration, db *mongo.Client) *FbConversationHandler {
	newHandler := new(FbConversationHandler)
	newHandler.FbConversationService = services.NewFbConversationService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create tạo mới một FbConversation
func (h *FbConversationHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.FbConversationCreateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	data, err := h.FbConversationService.ReviceData(context, input)
	h.HandleResponse(ctx, data, err)
}

// FindOne tìm một FbConversation theo ID
func (h *FbConversationHandler) FindOne(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	data, err := h.FbConversationService.FindOne(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các FbConversation với phân trang
func (h *FbConversationHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{}

	pageId := string(ctx.FormValue("pageId"))
	if pageId != "" {
		filter = bson.M{"pageId": pageId}
	}

	data, err := h.FbConversationService.FindAllWithPaginate(context, filter, page, limit)
	h.HandleResponse(ctx, data, err)
}

// FindAllSortByApiUpdate tìm tất cả các FbConversation với phân trang sắp xếp theo thời gian cập nhật của dữ liệu API
func (h *FbConversationHandler) FindAllSortByApiUpdate(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	context := context.Background()
	filter := bson.M{}

	pageId := string(ctx.FormValue("pageId"))
	if pageId != "" {
		filter = bson.M{"pageId": pageId}
	}

	data, err := h.FbConversationService.FindAllSortByApiUpdate(context, page, limit, filter)
	h.HandleResponse(ctx, data, err)
}

// Update cập nhật một FbConversation
func (h *FbConversationHandler) Update(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.HandleError(ctx, err)
		return
	}

	input := new(models.FbConversationCreateInput)
	if response := h.ParseRequestBody(ctx, input); response != nil {
		h.HandleError(ctx, nil)
		return
	}

	context := context.Background()
	conversation := models.FbConversation{
		ID:               objectID,
		PageId:           input.PageId,
		PageUsername:     input.PageUsername,
		PanCakeData:      input.PanCakeData,
		PanCakeUpdatedAt: time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
	}

	data, err := h.FbConversationService.Update(context, objectID, conversation)
	h.HandleResponse(ctx, data, err)
}

// Delete xóa một FbConversation
func (h *FbConversationHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	context := context.Background()
	err := h.FbConversationService.Delete(context, utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}
