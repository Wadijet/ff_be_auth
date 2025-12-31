package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient là wrapper cho http.Client với các tiện ích
type HTTPClient struct {
	client            *http.Client
	baseURL          string
	token            string
	activeRoleID     string // Header X-Active-Role-ID cho organization context
}

// NewHTTPClient tạo mới một HTTP client
func NewHTTPClient(baseURL string, timeout int) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		baseURL: baseURL,
	}
}

// SetToken thiết lập token cho các request tiếp theo
func (c *HTTPClient) SetToken(token string) {
	c.token = token
}

// SetActiveRoleID thiết lập active role ID cho organization context
// Header này sẽ được gửi kèm trong request để xác định organization context
func (c *HTTPClient) SetActiveRoleID(roleID string) {
	c.activeRoleID = roleID
}

// GetToken lấy token hiện tại
func (c *HTTPClient) GetToken() string {
	return c.token
}

// Request thực hiện HTTP request
func (c *HTTPClient) Request(method, path string, body interface{}) (*http.Response, []byte, error) {
	// Tạo URL đầy đủ
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	// Tạo request body nếu có
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("lỗi khi marshal JSON: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Tạo request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("lỗi khi tạo request: %v", err)
	}

	// Thêm headers
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	// Thêm header X-Active-Role-ID cho organization context
	if c.activeRoleID != "" {
		req.Header.Set("X-Active-Role-ID", c.activeRoleID)
	}

	// Thực hiện request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("lỗi khi thực hiện request: %v", err)
	}

	// Đọc response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("lỗi khi đọc response body: %v", err)
	}
	defer resp.Body.Close()

	return resp, respBody, nil
}

// GET request
func (c *HTTPClient) GET(path string) (*http.Response, []byte, error) {
	return c.Request(http.MethodGet, path, nil)
}

// POST request
func (c *HTTPClient) POST(path string, body interface{}) (*http.Response, []byte, error) {
	return c.Request(http.MethodPost, path, body)
}

// PUT request
func (c *HTTPClient) PUT(path string, body interface{}) (*http.Response, []byte, error) {
	return c.Request(http.MethodPut, path, body)
}

// DELETE request
func (c *HTTPClient) DELETE(path string) (*http.Response, []byte, error) {
	return c.Request(http.MethodDelete, path, nil)
}

// CheckResponse kiểm tra response status và trả về error nếu có
func CheckResponse(resp *http.Response, respBody []byte) error {
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
