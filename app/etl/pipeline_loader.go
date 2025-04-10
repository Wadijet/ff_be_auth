package etl

import (
	"meta_commerce/app/etl/dest"
	"meta_commerce/app/etl/source"
	"meta_commerce/app/etl/transform"
	"meta_commerce/app/registry"
	"time"
)

// PipelineBuilder giúp tạo pipeline một cách dễ dàng
type PipelineBuilder struct {
	registry *registry.ETLRegistry
}

// NewPipelineBuilder tạo một instance mới của PipelineBuilder
func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{
		registry: registry.GetETLRegistry(),
	}
}

// BuildAllPipelines tạo tất cả các pipelines được định nghĩa
func (b *PipelineBuilder) BuildAllPipelines() ([]*Pipeline, error) {
	var pipelines []*Pipeline

	// 1. User Sync Pipeline
	userPipeline, err := b.BuildUserSyncPipeline()
	if err != nil {
		return nil, err
	}
	pipelines = append(pipelines, userPipeline)

	// 2. Order Sync Pipeline
	orderPipeline, err := b.BuildOrderSyncPipeline()
	if err != nil {
		return nil, err
	}
	pipelines = append(pipelines, orderPipeline)

	// 3. Product Sync Pipeline
	productPipeline, err := b.BuildProductSyncPipeline()
	if err != nil {
		return nil, err
	}
	pipelines = append(pipelines, productPipeline)

	return pipelines, nil
}

// BuildUserSyncPipeline tạo một pipeline đồng bộ user
func (b *PipelineBuilder) BuildUserSyncPipeline() (*Pipeline, error) {
	// 1. Tạo REST API source
	sourceConfig := source.RestConfig{
		URL: "https://api.example.com/users",
		Headers: map[string]string{
			"Authorization": "Bearer your-token",
		},
		Method:  "GET",
		Timeout: 30 * time.Second,
	}
	src, err := b.registry.CreateSource("rest_api", sourceConfig)
	if err != nil {
		return nil, err
	}

	// 2. Tạo Field Mapper transformer
	transformConfig := transform.MapperConfig{
		Mappings: []transform.FieldMapping{
			{Source: "id", Target: "external_id", Type: "string"},
			{Source: "name", Target: "full_name", Type: "string"},
			{Source: "email", Target: "email", Type: "string"},
		},
	}
	transformer, err := b.registry.CreateTransformer("field_mapper", transformConfig)
	if err != nil {
		return nil, err
	}

	// 3. Tạo Internal API destination
	destConfig := dest.InternalAPIConfig{
		URL: "http://localhost:8080/api/users/sync",
		Headers: map[string]string{
			"X-API-Key":    "your-internal-key",
			"Content-Type": "application/json",
		},
		Method:  "POST",
		Timeout: 30 * time.Second,
	}
	dst, err := b.registry.CreateDestination("internal_api", destConfig)
	if err != nil {
		return nil, err
	}

	// 4. Tạo và trả về pipeline
	return NewPipeline(src, transformer, dst)
}

// BuildOrderSyncPipeline tạo một pipeline đồng bộ order
func (b *PipelineBuilder) BuildOrderSyncPipeline() (*Pipeline, error) {
	// 1. Tạo REST API source
	sourceConfig := source.RestConfig{
		URL: "https://api.example.com/orders",
		Headers: map[string]string{
			"Authorization": "Bearer your-token",
		},
		Method:  "GET",
		Timeout: 30 * time.Second,
	}
	src, err := b.registry.CreateSource("rest_api", sourceConfig)
	if err != nil {
		return nil, err
	}

	// 2. Tạo Field Mapper transformer
	transformConfig := transform.MapperConfig{
		Mappings: []transform.FieldMapping{
			{Source: "id", Target: "external_id", Type: "string"},
			{Source: "order_number", Target: "order_number", Type: "string"},
			{Source: "total_amount", Target: "amount", Type: "float"},
			{Source: "status", Target: "status", Type: "string"},
			{Source: "created_at", Target: "order_date", Type: "datetime"},
		},
	}
	transformer, err := b.registry.CreateTransformer("field_mapper", transformConfig)
	if err != nil {
		return nil, err
	}

	// 3. Tạo Internal API destination
	destConfig := dest.InternalAPIConfig{
		URL: "http://localhost:8080/api/orders/sync",
		Headers: map[string]string{
			"X-API-Key":    "your-internal-key",
			"Content-Type": "application/json",
		},
		Method:  "POST",
		Timeout: 30 * time.Second,
	}
	dst, err := b.registry.CreateDestination("internal_api", destConfig)
	if err != nil {
		return nil, err
	}

	// 4. Tạo và trả về pipeline
	return NewPipeline(src, transformer, dst)
}

// BuildProductSyncPipeline tạo một pipeline đồng bộ product
func (b *PipelineBuilder) BuildProductSyncPipeline() (*Pipeline, error) {
	// 1. Tạo REST API source
	sourceConfig := source.RestConfig{
		URL: "https://api.example.com/products",
		Headers: map[string]string{
			"Authorization": "Bearer your-token",
		},
		Method:  "GET",
		Timeout: 30 * time.Second,
	}
	src, err := b.registry.CreateSource("rest_api", sourceConfig)
	if err != nil {
		return nil, err
	}

	// 2. Tạo Field Mapper transformer
	transformConfig := transform.MapperConfig{
		Mappings: []transform.FieldMapping{
			{Source: "id", Target: "external_id", Type: "string"},
			{Source: "name", Target: "product_name", Type: "string"},
			{Source: "description", Target: "description", Type: "string"},
			{Source: "price", Target: "price", Type: "float"},
			{Source: "stock", Target: "quantity", Type: "integer"},
			{Source: "category", Target: "category", Type: "string"},
		},
	}
	transformer, err := b.registry.CreateTransformer("field_mapper", transformConfig)
	if err != nil {
		return nil, err
	}

	// 3. Tạo Internal API destination
	destConfig := dest.InternalAPIConfig{
		URL: "http://localhost:8080/api/products/sync",
		Headers: map[string]string{
			"X-API-Key":    "your-internal-key",
			"Content-Type": "application/json",
		},
		Method:  "POST",
		Timeout: 30 * time.Second,
	}
	dst, err := b.registry.CreateDestination("internal_api", destConfig)
	if err != nil {
		return nil, err
	}

	// 4. Tạo và trả về pipeline
	return NewPipeline(src, transformer, dst)
}
