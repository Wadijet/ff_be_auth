package etl

import (
	"context"
	"fmt"
	"meta_commerce/app/etl/types"
)

// Pipeline implement interface Pipeline để kết nối và điều phối các components
type Pipeline struct {
	source      types.DataSource
	transformer types.Transformer
	destination types.Destination
}

// NewPipeline tạo một instance mới của Pipeline
func NewPipeline(
	source types.DataSource,
	transformer types.Transformer,
	destination types.Destination,
) (*Pipeline, error) {
	// Validate các components
	if source == nil {
		return nil, fmt.Errorf("source không được để trống")
	}
	if transformer == nil {
		return nil, fmt.Errorf("transformer không được để trống")
	}
	if destination == nil {
		return nil, fmt.Errorf("destination không được để trống")
	}

	return &Pipeline{
		source:      source,
		transformer: transformer,
		destination: destination,
	}, nil
}

// Execute thực thi toàn bộ pipeline: fetch -> transform -> store
func (p *Pipeline) Execute(ctx context.Context) error {
	// 1. Fetch dữ liệu từ source
	sourceData, err := p.source.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("lỗi fetch dữ liệu từ source: %v", err)
	}

	// 2. Transform dữ liệu
	transformedData, err := p.transformer.Transform(sourceData)
	if err != nil {
		return fmt.Errorf("lỗi transform dữ liệu: %v", err)
	}

	// 3. Store dữ liệu vào destination
	if err := p.destination.Store(ctx, transformedData); err != nil {
		return fmt.Errorf("lỗi store dữ liệu vào destination: %v", err)
	}

	return nil
}

// GetComponents trả về thông tin cấu hình của tất cả components
func (p *Pipeline) GetComponents() map[string]interface{} {
	return map[string]interface{}{
		"source": map[string]interface{}{
			"type":   fmt.Sprintf("%T", p.source),
			"config": p.source.GetSourceConfig(),
		},
		"transformer": map[string]interface{}{
			"type":   fmt.Sprintf("%T", p.transformer),
			"config": p.transformer.GetTransformConfig(),
		},
		"destination": map[string]interface{}{
			"type":   fmt.Sprintf("%T", p.destination),
			"config": p.destination.GetDestConfig(),
		},
	}
}
