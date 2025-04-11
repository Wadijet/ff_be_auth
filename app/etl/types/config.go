package types

import "time"

// ComponentType định nghĩa loại component trong ETL
type ComponentType string

const (
	SourceREST     ComponentType = "rest_api"
	SourceGraphQL  ComponentType = "graphql"
	SourceGRPC     ComponentType = "grpc"
	SourceDatabase ComponentType = "database"
	SourceFile     ComponentType = "file"
)

// BaseConfig cấu hình cơ bản cho tất cả components
type BaseConfig struct {
	Type        ComponentType          // Loại component
	Name        string                 // Tên định danh
	Description string                 // Mô tả
	Metadata    map[string]interface{} // Metadata tùy chỉnh
}

// RequestConfig cấu hình chung cho HTTP/gRPC request
type RequestConfig struct {
	URL         string            // URL endpoint
	Method      string            // HTTP method
	Headers     map[string]string // Headers
	QueryParams map[string]string // Query parameters
	Timeout     time.Duration     // Request timeout
	Retry       struct {
		MaxAttempts int           // Số lần retry tối đa
		Delay       time.Duration // Thời gian giữa các lần retry
		Conditions  []string      // Điều kiện để retry
	}
}

// SecurityConfig cấu hình xác thực và bảo mật
type SecurityConfig struct {
	Type          string // bearer, basic, api_key
	TokenURL      string // URL lấy token
	ClientID      string
	ClientSecret  string
	RefreshToken  string
	TokenLifetime time.Duration
}

// PagingConfig cấu hình phân trang
type PagingConfig struct {
	Enabled     bool   // Bật/tắt phân trang
	Type        string // page, offset, cursor
	DefaultSize int    // Kích thước mặc định
	ItemsPath   string // Đường dẫn đến data trong response

	// Các tham số phân trang
	Params struct {
		// Page-based
		Page string // Tên tham số page
		Size string // Tên tham số size

		// Offset-based
		Offset string // Tên tham số offset
		Limit  string // Tên tham số limit

		// Cursor-based
		Cursor string // Tên tham số cursor
	}

	// Đường dẫn trong response
	Paths struct {
		Total      string // Tổng số records
		Pages      string // Tổng số trang
		NextCursor string // Cursor tiếp theo
		HasMore    string // Còn dữ liệu không
	}
}

// DataTransformConfig cấu hình transform
type DataTransformConfig struct {
	Rules []struct {
		Type   string                 // Loại transform
		Config map[string]interface{} // Cấu hình chi tiết
	}
}

// SourceConfig cấu hình chung cho data source
type SourceConfig struct {
	BaseConfig
	Request   *RequestConfig       // Cấu hình request
	Security  *SecurityConfig      // Cấu hình bảo mật
	Paging    *PagingConfig        // Cấu hình phân trang
	Transform *DataTransformConfig // Cấu hình transform
}

// NewSourceConfig tạo config mặc định cho source
func NewSourceConfig(sourceType ComponentType) *SourceConfig {
	config := &SourceConfig{
		BaseConfig: BaseConfig{
			Type: sourceType,
		},
		Request: &RequestConfig{
			Method:  "GET",
			Timeout: 30 * time.Second,
			Retry: struct {
				MaxAttempts int
				Delay       time.Duration
				Conditions  []string
			}{
				MaxAttempts: 3,
				Delay:       5 * time.Second,
			},
		},
		Paging: &PagingConfig{
			Enabled:     true,
			Type:        "page",
			DefaultSize: 20,
			ItemsPath:   "data",
			Params: struct {
				Page   string
				Size   string
				Offset string
				Limit  string
				Cursor string
			}{
				Page: "page",
				Size: "size",
			},
			Paths: struct {
				Total      string
				Pages      string
				NextCursor string
				HasMore    string
			}{
				Total: "metadata.total",
				Pages: "metadata.pages",
			},
		},
	}
	return config
}
