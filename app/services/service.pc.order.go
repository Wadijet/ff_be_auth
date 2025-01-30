package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"encoding/json"
	"errors"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AccessTokenService là cấu trúc chứa các phương thức liên quan đến người dùng
type PcOrderService struct {
	crudPcOrder RepositoryService
}

// Khởi tạo AccessTokenService với cấu hình và kết nối cơ sở dữ liệu
func NewPcOrderService(c *config.Configuration, db *mongo.Client) *PcOrderService {
	newService := new(PcOrderService)
	newService.crudPcOrder = *NewRepository(c, db, global.MongoDB_ColNames.PcOrders)
	return newService
}

// Nhận data từ Pancake và lưu vào cơ sở dữ liệu
func (h *PcOrderService) ReviceData(ctx *fasthttp.RequestCtx, credential *models.PcOrderCreateInput) (CreateResult interface{}, err error) {

	if credential.PanCakeData == nil {
		return nil, errors.New("ApiData is required")
	}

	// Lấy thông tin OrderID từ ApiData đưa vào biến
	pancakeOrderId := credential.PanCakeData["id"]
	pancakeOrderNumber := pancakeOrderId.(json.Number)
	pancakeOrderIdStr := pancakeOrderNumber.String()
	
	// Kiểm tra PcOrder đã tồn tại chưa
	filter := bson.M{"pancakeOrderId": pancakeOrderIdStr}
	checkResult, _ := h.crudPcOrder.FindOne(ctx, filter, nil)
	if checkResult == nil { // Nếu PcOrder chưa tồn tại thì tạo mới
		// Tạo một PcOrder mới
		newPcOrder := models.PcOrder{}
		newPcOrder.PancakeOrderId = pancakeOrderIdStr
		newPcOrder.PanCakeData = credential.PanCakeData

		// Thêm PcOrder vào cơ sở dữ liệu
		return h.crudPcOrder.InsertOne(ctx, newPcOrder)

	} else { // Nếu PcOrder đã tồn tại thì cập nhật thông tin mới
		// chuyển đổi checkResult từ interface{} sang models.PcOrder
		var oldPcOrder models.PcOrder
		bsonBytes, err := bson.Marshal(checkResult)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(bsonBytes, &oldPcOrder)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin mới
		oldPcOrder.PanCakeData = credential.PanCakeData

		CustomBson := &utility.CustomBson{}
		change, err := CustomBson.Set(oldPcOrder)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin mới vào cơ sở dữ liệu
		return h.crudPcOrder.UpdateOneById(ctx, oldPcOrder.ID, change)
	}
}

// Tìm một PcOrder theo ID
func (h *PcOrderService) FindOneById(ctx *fasthttp.RequestCtx, id string) (FindResult interface{}, err error) {
	return h.crudPcOrder.FindOneById(ctx, utility.String2ObjectID(id), nil)
}

// Tìm tất cả các PcOrder với phân trang
func (h *PcOrderService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (FindResult interface{}, err error) {

	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})
	return h.crudPcOrder.FindAllWithPaginate(ctx, nil, opts)
}
