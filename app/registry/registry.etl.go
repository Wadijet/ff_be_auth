package registry

import (
	"meta_commerce/app/etl/dest"
	"meta_commerce/app/etl/source"
	"meta_commerce/app/etl/transform"
	"meta_commerce/app/etl/types"
	"meta_commerce/app/utility"
)

// Định nghĩa các registry cho từng loại component
var (
	Sources      = NewRegistry[func(config interface{}) (types.DataSource, error)]()
	Transformers = NewRegistry[func(config interface{}) (types.Transformer, error)]()
	Destinations = NewRegistry[func(config interface{}) (types.Destination, error)]()
)

// InitETLRegistry khởi tạo và đăng ký các ETL components mặc định
func InitETLRegistry() error {
	// Đăng ký REST API source
	if err := Sources.Register("rest_api", func(config interface{}) (types.DataSource, error) {
		if cfg, ok := config.(*types.SourceConfig); ok {
			return source.NewRESTSource(cfg)
		}
		return nil, utility.NewError(
			utility.ErrCodeValidationInput,
			"Config không hợp lệ cho REST API source",
			utility.StatusBadRequest,
			nil,
		)
	}); err != nil {
		return err
	}

	// Đăng ký Field Mapper transformer
	if err := Transformers.Register("field_mapper", func(config interface{}) (types.Transformer, error) {
		if cfg, ok := config.(transform.MapperConfig); ok {
			return transform.NewFieldMapper(cfg)
		}
		return nil, utility.NewError(
			utility.ErrCodeValidationInput,
			"Config không hợp lệ cho Field Mapper transformer",
			utility.StatusBadRequest,
			nil,
		)
	}); err != nil {
		return err
	}

	// Đăng ký Internal API destination
	if err := Destinations.Register("internal_api", func(config interface{}) (types.Destination, error) {
		if cfg, ok := config.(dest.InternalAPIConfig); ok {
			return dest.NewInternalAPIDestination(cfg)
		}
		return nil, utility.NewError(
			utility.ErrCodeValidationInput,
			"Config không hợp lệ cho Internal API destination",
			utility.StatusBadRequest,
			nil,
		)
	}); err != nil {
		return err
	}

	return nil
}
