# Tổng Quan Kiến Trúc Dữ Liệu - Bức Tranh Toàn Cảnh

## 📋 Mục Lục

1. [Tổng Quan Hệ Thống](#tổng-quan-hệ-thống)
2. [Nguồn Dữ Liệu](#nguồn-dữ-liệu)
3. [Cấu Trúc Collections](#cấu-trúc-collections)
4. [Mối Quan Hệ Dữ Liệu](#mối-quan-hệ-dữ-liệu)
5. [Luồng Dữ Liệu](#luồng-dữ-liệu)
6. [Insights Kinh Doanh Tiềm Năng](#insights-kinh-doanh-tiềm-năng)
7. [Cơ Hội Ứng Dụng AI](#cơ-hội-ứng-dụng-ai)
8. [Phân Tích Dữ Liệu Thực Tế](#phân-tích-dữ-liệu-thực-tế)

---

## Tổng Quan Hệ Thống

Hệ thống **Folkform** là một nền tảng tích hợp quản lý khách hàng và bán hàng đa kênh, kết nối:
- **Facebook Pages** (qua Pancake API) - Quản lý tương tác khách hàng trên Facebook
- **Pancake POS** - Hệ thống quản lý bán hàng và kho hàng
- **Pancake API** - Quản lý conversations, messages, customers trên Facebook

### Thống Kê Hiện Tại (từ phân tích MongoDB)

- **Database**: `folkform_auth`
- **Tổng số collections**: 21
- **Tổng số documents**: 932,830
- **Collections lớn nhất**:
  - `fb_message_items`: 834,756 documents
  - `customers`: 33,110 documents
  - `fb_conversations`: 26,832 documents
  - `fb_messages`: 26,813 documents
  - `pc_pos_orders`: 2,633 documents
  - `pc_pos_products`: 401 documents

---

## Nguồn Dữ Liệu

### 1. Pancake API (Facebook Integration)

**Base URLs:**
- User's API: `https://pages.fm/api/v1`
- Page's API v1: `https://pages.fm/api/public_api/v1`
- Page's API v2: `https://pages.fm/api/public_api/v2`

**Dữ liệu thu thập:**
- **Pages**: Thông tin Facebook Pages
- **Conversations**: Cuộc hội thoại với khách hàng
- **Messages**: Tin nhắn trong conversations
- **Posts**: Bài đăng trên Facebook
- **Customers**: Thông tin khách hàng từ Facebook (PSID, name, phone, email)

**Đặc điểm:**
- Sử dụng `page_access_token` để xác thực
- Hỗ trợ pagination
- Real-time sync conversations và messages

### 2. Pancake POS API

**Base URL:** `https://pos.pages.fm/api/v1`

**Dữ liệu thu thập:**
- **Shops**: Thông tin cửa hàng
- **Warehouses**: Kho hàng
- **Products**: Sản phẩm và biến thể
- **Orders**: Đơn hàng từ POS
- **Customers**: Khách hàng trong hệ thống POS
- **Categories**: Danh mục sản phẩm

**Đặc điểm:**
- Sử dụng `api_key` để xác thực
- Tất cả endpoints theo format: `/shops/{SHOP_ID}/...`
- Hỗ trợ pagination và filtering

---

## Cấu Trúc Collections

### 1. Facebook Collections

#### `fb_pages`
**Mục đích**: Lưu thông tin Facebook Pages được kết nối

**Cấu trúc:**
```go
type FbPage struct {
    PageId          string  // Facebook Page ID (unique)
    PageName        string  // Tên trang
    PageUsername    string  // Username của trang
    IsSync          bool    // Trạng thái đồng bộ
    AccessToken     string  // Access token
    PageAccessToken string  // Page access token
    PanCakeData     map[string]interface{} // Dữ liệu gốc từ API
}
```

**Mối quan hệ:**
- 1 Page → N Conversations
- 1 Page → N Messages
- 1 Page → N Posts
- 1 Page → N Customers (qua PSID)

#### `fb_conversations`
**Mục đích**: Lưu metadata của các cuộc hội thoại

**Cấu trúc:**
```go
type FbConversation struct {
    ConversationId   string  // Facebook Conversation ID (unique)
    PageId          string  // Reference to FbPage
    CustomerId      string  // Facebook Customer ID
    PanCakeData     map[string]interface{} // Dữ liệu từ Pancake API
    PanCakeUpdatedAt int64  // Thời gian cập nhật từ Pancake
}
```

**Thống kê:**
- 26,832 conversations
- Top page: `109383448131220` với 10,062 conversations

**Mối quan hệ:**
- 1 Conversation → N Messages (qua `fb_message_items`)
- 1 Conversation → 1 Customer (qua `customerId`)
- 1 Conversation → 1 Page (qua `pageId`)

#### `fb_messages`
**Mục đích**: Lưu metadata của conversation (KHÔNG lưu messages[])

**Cấu trúc:**
```go
type FbMessage struct {
    ConversationId  string  // Reference to FbConversation
    PageId          string  // Reference to FbPage
    CustomerId      string  // Facebook Customer ID
    TotalMessages   int64   // Tổng số messages trong fb_message_items
    HasMore         bool    // Còn messages để sync không
    LastSyncedAt    int64   // Thời gian sync cuối cùng
    PanCakeData     map[string]interface{} // KHÔNG có messages[]
}
```

**Thống kê:**
- 26,813 message metadata records
- Top page: `109383448131220` với 10,056 messages

**Kiến trúc:**
- Messages được tách riêng vào `fb_message_items` để tránh document quá lớn
- `fb_messages` chỉ lưu metadata để query nhanh

#### `fb_message_items`
**Mục đích**: Lưu từng message riêng lẻ (1 message = 1 document)

**Cấu trúc:**
```go
type FbMessageItem struct {
    MessageId       string  // Message ID từ Pancake (unique)
    ConversationId  string  // Reference to FbConversation
    MessageData     map[string]interface{} // Toàn bộ dữ liệu message
    InsertedAt      int64   // Thời gian insert message
}
```

**Thống kê:**
- 834,756 message items (collection lớn nhất)
- Trung bình ~31 messages/conversation

**Mối quan hệ:**
- N Messages → 1 Conversation (qua `conversationId`)

#### `fb_posts`
**Mục đích**: Lưu thông tin bài đăng trên Facebook

**Cấu trúc:**
```go
type FbPost struct {
    PostId      string  // Facebook Post ID (unique)
    PageId      string  // Reference to FbPage
    InsertedAt  int64   // Thời gian insert bài viết
    PanCakeData map[string]interface{} // Dữ liệu từ Pancake API
}
```

**Thống kê:**
- 5,249 posts

**Mối quan hệ:**
- N Posts → 1 Page (qua `pageId`)
- 1 Post → N Conversations (comments trên post)

### 2. Pancake POS Collections

#### `pc_pos_shops`
**Mục đích**: Lưu thông tin cửa hàng từ Pancake POS

**Cấu trúc:**
```go
type PcPosShop struct {
    ShopId      int64   // Shop ID từ POS (unique)
    Name        string  // Tên cửa hàng
    AvatarUrl   string  // Link hình đại diện
    Pages       []interface{} // Thông tin các pages được gộp trong shop
    PanCakeData map[string]interface{} // Dữ liệu gốc từ API
}
```

**Thống kê:**
- 1 shop hiện tại: `860225178`

**Mối quan hệ:**
- 1 Shop → N Warehouses
- 1 Shop → N Products
- 1 Shop → N Orders
- 1 Shop → N Customers

#### `pc_pos_warehouses`
**Mục đích**: Lưu thông tin kho hàng

**Cấu trúc:**
```go
type PcPosWarehouse struct {
    WarehouseId string  // Warehouse ID (UUID)
    ShopId      int64   // Reference to Shop
    Name        string  // Tên kho hàng
    PhoneNumber string  // Số điện thoại
    FullAddress string  // Địa chỉ đầy đủ
    ProvinceId  string  // ID tỉnh/thành phố
    DistrictId  string  // ID quận/huyện
    CommuneId   string  // ID phường/xã
}
```

**Thống kê:**
- 13 warehouses

**Mối quan hệ:**
- N Warehouses → 1 Shop
- 1 Warehouse → N Orders (đơn hàng xuất từ kho)

#### `pc_pos_products`
**Mục đích**: Lưu thông tin sản phẩm

**Cấu trúc:**
```go
type PcPosProduct struct {
    ProductId         string  // Product ID (UUID)
    ShopId            int64   // Reference to Shop
    Name              string  // Tên sản phẩm
    CategoryIds       []int64 // Danh sách ID danh mục
    TagIds            []int64 // Danh sách ID tags
    IsHide            bool    // Trạng thái ẩn/hiện
    NoteProduct       string  // Ghi chú sản phẩm
    ProductAttributes []interface{} // Thuộc tính sản phẩm
    PosData           map[string]interface{} // Dữ liệu gốc
}
```

**Thống kê:**
- 401 products
- Tất cả từ shop `860225178`

**Mối quan hệ:**
- N Products → 1 Shop
- 1 Product → N Variations
- 1 Product → N OrderItems (trong orders)

#### `pc_pos_variations`
**Mục đích**: Lưu thông tin biến thể sản phẩm (màu, size, ...)

**Cấu trúc:**
```go
type PcPosVariation struct {
    VariationId    string  // Variation ID (UUID, unique)
    ProductId      string  // Reference to Product
    ShopId         int64   // Reference to Shop
    Sku            string  // Mã SKU
    RetailPrice    float64 // Giá bán lẻ
    PriceAtCounter float64 // Giá tại quầy
    Quantity       int64   // Số lượng tồn kho
    Weight         float64 // Trọng lượng
    Fields         []interface{} // Các trường thuộc tính (màu, size)
    Images         []string // Danh sách hình ảnh
}
```

**Thống kê:**
- 2,820 variations
- Trung bình ~7 variations/product

**Mối quan hệ:**
- N Variations → 1 Product
- 1 Variation → N OrderItems (trong orders)

#### `pc_pos_orders`
**Mục đích**: Lưu thông tin đơn hàng từ POS

**Cấu trúc:**
```go
type PcPosOrder struct {
    OrderId         int64   // Order ID từ POS
    SystemId        int64   // System ID
    ShopId          int64   // Reference to Shop
    Status          int     // Trạng thái đơn hàng (0-17)
    StatusName      string  // Tên trạng thái
    BillFullName    string  // Tên người thanh toán
    BillPhoneNumber string  // Số điện thoại
    BillEmail       string  // Email
    CustomerId      string  // Reference to Customer (UUID)
    WarehouseId     string  // Reference to Warehouse
    ShippingFee     float64 // Phí vận chuyển
    TotalDiscount   float64 // Tổng giảm giá
    PageId          string  // Facebook Page ID (nếu đơn từ Facebook)
    PostId          string  // Facebook Post ID (nếu đơn từ post)
    OrderItems      []interface{} // Danh sách sản phẩm
    ShippingAddress map[string]interface{} // Địa chỉ giao hàng
    WarehouseInfo   map[string]interface{} // Thông tin kho
    CustomerInfo    map[string]interface{} // Thông tin khách hàng
    PosData         map[string]interface{} // Dữ liệu gốc
}
```

**Thống kê:**
- 2,633 orders
- Tất cả có `status = 0` (Mới)
- Tất cả từ shop `860225178`

**Trạng thái đơn hàng:**
- 0: Mới
- 1: Đã xác nhận
- 2: Đã giao hàng
- 3: Đã nhận hàng
- 4: Đang trả hàng
- 5: Đã trả hàng
- 6: Đã hủy
- ... (xem chi tiết trong tài liệu POS API)

**Mối quan hệ:**
- N Orders → 1 Shop
- N Orders → 1 Customer (qua `customerId`)
- N Orders → 1 Warehouse
- 1 Order → N OrderItems (sản phẩm trong đơn)
- 1 Order → 0..1 Page (nếu đơn từ Facebook)
- 1 Order → 0..1 Post (nếu đơn từ post)

#### `pc_pos_categories`
**Mục đích**: Lưu danh mục sản phẩm

**Cấu trúc:**
```go
type PcPosCategory struct {
    CategoryId int64   // Category ID
    ShopId     int64   // Reference to Shop
    Name       string  // Tên danh mục
    PosData    map[string]interface{} // Dữ liệu gốc
}
```

**Thống kê:**
- 0 categories (chưa có dữ liệu)

**Mối quan hệ:**
- N Categories → 1 Shop
- 1 Category → N Products

### 3. Customer Collection (Multi-Source)

#### `customers`
**Mục đích**: Lưu thông tin khách hàng từ nhiều nguồn (Pancake + POS)

**Cấu trúc:**
```go
type Customer struct {
    // Common Fields (merge từ nhiều nguồn)
    CustomerId   string   // ID chung (unique, sparse)
    Name         string   // Ưu tiên POS > Pancake
    PhoneNumbers []string // Merge từ tất cả nguồn
    Email        string   // Ưu tiên POS > Pancake
    
    // Source-Specific Identifiers
    PanCakeCustomerId string // Pancake Customer ID
    Psid              string // Facebook PSID
    PageId            string // Facebook Page ID
    PosCustomerId     string // POS Customer ID (UUID, unique, sparse)
    
    // Extracted Fields
    Birthday          string // Ngày sinh
    Gender            string // Giới tính
    LivesIn           string // Nơi ở (Pancake)
    CustomerLevelId   string // Cấp độ khách hàng (POS)
    Point             int64  // Điểm tích lũy (POS)
    TotalOrder        int64  // Tổng đơn hàng (POS)
    TotalSpent        float64 // Tổng tiền đã mua (POS)
    SucceedOrderCount int64  // Số đơn hàng thành công (POS)
    TagIds            []interface{} // Tags (POS)
    PosLastOrderAt    int64  // Thời gian đơn hàng cuối (POS)
    PosAddresses      []interface{} // Địa chỉ (POS)
    PosReferralCode   string // Mã giới thiệu (POS)
    PosIsBlock        bool   // Trạng thái block (POS)
    
    // Source Data
    PanCakeData map[string]interface{} // Dữ liệu gốc từ Pancake
    PosData     map[string]interface{} // Dữ liệu gốc từ POS
    
    // Metadata
    Sources   []string // ["pancake", "pos"] - Track nguồn dữ liệu
    CreatedAt int64
    UpdatedAt int64
}
```

**Thống kê:**
- 33,110 customers
- Tất cả có `source = null` (cần cập nhật logic phân loại)

**Merge Strategy:**
- **Name**: Ưu tiên POS (priority=1) > Pancake (priority=2)
- **PhoneNumbers**: Merge array từ tất cả nguồn
- **Email**: Ưu tiên POS > Pancake
- **CustomerId**: Ưu tiên POS ID > Pancake ID

**Mối quan hệ:**
- 1 Customer → N Conversations (qua `psid` hoặc `customerId`)
- 1 Customer → N Orders (qua `customerId`)
- 1 Customer → N Messages (qua `customerId`)

---

## Mối Quan Hệ Dữ Liệu

### Sơ Đồ Quan Hệ Tổng Quan

```
┌─────────────┐
│  FbPage     │
│  (5 pages)  │
└──────┬──────┘
       │
       ├──→ N FbConversations (26,832)
       │         │
       │         ├──→ N FbMessageItems (834,756)
       │         │
       │         └──→ 1 Customer
       │
       ├──→ N FbPosts (5,249)
       │
       └──→ N Customers (via PSID)

┌─────────────┐
│ PcPosShop   │
│ (1 shop)    │
└──────┬──────┘
       │
       ├──→ N PcPosWarehouses (13)
       │
       ├──→ N PcPosProducts (401)
       │         │
       │         └──→ N PcPosVariations (2,820)
       │
       ├──→ N PcPosOrders (2,633)
       │         │
       │         ├──→ 1 Customer (via customerId)
       │         ├──→ 1 Warehouse
       │         └──→ N OrderItems → Variations
       │
       └──→ N Customers (via PosCustomerId)

┌─────────────┐
│  Customer   │
│ (33,110)    │
└──────┬──────┘
       │
       ├──→ N FbConversations (via psid/customerId)
       ├──→ N FbMessageItems (via customerId)
       └──→ N PcPosOrders (via customerId)
```

### Mối Quan Hệ Chi Tiết

#### 1. Customer ↔ Facebook Data
- **Customer.psid** ↔ **FbConversation.customerId**
- **Customer.pageId** ↔ **FbPage.pageId**
- **Customer.panCakeCustomerId** ↔ **Pancake API Customer ID**

#### 2. Customer ↔ POS Data
- **Customer.posCustomerId** ↔ **PcPosOrder.customerId**
- **Customer.customerId** (common) ↔ **PcPosOrder.customerId**

#### 3. Order ↔ Facebook
- **PcPosOrder.pageId** ↔ **FbPage.pageId**
- **PcPosOrder.postId** ↔ **FbPost.postId**
- **PcPosOrder.customerId** ↔ **Customer.customerId** ↔ **FbConversation.customerId**

#### 4. Product Hierarchy
- **PcPosShop** → **PcPosProduct** → **PcPosVariation**
- **PcPosOrder.orderItems** → **PcPosVariation**

#### 5. Message Flow
- **FbPage** → **FbConversation** → **FbMessage** (metadata) → **FbMessageItem** (actual messages)

---

## Luồng Dữ Liệu

### 1. Facebook → System

```
Pancake API
    ↓
FbPage (sync pages)
    ↓
FbConversation (sync conversations)
    ↓
FbMessage (metadata) + FbMessageItem (messages)
    ↓
Customer (extract từ conversation participants)
```

### 2. POS → System

```
Pancake POS API
    ↓
PcPosShop (sync shops)
    ↓
PcPosWarehouse (sync warehouses)
    ↓
PcPosProduct + PcPosVariation (sync products)
    ↓
PcPosOrder (sync orders)
    ↓
Customer (extract từ order customer info)
```

### 3. Customer Merge

```
Customer từ Pancake (PSID, name, phone, email)
    +
Customer từ POS (UUID, name, phone, email, points, orders)
    ↓
Customer (merged với priority rules)
```

### 4. Order Attribution

```
PcPosOrder
    ├── pageId → FbPage
    ├── postId → FbPost
    ├── customerId → Customer
    └── orderItems → PcPosVariation → PcPosProduct
```

---

## Insights Kinh Doanh Tiềm Năng

### 1. Customer Insights

#### Customer Lifetime Value (CLV)
- **Dữ liệu cần**: `Customer.totalSpent`, `Customer.totalOrder`, `Customer.succeedOrderCount`
- **Phân tích**: Tính CLV dựa trên tổng tiền đã mua và số đơn hàng
- **Action**: Phân loại khách hàng VIP, thường xuyên, mới

#### Customer Segmentation
- **Dữ liệu cần**: `Customer.totalOrder`, `Customer.totalSpent`, `Customer.point`, `Customer.customerLevelId`
- **Phân tích**: Phân nhóm khách hàng theo giá trị và tần suất mua
- **Action**: Chiến lược marketing phù hợp cho từng nhóm

#### Customer Churn Analysis
- **Dữ liệu cần**: `Customer.posLastOrderAt`, `PcPosOrder.insertedAt`
- **Phân tích**: Xác định khách hàng không mua trong X ngày
- **Action**: Chiến dịch win-back cho khách hàng có nguy cơ rời bỏ

#### Multi-Channel Customer Journey
- **Dữ liệu cần**: `Customer.sources`, `FbConversation`, `PcPosOrder`
- **Phân tích**: Theo dõi hành trình khách hàng từ Facebook → POS
- **Action**: Tối ưu conversion rate từ conversation → order

### 2. Sales Insights

#### Order Analysis
- **Dữ liệu cần**: `PcPosOrder.status`, `PcPosOrder.totalDiscount`, `PcPosOrder.shippingFee`
- **Phân tích**: 
  - Tỷ lệ đơn hàng theo trạng thái
  - Giá trị đơn hàng trung bình
  - Tỷ lệ giảm giá/shipping fee
- **Action**: Tối ưu quy trình xử lý đơn hàng

#### Product Performance
- **Dữ liệu cần**: `PcPosOrder.orderItems`, `PcPosProduct`, `PcPosVariation`
- **Phân tích**:
  - Sản phẩm bán chạy nhất
  - Biến thể phổ biến nhất
  - Sản phẩm có tỷ lệ return cao
- **Action**: Tối ưu inventory, marketing cho sản phẩm hot

#### Revenue Analysis
- **Dữ liệu cần**: `PcPosOrder`, `PcPosOrder.orderItems`
- **Phân tích**:
  - Doanh thu theo ngày/tuần/tháng
  - Doanh thu theo shop/warehouse
  - Doanh thu theo category
- **Action**: Kế hoạch kinh doanh và dự báo

#### Order Source Attribution
- **Dữ liệu cần**: `PcPosOrder.pageId`, `PcPosOrder.postId`
- **Phân tích**:
  - Đơn hàng đến từ Facebook Pages nào
  - Đơn hàng đến từ Posts nào
  - Conversion rate từ conversation → order
- **Action**: Tối ưu marketing trên Facebook

### 3. Inventory Insights

#### Stock Management
- **Dữ liệu cần**: `PcPosVariation.quantity`, `PcPosOrder.orderItems`
- **Phân tích**:
  - Sản phẩm sắp hết hàng
  - Sản phẩm tồn kho lâu
  - Dự báo nhu cầu
- **Action**: Cảnh báo hết hàng, đề xuất nhập hàng

#### Warehouse Performance
- **Dữ liệu cần**: `PcPosWarehouse`, `PcPosOrder.warehouseId`
- **Phân tích**:
  - Kho nào xuất hàng nhiều nhất
  - Kho nào có tồn kho cao nhất
- **Action**: Tối ưu phân bổ hàng hóa

### 4. Customer Service Insights

#### Response Time Analysis
- **Dữ liệu cần**: `FbMessageItem.insertedAt`, `FbConversation`
- **Phân tích**:
  - Thời gian phản hồi trung bình
  - Conversations chưa được trả lời
  - Peak hours cho customer service
- **Action**: Tối ưu đội ngũ CS, auto-response

#### Conversation Quality
- **Dữ liệu cần**: `FbMessageItem.messageData`, `FbConversation`
- **Phân tích**:
  - Sentiment analysis của messages
  - Topics được hỏi nhiều nhất
  - Customer satisfaction score
- **Action**: Cải thiện chất lượng phục vụ

#### Conversion from Conversation to Order
- **Dữ liệu cần**: `FbConversation`, `PcPosOrder.customerId`, `PcPosOrder.pageId`
- **Phân tích**:
  - Tỷ lệ chuyển đổi conversation → order
  - Thời gian từ conversation → order
  - Yếu tố ảnh hưởng đến conversion
- **Action**: Tối ưu sales process

### 5. Marketing Insights

#### Post Performance
- **Dữ liệu cần**: `FbPost`, `FbConversation.postId`, `PcPosOrder.postId`
- **Phân tích**:
  - Posts nào tạo nhiều conversations nhất
  - Posts nào tạo nhiều orders nhất
  - ROI của từng post
- **Action**: Tối ưu nội dung và timing của posts

#### Page Performance
- **Dữ liệu cần**: `FbPage`, `FbConversation`, `PcPosOrder.pageId`
- **Phân tích**:
  - Pages nào có engagement cao nhất
  - Pages nào có conversion tốt nhất
- **Action**: Tập trung marketing vào pages hiệu quả

---

## Cơ Hội Ứng Dụng AI

### 1. Customer Intelligence

#### AI-Powered Customer Segmentation
- **Input**: `Customer.totalSpent`, `Customer.totalOrder`, `Customer.point`, `PcPosOrder`
- **AI Model**: Clustering (K-means, DBSCAN)
- **Output**: Phân nhóm khách hàng tự động với đặc điểm riêng
- **Value**: Marketing cá nhân hóa, chiến lược pricing

#### Predictive Customer Lifetime Value
- **Input**: `Customer.totalSpent`, `Customer.totalOrder`, `Customer.posLastOrderAt`, `PcPosOrder`
- **AI Model**: Regression (Random Forest, XGBoost)
- **Output**: Dự đoán CLV trong tương lai
- **Value**: Ưu tiên nguồn lực cho khách hàng giá trị cao

#### Churn Prediction
- **Input**: `Customer.posLastOrderAt`, `Customer.totalOrder`, `PcPosOrder.insertedAt`
- **AI Model**: Classification (Logistic Regression, Neural Network)
- **Output**: Xác suất khách hàng rời bỏ
- **Value**: Can thiệp sớm để giữ chân khách hàng

#### Customer Matching (Pancake ↔ POS)
- **Input**: `Customer.phoneNumbers`, `Customer.name`, `Customer.email` từ cả 2 nguồn
- **AI Model**: Fuzzy Matching, Entity Resolution
- **Output**: Match khách hàng từ Facebook và POS
- **Value**: Unified customer view, không trùng lặp

### 2. Sales Intelligence

#### Sales Forecasting
- **Input**: `PcPosOrder.insertedAt`, `PcPosOrder.orderItems`, historical data
- **AI Model**: Time Series (ARIMA, Prophet, LSTM)
- **Output**: Dự báo doanh thu trong tương lai
- **Value**: Kế hoạch inventory, marketing budget

#### Product Recommendation
- **Input**: `PcPosOrder.orderItems`, `Customer.totalOrder`, `PcPosProduct`
- **AI Model**: Collaborative Filtering, Content-Based Filtering
- **Output**: Gợi ý sản phẩm cho khách hàng
- **Value**: Tăng cross-sell, upsell

#### Price Optimization
- **Input**: `PcPosVariation.retailPrice`, `PcPosOrder.orderItems`, `PcPosOrder.totalDiscount`
- **AI Model**: Reinforcement Learning, Optimization
- **Output**: Giá tối ưu để maximize revenue
- **Value**: Tăng lợi nhuận

#### Order Status Prediction
- **Input**: `PcPosOrder.status`, `PcPosOrder.insertedAt`, `PcPosOrder.orderItems`
- **AI Model**: Classification
- **Output**: Dự đoán đơn hàng có risk cao (cancel, return)
- **Value**: Can thiệp sớm để giảm tỷ lệ hủy/trả hàng

### 3. Customer Service Intelligence

#### Sentiment Analysis
- **Input**: `FbMessageItem.messageData.message` (text)
- **AI Model**: NLP (BERT, RoBERTa, Vietnamese models)
- **Output**: Sentiment score (positive/negative/neutral)
- **Value**: Phát hiện khách hàng không hài lòng sớm

#### Intent Classification
- **Input**: `FbMessageItem.messageData.message`
- **AI Model**: Text Classification (BERT-based)
- **Output**: Intent của khách hàng (hỏi giá, khiếu nại, đặt hàng, ...)
- **Value**: Route đến đúng bộ phận, auto-response

#### Auto-Response Generation
- **Input**: `FbMessageItem.messageData.message`, `FbConversation`, `Customer`
- **AI Model**: LLM (GPT, Claude, Vietnamese LLM)
- **Output**: Câu trả lời tự động phù hợp
- **Value**: Giảm workload CS, phản hồi nhanh

#### Conversation Quality Scoring
- **Input**: `FbMessageItem`, `FbConversation`, response time
- **AI Model**: Multi-factor scoring
- **Output**: Quality score cho conversation
- **Value**: Đánh giá hiệu quả CS team

#### Lead Scoring
- **Input**: `FbConversation`, `FbMessageItem`, `Customer`
- **AI Model**: Classification
- **Output**: Score khả năng chuyển đổi thành đơn hàng
- **Value**: Ưu tiên follow-up cho leads chất lượng

### 4. Marketing Intelligence

#### Content Performance Prediction
- **Input**: `FbPost`, historical performance
- **AI Model**: Regression, Classification
- **Output**: Dự đoán engagement/conversion của post
- **Value**: Tối ưu nội dung trước khi đăng

#### Optimal Posting Time
- **Input**: `FbPost.insertedAt`, `FbConversation.insertedAt`, engagement data
- **AI Model**: Time Series Analysis
- **Output**: Thời điểm đăng post tốt nhất
- **Value**: Tăng reach và engagement

#### Customer Journey Mapping
- **Input**: `FbConversation`, `FbMessageItem`, `PcPosOrder`, timestamps
- **AI Model**: Sequence Analysis, Graph Neural Networks
- **Output**: Map hành trình khách hàng từ awareness → purchase
- **Value**: Tối ưu touchpoints, giảm friction

#### A/B Testing Automation
- **Input**: `FbPost`, `FbConversation`, conversion data
- **AI Model**: Multi-armed Bandit, Bayesian Optimization
- **Output**: Tự động chọn variant tốt nhất
- **Value**: Tối ưu marketing campaigns tự động

### 5. Operational Intelligence

#### Inventory Demand Forecasting
- **Input**: `PcPosVariation.quantity`, `PcPosOrder.orderItems`, historical sales
- **AI Model**: Time Series Forecasting
- **Output**: Dự báo nhu cầu sản phẩm
- **Value**: Tối ưu inventory, giảm stockout/overstock

#### Anomaly Detection
- **Input**: `PcPosOrder`, `PcPosVariation.quantity`, `FbMessageItem`
- **AI Model**: Isolation Forest, Autoencoders
- **Output**: Phát hiện bất thường (đơn hàng lạ, inventory bất thường)
- **Value**: Phát hiện fraud, lỗi hệ thống sớm

#### Route Optimization (nếu có delivery)
- **Input**: `PcPosOrder.shippingAddress`, `PcPosWarehouse`
- **AI Model**: Optimization algorithms (TSP, VRP)
- **Output**: Route tối ưu cho delivery
- **Value**: Giảm chi phí vận chuyển

### 6. Data Quality & Integration

#### Data Cleaning & Deduplication
- **Input**: `Customer` từ nhiều nguồn
- **AI Model**: Entity Resolution, Fuzzy Matching
- **Output**: Customer records sạch, không trùng lặp
- **Value**: Data quality cao, insights chính xác

#### Missing Data Imputation
- **Input**: `Customer`, `PcPosOrder` với missing fields
- **AI Model**: Imputation (KNN, MICE, Deep Learning)
- **Output**: Điền đầy đủ dữ liệu thiếu
- **Value**: Phân tích đầy đủ hơn

---

## Phân Tích Dữ Liệu Thực Tế

### Thống Kê Hiện Tại

#### Collections Overview
- **fb_message_items**: 834,756 (89.5% tổng documents)
- **customers**: 33,110 (3.5%)
- **fb_conversations**: 26,832 (2.9%)
- **fb_messages**: 26,813 (2.9%)
- **fb_posts**: 5,249 (0.6%)
- **pc_pos_variations**: 2,820 (0.3%)
- **pc_pos_orders**: 2,633 (0.3%)
- **pc_pos_products**: 401 (0.04%)

#### Key Observations

1. **Message Volume**: 
   - 834K messages từ 26K conversations
   - Trung bình ~31 messages/conversation
   - Cho thấy conversations có độ sâu tốt

2. **Customer Base**:
   - 33K customers nhưng chỉ 2.6K orders
   - Tỷ lệ conversion: ~8% (cần cải thiện)
   - Tất cả customers có `source = null` → cần fix logic

3. **Order Status**:
   - Tất cả orders có `status = 0` (Mới)
   - Có thể là dữ liệu test hoặc cần sync status

4. **Product Catalog**:
   - 401 products với 2,820 variations
   - Trung bình ~7 variations/product
   - Cho thấy sản phẩm đa dạng về biến thể

5. **Facebook Engagement**:
   - 5 pages active
   - Top page: `109383448131220` với 10K+ conversations
   - Cho thấy tập trung vào một số pages chính

### Gaps & Opportunities

1. **Data Quality**:
   - Customer `source` field chưa được populate
   - Order status chưa được sync đầy đủ
   - Cần validation và cleaning

2. **Data Integration**:
   - Customer matching giữa Pancake và POS chưa rõ
   - Order attribution to Facebook posts chưa được track đầy đủ

3. **Analytics Ready**:
   - Dữ liệu đã đủ để phân tích cơ bản
   - Cần thêm calculated fields cho analytics
   - Cần time-series data cho forecasting

---

## Kết Luận

Hệ thống Folkform có một kiến trúc dữ liệu mạnh mẽ với:
- **Multi-source integration**: Pancake API + POS API
- **Scalable architecture**: Tách messages riêng để tránh document quá lớn
- **Rich data**: Đủ dữ liệu cho nhiều loại phân tích
- **Clear relationships**: Mối quan hệ rõ ràng giữa các entities

**Cơ hội lớn nhất**:
1. **AI-Powered Customer Intelligence**: Phân tích và dự đoán hành vi khách hàng
2. **Sales Optimization**: Tối ưu conversion, pricing, inventory
3. **Customer Service Automation**: AI chatbot, sentiment analysis
4. **Marketing Intelligence**: Content optimization, journey mapping

**Next Steps**:
1. Fix data quality issues (source field, order status)
2. Implement customer matching algorithm
3. Build analytics dashboard với insights cơ bản
4. Pilot AI models cho use cases ưu tiên
5. Scale AI solutions dựa trên ROI

---

## Tài Liệu Tham Khảo

- [Pancake API Context](./pancake-api-context.md)
- [Pancake POS API Context](./pancake-pos-api-context.md)
- [Customer Multi-Source Implementation](./customer-multi-source-implementation.md)
- [Database Schema](../02-architecture/database.md)




