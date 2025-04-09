package registry

import (
	"fmt"
	"meta_commerce/app/etl/dest"
	"meta_commerce/app/etl/source"
	"meta_commerce/app/etl/transform"
	"meta_commerce/app/etl/types"
	"sync"
)

// ETLRegistry quản lý các components của ETL system
type ETLRegistry struct {
	sources      map[string]func(config interface{}) (types.DataSource, error)
	transformers map[string]func(config interface{}) (types.Transformer, error)
	destinations map[string]func(config interface{}) (types.Destination, error)
	mu           sync.RWMutex
}

var (
	etlRegistry     *ETLRegistry
	etlRegistryOnce sync.Once
)

// GetETLRegistry trả về singleton instance của ETLRegistry
func GetETLRegistry() *ETLRegistry {
	etlRegistryOnce.Do(func() {
		etlRegistry = &ETLRegistry{
			sources:      make(map[string]func(config interface{}) (types.DataSource, error)),
			transformers: make(map[string]func(config interface{}) (types.Transformer, error)),
			destinations: make(map[string]func(config interface{}) (types.Destination, error)),
		}

		// Đăng ký các default components
		etlRegistry.registerDefaultComponents()
	})
	return etlRegistry
}

// registerDefaultComponents đăng ký các components mặc định
func (r *ETLRegistry) registerDefaultComponents() {
	// Đăng ký REST API source
	r.RegisterSource("rest_api", func(config interface{}) (types.DataSource, error) {
		if cfg, ok := config.(source.RestConfig); ok {
			return source.NewRestApiSource(cfg)
		}
		return nil, fmt.Errorf("invalid config type for rest_api source")
	})

	// Đăng ký Field Mapper transformer
	r.RegisterTransformer("field_mapper", func(config interface{}) (types.Transformer, error) {
		if cfg, ok := config.(transform.MapperConfig); ok {
			return transform.NewFieldMapper(cfg)
		}
		return nil, fmt.Errorf("invalid config type for field_mapper transformer")
	})

	// Đăng ký Internal API destination
	r.RegisterDestination("internal_api", func(config interface{}) (types.Destination, error) {
		if cfg, ok := config.(dest.InternalAPIConfig); ok {
			return dest.NewInternalAPIDestination(cfg)
		}
		return nil, fmt.Errorf("invalid config type for internal_api destination")
	})
}

// RegisterSource đăng ký một source mới
func (r *ETLRegistry) RegisterSource(name string, creator func(config interface{}) (types.DataSource, error)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sources[name] = creator
}

// RegisterTransformer đăng ký một transformer mới
func (r *ETLRegistry) RegisterTransformer(name string, creator func(config interface{}) (types.Transformer, error)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transformers[name] = creator
}

// RegisterDestination đăng ký một destination mới
func (r *ETLRegistry) RegisterDestination(name string, creator func(config interface{}) (types.Destination, error)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.destinations[name] = creator
}

// CreateSource tạo một instance của source từ type và config
func (r *ETLRegistry) CreateSource(sourceType string, config interface{}) (types.DataSource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	creator, exists := r.sources[sourceType]
	if !exists {
		return nil, fmt.Errorf("unknown source type: %s", sourceType)
	}
	return creator(config)
}

// CreateTransformer tạo một instance của transformer từ type và config
func (r *ETLRegistry) CreateTransformer(transformerType string, config interface{}) (types.Transformer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	creator, exists := r.transformers[transformerType]
	if !exists {
		return nil, fmt.Errorf("unknown transformer type: %s", transformerType)
	}
	return creator(config)
}

// CreateDestination tạo một instance của destination từ type và config
func (r *ETLRegistry) CreateDestination(destType string, config interface{}) (types.Destination, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	creator, exists := r.destinations[destType]
	if !exists {
		return nil, fmt.Errorf("unknown destination type: %s", destType)
	}
	return creator(config)
}
