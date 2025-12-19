package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// PcPosWarehouseService là cấu trúc chứa các phương thức liên quan đến Pancake POS Warehouse
type PcPosWarehouseService struct {
	*BaseServiceMongoImpl[models.PcPosWarehouse]
}

// NewPcPosWarehouseService tạo mới PcPosWarehouseService
func NewPcPosWarehouseService() (*PcPosWarehouseService, error) {
	warehouseCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcPosWarehouses)
	if !exist {
		return nil, fmt.Errorf("failed to get pc_pos_warehouses collection: %v", common.ErrNotFound)
	}

	return &PcPosWarehouseService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.PcPosWarehouse](warehouseCollection),
	}, nil
}
