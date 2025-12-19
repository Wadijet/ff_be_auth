package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// PcPosShopService là cấu trúc chứa các phương thức liên quan đến Pancake POS Shop
type PcPosShopService struct {
	*BaseServiceMongoImpl[models.PcPosShop]
}

// NewPcPosShopService tạo mới PcPosShopService
func NewPcPosShopService() (*PcPosShopService, error) {
	shopCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcPosShops)
	if !exist {
		return nil, fmt.Errorf("failed to get pc_pos_shops collection: %v", common.ErrNotFound)
	}

	return &PcPosShopService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.PcPosShop](shopCollection),
	}, nil
}
