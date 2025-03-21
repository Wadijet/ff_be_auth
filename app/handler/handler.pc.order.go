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
type PcOrderHandler struct {
	PcOrderService services.PcOrderService
}

// NewRoleHandler khởi tạo một RoleHandler mới
func NewPcOrderHandler(c *config.Configuration, db *mongo.Client) *PcOrderHandler {
	newHandler := new(PcOrderHandler)
	newHandler.PcOrderService = *services.NewPcOrderService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới PcOrder
func (h *PcOrderHandler) Create(ctx *fasthttp.RequestCtx) {
	utility.GenericHandler[models.PcOrderCreateInput](ctx, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		inputStruct := input.(*models.PcOrderCreateInput)
		return h.PcOrderService.ReviceData(ctx, inputStruct)
	})
}

// FindOneById tìm PcOrder theo ID
func (h *PcOrderHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	result, err := h.PcOrderService.FindOneById(ctx, id)
	response := utility.FinalResponse(result, err)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
	utility.JSON(ctx, response)
}

// Tìm tất cả các PcOrder với phân trang
func (h *PcOrderHandler) FindAll(ctx *fasthttp.RequestCtx) {

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

	response = utility.FinalResponse(h.PcOrderService.FindAll(ctx, page, limit))
	if response != nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	utility.JSON(ctx, response)
}
