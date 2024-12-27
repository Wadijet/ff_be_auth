package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"
	"errors"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService là cấu trúc chứa các phương thức liên quan đến người dùng
type RolePermissionService struct {
	crudRole           RepositoryService
	crudPermission     RepositoryService
	crudRolePermission RepositoryService
}

// Khởi tạo UserService với cấu hình và kết nối cơ sở dữ liệu
func NewRolePermissionService(c *config.Configuration, db *mongo.Client) *RolePermissionService {
	newService := new(RolePermissionService)
	newService.crudRole = *NewRepository(c, db, global.MongoDB_ColNames.Roles)
	newService.crudRole = *NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newService.crudRole = *NewRepository(c, db, global.MongoDB_ColNames.RolePermissions)
	return newService
}

// Tạo mới một RolePermission
func (h *RolePermissionService) Create(ctx *fasthttp.RequestCtx, credential *models.RolePermissionCreateInput) (CreateResult interface{}, err error) {
	// Kiểm tra Role có tồn tại không
	roleFilter := bson.M{"_id": credential.RoleID}
	roleResult, _ := h.crudRole.FindOne(ctx, roleFilter, nil)
	if roleResult == nil {
		return nil, errors.New("Role not found")
	}

	// Kiểm tra Permission có tồn tại không
	permissionFilter := bson.M{"_id": credential.PermissionID}
	permissionResult, _ := h.crudPermission.FindOne(ctx, permissionFilter, nil)
	if permissionResult == nil {
		return nil, errors.New("Permission not found")
	}

	// Kiểm tra RolePermission đã tồn tại chưa
	filter := bson.M{"role_id": credential.RoleID, "permission_id": credential.PermissionID, "scope": credential.Scope}
	checkResult, _ := h.crudRolePermission.FindOne(ctx, filter, nil)
	if checkResult != nil {
		return nil, errors.New("RolePermission already exists")
	}

	// Tạo mới RolePermission
	newRolePermission := models.RolePermission{
		RoleID:       credential.RoleID,
		PermissionID: credential.PermissionID,
		Scope:        credential.Scope,
	}

	// Thêm RolePermission vào cơ sở dữ liệu
	result, err := h.crudRolePermission.InsertOne(ctx, newRolePermission)
	if err != nil {
		return nil, err
	}

	return result, nil
}
