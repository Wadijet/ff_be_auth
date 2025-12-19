# Đề Xuất Sync Customer từ POS

## 📋 Phân Tích Data POS

Từ API response, data POS có cấu trúc:

```json
{
  "id": "b0110315-b102-436b-8b3b-ed8d16740327",           // UUID string
  "name": "Trần Văn Hoàng",
  "gender": "male",
  "emails": ["thudo@gmail.com"],                          // Array
  "phone_numbers": ["0999999999"],                       // Array
  "date_of_birth": "1999-09-01",                         // Format: YYYY-MM-DD
  "reward_point": 10,
  "is_discount_by_level": true,
  "tags": [],
  "is_block": false,
  "assigned_user_id": "cee3c05e-5f85-43c4-b27e-889b99c50097",
  "level_id": null,                                       // String hoặc null
  "notes": [],
  "shop_customer_address": [...],
  "order_count": 108,
  "purchased_amount": 0,
  "succeed_order_count": 8,
  "last_order_at": "2020-04-01T10:18:41Z",
  "conversation_link": "https://pancake.vn/...",
  "referral_code": "1nw4geGA",
  "fb_id": null,                                         // Có thể link với Pancake
  "customer_id": "96a8e283-3fba-492e-a35a-970f72a30a02"
}
```

---

## 🔄 Mapping Data POS → Customer Model

### 1. Identifiers

| POS Field | Customer Field | Extract Tag | Notes |
|-----------|----------------|-------------|-------|
| `id` | `PosCustomerId` | `PosData\\.id,converter=string,optional` | UUID string - ID của hệ thống POS |
| `customer_id` | - | - | Không lưu (đã có `id` mặc định của model) |
| `fb_id` | - | - | Có thể dùng để link với Pancake (nếu có) |

**Lưu ý:** 
- POS `id` là UUID string, đây là ID của hệ thống POS → lưu vào `PosCustomerId`
- Pancake `id` → lưu vào `PanCakeCustomerId`
- **Phương án: Lưu riêng** - Mỗi nguồn có ID riêng để identify
- Customer không cần lưu `shopId` (không có field `PosShopId`)
- Không cần `ExternalCustomerId` - đã có `id` mặc định của model để identify customer chung

### 2. Common Fields (Multi-Source)

| POS Field | Customer Field | Extract Tag | Merge Strategy |
|-----------|----------------|-------------|----------------|
| `name` | `Name` | `PosData\\.name,converter=string,optional,priority=1,merge=priority` | `priority` (POS priority=1, Pancake priority=2) |
| `phone_numbers` | `PhoneNumbers` | `PosData\\.phone_numbers,optional,priority=1,merge=merge_array` | `merge_array` (POS là array, không cần converter) |
| `emails` | `Email` | `PosData\\.emails,converter=array_first,optional,priority=1,merge=priority` | `priority` (lấy email đầu tiên từ array, POS priority=1) |

**Lưu ý:**
- **Ưu tiên POS hơn Pancake** vì sale thao tác và cập nhật trên POS
- POS `phone_numbers` là array → không cần converter
- POS `emails` là array → cần converter `array_first` để lấy email đầu tiên
- Hoặc có thể merge tất cả emails vào array riêng

### 3. POS-Specific Fields

| POS Field | Customer Field | Extract Tag | Merge Strategy |
|-----------|----------------|-------------|----------------|
| `date_of_birth` | `Birthday` | `PosData\\.date_of_birth,converter=string,optional,merge=keep_existing` | `keep_existing` (nếu Pancake đã có) |
| `gender` | `Gender` | `PosData\\.gender,converter=string,optional,merge=keep_existing` | `keep_existing` (nếu Pancake đã có) |
| `reward_point` | `Point` | `PosData\\.reward_point,converter=int64,optional,merge=overwrite` | `overwrite` (luôn cập nhật) |
| `level_id` | `CustomerLevelId` | `PosData\\.level_id,converter=string,optional,merge=overwrite` | `overwrite` (UUID string, không phải int64) |
| `tags` | `TagIds` | `PosData\\.tags,optional,merge=overwrite` | `overwrite` (array) |
| `order_count` | `TotalOrder` | `PosData\\.order_count,converter=int64,optional,merge=overwrite` | `overwrite` |
| `purchased_amount` | `TotalSpent` | `PosData\\.purchased_amount,converter=number,optional,merge=overwrite` | `overwrite` |
| `succeed_order_count` | - | - | Có thể thêm field mới nếu cần |
| `last_order_at` | `PosLastOrderAt` | `PosData\\.last_order_at,converter=time,format=2006-01-02T15:04:05Z,optional` | `overwrite` |
| `shop_customer_address` | `PosAddresses` | `PosData\\.shop_customer_address,optional,merge=overwrite` | `overwrite` (array) |
| `referral_code` | `PosReferralCode` | `PosData\\.referral_code,converter=string,optional,merge=overwrite` | `overwrite` |
| `is_block` | `PosIsBlock` | `PosData\\.is_block,converter=bool,optional,merge=overwrite` | `overwrite` |

---

## 🔍 Logic Identify Customer

### Khi Upsert từ POS

**Thứ tự ưu tiên tìm customer:**

1. **Theo `posCustomerId`** (ưu tiên nhất)
   ```go
   filter := bson.M{
       "posCustomerId": posData["id"],
   }
   ```

2. **Theo `fb_id`** (nếu POS có fb_id, link với Pancake)
   ```go
   if fbId, ok := posData["fb_id"].(string); ok && fbId != "" {
       filter := bson.M{
           "psid": fbId, // Link với Pancake PSID
       }
   }
   ```

3. **Theo `phone_numbers`** (tìm trong array)
   ```go
   if phoneNumbers, ok := posData["phone_numbers"].([]interface{}); ok && len(phoneNumbers) > 0 {
       filter := bson.M{
           "phoneNumbers": bson.M{
               "$in": phoneNumbers, // Tìm trong array
           },
       }
   }
   ```

4. **Theo `emails`** (lấy email đầu tiên)
   ```go
   if emails, ok := posData["emails"].([]interface{}); ok && len(emails) > 0 {
       if email, ok := emails[0].(string); ok && email != "" {
           filter := bson.M{
               "email": email,
           }
       }
   }
   ```

5. **Tạo mới** (nếu không tìm thấy)

---

## 📝 Cập Nhật Model Customer

```go
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer lưu thông tin khách hàng từ các nguồn (Pancake, POS, ...)
type Customer struct {
	ID                primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	
	// ===== COMMON FIELDS (Extract từ nhiều nguồn với conflict resolution) =====
	// Name: Ưu tiên POS (priority=1) hơn Pancake (priority=2) - vì sale thao tác trên POS
	Name              string                 `json:"name" bson:"name" index:"text" extract:"PosData\\.name,converter=string,optional,priority=1,merge=priority|PanCakeData\\.name,converter=string,optional,priority=2,merge=priority"`
	
	// PhoneNumbers: Merge từ tất cả nguồn vào array
	// - POS: phone_numbers (array) - ưu tiên
	// - Pancake: phone_numbers (array)
	PhoneNumbers      []string               `json:"phoneNumbers" bson:"phoneNumbers" index:"text" extract:"PosData\\.phone_numbers,optional,priority=1,merge=merge_array|PanCakeData\\.phone_numbers,optional,priority=2,merge=merge_array"`
	
	// Email: Ưu tiên POS (priority=1) hơn Pancake (priority=2)
	// - POS: emails (array) → lấy email đầu tiên
	// - Pancake: email (string)
	Email             string                 `json:"email" bson:"email" index:"text" extract:"PosData\\.emails,converter=array_first,optional,priority=1,merge=priority|PanCakeData\\.email,converter=string,optional,priority=2,merge=priority"`
	
	// ===== SOURCE-SPECIFIC IDENTIFIERS =====
	PanCakeCustomerId string                 `json:"panCakeCustomerId" bson:"panCakeCustomerId" index:"text" extract:"PanCakeData\\.id,converter=string,optional"` // Pancake ID (từ id)
	Psid              string                 `json:"psid" bson:"psid" index:"text" extract:"PanCakeData\\.psid,converter=string,optional"`
	PageId            string                 `json:"pageId" bson:"pageId" index:"text" extract:"PanCakeData\\.page_id,converter=string,optional"`
	
	PosCustomerId     string                 `json:"posCustomerId" bson:"posCustomerId" index:"text" extract:"PosData\\.id,converter=string,optional"` // UUID string - ID của hệ thống POS
	
	// ===== SOURCE-SPECIFIC DATA =====
	PanCakeData       map[string]interface{} `json:"panCakeData,omitempty" bson:"panCakeData,omitempty"`
	PosData           map[string]interface{} `json:"posData,omitempty" bson:"posData,omitempty"`
	
	// ===== EXTRACTED FIELDS (Từ các nguồn) =====
	// Common fields có thể có từ cả 2 nguồn - ưu tiên POS (priority=1)
	Birthday          string                 `json:"birthday,omitempty" bson:"birthday,omitempty" extract:"PosData\\.date_of_birth,converter=string,optional,priority=1,merge=priority|PanCakeData\\.birthday,converter=string,optional,priority=2,merge=priority"`
	Gender            string                 `json:"gender,omitempty" bson:"gender,omitempty" extract:"PosData\\.gender,converter=string,optional,priority=1,merge=priority|PanCakeData\\.gender,converter=string,optional,priority=2,merge=priority"`
	
	// Pancake-specific
	LivesIn           string                 `json:"livesIn,omitempty" bson:"livesIn,omitempty" extract:"PanCakeData\\.lives_in,converter=string,optional,merge=keep_existing"`
	PanCakeUpdatedAt  int64                  `json:"panCakeUpdatedAt" bson:"panCakeUpdatedAt" extract:"PanCakeData\\.updated_at,converter=time,format=2006-01-02T15:04:05.000000,optional"`
	
	// POS-specific
	CustomerLevelId   string                 `json:"customerLevelId,omitempty" bson:"customerLevelId,omitempty" extract:"PosData\\.level_id,converter=string,optional,merge=overwrite"` // UUID string
	Point             int64                  `json:"point,omitempty" bson:"point,omitempty" extract:"PosData\\.reward_point,converter=int64,optional,merge=overwrite"`
	TotalOrder        int64                  `json:"totalOrder,omitempty" bson:"totalOrder,omitempty" extract:"PosData\\.order_count,converter=int64,optional,merge=overwrite"`
	TotalSpent        float64                `json:"totalSpent,omitempty" bson:"totalSpent,omitempty" extract:"PosData\\.purchased_amount,converter=number,optional,merge=overwrite"`
	SucceedOrderCount int64                  `json:"succeedOrderCount,omitempty" bson:"succeedOrderCount,omitempty" extract:"PosData\\.succeed_order_count,converter=int64,optional,merge=overwrite"`
	TagIds            []interface{}          `json:"tagIds,omitempty" bson:"tagIds,omitempty" extract:"PosData\\.tags,optional,merge=overwrite"` // Array, có thể là string hoặc object
	PosLastOrderAt    int64                  `json:"posLastOrderAt,omitempty" bson:"posLastOrderAt,omitempty" extract:"PosData\\.last_order_at,converter=time,format=2006-01-02T15:04:05Z,optional"`
	PosAddresses      []interface{}          `json:"posAddresses,omitempty" bson:"posAddresses,omitempty" extract:"PosData\\.shop_customer_address,optional,merge=overwrite"`
	PosReferralCode   string                 `json:"posReferralCode,omitempty" bson:"posReferralCode,omitempty" extract:"PosData\\.referral_code,converter=string,optional,merge=overwrite"`
	PosIsBlock        bool                   `json:"posIsBlock,omitempty" bson:"posIsBlock,omitempty" extract:"PosData\\.is_block,converter=bool,optional,merge=overwrite"`
	
	// ===== METADATA =====
	Sources           []string               `json:"sources" bson:"sources"` // ["pancake", "pos"]
	CreatedAt         int64                  `json:"createdAt" bson:"createdAt"`
	UpdatedAt         int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

**Lưu ý quan trọng:**
- `PosCustomerId` là UUID string - ID của hệ thống POS (từ POS `id`)
- `PanCakeCustomerId` là ID của Pancake (từ Pancake `id`)
- **Phương án: Lưu riêng** - Mỗi nguồn có ID riêng, không cần `ExternalCustomerId` (đã có `id` mặc định của model)
- `CustomerLevelId` là UUID string (không phải int64)
- **Không có `PosShopId`** - Customer không cần lưu shopId
- **Ưu tiên POS hơn Pancake** cho thông tin cá nhân (POS priority=1, Pancake priority=2) vì sale thao tác trên POS
- Cần thêm converter `array_first` để lấy phần tử đầu tiên từ array

---

## 🔧 Implementation: Service Method

### UpsertFromPos

**File:** `api/core/api/services/service.customer.go`

```go
// UpsertFromPos upsert customer từ POS data
func (s *CustomerService) UpsertFromPos(ctx context.Context, posData map[string]interface{}) (models.Customer, error) {
	now := time.Now().UnixMilli()
	
	// 1. Identify customer (tìm customer hiện có)
	var existingCustomer models.Customer
	found := false
	
	// 1.1. Tìm theo posCustomerId (ưu tiên nhất)
	if posId, ok := posData["id"].(string); ok && posId != "" {
		filter := bson.M{
			"posCustomerId": posId,
		}
		err := s.collection.FindOne(ctx, filter).Decode(&existingCustomer)
		if err == nil {
			found = true
		}
	}
	
	// 1.2. Tìm theo fb_id (nếu có, link với Pancake)
	if !found {
		if fbId, ok := posData["fb_id"].(string); ok && fbId != "" {
			filter := bson.M{
				"psid": fbId, // Link với Pancake PSID
			}
			err := s.collection.FindOne(ctx, filter).Decode(&existingCustomer)
			if err == nil {
				found = true
			}
		}
	}
	
	// 1.3. Tìm theo phone_numbers
	if !found {
		if phoneNumbers, ok := posData["phone_numbers"].([]interface{}); ok && len(phoneNumbers) > 0 {
			// Convert sang []string
			phones := make([]string, 0, len(phoneNumbers))
			for _, p := range phoneNumbers {
				if phone, ok := p.(string); ok && phone != "" {
					phones = append(phones, phone)
				}
			}
			
			if len(phones) > 0 {
				filter := bson.M{
					"phoneNumbers": bson.M{
						"$in": phones,
					},
				}
				err := s.collection.FindOne(ctx, filter).Decode(&existingCustomer)
				if err == nil {
					found = true
				}
			}
		}
	}
	
	// 1.4. Tìm theo emails (lấy email đầu tiên)
	if !found {
		if emails, ok := posData["emails"].([]interface{}); ok && len(emails) > 0 {
			if email, ok := emails[0].(string); ok && email != "" {
				filter := bson.M{
					"email": email,
				}
				err := s.collection.FindOne(ctx, filter).Decode(&existingCustomer)
				if err == nil {
					found = true
				}
			}
		}
	}
	
	// 2. Prepare data
	if found {
		// Update existing customer
		// Merge posData
		if existingCustomer.PosData == nil {
			existingCustomer.PosData = make(map[string]interface{})
		}
		for k, v := range posData {
			existingCustomer.PosData[k] = v
		}
		
		// Update sources
		if !contains(existingCustomer.Sources, "pos") {
			existingCustomer.Sources = append(existingCustomer.Sources, "pos")
		}
		
		existingCustomer.UpdatedAt = now
		
		// Extract data tự động (qua struct tag)
		if err := utility.ExtractDataIfExists(&existingCustomer); err != nil {
			return models.Customer{}, fmt.Errorf("extract data failed: %w", err)
		}
		
		// Save
		filter := bson.M{"_id": existingCustomer.ID}
		update := bson.M{"$set": existingCustomer}
		_, err := s.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return models.Customer{}, err
		}
		
		return existingCustomer, nil
	} else {
		// Create new customer
		newCustomer := models.Customer{
			PosData:   posData,
			Sources:   []string{"pos"},
			CreatedAt: now,
			UpdatedAt: now,
		}
		
		// Extract data tự động (qua struct tag)
		if err := utility.ExtractDataIfExists(&newCustomer); err != nil {
			return models.Customer{}, fmt.Errorf("extract data failed: %w", err)
		}
		
		// Save
		result, err := s.collection.InsertOne(ctx, newCustomer)
		if err != nil {
			return models.Customer{}, err
		}
		
		newCustomer.ID = result.InsertedID.(primitive.ObjectID)
		return newCustomer, nil
	}
}
```

---

## 🔧 Implementation: Handler & Route

### Handler

**File:** `api/core/api/handler/handler.customer.go`

```go
// HandleUpsertFromPos xử lý upsert customer từ POS
func (h *CustomerHandler) HandleUpsertFromPos(c *fiber.Ctx) error {
	var input struct {
		PosData map[string]interface{} `json:"posData" validate:"required"`
	}
	
	if err := c.BodyParser(&input); err != nil {
		return h.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}
	
	// Validate
	if err := h.validator.Struct(input); err != nil {
		return h.SendError(c, fiber.StatusBadRequest, "Validation failed", err)
	}
	
	// Upsert
	customer, err := h.service.UpsertFromPos(c.Context(), input.PosData)
	if err != nil {
		return h.SendError(c, fiber.StatusInternalServerError, "Failed to upsert customer", err)
	}
	
	return h.SendSuccess(c, customer)
}
```

### Route

**File:** `api/core/api/router/routes.go`

```go
// Thêm route
customerGroup.Post("/upsert-from-pos", customerHandler.HandleUpsertFromPos)
```

---

## 📊 Indexes Cần Thêm

```go
// Trong init.go
// Index cho posCustomerId (sparse, unique)
indexes = append(indexes, mongo.IndexModel{
	Keys: bson.D{
		{Key: "posCustomerId", Value: 1},
	},
	Options: options.Index().SetUnique(true).SetSparse(true),
})
```

---

## ✅ Tóm Tắt

1. **Mapping:** POS data → Customer model với extract tags phù hợp
2. **Identify:** Tìm customer theo thứ tự ưu tiên (posCustomerId → fb_id → phone → email)
3. **Merge:** Tự động merge qua extract tags với conflict resolution
4. **Priority:** Ưu tiên POS (priority=1) hơn Pancake (priority=2) cho thông tin cá nhân vì sale thao tác trên POS
5. **Customer ID Strategy:**
   - **Lưu riêng:** `PanCakeCustomerId` (từ Pancake `id`) và `PosCustomerId` (từ POS `id`)
   - Không cần `ExternalCustomerId` - đã có `id` mặc định của model để identify customer chung
6. **Converter mới:** Cần thêm `array_first` converter để lấy phần tử đầu tiên từ array
7. **Indexes:** Thêm unique index cho `posCustomerId` (sparse)

