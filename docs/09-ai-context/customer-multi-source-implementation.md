# Customer Multi-Source Implementation Guide

## 📋 Tổng Quan

Khi customer có data từ nhiều nguồn (Pancake, POS), cần xử lý conflict và merge data. Tài liệu này mô tả cách implement multi-source extract với conflict resolution.

---

## 🎯 Phương Án: Extract Tag với Priority và Merge Strategy

**Nguyên tắc:**
- Một field có thể có nhiều extract tags từ nhiều nguồn (phân tách bằng `|`)
- Mỗi nguồn có thể có converter riêng (vì định dạng dữ liệu khác nhau)
- Mỗi nguồn có priority và merge strategy riêng
- Backend xử lý: có nguồn nào thì extract theo định nghĩa của nguồn đó

---

## 📝 Format Extract Tag

### Format Hiện Tại (Single Source)
```
extract:"PanCakeData\\.name,converter=string,optional"
```

### Format Mới (Multi-Source)
```
extract:"Source1\\.path,converter=type1,optional,priority=1,merge=strategy1|Source2\\.path,converter=type2,optional,priority=2,merge=strategy2"
```

**Ví dụ:**
```go
// Name: Ưu tiên theo priority
Name string `extract:"PanCakeData\\.name,converter=string,optional,priority=1,merge=priority|PosData\\.name,converter=string,optional,priority=2,merge=priority"`

// PhoneNumbers: Merge vào array (Pancake là array, POS là string)
PhoneNumbers []string `extract:"PanCakeData\\.phone_numbers,optional,priority=1,merge=merge_array|PosData\\.phone_number,converter=string,optional,priority=2,merge=merge_array"`
```

**Các tham số:**
- `converter`: Converter cho nguồn này (có thể khác nhau giữa các nguồn)
- `priority`: Độ ưu tiên (số càng nhỏ càng ưu tiên, dùng khi `merge=priority`)
- `merge`: Chiến lược merge (`merge_array`, `keep_existing`, `overwrite`, `priority`)

---

## 🔄 Merge Strategies

### 1. Strategy: `merge_array` (Merge vào array)

**Mô tả:** Merge tất cả giá trị từ các nguồn vào một array, loại bỏ duplicate.

**Logic:**
- Collect tất cả giá trị từ các nguồn
- Nếu giá trị là array → thêm từng phần tử
- Nếu giá trị là scalar → thêm trực tiếp
- Loại bỏ duplicate

**Ví dụ:**
```go
PhoneNumbers []string `extract:"PanCakeData\\.phone_numbers,optional,priority=1,merge=merge_array|PosData\\.phone_number,converter=string,optional,priority=2,merge=merge_array"`
```

**Kết quả:**
- Pancake: `["0912345678", "0987654321"]`
- POS: `"0911111111"`
- → `PhoneNumbers` = `["0912345678", "0987654321", "0911111111"]`

**Khi nào dùng:**
- PhoneNumbers, TagIds, Addresses (cần tổng hợp từ nhiều nguồn)

**Lưu ý:**
- Chỉ áp dụng cho slice/array fields
- Tự động loại bỏ duplicate

---

### 2. Strategy: `keep_existing` (Giữ giá trị hiện có)

**Mô tả:** Nếu field đã có giá trị (không rỗng), giữ nguyên. Nếu field rỗng, lấy từ nguồn có data.

**Logic:**
```go
if !targetField.IsZero() {
    return nil // Giữ nguyên giá trị hiện có
}
return setFieldValue(targetField, values[0].value) // Lấy từ nguồn đầu tiên
```

**Ví dụ:**
```go
Birthday string `extract:"PanCakeData\\.birthday,converter=string,optional,merge=keep_existing"`
```

**Kết quả:**
- Field đã có `"1990-01-01"` → Giữ nguyên
- Field rỗng → Lấy từ Pancake

**Khi nào dùng:**
- Birthday, Gender, LivesIn (dữ liệu static, ít thay đổi)

**Lưu ý:**
- Chỉ set giá trị khi field rỗng
- Không cập nhật nếu field đã có giá trị

---

### 3. Strategy: `overwrite` (Luôn ghi đè) - Mặc định

**Mô tả:** Luôn lấy giá trị mới nhất từ nguồn có data, ghi đè giá trị cũ.

**Logic:**
```go
return setFieldValue(targetField, values[0].value) // Lấy từ nguồn đầu tiên
```

**Ví dụ:**
```go
Point int64 `extract:"PosData\\.point,converter=int64,optional,merge=overwrite"`
```

**Kết quả:**
- `point: 100` (cũ) → sync POS `point: 500` → `point: 500`

**Khi nào dùng:**
- Point, TotalOrder, TotalSpent, TagIds (dữ liệu dynamic, luôn cập nhật)

**Lưu ý:**
- Đây là strategy mặc định nếu không chỉ định `merge`
- Luôn ghi đè, không giữ giá trị cũ

---

### 4. Strategy: `priority` (Ưu tiên theo priority)

**Mô tả:** Chọn giá trị từ nguồn có `priority` nhỏ nhất (ưu tiên cao nhất). Priority = 0 được coi là ưu tiên thấp nhất.

**Logic:**
```go
priorityValue := values[0]
for _, v := range values[1:] {
    priority1 := priorityValue.config.Priority
    priority2 := v.config.Priority
    
    // Priority = 0 → ưu tiên thấp nhất
    if priority1 == 0 { priority1 = 999999 }
    if priority2 == 0 { priority2 = 999999 }
    
    if priority2 < priority1 {
        priorityValue = v
    }
}
return setFieldValue(targetField, priorityValue.value)
```

**Ví dụ:**
```go
Name string `extract:"PanCakeData\\.name,converter=string,optional,priority=1,merge=priority|PosData\\.name,converter=string,optional,priority=2,merge=priority"`
```

**Kết quả:**
- Pancake: `priority=1` → ưu tiên cao
- POS: `priority=2` → ưu tiên thấp hơn
- → Chọn giá trị từ Pancake

**Khi nào dùng:**
- Name, Email (ưu tiên nguồn cụ thể, không phụ thuộc thời gian)

**Lưu ý:**
- Priority càng nhỏ = ưu tiên càng cao
- Priority = 0 → ưu tiên thấp nhất

---

## 📊 Bảng Chiến Lược Merge Đề Xuất

| Field | Strategy | Lý Do |
|-------|----------|-------|
| **Name** | `priority` | Ưu tiên Pancake (priority=1) hơn POS (priority=2) |
| **PhoneNumbers** | `merge_array` | Merge tất cả số điện thoại, không mất thông tin |
| **Email** | `priority` | Ưu tiên Pancake (priority=1) hơn POS (priority=2) |
| **Birthday** | `keep_existing` | Ngày sinh không thay đổi, giữ giá trị đầu tiên |
| **Gender** | `keep_existing` | Giới tính không thay đổi, giữ giá trị đầu tiên |
| **LivesIn** | `keep_existing` | Nơi ở ít thay đổi, giữ giá trị đầu tiên |
| **Point** | `overwrite` | Điểm tích lũy luôn cập nhật từ POS |
| **TotalOrder** | `overwrite` | Tổng đơn hàng luôn cập nhật từ POS |
| **TotalSpent** | `overwrite` | Tổng tiền luôn cập nhật từ POS |
| **TagIds** | `overwrite` | Tags luôn cập nhật từ POS |

---

## 📝 Cấu Trúc Model Customer

```go
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer lưu thông tin khách hàng từ các nguồn (Pancake, POS, ...)
type Customer struct {
	ID                primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	
	// ===== COMMON FIELDS (Extract từ nhiều nguồn với conflict resolution) =====
	// Name: Ưu tiên Pancake (priority=1) hơn POS (priority=2)
	Name              string                 `json:"name" bson:"name" index:"text" extract:"PanCakeData\\.name,converter=string,optional,priority=1,merge=priority|PosData\\.name,converter=string,optional,priority=2,merge=priority"`
	
	// PhoneNumbers: Merge từ tất cả nguồn vào array
	// - Pancake: phone_numbers (array) → không cần converter
	// - POS: phone_number (string) → cần converter=string
	PhoneNumbers      []string               `json:"phoneNumbers" bson:"phoneNumbers" index:"text" extract:"PanCakeData\\.phone_numbers,optional,priority=1,merge=merge_array|PosData\\.phone_number,converter=string,optional,priority=2,merge=merge_array"`
	
	// Email: Ưu tiên Pancake (priority=1) hơn POS (priority=2)
	Email             string                 `json:"email" bson:"email" index:"text" extract:"PanCakeData\\.email,converter=string,optional,priority=1,merge=priority|PosData\\.email,converter=string,optional,priority=2,merge=priority"`
	
	// ===== SOURCE-SPECIFIC IDENTIFIERS =====
	PanCakeCustomerId string                 `json:"panCakeCustomerId" bson:"panCakeCustomerId" index:"text" extract:"PanCakeData\\.id,converter=string,optional"`
	Psid              string                 `json:"psid" bson:"psid" index:"text" extract:"PanCakeData\\.psid,converter=string,optional"`
	PageId            string                 `json:"pageId" bson:"pageId" index:"text" extract:"PanCakeData\\.page_id,converter=string,optional"`
	
	PosCustomerId     int64                  `json:"posCustomerId" bson:"posCustomerId" index:"text" extract:"PosData\\.id,converter=int64,optional"`
	PosShopId         int64                  `json:"posShopId" bson:"posShopId" index:"text" extract:"PosData\\.shop_id,converter=int64,optional"`
	
	// ===== SOURCE-SPECIFIC DATA =====
	PanCakeData       map[string]interface{} `json:"panCakeData,omitempty" bson:"panCakeData,omitempty"`
	PosData           map[string]interface{} `json:"posData,omitempty" bson:"posData,omitempty"`
	
	// ===== EXTRACTED FIELDS (Từ các nguồn) =====
	// Pancake-specific (chỉ có từ Pancake, không conflict)
	Birthday          string                 `json:"birthday,omitempty" bson:"birthday,omitempty" extract:"PanCakeData\\.birthday,converter=string,optional,merge=keep_existing"`
	Gender            string                 `json:"gender,omitempty" bson:"gender,omitempty" extract:"PanCakeData\\.gender,converter=string,optional,merge=keep_existing"`
	LivesIn           string                 `json:"livesIn,omitempty" bson:"livesIn,omitempty" extract:"PanCakeData\\.lives_in,converter=string,optional,merge=keep_existing"`
	PanCakeUpdatedAt  int64                  `json:"panCakeUpdatedAt" bson:"panCakeUpdatedAt" extract:"PanCakeData\\.updated_at,converter=time,format=2006-01-02T15:04:05.000000,optional"`
	
	// POS-specific (chỉ có từ POS, không conflict)
	CustomerLevelId   int64                  `json:"customerLevelId,omitempty" bson:"customerLevelId,omitempty" extract:"PosData\\.customer_level_id,converter=int64,optional,merge=overwrite"`
	Point             int64                  `json:"point,omitempty" bson:"point,omitempty" extract:"PosData\\.point,converter=int64,optional,merge=overwrite"`
	TotalOrder        int64                  `json:"totalOrder,omitempty" bson:"totalOrder,omitempty" extract:"PosData\\.total_order,converter=int64,optional,merge=overwrite"`
	TotalSpent        float64                `json:"totalSpent,omitempty" bson:"totalSpent,omitempty" extract:"PosData\\.total_spent,converter=number,optional,merge=overwrite"`
	TagIds            []int64                `json:"tagIds,omitempty" bson:"tagIds,omitempty" extract:"PosData\\.tags,optional,merge=overwrite"`
	PosUpdatedAt      int64                  `json:"posUpdatedAt,omitempty" bson:"posUpdatedAt,omitempty" extract:"PosData\\.updated_at,converter=time,optional"`
	
	// ===== METADATA =====
	Sources           []string               `json:"sources" bson:"sources"` // ["pancake", "pos"]
	CreatedAt         int64                  `json:"createdAt" bson:"createdAt"`
	UpdatedAt         int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

**Lưu ý:**
- Bỏ field `LastSyncedAt` (không cần vì không dùng strategy `latest`)
- Mỗi nguồn có thể có converter khác nhau
- Backend tự động: có nguồn nào thì extract theo định nghĩa của nguồn đó

---

## 🔧 Implementation Plan

### Bước 1: Cập Nhật Parse Extract Tag

**File:** `api/core/utility/data.extract.go`

**Cập nhật struct `extractTagConfig`:**
```go
type extractTagConfig struct {
	SourcePath    []string // Path đến source field và nested path
	Converter     string   // Converter name (có thể khác nhau giữa các nguồn)
	Format        string   // Format cho time converter
	Default       string   // Giá trị mặc định
	Optional      bool     // Flag optional
	Required      bool     // Flag required
	Priority      int      // Độ ưu tiên (số càng nhỏ càng ưu tiên, 0 = mặc định = ưu tiên thấp nhất)
	MergeStrategy string   // Chiến lược merge: "merge_array", "keep_existing", "overwrite", "priority"
}
```

**Cập nhật `parseExtractTag` để parse nhiều extract tags:**
```go
// parseExtractTag parse tag extract thành config
// Format mới: "Source1\\.path,converter=type1,options|Source2\\.path,converter=type2,options"
func parseExtractTag(tag string) ([]*extractTagConfig, error) {
	// Kiểm tra xem có nhiều nguồn không (có dấu |)
	if !strings.Contains(tag, "|") {
		// Single source - backward compatible
		config, err := parseSingleSourceTag(tag)
		if err != nil {
			return nil, err
		}
		return []*extractTagConfig{config}, nil
	}
	
	// Multi-source: Split bằng | để tách các nguồn
	sources := strings.Split(tag, "|")
	configs := make([]*extractTagConfig, 0, len(sources))
	
	for _, sourceTag := range sources {
		sourceTag = strings.TrimSpace(sourceTag)
		if sourceTag == "" {
			continue
		}
		
		config, err := parseSingleSourceTag(sourceTag)
		if err != nil {
			return nil, fmt.Errorf("parse source tag '%s': %w", sourceTag, err)
		}
		configs = append(configs, config)
	}
	
	return configs, nil
}

// parseSingleSourceTag parse một extract tag từ một nguồn
func parseSingleSourceTag(tag string) (*extractTagConfig, error) {
	config := &extractTagConfig{
		Converter:     "string", // Default converter
		Format:        "2006-01-02T15:04:05",
		Priority:      0, // Mặc định = ưu tiên thấp nhất
		MergeStrategy: "overwrite", // Mặc định: ghi đè
	}
	
	// Parse logic tương tự như hiện tại
	// Thêm parse cho priority và merge
	// ...
	
	return config, nil
}
```

### Bước 2: Cập Nhật Extract Logic

**File:** `api/core/utility/data.extract.go`

**Cập nhật `extractDataIfExists` để xử lý nhiều configs:**
```go
// extractDataIfExists extract data từ source fields vào typed fields
func extractDataIfExists(s interface{}) error {
	// ... existing code ...
	
	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structType.Field(i)
		
		extractTag := fieldType.Tag.Get("extract")
		if extractTag == "" {
			continue
		}
		
		// Parse tag - có thể trả về nhiều configs (multi-source)
		configs, err := parseExtractTag(extractTag)
		if err != nil {
			return fmt.Errorf("parse extract tag cho field %s: %w", fieldType.Name, err)
		}
		
		// Nếu chỉ có 1 config, xử lý như cũ (backward compatible)
		if len(configs) == 1 {
			if err := extractFieldValue(structVal, field, configs[0]); err != nil {
				// ... error handling như hiện tại ...
			}
			continue
		}
		
		// Nếu có nhiều configs (multi-source), xử lý conflict
		if err := extractFieldValueMultiSource(structVal, field, configs); err != nil {
			// ... error handling ...
		}
	}
	
	return nil
}
```

### Bước 3: Thêm Function Xử Lý Multi-Source

**File:** `api/core/utility/data.extract.go`

```go
// extractFieldValueMultiSource extract giá trị từ nhiều nguồn với conflict resolution
func extractFieldValueMultiSource(structVal reflect.Value, targetField reflect.Value, configs []*extractTagConfig) error {
	if len(configs) == 0 {
		return fmt.Errorf("không có config nào")
	}
	
	// Extract giá trị từ tất cả các nguồn (mỗi nguồn có converter riêng)
	values := make([]extractedValue, 0, len(configs))
	
	for _, config := range configs {
		// Kiểm tra xem nguồn này có data không
		sourceFieldName := config.SourcePath[0]
		sourceField := structVal.FieldByName(sourceFieldName)
		if !sourceField.IsValid() {
			continue // Nguồn không tồn tại, bỏ qua
		}
		
		// Kiểm tra source field có data không
		if sourceField.Kind() != reflect.Map {
			continue // Không phải map, bỏ qua
		}
		
		sourceMap, ok := sourceField.Interface().(map[string]interface{})
		if !ok || sourceMap == nil || len(sourceMap) == 0 {
			// Nguồn không có data, bỏ qua (nếu optional)
			if config.Optional {
				continue
			}
			// Nếu required và không có data, return error
			if config.Required {
				return fmt.Errorf("source field %s là required nhưng không có data", sourceFieldName)
			}
			continue
		}
		
		// Extract giá trị từ nguồn này (với converter riêng của nguồn)
		value, err := extractValueFromSource(structVal, config)
		if err != nil {
			// Nếu optional và không tìm thấy, bỏ qua
			if config.Optional && strings.Contains(err.Error(), "không tìm thấy") {
				continue
			}
			// Nếu required và không tìm thấy, return error
			if config.Required && strings.Contains(err.Error(), "không tìm thấy") {
				return err
			}
			// Nếu optional và có lỗi convert, bỏ qua
			if config.Optional {
				continue
			}
			return err
		}
		
		values = append(values, extractedValue{
			value:  value,
			config: config,
		})
	}
	
	if len(values) == 0 {
		// Không có nguồn nào có data, kiểm tra default
		for _, config := range configs {
			if config.Default != "" {
				return setFieldValue(targetField, config.Default, config)
			}
		}
		// Nếu tất cả đều optional, bỏ qua
		allOptional := true
		for _, config := range configs {
			if !config.Optional {
				allOptional = false
				break
			}
		}
		if allOptional {
			return nil // Bỏ qua field này
		}
		return fmt.Errorf("không tìm thấy giá trị từ bất kỳ nguồn nào")
	}
	
	// Áp dụng merge strategy
	return applyMergeStrategy(targetField, values)
}

type extractedValue struct {
	value  interface{}
	config *extractTagConfig
}

// applyMergeStrategy áp dụng chiến lược merge
func applyMergeStrategy(targetField reflect.Value, values []extractedValue) error {
	if len(values) == 0 {
		return fmt.Errorf("không có giá trị nào")
	}
	
	// Lấy merge strategy từ config đầu tiên (tất cả configs nên có cùng strategy)
	strategy := values[0].config.MergeStrategy
	if strategy == "" {
		strategy = "overwrite" // Mặc định
	}
	
	switch strategy {
	case "merge_array":
		// Merge tất cả giá trị vào array (loại bỏ duplicate)
		// Chỉ áp dụng cho slice/array fields
		if targetField.Type().Kind() != reflect.Slice {
			return fmt.Errorf("merge_array chỉ áp dụng cho slice/array fields")
		}
		
		// Collect tất cả giá trị
		allValues := make([]interface{}, 0)
		for _, v := range values {
			val := reflect.ValueOf(v.value)
			if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
				for i := 0; i < val.Len(); i++ {
					allValues = append(allValues, val.Index(i).Interface())
				}
			} else {
				allValues = append(allValues, v.value)
			}
		}
		
		// Loại bỏ duplicate
		uniqueValues := removeDuplicates(allValues)
		
		// Tạo slice mới
		elemType := targetField.Type().Elem()
		newSlice := reflect.MakeSlice(targetField.Type(), len(uniqueValues), len(uniqueValues))
		for i, val := range uniqueValues {
			valVal := reflect.ValueOf(val)
			if valVal.Type().AssignableTo(elemType) {
				newSlice.Index(i).Set(valVal)
			} else if valVal.Type().ConvertibleTo(elemType) {
				newSlice.Index(i).Set(valVal.Convert(elemType))
			}
		}
		
		targetField.Set(newSlice)
		return nil
		
	case "keep_existing":
		// Giữ giá trị hiện có nếu đã có, nếu không lấy từ nguồn đầu tiên
		if !targetField.IsZero() {
			return nil // Giữ nguyên giá trị hiện có
		}
		return setFieldValue(targetField, values[0].value, values[0].config)
		
	case "priority":
		// Chọn giá trị từ nguồn có priority nhỏ nhất (ưu tiên cao nhất)
		priorityValue := values[0]
		for _, v := range values[1:] {
			priority1 := priorityValue.config.Priority
			priority2 := v.config.Priority
			
			// Priority = 0 → ưu tiên thấp nhất (số lớn)
			if priority1 == 0 {
				priority1 = 999999
			}
			if priority2 == 0 {
				priority2 = 999999
			}
			
			if priority2 < priority1 {
				priorityValue = v
			}
		}
		return setFieldValue(targetField, priorityValue.value, priorityValue.config)
		
	case "overwrite":
		fallthrough
	default:
		// Mặc định: ghi đè bằng giá trị từ nguồn đầu tiên
		return setFieldValue(targetField, values[0].value, values[0].config)
	}
}

// extractValueFromSource extract giá trị từ một nguồn (với converter riêng của nguồn)
func extractValueFromSource(structVal reflect.Value, config *extractTagConfig) (interface{}, error) {
	// Logic extract từ source field với nested path
	// Apply converter riêng cho nguồn này
	// ...
}

// removeDuplicates loại bỏ duplicate trong array
func removeDuplicates(values []interface{}) []interface{} {
	seen := make(map[string]bool)
	result := make([]interface{}, 0)
	
	for _, val := range values {
		key := fmt.Sprintf("%v", val) // Convert sang string để so sánh
		if !seen[key] {
			seen[key] = true
			result = append(result, val)
		}
	}
	
	return result
}
```

---

## 🔄 Logic Upsert với Multi-Source

### Khi Upsert từ Pancake

**Filter:** `{"panCakeCustomerId": "xxx"}` hoặc `{"psid": "xxx", "pageId": "yyy"}`

**Logic:**
1. Tìm customer theo filter
2. Update `panCakeData`
3. Update `sources[]` (thêm "pancake" nếu chưa có)
4. Extract data tự động:
   - Duyệt qua tất cả extract tags
   - Với mỗi tag có nhiều nguồn:
     - Kiểm tra nguồn nào có data
     - Extract từ nguồn có data (với converter riêng của nguồn)
     - Nếu có nhiều nguồn có data → Áp dụng merge strategy

### Khi Upsert từ POS

**Filter:** `{"posCustomerId": 123, "posShopId": 456}`

**Logic:**
1. Tìm customer theo filter (hoặc phone/email để link)
2. Update `posData`
3. Update `sources[]` (thêm "pos" nếu chưa có)
4. Extract data tự động:
   - Duyệt qua tất cả extract tags
   - Với mỗi tag có nhiều nguồn:
     - Kiểm tra nguồn nào có data
     - Extract từ nguồn có data (với converter riêng của nguồn)
     - Nếu có nhiều nguồn có data → Áp dụng merge strategy

**Ví dụ cụ thể:**

**Khi upsert từ POS:**
- `posData` có `phone_number: "0911111111"` (string)
- Extract tag: `PanCakeData\\.phone_numbers,optional,priority=1,merge=merge_array|PosData\\.phone_number,converter=string,optional,priority=2,merge=merge_array`
- Logic:
  1. Kiểm tra `PanCakeData` → có data không? (có thể có hoặc không)
  2. Kiểm tra `PosData` → có data (`phone_number: "0911111111"`)
  3. Extract từ `PosData` với `converter=string` → `"0911111111"`
  4. Nếu `PanCakeData` cũng có `phone_numbers: ["0912345678"]`:
     - Extract từ `PanCakeData` (không cần converter vì đã là array)
     - Áp dụng `merge=merge_array` → Merge 2 arrays: `["0912345678", "0911111111"]`
  5. Nếu chỉ có `PosData`:
     - Chỉ extract từ `PosData` → `["0911111111"]`

---

## ✅ Khuyến Nghị

**Phương án đơn giản:**

1. **Format:** Dùng `|` để phân tách nhiều nguồn trong cùng 1 extract tag
2. **Converter độc lập:** Mỗi nguồn có thể có converter riêng
3. **Default Strategy:** `overwrite` (ghi đè) nếu không chỉ định
4. **Merge Strategy:** Hỗ trợ `merge_array`, `keep_existing`, `overwrite`, `priority`
5. **Backward Compatible:** Extract tag single source vẫn hoạt động (không có `|`)
6. **Logic xử lý:** Có nguồn nào thì extract theo định nghĩa của nguồn đó

**Ưu điểm:**
- ✅ Đơn giản, dễ hiểu
- ✅ Backward compatible
- ✅ Linh hoạt: mỗi nguồn có converter riêng
- ✅ Tự động xử lý conflict
- ✅ Có nguồn nào extract nguồn đó (không cần tất cả nguồn đều có data)
