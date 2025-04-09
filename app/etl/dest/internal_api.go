package dest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// InternalAPIConfig chứa cấu hình cho Internal API destination
type InternalAPIConfig struct {
	URL     string            `json:"url" validate:"required"` // URL của internal API
	Method  string            `json:"method"`                  // HTTP method (mặc định là POST)
	Headers map[string]string `json:"headers"`                 // Headers tùy chọn
	Timeout time.Duration     `json:"timeout"`                 // Timeout cho request (mặc định 30s)
}

// InternalAPIDestination implement interface Destination cho việc gửi dữ liệu đến internal API
type InternalAPIDestination struct {
	config     InternalAPIConfig
	httpClient *http.Client
}

// NewInternalAPIDestination tạo một instance mới của InternalAPIDestination
func NewInternalAPIDestination(config InternalAPIConfig) (*InternalAPIDestination, error) {
	// Validate config
	if config.URL == "" {
		return nil, fmt.Errorf("URL không được để trống")
	}

	// Set giá trị mặc định
	if config.Method == "" {
		config.Method = http.MethodPost
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Khởi tạo HTTP client với timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	return &InternalAPIDestination{
		config:     config,
		httpClient: client,
	}, nil
}

// Store implement interface Destination.Store
func (d *InternalAPIDestination) Store(ctx context.Context, data []byte) error {
	// Tạo request với context
	req, err := http.NewRequestWithContext(ctx, d.config.Method, d.config.URL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("lỗi tạo request: %v", err)
	}

	// Set Content-Type header mặc định
	req.Header.Set("Content-Type", "application/json")

	// Thêm custom headers
	for key, value := range d.config.Headers {
		req.Header.Set(key, value)
	}

	// Thực hiện request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("lỗi gửi request: %v", err)
	}
	defer resp.Body.Close()

	// Đọc response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("lỗi đọc response: %v", err)
	}

	// Kiểm tra status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API trả về lỗi (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetDestConfig implement interface Destination.GetDestConfig
func (d *InternalAPIDestination) GetDestConfig() map[string]interface{} {
	return map[string]interface{}{
		"url":     d.config.URL,
		"method":  d.config.Method,
		"headers": d.config.Headers,
		"timeout": d.config.Timeout.String(),
	}
}
