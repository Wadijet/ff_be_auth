package registry

import (
	"meta_commerce/app/etl/types"
	"meta_commerce/app/utility"
)

// DestinationRegistryImpl quản lý việc đăng ký và tạo các destination
type DestinationRegistryImpl struct {
	*Registry[func(config interface{}) (types.Destination, error)]
}

// NewDestinationRegistry tạo một instance mới của DestinationRegistryImpl
func NewDestinationRegistry() *DestinationRegistryImpl {
	return &DestinationRegistryImpl{
		Registry: NewRegistry[func(config interface{}) (types.Destination, error)](),
	}
}

// Register implements types.DestinationRegistry
func (r *DestinationRegistryImpl) Register(name string, creator func(config interface{}) (types.Destination, error)) {
	r.Registry.Register(name, creator)
}

// Create tạo một destination mới từ config
func (r *DestinationRegistryImpl) Create(destType string, config interface{}) (types.Destination, error) {
	creator, exists := r.Get(destType)
	if !exists {
		return nil, utility.NewError(
			utility.ErrCodeValidation,
			"Không tìm thấy destination type: "+destType,
			utility.StatusNotFound,
			nil,
		)
	}
	return creator(config)
}
