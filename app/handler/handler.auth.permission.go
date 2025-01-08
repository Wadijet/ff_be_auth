package handler

import (
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// PermissionHandler là struct chứa các phương thức xử lý quyền
type PermissionHandler struct {
	//crud services.RepositoryService
	PermissionService *services.PermissionService
}

// NewPermissionHandler khởi tạo một PermissionHandler mới
func NewPermissionHandler(c *config.Configuration, db *mongo.Client) *PermissionHandler {
	newHandler := new(PermissionHandler)
	newHandler.PermissionService = services.NewPermissionService(c, db)
	return newHandler
}

// CRUD functions =========================================================================

// FindOneById tìm kiếm một quyền theo ID
func (h *PermissionHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ request
	id := ctx.UserValue("id").(string)
	result, err := h.PermissionService.FindOneById(ctx, id)
	response = utility.FinalResponse(result, err)

	if result == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

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

	result, err := h.PermissionService.FindAll(ctx, page, limit)
	response = utility.FinalResponse(result, err)

	if result == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	utility.JSON(ctx, response)
}

// Other functions =========================================================================
