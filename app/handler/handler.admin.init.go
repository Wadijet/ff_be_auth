package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"net/http"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitHandler là struct chứa các CRUD services và InitService
type InitHandler struct {
	UserCRUD       services.BaseService[models.User]
	PermissionCRUD services.BaseService[models.Permission]
	RoleCRUD       services.BaseService[models.Role]
	InitService    services.InitService
}

// NewInitHandler khởi tạo InitHandler mới
func NewInitHandler(c *config.Configuration, db *mongo.Client) *InitHandler {
	newHandler := new(InitHandler)

	// Khởi tạo các collection
	userCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Users)).Collection(global.MongoDB_ColNames.Users)
	permissionCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Permissions)).Collection(global.MongoDB_ColNames.Permissions)
	roleCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Roles)).Collection(global.MongoDB_ColNames.Roles)

	// Khởi tạo các service với BaseService
	newHandler.UserCRUD = services.NewBaseService[models.User](userCol)
	newHandler.PermissionCRUD = services.NewBaseService[models.Permission](permissionCol)
	newHandler.RoleCRUD = services.NewBaseService[models.Role](roleCol)
	newHandler.InitService = *services.NewInitService(c, db)
	return newHandler
}

// SetAdministrator tạo người dùng quản trị hệ thống
func (h *InitHandler) SetAdministrator(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	// GET ID
	id := ctx.UserValue("id").(string)

	result, err := h.InitService.SetAdministrator(utility.String2ObjectID(id))
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		response = utility.FinalResponse(nil, err)
	} else {
		ctx.SetStatusCode(http.StatusOK)
		response = utility.FinalResponse(result, nil)
	}

	utility.JSON(ctx, response)
}
