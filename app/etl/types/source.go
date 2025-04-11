package types

import (
	"time"
)

// SourceType định nghĩa các loại source được hỗ trợ
type SourceType string

const (
	RestAPI SourceType = "rest_api"
	GraphQL SourceType = "graphql"
	GRPC    SourceType = "grpc"
)

// HTTPConfig cấu hình chung cho các HTTP request
type HTTPConfig struct {
	Timeout    time.Duration     // timeout cho mỗi request
	MaxRetries int               // số lần retry tối đa
	RetryDelay time.Duration     // thời gian chờ giữa các lần retry
	Headers    map[string]string // headers cho request
}

// AuthConfig cấu hình xác thực
type AuthConfig struct {
	Type          string // kiểu xác thực: "bearer", "basic", "api_key"
	TokenURL      string // URL để lấy token mới
	ClientID      string
	ClientSecret  string
	RefreshToken  string
	TokenLifetime time.Duration // thời gian sống của token
}

// RESTSourceConfig cấu hình cho REST API source
type RESTSourceConfig struct {
	Type        SourceType        // loại source
	BaseURL     string            // URL cơ sở của API
	Method      string            // phương thức HTTP
	QueryParams map[string]string // query params cố định
	HTTPConfig  *HTTPConfig       // cấu hình HTTP
	AuthConfig  *AuthConfig       // cấu hình xác thực

	// Cấu hình phân trang
	Pagination struct {
		Enabled bool           // bật/tắt phân trang
		Type    PaginationType // kiểu phân trang

		// Cấu hình chung
		DefaultSize int    // kích thước mặc định mỗi trang
		ItemsPath   string // đường dẫn đến mảng dữ liệu trong response

		// Page-based config
		PageConfig *struct {
			PageParam      string // tên tham số page
			SizeParam      string // tên tham số size
			TotalPath      string // đường dẫn đến tổng số bản ghi
			TotalPagesPath string // đường dẫn đến tổng số trang
		}

		// Offset-based config
		OffsetConfig *struct {
			OffsetParam string // tên tham số offset
			LimitParam  string // tên tham số limit
			HasMorePath string // đường dẫn đến trường has_more
		}

		// Cursor-based config
		CursorConfig *struct {
			CursorParam    string // tên tham số cursor
			SizeParam      string // tên tham số size
			NextCursorPath string // đường dẫn đến cursor tiếp theo
			HasMorePath    string // đường dẫn đến trường has_more
		}
	}

	// Cấu hình transform dữ liệu sau khi fetch
	Transform struct {
		Enabled bool // bật/tắt transform
		// Các rules transform
		Rules []struct {
			Type   string                 // loại transform
			Config map[string]interface{} // cấu hình cho transform
		}
	}
}

// NewRESTSourceConfig tạo config mặc định cho REST API source
func NewRESTSourceConfig() *RESTSourceConfig {
	config := &RESTSourceConfig{
		Type:   RestAPI,
		Method: "GET",
		HTTPConfig: &HTTPConfig{
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			RetryDelay: 5 * time.Second,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	// Cấu hình phân trang mặc định (page-based)
	config.Pagination.Enabled = true
	config.Pagination.Type = PageBased
	config.Pagination.DefaultSize = 20
	config.Pagination.ItemsPath = "data"
	config.Pagination.PageConfig = &struct {
		PageParam      string
		SizeParam      string
		TotalPath      string
		TotalPagesPath string
	}{
		PageParam:      "page",
		SizeParam:      "size",
		TotalPath:      "metadata.total",
		TotalPagesPath: "metadata.pages",
	}

	return config
}

// GetPaginationConfig tạo đối tượng PaginationConfig dựa trên cấu hình source
func (c *RESTSourceConfig) GetPaginationConfig() PaginationConfig {
	if !c.Pagination.Enabled {
		return nil
	}

	switch c.Pagination.Type {
	case PageBased:
		if c.Pagination.PageConfig == nil {
			return nil
		}
		return &PageBasedConfig{
			BasePaginationConfig: BasePaginationConfig{
				ItemsPath:   c.Pagination.ItemsPath,
				DefaultSize: c.Pagination.DefaultSize,
			},
			PageParam:      c.Pagination.PageConfig.PageParam,
			SizeParam:      c.Pagination.PageConfig.SizeParam,
			CurrentPage:    1,
			TotalPath:      c.Pagination.PageConfig.TotalPath,
			TotalPagesPath: c.Pagination.PageConfig.TotalPagesPath,
		}

	case OffsetBased:
		if c.Pagination.OffsetConfig == nil {
			return nil
		}
		return &OffsetBasedConfig{
			BasePaginationConfig: BasePaginationConfig{
				ItemsPath:   c.Pagination.ItemsPath,
				DefaultSize: c.Pagination.DefaultSize,
			},
			OffsetParam: c.Pagination.OffsetConfig.OffsetParam,
			LimitParam:  c.Pagination.OffsetConfig.LimitParam,
			HasMorePath: c.Pagination.OffsetConfig.HasMorePath,
		}

	case CursorBased:
		if c.Pagination.CursorConfig == nil {
			return nil
		}
		return &CursorBasedConfig{
			BasePaginationConfig: BasePaginationConfig{
				ItemsPath:   c.Pagination.ItemsPath,
				DefaultSize: c.Pagination.DefaultSize,
			},
			CursorParam:    c.Pagination.CursorConfig.CursorParam,
			SizeParam:      c.Pagination.CursorConfig.SizeParam,
			NextCursorPath: c.Pagination.CursorConfig.NextCursorPath,
			HasMorePath:    c.Pagination.CursorConfig.HasMorePath,
		}
	}

	return nil
}
