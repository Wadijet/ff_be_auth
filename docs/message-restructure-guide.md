# Hướng Dẫn Tổ Chức Lại Collection Message

## 📋 Tổng Quan

Tài liệu này hướng dẫn tổ chức lại collection `message` để:
1. Chỉ lưu mảng `messages[]` từ `panCakeData.messages` (bỏ lớp vỏ ngoài)
2. Merge messages từ phân trang (Pancake API trả về 30 messages/lần)
3. Ghi đè các field khác với dữ liệu mới nhất từ Pancake API

---

## 🎯 Mục Tiêu

1. **Lưu toàn bộ messages** của một conversation (không chỉ 30 messages đầu)
2. **Tránh duplicate** khi merge messages mới
3. **Sắp xếp đúng thứ tự** theo thời gian (`inserted_at`)
4. **Đơn giản hóa** cấu trúc dữ liệu

---

## 📊 Phân Tích Dữ Liệu

### Dữ Liệu Hiện Tại

**Collection `message` hiện tại:**
- Lưu toàn bộ `panCakeData` (bao gồm `messages[]` và nhiều field khác)
- Có ~30 trường ngoài `messages[]`

**Collection `conversation`:**
- Đã có nhiều thông tin về conversation
- Một số field trùng với `message.panCakeData`

### Các Trường Sẽ Bị Thiếu (Nếu Chỉ Lưu `messages[]`)

**Tổng cộng: 30 trường** sẽ bị mất nếu không di chuyển sang `conversation.panCakeData`:

#### 🔴 **CAO - Nên Di Chuyển** (10 trường)
1. `read_watermarks` - Tracking đọc tin nhắn
2. `activities` - Tracking hoạt động
3. `ad_clicks` - Chi tiết click quảng cáo
4. `is_banned`, `banned_count`, `banned_by` - Moderation
5. `notes` - Ghi chú
6. `reports_by_phone`, `reported_count` - Moderation
7. `matched_wa_fb_customers` - Khớp WhatsApp-Facebook

#### ⚠️ **TRUNG BÌNH - Cân Nhắc** (15 trường)
- `last_commented_at`, `can_inbox`, `lives_in`, `global_id`
- `suggested_posts`, `available_for_report_phone_numbers`
- `conv_recent_phone_numbers`, `gender`, `profile_updated_at`
- `birthday`, `recent_phone_numbers`, `post`, `conv_phone_numbers`
- `conv_from`, `conv_customers`

#### ⚠️ **THẤP - Có Thể Bỏ Qua** (5 trường)
- `extra_info`, `app`, `allow_use_data_for_training_ai`
- `comment_count`, `success`

**Lưu ý:** Pancake API trả về đầy đủ các field này mỗi lần gọi, nên có thể ghi đè thay vì di chuyển.

---

## 💡 Giải Pháp: Merge Messages + Ghi Đè Field Khác

### Chiến Lược

**✅ MERGE:**
- `messages[]` - Tích lũy từ phân trang, tránh duplicate

**🔄 GHI ĐÈ:**
- Tất cả field khác - Vì Pancake API trả về đầy đủ mỗi lần

### Lý Do

1. **Pancake API trả về đầy đủ**: Mỗi lần gọi API trả về toàn bộ dữ liệu conversation
2. **Dữ liệu mới nhất**: Ghi đè đảm bảo luôn có dữ liệu mới nhất
3. **Đơn giản**: Dễ implement, dễ maintain
4. **Hiệu quả**: Không cần logic merge phức tạp

---

## 🏗️ Cấu Trúc Model Mới

```go
type FbMessage struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    PageId         string                 `json:"pageId" bson:"pageId" index:"text"`
    PageUsername   string                 `json:"pageUsername" bson:"pageUsername" index:"text"`
    ConversationId string                 `json:"conversationId" bson:"conversationId" index:"unique;text"`
    CustomerId     string                 `json:"customerId" bson:"customerId" index:"text"`
    
    // Chỉ lưu mảng messages (không có lớp vỏ panCakeData)
    Messages       []interface{}          `json:"messages" bson:"messages"`
    
    // Vẫn lưu panCakeData để giữ các field khác (ghi đè mỗi lần sync)
    PanCakeData    map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    
    // Metadata để tracking
    LastSyncedAt   int64                  `json:"lastSyncedAt" bson:"lastSyncedAt"`
    TotalMessages  int64                  `json:"totalMessages" bson:"totalMessages"`
    HasMore        bool                   `json:"hasMore" bson:"hasMore"`
    
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
    UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

---

## 🔄 Endpoint Merge Messages

### Endpoint

**POST** `/api/v1/facebook/message/merge-messages`

### Request

```json
{
  "conversationId": "157725629736743_9350439438393456",
  "pageId": "157725629736743",
  "pageUsername": "Folkformint",
  "customerId": "8b168fa9-4836-4648-a3fd-799c227675a1",
  "panCakeData": {
    "conv_from": {...},
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
    "read_watermarks": [...],
    "activities": [...],
    // ... các field khác
  },
  "hasMore": true
}
```

### Logic Merge

1. **Tìm document** theo `conversationId`
2. **Nếu chưa có** → Tạo mới với `messages[]` và `panCakeData`
3. **Nếu đã có** → Merge:
   - **Merge messages**: Lọc bỏ duplicate theo `message.id`, sắp xếp theo `inserted_at`
   - **Ghi đè panCakeData**: Thay thế toàn bộ (trừ `messages` đã merge)
   - Cập nhật `totalMessages`, `hasMore`, `lastSyncedAt`

---

## 💻 Implementation

### Service Method

```go
// MergeMessages merge messages mới và ghi đè các field khác
func (s *FbMessageService) MergeMessages(
    ctx context.Context,
    conversationId string,
    pageId string,
    pageUsername string,
    customerId string,
    newPanCakeData map[string]interface{}, // Toàn bộ panCakeData mới
    hasMore bool,
) (models.FbMessage, error) {
    filter := bson.M{"conversationId": conversationId}
    var existing models.FbMessage
    err := s.collection.FindOne(ctx, filter).Decode(&existing)
    
    now := time.Now().UnixMilli()
    
    // Extract messages từ newPanCakeData
    newMessages, _ := newPanCakeData["messages"].([]interface{})
    
    // 1. Nếu chưa có document → Tạo mới
    if err == mongo.ErrNoDocuments {
        newDoc := models.FbMessage{
            PageId:         pageId,
            PageUsername:   pageUsername,
            ConversationId: conversationId,
            CustomerId:     customerId,
            Messages:       newMessages,
            PanCakeData:    newPanCakeData,
            TotalMessages:  int64(len(newMessages)),
            HasMore:        hasMore,
            LastSyncedAt:   now,
            CreatedAt:      now,
            UpdatedAt:      now,
        }
        return s.InsertOne(ctx, newDoc)
    }
    
    if err != nil {
        return existing, err
    }
    
    // 2. Merge messages (tránh duplicate)
    existingMessages := existing.Messages
    existingMessageIds := make(map[string]bool)
    
    for _, msg := range existingMessages {
        if msgMap, ok := msg.(map[string]interface{}); ok {
            if id, ok := msgMap["id"].(string); ok {
                existingMessageIds[id] = true
            }
        }
    }
    
    // Lọc messages mới (chưa có)
    var uniqueNewMessages []interface{}
    for _, msg := range newMessages {
        if msgMap, ok := msg.(map[string]interface{}); ok {
            if id, ok := msgMap["id"].(string); ok {
                if !existingMessageIds[id] {
                    uniqueNewMessages = append(uniqueNewMessages, msg)
                }
            }
        }
    }
    
    // Merge messages
    mergedMessages := append(existingMessages, uniqueNewMessages...)
    
    // Sắp xếp theo inserted_at (từ cũ đến mới)
    sort.Slice(mergedMessages, func(i, j int) bool {
        msgI, okI := mergedMessages[i].(map[string]interface{})
        msgJ, okJ := mergedMessages[j].(map[string]interface{})
        if !okI || !okJ {
            return false
        }
        timeI, _ := parseTime(msgI["inserted_at"])
        timeJ, _ := parseTime(msgJ["inserted_at"])
        return timeI < timeJ
    })
    
    // 3. Cập nhật panCakeData: Ghi đè tất cả field khác, nhưng giữ messages đã merge
    updatedPanCakeData := make(map[string]interface{})
    
    // Copy tất cả field từ newPanCakeData
    for k, v := range newPanCakeData {
        updatedPanCakeData[k] = v
    }
    
    // Thay thế messages bằng mergedMessages
    updatedPanCakeData["messages"] = mergedMessages
    
    // 4. Update document
    update := bson.M{
        "$set": bson.M{
            "messages":       mergedMessages,
            "panCakeData":    updatedPanCakeData, // Ghi đè toàn bộ panCakeData
            "totalMessages":  int64(len(mergedMessages)),
            "hasMore":        hasMore,
            "lastSyncedAt":   now,
            "updatedAt":      now,
        },
    }
    
    opts := options.FindOneAndUpdate().
        SetReturnDocument(options.After)
    
    var updated models.FbMessage
    err = s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
    return updated, err
}
```

### Helper Function: Parse Time

```go
func parseTime(timeValue interface{}) (int64, error) {
    if timeStr, ok := timeValue.(string); ok {
        t, err := time.Parse("2006-01-02T15:04:05.000000", timeStr)
        if err != nil {
            return 0, err
        }
        return t.Unix(), nil
    }
    return 0, fmt.Errorf("invalid time format")
}
```

---

## 🔄 Flow Đồng Bộ Messages

### Bước 1: Lấy Messages Từ Pancake API

```go
currentCount := 0
hasMore := true

for hasMore {
    // Gọi Pancake API
    response := pancakeAPI.GetMessages(conversationId, currentCount)
    // → 30 messages mỗi lần
    
    // Bước 2: Merge vào collection
    result, err := messageService.MergeMessages(
        ctx,
        conversationId,
        pageId,
        pageUsername,
        customerId,
        response.PanCakeData, // Toàn bộ panCakeData
        response.HasMore,
    )
    
    // Bước 3: Kiểm tra còn messages không
    hasMore = result.HasMore
    if hasMore {
        currentCount += 30 // Lấy trang tiếp theo
    }
}
```

---

## 📋 Checklist Implementation

- [ ] Cập nhật model `FbMessage` (thêm `Messages`, `LastSyncedAt`, `TotalMessages`, `HasMore`)
- [ ] Tạo method `MergeMessages()` trong `FbMessageService`
- [ ] Tạo DTO `MergeMessagesRequest` và `MergeMessagesResponse`
- [ ] Tạo endpoint `POST /api/v1/facebook/message/merge-messages` trong handler
- [ ] Implement logic merge (tránh duplicate, sắp xếp)
- [ ] Thêm helper function `parseTime()` để sắp xếp messages
- [ ] Viết unit tests cho logic merge
- [ ] Tạo migration script để chuyển đổi dữ liệu cũ (nếu cần)

---

## 🔍 Lưu Ý Kỹ Thuật

### 1. Tránh Duplicate Messages

- Sử dụng `message.id` làm unique key
- Tạo map để check nhanh: `map[messageId]bool`
- Chỉ thêm messages chưa có trong map

### 2. Sắp Xếp Messages

- Sắp xếp theo `inserted_at` (từ cũ đến mới)
- Format: `"2006-01-02T15:04:05.000000"` (ISO 8601)
- Sử dụng `sort.Slice()` trong Go

### 3. Performance

- Sử dụng `FindOneAndUpdate` với `$set` thay vì load toàn bộ rồi update
- Index trên `conversationId` để query nhanh
- Cache `existingMessageIds` map trong memory

### 4. Error Handling

- Xử lý trường hợp `conversationId` không tồn tại
- Xử lý trường hợp `messages` rỗng
- Xử lý lỗi parse time
- Xử lý lỗi duplicate (nếu có unique constraint)

---

## ✅ Tóm Tắt

| Hành Động | Field | Lý Do |
|-----------|-------|-------|
| ✅ **MERGE** | `messages[]` | Phân trang, cần tích lũy |
| 🔄 **GHI ĐÈ** | Tất cả field khác trong `panCakeData` | Pancake API trả về đầy đủ mỗi lần |

**Kết luận:**
- Chỉ merge `messages[]` để tích lũy từ phân trang
- Ghi đè tất cả field khác vì Pancake API trả về đầy đủ
- Đơn giản, hiệu quả, dễ maintain

---

## 📚 Tài Liệu Liên Quan

- [Pancake API Documentation](09-ai-context/pancake-api-context.md)
- [Facebook Integration APIs](03-api/facebook.md)
- [Database Architecture](02-architecture/database.md)
