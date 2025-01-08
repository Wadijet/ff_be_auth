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
type AccessTokenHandler struct {
	AccessTokenService services.AccessTokenService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewAccessTokenHandler(c *config.Configuration, db *mongo.Client) *AccessTokenHandler {
	newHandler := new(AccessTokenHandler)
	newHandler.AccessTokenService = *services.NewAccessTokenService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo mới một Access Token
func (h *AccessTokenHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.AccessTokenCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.AccessTokenService.Create(ctx, inputStruct))
			ctx.SetStatusCode(fasthttp.StatusCreated)
		} else {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}

	utility.JSON(ctx, response)
}

// Tìm một Access Token theo ID
func (h *AccessTokenHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.AccessTokenService.FindOneById(ctx, id))

	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}

// Tìm tất cả các Access Token với phân trang
func (h *AccessTokenHandler) FindAll(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
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

	// Gọi hàm xử lý logic
	response = utility.FinalResponse(h.AccessTokenService.FindAll(ctx, page, limit))

	ctx.SetStatusCode(fasthttp.StatusOK)
	utility.JSON(ctx, response)
}

// Cập nhật một Access Token theo ID
func (h *AccessTokenHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.AccessTokenUpdateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.AccessTokenService.UpdateOneById(ctx, id, inputStruct))
			ctx.SetStatusCode(fasthttp.StatusOK)
		} else {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}

	utility.JSON(ctx, response)
}

// Xóa một Access Token theo ID
func (h *AccessTokenHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.AccessTokenService.DeleteOneById(ctx, id))

	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}
