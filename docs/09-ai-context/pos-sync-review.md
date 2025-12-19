# Rà Soát Đồng Bộ Thông Tin Từ POS

## 📋 Tổng Quan

Tài liệu này rà soát các thông tin đã và chưa được đồng bộ từ **Pancake POS API** (`pos.pages.fm/api/v1`) về **Folkform Backend**.

**Ngày rà soát:** 2025-01-XX  
**Trạng thái:** Đang triển khai

---

## ✅ ĐÃ ĐỒNG BỘ

### 1. Customer Model (Multi-Source)

**Model:** `api/core/api/models/mongodb/model.customer.go`

**Các field đã được extract từ POS:**

| POS Field | Customer Field | Extract Tag | Ghi Chú |
|-----------|----------------|-------------|---------|
| `id` | `PosCustomerId` | `PosData\\.id,converter=string,optional` | UUID string - ID của hệ thống POS |
| `id` | `CustomerId` | `PosData\\.id,converter=string,optional,priority=1` | ID chung để identify customer |
| `name` | `Name` | `PosData\\.name,converter=string,optional,priority=1,merge=priority` | Ưu tiên POS hơn Pancake |
| `phone_numbers` | `PhoneNumbers` | `PosData\\.phone_numbers,optional,priority=1,merge=merge_array` | Merge vào array |
| `emails` | `Email` | `PosData\\.emails,converter=array_first,optional,priority=1,merge=priority` | Lấy email đầu tiên |
| `date_of_birth` | `Birthday` | `PosData\\.date_of_birth,converter=string,optional,priority=1,merge=priority` | Ngày sinh |
| `gender` | `Gender` | `PosData\\.gender,converter=string,optional,priority=1,merge=priority` | Giới tính |
| `level_id` | `CustomerLevelId` | `PosData\\.level_id,converter=string,optional,merge=overwrite` | UUID string |
| `reward_point` | `Point` | `PosData\\.reward_point,converter=int64,optional,merge=overwrite` | Điểm tích lũy |
| `order_count` | `TotalOrder` | `PosData\\.order_count,converter=int64,optional,merge=overwrite` | Tổng đơn hàng |
| `purchased_amount` | `TotalSpent` | `PosData\\.purchased_amount,converter=number,optional,merge=overwrite` | Tổng tiền đã mua |
| `succeed_order_count` | `SucceedOrderCount` | `PosData\\.succeed_order_count,converter=int64,optional,merge=overwrite` | Số đơn hàng thành công |
| `tags` | `TagIds` | `PosData\\.tags,optional,merge=overwrite` | Tags (array) |
| `last_order_at` | `PosLastOrderAt` | `PosData\\.last_order_at,converter=time,format=2006-01-02T15:04:05Z,optional` | Thời gian đơn hàng cuối |
| `shop_customer_address` | `PosAddresses` | `PosData\\.shop_customer_address,optional,merge=overwrite` | Địa chỉ (array) |
| `referral_code` | `PosReferralCode` | `PosData\\.referral_code,converter=string,optional,merge=overwrite` | Mã giới thiệu |
| `is_block` | `PosIsBlock` | `PosData\\.is_block,converter=bool,optional,merge=overwrite` | Trạng thái block |

**Raw Data:**
- `PosData` - Lưu toàn bộ dữ liệu gốc từ POS API

---

### 2. PcPosShop Model

**Model:** `api/core/api/models/mongodb/model.pc.pos.shop.go`

**Các field đã được extract:**

| POS Field | Model Field | Extract Tag | Ghi Chú |
|-----------|-------------|-------------|---------|
| `id` | `ShopId` | `PanCakeData\\.id,converter=int64` | ID của shop trên Pancake POS |
| `name` | `Name` | `PanCakeData\\.name,converter=string,optional` | Tên cửa hàng |
| `avatar_url` | `AvatarUrl` | `PanCakeData\\.avatar_url,converter=string,optional` | Link hình đại diện |
| `pages` | `Pages` | `PanCakeData\\.pages,optional` | Thông tin các pages được gộp |

**Raw Data:**
- `PanCakeData` - Lưu toàn bộ dữ liệu gốc từ Pancake POS API

---

### 3. PcPosWarehouse Model

**Model:** `api/core/api/models/mongodb/model.pc.pos.warehouse.go`

**Các field đã được extract:**

| POS Field | Model Field | Extract Tag | Ghi Chú |
|-----------|-------------|-------------|---------|
| `id` | `WarehouseId` | `PanCakeData\\.id,converter=string` | UUID string - ID của warehouse |
| `shop_id` | `ShopId` | `PanCakeData\\.shop_id,converter=int64,optional` | ID của shop |
| `name` | `Name` | `PanCakeData\\.name,converter=string,optional` | Tên kho hàng |
| `phone_number` | `PhoneNumber` | `PanCakeData\\.phone_number,converter=string,optional` | Số điện thoại |
| `full_address` | `FullAddress` | `PanCakeData\\.full_address,converter=string,optional` | Địa chỉ đầy đủ |
| `province_id` | `ProvinceId` | `PanCakeData\\.province_id,converter=string,optional` | ID tỉnh/thành phố |
| `district_id` | `DistrictId` | `PanCakeData\\.district_id,converter=string,optional` | ID quận/huyện |
| `commune_id` | `CommuneId` | `PanCakeData\\.commune_id,converter=string,optional` | ID phường/xã |

**Raw Data:**
- `PanCakeData` - Lưu toàn bộ dữ liệu gốc từ Pancake POS API

---

## ❌ CHƯA ĐỒNG BỘ

### 1. Customer Model - Các Field Còn Thiếu

**Từ POS Customer API response, các field sau chưa được extract:**

| POS Field | Loại | Ghi Chú | Đề Xuất |
|-----------|------|---------|---------|
| `assigned_user_id` | UUID string | ID người dùng được gán cho customer | ⭐⭐⭐ Nên thêm `PosAssignedUserId` |
| `is_discount_by_level` | bool | Có được giảm giá theo cấp độ không | ⭐⭐ Có thể thêm `PosIsDiscountByLevel` |
| `notes` | Array | Ghi chú khách hàng (có thể có nhiều) | ⭐⭐⭐ Nên thêm `PosNotes []interface{}` |
| `conversation_link` | string | Link đến conversation trên Pancake | ⭐⭐ Có thể thêm `PosConversationLink` |
| `fb_id` | string/null | Facebook ID để link với Pancake | ⭐⭐⭐ Nên thêm `PosFbId` (dùng để identify customer) |
| `customer_id` | UUID string | ID khác với `id` (có thể là internal ID) | ⭐ Có thể bỏ qua (đã có `id`) |

**Lưu ý:**
- `fb_id` rất quan trọng để link customer từ POS với Pancake (qua PSID)
- `notes` cần thiết nếu muốn hiển thị ghi chú khách hàng
- `assigned_user_id` cần thiết nếu muốn quản lý người phụ trách customer

---

### 2. Các Module POS Chưa Có Model

#### ⭐⭐⭐⭐⭐ ƯU TIÊN CAO (Core Modules)

##### 2.1. Orders (Đơn hàng POS)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/orders` - Lấy danh sách đơn hàng
- `GET /shops/{SHOP_ID}/orders/{ORDER_ID}` - Lấy chi tiết đơn hàng
- `GET /shops/{SHOP_ID}/orders_returned` - Lấy đơn hàng đã trả

**Trạng thái:** ❌ Chưa có model `PcPosOrder`

**Cần implement:**
- Model `PcPosOrder` với các field chính:
  - `OrderId`, `SystemId`, `ShopId`
  - `Status`, `StatusName`
  - `BillFullName`, `BillPhoneNumber`, `BillEmail`
  - `CustomerId`, `WarehouseId`
  - `ShippingFee`, `TotalDiscount`
  - `InsertedAt`, `UpdatedAt`, `PaidAt`
  - `OrderItems` (array)
  - `ShippingAddress` (object)
- Service `PcPosOrderService`
- Handler `PcPosOrderHandler`
- Endpoints CRUD + Upsert

**Lý do ưu tiên cao:**
- Đơn hàng là core của hệ thống POS
- Cần thiết cho quản lý bán hàng và báo cáo
- Có nhiều filter và query phức tạp

---

##### 2.2. Products (Sản phẩm)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/products` - Lấy danh sách sản phẩm
- `GET /shops/{SHOP_ID}/products/{PRODUCT_ID}` - Lấy chi tiết sản phẩm
- `GET /shops/{SHOP_ID}/products/{PRODUCT_SKU}` - Lấy sản phẩm theo SKU
- `GET /shops/{SHOP_ID}/products/variations` - Lấy danh sách biến thể
- `GET /shops/{SHOP_ID}/categories` - Lấy danh mục
- `GET /shops/{SHOP_ID}/tags_products` - Lấy tags sản phẩm

**Trạng thái:** ❌ Chưa có model `PcPosProduct`, `PcPosVariation`, `PcPosCategory`

**Cần implement:**
- Model `PcPosProduct`:
  - `ProductId`, `ShopId`, `Name`
  - `CategoryIds`, `TagIds`
  - `IsHide`, `NoteProduct`
  - `ProductAttributes` (array)
  - `Variations` (array hoặc reference)
- Model `PcPosVariation`:
  - `VariationId`, `ProductId`, `ShopId`
  - `Sku`, `RetailPrice`, `PriceAtCounter`
  - `Quantity`, `Weight`
  - `Fields` (array - attributes)
  - `Images` (array)
- Model `PcPosCategory`:
  - `CategoryId`, `ShopId`, `Name`
- Services và Handlers tương ứng

**Lý do ưu tiên cao:**
- Sản phẩm là core của POS
- Cần thiết cho quản lý tồn kho và bán hàng
- Có nhiều biến thể và thuộc tính phức tạp

---

##### 2.3. Customer Levels (Cấp độ khách hàng)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/customer_levels` - Lấy danh sách cấp độ khách hàng

**Trạng thái:** ❌ Chưa có model `PcPosCustomerLevel`

**Cần implement:**
- Model `PcPosCustomerLevel`:
  - `LevelId`, `ShopId`, `Name`
  - `DiscountPercent` (nếu có)
  - `MinOrderAmount` (nếu có)
- Service và Handler tương ứng

**Lý do ưu tiên cao:**
- Cần thiết để hiển thị thông tin cấp độ khách hàng
- Customer model đã có `CustomerLevelId` nhưng chưa có model riêng để lưu thông tin level

---

#### ⭐⭐⭐⭐ ƯU TIÊN TRUNG BÌNH CAO

##### 2.4. Customer Point Logs (Lịch sử điểm tích lũy)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/customers/point_logs` - Lấy lịch sử điểm tích lũy

**Trạng thái:** ❌ Chưa có model `PcPosCustomerPointLog`

**Cần implement:**
- Model `PcPosCustomerPointLog`:
  - `LogId`, `CustomerId`, `ShopId`
  - `PointChange` (số điểm thay đổi, có thể âm)
  - `PointBefore`, `PointAfter`
  - `Reason` (lý do thay đổi)
  - `OrderId` (nếu liên quan đến đơn hàng)
  - `CreatedAt`
- Service và Handler tương ứng

**Lý do ưu tiên trung bình cao:**
- Cần thiết để theo dõi lịch sử điểm tích lũy của khách hàng
- Có thể query theo customer để hiển thị lịch sử

---

#### ⭐⭐⭐ ƯU TIÊN TRUNG BÌNH

##### 2.5. Purchases (Nhập hàng)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/purchases` - Lấy danh sách phiếu nhập
- `GET /shops/{SHOP_ID}/purchases/{PURCHASE_ID}` - Lấy chi tiết
- `GET /shops/{SHOP_ID}/supplier` - Lấy danh sách nhà cung cấp

**Trạng thái:** ❌ Chưa có model `PcPosPurchase`, `PcPosSupplier`

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý nhập hàng
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

##### 2.6. Transfers (Chuyển kho)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/transfers` - Lấy danh sách phiếu chuyển kho
- `GET /shops/{SHOP_ID}/transfers/{TRANSFER_ID}` - Lấy chi tiết

**Trạng thái:** ❌ Chưa có model `PcPosTransfer`

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý chuyển kho
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

##### 2.7. Stocktakings (Kiểm kê)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/stocktakings` - Lấy danh sách phiếu kiểm kê
- `GET /shops/{SHOP_ID}/stocktakings/{STOCKTAKING_ID}` - Lấy chi tiết

**Trạng thái:** ❌ Chưa có model `PcPosStocktaking`

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý kiểm kê
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

##### 2.8. Promotions (Khuyến mãi)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/promotion_advance` - Lấy danh sách khuyến mãi
- `GET /shops/{SHOP_ID}/promotion_advance/{PROMOTION_ID}` - Lấy chi tiết

**Trạng thái:** ❌ Chưa có model `PcPosPromotion`

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý khuyến mãi
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

##### 2.9. Vouchers

**API Endpoints:**
- `GET /shops/{SHOP_ID}/vouchers` - Lấy danh sách voucher
- `GET /shops/{SHOP_ID}/vouchers/{VOUCHER_ID}` - Lấy chi tiết

**Trạng thái:** ❌ Chưa có model `PcPosVoucher`

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý voucher
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

#### ⭐⭐ ƯU TIÊN THẤP

##### 2.10. Analytics (Phân tích)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/analytics/sale` - Phân tích bán hàng
- `GET /shops/{SHOP_ID}/inventory_analytics/inventory` - Phân tích tồn kho

**Trạng thái:** ❌ Chưa có model

**Lý do ưu tiên thấp:**
- Có thể gọi trực tiếp từ Pancake POS API khi cần
- Không cần lưu trữ lâu dài

---

##### 2.11. CRM

**API Endpoints:**
- `GET /shops/{SHOP_ID}/crm/tables` - Lấy danh sách bảng CRM
- `GET /shops/{SHOP_ID}/crm/{TABLE_NAME}/records` - Lấy records

**Trạng thái:** ❌ Chưa có model

**Lý do ưu tiên thấp:**
- Chỉ cần nếu quản lý CRM data
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

##### 2.12. Users (Người dùng POS)

**API Endpoints:**
- `GET /shops/{SHOP_ID}/users` - Lấy danh sách người dùng

**Trạng thái:** ❌ Chưa có model

**Lý do ưu tiên thấp:**
- POS users khác với FolkForm users (Auth module)
- Chỉ cần nếu cần quản lý users của POS
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

## 📊 Bảng Tổng Hợp

| Module/Field | Trạng Thái | Ưu Tiên | Ghi Chú |
|--------------|-----------|---------|---------|
| **Customer - Các field đã có** | ✅ Đã có | - | Đã extract đầy đủ các field chính |
| **Customer - assigned_user_id** | ❌ Chưa có | ⭐⭐⭐ | Nên thêm `PosAssignedUserId` |
| **Customer - notes** | ❌ Chưa có | ⭐⭐⭐ | Nên thêm `PosNotes []interface{}` |
| **Customer - fb_id** | ❌ Chưa có | ⭐⭐⭐ | Nên thêm `PosFbId` (quan trọng để link) |
| **Customer - is_discount_by_level** | ❌ Chưa có | ⭐⭐ | Có thể thêm `PosIsDiscountByLevel` |
| **Customer - conversation_link** | ❌ Chưa có | ⭐⭐ | Có thể thêm `PosConversationLink` |
| **PcPosShop** | ✅ Đã có | - | Đã implement đầy đủ |
| **PcPosWarehouse** | ✅ Đã có | - | Đã implement đầy đủ |
| **PcPosOrder** | ❌ Chưa có | ⭐⭐⭐⭐⭐ | Core module - cần làm ngay |
| **PcPosProduct** | ❌ Chưa có | ⭐⭐⭐⭐⭐ | Core module - cần làm ngay |
| **PcPosVariation** | ❌ Chưa có | ⭐⭐⭐⭐⭐ | Core module - cần làm ngay |
| **PcPosCategory** | ❌ Chưa có | ⭐⭐⭐⭐⭐ | Core module - cần làm ngay |
| **PcPosCustomerLevel** | ❌ Chưa có | ⭐⭐⭐⭐⭐ | Cần để hiển thị thông tin level |
| **PcPosCustomerPointLog** | ❌ Chưa có | ⭐⭐⭐⭐ | Cần để theo dõi lịch sử điểm |
| **PcPosPurchase** | ❌ Chưa có | ⭐⭐⭐ | Chỉ nếu quản lý nhập hàng |
| **PcPosTransfer** | ❌ Chưa có | ⭐⭐⭐ | Chỉ nếu quản lý chuyển kho |
| **PcPosStocktaking** | ❌ Chưa có | ⭐⭐⭐ | Chỉ nếu quản lý kiểm kê |
| **PcPosPromotion** | ❌ Chưa có | ⭐⭐⭐ | Chỉ nếu quản lý khuyến mãi |
| **PcPosVoucher** | ❌ Chưa có | ⭐⭐⭐ | Chỉ nếu quản lý voucher |
| **Analytics** | ❌ Chưa có | ⭐⭐ | Có thể gọi trực tiếp API |
| **CRM** | ❌ Chưa có | ⭐⭐ | Chỉ nếu quản lý CRM data |
| **Users** | ❌ Chưa có | ⭐⭐ | Chỉ nếu quản lý POS users |

---

## 🎯 Kế Hoạch Triển Khai

### Phase 1: Bổ Sung Customer Fields (Ưu tiên cao)

1. **Thêm các field còn thiếu vào Customer model:**
   - `PosFbId` - Quan trọng để link với Pancake
   - `PosAssignedUserId` - ID người phụ trách
   - `PosNotes` - Ghi chú khách hàng
   - `PosIsDiscountByLevel` - Có được giảm giá theo cấp độ
   - `PosConversationLink` - Link conversation

**File cần sửa:**
- `api/core/api/models/mongodb/model.customer.go`

---

### Phase 2: Core Modules (Ưu tiên cao)

1. **Orders Module:**
   - Model `PcPosOrder`
   - Service `PcPosOrderService`
   - Handler `PcPosOrderHandler`
   - Endpoints CRUD + Upsert

2. **Products Module:**
   - Models `PcPosProduct`, `PcPosVariation`, `PcPosCategory`
   - Services và Handlers tương ứng
   - Endpoints CRUD + Upsert

3. **Customer Levels Module:**
   - Model `PcPosCustomerLevel`
   - Service và Handler tương ứng
   - Endpoints CRUD + Upsert

---

### Phase 3: Supporting Modules (Ưu tiên trung bình)

4. **Customer Point Logs:**
   - Model `PcPosCustomerPointLog`
   - Service và Handler tương ứng

5. **Purchases, Transfers, Stocktakings:**
   - Models tương ứng (nếu cần)

6. **Promotions, Vouchers:**
   - Models tương ứng (nếu cần)

---

## 📝 Ghi Chú

1. **Customer Model:**
   - Hiện tại đã có đầy đủ các field chính từ POS
   - Còn thiếu một số field phụ nhưng quan trọng: `fb_id`, `notes`, `assigned_user_id`
   - Nên bổ sung các field này để đầy đủ thông tin

2. **Core Modules:**
   - Orders và Products là 2 module quan trọng nhất cần implement ngay
   - Customer Levels cần thiết để hiển thị thông tin đầy đủ về customer

3. **Supporting Modules:**
   - Các module như Purchases, Transfers, Stocktakings chỉ cần nếu thực sự cần quản lý
   - Có thể lưu trong `panCakeData` nếu không cần query riêng

4. **Pattern Implementation:**
   - Tất cả models đều follow pattern hiện tại:
     - Extract tags với converter
     - Lưu raw data trong `PanCakeData` hoặc `PosData`
     - Index cho các field quan trọng
     - Unique index cho identifier fields

---

## 🔗 Tài Liệu Tham Khảo

- [Pancake POS API Context](./pancake-pos-api-context.md)
- [Customer POS Sync Proposal](./customer-pos-sync-proposal.md)
- [Pancake POS Folkform Sync Analysis](./pancake-pos-folkform-sync-analysis.md)
- [Customer Multi-Source Implementation](./customer-multi-source-implementation.md)

---

**Ngày tạo:** 2025-01-XX  
**Phiên bản:** 1.0  
**Cập nhật lần cuối:** 2025-01-XX
