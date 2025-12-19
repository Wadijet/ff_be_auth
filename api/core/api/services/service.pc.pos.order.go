package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// PcPosOrderService là cấu trúc chứa các phương thức liên quan đến Pancake POS Order
type PcPosOrderService struct {
	*BaseServiceMongoImpl[models.PcPosOrder]
}

// NewPcPosOrderService tạo mới PcPosOrderService
func NewPcPosOrderService() (*PcPosOrderService, error) {
	orderCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcPosOrders)
	if !exist {
		return nil, fmt.Errorf("failed to get pc_pos_orders collection: %v", common.ErrNotFound)
	}

	return &PcPosOrderService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.PcPosOrder](orderCollection),
	}, nil
}
