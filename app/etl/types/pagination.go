package types

import (
	"fmt"
	"strings"
)

// PaginationType định nghĩa các loại phân trang được hỗ trợ
type PaginationType string

const (
	PageBased   PaginationType = "page"
	OffsetBased PaginationType = "offset"
	CursorBased PaginationType = "cursor"
)

// PaginationConfig interface cơ sở cho tất cả loại phân trang
type PaginationConfig interface {
	// BuildQueryParams tạo query params cho request tiếp theo
	BuildQueryParams() map[string]string
	// ParseResponse xử lý response và cập nhật trạng thái phân trang
	ParseResponse(response map[string]interface{}) error
	// HasNext kiểm tra còn dữ liệu tiếp theo không
	HasNext() bool
}

// BasePaginationConfig chứa các thông tin cơ bản cho phân trang
type BasePaginationConfig struct {
	ItemsPath   string // đường dẫn đến mảng dữ liệu trong response
	DefaultSize int    // kích thước mặc định mỗi trang
	HasMore     bool   // còn dữ liệu tiếp không
}

// PageBasedConfig cấu hình cho phân trang theo số trang
type PageBasedConfig struct {
	BasePaginationConfig
	PageParam      string // tên tham số page trong query
	SizeParam      string // tên tham số size trong query
	CurrentPage    int    // trang hiện tại
	TotalPages     int    // tổng số trang
	TotalPath      string // đường dẫn đến tổng số bản ghi trong response
	TotalPagesPath string // đường dẫn đến tổng số trang trong response
}

// NewPageBasedConfig tạo config mặc định cho page-based pagination
func NewPageBasedConfig() *PageBasedConfig {
	return &PageBasedConfig{
		BasePaginationConfig: BasePaginationConfig{
			ItemsPath:   "data",
			DefaultSize: 20,
		},
		PageParam:      "page",
		SizeParam:      "size",
		CurrentPage:    1,
		TotalPath:      "metadata.total",
		TotalPagesPath: "metadata.pages",
	}
}

// BuildQueryParams implements PaginationConfig
func (c *PageBasedConfig) BuildQueryParams() map[string]string {
	return map[string]string{
		c.PageParam: fmt.Sprintf("%d", c.CurrentPage),
		c.SizeParam: fmt.Sprintf("%d", c.DefaultSize),
	}
}

// ParseResponse implements PaginationConfig
func (c *PageBasedConfig) ParseResponse(response map[string]interface{}) error {
	// Parse tổng số trang từ response
	if pages, ok := getValueByPath(response, c.TotalPagesPath); ok {
		if p, ok := pages.(float64); ok {
			c.TotalPages = int(p)
		}
	}

	c.HasMore = c.CurrentPage < c.TotalPages
	return nil
}

// HasNext implements PaginationConfig
func (c *PageBasedConfig) HasNext() bool {
	return c.HasMore
}

// OffsetBasedConfig cấu hình cho phân trang theo offset
type OffsetBasedConfig struct {
	BasePaginationConfig
	OffsetParam   string // tên tham số offset trong query
	LimitParam    string // tên tham số limit trong query
	CurrentOffset int    // offset hiện tại
	HasMorePath   string // đường dẫn đến trường has_more trong response
}

// NewOffsetBasedConfig tạo config mặc định cho offset-based pagination
func NewOffsetBasedConfig() *OffsetBasedConfig {
	return &OffsetBasedConfig{
		BasePaginationConfig: BasePaginationConfig{
			ItemsPath:   "data",
			DefaultSize: 20,
		},
		OffsetParam: "offset",
		LimitParam:  "limit",
		HasMorePath: "has_more",
	}
}

// BuildQueryParams implements PaginationConfig
func (c *OffsetBasedConfig) BuildQueryParams() map[string]string {
	return map[string]string{
		c.OffsetParam: fmt.Sprintf("%d", c.CurrentOffset),
		c.LimitParam:  fmt.Sprintf("%d", c.DefaultSize),
	}
}

// ParseResponse implements PaginationConfig
func (c *OffsetBasedConfig) ParseResponse(response map[string]interface{}) error {
	if hasMore, ok := getValueByPath(response, c.HasMorePath); ok {
		if h, ok := hasMore.(bool); ok {
			c.HasMore = h
		}
	}

	// Tăng offset cho lần sau
	c.CurrentOffset += c.DefaultSize
	return nil
}

// HasNext implements PaginationConfig
func (c *OffsetBasedConfig) HasNext() bool {
	return c.HasMore
}

// CursorBasedConfig cấu hình cho phân trang theo cursor
type CursorBasedConfig struct {
	BasePaginationConfig
	CursorParam    string // tên tham số cursor trong query
	SizeParam      string // tên tham số size trong query
	CurrentCursor  string // cursor hiện tại
	NextCursorPath string // đường dẫn đến cursor tiếp theo trong response
	HasMorePath    string // đường dẫn đến trường has_more trong response
}

// NewCursorBasedConfig tạo config mặc định cho cursor-based pagination
func NewCursorBasedConfig() *CursorBasedConfig {
	return &CursorBasedConfig{
		BasePaginationConfig: BasePaginationConfig{
			ItemsPath:   "data",
			DefaultSize: 20,
		},
		CursorParam:    "cursor",
		SizeParam:      "size",
		NextCursorPath: "next_cursor",
		HasMorePath:    "has_next_page",
	}
}

// BuildQueryParams implements PaginationConfig
func (c *CursorBasedConfig) BuildQueryParams() map[string]string {
	params := map[string]string{
		c.SizeParam: fmt.Sprintf("%d", c.DefaultSize),
	}
	if c.CurrentCursor != "" {
		params[c.CursorParam] = c.CurrentCursor
	}
	return params
}

// ParseResponse implements PaginationConfig
func (c *CursorBasedConfig) ParseResponse(response map[string]interface{}) error {
	// Lấy cursor tiếp theo
	if nextCursor, ok := getValueByPath(response, c.NextCursorPath); ok {
		if cursor, ok := nextCursor.(string); ok {
			c.CurrentCursor = cursor
		}
	}

	// Kiểm tra còn dữ liệu không
	if hasMore, ok := getValueByPath(response, c.HasMorePath); ok {
		if h, ok := hasMore.(bool); ok {
			c.HasMore = h
		}
	}

	return nil
}

// HasNext implements PaginationConfig
func (c *CursorBasedConfig) HasNext() bool {
	return c.HasMore && c.CurrentCursor != ""
}

// getValueByPath lấy giá trị từ map theo đường dẫn (vd: "metadata.total")
func getValueByPath(data map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			val, ok := current[part]
			return val, ok
		}

		next, ok := current[part]
		if !ok {
			return nil, false
		}

		if nextMap, ok := next.(map[string]interface{}); ok {
			current = nextMap
		} else {
			return nil, false
		}
	}

	return nil, false
}
