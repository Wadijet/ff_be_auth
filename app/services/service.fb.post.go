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
type FbPostService struct {
	crudFbPost RepositoryService
}

// Khởi tạo AccessTokenService với cấu hình và kết nối cơ sở dữ liệu
func NewFbPostService(c *config.Configuration, db *mongo.Client) *FbPostService {
	newService := new(FbPostService)
	newService.crudFbPost = *NewRepository(c, db, global.MongoDB_ColNames.FbPosts)
	return newService
}

// Nhận data từ Facebook và lưu vào cơ sở dữ liệu
func (h *FbPostService) ReviceData(ctx *fasthttp.RequestCtx, credential *models.FbPostCreateInput) (CreateResult interface{}, err error) {

	if credential.PanCakeData == nil {
		return nil, errors.New("ApiData is required")
	}

	// Lấy thông tin PostId từ ApiData đưa vào biến
	pageId := credential.PanCakeData["page_id"].(string)
	postId := credential.PanCakeData["id"].(string)

	// Kiểm tra FbPost đã tồn tại chưa
	filter := bson.M{"postId": postId}
	checkResult, _ := h.crudFbPost.FindOne(ctx, filter, nil)
	if checkResult == nil { // Nếu FbPost chưa tồn tại thì tạo mới
		// Tạo một FbPost mới
		newFbPost := models.FbPost{}
		newFbPost.PageId = pageId
		newFbPost.PostId = postId
		newFbPost.PanCakeData = credential.PanCakeData

		// Thêm FbPost vào cơ sở dữ liệu
		return h.crudFbPost.InsertOne(ctx, newFbPost)

	} else { // Nếu FbPost đã tồn tại thì cập nhật thông tin mới
		// chuyển đổi checkResult từ interface{} sang models.FbPost
		var oldFbPost models.FbPost
		bsonBytes, err := bson.Marshal(checkResult)
		if err != nil {
			return nil, err
		}

		// chuyển đổi bsonBytes sang models.FbPost
		err = bson.Unmarshal(bsonBytes, &oldFbPost)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin mới
		oldFbPost.PanCakeData = credential.PanCakeData

		CustomBson := &utility.CustomBson{}
		change, err := CustomBson.Set(oldFbPost)
		if err != nil {
			return nil, err
		}

		// Cập nhật vào cơ sở dữ liệu
		return h.crudFbPost.UpdateOneById(ctx, oldFbPost.ID, change)
	}
}

// FindOneById tìm kiếm một bài viết theo ID
func (h *FbPostService) FindOneById(ctx *fasthttp.RequestCtx, Id string) (FindResult interface{}, err error) {
	filter := bson.M{"postId": utility.String2ObjectID(Id)}
	return h.crudFbPost.FindOne(ctx, filter, nil)
}

// FindAll tìm kiếm tất cả bài viết
func (h *FbPostService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (FindResult interface{}, err error) {
	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})

	return h.crudFbPost.FindAllWithPaginate(ctx, nil, opts)
}
