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
type FbMessageService struct {
	crudFbMessage RepositoryService
}

// Khởi tạo AccessTokenService với cấu hình và kết nối cơ sở dữ liệu
func NewFbMessageService(c *config.Configuration, db *mongo.Client) *FbMessageService {
	newService := new(FbMessageService)
	newService.crudFbMessage = *NewRepository(c, db, global.MongoDB_ColNames.FbMessages)
	return newService
}

// Nhận data từ Facebook và lưu vào cơ sở dữ liệu
func (h *FbMessageService) ReviceData(ctx *fasthttp.RequestCtx, credential *models.FbMessageCreateInput) (CreateResult interface{}, err error) {

	if credential.PanCakeData == nil {
		return nil, errors.New("ApiData is required")
	}

	// Lấy thông tin MessageId từ ApiData đưa vào biến
	conversationId := credential.PanCakeData["conversation_id"].(string)

	// Kiểm tra FbMessage đã tồn tại chưa
	filter := bson.M{"conversationId": conversationId, "customerId": credential.CustomerId}
	checkResult, _ := h.crudFbMessage.FindOne(ctx, filter, nil)
	if checkResult == nil { // Nếu FbMessage chưa tồn tại thì tạo mới
		// Tạo một FbMessage mới
		newFbMessage := models.FbMessage{}
		newFbMessage.PageId = credential.PageId
		newFbMessage.PageUsername = credential.PageUsername
		newFbMessage.PanCakeData = credential.PanCakeData
		newFbMessage.CustomerId = credential.CustomerId
		newFbMessage.ConversationId = conversationId

		// Thêm FbMessage vào cơ sở dữ liệu
		return h.crudFbMessage.InsertOne(ctx, newFbMessage)

	} else { // Nếu FbMessage đã tồn tại thì cập nhật thông tin mới
		// chuyển đổi checkResult từ interface{} sang models.FbMessage
		var oldFbMessage models.FbMessage
		bsonBytes, err := bson.Marshal(checkResult)
		if err != nil {
			return nil, err
		}

		// chuyển đổi bsonBytes sang models.FbMessage
		err = bson.Unmarshal(bsonBytes, &oldFbMessage)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin mới
		oldFbMessage.PanCakeData = credential.PanCakeData
		oldFbMessage.PageId = credential.PageId
		oldFbMessage.PageUsername = credential.PageUsername
		oldFbMessage.ConversationId = conversationId
		oldFbMessage.CustomerId = credential.CustomerId

		CustomBson := &utility.CustomBson{}
		change, err := CustomBson.Set(oldFbMessage)
		if err != nil {
			return nil, err
		}

		// Cập nhật vào cơ sở dữ liệu
		return h.crudFbMessage.UpdateOneById(ctx, oldFbMessage.ID, change)
	}
}

// Tìm một FbMessage theo ID
func (h *FbMessageService) FindOneById(ctx *fasthttp.RequestCtx, id string) (FindResult interface{}, err error) {
	return h.crudFbMessage.FindOneById(ctx, utility.String2ObjectID(id), nil)
}

// Tìm một FbMessage theo ConversationID
func (h *FbMessageService) FindOneByConversationID(ctx *fasthttp.RequestCtx, conversationID string) (FindResult interface{}, err error) {
	// Tạo điều kiện tìm kiếm
	filter := bson.M{"conversationId": conversationID}
	return h.crudFbMessage.FindOne(ctx, filter, nil)
}

// Tìm tất cả các FbMessage với phân trang
func (h *FbMessageService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (FindResult interface{}, err error) {

	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})

	return h.crudFbMessage.FindAllWithPaginate(ctx, nil, opts)
}
