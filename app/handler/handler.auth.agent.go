package handler

import (
	models "atk-go-server/app/models/mongodb"
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

// RoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
type AgentHandler struct {
	crud services.RepositoryService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewAgentHandler(c *config.Configuration, db *mongo.Client) *AgentHandler {
	newHandler := new(AgentHandler)
	newHandler.crud = *services.NewRepository(c, db, global.MongoDB_ColNames.Agents)
	return newHandler
}

// CRUD functions ==========================================================================

// Tạo mới một Agent
func (h *AgentHandler) Create(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.AgentCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.crud.InsertOne(ctx, inputStruct))
		}
	}

	utility.JSON(ctx, response)
}

// Tìm một Agent theo ID
func (h *AgentHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crud.FindOneById(ctx, utility.String2ObjectID(id), nil))

	utility.JSON(ctx, response)
}

// Tìm tất cả các Agent với phân trang
func (h *AgentHandler) FindAll(ctx *fasthttp.RequestCtx) {
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
	opts.SetSort(bson.D{{"updatedAt", 1}})

	response = utility.FinalResponse(h.crud.FindAllWithPaginate(ctx, bson.D{}, opts))

	utility.JSON(ctx, response)
}

// Cập nhật một Agent theo ID
func (h *AgentHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
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
			response = utility.FinalResponse(h.crud.UpdateOneById(ctx, utility.String2ObjectID(id), inputStruct))
		}
	}

	utility.JSON(ctx, response)
}

// Xóa một Agent theo ID
func (h *AgentHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.crud.DeleteOneById(ctx, utility.String2ObjectID(id)))

	utility.JSON(ctx, response)
}
