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
type FbPageHandler struct {
	FbPageHandlerService services.FbPageService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewFbPageHandler(c *config.Configuration, db *mongo.Client) *FbPageHandler {
	newHandler := new(FbPageHandler)
	newHandler.FbPageHandlerService = *services.NewFbPageService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo mới một FbPage
func (h *FbPageHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.FbPageCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.FbPageHandlerService.ReviceData(ctx, inputStruct))
			ctx.SetStatusCode(fasthttp.StatusOK)
		} else {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}

	utility.JSON(ctx, response)
}

// Tìm một FbPage theo ID
func (h *FbPageHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.FbPageHandlerService.FindOneById(ctx, id))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}

// Tìm một FbPage theo PageID
func (h *FbPageHandler) FindOneByPageID(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.FbPageHandlerService.FindOneByPageID(ctx, id))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)

}

// Tìm tất cả các FbPage với phân trang
func (h *FbPageHandler) FindAll(ctx *fasthttp.RequestCtx) {
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

	response = utility.FinalResponse(h.FbPageHandlerService.FindAll(ctx, page, limit))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}

// Cập nhật access token của một FbPage
func (h *FbPageHandler) UpdateToken(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.FbPageUpdateTokenInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.FbPageHandlerService.UpdateToken(ctx, inputStruct))
			ctx.SetStatusCode(fasthttp.StatusOK)
		} else {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}

	utility.JSON(ctx, response)
}
