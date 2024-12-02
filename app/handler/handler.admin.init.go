package handler

import (
	"atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitHandler là struct chứa các CRUD services và InitService
type InitHandler struct {
	UserCRUD       services.Repository
	PermissionCRUD services.Repository
	RoleCRUD       services.Repository
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

// InitPermission khởi tạo các quyền trong hệ thống
func (h *InitHandler) InitPermission(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	var Permissions []interface{}

	// Thêm các quyền vào danh sách Permissions
	Permissions = append(Permissions, models.PermissionCreateInput{"permission.create", "", "permission"})
	Permissions = append(Permissions, models.PermissionCreateInput{"permission.read", "", "permission"})
	Permissions = append(Permissions, models.PermissionCreateInput{"permission.read_all", "", "permission"})
	Permissions = append(Permissions, models.PermissionCreateInput{"permission.update", "", "permission"})
	Permissions = append(Permissions, models.PermissionCreateInput{"permission.update_all", "", "permission"})
	Permissions = append(Permissions, models.PermissionCreateInput{"permission.delete", "", "permission"})
	Permissions = append(Permissions, models.PermissionCreateInput{"permission.delete_all", "", "permission"})

	Permissions = append(Permissions, models.PermissionCreateInput{"role.create", "", "role"})
	Permissions = append(Permissions, models.PermissionCreateInput{"role.read", "", "role"})
	Permissions = append(Permissions, models.PermissionCreateInput{"role.read_all", "", "role"})
	Permissions = append(Permissions, models.PermissionCreateInput{"role.update", "", "role"})
	Permissions = append(Permissions, models.PermissionCreateInput{"role.update_all", "", "role"})
	Permissions = append(Permissions, models.PermissionCreateInput{"role.delete", "", "role"})
	Permissions = append(Permissions, models.PermissionCreateInput{"role.delete_all", "", "role"})

	Permissions = append(Permissions, models.PermissionCreateInput{"user.create", "", "user"})
	Permissions = append(Permissions, models.PermissionCreateInput{"user.read", "", "user"})
	Permissions = append(Permissions, models.PermissionCreateInput{"user.read_all", "", "user"})
	Permissions = append(Permissions, models.PermissionCreateInput{"user.update", "", "user"})
	Permissions = append(Permissions, models.PermissionCreateInput{"user.update_all", "", "user"})
	Permissions = append(Permissions, models.PermissionCreateInput{"user.delete", "", "user"})
	Permissions = append(Permissions, models.PermissionCreateInput{"user.delete_all", "", "user"})

	Permissions = append(Permissions, models.PermissionCreateInput{"mtservice.create", "", "mtservice"})
	Permissions = append(Permissions, models.PermissionCreateInput{"mtservice.read", "", "mtservice"})
	Permissions = append(Permissions, models.PermissionCreateInput{"mtservice.read_all", "", "mtservice"})
	Permissions = append(Permissions, models.PermissionCreateInput{"mtservice.update", "", "mtservice"})
	Permissions = append(Permissions, models.PermissionCreateInput{"mtservice.update_all", "", "mtservice"})
	Permissions = append(Permissions, models.PermissionCreateInput{"mtservice.delete", "", "mtservice"})
	Permissions = append(Permissions, models.PermissionCreateInput{"mtservice.delete_all", "", "mtservice"})

	Permissions = append(Permissions, models.PermissionCreateInput{"admin.set_role", "", "admin"})
	Permissions = append(Permissions, models.PermissionCreateInput{"admin.block_user", "", "admin"})
	Permissions = append(Permissions, models.PermissionCreateInput{"admin.unblock_user", "", "admin"})

	// Chèn các quyền vào cơ sở dữ liệu
	response = utility.FinalResponse(h.PermissionCRUD.InsertMany(ctx, Permissions))

	// Trả về kết quả dưới dạng JSON
	utility.JSON(ctx, response)
}

// CreatePermissions tạo danh sách chi tiết quyền từ danh sách quyền đầu vào
func CreatePermissions(listPermissions interface{}, inputPermissions []string, outputPermissions *[]models.PermissionDetail) {

	if inputPermissions == nil {
		for _, permission := range listPermissions.([]primitive.M) {
			permissionMap, err := utility.ToMap(permission)
			if err != nil {
				continue
			}

			newPermissonDetail := new(models.PermissionDetail)
			newPermissonDetail.ID = permissionMap["_id"].(primitive.ObjectID)
			newPermissonDetail.Name = permissionMap["name"].(string)

			*outputPermissions = append(*outputPermissions, *newPermissonDetail)

		}
	} else {
		for _, permission := range listPermissions.([]primitive.M) {
			permissionMap, err := utility.ToMap(permission)
			if err != nil {
				continue
			}

			var PermissionCheck = false
			for _, notUserPermission := range inputPermissions {
				if permissionMap["name"] == notUserPermission {
					PermissionCheck = true
					break
				}
			}

			if PermissionCheck == true {
				newPermissonDetail := new(models.PermissionDetail)
				newPermissonDetail.ID = permissionMap["_id"].(primitive.ObjectID)
				newPermissonDetail.Name = permissionMap["name"].(string)

				*outputPermissions = append(*outputPermissions, *newPermissonDetail)
			}
		}
	}
}

// InitRole khởi tạo các vai trò trong hệ thống
func (h *InitHandler) InitRole(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy tất cả các quyền trong hệ thống
	permissions, err := h.PermissionCRUD.FindAll(ctx, bson.D{}, nil)
	if err == nil {
		var Roles []interface{}

		// Tạo role ADMIN
		// Tạo danh sách quyền của role ADMIN (tất cả các quyền)
		{
			var admin_permissions []models.PermissionDetail
			CreatePermissions(permissions, nil, &admin_permissions)
			Roles = append(Roles, models.RoleInitInput{
				Name:        "Admin",
				Describe:    "System admin",
				Permissions: admin_permissions})
		}

		// Tạo role User
		// Tạo danh sách quyền của Role User
		{
			UserPermissions := []string{
				"user.create",
				"user.read",
				"user.update",
			}

			var user_permissions []models.PermissionDetail
			CreatePermissions(permissions, UserPermissions, &user_permissions)
			Roles = append(Roles, models.RoleInitInput{
				Name:        "User",
				Describe:    "User",
				Permissions: user_permissions})
		}

		// Tạo role MtService
		// Tạo danh sách quyền của Role MtService
		{
			MtServicePermissions := []string{
				"mtservice.create",
				"mtservice.read",
				"mtservice.update",
			}

			var MtService_permissions []models.PermissionDetail
			CreatePermissions(permissions, MtServicePermissions, &MtService_permissions)
			Roles = append(Roles, models.RoleInitInput{
				Name:        "MtService",
				Describe:    "MtService",
				Permissions: MtService_permissions})
		}

		// Chèn các vai trò vào cơ sở dữ liệu
		response = utility.FinalResponse(h.RoleCRUD.InsertMany(ctx, Roles))

	}

	// Trả về kết quả dưới dạng JSON
	utility.JSON(ctx, response)
}

// SetAdminStruct là struct chứa thông tin email của admin
type SetAdminStruct struct {
	Email string `json: "email" bson "email" validate:"required"`
}

// SetAdmin thiết lập người dùng thành admin
func (h *InitHandler) SetAdmin(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu từ request
	postValues := ctx.PostBody()
	inputStruct := new(SetAdminStruct)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic
			response = utility.FinalResponse(h.InitService.SetAdmin(ctx, inputStruct.Email))
		}
	}

	// Trả về kết quả dưới dạng JSON
	utility.JSON(ctx, response)
}
