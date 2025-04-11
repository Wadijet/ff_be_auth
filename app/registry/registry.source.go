package registry

import (
	"meta_commerce/app/etl/types"
	"meta_commerce/app/utility"
)

// SourceRegistryImpl quản lý việc đăng ký và tạo các source
type SourceRegistryImpl struct {
	*Registry[func(config interface{}) (types.DataSource, error)]
}

// NewSourceRegistry tạo một instance mới của SourceRegistryImpl
func NewSourceRegistry() *SourceRegistryImpl {
	return &SourceRegistryImpl{
		Registry: NewRegistry[func(config interface{}) (types.DataSource, error)](),
	}
}

// Register implements types.SourceRegistry
func (r *SourceRegistryImpl) Register(name string, creator func(config interface{}) (types.DataSource, error)) {
	r.Registry.Register(name, creator)
}

// Create tạo một source mới từ config
func (r *SourceRegistryImpl) Create(sourceType string, config interface{}) (types.DataSource, error) {
	creator, exists := r.Get(sourceType)
	if !exists {
		return nil, utility.NewError(
			utility.ErrCodeValidation,
			"Không tìm thấy source type: "+sourceType,
			utility.StatusNotFound,
			nil,
		)
	}
	return creator(config)
}
