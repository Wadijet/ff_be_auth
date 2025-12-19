# Phân Tích Phương Án Kiến Trúc Collection Message

## 📋 Vấn Đề

Khi merge messages vào 1 document, document sẽ ngày càng lớn:
- **100 messages** → ~500KB - 1MB
- **1000 messages** → ~5MB - 10MB  
- **10000 messages** → **Vượt quá 16MB (giới hạn MongoDB)!**

Cần quyết định: **Giữ collection cũ hay tạo collection mới?**

---

## 🎯 Các Phương Án

### Phương Án 1: Giữ Collection Cũ (Metadata) + Collection Mới (Messages) ⭐ KHUYẾN NGHỊ

**Cấu trúc:**
- **Collection `fb_messages` (cũ)**: Lưu metadata + panCakeData (không có messages)
- **Collection `fb_message_items` (mới)**: Lưu từng message riêng lẻ

**Model:**

```go
// Collection fb_messages (metadata)
type FbMessage struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    PageId         string                 `json:"pageId" bson:"pageId" index:"text"`
    PageUsername   string                 `json:"pageUsername" bson:"pageUsername" index:"text"`
    ConversationId string                 `json:"conversationId" bson:"conversationId" index:"unique;text"`
    CustomerId     string                 `json:"customerId" bson:"customerId" index:"text"`
    
    // PanCakeData KHÔNG có messages[] (chỉ các field khác)
    PanCakeData    map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    
    // Metadata tracking
    LastSyncedAt   int64                  `json:"lastSyncedAt" bson:"lastSyncedAt"`
    TotalMessages  int64                  `json:"totalMessages" bson:"totalMessages"` // Tổng số messages trong collection items
    HasMore        bool                   `json:"hasMore" bson:"hasMore"`
    
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
    UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`
}

// Collection fb_message_items (từng message)
type FbMessageItem struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    ConversationId string                 `json:"conversationId" bson:"conversationId" index:"text"` // Không unique
    MessageId      string                 `json:"messageId" bson:"messageId" index:"unique;text"` // ID từ Pancake (unique)
    MessageData    map[string]interface{} `json:"messageData" bson:"messageData"` // Dữ liệu message
    InsertedAt     int64                  `json:"insertedAt" bson:"insertedAt"` // Extract từ messageData.inserted_at
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
    UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

**Ưu điểm:**
- ✅ **Tương thích ngược**: Collection cũ vẫn hoạt động, chỉ bỏ messages[]
- ✅ **Migration dễ**: Có thể migrate từng bước
- ✅ **Rõ ràng**: Tách biệt metadata và messages
- ✅ **Scalable**: Không có giới hạn số lượng messages
- ✅ **Query linh hoạt**: Query metadata nhanh, query messages riêng

**Nhược điểm:**
- ⚠️ Cần query 2 collections khi cần cả metadata + messages
- ⚠️ Cần maintain 2 collections

**Index:**
- `fb_messages`: `conversationId` (unique)
- `fb_message_items`: `conversationId` + `insertedAt` (compound), `messageId` (unique)

---

### Phương Án 2: Thay Đổi Hoàn Toàn Collection Message

**Cấu trúc:**
- **Collection `fb_messages`**: Mỗi message là 1 document
- **Bỏ collection metadata** (hoặc di chuyển sang `fb_conversations`)

**Model:**

```go
// Mỗi message là 1 document
type FbMessage struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    PageId         string                 `json:"pageId" bson:"pageId" index:"text"`
    PageUsername   string                 `json:"pageUsername" bson:"pageUsername" index:"text"`
    ConversationId string                 `json:"conversationId" bson:"conversationId" index:"text"` // Không unique
    CustomerId     string                 `json:"customerId" bson:"customerId" index:"text"`
    MessageId      string                 `json:"messageId" bson:"messageId" index:"unique;text"` // Unique
    MessageData    map[string]interface{} `json:"messageData" bson:"messageData"`
    InsertedAt     int64                  `json:"insertedAt" bson:"insertedAt"`
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
    UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

**Ưu điểm:**
- ✅ **Đơn giản**: Chỉ 1 collection
- ✅ **Scalable**: Không có giới hạn
- ✅ **Query nhanh**: Có index trên conversationId

**Nhược điểm:**
- ❌ **Breaking change**: Phá vỡ cấu trúc hiện tại
- ❌ **Mất metadata**: Cần lưu metadata ở đâu? (có thể ở `fb_conversations`)
- ❌ **Migration phức tạp**: Cần migrate toàn bộ dữ liệu cũ

---

### Phương Án 3: Hybrid - Metadata + Messages Gần Đây

**Cấu trúc:**
- **Collection `fb_messages`**: Metadata + 100-200 messages gần đây nhất
- **Collection `fb_message_archive`**: Messages cũ (từng message riêng)

**Logic:**
- Khi merge, nếu messages > 200 → Di chuyển messages cũ vào archive
- Query: Lấy từ `fb_messages` + `fb_message_archive`

**Ưu điểm:**
- ✅ Messages gần đây query nhanh
- ✅ Có thể archive messages cũ

**Nhược điểm:**
- ⚠️ Logic phức tạp (cần di chuyển messages)
- ⚠️ Cần merge từ 2 collections khi query
- ⚠️ Vẫn có giới hạn 200 messages trong document

---

## 📊 So Sánh Chi Tiết

| Tiêu Chí | Phương Án 1 (Metadata + Items) | Phương Án 2 (Chỉ Items) | Phương Án 3 (Hybrid) |
|----------|-------------------------------|------------------------|---------------------|
| **Tương thích ngược** | ✅ Tốt | ❌ Breaking change | ⚠️ Cần migration |
| **Scalability** | ⭐⭐⭐ Tốt nhất | ⭐⭐⭐ Tốt nhất | ⭐⭐ Trung bình |
| **Độ phức tạp** | ⭐⭐ Trung bình | ⭐⭐⭐ Đơn giản | ⚠️ Phức tạp |
| **Query metadata** | ✅ Nhanh (1 collection) | ⚠️ Cần query messages | ✅ Nhanh |
| **Query messages** | ✅ Nhanh (có index) | ✅ Nhanh (có index) | ⚠️ Cần merge 2 collections |
| **Migration** | ✅ Dễ (từng bước) | ❌ Phức tạp (toàn bộ) | ⚠️ Trung bình |
| **Maintain** | ⚠️ 2 collections | ✅ 1 collection | ⚠️ 2 collections |
| **Storage** | ⚠️ Trùng lặp metadata | ✅ Tối ưu | ⚠️ Trùng lặp |
| **Khuyến nghị** | ✅ **Nên dùng** | ⚠️ Có thể dùng | ❌ Không nên |

---

## 💡 Đề Xuất: Phương Án 1 - Metadata + Items

### Lý Do

1. **Tương thích ngược**: Collection `fb_messages` cũ vẫn hoạt động, chỉ cần:
   - Bỏ `messages[]` khỏi `panCakeData`
   - Thêm tracking fields (`totalMessages`, `hasMore`)

2. **Migration dễ dàng**:
   - Bước 1: Tạo collection `fb_message_items` mới
   - Bước 2: Extract messages từ `fb_messages.panCakeData.messages[]` → `fb_message_items`
   - Bước 3: Xóa `messages[]` khỏi `fb_messages.panCakeData`
   - Có thể làm từng bước, không cần downtime

3. **Rõ ràng về mục đích**:
   - `fb_messages`: Metadata conversation (1 document/conversation)
   - `fb_message_items`: Messages (nhiều documents/conversation)

4. **Query linh hoạt**:
   - Query metadata: `fb_messages` (nhanh, document nhỏ)
   - Query messages: `fb_message_items` với filter `conversationId` (nhanh, có index)
   - Query cả 2: Join/aggregate khi cần

5. **Scalable**: Không có giới hạn số lượng messages

---

## 🏗️ Cấu Trúc Chi Tiết

### Collection `fb_messages` (Metadata)

```go
type FbMessage struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    PageId         string                 `json:"pageId" bson:"pageId" index:"text"`
    PageUsername   string                 `json:"pageUsername" bson:"pageUsername" index:"text"`
    ConversationId string                 `json:"conversationId" bson:"conversationId" index:"unique;text"`
    CustomerId     string                 `json:"customerId" bson:"customerId" index:"text"`
    
    // PanCakeData KHÔNG có messages[] (chỉ các field khác)
    PanCakeData    map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    // PanCakeData chứa:
    // - conv_from, read_watermarks, activities, ad_clicks
    // - is_banned, banned_count, banned_by, notes
    // - reports_by_phone, reported_count
    // - customers, conv_customers, ...
    // KHÔNG có: messages[]
    
    // Metadata tracking
    LastSyncedAt   int64                  `json:"lastSyncedAt" bson:"lastSyncedAt"`
    TotalMessages  int64                  `json:"totalMessages" bson:"totalMessages"` // Tổng số messages trong fb_message_items
    HasMore        bool                   `json:"hasMore" bson:"hasMore"`
    
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
    UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

**Kích thước document**: ~10-50KB (không có messages)

### Collection `fb_message_items` (Messages)

```go
type FbMessageItem struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    ConversationId string                 `json:"conversationId" bson:"conversationId" index:"text"` // Không unique
    MessageId      string                 `json:"messageId" bson:"messageId" index:"unique;text"` // ID từ Pancake (unique)
    MessageData    map[string]interface{} `json:"messageData" bson:"messageData"` // Toàn bộ dữ liệu message
    InsertedAt     int64                  `json:"insertedAt" bson:"insertedAt" index:"text"` // Extract từ messageData.inserted_at
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
    UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

**Kích thước document**: ~5-10KB/message

**Index:**
- `conversationId` + `insertedAt` (compound index) để query nhanh
- `messageId` (unique) để tránh duplicate

---

## 🔄 Flow Đồng Bộ

### Endpoint: Upsert Messages (API Bên Ngoài - Giữ Nguyên)

**POST** `/api/v1/facebook/message/upsert-messages`

**Request (Giữ nguyên như cũ):**
```json
{
  "conversationId": "157725629736743_9350439438393456",
  "pageId": "157725629736743",
  "pageUsername": "Folkformint",
  "customerId": "8b168fa9-4836-4648-a3fd-799c227675a1",
  "panCakeData": {
    "conv_from": {...},
    "read_watermarks": [...],
    "activities": [...],
    "messages": [
      {
        "id": "m_xxx1",
        "conversation_id": "157725629736743_9350439438393456",
        "message": "<div>Message 1</div>",
        "inserted_at": "2025-12-16T15:22:45.000000",
        ...
      },
      // ... 30 messages
    ],
    // ... các field khác
  },
  "hasMore": true
}
```

**⚠️ Lưu ý:** API bên ngoài vẫn gửi `panCakeData` đầy đủ (bao gồm `messages[]`), server sẽ tự động tách.

**Logic Xử Lý Nội Bộ (Trong Service Layer):**

1. **Tách messages[] khỏi panCakeData**:
   ```go
   messages := panCakeData["messages"].([]interface{})
   metadataPanCakeData := panCakeData // Copy
   delete(metadataPanCakeData, "messages") // Xóa messages[] khỏi metadata
   ```

2. **Upsert metadata** vào `fb_messages`:
   - Upsert theo `conversationId`
   - Ghi đè `panCakeData` (đã bỏ messages[])
   - Cập nhật `lastSyncedAt`, `hasMore`

3. **Upsert messages** vào `fb_message_items`:
   - Bulk upsert từng message theo `messageId`
   - Tự động tránh duplicate
   - Extract `insertedAt` từ `messageData.inserted_at`

4. **Cập nhật totalMessages**:
   - Count messages trong `fb_message_items` theo `conversationId`
   - Update vào `fb_messages.totalMessages`

---

## 📝 Query Messages

### Endpoint: Get Messages by Conversation

**GET** `/api/v1/facebook/message/find-by-conversation/:conversationId?page=1&limit=50&sort=insertedAt`

**Logic:**
- Query `fb_message_items` với filter `conversationId`
- Sort theo `insertedAt`
- Phân trang tự nhiên

**Response:**
```json
{
  "data": {
    "metadata": {
      "conversationId": "...",
      "totalMessages": 1000,
      "hasMore": false,
      "panCakeData": {...}
    },
    "messages": [
      {...},
      {...}
    ],
    "pagination": {
      "page": 1,
      "limit": 50,
      "total": 1000
    }
  }
}
```

---

## ✅ Ưu Điểm Phương Án 1

1. **Tương thích ngược**: Collection cũ vẫn hoạt động
2. **Migration dễ**: Có thể làm từng bước
3. **Scalable**: Không có giới hạn messages
4. **Performance tốt**: 
   - Query metadata: Document nhỏ, nhanh
   - Query messages: Có index, nhanh
5. **Rõ ràng**: Tách biệt metadata và messages
6. **Linh hoạt**: Có thể query riêng hoặc join

---

## ⚠️ Nhược Điểm & Giải Pháp

### 1. Cần Query 2 Collections

**Giải pháp:**
- Query metadata: Chỉ cần `fb_messages` (nhanh)
- Query messages: Chỉ cần `fb_message_items` (nhanh, có index)
- Query cả 2: Có thể cache metadata hoặc dùng aggregation

### 2. Maintain 2 Collections

**Giải pháp:**
- Logic rõ ràng: Metadata ở 1 nơi, messages ở 1 nơi
- Service methods tách biệt: `UpsertMetadata()`, `UpsertMessages()`
- Dễ test và maintain

### 3. Trùng Lặp Metadata

**Giải pháp:**
- Metadata nhỏ (~10-50KB), không ảnh hưởng nhiều
- Có thể cache metadata nếu cần

---

## 🎯 Kết Luận

### ✅ **KHUYẾN NGHỊ: Phương Án 1 - Metadata + Items**

**Lý do:**
1. Tương thích ngược tốt nhất
2. Migration dễ dàng
3. Scalable và performance tốt
4. Rõ ràng về mục đích
5. Linh hoạt trong query

**Cấu trúc:**
- `fb_messages`: Metadata (1 document/conversation)
- `fb_message_items`: Messages (nhiều documents/conversation)

**Trade-off:**
- Cần maintain 2 collections, nhưng logic rõ ràng
- Cần query 2 collections khi cần cả 2, nhưng có index nên nhanh

---

## 📋 Checklist Implementation

- [ ] Tạo model `FbMessageItem` cho collection mới
- [ ] Cập nhật model `FbMessage` (bỏ messages[], thêm tracking fields)
- [ ] Tạo collection `fb_message_items` trong init
- [ ] Tạo service `FbMessageItemService`
- [ ] Tạo method `UpsertMessages()` trong `FbMessageService` (upsert vào items)
- [ ] Tạo method `UpsertMetadata()` trong `FbMessageService` (upsert metadata)
- [ ] Tạo endpoint `POST /api/v1/facebook/message/upsert-messages`
- [ ] Tạo endpoint `GET /api/v1/facebook/message/find-by-conversation/:id`
- [ ] Tạo migration script để extract messages từ cũ sang mới
- [ ] Update index cho cả 2 collections
