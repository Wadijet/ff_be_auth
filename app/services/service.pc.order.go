package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PcOrderService là cấu trúc chứa các phương thức liên quan đến đơn hàng
type PcOrderService struct {
	*BaseServiceImpl[models.PcOrder]
}

// NewPcOrderService tạo mới PcOrderService
func NewPcOrderService(c *config.Configuration, db *mongo.Client) *PcOrderService {
	pcOrderCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.PcOrders)).Collection(global.MongoDB_ColNames.PcOrders)
	return &PcOrderService{
		BaseServiceImpl: NewBaseService[models.PcOrder](pcOrderCollection),
	}
}

// IsPancakeOrderIdExist kiểm tra ID đơn hàng Pancake có tồn tại hay không
func (s *PcOrderService) IsPancakeOrderIdExist(ctx context.Context, pancakeOrderId string) (bool, error) {
	filter := bson.M{"pancakeOrderId": pancakeOrderId}
	var pcOrder models.PcOrder
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&pcOrder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// FindOne tìm một document theo ObjectId
func (s *PcOrderService) FindOne(ctx context.Context, id primitive.ObjectID) (models.PcOrder, error) {
	filter := bson.M{"_id": id}
	var result models.PcOrder
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return models.PcOrder{}, err
	}
	return result, nil
}

// Delete xóa một document theo ObjectId
func (s *PcOrderService) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := s.BaseServiceImpl.collection.DeleteOne(ctx, filter)
	return err
}

// Update cập nhật một document theo ObjectId
func (s *PcOrderService) Update(ctx context.Context, id primitive.ObjectID, pcOrder models.PcOrder) (models.PcOrder, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": pcOrder}

	_, err := s.BaseServiceImpl.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.PcOrder{}, err
	}

	return s.FindOne(ctx, id)
}

// ReviceData nhận data từ Pancake và lưu vào cơ sở dữ liệu
func (s *PcOrderService) ReviceData(ctx context.Context, input *models.PcOrderCreateInput) (*models.PcOrder, error) {
	if input.PanCakeData == nil {
		return nil, errors.New("ApiData is required")
	}

	// Lấy thông tin OrderID từ ApiData đưa vào biến
	pancakeOrderId := input.PanCakeData["id"]
	pancakeOrderNumber := pancakeOrderId.(json.Number)
	pancakeOrderIdStr := pancakeOrderNumber.String()

	// Kiểm tra PcOrder đã tồn tại chưa
	exists, err := s.IsPancakeOrderIdExist(ctx, pancakeOrderIdStr)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Tạo một PcOrder mới
		pcOrder := &models.PcOrder{
			ID:             primitive.NewObjectID(),
			PancakeOrderId: pancakeOrderIdStr,
			PanCakeData:    input.PanCakeData,
			Status:         0,
			CreatedAt:      time.Now().Unix(),
			UpdatedAt:      time.Now().Unix(),
		}

		// Lưu PcOrder
		createdPcOrder, err := s.BaseServiceImpl.Create(ctx, *pcOrder)
		if err != nil {
			return nil, err
		}

		return &createdPcOrder, nil
	} else {
		// Lấy PcOrder hiện tại
		filter := bson.M{"pancakeOrderId": pancakeOrderIdStr}
		pcOrder, err := s.BaseServiceImpl.FindOneByFilter(ctx, filter, nil)
		if err != nil {
			return nil, err
		}

		// Cập nhật thông tin mới
		pcOrder.PanCakeData = input.PanCakeData
		pcOrder.UpdatedAt = time.Now().Unix()

		// Cập nhật PcOrder
		updatedPcOrder, err := s.BaseServiceImpl.Update(ctx, pcOrder.ID, pcOrder)
		if err != nil {
			return nil, err
		}

		return &updatedPcOrder, nil
	}
}

// FindOneById tìm một PcOrder theo ID
func (s *PcOrderService) FindOneById(ctx context.Context, id primitive.ObjectID) (models.PcOrder, error) {
	return s.BaseServiceImpl.FindOne(ctx, id)
}

// FindAll tìm tất cả các PcOrder với phân trang
func (s *PcOrderService) FindAll(ctx context.Context, page int64, limit int64) ([]models.PcOrder, error) {
	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit)
	return s.BaseServiceImpl.FindAll(ctx, nil, opts)
}
