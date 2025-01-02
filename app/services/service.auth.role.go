package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"
	"errors"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RoleService là cấu trúc chứa các phương thức liên quan đến người dùng
type RoleService struct {
	crudRole RepositoryService
}

// Khởi tạo RoleService với cấu hình và kết nối cơ sở dữ liệu
func NewRoleService(c *config.Configuration, db *mongo.Client) *RoleService {
	newService := new(RoleService)
	newService.crudRole = *NewRepository(c, db, global.MongoDB_ColNames.Roles)
	return newService
}

// Tạo mới một Role
func (h *RoleService) Create(ctx *fasthttp.RequestCtx, credential *models.RoleCreateInput) (CreateResult interface{}, err error) {
	// Kiểm tra tên của Role đã tồn tại chưa
	filter := bson.M{"name": credential.Name}
	checkResult, _ := h.crudRole.FindOne(ctx, filter, nil)
	if checkResult != nil {
		return nil, errors.New("Role already exists")
	}

	// Tạo mới Role
	newRole := models.Role{
		Name:     credential.Name,
		Describe: credential.Describe,
	}

	// Thêm Role vào cơ sở dữ liệu
	return h.crudRole.InsertOne(ctx, newRole)
}

// Tìm một Role theo ID
func (h *RoleService) FindOneById(ctx *fasthttp.RequestCtx, id string) (FindResult interface{}, err error) {
	return h.crudRole.FindOneById(ctx, id, nil)
}

// Tìm tất cả các Role với phân trang
func (h *RoleService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (FindResult interface{}, err error) {

	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})

	return h.crudRole.FindAllWithPaginate(ctx, bson.D{}, opts)
}

// Cập nhật một Role theo ID
func (h *RoleService) Update(ctx *fasthttp.RequestCtx, id string, credential *models.RoleUpdateInput) (UpdateResult interface{}, err error) {
	// Kiểm tra Role đã tồn tại chưa
	filter := bson.M{"_id": id}
	checkResult, _ := h.crudRole.FindOne(ctx, filter, nil)
	if checkResult == nil {
		return nil, errors.New("Role not found")
	}

	// Cập nhật Role
	update := bson.M{"$set": bson.M{
		"name":     credential.Name,
		"describe": credential.Describe,
	}}

	return h.crudRole.UpdateOneById(ctx, id, update)
}

// Xóa một Role theo ID
func (h *RoleService) Delete(ctx *fasthttp.RequestCtx, id string) (DeleteResult interface{}, err error) {
	return h.crudRole.DeleteOneById(ctx, id)
}
