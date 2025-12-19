package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// PcPosProductService là cấu trúc chứa các phương thức liên quan đến Pancake POS Product
type PcPosProductService struct {
	*BaseServiceMongoImpl[models.PcPosProduct]
}

// NewPcPosProductService tạo mới PcPosProductService
func NewPcPosProductService() (*PcPosProductService, error) {
	productCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcPosProducts)
	if !exist {
		return nil, fmt.Errorf("failed to get pc_pos_products collection: %v", common.ErrNotFound)
	}

	return &PcPosProductService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.PcPosProduct](productCollection),
	}, nil
}
