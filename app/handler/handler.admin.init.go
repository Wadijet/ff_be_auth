package handler

import (
	"meta_commerce/app/database/registry"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
	"net/http"

	"github.com/valyala/fasthttp"
)

// InitHandler là struct chứa các CRUD services và InitService
type InitHandler struct {
	UserCRUD       services.BaseServiceMongo[models.User]
	PermissionCRUD services.BaseServiceMongo[models.Permission]
	RoleCRUD       services.BaseServiceMongo[models.Role]
	InitService    services.InitService
}

// NewInitHandler khởi tạo InitHandler mới
func NewInitHandler() *InitHandler {
	newHandler := new(InitHandler)

	// Khởi tạo các collection từ registry
	userCol := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Users)
	permissionCol := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Permissions)
	roleCol := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Roles)

	// Khởi tạo các service với BaseService
	newHandler.UserCRUD = services.NewBaseServiceMongo[models.User](userCol)
	newHandler.PermissionCRUD = services.NewBaseServiceMongo[models.Permission](permissionCol)
	newHandler.RoleCRUD = services.NewBaseServiceMongo[models.Role](roleCol)
	newHandler.InitService = *services.NewInitService()
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
