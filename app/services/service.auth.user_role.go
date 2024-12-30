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
type UserRoleService struct {
	crudUser     RepositoryService
	crudRole     RepositoryService
	crudUserRole RepositoryService
}

// Khởi tạo UserRoleService với cấu hình và kết nối cơ sở dữ liệu
func NewUserRoleService(c *config.Configuration, db *mongo.Client) *UserRoleService {
	newService := new(UserRoleService)
	newService.crudUser = *NewRepository(c, db, global.MongoDB_ColNames.Users)
	newService.crudRole = *NewRepository(c, db, global.MongoDB_ColNames.Roles)
	newService.crudUserRole = *NewRepository(c, db, global.MongoDB_ColNames.UserRoles)
	return newService
}

// Tạo mới một UserRole
func (h *UserRoleService) Create(ctx *fasthttp.RequestCtx, credential *models.UserRoleCreateInput) (CreateResult interface{}, err error) {
	// Kiểm tra User có tồn tại không
	userFilter := bson.M{"_id": credential.UserID}
	userResult, _ := h.crudUser.FindOne(ctx, userFilter, nil)
	if userResult == nil {
		return nil, errors.New("User not found")
	}

	// Kiểm tra Role có tồn tại không
	roleFilter := bson.M{"_id": credential.RoleID}
	roleResult, _ := h.crudRole.FindOne(ctx, roleFilter, nil)
	if roleResult == nil {
		return nil, errors.New("Role not found")
	}

	// Kiểm tra UserRole đã tồn tại chưa
	filter := bson.M{"user_id": credential.UserID, "role_id": credential.RoleID}
	checkResult, _ := h.crudUserRole.FindOne(ctx, filter, nil)
	if checkResult != nil {
		return nil, errors.New("UserRole already exists")
	}

	// Tạo mới UserRole
	newUserRole := models.UserRole{
		UserID: credential.UserID,
		RoleID: credential.RoleID,
	}

	// Thêm UserRole vào cơ sở dữ liệu
	insertResult, err := h.crudUserRole.InsertOne(ctx, newUserRole)
	if err != nil {
		return nil, err
	}

	return insertResult, nil
}

// Xóa một UserRole
func (h *UserRoleService) Delete(ctx *fasthttp.RequestCtx, id string) (DeleteResult interface{}, err error) {
	// Kiểm tra UserRole có tồn tại không
	filter := bson.M{"_id": id}
	checkResult, _ := h.crudUserRole.FindOne(ctx, filter, nil)
	if checkResult == nil {
		return nil, errors.New("UserRole not found")
	}

	// Xóa UserRole khỏi cơ sở dữ liệu
	deleteResult, err := h.crudUserRole.DeleteOneById(ctx, id)
	if err != nil {
		return nil, err
	}

	return deleteResult, nil
}
