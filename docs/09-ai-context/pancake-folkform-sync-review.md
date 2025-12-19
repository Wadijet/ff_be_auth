# Rà Soát Đồng Bộ Dữ Liệu Pancake & Folkform

## 📋 Tổng Quan

Tài liệu này rà soát chi tiết các loại dữ liệu cần đồng bộ giữa:
- **Pancake API** (pages.fm) - Quản lý Facebook Pages, Conversations, Messages
- **Pancake POS API** (pos.pages.fm) - Quản lý đơn hàng, sản phẩm, kho hàng
- **FolkForm API** - Hệ thống backend hiện tại

---

## ✅ ĐÃ ĐỒNG BỘ (Đã Implement)

### 1. Pancake API (pages.fm) - Facebook Integration

#### ✅ Facebook Pages (FbPage)
- **Model**: `models.FbPage`
- **Service**: `FbPageService`
- **Handler**: `FbPageHandler`
- **Endpoints**: CRUD + `find-by-page-id`, `update-token`
- **Data Extraction**: Tự động extract `pageId`, `pageName`, `pageUsername` từ `panCakeData`

#### ✅ Facebook Posts (FbPost)
- **Model**: `models.FbPost`
- **Service**: `FbPostService`
- **Handler**: `FbPostHandler`
- **Endpoints**: CRUD + `find-by-post-id`
- **Data Extraction**: Tự động extract `pageId`, `postId`, `insertedAt` từ `panCakeData`

#### ✅ Facebook Conversations (FbConversation)
- **Model**: `models.FbConversation`
- **Service**: `FbConversationService`
- **Handler**: `FbConversationHandler`
- **Endpoints**: CRUD + `sort-by-api-update`
- **Data Extraction**: Tự động extract `conversationId`, `customerId`, `panCakeUpdatedAt` từ `panCakeData`

#### ✅ Facebook Messages (FbMessage + FbMessageItem)
- **Model**: `models.FbMessage`, `models.FbMessageItem`
- **Service**: `FbMessageService`, `FbMessageItemService`
- **Handler**: `FbMessageHandler`, `FbMessageItemHandler`
- **Endpoints**: 
  - CRUD cho FbMessage
  - CRUD cho FbMessageItem
  - `upsert-messages` (tự động tách messages vào collection riêng)
- **Data Extraction**: Tự động extract `conversationId` từ `panCakeData`
- **Đặc biệt**: Logic tự động tách `messages[]` ra khỏi `panCakeData` và lưu vào 2 collections

#### ✅ Pancake Orders (PcOrder)
- **Model**: `models.PcOrder`
- **Service**: `PcOrderService`
- **Handler**: `PcOrderHandler`
- **Endpoints**: CRUD operations
- **Data Extraction**: Tự động extract `pancakeOrderId` từ `panCakeData`

#### ✅ Access Tokens (PcAccessToken)
- **Model**: `models.PcAccessToken`
- **Service**: `PcAccessTokenService`
- **Handler**: `PcAccessTokenHandler`
- **Endpoints**: CRUD operations

---

## ❌ CHƯA ĐỒNG BỘ - Pancake API (pages.fm)

### 1. Statistics Module ⚠️ Ưu tiên trung bình
**Pancake API có:**
- Ads Campaign Statistics
- Ads Statistics
- Customer Engagement Statistics
- Page Statistics
- Tag Statistics
- User Statistics

**FolkForm chưa có:**
- ❌ Model `PcStatistics` hoặc `FbStatistics`
- ❌ Service và Handler để quản lý statistics
- ❌ Endpoints để lưu trữ và truy vấn statistics

**Khuyến nghị:**
- Nếu cần phân tích và báo cáo → Nên implement
- Có thể lưu dưới dạng `panCakeData` với các trường extract như `pageId`, `statType`, `period`

### 2. Customers Module ⚠️ Ưu tiên cao
**Pancake API có:**
- Get Page Customers
- Update Customer
- Add Customer Note
- Update Customer Note
- Delete Customer Note

**FolkForm chưa có:**
- ❌ Model `PcCustomer` hoặc `FbCustomer`
- ❌ Service và Handler để quản lý customers
- ❌ Endpoints để lưu trữ customer data từ Pancake

**Khuyến nghị:**
- **Nên implement** vì customer data quan trọng cho CRM và phân tích
- Model nên có: `customerId`, `pageId`, `name`, `phone`, `email`, `panCakeData`
- Cần extract: `psid`, `name`, `phone_numbers`, `email`, `birthday`, `gender`, `lives_in`

### 3. Export Data Module ⚠️ Ưu tiên thấp
**Pancake API có:**
- Export Conversations from Ads

**FolkForm chưa có:**
- ❌ Endpoint để trigger export từ Pancake
- ❌ Endpoint để nhận và lưu trữ exported data

**Khuyến nghị:**
- Có thể không cần nếu đã có sync conversations thông qua API thông thường
- Nếu cần export hàng loạt → Có thể implement như một job/background task

### 4. Call Logs Module ⚠️ Ưu tiên thấp
**Pancake API có:**
- Retrieve Call Logs (SIP Call Logs)

**FolkForm chưa có:**
- ❌ Model `PcCallLog` hoặc `SipCallLog`
- ❌ Service và Handler để quản lý call logs
- ❌ Endpoints để lưu trữ call logs từ Pancake

**Khuyến nghị:**
- Chỉ cần nếu tích hợp SIP/VoIP
- Nếu không dùng SIP → Có thể bỏ qua

### 5. Tags Module ⚠️ Ưu tiên thấp
**Pancake API có:**
- Get List Tags
- Tag Conversation (đã có trong Conversation API)

**FolkForm chưa có:**
- ❌ Model `PcTag` hoặc `FbTag`
- ❌ Service và Handler để quản lý tags
- ❌ Endpoints để lưu trữ tags từ Pancake

**Khuyến nghị:**
- Tags có thể lưu trong `panCakeData` của conversations
- Nếu cần query/filter theo tags → Nên implement riêng
- Model nên có: `tagId`, `pageId`, `text`, `color`, `lightenColor`

### 6. Users Module ⚠️ Ưu tiên thấp
**Pancake API có:**
- Get List of Users
- Update Round Robin Users

**FolkForm chưa có:**
- ❌ Model `PcUser` (khác với User trong Auth module)
- ❌ Service và Handler để quản lý Pancake users
- ❌ Endpoints để lưu trữ Pancake user data

**Khuyến nghị:**
- Pancake users khác với FolkForm users (Auth module)
- Chỉ cần nếu cần quản lý users của Pancake (assign conversations, round robin)
- Có thể lưu trong `panCakeData` nếu không cần query riêng

### 7. Page's Contents Module ⚠️ Ưu tiên thấp
**Pancake API có:**
- Upload Media Content

**FolkForm chưa có:**
- ❌ Model `PcContent` hoặc `FbContent`
- ❌ Service và Handler để quản lý uploaded contents
- ❌ Endpoints để lưu trữ content metadata từ Pancake

**Khuyến nghị:**
- Chỉ cần lưu metadata (content_id, attachment_type)
- File thực tế được lưu trên Pancake/CDN
- Có thể lưu trong `panCakeData` của messages nếu không cần query riêng

### 8. Conversation Actions ⚠️ Ưu tiên trung bình
**Pancake API có:**
- Tag Conversation
- Assign Conversation
- Mark as Read
- Mark as Unread

**FolkForm chưa có:**
- ❌ Endpoints để thực hiện các actions này

**Khuyến nghị:**
- Có thể gọi Pancake API trực tiếp từ frontend
- Hoặc tạo proxy endpoints trong FolkForm để gọi Pancake API
- Nếu cần lưu trạng thái → Cập nhật vào `FbConversation` model

### 9. Webhooks ⚠️ Ưu tiên cao
**Pancake API:**
- ❓ Cần kiểm tra Pancake có hỗ trợ webhook không

**FolkForm chưa có:**
- ❌ Webhook endpoints để nhận dữ liệu từ Pancake
- ❌ Webhook verification middleware

**Khuyến nghị:**
- Nếu Pancake hỗ trợ webhook → Nên implement để real-time sync
- Webhook handlers sẽ gọi `Upsert()` với filter và data từ Pancake
- Cần middleware để verify webhook signature (nếu Pancake hỗ trợ)

---

## ❌ CHƯA ĐỒNG BỘ - Pancake POS API (pos.pages.fm)

### 1. Shop (Cửa hàng) ⚠️ Ưu tiên cao
**Pancake POS API có:**
- Get Shops
- Get Shop Details

**FolkForm chưa có:**
- ❌ Model `PcPosShop`
- ❌ Service và Handler để quản lý shops
- ❌ Endpoints để lưu trữ shop data

**Khuyến nghị:**
- **Nên implement** vì shop là entity cơ bản trong POS
- Model nên có: `shopId`, `name`, `avatarUrl`, `panCakeData`
- Cần extract: `id`, `name`, `avatar_url`, `pages[]`

### 2. Geo (Địa lý) ⚠️ Ưu tiên thấp
**Pancake POS API có:**
- Get Provinces
- Get Districts
- Get Communes

**FolkForm chưa có:**
- ❌ Model `PcGeoProvince`, `PcGeoDistrict`, `PcGeoCommune`
- ❌ Service và Handler để quản lý địa lý
- ❌ Endpoints để lưu trữ địa lý data

**Khuyến nghị:**
- Có thể cache tạm thời hoặc gọi trực tiếp từ Pancake POS API
- Chỉ cần implement nếu cần query/filter theo địa lý thường xuyên
- Hoặc có thể lưu trong `panCakeData` của orders/customers

### 3. Warehouses (Kho hàng) ⚠️ Ưu tiên cao
**Pancake POS API có:**
- Get Warehouses
- Get Warehouse Details
- Get Inventory Histories

**FolkForm chưa có:**
- ❌ Model `PcPosWarehouse`
- ❌ Service và Handler để quản lý warehouses
- ❌ Endpoints để lưu trữ warehouse data

**Khuyến nghị:**
- **Nên implement** nếu cần quản lý tồn kho
- Model nên có: `warehouseId`, `shopId`, `name`, `address`, `panCakeData`
- Cần extract: `id`, `name`, `phone_number`, `full_address`, `province_id`, `district_id`, `commune_id`

### 4. Orders (Đơn hàng POS) ⚠️ Ưu tiên cao
**Pancake POS API có:**
- Get Orders (với nhiều filter)
- Get Order Details
- Get Order Sources
- Get Order Tags
- Get Tracking URL
- Get Returned Orders

**FolkForm chưa có:**
- ❌ Model `PcPosOrder` (khác với `PcOrder` từ Pancake API)
- ❌ Service và Handler để quản lý POS orders
- ❌ Endpoints để lưu trữ POS order data

**Khuyến nghị:**
- **Nên implement** vì đơn hàng là core của POS
- Model nên có: `orderId`, `shopId`, `status`, `customerId`, `panCakeData`
- Cần extract: `id`, `system_id`, `shop_id`, `status`, `inserted_at`, `updated_at`, `bill_full_name`, `bill_phone_number`, `total_discount`, `shipping_fee`, `warehouse_id`, `customer`, `order_items[]`, `shipping_address`

### 5. Customers (Khách hàng POS) ⚠️ Ưu tiên cao
**Pancake POS API có:**
- Get Customers
- Get Customer Details
- Get Point Logs
- Get/Add Customer Notes
- Get Customer Levels

**FolkForm chưa có:**
- ❌ Model `PcPosCustomer` (khác với `PcCustomer` từ Pancake API)
- ❌ Service và Handler để quản lý POS customers
- ❌ Endpoints để lưu trữ POS customer data

**Khuyến nghị:**
- **Nên implement** vì customer là core của CRM
- Model nên có: `customerId`, `shopId`, `name`, `phone`, `email`, `point`, `totalOrder`, `totalSpent`, `panCakeData`
- Cần extract: `id`, `name`, `phone_number`, `email`, `customer_level_id`, `point`, `total_order`, `total_spent`, `tags[]`

### 6. Products (Sản phẩm) ⚠️ Ưu tiên cao
**Pancake POS API có:**
- Get Products
- Create Product
- Get Product Details
- Get Product by SKU
- Update Quantity
- Update Hide Status
- Get Product Tags
- Get Categories
- Get Materials
- Get Measurements

**FolkForm chưa có:**
- ❌ Model `PcPosProduct`, `PcPosVariation`, `PcPosCategory`
- ❌ Service và Handler để quản lý products
- ❌ Endpoints để lưu trữ product data

**Khuyến nghị:**
- **Nên implement** vì sản phẩm là core của POS
- Model nên có: `productId`, `shopId`, `name`, `categoryIds[]`, `tags[]`, `variations[]`, `panCakeData`
- Cần extract: `id`, `name`, `category_ids[]`, `tags[]`, `variations[]` (với `id`, `fields[]`, `images[]`, `retail_price`, `price_at_counter`, `sku`, `quantity`)

### 7. Purchases (Nhập hàng) ⚠️ Ưu tiên trung bình
**Pancake POS API có:**
- Get Purchases
- Get Purchase Details
- Separate Purchase
- Get Suppliers

**FolkForm chưa có:**
- ❌ Model `PcPosPurchase`, `PcPosSupplier`
- ❌ Service và Handler để quản lý purchases
- ❌ Endpoints để lưu trữ purchase data

**Khuyến nghị:**
- Nếu cần quản lý nhập hàng → Nên implement
- Model nên có: `purchaseId`, `shopId`, `supplierId`, `warehouseId`, `status`, `panCakeData`
- Cần extract: `id`, `supplier_id`, `warehouse_id`, `status`, `inserted_at`, `purchase_items[]`

### 8. Transfers (Chuyển kho) ⚠️ Ưu tiên trung bình
**Pancake POS API có:**
- Get Transfers
- Create Transfer
- Get Transfer Details
- Get Transfer Status History

**FolkForm chưa có:**
- ❌ Model `PcPosTransfer`
- ❌ Service và Handler để quản lý transfers
- ❌ Endpoints để lưu trữ transfer data

**Khuyến nghị:**
- Nếu cần quản lý chuyển kho → Nên implement
- Model nên có: `transferId`, `shopId`, `fromWarehouseId`, `toWarehouseId`, `status`, `panCakeData`
- Cần extract: `id`, `from_warehouse_id`, `to_warehouse_id`, `status`, `inserted_at`, `transfer_items[]`

### 9. Stocktakings (Kiểm kê) ⚠️ Ưu tiên trung bình
**Pancake POS API có:**
- Get Stocktakings
- Get Stocktaking Details

**FolkForm chưa có:**
- ❌ Model `PcPosStocktaking`
- ❌ Service và Handler để quản lý stocktakings
- ❌ Endpoints để lưu trữ stocktaking data

**Khuyến nghị:**
- Nếu cần quản lý kiểm kê → Nên implement
- Model nên có: `stocktakingId`, `shopId`, `warehouseId`, `status`, `panCakeData`
- Cần extract: `id`, `warehouse_id`, `status`, `inserted_at`, `stocktaking_items[]`

### 10. Promotions (Khuyến mãi) ⚠️ Ưu tiên trung bình
**Pancake POS API có:**
- Get Promotions
- Get Promotion Details
- Create Multiple Promotions
- Delete Multiple Promotions

**FolkForm chưa có:**
- ❌ Model `PcPosPromotion`
- ❌ Service và Handler để quản lý promotions
- ❌ Endpoints để lưu trữ promotion data

**Khuyến nghị:**
- Nếu cần quản lý khuyến mãi → Nên implement
- Model nên có: `promotionId`, `shopId`, `name`, `status`, `panCakeData`
- Cần extract: `id`, `name`, `status`, `start_date`, `end_date`, `discount_type`, `discount_value`

### 11. Vouchers ⚠️ Ưu tiên trung bình
**Pancake POS API có:**
- Get Vouchers
- Get Voucher Details
- Create Multiple Vouchers

**FolkForm chưa có:**
- ❌ Model `PcPosVoucher`
- ❌ Service và Handler để quản lý vouchers
- ❌ Endpoints để lưu trữ voucher data

**Khuyến nghị:**
- Nếu cần quản lý voucher → Nên implement
- Model nên có: `voucherId`, `shopId`, `code`, `status`, `panCakeData`
- Cần extract: `id`, `code`, `status`, `discount_type`, `discount_value`, `start_date`, `end_date`

### 12. Combo Products ⚠️ Ưu tiên thấp
**Pancake POS API có:**
- Get Combo Products

**FolkForm chưa có:**
- ❌ Model `PcPosComboProduct`
- ❌ Service và Handler để quản lý combo products
- ❌ Endpoints để lưu trữ combo product data

**Khuyến nghị:**
- Có thể lưu trong `panCakeData` của products nếu không cần query riêng
- Nếu cần query/filter combo products → Nên implement riêng

### 13. Analytics (Phân tích) ⚠️ Ưu tiên trung bình
**Pancake POS API có:**
- Sale Analytics
- Inventory Analytics
- Get List Formula
- Get Analytic Fields

**FolkForm chưa có:**
- ❌ Model `PcPosAnalytics`
- ❌ Service và Handler để quản lý analytics
- ❌ Endpoints để lưu trữ analytics data

**Khuyến nghị:**
- Nếu cần lưu trữ và phân tích dữ liệu → Nên implement
- Có thể lưu dưới dạng `panCakeData` với các trường extract như `shopId`, `analyticsType`, `period`, `data`

### 14. Users (Người dùng POS) ⚠️ Ưu tiên thấp
**Pancake POS API có:**
- Get Users

**FolkForm chưa có:**
- ❌ Model `PcPosUser` (khác với User trong Auth module)
- ❌ Service và Handler để quản lý POS users
- ❌ Endpoints để lưu trữ POS user data

**Khuyến nghị:**
- POS users khác với FolkForm users (Auth module)
- Chỉ cần nếu cần quản lý users của POS
- Có thể lưu trong `panCakeData` nếu không cần query riêng

### 15. CRM ⚠️ Ưu tiên trung bình
**Pancake POS API có:**
- Get CRM Tables
- Get CRM Profile
- Get CRM Records
- Get CRM History

**FolkForm chưa có:**
- ❌ Model `PcPosCrmTable`, `PcPosCrmRecord`
- ❌ Service và Handler để quản lý CRM
- ❌ Endpoints để lưu trữ CRM data

**Khuyến nghị:**
- Nếu cần quản lý CRM data → Nên implement
- Model nên có: `tableName`, `shopId`, `recordId`, `panCakeData`
- Cần extract: `id`, `table_name`, `fields[]`, `inserted_at`, `updated_at`

### 16. Các API khác ⚠️ Ưu tiên thấp
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

**FolkForm chưa có:**
- ❌ Các models và services tương ứng

**Khuyến nghị:**
- Chỉ implement nếu thực sự cần
- Có thể lưu trong `panCakeData` của orders/customers nếu không cần query riêng

---

## 📊 Bảng Tổng Hợp

### Pancake API (pages.fm)

| Module | Trạng Thái | Ưu Tiên | Ghi Chú |
|--------|-----------|---------|---------|
| **Pages** | ✅ Đã có | - | Hoàn chỉnh |
| **Posts** | ✅ Đã có | - | Hoàn chỉnh |
| **Conversations** | ✅ Đã có | - | Hoàn chỉnh |
| **Messages** | ✅ Đã có | - | Hoàn chỉnh (có logic tách messages) |
| **Orders** | ✅ Đã có | - | Hoàn chỉnh |
| **Access Tokens** | ✅ Đã có | - | Hoàn chỉnh |
| **Customers** | ❌ Chưa có | ⚠️ Cao | Quan trọng cho CRM |
| **Statistics** | ❌ Chưa có | ⚠️ Trung bình | Nếu cần phân tích |
| **Tags** | ❌ Chưa có | ⚠️ Thấp | Có thể lưu trong panCakeData |
| **Users** | ❌ Chưa có | ⚠️ Thấp | Có thể lưu trong panCakeData |
| **Page Contents** | ❌ Chưa có | ⚠️ Thấp | Có thể lưu trong panCakeData |
| **Call Logs** | ❌ Chưa có | ⚠️ Thấp | Chỉ nếu dùng SIP |
| **Export Data** | ❌ Chưa có | ⚠️ Thấp | Có thể không cần |
| **Conversation Actions** | ❌ Chưa có | ⚠️ Trung bình | Có thể proxy Pancake API |
| **Webhooks** | ❌ Chưa có | ⚠️ Cao | Nếu Pancake hỗ trợ |

### Pancake POS API (pos.pages.fm)

| Module | Trạng Thái | Ưu Tiên | Ghi Chú |
|--------|-----------|---------|---------|
| **Shop** | ❌ Chưa có | ⚠️ Cao | Entity cơ bản |
| **Orders** | ❌ Chưa có | ⚠️ Cao | Core của POS |
| **Customers** | ❌ Chưa có | ⚠️ Cao | Core của CRM |
| **Products** | ❌ Chưa có | ⚠️ Cao | Core của POS |
| **Warehouses** | ❌ Chưa có | ⚠️ Cao | Nếu cần quản lý kho |
| **Purchases** | ❌ Chưa có | ⚠️ Trung bình | Nếu cần quản lý nhập hàng |
| **Transfers** | ❌ Chưa có | ⚠️ Trung bình | Nếu cần quản lý chuyển kho |
| **Stocktakings** | ❌ Chưa có | ⚠️ Trung bình | Nếu cần quản lý kiểm kê |
| **Promotions** | ❌ Chưa có | ⚠️ Trung bình | Nếu cần quản lý khuyến mãi |
| **Vouchers** | ❌ Chưa có | ⚠️ Trung bình | Nếu cần quản lý voucher |
| **Analytics** | ❌ Chưa có | ⚠️ Trung bình | Nếu cần phân tích |
| **CRM** | ❌ Chưa có | ⚠️ Trung bình | Nếu cần quản lý CRM |
| **Geo** | ❌ Chưa có | ⚠️ Thấp | Có thể cache hoặc gọi trực tiếp |
| **Combo Products** | ❌ Chưa có | ⚠️ Thấp | Có thể lưu trong panCakeData |
| **Users** | ❌ Chưa có | ⚠️ Thấp | Có thể lưu trong panCakeData |
| **Các API khác** | ❌ Chưa có | ⚠️ Thấp | Chỉ nếu thực sự cần |

---

## 🎯 Khuyến Nghị Ưu Tiên

### Ưu Tiên Cao (Cần làm ngay)

1. **Pancake API - Customers Module**
   - Quan trọng cho CRM và phân tích
   - Model: `PcCustomer` hoặc `FbCustomer`
   - Extract: `psid`, `name`, `phone_numbers[]`, `email`, `birthday`, `gender`, `lives_in`

2. **Pancake API - Webhooks** (nếu Pancake hỗ trợ)
   - Real-time sync thay vì polling
   - Webhook handlers cho Pages, Posts, Conversations, Messages, Customers

3. **Pancake POS API - Shop Module**
   - Entity cơ bản, cần cho các module khác
   - Model: `PcPosShop`

4. **Pancake POS API - Orders Module**
   - Core của POS, cần thiết cho quản lý đơn hàng
   - Model: `PcPosOrder` (khác với `PcOrder` từ Pancake API)

5. **Pancake POS API - Customers Module**
   - Core của CRM, cần thiết cho quản lý khách hàng
   - Model: `PcPosCustomer` (khác với `PcCustomer` từ Pancake API)

6. **Pancake POS API - Products Module**
   - Core của POS, cần thiết cho quản lý sản phẩm
   - Model: `PcPosProduct`, `PcPosVariation`, `PcPosCategory`

7. **Pancake POS API - Warehouses Module**
   - Cần thiết nếu quản lý tồn kho
   - Model: `PcPosWarehouse`

### Ưu Tiên Trung Bình (Nếu cần)

8. **Pancake API - Statistics Module**
   - Nếu cần phân tích và báo cáo
   - Model: `PcStatistics` hoặc `FbStatistics`

9. **Pancake API - Conversation Actions**
   - Proxy endpoints để gọi Pancake API
   - Hoặc cập nhật trạng thái vào `FbConversation`

10. **Pancake POS API - Purchases, Transfers, Stocktakings**
    - Nếu cần quản lý nhập hàng, chuyển kho, kiểm kê

11. **Pancake POS API - Promotions, Vouchers**
    - Nếu cần quản lý khuyến mãi và voucher

12. **Pancake POS API - Analytics, CRM**
    - Nếu cần phân tích và quản lý CRM

### Ưu Tiên Thấp (Có thể bỏ qua)

13. **Pancake API - Tags, Users, Page Contents, Call Logs, Export Data**
    - Có thể lưu trong `panCakeData` nếu không cần query riêng

14. **Pancake POS API - Geo, Combo Products, Users, Các API khác**
    - Chỉ implement nếu thực sự cần

---

## 📝 Pattern Implementation

### Cách Implement Module Mới

1. **Tạo Model** với struct tag `extract`:
```go
type PcCustomer struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    CustomerId     string                 `json:"customerId" bson:"customerId" index:"unique" extract:"PanCakeData\\.psid,converter=string"`
    PageId        string                 `json:"pageId" bson:"pageId" extract:"PanCakeData\\.page_id,converter=string"`
    Name           string                 `json:"name" bson:"name" extract:"PanCakeData\\.name,converter=string,optional"`
    PhoneNumbers   []string               `json:"phoneNumbers" bson:"phoneNumbers" extract:"PanCakeData\\.phone_numbers,optional"`
    Email          string                 `json:"email" bson:"email" extract:"PanCakeData\\.email,converter=string,optional"`
    PanCakeData    map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
    UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

2. **Tạo Service** kế thừa `BaseServiceMongoImpl`:
```go
type PcCustomerService struct {
    *BaseServiceMongoImpl[models.PcCustomer]
}

func NewPcCustomerService() (*PcCustomerService, error) {
    collection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.PcCustomers)
    if !exist {
        return nil, fmt.Errorf("failed to get pc_customers collection")
    }
    return &PcCustomerService{
        BaseServiceMongoImpl: NewBaseServiceMongo[models.PcCustomer](collection),
    }, nil
}
```

3. **Tạo Handler** với CRUD + Upsert endpoints:
```go
type PcCustomerHandler struct {
    service *services.PcCustomerService
}

func (h *PcCustomerHandler) HandleUpsertOne(c *fiber.Ctx) error {
    // Parse filter từ query string
    // Parse body với panCakeData
    // Gọi service.Upsert() với filter và data
    // Data extraction tự động chạy qua struct tag extract
}
```

4. **Đăng ký Routes** trong `routes.go`:
```go
pcCustomerHandler := handlers.NewPcCustomerHandler(pcCustomerService)
apiV1.Post("/pancake/customer/upsert-one", pcCustomerHandler.HandleUpsertOne)
```

5. **Đăng ký Collection** trong `init.go` và `init.registry.go`

---

## 🔄 Webhook Implementation (Nếu Pancake hỗ trợ)

### Webhook Handler Pattern

```go
type PancakeWebhookHandler struct {
    fbPageService        *services.FbPageService
    fbPostService       *services.FbPostService
    fbConversationService *services.FbConversationService
    fbMessageService    *services.FbMessageService
    pcCustomerService   *services.PcCustomerService
}

func (h *PancakeWebhookHandler) HandlePageWebhook(c *fiber.Ctx) error {
    // 1. Verify webhook signature (nếu có)
    // 2. Parse payload từ Pancake
    // 3. Tạo filter: {"pageId": payload["id"]}
    // 4. Gọi service.Upsert() với filter và payload
    //    - Data extraction tự động chạy qua struct tag extract
}
```

### Webhook Routes

```go
pancakeWebhookHandler := handlers.NewPancakeWebhookHandler(...)
apiV1.Post("/pancake/webhook/page", pancakeWebhookHandler.HandlePageWebhook)
apiV1.Post("/pancake/webhook/post", pancakeWebhookHandler.HandlePostWebhook)
apiV1.Post("/pancake/webhook/conversation", pancakeWebhookHandler.HandleConversationWebhook)
apiV1.Post("/pancake/webhook/message", pancakeWebhookHandler.HandleMessageWebhook)
apiV1.Post("/pancake/webhook/customer", pancakeWebhookHandler.HandleCustomerWebhook)
```

---

## 📚 Tài Liệu Tham Khảo

- [Pancake API Context](./pancake-api-context.md)
- [Pancake POS API Context](./pancake-pos-api-context.md)
- [Pancake Integration Review](./pancake-integration-review.md)
- [FolkForm API Context](./folkform-api-context.md)

---

**Ngày tạo**: 2025-01-XX  
**Phiên bản**: 1.0  
**Cập nhật lần cuối**: 2025-01-XX
