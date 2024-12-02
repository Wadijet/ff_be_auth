package handler

import (
	"atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PermissionHandler là struct chứa các phương thức xử lý quyền
type PermissionHandler struct {
	crud services.Repository
}

// NewPermissionHandler khởi tạo một PermissionHandler mới
func NewPermissionHandler(c *config.Configuration, db *mongo.Client) *PermissionHandler {
	newHandler := new(PermissionHandler)
	newHandler.crud = *services.NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	return newHandler
}

// CRUD functions =========================================================================

// Create tạo mới một quyền
func (h *PermissionHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ request
	postValues := ctx.PostBody()
	inputStruct := new(models.PermissionCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.crud.InsertOne(ctx, inputStruct))
		}
	}

	utility.JSON(ctx, response)
}

// FindOneById tìm kiếm một quyền theo ID
func (h *PermissionHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ request
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crud.FindOneById(ctx, id, nil))

	utility.JSON(ctx, response)
}

// FindAll tìm kiếm tất cả các quyền với phân trang
func (h *PermissionHandler) FindAll(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu phân trang từ request
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
	opts.SetSort(bson.D{{"updatedAt", 1}})

	response = utility.FinalResponse(h.crud.FindAllWithPaginate(ctx, bson.D{}, opts))

	utility.JSON(ctx, response)
}

// UpdateOneById cập nhật một quyền theo ID
func (h *PermissionHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ request
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crud.FindOneById(ctx, id, nil))

	// Lấy dữ liệu từ request
	postValues := ctx.PostBody()
	inputStruct := new(models.PermissionUpdateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm tạo json changes
			var change map[string]interface{}
			response = utility.CreateChangeMap(inputStruct, &change)
			if response == nil { // Gọi hàm xử lý logic
				response = utility.FinalResponse(h.crud.UpdateOneById(ctx, id, change))
			}
		}
	}

	utility.JSON(ctx, response)
}

// DeleteOneById xóa một quyền theo ID
func (h *PermissionHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ request
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crud.DeleteOneById(ctx, id))

	utility.JSON(ctx, response)
}

// Other functions =========================================================================
