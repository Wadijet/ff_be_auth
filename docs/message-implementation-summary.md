# Tóm Tắt Implementation: Endpoint Upsert Messages

## ✅ Đã Hoàn Thành

### 1. Model Mới

**File:** `api/core/api/models/mongodb/model.fb.message.item.go`
- Tạo model `FbMessageItem` cho collection `fb_message_items`
- Mỗi message là 1 document riêng
- Index: `messageId` (unique), `conversationId` + `insertedAt` (compound)

**File:** `api/core/api/models/mongodb/model.fb.message.go`
- Cập nhật model `FbMessage` (thêm `LastSyncedAt`, `TotalMessages`, `HasMore`)
- `PanCakeData` không còn lưu `messages[]` (chỉ metadata)

### 2. Service Layer

**File:** `api/core/api/services/service.fb.message.item.go`
- Tạo `FbMessageItemService` với các method:
  - `UpsertMessages()`: Bulk upsert messages (tự động tránh duplicate)
  - `FindByConversationId()`: Query messages với phân trang
  - `CountByConversationId()`: Đếm số lượng messages

**File:** `api/core/api/services/service.fb.message.go`
- Thêm `fbMessageItemService` vào `FbMessageService`
- Tạo method `UpsertMessages()`:
  - **Logic nội bộ**: Tự động tách `messages[]` ra khỏi `panCakeData`
  - Upsert metadata vào `fb_messages` (không có messages[])
  - Upsert messages vào `fb_message_items` (từng message riêng)
  - Cập nhật `totalMessages`

### 3. Handler & Route

**File:** `api/core/api/dto/dto.fb.message.go`
- Tạo `FbMessageUpsertMessagesInput` DTO

**File:** `api/core/api/handler/handler.fb.message.go`
- Tạo `HandleUpsertMessages()` handler method

**File:** `api/core/api/router/routes.go`
- Đăng ký route: `POST /api/v1/facebook/message/upsert-messages`
- Permission: `FbMessage.Update`

### 4. Database Setup

**File:** `api/core/global/global.vars.go`
- Thêm `FbMessageItems` collection name

**File:** `api/cmd/server/init.go`
- Khởi tạo collection name
- Tạo index cho `fb_message_items`

**File:** `api/cmd/server/init.registry.go`
- Đăng ký collection `fb_message_items` vào registry

---

## 🔄 Endpoint Đã Tạo

### POST `/api/v1/facebook/message/upsert-messages`

**Request (Giữ nguyên như cũ - API bên ngoài không cần thay đổi):**
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

**Logic Xử Lý Nội Bộ (Trong Service):**
1. Tách `messages[]` ra khỏi `panCakeData`
2. Upsert metadata (không có messages[]) vào `fb_messages`
3. Upsert từng message vào `fb_message_items` (bulk upsert, tự động tránh duplicate)
4. Cập nhật `totalMessages` trong `fb_messages`

**Response:**
```json
{
  "data": {
    "id": "...",
    "conversationId": "157725629736743_9350439438393456",
    "panCakeData": {
      // Không có messages[]
      "conv_from": {...},
      "read_watermarks": [...],
      ...
    },
    "totalMessages": 30,
    "hasMore": true,
    "lastSyncedAt": 1765898960082
  }
}
```

---

## 📊 Cấu Trúc Dữ Liệu

### Collection `fb_messages` (Metadata)
- 1 document/conversation
- Kích thước: ~10-50KB (không có messages[])
- Chứa: metadata, panCakeData (không có messages[])

### Collection `fb_message_items` (Messages)
- Nhiều documents/conversation (mỗi message là 1 document)
- Kích thước: ~5-10KB/message
- Chứa: từng message riêng lẻ

---

## ✅ Lợi Ích

1. **API bên ngoài không cần thay đổi**: Vẫn gửi `panCakeData` đầy đủ
2. **Logic tách tự động**: Server tự động tách và lưu vào 2 collections
3. **Scalable**: Không có giới hạn số lượng messages
4. **Performance tốt**: Query nhanh với index
5. **Tương thích ngược**: Collection cũ vẫn hoạt động

---

## 📝 Cần Làm Thêm (Optional)

- [ ] Tạo endpoint `GET /api/v1/facebook/message/find-by-conversation/:id` để query messages
- [ ] Tạo migration script để extract messages từ dữ liệu cũ
- [ ] Viết unit tests
