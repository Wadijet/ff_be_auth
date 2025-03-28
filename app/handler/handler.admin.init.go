package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
	"meta_commerce/config"
	"meta_commerce/global"
	"net/http"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitHandler là struct chứa các CRUD services và InitService
type InitHandler struct {
	UserCRUD       services.BaseServiceMongo[models.User]
	PermissionCRUD services.BaseServiceMongo[models.Permission]
	RoleCRUD       services.BaseServiceMongo[models.Role]
	InitService    services.InitService
}

// NewInitHandler khởi tạo InitHandler mới
func NewInitHandler(c *config.Configuration, db *mongo.Client) *InitHandler {
	newHandler := new(InitHandler)

	// Khởi tạo các collection
	userCol := db.Database(services.GetDBNameFromCollectionName(c, global.MongoDB_ColNames.Users)).Collection(global.MongoDB_ColNames.Users)
	permissionCol := db.Database(services.GetDBNameFromCollectionName(c, global.MongoDB_ColNames.Permissions)).Collection(global.MongoDB_ColNames.Permissions)
	roleCol := db.Database(services.GetDBNameFromCollectionName(c, global.MongoDB_ColNames.Roles)).Collection(global.MongoDB_ColNames.Roles)

	// Khởi tạo các service với BaseService
	newHandler.UserCRUD = services.NewBaseServiceMongo[models.User](userCol)
	newHandler.PermissionCRUD = services.NewBaseServiceMongo[models.Permission](permissionCol)
	newHandler.RoleCRUD = services.NewBaseServiceMongo[models.Role](roleCol)
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
