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
type FbPostHandler struct {
	FbPostService services.FbPostService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewFbPostHandler(c *config.Configuration, db *mongo.Client) *FbPostHandler {
	newHandler := new(FbPostHandler)
	newHandler.FbPostService = *services.NewFbPostService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo mới một FbPost
func (h *FbPostHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.FbPostCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.FbPostService.ReviceData(ctx, inputStruct))
			ctx.SetStatusCode(fasthttp.StatusOK)
		} else {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}

	utility.JSON(ctx, response)
}

// Tìm một FbPost theo ID
func (h *FbPostHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.FbPostService.FindOneById(ctx, id))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}

// Tìm tất cả các FbPost với phân trang
func (h *FbPostHandler) FindAll(ctx *fasthttp.RequestCtx) {
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

	response = utility.FinalResponse(h.FbPostService.FindAll(ctx, page, limit))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}
