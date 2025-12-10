# User Management APIs

Tài liệu về các API endpoints quản lý người dùng (CRUD operations).

## 📋 Tổng Quan

Tất cả các API user management đều nằm dưới prefix `/api/v1/user/` và sử dụng CRUD pattern.

## 🔐 Endpoints CRUD

### 1. Insert One

Tạo một user mới.

**Endpoint:** `POST /api/v1/user/insert-one`

**Authentication:** Cần (Permission: `User.Insert`)

**Request Body:**
```json
{
  "firebaseUid": "firebase-user-uid",
  "email": "user@example.com",
  "name": "User Name"
}
```

**Response 200:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "firebaseUid": "firebase-user-uid",
    "email": "user@example.com",
    "name": "User Name"
  },
  "error": null
}
```

### 2. Find

Tìm tất cả users với filter.

**Endpoint:** `GET /api/v1/user/find`

**Authentication:** Cần (Permission: `User.Read`)

**Query Parameters:**
- `filter`: JSON string filter (MongoDB query)
- `sort`: JSON string sort
- `limit`: Số lượng kết quả
- `skip`: Số lượng bỏ qua

**Response 200:**
```json
{
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "firebaseUid": "firebase-user-uid",
      "email": "user@example.com",
      "name": "User Name"
    }
  ],
  "error": null
}
```

### 3. Find One

Tìm một user với filter.

**Endpoint:** `GET /api/v1/user/find-one`

**Authentication:** Cần (Permission: `User.Read`)

**Query Parameters:** Tương tự Find

**Response 200:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "firebaseUid": "firebase-user-uid",
    "email": "user@example.com",
    "name": "User Name"
  },
  "error": null
}
```

### 4. Find By ID

Tìm user theo ID.

**Endpoint:** `GET /api/v1/user/find-by-id/:id`

**Authentication:** Cần (Permission: `User.Read`)

**Path Parameters:**
- `id`: User ID (MongoDB ObjectID)

**Response 200:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "firebaseUid": "firebase-user-uid",
    "email": "user@example.com",
    "name": "User Name"
  },
  "error": null
}
```

### 5. Find By IDs

Tìm nhiều users theo danh sách IDs.

**Endpoint:** `POST /api/v1/user/find-by-ids`

**Authentication:** Cần (Permission: `User.Read`)

**Request Body:**
```json
{
  "ids": ["507f1f77bcf86cd799439011", "507f1f77bcf86cd799439012"]
}
```

**Response 200:**
```json
{
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "firebaseUid": "firebase-user-uid-1",
      "email": "user1@example.com"
    },
    {
      "_id": "507f1f77bcf86cd799439012",
      "firebaseUid": "firebase-user-uid-2",
      "email": "user2@example.com"
    }
  ],
  "error": null
}
```

### 6. Find With Pagination

Tìm users với phân trang.

**Endpoint:** `GET /api/v1/user/find-with-pagination`

**Authentication:** Cần (Permission: `User.Read`)

**Query Parameters:**
- `page`: Số trang (default: 1)
- `limit`: Số lượng mỗi trang (default: 10)
- `filter`: JSON string filter
- `sort`: JSON string sort

**Response 200:**
```json
{
  "data": {
    "items": [
      {
        "_id": "507f1f77bcf86cd799439011",
        "firebaseUid": "firebase-user-uid",
        "email": "user@example.com"
      }
    ],
    "total": 100,
    "page": 1,
    "limit": 10,
    "totalPages": 10
  },
  "error": null
}
```

### 7. Update By ID

Cập nhật user theo ID.

**Endpoint:** `PUT /api/v1/user/update-by-id/:id`

**Authentication:** Cần (Permission: `User.Update`)

**Path Parameters:**
- `id`: User ID

**Request Body:**
```json
{
  "name": "New Name",
  "email": "newemail@example.com"
}
```

**Response 200:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "name": "New Name",
    "email": "newemail@example.com",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "error": null
}
```

### 8. Delete By ID

Xóa user theo ID.

**Endpoint:** `DELETE /api/v1/user/delete-by-id/:id`

**Authentication:** Cần (Permission: `User.Delete`)

**Path Parameters:**
- `id`: User ID

**Response 200:**
```json
{
  "data": {
    "message": "User deleted successfully"
  },
  "error": null
}
```

### 9. Count Documents

Đếm số lượng users.

**Endpoint:** `GET /api/v1/user/count`

**Authentication:** Cần (Permission: `User.Read`)

**Query Parameters:**
- `filter`: JSON string filter

**Response 200:**
```json
{
  "data": {
    "count": 100
  },
  "error": null
}
```

## 📝 Lưu Ý

- Tất cả endpoints đều yêu cầu authentication
- Mỗi endpoint yêu cầu permission tương ứng
- User collection là read-only trong CRUD (chỉ có thể đọc, không thể tạo/sửa/xóa qua CRUD)
- Để tạo user, sử dụng Firebase login endpoint
- Để cập nhật user, sử dụng profile endpoint hoặc admin endpoints

## 📚 Tài Liệu Liên Quan

- [Authentication APIs](authentication.md)
- [Admin APIs](admin.md)
- [RBAC APIs](rbac.md)

