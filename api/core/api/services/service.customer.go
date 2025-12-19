package services

import (
	"fmt"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
)

// CustomerService là cấu trúc chứa các phương thức liên quan đến Customer
type CustomerService struct {
	*BaseServiceMongoImpl[models.Customer]
}

// NewCustomerService tạo mới CustomerService
func NewCustomerService() (*CustomerService, error) {
	customerCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.Customers)
	if !exist {
		return nil, fmt.Errorf("failed to get customers collection: %v", common.ErrNotFound)
	}

	return &CustomerService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Customer](customerCollection),
	}, nil
}
