package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// FbCustomerService là cấu trúc chứa các phương thức liên quan đến Facebook Customer
type FbCustomerService struct {
	*BaseServiceMongoImpl[models.FbCustomer]
}

// NewFbCustomerService tạo mới FbCustomerService
func NewFbCustomerService() (*FbCustomerService, error) {
	fbCustomerCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.FbCustomers)
	if !exist {
		return nil, fmt.Errorf("failed to get fb_customers collection: %v", common.ErrNotFound)
	}

	return &FbCustomerService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.FbCustomer](fbCustomerCollection),
	}, nil
}
