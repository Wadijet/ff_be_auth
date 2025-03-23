package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FbConversationHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Conversation
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbConversationHandler struct {
	BaseHandler
	FbConversationService services.FbConversationService
}

// NewFbConversationHandler khởi tạo một FbConversationHandler mới
func NewFbConversationHandler(c *config.Configuration, db *mongo.Client) *FbConversationHandler {
	newHandler := new(FbConversationHandler)
	newHandler.FbConversationService = *services.NewFbConversationService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create tạo mới một FbConversation
func (h *FbConversationHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.FbConversationCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		return h.FbConversationService.ReviceData(ctx, input.(*models.FbConversationCreateInput))
	})
}

// FindOneById tìm một FbConversation theo ID
func (h *FbConversationHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.FbConversationService.FindOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các FbConversation với phân trang
func (h *FbConversationHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	filter := bson.M{}

	pageId := string(ctx.FormValue("pageId"))
	if pageId != "" {
		filter = bson.M{"pageId": pageId}
	}

	data, err := h.FbConversationService.FindAll(ctx, page, limit, filter)
	h.HandleResponse(ctx, data, err)
}

// FindAllSortByApiUpdate tìm tất cả các FbConversation với phân trang sắp xếp theo thời gian cập nhật của dữ liệu API
func (h *FbConversationHandler) FindAllSortByApiUpdate(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	filter := bson.M{}

	pageId := string(ctx.FormValue("pageId"))
	if pageId != "" {
		filter = bson.M{"pageId": pageId}
	}

	data, err := h.FbConversationService.FindAllSortByApiUpdate(ctx, page, limit, filter)
	h.HandleResponse(ctx, data, err)
}
