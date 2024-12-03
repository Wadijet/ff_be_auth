package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
type OrganizationHandler struct {
	crud services.Repository
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewOrganizationHandler(c *config.Configuration, db *mongo.Client) *OrganizationHandler {
	newHandler := new(OrganizationHandler)
	newHandler.crud = *services.NewRepository(c, db, global.MongoDB_ColNames.Organizations)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo mới một tổ chức
func (h *OrganizationHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.OrganizationCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.crud.InsertOne(ctx, inputStruct))
		}
	}

	utility.JSON(ctx, response)
}

// Tìm một tổ chức theo ID
func (h *OrganizationHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crud.FindOneById(ctx, id, nil))

	utility.JSON(ctx, response)
}

// Tìm tất cả các tổ chức với phân trang
func (h *OrganizationHandler) FindAll(ctx *fasthttp.RequestCtx) {
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

	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)

	response = utility.FinalResponse(h.crud.FindAll(ctx, nil, opts))

	utility.JSON(ctx, response)
}

// Cập nhật một tổ chức theo ID
func (h *OrganizationHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.OrganizationUpdateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.crud.UpdateOneById(ctx, id, inputStruct))
		}
	}

	utility.JSON(ctx, response)
}

// Xóa một tổ chức theo ID
func (h *OrganizationHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crud.DeleteOneById(ctx, id))

	utility.JSON(ctx, response)
}
