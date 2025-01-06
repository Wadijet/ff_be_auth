package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"errors"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AccessTokenService là cấu trúc chứa các phương thức liên quan đến người dùng
type FbPageService struct {
	crudFbPage RepositoryService
}

// Khởi tạo AccessTokenService với cấu hình và kết nối cơ sở dữ liệu
func NewFbPageService(c *config.Configuration, db *mongo.Client) *FbPageService {
	newService := new(FbPageService)
	newService.crudFbPage = *NewRepository(c, db, global.MongoDB_ColNames.FbPages)
	return newService
}

// Tạo mới một FbPage
func (h *FbPageService) ReviceData(ctx *fasthttp.RequestCtx, credential *models.FbPageCreateInput) (CreateResult interface{}, err error) {

	if credential.ApiData == nil {
		return nil, errors.New("ApiData is required")
	}

	// Lấy thông tin PageID từ ApiData đưa vào biến
	pageID := credential.ApiData["id"].(string)

	// Kiểm tra FbPage đã tồn tại chưa
	filter := bson.M{"pageID": pageID}
	checkResult, _ := h.crudFbPage.FindOne(ctx, filter, nil)
	if checkResult == nil { // Nếu FbPage chưa tồn tại thì tạo mới
		// Tạo một FbPage mới
		newFbPage := models.FbPage{}
		newFbPage.ApiData = credential.ApiData
		newFbPage.PageName = credential.ApiData["name"].(string)
		newFbPage.PageId = credential.ApiData["id"].(string)

		// Thêm FbPage vào cơ sở dữ liệu
		return h.crudFbPage.InsertOne(ctx, newFbPage)
	} else { // Nếu FbPage đã tồn tại thì cập nhật thông tin mới
		// chuyển đổi checkResult từ interface{} sang models.FbPage
		var oldFbPage models.FbPage
		bsonBytes, err := bson.Marshal(checkResult)
		if err != nil {
			return nil, err
		}

		err = bson.Unmarshal(bsonBytes, &oldFbPage)
		if err != nil {
			return nil, err
		}

		oldFbPage.ApiData = credential.ApiData
		oldFbPage.PageName = credential.ApiData["name"].(string)
		//oldFbPage.PageID = newApiData["id"].(string)

		CustomBson := &utility.CustomBson{}
		change, err := CustomBson.Set(oldFbPage)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin FbPage
		return h.crudFbPage.UpdateOneById(ctx, oldFbPage.ID, change)
	}

}

// Tìm một FbPage theo ID
func (h *FbPageService) FindOneById(ctx *fasthttp.RequestCtx, id string) (FindResult interface{}, err error) {
	return h.crudFbPage.FindOneById(ctx, utility.String2ObjectID(id), nil)
}

// Tìm một FbPage theo PageID
func (h *FbPageService) FindOneByPageID(ctx *fasthttp.RequestCtx, pageID string) (FindResult interface{}, err error) {
	// Tạo điều kiện tìm kiếm
	filter := bson.M{"pageID": pageID}
	return h.crudFbPage.FindOne(ctx, filter, nil)
}

// Tìm tất cả các FbPage với phân trang
func (h *FbPageService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (FindResult interface{}, err error) {
	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})
	return h.crudFbPage.FindAll(ctx, nil, opts)
}
