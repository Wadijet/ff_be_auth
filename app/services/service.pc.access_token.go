package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AccessTokenService là cấu trúc chứa các phương thức liên quan đến người dùng
type AccessTokenService struct {
	crudAccessToken RepositoryService
}

// Khởi tạo AccessTokenService với cấu hình và kết nối cơ sở dữ liệu
func NewAccessTokenService(c *config.Configuration, db *mongo.Client) *AccessTokenService {
	newService := new(AccessTokenService)
	newService.crudAccessToken = *NewRepository(c, db, global.MongoDB_ColNames.AccessTokens)
	return newService
}

// Tạo mới một Access Token
func (h *AccessTokenService) Create(ctx *fasthttp.RequestCtx, credential *models.AccessTokenCreateInput) (CreateResult interface{}, err error) {

	// Chuyển đổi credential.AssignedUsers từ mảng []string sang mảng []ObjectID
	newAssignedUsers := make([]primitive.ObjectID, 0)
	for _, userID := range credential.AssignedUsers {
		newAssignedUsers = append(newAssignedUsers, utility.String2ObjectID(userID))
	}

	// Tạo mới Access Token
	newAccessToken := models.AccessToken{
		Name:          credential.Name,
		Describe:      credential.Describe,
		System:        credential.System,
		Value:         credential.Value,
		AssignedUsers: newAssignedUsers,
		Status:        0,
	}

	// Thêm Access Token vào cơ sở dữ liệu
	return h.crudAccessToken.InsertOne(ctx, newAccessToken)
}

// Tìm một Access Token theo ID
func (h *AccessTokenService) FindOneById(ctx *fasthttp.RequestCtx, id string) (FindResult interface{}, err error) {
	return h.crudAccessToken.FindOneById(ctx, utility.String2ObjectID(id), nil)
}

// Tìm tất cả các Access Token với phân trang
func (h *AccessTokenService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (FindResult interface{}, err error) {

	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})
	return h.crudAccessToken.FindAllWithPaginate(ctx, nil, opts)
}

// Cập nhật một Access Token theo ID
func (h *AccessTokenService) UpdateOneById(ctx *fasthttp.RequestCtx, id string, credential *models.AccessTokenUpdateInput) (UpdateResult interface{}, err error) {

	// Chuyển đổi credential.AssignedUsers từ mảng []string sang mảng []ObjectID
	newAssignedUsers := make([]primitive.ObjectID, 0)
	for _, userID := range credential.AssignedUsers {
		newAssignedUsers = append(newAssignedUsers, utility.String2ObjectID(userID))
	}

	// Cập nhật Access Token
	updateFields := bson.M{
		"name":          credential.Name,
		"describe":      credential.Describe,
		"system":        credential.System,
		"value":         credential.Value,
		"assignedUsers": newAssignedUsers,
	}

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(updateFields)
	if err != nil {
		return nil, err
	}

	return h.crudAccessToken.UpdateOneById(ctx, utility.String2ObjectID(id), change)
}

// Xóa một Access Token theo ID
func (h *AccessTokenService) DeleteOneById(ctx *fasthttp.RequestCtx, id string) (DeleteResult interface{}, err error) {
	return h.crudAccessToken.DeleteOneById(ctx, utility.String2ObjectID(id))
}
