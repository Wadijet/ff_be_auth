package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"meta_commerce/app/etl/types"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// RESTSource implementation của DataSource cho REST API
type RESTSource struct {
	config *types.SourceConfig
	client *http.Client
	state  struct {
		currentPage   int
		currentOffset int
		currentCursor string
		hasMore       bool
	}
}

// NewRESTSource tạo instance mới của RESTSource
func NewRESTSource(config *types.SourceConfig) (*RESTSource, error) {
	if config.Request == nil {
		return nil, fmt.Errorf("thiếu cấu hình Request")
	}

	if config.Type != types.SourceREST {
		return nil, fmt.Errorf("loại source không hợp lệ, cần SourceREST")
	}

	return &RESTSource{
		config: config,
		client: &http.Client{
			Timeout: config.Request.Timeout,
		},
	}, nil
}

// Fetch implements types.DataSource
func (s *RESTSource) Fetch(ctx context.Context) ([]byte, error) {
	// Tạo request
	req, err := s.buildRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("lỗi tạo request: %v", err)
	}

	// Thực hiện request với retry
	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt <= s.config.Request.Retry.MaxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(s.config.Request.Retry.Delay)
		}

		resp, err = s.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}

		lastErr = err
		if resp != nil {
			resp.Body.Close()
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("lỗi sau %d lần retry: %v", s.config.Request.Retry.MaxAttempts, lastErr)
	}
	defer resp.Body.Close()

	// Đọc response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc response: %v", err)
	}

	// Xử lý phân trang
	if s.config.Paging != nil && s.config.Paging.Enabled {
		if err := s.updatePaginationState(body); err != nil {
			return nil, fmt.Errorf("lỗi xử lý phân trang: %v", err)
		}
	}

	return body, nil
}

// buildRequest tạo HTTP request với các tham số cần thiết
func (s *RESTSource) buildRequest(ctx context.Context) (*http.Request, error) {
	// Parse URL và thêm query params
	baseURL, err := url.Parse(s.config.Request.URL)
	if err != nil {
		return nil, err
	}

	query := baseURL.Query()
	// Thêm query params từ config
	for k, v := range s.config.Request.QueryParams {
		query.Set(k, v)
	}

	// Thêm params phân trang nếu có
	if s.config.Paging != nil && s.config.Paging.Enabled {
		switch s.config.Paging.Type {
		case "page":
			query.Set(s.config.Paging.Params.Page, strconv.Itoa(s.state.currentPage))
			query.Set(s.config.Paging.Params.Size, strconv.Itoa(s.config.Paging.DefaultSize))
		case "offset":
			query.Set(s.config.Paging.Params.Offset, strconv.Itoa(s.state.currentOffset))
			query.Set(s.config.Paging.Params.Limit, strconv.Itoa(s.config.Paging.DefaultSize))
		case "cursor":
			if s.state.currentCursor != "" {
				query.Set(s.config.Paging.Params.Cursor, s.state.currentCursor)
			}
		}
	}

	baseURL.RawQuery = query.Encode()

	// Tạo request
	req, err := http.NewRequestWithContext(ctx, s.config.Request.Method, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// Thêm headers
	for k, v := range s.config.Request.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// updatePaginationState cập nhật trạng thái phân trang từ response
func (s *RESTSource) updatePaginationState(response []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(response, &data); err != nil {
		return err
	}

	switch s.config.Paging.Type {
	case "page":
		s.state.currentPage++
		if pages, ok := getValueByPath(data, s.config.Paging.Paths.Pages); ok {
			if p, ok := pages.(float64); ok {
				s.state.hasMore = s.state.currentPage <= int(p)
			}
		}

	case "offset":
		s.state.currentOffset += s.config.Paging.DefaultSize
		if hasMore, ok := getValueByPath(data, s.config.Paging.Paths.HasMore); ok {
			if h, ok := hasMore.(bool); ok {
				s.state.hasMore = h
			}
		}

	case "cursor":
		if cursor, ok := getValueByPath(data, s.config.Paging.Paths.NextCursor); ok {
			if c, ok := cursor.(string); ok {
				s.state.currentCursor = c
				s.state.hasMore = c != ""
			}
		}
	}

	return nil
}

// GetSourceConfig implements types.DataSource
func (s *RESTSource) GetSourceConfig() map[string]interface{} {
	return map[string]interface{}{
		"type":         s.config.Type,
		"url":          s.config.Request.URL,
		"method":       s.config.Request.Method,
		"pagination":   s.config.Paging,
		"has_more":     s.state.hasMore,
		"current_page": s.state.currentPage,
	}
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
