package handler

import (
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitHandler là struct chứa các CRUD services và InitService
type InitHandler struct {
	UserCRUD       services.RepositoryService
	PermissionCRUD services.RepositoryService
	RoleCRUD       services.RepositoryService
	InitService    services.InitService
}

// NewInitHandler khởi tạo InitHandler mới
func NewInitHandler(c *config.Configuration, db *mongo.Client) *InitHandler {
	newHandler := new(InitHandler)
	newHandler.UserCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Users)
	newHandler.PermissionCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newHandler.RoleCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Roles)
	newHandler.InitService = *services.NewInitService(c, db)
	return newHandler
}

// SetAdministrator tạo người dùng quản trị hệ thống
func (h *InitHandler) SetAdministrator(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	// GET ID
	id := ctx.UserValue("id").(string)

	response = utility.FinalResponse(h.InitService.SetAdministrator(id))

	utility.JSON(ctx, response)

}
