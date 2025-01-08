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
type AgentHandler struct {
	AgentService services.AgentService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewAgentHandler(c *config.Configuration, db *mongo.Client) *AgentHandler {
	newHandler := new(AgentHandler)
	newHandler.AgentService = *services.NewAgentService(c, db)
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
			response = utility.FinalResponse(h.AgentService.Create(ctx, inputStruct))
		}
	}

	utility.JSON(ctx, response)
}

// Tìm một Agent theo ID
func (h *AgentHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.AgentService.FindOneById(ctx, id))

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

	response = utility.FinalResponse(h.AgentService.FindAll(ctx, page, limit))

	utility.JSON(ctx, response)
}

// Cập nhật một Agent theo ID
func (h *AgentHandler) UpdateOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)

	// Lấy dữ liệu từ yêu cầu
	postValues := ctx.PostBody()
	inputStruct := new(models.AgentUpdateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.AgentService.Update(ctx, id, inputStruct))
		}
	}

	utility.JSON(ctx, response)
}

// Xóa một Agent theo ID
func (h *AgentHandler) DeleteOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.AgentService.Delete(ctx, id))

	utility.JSON(ctx, response)
}
