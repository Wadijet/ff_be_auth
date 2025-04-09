package source

import (
	"context"
	"fmt"
	"io"
	"meta_commerce/app/etl/types"
	"net/http"
	"time"
)

// RestConfig chứa các cấu hình cho REST API source
type RestConfig struct {
	URL     string            `json:"url" validate:"required,url"` // URL của API
	Method  string            `json:"method" validate:"required"`  // HTTP method (GET, POST, etc.)
	Headers map[string]string `json:"headers"`                     // Headers tùy chọn
	Timeout time.Duration     `json:"timeout" validate:"min=1"`    // Timeout cho request (seconds)
}

// RestApiSource implement interface DataSource cho REST API
type RestApiSource struct {
	config RestConfig
	client *http.Client
}

// NewRestApiSource tạo một instance mới của RestApiSource
func NewRestApiSource(config RestConfig) (*RestApiSource, error) {
	// Validate config
	if config.URL == "" {
		return nil, fmt.Errorf("URL không được để trống")
	}
	if config.Method == "" {
		config.Method = "GET" // Default method
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second // Default timeout
	}

	// Tạo HTTP client với timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	return &RestApiSource{
		config: config,
		client: client,
	}, nil
}

// Fetch implement interface DataSource.Fetch
func (s *RestApiSource) Fetch(ctx context.Context) ([]byte, error) {
	// Tạo request với context
	req, err := http.NewRequestWithContext(ctx, s.config.Method, s.config.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("lỗi tạo request: %v", err)
	}

	// Thêm headers
	for key, value := range s.config.Headers {
		req.Header.Set(key, value)
	}

	// Thực hiện request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lỗi thực hiện request: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API trả về status code không hợp lệ: %d", resp.StatusCode)
	}

	// Đọc response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc response body: %v", err)
	}

	return body, nil
}

// GetSourceConfig implement interface DataSource.GetSourceConfig
func (s *RestApiSource) GetSourceConfig() map[string]interface{} {
	return map[string]interface{}{
		"url":     s.config.URL,
		"method":  s.config.Method,
		"headers": s.config.Headers,
		"timeout": s.config.Timeout,
	}
}

// Đảm bảo RestApiSource implement interface DataSource
var _ types.DataSource = (*RestApiSource)(nil)
