# Đề Xuất Tách Riêng Customer: Pancake vs POS

## 📋 Tổng Quan

Tài liệu này đề xuất phương án **tách riêng** customer từ **Pancake (Facebook)** và **POS** thành 2 collections riêng biệt thay vì merge vào một collection như hiện tại.

---

## 🔍 Phân Tích Vấn Đề Hiện Tại

### Cấu Trúc Hiện Tại

Hiện tại, hệ thống đang sử dụng một model `Customer` duy nhất để lưu trữ dữ liệu từ cả 2 nguồn:

```go
type Customer struct {
    // Common fields với merge strategy phức tạp
    Name         string   // Extract từ cả PosData và PanCakeData với priority
    PhoneNumbers []string // Merge array từ cả 2 nguồn
    Email        string   // Priority resolution
    
    // Source-specific identifiers
    PanCakeCustomerId string
    Psid              string
    PageId            string
    PosCustomerId     string
    
    // Source-specific data
    PanCakeData map[string]interface{}
    PosData     map[string]interface{}
    
    // Extracted fields với merge strategies
    // ...
}
```

### Vấn Đề Của Cách Tiếp Cận Hiện Tại

#### 1. **Phức Tạp Về Logic Merge** ❌
- Cần xử lý nhiều merge strategies: `priority`, `merge_array`, `keep_existing`, `overwrite`
- Logic extract phức tạp với nhiều nguồn trong cùng một field
- Khó debug khi có conflict giữa các nguồn

#### 2. **Không Rõ Ràng Về Nguồn Dữ Liệu** ❌
- Không biết field nào đến từ nguồn nào
- Khó trace back khi có vấn đề về dữ liệu
- Khó maintain và update

#### 3. **Conflict Resolution Phức Tạp** ❌
- Cần quyết định priority cho từng field
- Logic merge có thể gây mất dữ liệu nếu không cẩn thận
- Khó test các edge cases

#### 4. **Performance Issues** ❌
- Document lớn hơn (có cả PanCakeData và PosData)
- Index phức tạp hơn (cần index cho cả 2 nguồn)
- Query có thể chậm hơn khi cần filter theo nguồn

#### 5. **Khó Mở Rộng** ❌
- Nếu thêm nguồn mới (ví dụ: Shopee, Lazada) sẽ phức tạp hơn nhiều
- Cần update logic merge cho tất cả fields
- Risk cao khi thay đổi

#### 6. **Không Phù Hợp Với Use Cases** ❌
- **Pancake Customer**: Chủ yếu dùng cho Facebook conversations, messages
- **POS Customer**: Chủ yếu dùng cho orders, points, loyalty programs
- Hai use cases này khác nhau, không cần merge

#### 7. **Data Integrity Issues** ❌
- Một customer có thể có data từ Pancake nhưng chưa có từ POS (hoặc ngược lại)
- Khó validate dữ liệu khi merge
- Có thể gây confusion khi một số fields có data, một số không

---

## ✅ Phương Án Đề Xuất: Tách Riêng

### Kiến Trúc Mới

Tách thành **2 collections riêng biệt**:

1. **`fb_customers`** - Customer từ Pancake (Facebook)
2. **`pc_pos_customers`** - Customer từ POS

### 1. FB Customer (Pancake/Facebook)

```go
// api/core/api/models/mongodb/model.fb.customer.go
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FbCustomer lưu thông tin khách hàng từ Pancake API (Facebook)
type FbCustomer struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	// ===== IDENTIFIERS =====
	CustomerId string `json:"customerId" bson:"customerId" index:"text,unique" extract:"PanCakeData\\.id,converter=string"` // Pancake Customer ID
	Psid       string `json:"psid" bson:"psid" index:"text,unique,sparse" extract:"PanCakeData\\.psid,converter=string,optional"` // Page Scoped ID (Facebook)
	PageId     string `json:"pageId" bson:"pageId" index:"text" extract:"PanCakeData\\.page_id,converter=string,optional"` // Facebook Page ID

	// ===== BASIC INFO =====
	Name         string   `json:"name" bson:"name" index:"text" extract:"PanCakeData\\.name,converter=string,optional"`
	PhoneNumbers []string `json:"phoneNumbers" bson:"phoneNumbers" index:"text" extract:"PanCakeData\\.phone_numbers,optional"`
	Email        string   `json:"email" bson:"email" index:"text" extract:"PanCakeData\\.email,converter=string,optional"`

	// ===== ADDITIONAL INFO =====
	Birthday string `json:"birthday,omitempty" bson:"birthday,omitempty" extract:"PanCakeData\\.birthday,converter=string,optional"`
	Gender   string `json:"gender,omitempty" bson:"gender,omitempty" extract:"PanCakeData\\.gender,converter=string,optional"`
	LivesIn  string `json:"livesIn,omitempty" bson:"livesIn,omitempty" extract:"PanCakeData\\.lives_in,converter=string,optional"`

	// ===== SOURCE DATA =====
	PanCakeData map[string]interface{} `json:"panCakeData,omitempty" bson:"panCakeData,omitempty"` // Dữ liệu gốc từ Pancake API

	// ===== METADATA =====
	PanCakeUpdatedAt int64 `json:"panCakeUpdatedAt" bson:"panCakeUpdatedAt" extract:"PanCakeData\\.updated_at,converter=time,format=2006-01-02T15:04:05.000000,optional"`
	CreatedAt        int64 `json:"createdAt" bson:"createdAt"`
	UpdatedAt        int64 `json:"updatedAt" bson:"updatedAt"`
}
```

**Collection Name:** `fb_customers`

**Unique Indexes:**
- `{customerId: 1}` - Unique
- `{psid: 1}` - Unique, sparse (vì không phải customer nào cũng có PSID)

**Use Cases:**
- Link với `fb_conversations` qua `psid` hoặc `customerId`
- Link với `fb_messages` qua `customerId`
- Hiển thị thông tin khách hàng trong Facebook conversations

---

### 2. POS Customer

```go
// api/core/api/models/mongodb/model.pc.pos.customer.go
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PcPosCustomer lưu thông tin khách hàng từ Pancake POS API
type PcPosCustomer struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	// ===== IDENTIFIERS =====
	CustomerId string `json:"customerId" bson:"customerId" index:"text,unique" extract:"PosData\\.id,converter=string"` // UUID string - POS Customer ID
	ShopId     int64  `json:"shopId" bson:"shopId" index:"text" extract:"PosData\\.shop_id,converter=int64,optional"` // Shop ID

	// ===== BASIC INFO =====
	Name         string   `json:"name" bson:"name" index:"text" extract:"PosData\\.name,converter=string,optional"`
	PhoneNumbers []string `json:"phoneNumbers" bson:"phoneNumbers" index:"text" extract:"PosData\\.phone_numbers,optional"`
	Emails       []string `json:"emails" bson:"emails" index:"text" extract:"PosData\\.emails,optional"` // POS có thể có nhiều emails

	// ===== ADDITIONAL INFO =====
	DateOfBirth string `json:"dateOfBirth,omitempty" bson:"dateOfBirth,omitempty" extract:"PosData\\.date_of_birth,converter=string,optional"`
	Gender      string `json:"gender,omitempty" bson:"gender,omitempty" extract:"PosData\\.gender,converter=string,optional"`

	// ===== POS-SPECIFIC FIELDS =====
	CustomerLevelId   string        `json:"customerLevelId,omitempty" bson:"customerLevelId,omitempty" extract:"PosData\\.level_id,converter=string,optional"` // UUID string
	Point             int64         `json:"point,omitempty" bson:"point,omitempty" extract:"PosData\\.reward_point,converter=int64,optional"`                // Điểm tích lũy
	TotalOrder        int64         `json:"totalOrder,omitempty" bson:"totalOrder,omitempty" extract:"PosData\\.order_count,converter=int64,optional"`       // Tổng đơn hàng
	TotalSpent        float64       `json:"totalSpent,omitempty" bson:"totalSpent,omitempty" extract:"PosData\\.purchased_amount,converter=number,optional"` // Tổng tiền đã mua
	SucceedOrderCount int64         `json:"succeedOrderCount,omitempty" bson:"succeedOrderCount,omitempty" extract:"PosData\\.succeed_order_count,converter=int64,optional"` // Số đơn hàng thành công
	TagIds            []interface{} `json:"tagIds,omitempty" bson:"tagIds,omitempty" extract:"PosData\\.tags,optional"`                                        // Tags (array)
	LastOrderAt       int64         `json:"lastOrderAt,omitempty" bson:"lastOrderAt,omitempty" extract:"PosData\\.last_order_at,converter=time,format=2006-01-02T15:04:05Z,optional"` // Thời gian đơn hàng cuối
	Addresses         []interface{} `json:"addresses,omitempty" bson:"addresses,omitempty" extract:"PosData\\.shop_customer_address,optional"`              // Địa chỉ (array)
	ReferralCode      string        `json:"referralCode,omitempty" bson:"referralCode,omitempty" extract:"PosData\\.referral_code,converter=string,optional"` // Mã giới thiệu
	IsBlock           bool          `json:"isBlock,omitempty" bson:"isBlock,omitempty" extract:"PosData\\.is_block,converter=bool,optional"`               // Trạng thái block

	// ===== SOURCE DATA =====
	PosData map[string]interface{} `json:"posData,omitempty" bson:"posData,omitempty"` // Dữ liệu gốc từ POS API

	// ===== METADATA =====
	PosUpdatedAt int64 `json:"posUpdatedAt" bson:"posUpdatedAt" extract:"PosData\\.updated_at,converter=time,format=2006-01-02T15:04:05Z,optional"`
	CreatedAt    int64 `json:"createdAt" bson:"createdAt"`
	UpdatedAt    int64 `json:"updatedAt" bson:"updatedAt"`
}
```

**Collection Name:** `pc_pos_customers`

**Unique Indexes:**
- `{customerId: 1}` - Unique (UUID string)

**Use Cases:**
- Link với `pc_pos_orders` qua `customerId`
- Hiển thị thông tin khách hàng trong orders
- Phân tích customer lifetime value, segmentation
- Quản lý điểm tích lũy, loyalty programs

---

## 🔗 Linking Giữa 2 Collections (Nếu Cần)

Nếu cần link giữa FB Customer và POS Customer, có thể:

### Phương Án 1: Reference Field (Đơn Giản)

Thêm field reference vào mỗi collection:

```go
// Trong FbCustomer
LinkedPosCustomerId string `json:"linkedPosCustomerId,omitempty" bson:"linkedPosCustomerId,omitempty"` // Reference to PcPosCustomer.customerId

// Trong PcPosCustomer
LinkedFbCustomerId string `json:"linkedFbCustomerId,omitempty" bson:"linkedFbCustomerId,omitempty"` // Reference to FbCustomer.customerId
```

**Cách link:** Dựa trên `phoneNumbers` hoặc `email` matching (có thể tự động hoặc manual)

### Phương Án 2: Separate Linking Collection (Linh Hoạt Hơn)

Tạo collection riêng để link:

```go
// api/core/api/models/mongodb/model.customer.link.go
type CustomerLink struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	
	FbCustomerId  string `json:"fbCustomerId" bson:"fbCustomerId" index:"text"` // Reference to FbCustomer.customerId
	PosCustomerId string `json:"posCustomerId" bson:"posCustomerId" index:"text"` // Reference to PcPosCustomer.customerId
	
	// Matching criteria
	MatchedBy     string   `json:"matchedBy" bson:"matchedBy"` // "phone", "email", "manual"
	MatchedFields []string `json:"matchedFields" bson:"matchedFields"` // ["phone", "email"]
	Confidence    float64  `json:"confidence" bson:"confidence"` // 0.0 - 1.0
	
	CreatedAt int64 `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt"`
}
```

**Collection Name:** `customer_links`

**Indexes:**
- `{fbCustomerId: 1}`
- `{posCustomerId: 1}`
- `{fbCustomerId: 1, posCustomerId: 1}` - Compound unique

**Ưu điểm:**
- Linh hoạt hơn, có thể có nhiều links
- Có thể track matching confidence
- Dễ query và maintain

---

## 📊 So Sánh: Merge vs Tách Riêng

| Tiêu Chí | Merge (Hiện Tại) | Tách Riêng (Đề Xuất) |
|----------|------------------|----------------------|
| **Độ Phức Tạp** | ⚠️ Cao (merge logic phức tạp) | ✅ Thấp (đơn giản, rõ ràng) |
| **Maintainability** | ⚠️ Khó maintain | ✅ Dễ maintain |
| **Clarity** | ⚠️ Không rõ nguồn dữ liệu | ✅ Rõ ràng từng nguồn |
| **Performance** | ⚠️ Document lớn, index phức tạp | ✅ Document nhỏ, index đơn giản |
| **Scalability** | ⚠️ Khó mở rộng thêm nguồn | ✅ Dễ thêm nguồn mới |
| **Use Case Fit** | ⚠️ Không phù hợp (2 use cases khác nhau) | ✅ Phù hợp (mỗi collection cho 1 use case) |
| **Data Integrity** | ⚠️ Khó validate | ✅ Dễ validate |
| **Testing** | ⚠️ Khó test edge cases | ✅ Dễ test |
| **Query Performance** | ⚠️ Có thể chậm | ✅ Nhanh hơn (document nhỏ) |
| **Linking** | ✅ Tự động (cùng document) | ⚠️ Cần logic link riêng (nhưng linh hoạt hơn) |

---

## 🚀 Kế Hoạch Migration

### Phase 1: Tạo Models Mới

1. Tạo `model.fb.customer.go`
2. Tạo `model.pc.pos.customer.go`
3. Tạo `model.customer.link.go` (nếu dùng linking collection)

### Phase 2: Tạo Services & Handlers

1. Tạo `service.fb.customer.go` và `handler.fb.customer.go`
2. Tạo `service.pc.pos.customer.go` và `handler.pc.pos.customer.go`
3. Tạo `service.customer.link.go` và `handler.customer.link.go` (nếu cần)

### Phase 3: Migration Data

1. **Script Migration:**
   ```go
   // scripts/migrate_customers.go
   // 1. Đọc từ collection `customers` cũ
   // 2. Tách thành FbCustomer và PcPosCustomer
   // 3. Insert vào collections mới
   // 4. Tạo links nếu có matching
   ```

2. **Migration Logic:**
   - Nếu có `PanCakeData` → Tạo `FbCustomer`
   - Nếu có `PosData` → Tạo `PcPosCustomer`
   - Nếu có cả 2 → Tạo cả 2 và link
   - Link dựa trên `phoneNumbers` hoặc `email`

### Phase 4: Update References

1. Update `fb_conversations` để reference `fb_customers` thay vì `customers`
2. Update `pc_pos_orders` để reference `pc_pos_customers` thay vì `customers`
3. Update các handlers/services khác

### Phase 5: Deprecate Old Model

1. Mark `Customer` model as deprecated
2. Keep collection `customers` để backup (không xóa ngay)
3. Sau 1-2 tháng, có thể archive hoặc xóa

---

## 📝 Implementation Details

### 1. Service Structure

```go
// api/core/api/services/service.fb.customer.go
type FbCustomerService struct {
	*BaseServiceMongoImpl[models.FbCustomer]
}

// api/core/api/services/service.pc.pos.customer.go
type PcPosCustomerService struct {
	*BaseServiceMongoImpl[models.PcPosCustomer]
}
```

### 2. Handler Structure

```go
// api/core/api/handler/handler.fb.customer.go
type FbCustomerHandler struct {
	BaseHandler[models.FbCustomer, dto.FbCustomerCreateInput, dto.FbCustomerUpdateInput]
	FbCustomerService *services.FbCustomerService
}

// api/core/api/handler/handler.pc.pos.customer.go
type PcPosCustomerHandler struct {
	BaseHandler[models.PcPosCustomer, dto.PcPosCustomerCreateInput, dto.PcPosCustomerUpdateInput]
	PcPosCustomerService *services.PcPosCustomerService
}
```

### 3. Routes

```go
// api/core/api/router/routes.go

// FB Customer routes
fbCustomerHandler := handlers.NewFbCustomerHandler()
apiV1.Post("/fb-customer/upsert-one", fbCustomerHandler.Upsert)
apiV1.Get("/fb-customer/find", fbCustomerHandler.Find)
// ... other CRUD operations

// POS Customer routes
pcPosCustomerHandler := handlers.NewPcPosCustomerHandler()
apiV1.Post("/pc-pos-customer/upsert-one", pcPosCustomerHandler.Upsert)
apiV1.Get("/pc-pos-customer/find", pcPosCustomerHandler.Find)
// ... other CRUD operations

// Customer Link routes (nếu cần)
customerLinkHandler := handlers.NewCustomerLinkHandler()
apiV1.Post("/customer-link/create", customerLinkHandler.InsertOne)
apiV1.Get("/customer-link/find-by-fb", customerLinkHandler.FindByFbCustomer)
apiV1.Get("/customer-link/find-by-pos", customerLinkHandler.FindByPosCustomer)
```

### 4. Collection Registration

```go
// api/core/global/global.vars.go
type MongoDB_ColNames struct {
	// ... existing collections
	FbCustomers    string // "fb_customers"
	PcPosCustomers string // "pc_pos_customers"
	CustomerLinks  string // "customer_links" (nếu dùng)
}
```

---

## 🎯 Use Cases Sau Khi Tách

### Use Case 1: Hiển Thị Customer Trong Facebook Conversation

```go
// Lấy customer từ conversation
conversation := getFbConversation(conversationId)
fbCustomer := fbCustomerService.FindOneByCustomerId(conversation.CustomerId)

// Nếu cần thông tin POS (nếu có link)
if link := customerLinkService.FindByFbCustomer(fbCustomer.CustomerId); link != nil {
	posCustomer := pcPosCustomerService.FindOneByCustomerId(link.PosCustomerId)
	// Merge data để hiển thị
}
```

### Use Case 2: Hiển Thị Customer Trong POS Order

```go
// Lấy customer từ order
order := getPcPosOrder(orderId)
posCustomer := pcPosCustomerService.FindOneByCustomerId(order.CustomerId)

// Nếu cần thông tin Facebook (nếu có link)
if link := customerLinkService.FindByPosCustomer(posCustomer.CustomerId); link != nil {
	fbCustomer := fbCustomerService.FindOneByCustomerId(link.FbCustomerId)
	// Merge data để hiển thị
}
```

### Use Case 3: Customer Matching (Tự Động hoặc Manual)

```go
// Tự động match dựa trên phone/email
func AutoMatchCustomers() {
	fbCustomers := fbCustomerService.FindAll()
	posCustomers := pcPosCustomerService.FindAll()
	
	for _, fb := range fbCustomers {
		for _, pos := range posCustomers {
			if matchPhoneOrEmail(fb, pos) {
				// Tạo link
				customerLinkService.CreateLink(fb.CustomerId, pos.CustomerId, "auto", 0.9)
			}
		}
	}
}
```

---

## ✅ Kết Luận

### Ưu Điểm Của Phương Án Tách Riêng

1. ✅ **Đơn giản, rõ ràng**: Mỗi collection có mục đích riêng
2. ✅ **Dễ maintain**: Không cần logic merge phức tạp
3. ✅ **Performance tốt hơn**: Document nhỏ hơn, index đơn giản hơn
4. ✅ **Dễ mở rộng**: Thêm nguồn mới chỉ cần tạo collection mới
5. ✅ **Phù hợp use cases**: Mỗi collection phục vụ use case riêng
6. ✅ **Data integrity tốt hơn**: Dễ validate và kiểm tra

### Nhược Điểm

1. ⚠️ **Cần logic link riêng**: Nhưng linh hoạt và có thể control được
2. ⚠️ **Cần migration**: Nhưng chỉ làm 1 lần

### Khuyến Nghị

**Nên tách riêng** vì:
- Phù hợp với kiến trúc hiện tại (đã tách riêng các collections khác)
- Đơn giản hơn, dễ maintain hơn
- Performance tốt hơn
- Dễ mở rộng trong tương lai

---

## 📚 Tài Liệu Tham Khảo

- [Customer Multi-Source Implementation](./customer-multi-source-implementation.md)
- [Data Architecture Overview](./data-architecture-overview.md)
- [Pancake POS Folkform Sync Analysis](./pancake-pos-folkform-sync-analysis.md)

---

**Ngày tạo**: 2025-01-XX  
**Phiên bản**: 1.0  
**Tác giả**: AI Assistant  
**Trạng thái**: Đề xuất
