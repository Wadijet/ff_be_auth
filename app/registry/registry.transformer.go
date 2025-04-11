package registry

import (
	"meta_commerce/app/etl/types"
	"meta_commerce/app/utility"
)

// TransformerRegistryImpl quản lý việc đăng ký và tạo các transformer
type TransformerRegistryImpl struct {
	*Registry[func(config interface{}) (types.Transformer, error)]
}

// NewTransformerRegistry tạo một instance mới của TransformerRegistryImpl
func NewTransformerRegistry() *TransformerRegistryImpl {
	return &TransformerRegistryImpl{
		Registry: NewRegistry[func(config interface{}) (types.Transformer, error)](),
	}
}

// Register implements types.TransformerRegistry
func (r *TransformerRegistryImpl) Register(name string, creator func(config interface{}) (types.Transformer, error)) {
	r.Registry.Register(name, creator)
}

// Create tạo một transformer mới từ config
func (r *TransformerRegistryImpl) Create(transformType string, config interface{}) (types.Transformer, error) {
	creator, exists := r.Get(transformType)
	if !exists {
		return nil, utility.NewError(
			utility.ErrCodeValidation,
			"Không tìm thấy transformer type: "+transformType,
			utility.StatusNotFound,
			nil,
		)
	}
	return creator(config)
}
