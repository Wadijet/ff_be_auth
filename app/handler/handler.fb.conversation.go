package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FbConversationHandler là cấu trúc xử lý các yêu cầu liên quan đến Facebook Conversation
// Kế thừa từ BaseHandler để sử dụng các phương thức xử lý chung
type FbConversationHandler struct {
	BaseHandler[models.FbConversation, models.FbConversationCreateInput, models.FbConversationCreateInput]
	FbConversationService *services.FbConversationService
}

// NewFbConversationHandler khởi tạo một FbConversationHandler mới
func NewFbConversationHandler(c *config.Configuration, db *mongo.Client) *FbConversationHandler {
	newHandler := new(FbConversationHandler)
	newHandler.FbConversationService = services.NewFbConversationService(c, db)
	// Không cần gán service cho BaseHandler vì chúng ta sẽ sử dụng FbConversationService trực tiếp
	return newHandler
}

// Các hàm đặc thù của FbConversation (nếu có) sẽ được thêm vào đây

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
