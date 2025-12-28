# AI Context - FolkForm Auth Backend API

> **Tài liệu này cung cấp thông tin đầy đủ về hệ thống API Backend để Frontend có thể tích hợp và sử dụng.**

## 📋 Mục Lục

1. [Tổng Quan Hệ Thống](#tổng-quan-hệ-thống)
2. [Base URL và Cấu Trúc API](#base-url-và-cấu-trúc-api)
3. [Authentication Flow](#authentication-flow)
4. [Response Format](#response-format)
5. [Error Handling](#error-handling)
6. [CRUD Operations Pattern](#crud-operations-pattern)
7. [Organization Context](#organization-context)
8. [RBAC và Permissions](#rbac-và-permissions)
8. [Các Module Chính](#các-module-chính)
9. [Ví Dụ Sử Dụng](#ví-dụ-sử-dụng)

---

## Tổng Quan Hệ Thống

### Mô Tả

FolkForm Auth Backend là hệ thống backend cung cấp:
- 🔐 **Firebase Authentication**: Đăng nhập đa phương thức (Email/Password, Google, Facebook, Phone OTP)
- 👥 **Quản lý Người Dùng**: Tự động tạo user từ Firebase, quản lý profile
- 🔑 **RBAC (Role-Based Access Control)**: Hệ thống phân quyền theo vai trò và tổ chức
- 🏢 **Quản lý Tổ chức**: Cấu trúc tổ chức theo cây (Organization Tree)
- 📱 **Tích hợp Facebook**: Quản lý pages, posts, conversations, messages
- 🛒 **Tích hợp Pancake**: Quản lý đơn hàng
- 🤖 **Quản lý Agent**: Hệ thống trợ lý tự động với check-in/check-out
- 📬 **Notification System**: Hệ thống thông báo đa kênh

### Công Nghệ

- **Framework**: Fiber v3 (Go)
- **Database**: MongoDB
- **Authentication**: Firebase Authentication + JWT
- **API Version**: v1
- **Base Path**: `/api/v1`

---

## Base URL và Cấu Trúc API

### Base URL

```
Development: http://localhost:8080
Production: https://api.folkform.com
```

### Cấu Trúc Endpoint

Tất cả API endpoints đều có prefix: `/api/v1`

**Ví dụ:**
- Authentication: `/api/v1/auth/*`
- RBAC: `/api/v1/user`, `/api/v1/role`, `/api/v1/permission`
- Facebook: `/api/v1/facebook/*`
- Notification: `/api/v1/notification/*`

---

## Authentication Flow

### 1. Đăng Nhập với Firebase

**Flow:**
```
1. Frontend: User đăng nhập bằng Firebase SDK (Email/Google/Facebook/Phone)
2. Firebase: Trả về Firebase ID Token
3. Frontend: Gửi ID Token đến backend
4. Backend: Verify token, tạo/update user trong MongoDB, trả về JWT
5. Frontend: Lưu JWT token để sử dụng cho các request tiếp theo
```

**Endpoint:** `POST /api/v1/auth/login/firebase`

**Request:**
```json
{
  "idToken": "firebase-id-token-here",
  "hwid": "hardware-id-optional"
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Success",
  "status": "success",
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "firebaseUid": "firebase-user-uid",
    "email": "user@example.com",
    "name": "User Name",
    "token": "jwt-token-here",
    "roles": ["role-id-1", "role-id-2"]
  }
}
```

**Lưu ý:**
- `token` trong response là JWT token cần lưu lại
- Sử dụng JWT token này cho tất cả request tiếp theo

### 2. Sử Dụng JWT Token

Tất cả các API (trừ login) yêu cầu header:

```
Authorization: Bearer <jwt-token>
```

**Ví dụ:**
```javascript
fetch('/api/v1/auth/profile', {
  headers: {
    'Authorization': `Bearer ${jwtToken}`,
    'Content-Type': 'application/json'
  }
})
```

### 3. Đăng Xuất

**Endpoint:** `POST /api/v1/auth/logout`

**Request:**
```json
{
  "hwid": "hardware-id-optional"
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Success",
  "status": "success",
  "data": {
    "message": "Logged out successfully"
  }
}
```

**Lưu ý:** Sau khi logout, JWT token sẽ bị vô hiệu hóa. Frontend cần xóa token khỏi storage.

---

## Response Format

### Format Thành Công

Tất cả response thành công đều theo format:

```json
{
  "code": 200,
  "message": "Success",
  "status": "success",
  "data": <response-data>
}
```

**Ví dụ:**
```json
{
  "code": 200,
  "message": "Success",
  "status": "success",
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "name": "Example",
    "email": "example@test.com"
  }
}
```

### Format Lỗi

Tất cả response lỗi đều theo format:

```json
{
  "code": "ERROR_CODE",
  "message": "Error message",
  "status": "error",
  "details": {}
}
```

**Ví dụ:**
```json
{
  "code": "ErrCodeAuth",
  "message": "Token không hợp lệ",
  "status": "error",
  "details": {}
}
```

### HTTP Status Codes

- `200`: Thành công
- `400`: Bad Request (lỗi validation, format)
- `401`: Unauthorized (chưa đăng nhập)
- `403`: Forbidden (không có quyền)
- `404`: Not Found
- `500`: Internal Server Error

---

## Error Handling

### Error Codes Phổ Biến

| Code | Mô Tả |
|------|-------|
| `ErrCodeAuth` | Lỗi xác thực |
| `ErrCodeAuthCredentials` | Sai thông tin đăng nhập |
| `ErrCodeAuthRole` | Không có quyền truy cập |
| `ErrCodeValidationFormat` | Lỗi format input |
| `ErrCodeDatabase` | Lỗi database |
| `ErrCodeInternalServer` | Lỗi server |

### Xử Lý Lỗi trong Frontend

```javascript
async function apiCall(url, options = {}) {
  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Authorization': `Bearer ${getJWTToken()}`,
        'Content-Type': 'application/json',
        ...options.headers
      }
    });
    
    const data = await response.json();
    
    if (data.status === 'error') {
      // Xử lý lỗi
      if (data.code === 'ErrCodeAuth' || response.status === 401) {
        // Token hết hạn hoặc không hợp lệ
        // Redirect đến trang đăng nhập
        redirectToLogin();
        return;
      }
      
      if (response.status === 403) {
        // Không có quyền
        showError('Bạn không có quyền thực hiện thao tác này');
        return;
      }
      
      // Lỗi khác
      showError(data.message);
      return;
    }
    
    return data.data;
  } catch (error) {
    console.error('API Error:', error);
    showError('Có lỗi xảy ra. Vui lòng thử lại.');
  }
}
```

---

## CRUD Operations Pattern

Hệ thống sử dụng pattern CRUD chuẩn cho tất cả các module. Mỗi module có các endpoint sau:

### Cấu Trúc Endpoint

```
/api/v1/{module}/{operation}
```

### Các Operations

#### 1. Create (Tạo)

**Insert One:**
- `POST /api/v1/{module}/insert-one`
- Body: `{ <trường dữ liệu> }`
- Permission: `{Module}.Insert`

**Insert Many:**
- `POST /api/v1/{module}/insert-many`
- Body: `[ { <trường dữ liệu> }, ... ]`
- Permission: `{Module}.Insert`

#### 2. Read (Đọc)

**Find (Tìm kiếm):**
- `GET /api/v1/{module}/find?filter={...}&options={...}`
- Query params:
  - `filter`: Điều kiện lọc (JSON string)
  - `options`: Tùy chọn MongoDB (JSON string: projection, sort, ...)
- Permission: `{Module}.Read`

**Find One:**
- `GET /api/v1/{module}/find-one?filter={...}`
- Permission: `{Module}.Read`

**Find By ID:**
- `GET /api/v1/{module}/find-by-id/:id`
- Permission: `{Module}.Read`

**Find By IDs:**
- `POST /api/v1/{module}/find-by-ids`
- Body: `{ "ids": ["id1", "id2", ...] }`
- Permission: `{Module}.Read`

**Find With Pagination:**
- `GET /api/v1/{module}/find-with-pagination?filter={...}&page=1&limit=20`
- Query params:
  - `filter`: Điều kiện lọc (JSON string)
  - `page`: Số trang (mặc định: 1)
  - `limit`: Số lượng mỗi trang (mặc định: 10)
- Permission: `{Module}.Read`
- Response:
```json
{
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "limit": 20,
    "totalPages": 5
  }
}
```

#### 3. Update (Cập Nhật)

**Update One:**
- `PUT /api/v1/{module}/update-one`
- Body: `{ "filter": {...}, "update": {...} }`
- Permission: `{Module}.Update`

**Update Many:**
- `PUT /api/v1/{module}/update-many`
- Body: `[ { "filter": {...}, "update": {...} }, ... ]`
- Permission: `{Module}.Update`

**Update By ID:**
- `PUT /api/v1/{module}/update-by-id/:id`
- Body: `{ <trường cần cập nhật> }`
- Permission: `{Module}.Update`

**Find One And Update:**
- `PUT /api/v1/{module}/find-one-and-update`
- Body: `{ "filter": {...}, "update": {...} }`
- Permission: `{Module}.Update`

#### 4. Delete (Xóa)

**Delete One:**
- `DELETE /api/v1/{module}/delete-one?filter={...}`
- Permission: `{Module}.Delete`

**Delete Many:**
- `DELETE /api/v1/{module}/delete-many?filter={...}`
- Permission: `{Module}.Delete`

**Delete By ID:**
- `DELETE /api/v1/{module}/delete-by-id/:id`
- Permission: `{Module}.Delete`

**Find One And Delete:**
- `DELETE /api/v1/{module}/find-one-and-dedate?filter={...}`
- Permission: `{Module}.Delete`

#### 5. Other Operations

**Count:**
- `GET /api/v1/{module}/count?filter={...}`
- Response: `{ "data": 100 }`
- Permission: `{Module}.Read`

**Distinct:**
- `GET /api/v1/{module}/distinct?field=name&filter={...}`
- Query params:
  - `field`: Tên trường cần lấy giá trị duy nhất
  - `filter`: Điều kiện lọc (optional)
- Permission: `{Module}.Read`

**Upsert One:**
- `POST /api/v1/{module}/upsert-one`
- Body: `{ "filter": {...}, "update": {...} }`
- Permission: `{Module}.Update`

**Upsert Many:**
- `POST /api/v1/{module}/upsert-many`
- Body: `[ { "filter": {...}, "update": {...} }, ... ]`
- Permission: `{Module}.Update`

**Exists:**
- `GET /api/v1/{module}/exists?filter={...}`
- Response: `{ "data": true/false }`
- Permission: `{Module}.Read`

### Ví Dụ Sử Dụng CRUD

```javascript
// 1. Tạo mới
const newRole = await apiCall('/api/v1/role/insert-one', {
  method: 'POST',
  body: JSON.stringify({
    name: 'Manager',
    code: 'MANAGER',
    organizationId: 'org-id'
  })
});

// 2. Tìm kiếm với filter
const roles = await apiCall(
  `/api/v1/role/find?filter=${encodeURIComponent(JSON.stringify({ organizationId: 'org-id' }))}`
);

// 3. Tìm với phân trang
const paginated = await apiCall(
  `/api/v1/role/find-with-pagination?filter=${encodeURIComponent(JSON.stringify({}))}&page=1&limit=10`
);

// 4. Cập nhật theo ID
const updated = await apiCall(`/api/v1/role/update-by-id/${roleId}`, {
  method: 'PUT',
  body: JSON.stringify({
    name: 'Updated Name'
  })
});

// 5. Xóa theo ID
await apiCall(`/api/v1/role/delete-by-id/${roleId}`, {
  method: 'DELETE'
});
```

---

## Organization Context

### Tổng Quan

Hệ thống hỗ trợ multi-tenant với Organization Context. Mỗi request có thể được filter theo organization.

### Organization Header

Một số API yêu cầu header để xác định organization context:

```
X-Organization-Id: <organization-id>
```

**Lưu ý:** Header này được tự động xử lý bởi `OrganizationContextMiddleware`. Frontend có thể gửi header này nếu cần filter theo organization cụ thể.

### Organization Tree

Organizations được tổ chức theo cấu trúc cây:

```json
{
  "_id": "org-id",
  "name": "Parent Organization",
  "parentId": null,
  "children": [
    {
      "_id": "child-org-id",
      "name": "Child Organization",
      "parentId": "org-id"
    }
  ]
}
```

---

## RBAC và Permissions

### Cấu Trúc RBAC

```
User
  ├── UserRole (nhiều)
  │     └── Role
  │           ├── RolePermission (nhiều)
  │           │     └── Permission
  │           └── Organization
```

### Permission Format

Permission có format: `<Module>.<Action>`

**Ví dụ:**
- `User.Read` - Đọc thông tin user
- `User.Insert` - Tạo user
- `User.Update` - Cập nhật user
- `User.Delete` - Xóa user
- `Role.Read` - Đọc thông tin role
- `Role.Update` - Cập nhật role
- `FbPage.Read` - Đọc thông tin Facebook page
- `Notification.Trigger` - Kích hoạt thông báo

### Permission Scope

Mỗi permission có scope (mức độ quyền):
- `0`: Read (Đọc)
- `1`: Write (Ghi)
- `2`: Delete (Xóa)

### Kiểm Tra Quyền

Tất cả API endpoints (trừ public endpoints) đều yêu cầu permission cụ thể. Nếu user không có permission, sẽ nhận lỗi `403 Forbidden`.

### Lấy Roles của User

**Endpoint:** `GET /api/v1/auth/roles`

**Response:**
```json
{
  "data": [
    {
      "_id": "role-id",
      "name": "Administrator",
      "code": "ADMIN",
      "organizationId": "org-id",
      "permissions": ["permission-id-1", "permission-id-2"]
    }
  ]
}
```

---

## Các Module Chính

### 1. Authentication Module

**Base Path:** `/api/v1/auth`

**Endpoints:**
- `POST /auth/login/firebase` - Đăng nhập với Firebase
- `POST /auth/logout` - Đăng xuất
- `GET /auth/profile` - Lấy profile
- `PUT /auth/profile` - Cập nhật profile
- `GET /auth/roles` - Lấy roles của user

### 2. RBAC Module

**Base Path:** `/api/v1`

**Modules:**
- `/user` - Quản lý người dùng (Read-only)
- `/permission` - Quản lý quyền (Read-only)
- `/role` - Quản lý vai trò (CRUD)
- `/role-permission` - Mapping role-permission (CRUD)
- `/user-role` - Mapping user-role (CRUD)
- `/organization` - Quản lý tổ chức (CRUD)
- `/agent` - Quản lý agent (CRUD)

**Endpoints Đặc Biệt:**
- `PUT /role-permission/update-role` - Cập nhật quyền cho role
- `PUT /user-role/update-user-roles` - Cập nhật roles cho user
- `POST /agent/check-in/:id` - Agent check-in
- `POST /agent/check-out/:id` - Agent check-out

### 3. Facebook Module

**Base Path:** `/api/v1/facebook`

**Modules:**
- `/access-token` - Quản lý access token (CRUD)
- `/page` - Quản lý Facebook pages (CRUD)
- `/post` - Quản lý Facebook posts (CRUD)
- `/conversation` - Quản lý conversations (CRUD)
- `/message` - Quản lý messages (CRUD)
- `/message-item` - Quản lý message items (CRUD)

**Endpoints Đặc Biệt:**
- `GET /facebook/page/find-by-page-id/:id` - Tìm page theo PageID
- `PUT /facebook/page/update-token` - Cập nhật token của page
- `GET /facebook/post/find-by-post-id/:id` - Tìm post theo PostID
- `GET /facebook/conversation/sort-by-api-update` - Lấy conversations sắp xếp theo thời gian cập nhật
- `POST /facebook/message/upsert-messages` - Upsert messages (tự động tách messages)
- `GET /facebook/message-item/find-by-conversation/:conversationId` - Lấy message items theo conversation
- `GET /facebook/message-item/find-by-message-id/:messageId` - Tìm message item theo messageId

### 4. Pancake Module

**Base Path:** `/api/v1/pancake` và `/api/v1/pancake-pos`

**Modules:**
- `/pancake/order` - Quản lý đơn hàng Pancake (CRUD)
- `/pancake-pos/shop` - Quản lý cửa hàng (CRUD)
- `/pancake-pos/warehouse` - Quản lý kho (CRUD)
- `/pancake-pos/product` - Quản lý sản phẩm (CRUD)
- `/pancake-pos/variation` - Quản lý biến thể (CRUD)
- `/pancake-pos/category` - Quản lý danh mục (CRUD)
- `/pancake-pos/order` - Quản lý đơn hàng POS (CRUD)

### 5. Customer Module

**Base Path:** `/api/v1`

**Modules:**
- `/customer` - Customer (deprecated, dùng fb-customer và pc-pos-customer)
- `/fb-customer` - Facebook customer (CRUD)
- `/pc-pos-customer` - Pancake POS customer (CRUD)

### 6. Notification Module

**Base Path:** `/api/v1/notification`

**Modules:**
- `/sender` - Quản lý sender (CRUD)
- `/channel` - Quản lý channel (CRUD)
- `/template` - Quản lý template (CRUD)
- `/routing` - Quản lý routing rules (CRUD)
- `/history` - Lịch sử thông báo (Read-only)

**Endpoints Đặc Biệt:**
- `POST /notification/trigger` - Kích hoạt thông báo
- `GET /notification/track/open/:historyId` - Track mở thông báo (public)
- `GET /notification/track/:historyId/:ctaIndex` - Track click CTA (public)
- `GET /notification/confirm/:historyId` - Xác nhận thông báo (public)

### 7. Admin Module

**Base Path:** `/api/v1/admin`

**Endpoints:**
- `POST /admin/user/block` - Khóa user
- `POST /admin/user/unblock` - Mở khóa user
- `POST /admin/user/role` - Gán role cho user
- `POST /admin/user/set-administrator/:id` - Thiết lập administrator
- `POST /admin/sync-administrator-permissions` - Đồng bộ quyền administrator

### 8. System Module

**Base Path:** `/api/v1/system`

**Endpoints:**
- `GET /system/health` - Kiểm tra sức khỏe hệ thống (public)

### 9. Init Module

**Base Path:** `/api/v1/init`

**Lưu ý:** Các endpoint này chỉ hoạt động khi hệ thống chưa có administrator.

**Endpoints:**
- `GET /init/status` - Kiểm tra trạng thái init
- `POST /init/organization` - Khởi tạo organization
- `POST /init/permissions` - Khởi tạo permissions
- `POST /init/roles` - Khởi tạo roles
- `POST /init/admin-user` - Khởi tạo admin user
- `POST /init/all` - Khởi tạo tất cả (one-click setup)
- `POST /init/set-administrator/:id` - Thiết lập administrator lần đầu

---

## Ví Dụ Sử Dụng

### 1. Đăng Nhập và Lưu Token

```javascript
// Đăng nhập với Firebase
async function loginWithFirebase(firebaseIdToken) {
  const response = await fetch('/api/v1/auth/login/firebase', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      idToken: firebaseIdToken,
      hwid: getHardwareId() // optional
    })
  });
  
  const data = await response.json();
  
  if (data.status === 'success') {
    // Lưu JWT token
    localStorage.setItem('jwt_token', data.data.token);
    localStorage.setItem('user', JSON.stringify(data.data));
    return data.data;
  } else {
    throw new Error(data.message);
  }
}
```

### 2. Lấy Profile

```javascript
async function getProfile() {
  const token = localStorage.getItem('jwt_token');
  
  const response = await fetch('/api/v1/auth/profile', {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  const data = await response.json();
  
  if (data.status === 'success') {
    return data.data;
  } else {
    throw new Error(data.message);
  }
}
```

### 3. Tìm Kiếm với Filter

```javascript
async function searchRoles(organizationId) {
  const token = localStorage.getItem('jwt_token');
  
  const filter = JSON.stringify({ organizationId });
  const url = `/api/v1/role/find?filter=${encodeURIComponent(filter)}`;
  
  const response = await fetch(url, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  const data = await response.json();
  
  if (data.status === 'success') {
    return data.data;
  } else {
    throw new Error(data.message);
  }
}
```

### 4. Tạo Mới với CRUD

```javascript
async function createRole(roleData) {
  const token = localStorage.getItem('jwt_token');
  
  const response = await fetch('/api/v1/role/insert-one', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(roleData)
  });
  
  const data = await response.json();
  
  if (data.status === 'success') {
    return data.data;
  } else {
    throw new Error(data.message);
  }
}
```

### 5. Phân Trang

```javascript
async function getRolesPaginated(page = 1, limit = 20, filter = {}) {
  const token = localStorage.getItem('jwt_token');
  
  const filterStr = encodeURIComponent(JSON.stringify(filter));
  const url = `/api/v1/role/find-with-pagination?filter=${filterStr}&page=${page}&limit=${limit}`;
  
  const response = await fetch(url, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  const data = await response.json();
  
  if (data.status === 'success') {
    return data.data; // { items: [...], total: 100, page: 1, limit: 20, totalPages: 5 }
  } else {
    throw new Error(data.message);
  }
}
```

### 6. Cập Nhật

```javascript
async function updateRole(roleId, updateData) {
  const token = localStorage.getItem('jwt_token');
  
  const response = await fetch(`/api/v1/role/update-by-id/${roleId}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(updateData)
  });
  
  const data = await response.json();
  
  if (data.status === 'success') {
    return data.data;
  } else {
    throw new Error(data.message);
  }
}
```

### 7. Xóa

```javascript
async function deleteRole(roleId) {
  const token = localStorage.getItem('jwt_token');
  
  const response = await fetch(`/api/v1/role/delete-by-id/${roleId}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  const data = await response.json();
  
  if (data.status === 'success') {
    return data.data;
  } else {
    throw new Error(data.message);
  }
}
```

### 8. API Client Wrapper (Recommended)

```javascript
class APIClient {
  constructor(baseURL = '/api/v1') {
    this.baseURL = baseURL;
  }
  
  getToken() {
    return localStorage.getItem('jwt_token');
  }
  
  async request(endpoint, options = {}) {
    const token = this.getToken();
    const url = `${this.baseURL}${endpoint}`;
    
    const config = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...options.headers
      }
    };
    
    if (options.body && typeof options.body === 'object') {
      config.body = JSON.stringify(options.body);
    }
    
    try {
      const response = await fetch(url, config);
      const data = await response.json();
      
      if (data.status === 'error') {
        if (response.status === 401) {
          // Token hết hạn
          localStorage.removeItem('jwt_token');
          window.location.href = '/login';
          throw new Error('Unauthorized');
        }
        
        throw new Error(data.message || 'API Error');
      }
      
      return data.data;
    } catch (error) {
      console.error('API Request Error:', error);
      throw error;
    }
  }
  
  // CRUD helpers
  async find(module, filter = {}, options = {}) {
    const filterStr = encodeURIComponent(JSON.stringify(filter));
    const optionsStr = options ? encodeURIComponent(JSON.stringify(options)) : '';
    let url = `/${module}/find?filter=${filterStr}`;
    if (optionsStr) url += `&options=${optionsStr}`;
    return this.request(url);
  }
  
  async findById(module, id) {
    return this.request(`/${module}/find-by-id/${id}`);
  }
  
  async findWithPagination(module, page = 1, limit = 20, filter = {}) {
    const filterStr = encodeURIComponent(JSON.stringify(filter));
    return this.request(`/${module}/find-with-pagination?filter=${filterStr}&page=${page}&limit=${limit}`);
  }
  
  async insertOne(module, data) {
    return this.request(`/${module}/insert-one`, {
      method: 'POST',
      body: data
    });
  }
  
  async updateById(module, id, data) {
    return this.request(`/${module}/update-by-id/${id}`, {
      method: 'PUT',
      body: data
    });
  }
  
  async deleteById(module, id) {
    return this.request(`/${module}/delete-by-id/${id}`, {
      method: 'DELETE'
    });
  }
}

// Sử dụng
const api = new APIClient();

// Đăng nhập
const user = await api.request('/auth/login/firebase', {
  method: 'POST',
  body: { idToken: firebaseToken }
});

// Lấy roles
const roles = await api.find('role', { organizationId: 'org-id' });

// Phân trang
const paginated = await api.findWithPagination('role', 1, 20, {});
```

---

## Lưu Ý Quan Trọng

### 1. JWT Token

- JWT token được trả về sau khi đăng nhập thành công
- Token cần được lưu và gửi trong header `Authorization: Bearer <token>` cho tất cả request
- Token có thể hết hạn, cần xử lý refresh hoặc redirect đến trang đăng nhập

### 2. Organization Context

- Một số API yêu cầu organization context
- Có thể gửi header `X-Organization-Id` để filter theo organization
- Roles và permissions được gắn với organization

### 3. Permissions

- Mỗi API endpoint yêu cầu permission cụ thể
- Nếu user không có permission, sẽ nhận lỗi `403 Forbidden`
- Frontend nên kiểm tra permissions trước khi hiển thị UI

### 4. Error Handling

- Luôn kiểm tra `data.status` trong response
- Xử lý các lỗi phổ biến: 401 (unauthorized), 403 (forbidden), 400 (bad request)
- Hiển thị thông báo lỗi thân thiện cho user

### 5. Filter Format

- Filter phải là JSON string được encode trong query param
- Sử dụng `encodeURIComponent(JSON.stringify(filter))` khi gửi filter

### 6. Pagination

- Sử dụng `find-with-pagination` cho danh sách lớn
- Response trả về: `{ items: [...], total: 100, page: 1, limit: 20, totalPages: 5 }`

---

## Tài Liệu Tham Khảo

- [Authentication Flow](../02-architecture/authentication.md)
- [RBAC System](../02-architecture/rbac.md)
- [API Documentation](../03-api/)
- [Testing Guide](../06-testing/)

---

**Cập nhật lần cuối:** 2025-01-27

