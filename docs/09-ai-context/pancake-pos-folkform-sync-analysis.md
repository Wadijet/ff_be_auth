# Phân Tích Đồng Bộ Pancake POS API với Folkform

## 📋 Tổng Quan

Tài liệu này phân tích chi tiết các module cần đồng bộ từ **Pancake POS API** (`pos.pages.fm/api/v1`) về **Folkform Backend** dựa trên tài liệu `pancake-pos-api-context.md`.

**Base URL Pancake POS API:** `https://pos.pages.fm/api/v1`  
**Authentication:** API Key (truyền qua query parameter `api_key`)

---

## ✅ ĐÃ ĐỒNG BỘ (Hiện Trạng)

### Hiện tại Folkform CHƯA có bất kỳ module nào từ Pancake POS API

**Lý do:**
- `PcOrder` model hiện tại là từ **Pancake API** (pages.fm), không phải từ **Pancake POS API** (pos.pages.fm)
- `Customer` model hiện tại chỉ đồng bộ từ **Pancake API** (Facebook customers), chưa có data từ **Pancake POS API**

**Kết luận:** Cần implement toàn bộ các module từ Pancake POS API nếu muốn tích hợp.

---

## ❌ CHƯA ĐỒNG BỘ - Các Module Cần Implement

### 🎯 ƯU TIÊN CAO (Core Modules - Cần làm ngay)

#### 1. Shop (Cửa hàng) ⭐⭐⭐⭐⭐

**Pancake POS API có:**
- `GET /shops` - Lấy danh sách shops
- `GET /shops/{SHOP_ID}` - Lấy chi tiết shop

**Folkform cần:**
- ❌ Model `PcPosShop`
- ❌ Service `PcPosShopService`
- ❌ Handler `PcPosShopHandler`
- ❌ Endpoints CRUD + Upsert

**Dữ liệu cần extract:**
```go
type PcPosShop struct {
    ID          primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    ShopId      int64                 `json:"shopId" bson:"shopId" index:"unique" extract:"PanCakeData\\.id,converter=int64"`
    Name        string                `json:"name" bson:"name" extract:"PanCakeData\\.name,converter=string,optional"`
    AvatarUrl   string                `json:"avatarUrl" bson:"avatarUrl" extract:"PanCakeData\\.avatar_url,converter=string,optional"`
    Pages       []interface{}         `json:"pages" bson:"pages" extract:"PanCakeData\\.pages,optional"`
    PanCakeData map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt   int64                 `json:"createdAt" bson:"createdAt"`
    UpdatedAt   int64                 `json:"updatedAt" bson:"updatedAt"`
}
```

**Unique Index:** `{shopId: 1}`

**Lý do ưu tiên cao:**
- Shop là entity cơ bản nhất trong POS
- Các module khác đều cần `shopId` để filter
- Cần có shop trước khi sync các module khác

---

#### 2. Orders (Đơn hàng POS) ⭐⭐⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/orders` - Lấy danh sách đơn hàng (với nhiều filter)
- `GET /shops/{SHOP_ID}/orders/{ORDER_ID}` - Lấy chi tiết đơn hàng
- `GET /shops/{SHOP_ID}/order_source` - Lấy nguồn đơn hàng
- `GET /shops/{SHOP_ID}/orders/tags` - Lấy tags đơn hàng
- `GET /shops/{SHOP_ID}/orders/get_tracking_url` - Lấy URL tracking
- `GET /shops/{SHOP_ID}/orders_returned` - Lấy đơn hàng đã trả

**Folkform cần:**
- ❌ Model `PcPosOrder` (khác với `PcOrder` từ Pancake API)
- ❌ Service `PcPosOrderService`
- ❌ Handler `PcPosOrderHandler`
- ❌ Endpoints CRUD + Upsert + các endpoints đặc biệt

**Dữ liệu cần extract:**
```go
type PcPosOrder struct {
    ID              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    OrderId         int64                 `json:"orderId" bson:"orderId" index:"text" extract:"PanCakeData\\.id,converter=int64"`
    SystemId        int64                 `json:"systemId" bson:"systemId" extract:"PanCakeData\\.system_id,converter=int64,optional"`
    ShopId          int64                 `json:"shopId" bson:"shopId" index:"text" extract:"PanCakeData\\.shop_id,converter=int64"`
    Status          int                   `json:"status" bson:"status" extract:"PanCakeData\\.status,converter=int,optional"`
    StatusName      string                `json:"statusName" bson:"statusName" extract:"PanCakeData\\.status_name,converter=string,optional"`
    BillFullName    string                `json:"billFullName" bson:"billFullName" extract:"PanCakeData\\.bill_full_name,converter=string,optional"`
    BillPhoneNumber string                `json:"billPhoneNumber" bson:"billPhoneNumber" extract:"PanCakeData\\.bill_phone_number,converter=string,optional"`
    BillEmail       string                `json:"billEmail" bson:"billEmail" extract:"PanCakeData\\.bill_email,converter=string,optional"`
    CustomerId      int64                 `json:"customerId" bson:"customerId" extract:"PanCakeData\\.customer\\.id,converter=int64,optional"`
    WarehouseId     string                `json:"warehouseId" bson:"warehouseId" extract:"PanCakeData\\.warehouse_id,converter=string,optional"`
    ShippingFee     float64               `json:"shippingFee" bson:"shippingFee" extract:"PanCakeData\\.shipping_fee,converter=number,optional"`
    TotalDiscount   float64               `json:"totalDiscount" bson:"totalDiscount" extract:"PanCakeData\\.total_discount,converter=number,optional"`
    InsertedAt      int64                 `json:"insertedAt" bson:"insertedAt" extract:"PanCakeData\\.inserted_at,converter=time,optional"`
    UpdatedAt       int64                 `json:"updatedAt" bson:"updatedAt" extract:"PanCakeData\\.updated_at,converter=time,optional"`
    PaidAt          int64                 `json:"paidAt" bson:"paidAt" extract:"PanCakeData\\.paid_at,converter=time,optional"`
    PanCakeData     map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt       int64                 `json:"createdAt" bson:"createdAt"`
    UpdatedAt       int64                 `json:"updatedAt" bson:"updatedAt"`
}
```

**Unique Index:** `{orderId: 1, shopId: 1}` (compound unique)

**Lý do ưu tiên cao:**
- Đơn hàng là core của hệ thống POS
- Cần thiết cho quản lý bán hàng và báo cáo
- Có nhiều filter và query phức tạp

---

#### 3. Customers (Khách hàng POS) ⭐⭐⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/customers` - Lấy danh sách khách hàng
- `GET /shops/{SHOP_ID}/customers/{CUSTOMER_ID}` - Lấy chi tiết khách hàng
- `GET /shops/{SHOP_ID}/customers/point_logs` - Lấy lịch sử điểm tích lũy
- `GET /shops/{SHOP_ID}/customers/{CUSTOMER_ID}/load_customer_notes` - Lấy ghi chú
- `POST /shops/{SHOP_ID}/customers/{CUSTOMER_ID}/create_note` - Tạo ghi chú
- `GET /shops/{SHOP_ID}/customer_levels` - Lấy danh sách cấp độ khách hàng

**Folkform cần:**
- ❌ Model `PcPosCustomer` (hoặc mở rộng `Customer` model hiện tại)
- ❌ Service `PcPosCustomerService`
- ❌ Handler `PcPosCustomerHandler`
- ❌ Endpoints CRUD + Upsert

**Lưu ý quan trọng:**
- Hiện tại `Customer` model chỉ có data từ Pancake API (Facebook)
- Có 2 phương án:
  1. **Tách riêng:** Tạo `PcPosCustomer` riêng (đơn giản, rõ ràng)
  2. **Unified:** Mở rộng `Customer` model thêm `PosData` (phức tạp hơn, cần logic merge)

**Khuyến nghị:** Tách riêng `PcPosCustomer` để đơn giản và tương thích với pattern hiện tại.

**Dữ liệu cần extract:**
```go
type PcPosCustomer struct {
    ID              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    CustomerId      int64                 `json:"customerId" bson:"customerId" index:"text" extract:"PanCakeData\\.id,converter=int64"`
    ShopId          int64                 `json:"shopId" bson:"shopId" index:"text" extract:"PanCakeData\\.shop_id,converter=int64"`
    Name            string                `json:"name" bson:"name" extract:"PanCakeData\\.name,converter=string,optional"`
    PhoneNumber     string                `json:"phoneNumber" bson:"phoneNumber" extract:"PanCakeData\\.phone_number,converter=string,optional"`
    Email           string                `json:"email" bson:"email" extract:"PanCakeData\\.email,converter=string,optional"`
    CustomerLevelId  int64                 `json:"customerLevelId" bson:"customerLevelId" extract:"PanCakeData\\.customer_level_id,converter=int64,optional"`
    Point           int64                 `json:"point" bson:"point" extract:"PanCakeData\\.point,converter=int64,optional"`
    TotalOrder      int64                 `json:"totalOrder" bson:"totalOrder" extract:"PanCakeData\\.total_order,converter=int64,optional"`
    TotalSpent      float64               `json:"totalSpent" bson:"totalSpent" extract:"PanCakeData\\.total_spent,converter=number,optional"`
    TagIds          []int64               `json:"tagIds" bson:"tagIds" extract:"PanCakeData\\.tags,optional"`
    PanCakeData     map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt       int64                 `json:"createdAt" bson:"createdAt"`
    UpdatedAt       int64                 `json:"updatedAt" bson:"updatedAt"`
}
```

**Unique Index:** `{customerId: 1, shopId: 1}` (compound unique)

**Lý do ưu tiên cao:**
- Khách hàng là core của CRM
- Cần thiết cho phân tích và marketing
- Có thể link với Facebook Customer qua phone/email (sau này)

---

#### 4. Products (Sản phẩm) ⭐⭐⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/products` - Lấy danh sách sản phẩm
- `POST /shops/{SHOP_ID}/products` - Tạo sản phẩm
- `GET /shops/{SHOP_ID}/products/{PRODUCT_ID}` - Lấy chi tiết sản phẩm
- `GET /shops/{SHOP_ID}/products/{PRODUCT_SKU}` - Lấy sản phẩm theo SKU
- `PUT /shops/{SHOP_ID}/variations/{VARIATION_ID}/update_quantity` - Cập nhật số lượng
- `PUT /shops/{SHOP_ID}/variations/update_quantity` - Cập nhật số lượng nhiều biến thể
- `PUT /shops/{SHOP_ID}/products/update_hide` - Cập nhật trạng thái ẩn/hiện
- `GET /shops/{SHOP_ID}/tags_products` - Lấy tags sản phẩm
- `GET /shops/{SHOP_ID}/categories` - Lấy danh mục
- `GET /shops/{SHOP_ID}/materials_products` - Lấy nguyên liệu
- `GET /shops/{SHOP_ID}/product_measurements/get_measure` - Lấy đơn vị đo lường

**Folkform cần:**
- ❌ Model `PcPosProduct`, `PcPosVariation`, `PcPosCategory`
- ❌ Service `PcPosProductService`, `PcPosVariationService`, `PcPosCategoryService`
- ❌ Handler tương ứng
- ❌ Endpoints CRUD + Upsert

**Dữ liệu cần extract:**
```go
type PcPosProduct struct {
    ID              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    ProductId       int64                 `json:"productId" bson:"productId" index:"text" extract:"PanCakeData\\.id,converter=int64"`
    ShopId          int64                 `json:"shopId" bson:"shopId" index:"text" extract:"PanCakeData\\.shop_id,converter=int64"`
    Name            string                `json:"name" bson:"name" extract:"PanCakeData\\.name,converter=string,optional"`
    CategoryIds     []int64               `json:"categoryIds" bson:"categoryIds" extract:"PanCakeData\\.category_ids,optional"`
    TagIds          []int64               `json:"tagIds" bson:"tagIds" extract:"PanCakeData\\.tags,optional"`
    IsHide          bool                  `json:"isHide" bson:"isHide" extract:"PanCakeData\\.is_hide,converter=bool,optional"`
    NoteProduct     string                `json:"noteProduct" bson:"noteProduct" extract:"PanCakeData\\.note_product,converter=string,optional"`
    PanCakeData     map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt       int64                 `json:"createdAt" bson:"createdAt"`
    UpdatedAt       int64                 `json:"updatedAt" bson:"updatedAt"`
}

type PcPosVariation struct {
    ID              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    VariationId     string                `json:"variationId" bson:"variationId" index:"text" extract:"PanCakeData\\.id,converter=string"`
    ProductId       int64                 `json:"productId" bson:"productId" index:"text" extract:"PanCakeData\\.product_id,converter=int64,optional"`
    ShopId          int64                 `json:"shopId" bson:"shopId" index:"text" extract:"PanCakeData\\.shop_id,converter=int64,optional"`
    Sku             string                `json:"sku" bson:"sku" extract:"PanCakeData\\.sku,converter=string,optional"`
    RetailPrice     float64               `json:"retailPrice" bson:"retailPrice" extract:"PanCakeData\\.retail_price,converter=number,optional"`
    PriceAtCounter  float64               `json:"priceAtCounter" bson:"priceAtCounter" extract:"PanCakeData\\.price_at_counter,converter=number,optional"`
    Quantity        int64                 `json:"quantity" bson:"quantity" extract:"PanCakeData\\.quantity,converter=int64,optional"`
    PanCakeData     map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt       int64                 `json:"createdAt" bson:"createdAt"`
    UpdatedAt       int64                 `json:"updatedAt" bson:"updatedAt"`
}

type PcPosCategory struct {
    ID              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    CategoryId      int64                 `json:"categoryId" bson:"categoryId" index:"text" extract:"PanCakeData\\.id,converter=int64"`
    ShopId          int64                 `json:"shopId" bson:"shopId" index:"text" extract:"PanCakeData\\.shop_id,converter=int64"`
    Name            string                `json:"name" bson:"name" extract:"PanCakeData\\.name,converter=string,optional"`
    PanCakeData     map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt       int64                 `json:"createdAt" bson:"createdAt"`
    UpdatedAt       int64                 `json:"updatedAt" bson:"updatedAt"`
}
```

**Unique Indexes:**
- `PcPosProduct`: `{productId: 1, shopId: 1}` (compound unique)
- `PcPosVariation`: `{variationId: 1}` (unique)
- `PcPosCategory`: `{categoryId: 1, shopId: 1}` (compound unique)

**Lý do ưu tiên cao:**
- Sản phẩm là core của POS
- Cần thiết cho quản lý tồn kho và bán hàng
- Có nhiều biến thể và thuộc tính phức tạp

---

#### 5. Warehouses (Kho hàng) ⭐⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/warehouses` - Lấy danh sách kho hàng
- `GET /shops/{SHOP_ID}/warehouses/{WAREHOUSE_ID}` - Lấy chi tiết kho hàng
- `GET /shops/{SHOP_ID}/inventory_histories` - Lấy lịch sử tồn kho

**Folkform cần:**
- ❌ Model `PcPosWarehouse`
- ❌ Service `PcPosWarehouseService`
- ❌ Handler `PcPosWarehouseHandler`
- ❌ Endpoints CRUD + Upsert

**Dữ liệu cần extract:**
```go
type PcPosWarehouse struct {
    ID              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    WarehouseId     string                `json:"warehouseId" bson:"warehouseId" index:"text" extract:"PanCakeData\\.id,converter=string"`
    ShopId          int64                 `json:"shopId" bson:"shopId" index:"text" extract:"PanCakeData\\.shop_id,converter=int64"`
    Name            string                `json:"name" bson:"name" extract:"PanCakeData\\.name,converter=string,optional"`
    PhoneNumber     string                `json:"phoneNumber" bson:"phoneNumber" extract:"PanCakeData\\.phone_number,converter=string,optional"`
    FullAddress     string                `json:"fullAddress" bson:"fullAddress" extract:"PanCakeData\\.full_address,converter=string,optional"`
    ProvinceId      string                `json:"provinceId" bson:"provinceId" extract:"PanCakeData\\.province_id,converter=string,optional"`
    DistrictId      string                `json:"districtId" bson:"districtId" extract:"PanCakeData\\.district_id,converter=string,optional"`
    CommuneId       string                `json:"communeId" bson:"communeId" extract:"PanCakeData\\.commune_id,converter=string,optional"`
    PanCakeData     map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt       int64                 `json:"createdAt" bson:"createdAt"`
    UpdatedAt       int64                 `json:"updatedAt" bson:"updatedAt"`
}
```

**Unique Index:** `{warehouseId: 1, shopId: 1}` (compound unique)

**Lý do ưu tiên cao:**
- Cần thiết nếu quản lý tồn kho
- Liên quan đến orders và products
- Cần cho quản lý nhập hàng, chuyển kho

---

### ⚠️ ƯU TIÊN TRUNG BÌNH (Nếu cần)

#### 6. Purchases (Nhập hàng) ⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/purchases` - Lấy danh sách phiếu nhập
- `GET /shops/{SHOP_ID}/purchases/{PURCHASE_ID}` - Lấy chi tiết phiếu nhập
- `POST /shops/{SHOP_ID}/purchases/separate` - Tách phiếu nhập
- `GET /shops/{SHOP_ID}/supplier` - Lấy danh sách nhà cung cấp

**Folkform cần:**
- ❌ Model `PcPosPurchase`, `PcPosSupplier`
- ❌ Service và Handler tương ứng

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý nhập hàng
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

#### 7. Transfers (Chuyển kho) ⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/transfers` - Lấy danh sách phiếu chuyển kho
- `POST /shops/{SHOP_ID}/transfers/multi` - Tạo phiếu chuyển kho
- `GET /shops/{SHOP_ID}/transfers/{TRANSFER_ID}` - Lấy chi tiết
- `GET /shops/{SHOP_ID}/transfers/get_status_history/{TRANSFER_ID}` - Lấy lịch sử trạng thái

**Folkform cần:**
- ❌ Model `PcPosTransfer`
- ❌ Service và Handler tương ứng

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý chuyển kho
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

#### 8. Stocktakings (Kiểm kê) ⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/stocktakings` - Lấy danh sách phiếu kiểm kê
- `GET /shops/{SHOP_ID}/stocktakings/{STOCKTAKING_ID}` - Lấy chi tiết

**Folkform cần:**
- ❌ Model `PcPosStocktaking`
- ❌ Service và Handler tương ứng

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý kiểm kê
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

#### 9. Promotions (Khuyến mãi) ⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/promotion_advance` - Lấy danh sách khuyến mãi
- `GET /shops/{SHOP_ID}/promotion_advance/{PROMOTION_ID}` - Lấy chi tiết
- `POST /shops/{SHOP_ID}/promotion_advance/create_multi` - Tạo nhiều khuyến mãi
- `POST /shops/{SHOP_ID}/promotion_advance/delete_multi` - Xóa nhiều khuyến mãi

**Folkform cần:**
- ❌ Model `PcPosPromotion`
- ❌ Service và Handler tương ứng

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý khuyến mãi
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

#### 10. Vouchers ⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/vouchers` - Lấy danh sách voucher
- `GET /shops/{SHOP_ID}/vouchers/{VOUCHER_ID}` - Lấy chi tiết
- `POST /shops/{SHOP_ID}/vouchers/create_multi` - Tạo nhiều voucher

**Folkform cần:**
- ❌ Model `PcPosVoucher`
- ❌ Service và Handler tương ứng

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý voucher
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

#### 11. Analytics (Phân tích) ⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/analytics/sale` - Phân tích bán hàng
- `GET /shops/{SHOP_ID}/analytics/get_list_formula` - Lấy danh sách công thức
- `GET /shops/{SHOP_ID}/analytics/get_analytic_fields` - Lấy các trường phân tích
- `GET /shops/{SHOP_ID}/inventory_analytics/inventory` - Phân tích tồn kho
- `GET /shops/{SHOP_ID}/inventory_analytics/inventory_by_product` - Phân tích tồn kho theo sản phẩm

**Folkform cần:**
- ❌ Model `PcPosAnalytics`
- ❌ Service và Handler tương ứng

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu cần lưu trữ và phân tích dữ liệu
- Có thể gọi trực tiếp từ Pancake POS API khi cần
- Có thể lưu dưới dạng `panCakeData` với các trường extract

---

#### 12. CRM ⭐⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/crm/tables` - Lấy danh sách bảng CRM
- `GET /shops/{SHOP_ID}/crm/profile` - Lấy profile CRM
- `GET /shops/{SHOP_ID}/crm/{TABLE_NAME}/records` - Lấy records từ bảng
- `GET /shops/{SHOP_ID}/crm/{TABLE_NAME}/history` - Lấy lịch sử bảng

**Folkform cần:**
- ❌ Model `PcPosCrmTable`, `PcPosCrmRecord`
- ❌ Service và Handler tương ứng

**Lý do ưu tiên trung bình:**
- Chỉ cần nếu quản lý CRM data
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

### 📉 ƯU TIÊN THẤP (Có thể bỏ qua)

#### 13. Geo (Địa lý) ⭐⭐

**Pancake POS API có:**
- `GET /geo/provinces` - Lấy danh sách tỉnh/thành phố
- `GET /geo/districts?province_id={PROVINCE_ID}` - Lấy danh sách quận/huyện
- `GET /geo/communes?district_id={DISTRICT_ID}` - Lấy danh sách phường/xã

**Khuyến nghị:**
- Có thể cache tạm thời hoặc gọi trực tiếp từ Pancake POS API
- Chỉ cần implement nếu cần query/filter theo địa lý thường xuyên
- Hoặc có thể lưu trong `panCakeData` của orders/customers/warehouses

---

#### 14. Combo Products ⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/combo_products` - Lấy danh sách combo sản phẩm

**Khuyến nghị:**
- Có thể lưu trong `panCakeData` của products nếu không cần query riêng
- Nếu cần query/filter combo products → Nên implement riêng

---

#### 15. Users (Người dùng POS) ⭐⭐

**Pancake POS API có:**
- `GET /shops/{SHOP_ID}/users` - Lấy danh sách người dùng

**Khuyến nghị:**
- POS users khác với FolkForm users (Auth module)
- Chỉ cần nếu cần quản lý users của POS
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

#### 16. Các API khác ⭐

**Pancake POS API có:**
- Logistics Shipping Document
- Bank Payments
- Order Call Laters
- Debt
- Transactions
- Adv Costs
- Payment Histories
- Export Data
- Marketplace Account Info
- Shopee Evaluate/Reverse Order
- Partners
- E-Invoices

**Khuyến nghị:**
- Chỉ implement nếu thực sự cần
- Có thể lưu trong `panCakeData` của orders/customers nếu không cần query riêng

---

## 📊 Bảng Tổng Hợp

| Module | Trạng Thái | Ưu Tiên | Ghi Chú |
|--------|-----------|---------|---------|
| **Shop** | ❌ Chưa có | ⭐⭐⭐⭐⭐ | Entity cơ bản, cần làm đầu tiên |
| **Orders** | ❌ Chưa có | ⭐⭐⭐⭐⭐ | Core của POS |
| **Customers** | ❌ Chưa có | ⭐⭐⭐⭐⭐ | Core của CRM |
| **Products** | ❌ Chưa có | ⭐⭐⭐⭐⭐ | Core của POS |
| **Warehouses** | ❌ Chưa có | ⭐⭐⭐⭐ | Nếu cần quản lý kho |
| **Purchases** | ❌ Chưa có | ⭐⭐⭐ | Nếu cần quản lý nhập hàng |
| **Transfers** | ❌ Chưa có | ⭐⭐⭐ | Nếu cần quản lý chuyển kho |
| **Stocktakings** | ❌ Chưa có | ⭐⭐⭐ | Nếu cần quản lý kiểm kê |
| **Promotions** | ❌ Chưa có | ⭐⭐⭐ | Nếu cần quản lý khuyến mãi |
| **Vouchers** | ❌ Chưa có | ⭐⭐⭐ | Nếu cần quản lý voucher |
| **Analytics** | ❌ Chưa có | ⭐⭐⭐ | Nếu cần phân tích |
| **CRM** | ❌ Chưa có | ⭐⭐⭐ | Nếu cần quản lý CRM |
| **Geo** | ❌ Chưa có | ⭐⭐ | Có thể cache hoặc gọi trực tiếp |
| **Combo Products** | ❌ Chưa có | ⭐⭐ | Có thể lưu trong panCakeData |
| **Users** | ❌ Chưa có | ⭐⭐ | Có thể lưu trong panCakeData |
| **Các API khác** | ❌ Chưa có | ⭐ | Chỉ nếu thực sự cần |

---

## 🎯 Kế Hoạch Implementation

### Phase 1: Core Modules (Ưu tiên cao)

1. **Shop Module**
   - Model `PcPosShop`
   - Service `PcPosShopService`
   - Handler `PcPosShopHandler`
   - Endpoints CRUD + Upsert

2. **Orders Module**
   - Model `PcPosOrder`
   - Service `PcPosOrderService`
   - Handler `PcPosOrderHandler`
   - Endpoints CRUD + Upsert + các endpoints đặc biệt

3. **Customers Module**
   - Model `PcPosCustomer` (tách riêng)
   - Service `PcPosCustomerService`
   - Handler `PcPosCustomerHandler`
   - Endpoints CRUD + Upsert

4. **Products Module**
   - Models `PcPosProduct`, `PcPosVariation`, `PcPosCategory`
   - Services tương ứng
   - Handlers tương ứng
   - Endpoints CRUD + Upsert

5. **Warehouses Module**
   - Model `PcPosWarehouse`
   - Service `PcPosWarehouseService`
   - Handler `PcPosWarehouseHandler`
   - Endpoints CRUD + Upsert

### Phase 2: Supporting Modules (Ưu tiên trung bình)

6. Purchases, Transfers, Stocktakings
7. Promotions, Vouchers
8. Analytics, CRM

### Phase 3: Optional Modules (Ưu tiên thấp)

9. Geo, Combo Products, Users
10. Các API khác (nếu cần)

---

## 📝 Pattern Implementation

### Cách Implement Module Mới

1. **Tạo Model** với struct tag `extract`:
```go
type PcPosShop struct {
    ID          primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    ShopId      int64                 `json:"shopId" bson:"shopId" index:"unique" extract:"PanCakeData\\.id,converter=int64"`
    Name        string                `json:"name" bson:"name" extract:"PanCakeData\\.name,converter=string,optional"`
    PanCakeData map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt   int64                 `json:"createdAt" bson:"createdAt"`
    UpdatedAt   int64                 `json:"updatedAt" bson:"updatedAt"`
}
```

2. **Tạo Service** kế thừa `BaseServiceMongoImpl`:
```go
type PcPosShopService struct {
    *BaseServiceMongoImpl[models.PcPosShop]
}

func NewPcPosShopService() (*PcPosShopService, error) {
    collection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcPosShops)
    if !exist {
        return nil, fmt.Errorf("failed to get pc_pos_shops collection")
    }
    return &PcPosShopService{
        BaseServiceMongoImpl: NewBaseServiceMongo[models.PcPosShop](collection),
    }, nil
}
```

3. **Tạo Handler** với CRUD + Upsert endpoints:
```go
type PcPosShopHandler struct {
    service *services.PcPosShopService
}

func (h *PcPosShopHandler) HandleUpsertOne(c *fiber.Ctx) error {
    // Parse filter từ query string: {"shopId": 123}
    // Parse body với panCakeData
    // Gọi service.Upsert() với filter và data
    // Data extraction tự động chạy qua struct tag extract
}
```

4. **Đăng ký Routes** trong `routes.go`:
```go
pcPosShopHandler := handlers.NewPcPosShopHandler(pcPosShopService)
apiV1.Post("/pancake-pos/shop/upsert-one", pcPosShopHandler.HandleUpsertOne)
```

5. **Đăng ký Collection** trong `init.go` và `init.registry.go`

---

## 🔄 Sync Strategy

### 1. Initial Sync (Lần đầu)

**Ví dụ: Sync Shops**
```bash
# Lấy shops từ Pancake POS API
GET https://pos.pages.fm/api/v1/shops?api_key=YOUR_API_KEY

# Upsert vào FolkForm
POST /api/v1/pancake-pos/shop/upsert-one?filter={"shopId":123}
{
  "panCakeData": { ... }
}
```

### 2. Incremental Sync (Định kỳ)

- Sync data mới/updated từ `inserted_at` hoặc `updated_at`
- Query với filter `updated_at >= last_sync_time`

### 3. Webhook (Nếu Pancake POS hỗ trợ)

- Webhook handlers sẽ gọi `Upsert()` với filter và data từ Pancake POS
- Cần middleware để verify webhook signature (nếu có)

---

## ❓ Câu Hỏi Cần Bàn Bạc

1. **Customer Model Strategy:**
   - Tách riêng `PcPosCustomer` hay mở rộng `Customer` hiện tại?
   - Khuyến nghị: Tách riêng để đơn giản

2. **Sync Frequency:**
   - Real-time (webhook) hay polling định kỳ?
   - Tần suất polling nếu dùng polling?

3. **Data Retention:**
   - Có cần lưu lịch sử thay đổi không?
   - Có cần soft delete không?

4. **Error Handling:**
   - Xử lý lỗi khi sync như thế nào?
   - Có cần retry mechanism không?

5. **Performance:**
   - Có cần cache không?
   - Có cần pagination cho sync không?

---

## 📚 Tài Liệu Tham Khảo

- [Pancake POS API Context](./pancake-pos-api-context.md)
- [Pancake Folkform Sync Review](./pancake-folkform-sync-review.md)
- [Customer Sync Proposal](./customer-sync-proposal.md)
- [FolkForm API Context](./folkform-api-context.md)

---

**Ngày tạo**: 2025-01-XX  
**Phiên bản**: 1.0  
**Cập nhật lần cuối**: 2025-01-XX
