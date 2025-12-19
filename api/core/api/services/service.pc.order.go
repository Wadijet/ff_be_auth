package services

import (
	"context"
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PcOrderService là cấu trúc chứa các phương thức liên quan đến đơn hàng
type PcOrderService struct {
	*BaseServiceMongoImpl[models.PcOrder]
}

// NewPcOrderService tạo mới PcOrderService
func NewPcOrderService() (*PcOrderService, error) {
	orderCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcOrders)
	if !exist {
		return nil, fmt.Errorf("failed to get pc_orders collection: %v", common.ErrNotFound)
	}

	return &PcOrderService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.PcOrder](orderCollection),
	}, nil
}

// IsPancakeOrderIdExist kiểm tra ID đơn hàng Pancake có tồn tại hay không
func (s *PcOrderService) IsPancakeOrderIdExist(ctx context.Context, pancakeOrderId string) (bool, error) {
	filter := bson.M{"pancakeOrderId": pancakeOrderId}
	var pcOrder models.PcOrder
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&pcOrder)
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
	err := s.BaseServiceMongoImpl.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return models.PcOrder{}, err
	}
	return result, nil
}

// Delete xóa một document theo ObjectId
func (s *PcOrderService) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := s.BaseServiceMongoImpl.collection.DeleteOne(ctx, filter)
	return err
}

// Update cập nhật một document theo ObjectId
func (s *PcOrderService) Update(ctx context.Context, id primitive.ObjectID, pcOrder models.PcOrder) (models.PcOrder, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": pcOrder}

	_, err := s.BaseServiceMongoImpl.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.PcOrder{}, err
	}

	return s.FindOne(ctx, id)
}
