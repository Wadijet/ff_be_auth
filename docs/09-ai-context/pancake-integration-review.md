# Rà Soát Tích Hợp Pancake API với FolkForm API

## 📋 Tổng Quan

Tài liệu này rà soát các tính năng đã implement và còn thiếu khi tích hợp Pancake API vào FolkForm API.

---

## ✅ ĐÃ CÓ (Đã Implement)

### Cách Thức Nhận Dữ Liệu Từ Pancake

**Pattern hiện tại**: Sử dụng **Upsert CRUD** kết hợp **Data Extraction qua struct tag**

1. **Upsert Endpoint**: `POST /api/v1/{collection}/upsert-one?filter={...}`
   - Filter xác định unique key (ví dụ: `{"pageId": "xxx"}`, `{"conversationId": "xxx"}`)
   - Request body chứa data với `panCakeData`
   - Tự động tạo mới nếu chưa có, cập nhật nếu đã có

2. **Data Extraction**: Tự động extract qua struct tag `extract`
   - Ví dụ: `extract:"PanCakeData\\.id"` → extract từ `panCakeData["id"]`
   - Chạy tự động khi upsert/insert/update

3. **Không cần method `ReviceData()` riêng** - Tất cả đều dùng CRUD operations

### 1. Facebook Pages (FbPage)
- ✅ **Model**: `models.FbPage` với struct tag `extract` cho `pageId`, `pageName`, `pageUsername`
- ✅ **Service**: `FbPageService` với CRUD operations
- ✅ **Handler**: `FbPageHandler` với các endpoint CRUD + Upsert
- ✅ **Endpoints**:
  - `POST /api/v1/facebook/page/upsert-one?filter={"pageId":"xxx"}` - Upsert page
  - `GET /api/v1/facebook/page/find-by-page-id/:id` - Tìm page theo Facebook PageID
  - `PUT /api/v1/facebook/page/update-token` - Cập nhật Page Access Token
- ✅ **Data Extraction**: Tự động extract từ `panCakeData` qua struct tag

### 2. Facebook Posts (FbPost)
- ✅ **Model**: `models.FbPost` với các trường cần thiết
- ✅ **Service**: `FbPostService` với CRUD operations
- ✅ **Handler**: `FbPostHandler` với các endpoint CRUD + Upsert
- ✅ **Endpoints**:
  - `POST /api/v1/facebook/post/upsert-one?filter={"postId":"xxx"}` - Upsert post
  - `GET /api/v1/facebook/post/find-by-post-id/:id` - Tìm post theo Facebook PostID
- ✅ **Data Extraction**: Tự động extract từ `panCakeData`

### 3. Facebook Conversations (FbConversation)
- ✅ **Model**: `models.FbConversation` với struct tag `extract` cho `conversationId`, `customerId`, `panCakeUpdatedAt`
- ✅ **Service**: `FbConversationService` với CRUD operations + `FindAllSortByApiUpdate()`
- ✅ **Handler**: `FbConversationHandler` với các endpoint CRUD + Upsert
- ✅ **Endpoints**:
  - `POST /api/v1/facebook/conversation/upsert-one?filter={"conversationId":"xxx"}` - Upsert conversation
  - `GET /api/v1/facebook/conversation/sort-by-api-update` - Lấy conversations sắp xếp theo thời gian cập nhật API
- ✅ **Data Extraction**: Tự động extract từ `panCakeData` qua struct tag

### 4. Facebook Messages (FbMessage)
- ✅ **Model**: `models.FbMessage` với struct tag `extract` cho `conversationId`
- ✅ **Service**: `FbMessageService` với CRUD operations
- ✅ **Handler**: `FbMessageHandler` với các endpoint CRUD + Upsert
- ✅ **Endpoints**:
  - `POST /api/v1/facebook/message/upsert-one?filter={"conversationId":"xxx","customerId":"yyy"}` - Upsert message
- ✅ **Data Extraction**: Tự động extract từ `panCakeData` qua struct tag

### 5. Pancake Orders (PcOrder)
- ✅ **Model**: `models.PcOrder` với các trường cần thiết
- ✅ **Service**: `PcOrderService` với CRUD operations
- ✅ **Handler**: `PcOrderHandler` với các endpoint CRUD + Upsert
- ✅ **Endpoints**:
  - `POST /api/v1/pancake/order/upsert-one?filter={"pancakeOrderId":"xxx"}` - Upsert order

---

## ⚠️ CÒN THIẾU (Chưa Implement)

### 1. Webhook/Callback Endpoints
**Vấn đề**: Không có endpoint để Pancake gửi dữ liệu đến FolkForm qua webhook/callback.

**Cần thêm**:
- ❌ `POST /api/v1/pancake/webhook/page` - Nhận webhook cho Page updates
- ❌ `POST /api/v1/pancake/webhook/post` - Nhận webhook cho Post updates
- ❌ `POST /api/v1/pancake/webhook/conversation` - Nhận webhook cho Conversation updates
- ❌ `POST /api/v1/pancake/webhook/message` - Nhận webhook cho Message updates
- ❌ `POST /api/v1/pancake/webhook/order` - Nhận webhook cho Order updates

**Cách implement**:
- Webhook handler sẽ parse payload từ Pancake
- Tạo filter dựa trên unique key (pageId, postId, conversationId, etc.)
- Gọi endpoint upsert-one tương ứng hoặc gọi service.Upsert() trực tiếp
- Cần thêm middleware để verify webhook signature từ Pancake (nếu có)

### 2. Statistics Module
**Pancake API có**:
- Ads Campaign Statistics
- Ads Statistics
- Customer Engagement Statistics
- Page Statistics
- Tag Statistics
- User Statistics

**FolkForm chưa có**:
- ❌ Model, Service, Handler cho Statistics
- ❌ Endpoints để lưu trữ và truy vấn statistics từ Pancake

### 3. Customers Module
**Pancake API có**:
- Get Page Customers
- Update Customer
- Add Customer Note
- Update Customer Note
- Delete Customer Note

**FolkForm chưa có**:
- ❌ Model `FbCustomer` hoặc `PcCustomer`
- ❌ Service và Handler để quản lý customers
- ❌ Endpoints để lưu trữ customer data từ Pancake

### 4. Export Data Module
**Pancake API có**:
- Export Conversations from Ads

**FolkForm chưa có**:
- ❌ Endpoint để trigger export từ Pancake
- ❌ Endpoint để nhận và lưu trữ exported data

### 5. Call Logs Module
**Pancake API có**:
- Retrieve Call Logs (SIP Call Logs)

**FolkForm chưa có**:
- ❌ Model `PcCallLog` hoặc `SipCallLog`
- ❌ Service và Handler để quản lý call logs
- ❌ Endpoints để lưu trữ call logs từ Pancake

### 6. Tags Module
**Pancake API có**:
- Get List Tags
- Tag Conversation (đã có trong Conversation API)

**FolkForm chưa có**:
- ❌ Model `PcTag` hoặc `FbTag`
- ❌ Service và Handler để quản lý tags
- ❌ Endpoints để lưu trữ tags từ Pancake

### 7. Users Module
**Pancake API có**:
- Get List of Users
- Update Round Robin Users

**FolkForm chưa có**:
- ❌ Model `PcUser` (khác với User trong Auth module)
- ❌ Service và Handler để quản lý Pancake users
- ❌ Endpoints để lưu trữ Pancake user data

### 8. Page's Contents Module
**Pancake API có**:
- Upload Media Content

**FolkForm chưa có**:
- ❌ Model `PcContent` hoặc `FbContent`
- ❌ Service và Handler để quản lý uploaded contents
- ❌ Endpoints để lưu trữ content metadata từ Pancake

### 9. Conversation Actions
**Pancake API có**:
- Tag Conversation
- Assign Conversation
- Mark as Read
- Mark as Unread

**FolkForm chưa có**:
- ❌ Endpoints để thực hiện các actions này (có thể gọi Pancake API trực tiếp hoặc lưu trạng thái)

---

## 🔧 CẦN BỔ SUNG

### 1. Webhook Handlers (Ưu tiên cao)
Cần thêm các handler để nhận webhook từ Pancake:

```go
// handler.pancake.webhook.go
func (h *PancakeWebhookHandler) HandlePageWebhook(c fiber.Ctx) error {
    // 1. Verify webhook signature (nếu có)
    // 2. Parse payload từ Pancake
    // 3. Tạo filter: {"pageId": payload["id"]}
    // 4. Gọi FbPageService.Upsert() với filter và payload
    //    - Data extraction sẽ tự động chạy qua struct tag extract
}

func (h *PancakeWebhookHandler) HandlePostWebhook(c fiber.Ctx) error {
    // Similar logic với filter: {"postId": payload["id"]}
}

func (h *PancakeWebhookHandler) HandleConversationWebhook(c fiber.Ctx) error {
    // Similar logic với filter: {"conversationId": payload["id"]}
}

func (h *PancakeWebhookHandler) HandleMessageWebhook(c fiber.Ctx) error {
    // Similar logic với filter: {"conversationId": payload["conversation_id"], "customerId": payload["customer_id"]}
}

func (h *PancakeWebhookHandler) HandleOrderWebhook(c fiber.Ctx) error {
    // Similar logic với filter: {"pancakeOrderId": payload["id"]}
}
```

**Lưu ý**: 
- Không cần method `ReviceData()` riêng nữa
- Dùng `Upsert()` từ BaseService với filter và data
- Data extraction tự động qua struct tag `extract`

### 3. Statistics Module (Nếu cần)
Nếu cần lưu trữ statistics từ Pancake:
- Model: `PcStatistics` hoặc `FbStatistics`
- Service: `PcStatisticsService`
- Handler: `PcStatisticsHandler`
- Endpoints: CRUD + webhook để nhận statistics

### 4. Customers Module (Nếu cần)
Nếu cần lưu trữ customer data từ Pancake:
- Model: `PcCustomer` hoặc `FbCustomer`
- Service: `PcCustomerService`
- Handler: `PcCustomerHandler`
- Endpoints: CRUD + webhook để nhận customer updates

### 5. Webhook Verification Middleware
Cần middleware để verify webhook signature từ Pancake (nếu Pancake hỗ trợ):

```go
// middleware.pancake.webhook.go
func VerifyPancakeWebhook(c fiber.Ctx) error {
    // Verify signature
    // Validate payload
    // Continue to handler
}
```

---

## 📊 Bảng So Sánh Chi Tiết

| Module | Pancake API | FolkForm API | Trạng Thái |
|--------|-------------|--------------|------------|
| **Pages** | ✅ List, Generate Token | ✅ CRUD + Find by PageID + Update Token | ✅ **OK** |
| **Posts** | ✅ Get Posts | ✅ CRUD + Find by PostID | ✅ **OK** |
| **Conversations** | ✅ List, Tag, Assign, Read/Unread | ✅ CRUD + Sort by API Update | ⚠️ **Thiếu ReviceData** |
| **Messages** | ✅ Get, Send | ✅ CRUD | ⚠️ **Thiếu ReviceData** |
| **Orders** | ✅ (Không có trong Pancake API doc) | ✅ CRUD | ✅ **OK** |
| **Statistics** | ✅ 6 loại statistics | ❌ Chưa có | ❌ **Thiếu** |
| **Customers** | ✅ CRUD + Notes | ❌ Chưa có | ❌ **Thiếu** |
| **Export Data** | ✅ Export Conversations | ❌ Chưa có | ❌ **Thiếu** |
| **Call Logs** | ✅ Retrieve Call Logs | ❌ Chưa có | ❌ **Thiếu** |
| **Tags** | ✅ Get List Tags | ❌ Chưa có | ❌ **Thiếu** |
| **Users** | ✅ Get List, Update Round Robin | ❌ Chưa có | ❌ **Thiếu** |
| **Page Contents** | ✅ Upload Media | ❌ Chưa có | ❌ **Thiếu** |
| **Webhooks** | ❓ (Cần kiểm tra Pancake có hỗ trợ không) | ❌ Chưa có | ❌ **Thiếu** |

---

## 🎯 Khuyến Nghị

### Ưu Tiên Cao (Cần làm ngay)
1. ✅ **Thêm Webhook Endpoints**: Tạo các endpoint để nhận dữ liệu từ Pancake
   - Sử dụng `Upsert()` với filter dựa trên unique key
   - Data extraction tự động qua struct tag `extract`
2. ✅ **Webhook Verification**: Thêm middleware để verify webhook signature (nếu Pancake hỗ trợ)

### Ưu Tiên Trung Bình (Nếu cần)
4. ⚠️ **Customers Module**: Nếu cần lưu trữ customer data
5. ⚠️ **Statistics Module**: Nếu cần lưu trữ và phân tích statistics

### Ưu Tiên Thấp (Có thể bỏ qua)
6. ⚠️ **Call Logs**: Chỉ cần nếu tích hợp SIP
7. ⚠️ **Tags**: Có thể lưu trong panCakeData
8. ⚠️ **Users**: Có thể không cần lưu riêng
9. ⚠️ **Page Contents**: Có thể chỉ cần lưu metadata

---

## 📝 Ghi Chú

1. **Data Extraction**: Hệ thống sử dụng struct tag `extract` để tự động extract dữ liệu từ `panCakeData`:
   - Format: `extract:"PanCakeData\\.field_path[,converter=name][,optional]"`
   - Tự động chạy khi upsert/insert/update
   - Ví dụ: `extract:"PanCakeData\\.id"` → extract `panCakeData["id"]`

2. **Upsert Pattern**: 
   - Dùng endpoint `POST /api/v1/{collection}/upsert-one?filter={...}`
   - Filter xác định unique key để tìm document
   - Tự động tạo mới nếu chưa có, cập nhật nếu đã có
   - Không cần method `ReviceData()` riêng

3. **Webhook vs Polling**: 
   - Nếu Pancake không hỗ trợ webhook, có thể sử dụng polling (gọi Pancake API định kỳ)
   - Endpoint `sort-by-api-update` đã hỗ trợ việc này cho Conversations
   - Khi polling, gọi `upsert-one` với filter và data từ Pancake

4. **Authentication**: Cần xác định cách Pancake xác thực khi gửi webhook (API key, signature, etc.)

5. **Ví dụ sử dụng Upsert**:
   ```bash
   # Upsert Conversation
   POST /api/v1/facebook/conversation/upsert-one?filter={"conversationId":"conv_123"}
   {
     "pageId": "page_123",
     "pageUsername": "my_page",
     "panCakeData": {
       "id": "conv_123",
       "customer_id": "customer_456",
       "updated_at": "2019-08-24T14:15:22.000000",
       "type": "INBOX"
     }
   }
   # → conversationId, customerId, panCakeUpdatedAt sẽ tự động extract từ panCakeData
   ```

---

**Ngày tạo**: 2025-01-XX  
**Phiên bản**: 1.0
