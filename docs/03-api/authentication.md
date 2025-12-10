# Authentication APIs

Tài liệu về các API endpoints liên quan đến xác thực và quản lý profile người dùng.

## 📋 Tổng Quan

Tất cả các API authentication đều nằm dưới prefix `/api/v1/auth/`.

## 🔐 Endpoints

### 1. Đăng Nhập với Firebase

Đăng nhập bằng Firebase ID token và nhận JWT token của hệ thống.

**Endpoint:** `POST /api/v1/auth/login/firebase`

**Authentication:** Không cần

**Request Body:**
```json
{
  "idToken": "string",  // Firebase ID token
  "hwid": "string"      // Hardware ID (optional)
}
```

**Response 200:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "firebaseUid": "firebase-user-uid",
    "email": "user@example.com",
    "name": "User Name",
    "token": "jwt-token-here",
    "roles": ["role-id-1", "role-id-2"]
  },
  "error": null
}
```

**Lỗi:**
- `400`: Invalid input
- `401`: Invalid Firebase token

### 2. Đăng Xuất

Đăng xuất và xóa JWT token.

**Endpoint:** `POST /api/v1/auth/logout`

**Authentication:** Cần (Bearer Token)

**Request Body:**
```json
{
  "hwid": "string"  // Optional
}
```

**Response 200:**
```json
{
  "data": {
    "message": "Logged out successfully"
  },
  "error": null
}
```

### 3. Lấy Profile

Lấy thông tin profile của người dùng hiện tại.

**Endpoint:** `GET /api/v1/auth/profile`

**Authentication:** Cần (Bearer Token)

**Response 200:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "firebaseUid": "firebase-user-uid",
    "email": "user@example.com",
    "name": "User Name",
    "phone": "+84123456789",
    "avatarUrl": "https://example.com/avatar.jpg",
    "emailVerified": true,
    "phoneVerified": false,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "error": null
}
```

**Lưu ý:** Response không bao gồm password, salt, và tokens.

### 4. Cập Nhật Profile

Cập nhật thông tin profile của người dùng hiện tại.

**Endpoint:** `PUT /api/v1/auth/profile`

**Authentication:** Cần (Bearer Token)

**Request Body:**
```json
{
  "name": "New Name"  // Các trường khác tùy chọn
}
```

**Response 200:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "firebaseUid": "firebase-user-uid",
    "email": "user@example.com",
    "name": "New Name",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "error": null
}
```

### 5. Lấy Roles của User

Lấy danh sách tất cả các role mà user hiện tại có.

**Endpoint:** `GET /api/v1/auth/roles`

**Authentication:** Cần (Bearer Token)

**Response 200:**
```json
{
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "name": "Administrator",
      "code": "ADMIN",
      "organizationId": "507f1f77bcf86cd799439012",
      "permissions": ["permission-id-1", "permission-id-2"]
    }
  ],
  "error": null
}
```

## 🔒 Authentication Header

Tất cả các endpoint (trừ login) yêu cầu header:

```
Authorization: Bearer <jwt-token>
```

## 📝 Response Format

Tất cả responses đều theo format:

```json
{
  "data": <response-data>,
  "error": <error-object-or-null>
}
```

**Error Object:**
```json
{
  "code": "ERROR_CODE",
  "message": "Error message",
  "status": 400,
  "details": {}
}
```

## 🐛 Error Codes

- `ErrCodeAuth`: Lỗi xác thực
- `ErrCodeValidationFormat`: Lỗi format input
- `ErrCodeInternalServer`: Lỗi server

## 📚 Tài Liệu Liên Quan

- [Authentication Flow](../02-architecture/authentication.md)
- [Firebase Authentication](../firebase-auth-voi-database.md)

