package transform

import (
	"encoding/json"
	"fmt"
	"meta_commerce/app/etl/types"
)

// FieldMapping định nghĩa mapping từ field nguồn sang field đích
type FieldMapping struct {
	Source string `json:"source"`         // Tên field trong dữ liệu nguồn
	Target string `json:"target"`         // Tên field trong dữ liệu đích
	Type   string `json:"type,omitempty"` // Kiểu dữ liệu (string, number, boolean, etc.)
}

// MapperConfig chứa cấu hình cho field mapper
type MapperConfig struct {
	Mappings []FieldMapping `json:"mappings" validate:"required,min=1"` // Danh sách các mapping
}

// FieldMapper implement interface Transformer cho việc map fields
type FieldMapper struct {
	config MapperConfig
}

// NewFieldMapper tạo một instance mới của FieldMapper
func NewFieldMapper(config MapperConfig) (*FieldMapper, error) {
	// Validate config
	if len(config.Mappings) == 0 {
		return nil, fmt.Errorf("cần ít nhất một mapping")
	}

	// Validate từng mapping
	for _, m := range config.Mappings {
		if m.Source == "" {
			return nil, fmt.Errorf("source field không được để trống")
		}
		if m.Target == "" {
			return nil, fmt.Errorf("target field không được để trống")
		}
	}

	return &FieldMapper{
		config: config,
	}, nil
}

// Transform implement interface Transformer.Transform
func (m *FieldMapper) Transform(data []byte) ([]byte, error) {
	// Parse input JSON
	var sourceData map[string]interface{}
	if err := json.Unmarshal(data, &sourceData); err != nil {
		return nil, fmt.Errorf("lỗi parse dữ liệu nguồn: %v", err)
	}

	// Tạo map kết quả
	result := make(map[string]interface{})

	// Áp dụng mappings
	for _, mapping := range m.config.Mappings {
		// Lấy giá trị từ source
		value, exists := sourceData[mapping.Source]
		if !exists {
			// Skip nếu field không tồn tại trong source
			continue
		}

		// Kiểm tra và chuyển đổi kiểu dữ liệu nếu cần
		if mapping.Type != "" {
			var err error
			value, err = m.convertType(value, mapping.Type)
			if err != nil {
				return nil, fmt.Errorf("lỗi chuyển đổi kiểu dữ liệu cho field %s: %v", mapping.Source, err)
			}
		}

		// Gán vào target
		result[mapping.Target] = value
	}

	// Convert kết quả về JSON
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("lỗi chuyển đổi kết quả sang JSON: %v", err)
	}

	return resultJSON, nil
}

// convertType chuyển đổi giá trị sang kiểu dữ liệu mong muốn
func (m *FieldMapper) convertType(value interface{}, targetType string) (interface{}, error) {
	switch targetType {
	case "string":
		return fmt.Sprintf("%v", value), nil
	case "number":
		// Thử chuyển về float64
		switch v := value.(type) {
		case float64:
			return v, nil
		case string:
			var f float64
			if err := json.Unmarshal([]byte(v), &f); err != nil {
				return nil, fmt.Errorf("không thể chuyển đổi sang number: %v", err)
			}
			return f, nil
		default:
			return nil, fmt.Errorf("kiểu dữ liệu không hỗ trợ chuyển đổi sang number")
		}
	case "boolean":
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			var b bool
			if err := json.Unmarshal([]byte(v), &b); err != nil {
				return nil, fmt.Errorf("không thể chuyển đổi sang boolean: %v", err)
			}
			return b, nil
		default:
			return nil, fmt.Errorf("kiểu dữ liệu không hỗ trợ chuyển đổi sang boolean")
		}
	default:
		return value, nil // Giữ nguyên giá trị nếu không có yêu cầu chuyển đổi
	}
}

// GetTransformConfig implement interface Transformer.GetTransformConfig
func (m *FieldMapper) GetTransformConfig() map[string]interface{} {
	return map[string]interface{}{
		"mappings": m.config.Mappings,
	}
}

// Đảm bảo FieldMapper implement interface Transformer
var _ types.Transformer = (*FieldMapper)(nil)
