# Phân Tích Hiệu Suất: Vấn Đề Document Quá Lớn

## ⚠️ Vấn Đề

Khi merge messages vào một document, document sẽ ngày càng lớn:
- **100 messages** → ~500KB - 1MB
- **1000 messages** → ~5MB - 10MB
- **10000 messages** → **Vượt quá giới hạn 16MB của MongoDB!**

### Hậu Quả

1. **Giới hạn MongoDB**: Document không thể vượt quá **16MB**
2. **Hiệu suất query**: Load document lớn rất chậm
3. **Memory**: Tốn nhiều RAM khi load document
4. **Network**: Transfer document lớn chậm
5. **Update performance**: Update document lớn rất chậm

---

## 💡 Giải Pháp Đề Xuất

### Phương Án 1: Lưu Từng Message Riêng Lẻ ⭐ KHUYẾN NGHỊ

**Cấu trúc:**
- Mỗi message là **1 document riêng** trong collection `fb_messages`
- Collection `fb_conversations` chỉ lưu metadata (không lưu messages)

**Ưu điểm:**
- ✅ Không có giới hạn số lượng messages
- ✅ Query nhanh (có thể index theo conversationId, inserted_at)
- ✅ Update/Delete message đơn lẻ dễ dàng
- ✅ Phân trang tự nhiên
- ✅ Không ảnh hưởng hiệu suất khi có nhiều messages

**Nhược điểm:**
- ⚠️ Nhiều documents hơn (nhưng MongoDB xử lý tốt)
- ⚠️ Cần join/aggregate để lấy tất cả messages của conversation

**Model:**
```go
type FbMessage struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    PageId         string                 `json:"pageId" bson:"pageId" index:"text"`
    PageUsername   string                 `json:"pageUsername" bson:"pageUsername" index:"text"`
    ConversationId string                 `json:"conversationId" bson:"conversationId" index:"text"` // Không unique nữa
    CustomerId     string                 `json:"customerId" bson:"customerId" index:"text"`
    
    // Mỗi message là 1 document
    MessageId      string                 `json:"messageId" bson:"messageId" index:"unique"` // ID của message từ Pancake
    MessageData    map[string]interface{} `json:"messageData" bson:"messageData"` // Dữ liệu của message
    
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
    UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

**Index:**
- `conversationId` + `inserted_at` (compound index) để query nhanh
- `messageId` (unique) để tránh duplicate

---

### Phương Án 2: Hybrid - Messages Gần Đây + Archive Cũ

**Cấu trúc:**
- Collection `fb_messages`: Lưu **100-200 messages gần đây nhất** trong document
- Collection `fb_message_archive`: Lưu messages cũ (từng message riêng)

**Logic:**
- Khi merge, nếu messages > 200 → Di chuyển messages cũ vào archive
- Query: Lấy từ `fb_messages` + `fb_message_archive` và merge

**Ưu điểm:**
- ✅ Document không quá lớn
- ✅ Messages gần đây query nhanh
- ✅ Có thể archive messages cũ

**Nhược điểm:**
- ⚠️ Logic phức tạp hơn
- ⚠️ Cần merge từ 2 collections khi query

---

### Phương Án 3: Pagination Trong Messages

**Cấu trúc:**
- Chia messages thành nhiều "chunks" (mỗi chunk 100-200 messages)
- Mỗi chunk là 1 document với `chunkIndex`

**Model:**
```go
type FbMessageChunk struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    ConversationId string                 `json:"conversationId" bson:"conversationId" index:"text"`
    ChunkIndex     int                    `json:"chunkIndex" bson:"chunkIndex"` // 0, 1, 2, ...
    Messages       []interface{}          `json:"messages" bson:"messages"` // 100-200 messages
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
}
```

**Ưu điểm:**
- ✅ Document không quá lớn
- ✅ Có thể query từng chunk

**Nhược điểm:**
- ⚠️ Cần query nhiều chunks để lấy toàn bộ messages
- ⚠️ Logic phức tạp hơn

---

## 🎯 So Sánh Các Phương Án

| Tiêu Chí | Phương Án 1 (Từng Message) | Phương Án 2 (Hybrid) | Phương Án 3 (Chunks) |
|----------|---------------------------|---------------------|---------------------|
| **Scalability** | ⭐⭐⭐ Tốt nhất | ⭐⭐ Trung bình | ⭐⭐ Trung bình |
| **Query Performance** | ⭐⭐⭐ Tốt (có index) | ⭐⭐ Trung bình | ⭐⭐ Trung bình |
| **Độ phức tạp** | ⭐⭐⭐ Đơn giản | ⚠️ Phức tạp | ⚠️ Phức tạp |
| **Update message** | ⭐⭐⭐ Dễ dàng | ⚠️ Khó | ⚠️ Khó |
| **Phân trang** | ⭐⭐⭐ Tự nhiên | ⚠️ Cần merge | ⚠️ Cần merge |
| **Storage** | ⚠️ Nhiều documents | ⭐⭐ Trung bình | ⭐⭐ Trung bình |
| **Khuyến nghị** | ✅ **Nên dùng** | ⚠️ Có thể dùng | ❌ Không nên |

---

## ✅ Kết Luận & Khuyến Nghị

### **KHUYẾN NGHỊ: Phương Án 1 - Lưu Từng Message Riêng Lẻ**

**Lý do:**
1. **Scalability tốt nhất**: Không có giới hạn số lượng messages
2. **Performance tốt**: Query nhanh với index, không cần load document lớn
3. **Đơn giản**: Logic rõ ràng, dễ maintain
4. **Linh hoạt**: Dễ update/delete message đơn lẻ
5. **Phân trang tự nhiên**: MongoDB hỗ trợ tốt

**Trade-off:**
- Nhiều documents hơn, nhưng MongoDB được thiết kế để xử lý điều này
- Cần query với filter `conversationId` để lấy messages, nhưng có index nên nhanh

---

## 📝 Implementation Mới

### Model Mới

```go
type FbMessage struct {
    ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
    PageId         string                 `json:"pageId" bson:"pageId" index:"text"`
    PageUsername   string                 `json:"pageUsername" bson:"pageUsername" index:"text"`
    ConversationId string                 `json:"conversationId" bson:"conversationId" index:"text"` // Không unique
    CustomerId     string                 `json:"customerId" bson:"customerId" index:"text"`
    
    // ID của message từ Pancake (unique)
    MessageId      string                 `json:"messageId" bson:"messageId" index:"unique;text" extract:"MessageData\\.id"`
    
    // Dữ liệu của message (từ panCakeData.messages[])
    MessageData    map[string]interface{} `json:"messageData" bson:"messageData"`
    
    // Extract inserted_at để sort
    InsertedAt     int64                  `json:"insertedAt" bson:"insertedAt" extract:"MessageData\\.inserted_at,converter=time,format=2006-01-02T15:04:05.000000"`
    
    CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
    UpdatedAt      int64                  `json:"updatedAt" bson:"updatedAt"`
}
```

### Service Method: Upsert Messages

```go
// UpsertMessages upsert nhiều messages (mỗi message là 1 document)
func (s *FbMessageService) UpsertMessages(
    ctx context.Context,
    conversationId string,
    pageId string,
    pageUsername string,
    customerId string,
    messages []interface{}, // Mảng messages từ panCakeData.messages
) (int, error) {
    if len(messages) == 0 {
        return 0, nil
    }
    
    var operations []mongo.WriteModel
    now := time.Now().UnixMilli()
    
    for _, msg := range messages {
        msgMap, ok := msg.(map[string]interface{})
        if !ok {
            continue
        }
        
        // Extract messageId
        messageId, ok := msgMap["id"].(string)
        if !ok || messageId == "" {
            continue
        }
        
        // Tạo document cho message
        doc := models.FbMessage{
            PageId:         pageId,
            PageUsername:   pageUsername,
            ConversationId: conversationId,
            CustomerId:     customerId,
            MessageId:      messageId,
            MessageData:    msgMap,
            CreatedAt:      now,
            UpdatedAt:      now,
        }
        
        // Extract inserted_at
        if insertedAtStr, ok := msgMap["inserted_at"].(string); ok {
            if t, err := time.Parse("2006-01-02T15:04:05.000000", insertedAtStr); err == nil {
                doc.InsertedAt = t.Unix()
            }
        }
        
        // Convert to map
        docMap, err := utility.ToMap(doc)
        if err != nil {
            continue
        }
        
        // Tạo upsert operation
        filter := bson.M{"messageId": messageId}
        update := bson.M{
            "$set": docMap,
            "$setOnInsert": bson.M{
                "createdAt": now,
            },
        }
        
        operation := mongo.NewUpdateOneModel().
            SetFilter(filter).
            SetUpdate(update).
            SetUpsert(true)
        
        operations = append(operations, operation)
    }
    
    if len(operations) == 0 {
        return 0, nil
    }
    
    // Bulk write
    opts := options.BulkWrite().SetOrdered(false)
    result, err := s.collection.BulkWrite(ctx, operations, opts)
    if err != nil {
        return 0, common.ConvertMongoError(err)
    }
    
    return int(result.UpsertedCount + result.ModifiedCount), nil
}
```

### Endpoint Mới

**POST** `/api/v1/facebook/message/upsert-messages`

**Request:**
```json
{
  "conversationId": "157725629736743_9350439438393456",
  "pageId": "157725629736743",
  "pageUsername": "Folkformint",
  "customerId": "8b168fa9-4836-4648-a3fd-799c227675a1",
  "messages": [
    {
      "id": "m_xxx1",
      "conversation_id": "157725629736743_9350439438393456",
      "message": "<div>Message 1</div>",
      "inserted_at": "2025-12-16T15:22:45.000000",
      ...
    },
    // ... 30 messages
  ]
}
```

**Logic:**
- Upsert từng message (mỗi message là 1 document)
- Sử dụng `messageId` làm unique key
- Tự động tránh duplicate

### Query Messages

**GET** `/api/v1/facebook/message/find-by-conversation/:conversationId?page=1&limit=50`

**Logic:**
- Query với filter `conversationId`
- Sort theo `insertedAt`
- Phân trang tự nhiên

---

## 📊 So Sánh Performance

### Scenario: Conversation có 10,000 messages

| Phương Án | Document Size | Query Time | Memory Usage |
|-----------|--------------|------------|--------------|
| **Merge vào 1 document** | ❌ **>16MB (FAIL!)** | ❌ Rất chậm | ❌ Rất cao |
| **Từng message riêng** | ✅ ~5-10KB/doc | ✅ Nhanh (có index) | ✅ Thấp |
| **Hybrid** | ✅ ~1-2MB | ⚠️ Trung bình | ⚠️ Trung bình |
| **Chunks** | ✅ ~1-2MB/chunk | ⚠️ Trung bình | ⚠️ Trung bình |

---

## ✅ Kết Luận Cuối Cùng

**Nên thay đổi chiến lược:**

1. **Lưu từng message riêng lẻ** thay vì merge vào 1 document
2. **Upsert messages** thay vì merge
3. **Query với filter** `conversationId` để lấy messages
4. **Index** trên `conversationId` + `insertedAt` để query nhanh

**Lợi ích:**
- ✅ Không có giới hạn số lượng messages
- ✅ Performance tốt ngay cả với hàng ngàn messages
- ✅ Dễ dàng update/delete message đơn lẻ
- ✅ Phân trang tự nhiên
- ✅ Scalable cho tương lai
