# Phân Tích Cấu Trúc Dữ Liệu Thực Tế

## 📋 Mục Lục

1. [Tổng Quan](#tổng-quan)
2. [Customers - Cấu Trúc Thực Tế](#customers---cấu-trúc-thực-tế)
3. [Pancake POS Orders - Cấu Trúc Thực Tế](#pancake-pos-orders---cấu-trúc-thực-tế)
4. [Facebook Conversations - Cấu Trúc Thực Tế](#facebook-conversations---cấu-trúc-thực-tế)
5. [Facebook Messages - Cấu Trúc Thực Tế](#facebook-messages---cấu-trúc-thực-tế)
6. [So Sánh với API Documentation](#so-sánh-với-api-documentation)
7. [Gaps & Recommendations](#gaps--recommendations)

---

## Tổng Quan

Tài liệu này phân tích cấu trúc dữ liệu **thực tế** trong MongoDB dựa trên documents mẫu đã export, so sánh với tài liệu API nguồn để xác định:
- ✅ Fields đã được sync
- ⚠️ Fields có trong API nhưng chưa sync
- ❌ Fields thiếu quan trọng

**Nguồn dữ liệu mẫu**: `docs/09-ai-context/sample-data/*.json`

---

## Customers - Cấu Trúc Thực Tế

### Cấu Trúc Document Thực Tế

```json
{
  "_id": "ObjectId",
  "customerId": "ef40d9c7-bc33-481a-aa20-f31f618c081b",  // Common ID
  "name": "Nam Tống",
  "gender": "male",
  "pageId": "102039018873979",
  "panCakeCustomerId": "ef40d9c7-bc33-481a-aa20-f31f618c081b",
  "psid": "25765649366371270",
  "posCustomerId": "",  // ⚠️ Thường rỗng - chưa match với POS
  "phoneNumbers": ["[0399808840]"],  // ⚠️ Format lạ: có dấu []
  "sources": null,  // ❌ Chưa được populate
  "panCakeData": {
    "id": "ef40d9c7-bc33-481a-aa20-f31f618c081b",
    "customer_id": "c5219f4a-645c-4c5e-8975-bd5bfc32201e",
    "name": "Nam Tống",
    "psid": "25765649366371270",
    "page_id": "102039018873979",
    "thread_id": "102039018873979_25765649366371270",
    "gender": "male",
    "phone_numbers": ["0399808840"],
    "birthday": null,
    "lives_in": null,
    "can_inbox": true,
    "inserted_at": "2025-10-30T03:20:59",
    "updated_at": "2025-12-08T08:41:09",
    "notes": null,
    "recent_orders": null
  },
  "panCakeUpdatedAt": 1765183269000,
  "createdAt": 1766110022608,
  "updatedAt": 1766110022608
}
```

### So Sánh với Pancake API Documentation

#### ✅ Fields Đã Sync (từ Pancake API)

| Field trong DB | Field trong API | Status | Notes |
|---------------|----------------|--------|-------|
| `panCakeData.id` | `id` | ✅ | Pancake Customer ID |
| `panCakeData.customer_id` | `customer_id` | ✅ | Internal customer ID |
| `panCakeData.name` | `name` | ✅ | Tên khách hàng |
| `panCakeData.psid` | `psid` | ✅ | Facebook PSID |
| `panCakeData.page_id` | `page_id` | ✅ | Facebook Page ID |
| `panCakeData.thread_id` | `thread_id` | ✅ | Thread ID |
| `panCakeData.gender` | `gender` | ✅ | Giới tính |
| `panCakeData.phone_numbers` | `phone_numbers` | ✅ | Số điện thoại (array) |
| `panCakeData.birthday` | `birthday` | ✅ | Ngày sinh (có thể null) |
| `panCakeData.lives_in` | `lives_in` | ✅ | Nơi ở (có thể null) |
| `panCakeData.can_inbox` | `can_inbox` | ✅ | Có thể inbox |
| `panCakeData.inserted_at` | `inserted_at` | ✅ | Thời gian tạo |
| `panCakeData.updated_at` | `updated_at` | ✅ | Thời gian cập nhật |
| `panCakeData.notes` | `notes` | ✅ | Ghi chú (có thể null) |
| `panCakeData.recent_orders` | `recent_orders` | ✅ | Đơn hàng gần đây (có thể null) |

#### ⚠️ Fields Có Trong API Nhưng Chưa Sync Đầy Đủ

| Field trong API | Field trong DB | Status | Notes |
|----------------|----------------|--------|-------|
| `email` | `panCakeData.email` | ⚠️ | Có trong API nhưng không thấy trong sample |
| `tags` | `panCakeData.tags` | ⚠️ | Tags của customer (có thể có trong API) |

#### ❌ Fields Thiếu Quan Trọng

1. **`sources`**: Field này luôn `null` - cần populate để track nguồn dữ liệu
2. **`posCustomerId`**: Thường rỗng - cần logic matching với POS customers
3. **`phoneNumbers` format**: Có dấu `[]` trong string - cần fix extract logic

### So Sánh với Pancake POS API Documentation

#### ✅ Fields Đã Sync (từ POS API - nếu có)

Hiện tại trong sample không có customer nào có `posData`. Điều này cho thấy:
- ❌ Chưa sync customers từ POS API
- ❌ Chưa có logic merge customers giữa Pancake và POS

#### ⚠️ Fields Cần Sync Từ POS API

Theo [Pancake POS API Documentation](./pancake-pos-api-context.md), Customer schema có:

| Field trong POS API | Field trong DB | Status | Priority |
|---------------------|----------------|--------|----------|
| `id` | `posCustomerId` | ❌ | High |
| `name` | `name` (merge) | ⚠️ | High |
| `phone_numbers` | `phoneNumbers` (merge) | ⚠️ | High |
| `emails` | `email` (merge) | ❌ | Medium |
| `date_of_birth` | `birthday` (merge) | ❌ | Medium |
| `gender` | `gender` (merge) | ⚠️ | Medium |
| `level_id` | `customerLevelId` | ❌ | Low |
| `reward_point` | `point` | ❌ | Low |
| `order_count` | `totalOrder` | ❌ | Low |
| `purchased_amount` | `totalSpent` | ❌ | Low |
| `succeed_order_count` | `succeedOrderCount` | ❌ | Low |
| `tags` | `tagIds` | ❌ | Low |
| `last_order_at` | `posLastOrderAt` | ❌ | Medium |
| `shop_customer_addresses` | `posAddresses` | ❌ | Medium |
| `referral_code` | `posReferralCode` | ❌ | Low |
| `is_block` | `posIsBlock` | ❌ | Low |

---

## Pancake POS Orders - Cấu Trúc Thực Tế

### Cấu Trúc Document Thực Tế

```json
{
  "_id": "ObjectId",
  "orderId": 3037,
  "systemId": 3037,
  "shopId": 860225178,
  "status": 0,
  "statusName": "submitted",
  "billFullName": "My Dung Truong",
  "billPhoneNumber": "0944252001",
  "billEmail": "",
  "customerId": "f87b4bd9-5182-4fda-84be-1e6b93ae6208",
  "warehouseId": "29d809c3-b0ad-4aa8-94b3-4e5d7f27175d",
  "pageId": "109383448131220",
  "postId": "109383448131220_122260255748023280",
  "shippingFee": 0,
  "totalDiscount": 0,
  "note": "",
  "insertedAt": 1766060613022,
  "posUpdatedAt": 1766060615349,
  "paidAt": 0,
  "orderItems": null,  // ⚠️ Chưa extract
  "shippingAddress": null,  // ⚠️ Chưa extract
  "warehouseInfo": null,  // ⚠️ Chưa extract
  "customerInfo": null,  // ⚠️ Chưa extract
  "posData": {
    // ... rất nhiều fields từ POS API
    "id": 3037,
    "system_id": 3037,
    "shop_id": 860225178,
    "status": 1,
    "status_name": "submitted",
    "bill_full_name": "My Dung Truong",
    "bill_phone_number": "0944252001",
    "bill_email": null,
    "customer": { /* full customer object */ },
    "warehouse_info": { /* full warehouse object */ },
    "shipping_address": { /* full address object */ },
    "items": [ /* array of order items */ ],
    "conversation_id": "109383448131220_25860307226895435",
    "page_id": "109383448131220",
    "post_id": "109383448131220_122260255748023280",
    "ad_id": "120233654668590705",
    // ... nhiều fields khác
  }
}
```

### So Sánh với Pancake POS API Documentation

#### ✅ Fields Đã Extract

| Field trong DB | Field trong API | Status | Notes |
|---------------|----------------|--------|-------|
| `orderId` | `posData.id` | ✅ | Order ID |
| `systemId` | `posData.system_id` | ✅ | System ID |
| `shopId` | `posData.shop_id` | ✅ | Shop ID |
| `status` | `posData.status` | ✅ | Status code |
| `statusName` | `posData.status_name` | ✅ | Status name |
| `billFullName` | `posData.bill_full_name` | ✅ | Tên người thanh toán |
| `billPhoneNumber` | `posData.bill_phone_number` | ✅ | SĐT người thanh toán |
| `billEmail` | `posData.bill_email` | ✅ | Email (có thể null) |
| `customerId` | `posData.customer.id` | ✅ | Customer ID |
| `warehouseId` | `posData.warehouse_id` | ✅ | Warehouse ID |
| `pageId` | `posData.page_id` | ✅ | Facebook Page ID |
| `postId` | `posData.post_id` | ✅ | Facebook Post ID |
| `shippingFee` | `posData.shipping_fee` | ✅ | Phí vận chuyển |
| `totalDiscount` | `posData.total_discount` | ✅ | Tổng giảm giá |
| `note` | `posData.note` | ✅ | Ghi chú |
| `insertedAt` | `posData.inserted_at` | ✅ | Thời gian tạo |
| `posUpdatedAt` | `posData.updated_at` | ✅ | Thời gian cập nhật |
| `paidAt` | `posData.paid_at` | ✅ | Thời gian thanh toán |

#### ⚠️ Fields Có Trong `posData` Nhưng Chưa Extract

| Field trong API | Field trong DB | Status | Priority | Notes |
|----------------|----------------|--------|----------|-------|
| `items` | `orderItems` | ❌ | **High** | Danh sách sản phẩm trong đơn - **QUAN TRỌNG** |
| `shipping_address` | `shippingAddress` | ❌ | **High** | Địa chỉ giao hàng - **QUAN TRỌNG** |
| `warehouse_info` | `warehouseInfo` | ❌ | Medium | Thông tin kho hàng |
| `customer` | `customerInfo` | ❌ | Medium | Thông tin khách hàng đầy đủ |
| `conversation_id` | - | ❌ | **High** | Link với Facebook conversation |
| `ad_id` | - | ❌ | Medium | Facebook Ad ID |
| `total_price` | - | ❌ | **High** | Tổng giá trị đơn hàng |
| `money_to_collect` | - | ❌ | **High** | Số tiền cần thu |
| `cod` | - | ❌ | Medium | Tiền COD |
| `order_link` | - | ❌ | Low | Link đến order trên POS |
| `tracking_link` | - | ❌ | Medium | Link tracking đơn hàng |
| `tags` | - | ❌ | Low | Tags của đơn hàng |
| `assigning_seller` | - | ❌ | Medium | Người bán được assign |
| `creator` | - | ❌ | Low | Người tạo đơn |
| `status_history` | - | ❌ | Medium | Lịch sử thay đổi status |
| `payment_purchase_histories` | - | ❌ | Low | Lịch sử thanh toán |
| `activated_promotion_advances` | - | ❌ | Low | Khuyến mãi đã áp dụng |
| `activated_combo_products` | - | ❌ | Low | Combo products đã áp dụng |

#### 📊 Order Items Structure (trong `posData.items`)

```json
{
  "id": 11215990254,
  "product_id": "14b5b2db-719a-4e34-aeb4-d0e9daaaa14f",
  "variation_id": "uuid",
  "quantity": 1,
  "price": 700000,
  "total": 700000,
  "discount_each_product": 0,
  "note": null,
  "note_product": "",
  "is_bonus_product": false,
  "is_wholesale": false,
  "return_quantity": 0,
  "returned_count": 0,
  // ... nhiều fields khác
}
```

**Cần extract:**
- `product_id` → Link với `pc_pos_products`
- `variation_id` → Link với `pc_pos_variations`
- `quantity`, `price`, `total` → Tính toán doanh thu
- `discount_each_product` → Phân tích discount

#### 📊 Shipping Address Structure (trong `posData.shipping_address`)

```json
{
  "full_name": "My Dung Truong",
  "phone_number": "0944252001",
  "full_address": "C ở 18, Phường Láng, Hà Nội",
  "address": "C ở 18",
  "province_id": "84_VN101",
  "province_name": "Hà Nội",
  "district_id": null,
  "district_name": null,
  "commune_id": "84_VN10111",
  "commune_name": "Phường Láng",
  "country_code": "84",
  "post_code": null
}
```

**Cần extract để:**
- Phân tích địa lý (tỉnh/thành nào bán nhiều nhất)
- Tối ưu logistics
- Phân tích customer location

---

## Facebook Conversations - Cấu Trúc Thực Tế

### Cấu Trúc Document Thực Tế

```json
{
  "_id": "ObjectId",
  "conversationId": "102039018873979_9570176223069085",
  "pageId": "102039018873979",
  "pageUsername": "Folkformint6",
  "customerId": "bb6dac25-2c05-412a-8d66-6b916b33c1c7",
  "panCakeData": {
    "id": "102039018873979_9570176223069085",
    "page_id": "102039018873979",
    "customer_id": "bb6dac25-2c05-412a-8d66-6b916b33c1c7",
    "from": {
      "email": "9570176223069085@facebook.com",
      "id": "9570176223069085",
      "name": "Vicky Hà My"
    },
    "customers": [ /* array of customer objects */ ],
    "page_customer": { /* customer info */ },
    "inserted_at": "2025-04-14T02:40:07.000000",
    "message_count": 66,
    "snippet": "Nếu chị thích mẫu khăn hay cần tư vấn cách cột phù hợp, em có...",
    "seen": true,
    "has_phone": false,
    "recent_phone_numbers": [],
    "post_id": null,
    "ad_ids": ["120219624292870241"],
    "ads": [ /* array of ad objects */ ],
    "tag_histories": [ /* array of tag history */ ],
    "assignee_histories": [],
    "assignee_ids": [],
    "current_assign_users": [],
    "last_sent_by": { /* user object */ }
  },
  "panCakeUpdatedAt": 1765183269000,
  "createdAt": 1765994113126,
  "updatedAt": 1765994113126
}
```

### So Sánh với Pancake API Documentation

#### ✅ Fields Đã Sync

| Field trong DB | Field trong API | Status | Notes |
|---------------|----------------|--------|-------|
| `conversationId` | `panCakeData.id` | ✅ | Conversation ID |
| `pageId` | `panCakeData.page_id` | ✅ | Page ID |
| `customerId` | `panCakeData.customer_id` | ✅ | Customer ID |
| `panCakeData.from` | `from` | ✅ | Người gửi |
| `panCakeData.customers` | `customers` | ✅ | Danh sách customers |
| `panCakeData.page_customer` | `page_customer` | ✅ | Customer info |
| `panCakeData.inserted_at` | `inserted_at` | ✅ | Thời gian tạo |
| `panCakeData.message_count` | `message_count` | ✅ | Số lượng messages |
| `panCakeData.snippet` | `snippet` | ✅ | Snippet tin nhắn cuối |
| `panCakeData.seen` | `seen` | ✅ | Đã xem |
| `panCakeData.has_phone` | `has_phone` | ✅ | Có số điện thoại |
| `panCakeData.recent_phone_numbers` | `recent_phone_numbers` | ✅ | SĐT gần đây |
| `panCakeData.post_id` | `post_id` | ✅ | Post ID (nếu từ post) |
| `panCakeData.ad_ids` | `ad_ids` | ✅ | Ad IDs |
| `panCakeData.ads` | `ads` | ✅ | Ad objects |
| `panCakeData.tag_histories` | `tag_histories` | ✅ | Lịch sử tags |
| `panCakeData.assignee_histories` | `assignee_histories` | ✅ | Lịch sử assign |
| `panCakeData.assignee_ids` | `assignee_ids` | ✅ | Assignee IDs |
| `panCakeData.current_assign_users` | `current_assign_users` | ✅ | Users hiện tại được assign |
| `panCakeData.last_sent_by` | `last_sent_by` | ✅ | Người gửi cuối |

#### ⚠️ Fields Có Trong API Nhưng Cần Kiểm Tra

Theo [Pancake API Documentation](./pancake-api-context.md), Conversation schema có:

| Field trong API | Field trong DB | Status | Notes |
|----------------|----------------|--------|-------|
| `type` | `panCakeData.type` | ⚠️ | INBOX, COMMENT, LIVESTREAM - cần kiểm tra |
| `updated_at` | `panCakeData.updated_at` | ⚠️ | Cần kiểm tra |
| `tags` | `panCakeData.tags` | ⚠️ | Tags hiện tại - cần kiểm tra |
| `last_message` | `panCakeData.last_message` | ⚠️ | Tin nhắn cuối - cần kiểm tra |
| `participants` | `panCakeData.participants` | ⚠️ | Participants - cần kiểm tra |

#### ❌ Fields Thiếu Quan Trọng

1. **`type`**: Cần extract để phân biệt INBOX, COMMENT, LIVESTREAM
2. **`updated_at`**: Cần extract để track thời gian cập nhật
3. **`tags`**: Tags hiện tại (khác với `tag_histories`)
4. **`last_message`**: Thông tin tin nhắn cuối (có thể dùng thay cho `snippet`)

---

## Facebook Messages - Cấu Trúc Thực Tế

### Cấu Trúc Document Thực Tế (fb_message_items)

```json
{
  "_id": "ObjectId",
  "messageId": "m_dxBgYEOzWq2GwmxPrqeAHQMBLexzlInjdlLHZU4paUsr7Gla4cODvr0q8T9xMxXVARBI9gXkvBBBw9pdD_Eauw",
  "conversationId": "102039018873979_9570176223069085",
  "insertedAt": 1746853493,
  "messageData": {
    "id": "m_dxBgYEOzWq2GwmxPrqeAHQMBLexzlInjdlLHZU4paUsr7Gla4cODvr0q8T9xMxXVARBI9gXkvBBBw9pdD_Eauw",
    "conversation_id": "102039018873979_9570176223069085",
    "page_id": "102039018873979",
    "from": {
      "email": "102039018873979@facebook.com",
      "id": "102039018873979",
      "name": "Folk Form"
    },
    "message": "<div></div>",
    "type": "INBOX",
    "inserted_at": "2025-05-10T05:04:53.000000",
    "is_hidden": false,
    "is_removed": false,
    "has_phone": false,
    "seen": true,
    "attachments": [ /* array of attachment objects */ ],
    "can_comment": false,
    "can_hide": false,
    "can_like": false,
    "can_remove": false,
    "can_reply_privately": false,
    "comment_count": null,
    "edit_history": null,
    "like_count": null,
    "message_tags": [],
    "original_message": "",
    "parent_id": null,
    "phone_info": [],
    "private_reply_conversation": null,
    "removed_by": null,
    "rich_message": null,
    "show_info": false,
    "user_likes": false,
    "is_livestream_order": null,
    "is_parent": false,
    "is_parent_hidden": false
  },
  "createdAt": 1765994113544,
  "updatedAt": 1765994113544
}
```

### So Sánh với Pancake API Documentation

#### ✅ Fields Đã Sync

| Field trong DB | Field trong API | Status | Notes |
|---------------|----------------|--------|-------|
| `messageId` | `messageData.id` | ✅ | Message ID |
| `conversationId` | `messageData.conversation_id` | ✅ | Conversation ID |
| `insertedAt` | `messageData.inserted_at` | ✅ | Thời gian insert |
| `messageData.from` | `from` | ✅ | Người gửi |
| `messageData.message` | `message` | ✅ | Nội dung tin nhắn |
| `messageData.type` | `type` | ✅ | Loại tin nhắn |
| `messageData.page_id` | `page_id` | ✅ | Page ID |
| `messageData.is_hidden` | `is_hidden` | ✅ | Đã ẩn |
| `messageData.is_removed` | `is_removed` | ✅ | Đã xóa |
| `messageData.has_phone` | `has_phone` | ✅ | Có số điện thoại |
| `messageData.seen` | `seen` | ✅ | Đã xem |
| `messageData.attachments` | `attachments` | ✅ | File đính kèm |
| `messageData.can_*` | `can_*` | ✅ | Permissions |
| `messageData.comment_count` | `comment_count` | ✅ | Số comment |
| `messageData.like_count` | `like_count` | ✅ | Số like |
| `messageData.message_tags` | `message_tags` | ✅ | Tags trong message |
| `messageData.parent_id` | `parent_id` | ✅ | Parent message ID |
| `messageData.phone_info` | `phone_info` | ✅ | Thông tin phone |

#### 📊 Attachments Structure

```json
{
  "id": "1661458518578441",
  "type": "photo",
  "mime_type": "image/jpeg",
  "name": "image-1661458518578441",
  "size": 275856,
  "url": "https://content.pancake.vn/...",
  "can_download": true,
  "image_data": {
    "height": 2048,
    "width": 1365,
    "url": "https://scontent.fdad5-1.fna.fbcdn.net/...",
    "preview_url": "https://scontent.fdad5-1.fna.fbcdn.net/..."
  }
}
```

**Cần extract để:**
- Phân tích loại content (ảnh, video, file)
- Download và lưu trữ media
- Phân tích visual content với AI

---

## So Sánh với API Documentation

### Pancake API - Coverage Analysis

| Endpoint | Collection | Fields Sync | Status | Notes |
|----------|-----------|-------------|--------|-------|
| `/pages/{page_id}/conversations` | `fb_conversations` | ~90% | ✅ Good | Thiếu một số fields như `type`, `updated_at` |
| `/pages/{page_id}/conversations/{id}/messages` | `fb_message_items` | ~95% | ✅ Excellent | Gần như đầy đủ |
| `/pages/{page_id}/page_customers` | `customers.panCakeData` | ~85% | ⚠️ Fair | Thiếu `email`, `tags` |
| `/pages/{page_id}/posts` | `fb_posts` | ? | ⚠️ | Cần kiểm tra |

### Pancake POS API - Coverage Analysis

| Endpoint | Collection | Fields Sync | Status | Notes |
|----------|-----------|-------------|--------|-------|
| `/shops/{shop_id}/orders` | `pc_pos_orders` | ~60% | ❌ **Poor** | **Thiếu quan trọng**: `items`, `shipping_address`, `total_price` |
| `/shops/{shop_id}/customers` | `customers.posData` | 0% | ❌ **Missing** | **Chưa sync customers từ POS** |
| `/shops/{shop_id}/products` | `pc_pos_products` | ? | ⚠️ | Cần kiểm tra |
| `/shops/{shop_id}/products/variations` | `pc_pos_variations` | ? | ⚠️ | Cần kiểm tra |
| `/shops/{shop_id}/warehouses` | `pc_pos_warehouses` | ? | ⚠️ | Cần kiểm tra |

---

## Gaps & Recommendations

### 🔴 Critical Gaps (Cần Fix Ngay)

1. **Orders - Missing Order Items**
   - **Impact**: Không thể phân tích sản phẩm bán chạy, doanh thu theo sản phẩm
   - **Fix**: Extract `posData.items` → `orderItems`
   - **Priority**: **HIGH**

2. **Orders - Missing Shipping Address**
   - **Impact**: Không thể phân tích địa lý, logistics
   - **Fix**: Extract `posData.shipping_address` → `shippingAddress`
   - **Priority**: **HIGH**

3. **Orders - Missing Total Price**
   - **Impact**: Không thể tính doanh thu, AOV
   - **Fix**: Extract `posData.total_price`, `posData.money_to_collect`
   - **Priority**: **HIGH**

4. **Customers - No POS Data**
   - **Impact**: Không có unified customer view, thiếu data từ POS
   - **Fix**: Sync customers từ POS API và merge với Pancake customers
   - **Priority**: **HIGH**

5. **Customers - Sources Field Null**
   - **Impact**: Không biết customer đến từ nguồn nào
   - **Fix**: Populate `sources` field khi sync
   - **Priority**: **MEDIUM**

### 🟡 Important Gaps (Nên Fix Sớm)

6. **Orders - Missing Conversation Link**
   - **Impact**: Không thể track journey từ conversation → order
   - **Fix**: Extract `posData.conversation_id`
   - **Priority**: **MEDIUM**

7. **Conversations - Missing Type**
   - **Impact**: Không phân biệt được INBOX, COMMENT, LIVESTREAM
   - **Fix**: Extract `panCakeData.type`
   - **Priority**: **MEDIUM**

8. **Orders - Missing Warehouse Info**
   - **Impact**: Không thể phân tích theo kho hàng
   - **Fix**: Extract `posData.warehouse_info` → `warehouseInfo`
   - **Priority**: **MEDIUM**

9. **Orders - Missing Customer Info**
   - **Impact**: Thiếu thông tin customer đầy đủ trong order
   - **Fix**: Extract `posData.customer` → `customerInfo`
   - **Priority**: **MEDIUM**

### 🟢 Nice to Have

10. **Orders - Missing Tags, Status History, etc.**
    - **Impact**: Thiếu metadata cho phân tích nâng cao
    - **Fix**: Extract các fields metadata
    - **Priority**: **LOW**

11. **Messages - Extract Attachment URLs**
    - **Impact**: Không thể download và phân tích media
    - **Fix**: Extract và lưu attachment URLs
    - **Priority**: **LOW**

---

## Action Plan

### Phase 1: Critical Fixes (1-2 tuần)

1. ✅ Extract `orderItems` từ `posData.items`
2. ✅ Extract `shippingAddress` từ `posData.shipping_address`
3. ✅ Extract `total_price`, `money_to_collect` từ `posData`
4. ✅ Fix `phoneNumbers` format (bỏ dấu `[]`)
5. ✅ Populate `sources` field

### Phase 2: Important Fixes (2-4 tuần)

6. ✅ Sync customers từ POS API
7. ✅ Implement customer matching logic
8. ✅ Extract `conversation_id` từ orders
9. ✅ Extract `type` từ conversations
10. ✅ Extract `warehouseInfo`, `customerInfo` từ orders

### Phase 3: Enhancements (1-2 tháng)

11. ✅ Extract metadata fields (tags, status_history, etc.)
12. ✅ Extract attachment URLs và implement media storage
13. ✅ Implement data validation và cleaning
14. ✅ Build analytics-ready calculated fields

---

## Kết Luận

Hệ thống đã sync được **phần lớn dữ liệu** từ các API nguồn, nhưng còn **thiếu một số fields quan trọng** đặc biệt là:

1. **Order Items** - Critical cho phân tích sản phẩm
2. **Shipping Address** - Critical cho phân tích địa lý
3. **POS Customers** - Critical cho unified customer view
4. **Total Price** - Critical cho phân tích doanh thu

Với việc fix các gaps này, hệ thống sẽ có đủ dữ liệu để:
- ✅ Phân tích doanh thu và sản phẩm
- ✅ Phân tích customer journey
- ✅ Phân tích địa lý và logistics
- ✅ Ứng dụng AI cho insights kinh doanh

