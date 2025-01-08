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
type RoleHandler struct {
	RoleService *services.RoleService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewRoleHandler(c *config.Configuration, db *mongo.Client) *RoleHandler {
	newHandler := new(RoleHandler)
	newHandler.RoleService = services.NewRoleService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo mới một vai trò
func (h *RoleHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.RoleCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.RoleService.Create(ctx, inputStruct))
			ctx.SetStatusCode(fasthttp.StatusCreated)
		} else {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}

	utility.JSON(ctx, response)
}

// Tìm một vai trò theo ID
func (h *RoleHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.RoleService.FindOneById(ctx, id))
	if response["error"] != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	utility.JSON(ctx, response)
}

// Tìm tất cả các vai trò với phân trang
func (h *RoleHandler) FindAll(ctx *fasthttp.RequestCtx) {
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

	response = utility.FinalResponse(h.RoleService.FindAll(ctx, page, limit))
	ctx.SetStatusCode(fasthttp.StatusOK)

	utility.JSON(ctx, response)
}

// Cập nhật một vai trò theo ID
func (h *RoleHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.RoleUpdateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.RoleService.Update(ctx, id, inputStruct))
			ctx.SetStatusCode(fasthttp.StatusOK)
		} else {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}

	utility.JSON(ctx, response)
}

// Xóa một vai trò theo ID
func (h *RoleHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)

	response = utility.FinalResponse(h.RoleService.Delete(ctx, id))
	if response["error"] != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	utility.JSON(ctx, response)
}

// Other functions =========================================================================
