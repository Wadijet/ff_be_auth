# Tách Biệt Endpoint: CRUD vs Endpoint Đặc Biệt

## 📋 Tổng Quan

Hệ thống có **2 loại endpoint** cho Facebook Messages:

1. **CRUD Routes** (Logic chung) - Không tách messages
2. **Endpoint Đặc Biệt** (Logic tách messages) - Tự động tách messages vào collection riêng

---

## 🔄 CRUD Routes (Logic Chung)

### Endpoints
- `POST /api/v1/facebook/message/insert-one` - Tạo mới
- `PUT /api/v1/facebook/message/update-one` - Cập nhật
- `GET /api/v1/facebook/message/find` - Tìm kiếm
- `GET /api/v1/facebook/message/find-by-id/:id` - Tìm theo ID
- `DELETE /api/v1/facebook/message/delete-one` - Xóa
- ... (các endpoint CRUD khác)

### DTO
```go
type FbMessageCreateInput struct {
    PageId         string
    PageUsername   string
    ConversationId string
    CustomerId     string
    PanCakeData    map[string]interface{} // Có thể chứa messages[]
}
```

### Đặc Điểm
- ✅ **Không có logic tách messages**
- ✅ **Lưu nguyên panCakeData** (có thể có messages[])
- ✅ **Tương thích ngược** với dữ liệu cũ
- ✅ **Logic chung** từ BaseHandler, không thay đổi

### Khi Nào Dùng
- Tạo/cập nhật message thủ công
- Import dữ liệu từ nguồn khác
- Các thao tác CRUD chuẩn

---

## ⚡ Endpoint Đặc Biệt (Logic Tách Messages)

### Endpoint
- `POST /api/v1/facebook/message/upsert-messages` - Upsert và tự động tách messages

### DTO
```go
type FbMessageUpsertMessagesInput struct {
    PageId         string
    PageUsername   string
    ConversationId string
    CustomerId     string
    PanCakeData    map[string]interface{} // Đầy đủ (bao gồm messages[])
    HasMore        bool                   // Còn messages để sync không
}
```

### Đặc Điểm
- ✅ **Tự động tách messages[]** ra khỏi panCakeData
- ✅ **Lưu vào 2 collections**:
  - `fb_messages`: Metadata (không có messages[])
  - `fb_message_items`: Từng message riêng lẻ
- ✅ **Bulk upsert** messages (tự động tránh duplicate)
- ✅ **Cập nhật totalMessages** tự động

### Logic Xử Lý
1. Tách `messages[]` ra khỏi `panCakeData`
2. Upsert metadata vào `fb_messages` (không có messages[])
3. Upsert từng message vào `fb_message_items` (bulk upsert)
4. Cập nhật `totalMessages` trong `fb_messages`

### Khi Nào Dùng
- Sync messages từ Pancake API
- Đồng bộ dữ liệu tự động
- Xử lý số lượng messages lớn

---

## 📊 So Sánh

| Tiêu Chí | CRUD Routes | Endpoint Đặc Biệt |
|----------|-------------|-------------------|
| **Endpoint** | `/insert-one`, `/update-one`, ... | `/upsert-messages` |
| **DTO** | `FbMessageCreateInput` | `FbMessageUpsertMessagesInput` |
| **Tách messages** | ❌ Không | ✅ Có (tự động) |
| **Lưu messages[]** | ✅ Trong panCakeData | ❌ Tách ra collection riêng |
| **Collections** | Chỉ `fb_messages` | `fb_messages` + `fb_message_items` |
| **HasMore field** | ❌ Không có | ✅ Có |
| **Bulk upsert** | ❌ Không | ✅ Có |
| **Tự động count** | ❌ Không | ✅ Có (totalMessages) |
| **Khi nào dùng** | CRUD thủ công | Sync từ Pancake API |

---

## 🎯 Ví Dụ Sử Dụng

### CRUD Route (Không tách messages)

```bash
POST /api/v1/facebook/message/insert-one
{
  "pageId": "123",
  "conversationId": "456",
  "panCakeData": {
    "conv_from": {...},
    "messages": [
      {"id": "m1", "message": "Hello"},
      {"id": "m2", "message": "World"}
    ]
  }
}
```

**Kết quả:**
- Lưu vào `fb_messages` với `panCakeData.messages[]` còn nguyên
- Không tách messages ra collection riêng

### Endpoint Đặc Biệt (Tự động tách messages)

```bash
POST /api/v1/facebook/message/upsert-messages
{
  "pageId": "123",
  "conversationId": "456",
  "panCakeData": {
    "conv_from": {...},
    "messages": [
      {"id": "m1", "message": "Hello"},
      {"id": "m2", "message": "World"}
    ]
  },
  "hasMore": false
}
```

**Kết quả:**
- Lưu metadata vào `fb_messages` (không có messages[])
- Lưu 2 messages vào `fb_message_items` (mỗi message là 1 document)
- Cập nhật `totalMessages = 2`

---

## ✅ Lợi Ích Tách Biệt

1. **CRUD Routes không bị ảnh hưởng**
   - Logic chung vẫn hoạt động bình thường
   - Không cần thay đổi code CRUD hiện có
   - Tương thích ngược với dữ liệu cũ

2. **Endpoint đặc biệt xử lý riêng**
   - Logic tách messages chỉ ở 1 nơi
   - Dễ maintain và test
   - Không ảnh hưởng đến CRUD

3. **Linh hoạt trong sử dụng**
   - Dùng CRUD cho thao tác thủ công
   - Dùng endpoint đặc biệt cho sync tự động

---

## 🔒 Đảm Bảo Tách Biệt

### Code Structure
```
handler.fb.message.go
├── BaseHandler (CRUD methods - không tách messages)
└── HandleUpsertMessages() (Endpoint đặc biệt - có tách messages)

dto.fb.message.go
├── FbMessageCreateInput (DTO cho CRUD)
└── FbMessageUpsertMessagesInput (DTO cho endpoint đặc biệt)

routes.go
├── registerCRUDRoutes() (CRUD routes)
└── router.Post("/upsert-messages") (Endpoint đặc biệt)
```

### Validation
- CRUD routes: Dùng `FbMessageCreateInput` (không có `HasMore`)
- Endpoint đặc biệt: Dùng `FbMessageUpsertMessagesInput` (có `HasMore`)
- Không thể nhầm lẫn giữa 2 loại endpoint

---

## 📝 Kết Luận

- ✅ **CRUD Routes**: Giữ nguyên logic chung, không tách messages
- ✅ **Endpoint Đặc Biệt**: Logic tách messages riêng biệt
- ✅ **Tách biệt hoàn toàn**: Không ảnh hưởng lẫn nhau
- ✅ **Dễ maintain**: Logic rõ ràng, dễ hiểu
