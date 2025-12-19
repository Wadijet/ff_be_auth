# Rà Soát Đồng Bộ Pancake API - Báo Cáo Chi Tiết

**Ngày rà soát:** $(date)  
**Tài liệu tham khảo:** `docs/09-ai-context/pancake-api-context.md`

---

## 📊 Tổng Quan

### ✅ Đã Đồng Bộ (Implemented)

| Module | API Pancake | Trạng Thái FolkForm | Ghi Chú |
|--------|-------------|---------------------|---------|
| **Pages** | List Pages, Generate Page Access Token | ✅ Đã có | Model: `FbPage`, Service, Handler đầy đủ |
| **Posts** | Get Posts | ✅ Đã có | Model: `FbPost`, Service, Handler đầy đủ |
| **Conversations** | List Conversations | ✅ Đã có | Model: `FbConversation`, Service, Handler đầy đủ |
| **Messages** | Get Messages, Send Message | ✅ Đã có | Model: `FbMessage`, `FbMessageItem`, Service, Handler đầy đủ |
| **Customers** | Get Page Customers | ✅ Đã có | Model: `Customer`, Service, Handler đầy đủ |
| **Orders** | - | ✅ Đã có | Model: `PcOrder`, Service, Handler đầy đủ |
| **Access Token** | - | ✅ Đã có | Model: `PcAccessToken`, Service, Handler đầy đủ |

---

## ❌ Chưa Đồng Bộ (Missing)

### 1. Customer Notes (Ghi Chú Khách Hàng) ⚠️ Ưu Tiên Trung Bình

**Pancake API có:**
- ✅ `POST /pages/{page_id}/page_customers/{page_customer_id}/notes` - Thêm ghi chú
- ✅ `PUT /pages/{page_id}/page_customers/{page_customer_id}/notes` - Cập nhật ghi chú
- ✅ `DELETE /pages/{page_id}/page_customers/{page_customer_id}/notes` - Xóa ghi chú

**FolkForm chưa có:**
- ❌ Model `CustomerNote` để lưu trữ ghi chú
- ❌ Service và Handler để quản lý customer notes
- ❌ Endpoints để CRUD customer notes

**Cấu trúc đề xuất:**
```go
type CustomerNote struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    CustomerId  primitive.ObjectID `json:"customerId" bson:"customerId" index:"text"`
    NoteId      string             `json:"noteId" bson:"noteId" index:"unique;text" extract:"PanCakeData\\.id"`
    Message     string             `json:"message" bson:"message" extract:"PanCakeData\\.message"`
    OrderId     string             `json:"orderId" bson:"orderId" extract:"PanCakeData\\.order_id,optional"`
    Images      []string           `json:"images" bson:"images" extract:"PanCakeData\\.images,optional"`
    Links       []string           `json:"links" bson:"links" extract:"PanCakeData\\.links,optional"`
    CreatedBy   map[string]interface{} `json:"createdBy" bson:"createdBy" extract:"PanCakeData\\.created_by,optional"`
    CreatedAt   int64              `json:"createdAt" bson:"createdAt" extract:"PanCakeData\\.created_at,converter=time"`
    UpdatedAt   int64              `json:"updatedAt" bson:"updatedAt" extract:"PanCakeData\\.updated_at,converter=time"`
    RemovedAt   int64              `json:"removedAt" bson:"removedAt" extract:"PanCakeData\\.removed_at,converter=time,optional"`
    PanCakeData map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
}
```

**Endpoints đề xuất:**
- `POST /api/v1/customer-note/upsert-one?filter={"customerId":"xxx","noteId":"yyy"}` - Upsert note
- `GET /api/v1/customer-note/find-by-customer/:customerId` - Lấy tất cả notes của customer
- `DELETE /api/v1/customer-note/delete-by-id/:id` - Xóa note

---

### 2. Statistics (Thống Kê) ⚠️ Ưu Tiên Trung Bình

**Pancake API có:**
- ✅ `GET /pages/{page_id}/statistics/pages_campaign` - Thống kê chiến dịch quảng cáo
- ✅ `GET /pages/{page_id}/statistics/ads` - Thống kê quảng cáo
- ✅ `GET /pages/{page_id}/statistics/customer_engagements` - Thống kê tương tác khách hàng
- ✅ `GET /pages/{page_id}/statistics/pages` - Thống kê trang
- ✅ `GET /pages/{page_id}/statistics/tags` - Thống kê tags
- ✅ `GET /pages/{page_id}/statistics/users` - Thống kê người dùng

**FolkForm chưa có:**
- ❌ Model `PcStatistics` hoặc các model riêng cho từng loại statistics
- ❌ Service và Handler để lưu trữ và truy vấn statistics
- ❌ Endpoints để sync statistics từ Pancake

**Cấu trúc đề xuất:**
```go
type PcStatistics struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    PageId      string             `json:"pageId" bson:"pageId" index:"text"`
    Type        string             `json:"type" bson:"type" index:"text"` // "pages_campaign", "ads", "customer_engagements", "pages", "tags", "users"
    Date        string             `json:"date" bson:"date" index:"text"` // YYYY-MM-DD
    PanCakeData map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt   int64              `json:"createdAt" bson:"createdAt"`
    UpdatedAt   int64              `json:"updatedAt" bson:"updatedAt"`
}
```

**Endpoints đề xuất:**
- `POST /api/v1/pancake/statistics/upsert-one?filter={"pageId":"xxx","type":"ads","date":"2024-01-01"}` - Upsert statistics
- `GET /api/v1/pancake/statistics/find?pageId=xxx&type=ads&date=2024-01-01` - Lấy statistics

---

### 3. Conversation Actions (Hành Động Cuộc Hội Thoại) ⚠️ Ưu Tiên Trung Bình

**Pancake API có:**
- ✅ `POST /pages/{page_id}/conversations/{conversation_id}/tags` - Gán tag cho conversation
- ✅ `POST /pages/{page_id}/conversations/{conversation_id}/assign` - Gán conversation cho user
- ✅ `POST /pages/{page_id}/conversations/{conversation_id}/read` - Đánh dấu đã đọc
- ✅ `POST /pages/{page_id}/conversations/{conversation_id}/unread` - Đánh dấu chưa đọc

**FolkForm chưa có:**
- ❌ Endpoints để thực hiện các actions này (có thể gọi Pancake API trực tiếp hoặc lưu trạng thái)

**Endpoints đề xuất:**
- `POST /api/v1/facebook/conversation/:conversationId/tag` - Gán tag
- `POST /api/v1/facebook/conversation/:conversationId/assign` - Gán cho user
- `POST /api/v1/facebook/conversation/:conversationId/mark-read` - Đánh dấu đã đọc
- `POST /api/v1/facebook/conversation/:conversationId/mark-unread` - Đánh dấu chưa đọc

**Lưu ý:** Các endpoints này có thể:
1. Gọi trực tiếp Pancake API và cập nhật local database
2. Hoặc chỉ cập nhật local database nếu đã có webhook từ Pancake

---

### 4. Tags (Thẻ) ⚠️ Ưu Tiên Thấp

**Pancake API có:**
- ✅ `GET /pages/{page_id}/tags` - Lấy danh sách tags

**FolkForm chưa có:**
- ❌ Model `PcTag` hoặc `FbTag`
- ❌ Service và Handler để quản lý tags
- ❌ Endpoints để lưu trữ tags từ Pancake

**Cấu trúc đề xuất:**
```go
type FbTag struct {
    ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    PageId       string             `json:"pageId" bson:"pageId" index:"text"`
    TagId        int                `json:"tagId" bson:"tagId" index:"unique" extract:"PanCakeData\\.id"`
    Text         string             `json:"text" bson:"text" extract:"PanCakeData\\.text"`
    Color        string             `json:"color" bson:"color" extract:"PanCakeData\\.color"`
    LightenColor string             `json:"lightenColor" bson:"lightenColor" extract:"PanCakeData\\.lighten_color"`
    PanCakeData  map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt    int64              `json:"createdAt" bson:"createdAt"`
    UpdatedAt    int64              `json:"updatedAt" bson:"updatedAt"`
}
```

**Endpoints đề xuất:**
- `POST /api/v1/facebook/tag/upsert-one?filter={"pageId":"xxx","tagId":123}` - Upsert tag
- `GET /api/v1/facebook/tag/find-by-page/:pageId` - Lấy tất cả tags của page

**Khuyến nghị:**
- Tags có thể lưu trong `panCakeData` của conversations nếu không cần query riêng
- Nếu cần query/filter theo tags → Nên implement riêng

---

### 5. Users (Người Dùng Pancake) ⚠️ Ưu Tiên Thấp

**Pancake API có:**
- ✅ `GET /pages/{page_id}/users` - Lấy danh sách users
- ✅ `POST /pages/{page_id}/round_robin_users` - Cập nhật round robin users

**FolkForm chưa có:**
- ❌ Model `PcUser` (khác với User trong Auth module)
- ❌ Service và Handler để quản lý Pancake users
- ❌ Endpoints để lưu trữ Pancake user data

**Cấu trúc đề xuất:**
```go
type PcUser struct {
    ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    PageId          string             `json:"pageId" bson:"pageId" index:"text"`
    PancakeUserId   string             `json:"pancakeUserId" bson:"pancakeUserId" index:"unique;text" extract:"PanCakeData\\.id"`
    Name            string             `json:"name" bson:"name" extract:"PanCakeData\\.name"`
    FbId            string             `json:"fbId" bson:"fbId" extract:"PanCakeData\\.fb_id"`
    Status          string             `json:"status" bson:"status" extract:"PanCakeData\\.status"`
    StatusInPage    string             `json:"statusInPage" bson:"statusInPage" extract:"PanCakeData\\.status_in_page"`
    IsOnline        bool               `json:"isOnline" bson:"isOnline" extract:"PanCakeData\\.is_online"`
    PagePermissions map[string]interface{} `json:"pagePermissions" bson:"pagePermissions" extract:"PanCakeData\\.page_permissions,optional"`
    PanCakeData     map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt       int64              `json:"createdAt" bson:"createdAt"`
    UpdatedAt       int64              `json:"updatedAt" bson:"updatedAt"`
}
```

**Endpoints đề xuất:**
- `POST /api/v1/pancake/user/upsert-one?filter={"pageId":"xxx","pancakeUserId":"yyy"}` - Upsert user
- `GET /api/v1/pancake/user/find-by-page/:pageId` - Lấy tất cả users của page
- `POST /api/v1/pancake/user/update-round-robin` - Cập nhật round robin users

**Khuyến nghị:**
- Pancake users khác với FolkForm users (Auth module)
- Chỉ cần nếu cần quản lý users của Pancake (assign conversations, round robin)
- Có thể lưu trong `panCakeData` nếu không cần query riêng

---

### 6. Export Data (Xuất Dữ Liệu) ⚠️ Ưu Tiên Thấp

**Pancake API có:**
- ✅ `GET /pages/{page_id}/export_data?action=conversations_from_ads&since=xxx&until=yyy&offset=0` - Export conversations từ ads

**FolkForm chưa có:**
- ❌ Endpoint để trigger export từ Pancake
- ❌ Endpoint để nhận và lưu trữ exported data

**Khuyến nghị:**
- Có thể không cần nếu đã có sync conversations thông qua API thông thường
- Nếu cần export hàng loạt → Có thể implement như một job/background task
- Endpoint có thể gọi Pancake API và tự động upsert conversations vào database

**Endpoint đề xuất:**
- `POST /api/v1/pancake/export/conversations-from-ads` - Trigger export và sync conversations

---

### 7. Call Logs (Nhật Ký Cuộc Gọi) ⚠️ Ưu Tiên Thấp

**Pancake API có:**
- ✅ `GET /pages/{page_id}/sip_call_logs?id=SIP_PACKAGE_ID&page_number=1&page_size=30&since=xxx&until=yyy` - Lấy call logs

**FolkForm chưa có:**
- ❌ Model `PcCallLog` hoặc `SipCallLog`
- ❌ Service và Handler để quản lý call logs
- ❌ Endpoints để lưu trữ call logs từ Pancake

**Cấu trúc đề xuất:**
```go
type PcCallLog struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    PageId      string             `json:"pageId" bson:"pageId" index:"text"`
    CallId      string             `json:"callId" bson:"callId" index:"unique;text" extract:"PanCakeData\\.call_id"`
    Caller      string             `json:"caller" bson:"caller" extract:"PanCakeData\\.caller"`
    Callee      string             `json:"callee" bson:"callee" extract:"PanCakeData\\.callee"`
    StartTime   int64              `json:"startTime" bson:"startTime" extract:"PanCakeData\\.start_time,converter=time"`
    Duration    int                `json:"duration" bson:"duration" extract:"PanCakeData\\.duration"`
    Status      string             `json:"status" bson:"status" extract:"PanCakeData\\.status"`
    PanCakeData map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt   int64              `json:"createdAt" bson:"createdAt"`
    UpdatedAt   int64              `json:"updatedAt" bson:"updatedAt"`
}
```

**Endpoints đề xuất:**
- `POST /api/v1/pancake/call-log/upsert-one?filter={"pageId":"xxx","callId":"yyy"}` - Upsert call log
- `GET /api/v1/pancake/call-log/find-by-page/:pageId` - Lấy call logs của page

**Khuyến nghị:**
- Chỉ cần nếu tích hợp SIP/VoIP
- Nếu không dùng SIP → Có thể bỏ qua

---

### 8. Page's Contents (Nội Dung Trang) ⚠️ Ưu Tiên Thấp

**Pancake API có:**
- ✅ `POST /pages/{page_id}/upload_contents` - Upload media content (hình ảnh, video)

**FolkForm chưa có:**
- ❌ Model `PcContent` hoặc `FbContent`
- ❌ Service và Handler để quản lý uploaded contents
- ❌ Endpoints để lưu trữ content metadata từ Pancake

**Cấu trúc đề xuất:**
```go
type PcContent struct {
    ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    PageId        string             `json:"pageId" bson:"pageId" index:"text"`
    ContentId     string             `json:"contentId" bson:"contentId" index:"unique;text" extract:"PanCakeData\\.id"`
    AttachmentType string            `json:"attachmentType" bson:"attachmentType" extract:"PanCakeData\\.attachment_type"`
    PanCakeData   map[string]interface{} `json:"panCakeData" bson:"panCakeData"`
    CreatedAt     int64              `json:"createdAt" bson:"createdAt"`
    UpdatedAt     int64              `json:"updatedAt" bson:"updatedAt"`
}
```

**Endpoints đề xuất:**
- `POST /api/v1/pancake/content/upsert-one?filter={"pageId":"xxx","contentId":"yyy"}` - Upsert content metadata
- `GET /api/v1/pancake/content/find-by-page/:pageId` - Lấy contents của page

**Khuyến nghị:**
- Chỉ cần lưu metadata (content_id, attachment_type)
- File thực tế được lưu trên Pancake server
- Có thể không cần nếu không cần query riêng

---

### 9. Webhooks (Webhook Handlers) ⚠️ Ưu Tiên Cao

**Pancake API có:**
- ❓ Cần kiểm tra Pancake có hỗ trợ webhook không

**FolkForm chưa có:**
- ❌ Webhook handlers để nhận dữ liệu từ Pancake
- ❌ Middleware để verify webhook signature

**Endpoints đề xuất:**
- `POST /api/v1/pancake/webhook/page` - Nhận webhook cho Page updates
- `POST /api/v1/pancake/webhook/post` - Nhận webhook cho Post updates
- `POST /api/v1/pancake/webhook/conversation` - Nhận webhook cho Conversation updates
- `POST /api/v1/pancake/webhook/message` - Nhận webhook cho Message updates
- `POST /api/v1/pancake/webhook/customer` - Nhận webhook cho Customer updates
- `POST /api/v1/pancake/webhook/order` - Nhận webhook cho Order updates

**Cách implement:**
- Webhook handler sẽ parse payload từ Pancake
- Tạo filter dựa trên unique key (pageId, postId, conversationId, etc.)
- Gọi endpoint upsert-one tương ứng hoặc gọi service.Upsert() trực tiếp
- Cần thêm middleware để verify webhook signature từ Pancake (nếu có)

---

## 📋 Tóm Tắt Ưu Tiên

### Ưu Tiên Cao (Cần làm ngay)
1. ✅ **Webhook Handlers** - Tạo các endpoint để nhận dữ liệu từ Pancake
   - Sử dụng `Upsert()` với filter dựa trên unique key
   - Data extraction tự động qua struct tag `extract`

### Ưu Tiên Trung Bình (Nếu cần)
2. ⚠️ **Customer Notes** - Nếu cần quản lý ghi chú khách hàng
3. ⚠️ **Statistics** - Nếu cần lưu trữ và phân tích statistics
4. ⚠️ **Conversation Actions** - Nếu cần thực hiện các actions (tag, assign, read/unread)

### Ưu Tiên Thấp (Có thể bỏ qua)
5. ⚠️ **Tags** - Có thể lưu trong panCakeData của conversations
6. ⚠️ **Users** - Có thể không cần lưu riêng (khác với Auth users)
7. ⚠️ **Call Logs** - Chỉ cần nếu tích hợp SIP
8. ⚠️ **Page Contents** - Có thể chỉ cần lưu metadata
9. ⚠️ **Export Data** - Có thể không cần nếu đã có sync thông thường

---

## 📝 Ghi Chú

1. **Data Extraction**: Hệ thống sử dụng struct tag `extract` để tự động extract dữ liệu từ `panCakeData`:
   - Format: `extract:"PanCakeData\\.field_path[,converter=name][,optional]"`
   - Tự động chạy khi upsert/insert/update
   - Ví dụ: `extract:"PanCakeData\\.id"` → extract `panCakeData["id"]`

2. **Upsert Pattern**: 
   - Dùng endpoint `POST /api/v1/{collection}/upsert-one?filter={...}`
   - Filter dựa trên unique key (pageId, postId, conversationId, etc.)
   - Body chứa `panCakeData` đầy đủ từ Pancake API

3. **Permissions**: Cần thêm permissions mới cho các module mới:
   - `CustomerNote.*`
   - `PcStatistics.*`
   - `FbTag.*`
   - `PcUser.*`
   - `PcCallLog.*`
   - `PcContent.*`

---

## 🔗 Liên Kết

- **Pancake API Documentation:** https://developer.pancake.biz/
- **Tài liệu Pancake API Context:** `docs/09-ai-context/pancake-api-context.md`
- **Tài liệu FolkForm API Context:** `docs/09-ai-context/folkform-api-context.md`
- **Review Integration:** `docs/09-ai-context/pancake-integration-review.md`
- **Sync Review:** `docs/09-ai-context/pancake-folkform-sync-review.md`


