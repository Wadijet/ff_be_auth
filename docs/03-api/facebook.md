# Facebook Integration APIs

Tài liệu về các API endpoints tích hợp Facebook (Pages, Posts, Conversations, Messages).

## 📋 Tổng Quan

Tất cả các API Facebook đều nằm dưới prefix `/api/v1/facebook/` hoặc `/api/v1/access-token/`.

## 🔐 Access Token APIs

Quản lý Facebook Access Tokens.

**Prefix:** `/api/v1/access-token/`

**Endpoints (Full CRUD):**
- `POST /api/v1/access-token/insert-one` - Tạo access token (Permission: `AccessToken.Insert`)
- `GET /api/v1/access-token/find` - Tìm access tokens (Permission: `AccessToken.Read`)
- `GET /api/v1/access-token/find-by-id/:id` - Tìm theo ID (Permission: `AccessToken.Read`)
- `PUT /api/v1/access-token/update-by-id/:id` - Cập nhật (Permission: `AccessToken.Update`)
- `DELETE /api/v1/access-token/delete-by-id/:id` - Xóa (Permission: `AccessToken.Delete`)

## 🔐 Facebook Page APIs

Quản lý Facebook Pages.

**Prefix:** `/api/v1/facebook/page/`

**Endpoints (Full CRUD):**
- `POST /api/v1/facebook/page/insert-one` - Tạo page (Permission: `FbPage.Insert`)
- `GET /api/v1/facebook/page/find` - Tìm pages (Permission: `FbPage.Read`)
- `GET /api/v1/facebook/page/find-by-id/:id` - Tìm theo ID (Permission: `FbPage.Read`)
- `PUT /api/v1/facebook/page/update-by-id/:id` - Cập nhật (Permission: `FbPage.Update`)
- `DELETE /api/v1/facebook/page/delete-by-id/:id` - Xóa (Permission: `FbPage.Delete`)

**Endpoints Đặc Biệt:**
- `GET /api/v1/facebook/page/find-by-page-id/:id` - Tìm page theo Facebook PageID (Permission: `FbPage.Read`)
- `PUT /api/v1/facebook/page/update-token` - Cập nhật Page Access Token (Permission: `FbPage.Update`)

**Request Body cho update-token:**
```json
{
  "pageId": "facebook-page-id",
  "pageAccessToken": "new-page-access-token"
}
```

## 🔐 Facebook Post APIs

Quản lý Facebook Posts.

**Prefix:** `/api/v1/facebook/post/`

**Endpoints (Full CRUD):**
- `POST /api/v1/facebook/post/insert-one` - Tạo post (Permission: `FbPost.Insert`)
- `GET /api/v1/facebook/post/find` - Tìm posts (Permission: `FbPost.Read`)
- `GET /api/v1/facebook/post/find-by-id/:id` - Tìm theo ID (Permission: `FbPost.Read`)
- `PUT /api/v1/facebook/post/update-by-id/:id` - Cập nhật (Permission: `FbPost.Update`)
- `DELETE /api/v1/facebook/post/delete-by-id/:id` - Xóa (Permission: `FbPost.Delete`)

**Endpoints Đặc Biệt:**
- `GET /api/v1/facebook/post/find-by-post-id/:id` - Tìm post theo Facebook PostID (Permission: `FbPost.Read`)
- `PUT /api/v1/facebook/post/update-token` - Cập nhật token của post (Permission: `FbPost.Update`)

**Request Body cho update-token:**
```json
{
  "postId": "facebook-post-id",
  "panCakeData": { /* dữ liệu từ Pancake API */ }
}
```

## 🔐 Facebook Conversation APIs

Quản lý Facebook Conversations.

**Prefix:** `/api/v1/facebook/conversation/`

**Endpoints (Full CRUD):**
- `POST /api/v1/facebook/conversation/insert-one` - Tạo conversation (Permission: `FbConversation.Insert`)
- `GET /api/v1/facebook/conversation/find` - Tìm conversations (Permission: `FbConversation.Read`)
- `GET /api/v1/facebook/conversation/find-by-id/:id` - Tìm theo ID (Permission: `FbConversation.Read`)
- `PUT /api/v1/facebook/conversation/update-by-id/:id` - Cập nhật (Permission: `FbConversation.Update`)
- `DELETE /api/v1/facebook/conversation/delete-by-id/:id` - Xóa (Permission: `FbConversation.Delete`)

### Endpoint Đặc Biệt: Sort By API Update

Lấy conversations sắp xếp theo thời gian cập nhật API.

**Endpoint:** `GET /api/v1/facebook/conversation/sort-by-api-update`

**Authentication:** Cần (Permission: `FbConversation.Read`)

**Response:**
```json
{
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "pageId": "page-id",
      "conversationId": "conversation-id",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## 🔐 Facebook Message APIs

Quản lý Facebook Messages.

**Prefix:** `/api/v1/facebook/message/`

**Endpoints (Full CRUD):**
- `POST /api/v1/facebook/message/insert-one` - Tạo message (Permission: `FbMessage.Insert`)
- `GET /api/v1/facebook/message/find` - Tìm messages (Permission: `FbMessage.Read`)
- `GET /api/v1/facebook/message/find-by-id/:id` - Tìm theo ID (Permission: `FbMessage.Read`)
- `PUT /api/v1/facebook/message/update-by-id/:id` - Cập nhật (Permission: `FbMessage.Update`)
- `DELETE /api/v1/facebook/message/delete-by-id/:id` - Xóa (Permission: `FbMessage.Delete`)

## 📝 Lưu Ý

- Tất cả endpoints đều yêu cầu authentication
- Mỗi endpoint yêu cầu permission tương ứng
- Tất cả collections đều có full CRUD operations
- Facebook integration sử dụng Facebook Graph API

## 📚 Tài Liệu Liên Quan

- [Pancake Integration APIs](pancake.md)
- [Agent Management APIs](agent.md)

