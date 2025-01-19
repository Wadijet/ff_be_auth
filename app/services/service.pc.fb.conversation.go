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
type FbConversationService struct {
	crudFbConversation RepositoryService
}

// Khởi tạo AccessTokenService với cấu hình và kết nối cơ sở dữ liệu
func NewFbConversationService(c *config.Configuration, db *mongo.Client) *FbConversationService {
	newService := new(FbConversationService)
	newService.crudFbConversation = *NewRepository(c, db, global.MongoDB_ColNames.FbConvesations)
	return newService
}

// Nhận data từ Facebook và lưu vào cơ sở dữ liệu
func (h *FbConversationService) ReviceData(ctx *fasthttp.RequestCtx, credential *models.FbConversationCreateInput) (CreateResult interface{}, err error) {

	if credential.PanCakeData == nil {
		return nil, errors.New("ApiData is required")
	}

	// Lấy thông tin ConversationID từ ApiData đưa vào biến
	conversationId := credential.PanCakeData["id"].(string)
	customerId := credential.PanCakeData["customer_id"].(string)

	// Kiểm tra FbConversation đã tồn tại chưa
	filter := bson.M{"conversationId": conversationId}
	checkResult, _ := h.crudFbConversation.FindOne(ctx, filter, nil)
	if checkResult == nil { // Nếu FbConversation chưa tồn tại thì tạo mới
		// Tạo một FbConversation mới
		newFbConversation := models.FbConversation{}
		newFbConversation.PageId = credential.PageId
		newFbConversation.PageUsername = credential.PageUsername
		newFbConversation.PanCakeData = credential.PanCakeData
		newFbConversation.ConversationId = conversationId
		newFbConversation.CustomerId = customerId

		// Thêm FbConversation vào cơ sở dữ liệu
		return h.crudFbConversation.InsertOne(ctx, newFbConversation)

	} else { // Nếu FbConversation đã tồn tại thì cập nhật thông tin mới
		// chuyển đổi checkResult từ interface{} sang models.FbConversation
		var oldFbConversation models.FbConversation
		bsonBytes, err := bson.Marshal(checkResult)
		if err != nil {
			return nil, err
		}

		err = bson.Unmarshal(bsonBytes, &oldFbConversation)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin mới
		oldFbConversation.PanCakeData = credential.PanCakeData
		oldFbConversation.PageId = credential.PageId
		oldFbConversation.PageUsername = credential.PageUsername
		oldFbConversation.ConversationId = conversationId
		oldFbConversation.CustomerId = customerId

		CustomBson := &utility.CustomBson{}
		change, err := CustomBson.Set(oldFbConversation)
		if err != nil {
			return nil, err
		}

		// Cập nhật FbConversation vào cơ sở dữ liệu
		return h.crudFbConversation.UpdateOneById(ctx, oldFbConversation.ID, change)
	}
}

// Tìm một FbConversation theo ID
func (h *FbConversationService) FindOneByConversationID(ctx *fasthttp.RequestCtx, id string) (FindResult interface{}, err error) {
	return h.crudFbConversation.FindOneById(ctx, utility.String2ObjectID(id), nil)
}

// Tìm tất cả các FbConversation với phân trang
func (h *FbConversationService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (FindResult interface{}, err error) {

	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})

	return h.crudFbConversation.FindAllWithPaginate(ctx, nil, opts)
}
