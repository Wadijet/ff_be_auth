package services

import (
	"context"
	"fmt"
	"meta_commerce/core/common"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RelationshipDefinition định nghĩa một quan hệ từ struct tag
type RelationshipDefinition struct {
	// CollectionName: Tên collection cần kiểm tra
	CollectionName string
	// FieldName: Tên field trong collection đó trỏ tới record hiện tại
	FieldName string
	// ErrorMessage: Thông báo lỗi khi tìm thấy quan hệ (có thể dùng %d để thay thế số lượng)
	ErrorMessage string
	// Optional: Nếu true, sẽ không trả về lỗi nếu không tìm thấy collection
	Optional bool
	// Cascade: Nếu true, cho phép xóa cascade (xóa các record con trước)
	Cascade bool
}

// ParseRelationshipTag phân tích struct tag `relationship` để lấy các định nghĩa quan hệ
//
// Format: relationship:"collection:user_roles,field:roleId,message:Không thể xóa role vì có %d user đang sử dụng"
// Hoặc nhiều quan hệ: relationship:"collection:user_roles,field:roleId,message:...|collection:role_permissions,field:roleId,message:..."
//
// Có thể đặt tag trên:
// 1. Field ẩn `_Relationships` (khuyến nghị)
// 2. Bất kỳ field nào trong struct
//
// Parameters:
//   - structType: Type của struct cần phân tích
//
// Returns:
//   - []RelationshipDefinition: Danh sách các quan hệ được định nghĩa
func ParseRelationshipTag(structType reflect.Type) []RelationshipDefinition {
	var relationships []RelationshipDefinition

	// Kiểm tra field ẩn `_Relationships` (khuyến nghị)
	if field, ok := structType.FieldByName("_Relationships"); ok {
		if tag := field.Tag.Get("relationship"); tag != "" {
			relationships = append(relationships, parseRelationshipTagValue(tag)...)
		}
	}

	// Kiểm tra các field khác trong struct (fallback)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		// Bỏ qua field _Relationships vì đã xử lý ở trên
		if field.Name == "_Relationships" {
			continue
		}
		tag := field.Tag.Get("relationship")
		if tag == "" {
			continue
		}

		relationships = append(relationships, parseRelationshipTagValue(tag)...)
	}

	return relationships
}

// parseRelationshipTagValue phân tích giá trị của relationship tag
// Hỗ trợ nhiều quan hệ phân tách bởi dấu |
func parseRelationshipTagValue(tagValue string) []RelationshipDefinition {
	var relationships []RelationshipDefinition

	// Tách các quan hệ (phân tách bởi |)
	parts := strings.Split(tagValue, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		rel := RelationshipDefinition{}

		// Parse các key-value pairs
		pairs := strings.Split(part, ",")
		for _, pair := range pairs {
			pair = strings.TrimSpace(pair)
			if pair == "" {
				continue
			}

			kv := strings.SplitN(pair, ":", 2)
			if len(kv) != 2 {
				continue
			}

			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			switch key {
			case "collection":
				rel.CollectionName = value
			case "field":
				rel.FieldName = value
			case "message", "msg":
				rel.ErrorMessage = value
			case "optional":
				rel.Optional = value == "true" || value == "1"
			case "cascade":
				rel.Cascade = value == "true" || value == "1"
			}
		}

		// Chỉ thêm nếu có đủ thông tin cần thiết
		if rel.CollectionName != "" && rel.FieldName != "" {
			// Nếu không có message, tạo message mặc định
			if rel.ErrorMessage == "" {
				rel.ErrorMessage = fmt.Sprintf("Không thể xóa record vì có %%d record trong collection '%s' đang tham chiếu tới.", rel.CollectionName)
			}
			relationships = append(relationships, rel)
		}
	}

	return relationships
}

// ValidateRelationships kiểm tra các quan hệ được định nghĩa trong struct tag
//
// Parameters:
//   - ctx: Context
//   - recordID: ID của record cần kiểm tra
//   - structType: Type của struct
//
// Returns:
//   - error: Lỗi nếu tìm thấy quan hệ hoặc lỗi khác
func ValidateRelationships(ctx context.Context, recordID primitive.ObjectID, structType reflect.Type) error {
	relationships := ParseRelationshipTag(structType)
	if len(relationships) == 0 {
		return nil
	}

	// Chuyển đổi sang RelationshipCheck để sử dụng hàm có sẵn
	checks := make([]RelationshipCheck, 0, len(relationships))
	for _, rel := range relationships {
		// Bỏ qua nếu có cascade flag (sẽ xóa cascade)
		if rel.Cascade {
			continue
		}

		checks = append(checks, RelationshipCheck{
			CollectionName: rel.CollectionName,
			FieldName:      rel.FieldName,
			ErrorMessage:   rel.ErrorMessage,
			Optional:       rel.Optional,
		})
	}

	if len(checks) > 0 {
		return CheckRelationshipExists(ctx, recordID, checks)
	}

	return nil
}

// ValidateRelationshipsFromValue kiểm tra quan hệ từ một giá trị struct
//
// Parameters:
//   - ctx: Context
//   - record: Record cần kiểm tra (phải có field ID)
//   - structType: Type của struct (có thể là nil, sẽ tự động lấy từ record)
//
// Returns:
//   - error: Lỗi nếu tìm thấy quan hệ hoặc lỗi khác
func ValidateRelationshipsFromValue(ctx context.Context, record interface{}, structType reflect.Type) error {
	// Lấy ID từ record
	recordID, ok := getIDFromModel(record)
	if !ok {
		return common.NewError(
			common.ErrCodeValidation,
			"Record không có field ID",
			common.StatusBadRequest,
			nil,
		)
	}

	// Lấy type nếu chưa có
	if structType == nil {
		val := reflect.ValueOf(record)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		structType = val.Type()
	}

	return ValidateRelationships(ctx, recordID, structType)
}

// GetRelationshipDefinitions lấy danh sách các quan hệ được định nghĩa trong struct
//
// Parameters:
//   - structType: Type của struct
//
// Returns:
//   - []RelationshipDefinition: Danh sách các quan hệ
func GetRelationshipDefinitions(structType reflect.Type) []RelationshipDefinition {
	return ParseRelationshipTag(structType)
}
