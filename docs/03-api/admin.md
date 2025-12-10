# Admin APIs

Tài liệu về các API endpoints quản trị hệ thống.

## 📋 Tổng Quan

Tất cả các API admin đều nằm dưới prefix `/api/v1/admin/` và yêu cầu quyền admin.

## 🔐 Endpoints

### 1. Block User

Chặn một user.

**Endpoint:** `POST /api/v1/admin/user/block`

**Authentication:** Cần (Permission: `User.Block`)

**Request Body:**
```json
{
  "userId": "507f1f77bcf86cd799439011",
  "reason": "Violation of terms"
}
```

**Response 200:**
```json
{
  "data": {
    "message": "User blocked successfully"
  },
  "error": null
}
```

### 2. Unblock User

Bỏ chặn một user.

**Endpoint:** `POST /api/v1/admin/user/unblock`

**Authentication:** Cần (Permission: `User.Block`)

**Request Body:**
```json
{
  "userId": "507f1f77bcf86cd799439011"
}
```

**Response 200:**
```json
{
  "data": {
    "message": "User unblocked successfully"
  },
  "error": null
}
```

### 3. Set Role

Gán role cho user.

**Endpoint:** `POST /api/v1/admin/user/role`

**Authentication:** Cần (Permission: `User.SetRole`)

**Request Body:**
```json
{
  "userId": "507f1f77bcf86cd799439011",
  "roleId": "507f1f77bcf86cd799439012"
}
```

**Response 200:**
```json
{
  "data": {
    "message": "Role assigned successfully"
  },
  "error": null
}
```

### 4. Set Administrator

Thiết lập user làm administrator (khi đã có admin).

**Endpoint:** `POST /api/v1/admin/user/set-administrator/:id`

**Authentication:** Cần (Permission: `Init.SetAdmin`)

**Path Parameters:**
- `id`: User ID

**Response 200:**
```json
{
  "data": {
    "message": "Administrator set successfully"
  },
  "error": null
}
```

**Lưu ý:** Endpoint này chỉ hoạt động khi đã có admin trong hệ thống.

## 🔐 Init Endpoints

Các endpoint khởi tạo hệ thống (chỉ hoạt động khi chưa có admin).

### 1. Init Status

Kiểm tra trạng thái khởi tạo.

**Endpoint:** `GET /api/v1/init/status`

**Authentication:** Không cần

**Response 200:**
```json
{
  "data": {
    "hasOrganization": true,
    "hasPermissions": true,
    "hasRoles": true,
    "hasAdmin": false
  }
}
```

### 2. Init Organization

Khởi tạo Organization Root.

**Endpoint:** `POST /api/v1/init/organization`

**Authentication:** Không cần (chỉ khi chưa có admin)

**Response 200:**
```json
{
  "data": {
    "message": "Organization Root đã được khởi tạo thành công"
  }
}
```

### 3. Init Permissions

Khởi tạo Permissions.

**Endpoint:** `POST /api/v1/init/permissions`

**Authentication:** Không cần (chỉ khi chưa có admin)

**Response 200:**
```json
{
  "data": {
    "message": "Permissions đã được khởi tạo thành công"
  }
}
```

### 4. Init Roles

Khởi tạo Roles.

**Endpoint:** `POST /api/v1/init/roles`

**Authentication:** Không cần (chỉ khi chưa có admin)

**Response 200:**
```json
{
  "data": {
    "message": "Roles đã được khởi tạo thành công"
  }
}
```

### 5. Init Admin User

Tạo admin user từ Firebase UID.

**Endpoint:** `POST /api/v1/init/admin-user`

**Authentication:** Không cần (chỉ khi chưa có admin)

**Request Body:**
```json
{
  "firebaseUid": "firebase-user-uid"
}
```

**Response 200:**
```json
{
  "data": {
    "message": "Admin user đã được khởi tạo thành công"
  }
}
```

### 6. Init All

Khởi tạo tất cả (one-click setup).

**Endpoint:** `POST /api/v1/init/all`

**Authentication:** Không cần (chỉ khi chưa có admin)

**Response 200:**
```json
{
  "data": {
    "organization": {"status": "success"},
    "permissions": {"status": "success"},
    "roles": {"status": "success"}
  }
}
```

### 7. Set Administrator (Lần Đầu)

Thiết lập user làm administrator lần đầu (không cần quyền).

**Endpoint:** `POST /api/v1/init/set-administrator/:id`

**Authentication:** Không cần (chỉ khi chưa có admin)

**Path Parameters:**
- `id`: User ID

**Response 200:**
```json
{
  "data": {
    "message": "Administrator set successfully"
  }
}
```

## 📝 Lưu Ý

- Init endpoints chỉ hoạt động khi chưa có admin
- Khi đã có admin, init endpoints trả về 404
- Admin endpoints yêu cầu quyền tương ứng
- Set administrator endpoint có 2 phiên bản:
  - `/init/set-administrator/:id` - Khi chưa có admin (không cần quyền)
  - `/admin/user/set-administrator/:id` - Khi đã có admin (cần quyền `Init.SetAdmin`)

## 📚 Tài Liệu Liên Quan

- [Khởi Tạo Hệ Thống](../01-getting-started/khoi-tao.md)
- [RBAC APIs](rbac.md)
- [User Management APIs](user-management.md)

