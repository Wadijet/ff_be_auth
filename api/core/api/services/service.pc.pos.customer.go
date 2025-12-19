package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// PcPosCustomerService là cấu trúc chứa các phương thức liên quan đến Pancake POS Customer
type PcPosCustomerService struct {
	*BaseServiceMongoImpl[models.PcPosCustomer]
}

// NewPcPosCustomerService tạo mới PcPosCustomerService
func NewPcPosCustomerService() (*PcPosCustomerService, error) {
	pcPosCustomerCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcPosCustomers)
	if !exist {
		return nil, fmt.Errorf("failed to get pc_pos_customers collection: %v", common.ErrNotFound)
	}

	return &PcPosCustomerService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.PcPosCustomer](pcPosCustomerCollection),
	}, nil
}
