package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// PcPosCategoryService là cấu trúc chứa các phương thức liên quan đến Pancake POS Category
type PcPosCategoryService struct {
	*BaseServiceMongoImpl[models.PcPosCategory]
}

// NewPcPosCategoryService tạo mới PcPosCategoryService
func NewPcPosCategoryService() (*PcPosCategoryService, error) {
	categoryCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcPosCategories)
	if !exist {
		return nil, fmt.Errorf("failed to get pc_pos_categories collection: %v", common.ErrNotFound)
	}

	return &PcPosCategoryService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.PcPosCategory](categoryCollection),
	}, nil
}
