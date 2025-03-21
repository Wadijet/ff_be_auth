package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRoleHandler là cấu trúc xử lý các yêu cầu liên quan đến vai trò
type UserRoleHandler struct {
	UserRoleService *services.UserRoleService
}

// NewUserRoleHandler khởi tạo một UserRoleHandler mới
func NewUserRoleHandler(c *config.Configuration, db *mongo.Client) *UserRoleHandler {
	newHandler := new(UserRoleHandler)
	newHandler.UserRoleService = services.NewUserRoleService(c, db)
	return newHandler
}

// CRUD functions ==========================================================================

// Create xử lý tạo mới UserRole
func (h *UserRoleHandler) Create(ctx *fasthttp.RequestCtx) {
	utility.GenericHandler[models.UserRoleCreateInput](ctx, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		inputStruct := input.(*models.UserRoleCreateInput)
		return h.UserRoleService.Create(ctx, inputStruct)
	})
}

// Xóa một UserRole
func (h *UserRoleHandler) Delete(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy ID từ yêu cầu
	id := ctx.UserValue("id").(string)
	response = utility.FinalResponse(h.UserRoleService.Delete(ctx, id))
	if response["error"] != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // Set status code to 400 Bad Request if there's an error
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK) // Set status code to 200 OK
	}
	utility.JSON(ctx, response)
}
