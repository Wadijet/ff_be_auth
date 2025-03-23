package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// FbPostHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Post
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbPostHandler struct {
	BaseHandler
	FbPostService services.FbPostService
}

// NewFbPostHandler khởi tạo một FbPostHandler mới
func NewFbPostHandler(c *config.Configuration, db *mongo.Client) *FbPostHandler {
	newHandler := new(FbPostHandler)
	newHandler.FbPostService = *services.NewFbPostService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create tạo mới một FbPost
func (h *FbPostHandler) Create(ctx *fasthttp.RequestCtx) {
	input := new(models.FbPostCreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		return h.FbPostService.ReviceData(ctx, input.(*models.FbPostCreateInput))
	})
}

// FindOneById tìm một FbPost theo ID
func (h *FbPostHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.FbPostService.FindOneById(ctx, id)
	h.HandleResponse(ctx, data, err)
}

// FindAll tìm tất cả các FbPost với phân trang
func (h *FbPostHandler) FindAll(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	data, err := h.FbPostService.FindAll(ctx, page, limit)
	h.HandleResponse(ctx, data, err)
}
