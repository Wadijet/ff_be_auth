package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
type FbConversationHandler struct {
	FbConversationService services.FbConversationService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewFbConversationHandler(c *config.Configuration, db *mongo.Client) *FbConversationHandler {
	newHandler := new(FbConversationHandler)
	newHandler.FbConversationService = *services.NewFbConversationService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo mới một FbConversation
func (h *FbConversationHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.FbConversationCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.FbConversationService.ReviceData(ctx, inputStruct))
			ctx.SetStatusCode(fasthttp.StatusOK)
		} else {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}

	utility.JSON(ctx, response)
}

// Tìm một FbConversation theo ID
func (h *FbConversationHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.FbConversationService.FindOneById(ctx, id))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}

// Tìm tất cả các FbConversation với phân trang
func (h *FbConversationHandler) FindAll(ctx *fasthttp.RequestCtx) {

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

	filter := bson.M{}

	pageId := string(ctx.FormValue("pageId"))
	if pageId != "" {
		filter = bson.M{"pageId": pageId}
	}

	response = utility.FinalResponse(h.FbConversationService.FindAll(ctx, page, limit, filter))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}

// Tìm tất cả các FbConversation với phân trang sắp xếp theo thời gian cập nhật của dữ liệu API
func (h *FbConversationHandler) FindAllSortByApiUpdate(ctx *fasthttp.RequestCtx) {

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

	filter := bson.M{}

	pageId := string(ctx.FormValue("pageId"))
	if pageId != "" {
		filter = bson.M{"pageId": pageId}
	}

	response = utility.FinalResponse(h.FbConversationService.FindAllSortByApiUpdate(ctx, page, limit, filter))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}
