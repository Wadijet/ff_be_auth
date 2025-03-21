package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
type FbMessageHandler struct {
	FbMessageService services.FbMessageService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewFbMessageHandler(c *config.Configuration, db *mongo.Client) *FbMessageHandler {
	newHandler := new(FbMessageHandler)
	newHandler.FbMessageService = *services.NewFbMessageService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới FbMessage
func (h *FbMessageHandler) Create(ctx *fasthttp.RequestCtx) {
	utility.GenericHandler[models.FbMessageCreateInput](ctx, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		inputStruct := input.(*models.FbMessageCreateInput)
		return h.FbMessageService.ReviceData(ctx, inputStruct)
	})
}

// Tìm một FbMessage theo ID
func (h *FbMessageHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.FbMessageService.FindOneById(ctx, id))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}

//

// Tìm tất cả các FbMessage với phân trang
func (h *FbMessageHandler) FindAll(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	buf := string(ctx.FormValue("limit"))
	limit, err := strconv.ParseInt(buf, 10, 64)
	if err != nil {
		limit = 10
	}

	buf = string(ctx.FormValue("page"))
	page, err := strconv.ParseInt(buf, 10, 64)
	if err != nil {
		page = 0
	}

	response = utility.FinalResponse(h.FbMessageService.FindAll(ctx, page, limit))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}
