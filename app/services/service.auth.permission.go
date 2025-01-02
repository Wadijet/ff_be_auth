package services

import (
	"atk-go-server/config"
	"atk-go-server/global"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PermissionService là cấu trúc chứa các phương thức liên quan đến người dùng
type PermissionService struct {
	crudRole RepositoryService
}

// Khởi tạo PermissionService với cấu hình và kết nối cơ sở dữ liệu
func NewPermissionService(c *config.Configuration, db *mongo.Client) *PermissionService {
	newService := new(PermissionService)
	newService.crudRole = *NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	return newService
}

// Tìm một Permission theo ID
func (h *PermissionService) FindOneById(ctx *fasthttp.RequestCtx, id string) (FindResult interface{}, err error) {
	return h.crudRole.FindOneById(ctx, id, nil)
}

// Tìm tất cả các Permission với phân trang
func (h *PermissionService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (FindResult interface{}, err error) {

	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})

	return h.crudRole.FindAllWithPaginate(ctx, bson.D{}, opts)
}
