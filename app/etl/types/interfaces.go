package types

import (
	"context"
)

// DataSource định nghĩa interface cho việc lấy dữ liệu từ nguồn
// Implementations: REST API, Database, File System, etc.
type DataSource interface {
	// Fetch lấy dữ liệu từ nguồn
	// Context được sử dụng để timeout và cancel
	// Trả về dữ liệu dạng []byte và error nếu có
	Fetch(ctx context.Context) ([]byte, error)

	// GetSourceConfig trả về cấu hình của nguồn dữ liệu
	GetSourceConfig() map[string]interface{}
}

// Transformer định nghĩa interface cho việc biến đổi dữ liệu
// Implementations: Field Mapper, Data Converter, etc.
type Transformer interface {
	// Transform biến đổi dữ liệu từ dạng này sang dạng khác
	// Input và output dạng []byte để linh hoạt với các kiểu dữ liệu
	Transform(data []byte) ([]byte, error)

	// GetTransformConfig trả về cấu hình của transformer
	GetTransformConfig() map[string]interface{}
}

// Destination định nghĩa interface cho việc lưu trữ dữ liệu
// Implementations: API, Database, File System, etc.
type Destination interface {
	// Store lưu trữ dữ liệu đã được transform
	// Context được sử dụng để timeout và cancel
	Store(ctx context.Context, data []byte) error

	// GetDestConfig trả về cấu hình của destination
	GetDestConfig() map[string]interface{}
}

// Pipeline định nghĩa interface cho một luồng xử lý ETL hoàn chỉnh
type Pipeline interface {
	// Execute thực thi toàn bộ quá trình ETL
	// Context được sử dụng để timeout và cancel
	Execute(ctx context.Context) error

	// GetComponents trả về các thành phần của pipeline
	GetComponents() (DataSource, Transformer, Destination)
}
