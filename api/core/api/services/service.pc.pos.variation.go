package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// PcPosVariationService là cấu trúc chứa các phương thức liên quan đến Pancake POS Variation
type PcPosVariationService struct {
	*BaseServiceMongoImpl[models.PcPosVariation]
}

// NewPcPosVariationService tạo mới PcPosVariationService
func NewPcPosVariationService() (*PcPosVariationService, error) {
	variationCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcPosVariations)
	if !exist {
		return nil, fmt.Errorf("failed to get pc_pos_variations collection: %v", common.ErrNotFound)
	}

	return &PcPosVariationService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.PcPosVariation](variationCollection),
	}, nil
}
