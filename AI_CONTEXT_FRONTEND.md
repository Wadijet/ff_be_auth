# AI Context - Thông Tin Server API cho Frontend Development

## 📋 Tổng Quan Hệ Thống

### Thông Tin Cơ Bản
- **Framework Backend**: Go (Golang) với Fiber v3
- **Database**: MongoDB
- **Base URL**: `http://localhost:8080/api/v1`
- **Authentication**: Firebase Authentication + JWT Token (Bearer Token)
- **Response Format**: JSON

### Mục Đích Hệ Thống
Hệ thống **FolkForm Auth Backend** là một hệ thống quản lý xác thực và phân quyền (RBAC) với các tính năng:
- **Firebase Authentication**: Đăng nhập bằng Firebase (Email, Phone OTP, Google, Facebook)
- Cấp quyền theo vai trò (Role-Based Access Control)
- Quản lý tổ chức (Organization) theo cấu trúc cây
- Tích hợp với Facebook (quản lý pages, posts, conversations, messages)
- Tích hợp với Pancake (quản lý đơn hàng)
- Quản lý Agent (trợ lý tự động) với check-in/check-out

---

## 🔐 Authentication & Authorization

### Cách Xác Thực
Tất cả các API (trừ auth endpoints) yêu cầu header:
```
Authorization: Bearer <token>
```

**Firebase Authentication Flow:**
1. Frontend sử dụng Firebase Client SDK để đăng nhập (Email/Password, Phone OTP, Google, Facebook)
2. Firebase trả về **Firebase ID Token**
3. Frontend gửi Firebase ID Token đến backend endpoint `/auth/login/firebase`
4. Backend verify Firebase ID Token và trả về **JWT Token** của hệ thống
5. Lưu JWT Token để sử dụng cho các request tiếp theo

**Lưu ý:** User được tạo tự động trong MongoDB khi đăng nhập lần đầu với Firebase.

### Permission System
Hệ thống sử dụng RBAC (Role-Based Access Control):
- **Permission**: Quyền cụ thể (ví dụ: `User.Read`, `Role.Update`)
- **Role**: Vai trò chứa nhiều permissions, thuộc về một Organization
- **User**: Người dùng có nhiều roles
- **Scope**: Mức độ quyền (0: Read, 1: Write, 2: Delete)

Format permission: `<Module>.<Action>`
- Module: User, Role, Permission, Agent, FbPage, FbPost, etc.
- Action: Read, Insert, Update, Delete, Block, CheckIn, CheckOut, etc.

---

## 📡 Cấu Trúc Response

### Response Thành Công
```json
{
  "code": 200,
  "message": "Thao tác thành công",
  "data": { /* dữ liệu trả về */ },
  "status": "success"
}
```

### Response Lỗi
```json
{
  "code": "AUTH_001",
  "message": "Thông báo lỗi",
  "details": { /* chi tiết lỗi (nếu có) */ },
  "status": "error"
}
```

### HTTP Status Codes
- `200` - Thành công
- `201` - Tạo mới thành công
- `400` - Yêu cầu không hợp lệ
- `401` - Chưa xác thực
- `403` - Không có quyền truy cập
- `404` - Không tìm thấy
- `409` - Xung đột dữ liệu
- `500` - Lỗi server

---

## 📚 Mô Tả Collections & Tính Năng

### 1. Authentication Module (BẮT BUỘC)

#### User Collection
**Ý nghĩa**: Quản lý thông tin người dùng trong hệ thống
**Tính năng**:
- Đăng ký, đăng nhập, đăng xuất
- Quản lý profile (xem, cập nhật)
- Đổi mật khẩu
- Quản lý tokens (mỗi thiết bị có một token riêng dựa trên HWID)
- Block/Unblock user (chỉ admin)

**Cần thiết**: ⭐⭐⭐⭐⭐ (BẮT BUỘC - Core của hệ thống)

**Model:**
```typescript
interface User {
  id: string;
  firebaseUid: string; // Firebase User ID (unique)
  name: string;
  email?: string; // Optional - có thể đăng nhập bằng phone
  emailVerified: boolean;
  phone?: string; // Optional - có thể đăng nhập bằng email
  phoneVerified: boolean;
  avatarUrl?: string; // URL avatar từ Firebase
  token: string; // JWT token hiện tại
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/user/*` - CRUD operations (Read-only cho user thường)
- `/api/v1/auth/login/firebase` - Đăng nhập bằng Firebase ID Token
- `/api/v1/auth/logout` - Đăng xuất
- `/api/v1/auth/profile` - Xem/Cập nhật profile
- `/api/v1/auth/roles` - Lấy danh sách roles của user

---

### 2. RBAC Module (BẮT BUỘC)

#### Permission Collection
**Ý nghĩa**: Định nghĩa các quyền trong hệ thống (tương ứng với các API endpoints)
**Tính năng**:
- Quản lý danh sách quyền (Read-only, được tạo tự động khi khởi tạo hệ thống)
- Mỗi quyền có format: `Module.Action` (ví dụ: `User.Read`, `Role.Update`)
- Phân loại theo Category và Group

**Cần thiết**: ⭐⭐⭐⭐⭐ (BẮT BUỘC - Core của hệ thống phân quyền)

**Model:**
```typescript
interface Permission {
  id: string;
  name: string; // Format: "Module.Action"
  describe: string;
  category: string; // Auth, Pancake, etc.
  group: string; // User, Role, FbPage, etc.
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/permission/*` - CRUD operations (Read-only)

---

#### Role Collection
**Ý nghĩa**: Định nghĩa các vai trò trong hệ thống, mỗi role thuộc về một Organization
**Tính năng**:
- Tạo, sửa, xóa vai trò
- Mỗi role thuộc về một Organization (bắt buộc)
- Tên role phải unique trong mỗi Organization
- Gán permissions cho role thông qua RolePermission

**Cần thiết**: ⭐⭐⭐⭐⭐ (BẮT BUỘC - Core của hệ thống phân quyền)

**Model:**
```typescript
interface Role {
  id: string;
  name: string;
  describe: string;
  organizationId: string; // BẮT BUỘC - Role thuộc Organization nào
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/role/*` - Full CRUD operations

---

#### RolePermission Collection
**Ý nghĩa**: Liên kết giữa Role và Permission, định nghĩa quyền của từng role
**Tính năng**:
- Gán permissions cho role với scope (0: Read, 1: Write, 2: Delete)
- Cập nhật hàng loạt permissions của một role
- Quản lý quyền chi tiết cho từng role

**Cần thiết**: ⭐⭐⭐⭐⭐ (BẮT BUỘC - Core của hệ thống phân quyền)

**Model:**
```typescript
interface RolePermission {
  id: string;
  roleId: string;
  permissionId: string;
  scope: number; // 0: Read, 1: Write, 2: Delete
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/role-permission/*` - Full CRUD operations
- `/api/v1/role-permission/update-role` - Cập nhật hàng loạt permissions của role

---

#### UserRole Collection
**Ý nghĩa**: Liên kết giữa User và Role, định nghĩa user có những roles nào
**Tính năng**:
- Gán roles cho user
- Một user có thể có nhiều roles
- Quản lý vai trò của từng user

**Cần thiết**: ⭐⭐⭐⭐⭐ (BẮT BUỘC - Core của hệ thống phân quyền)

**Model:**
```typescript
interface UserRole {
  id: string;
  userId: string;
  roleId: string;
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/user-role/*` - Full CRUD operations

---

#### Organization Collection
**Ý nghĩa**: Quản lý cấu trúc tổ chức theo dạng cây (Tập đoàn → Công ty → Phòng ban → Bộ phận → Team)
**Tính năng**:
- Quản lý cấu trúc tổ chức phân cấp
- Hỗ trợ 5 loại: Group, Company, Department, Division, Team
- Mỗi organization có parent (null nếu là root)
- Lưu path và level để truy vấn nhanh
- Roles thuộc về Organization

**Cần thiết**: ⭐⭐⭐⭐ (RẤT QUAN TRỌNG - Nếu hệ thống cần phân quyền theo tổ chức)

**Model:**
```typescript
interface Organization {
  id: string;
  name: string;
  code: string; // Unique code
  type: string; // "group" | "company" | "department" | "division" | "team"
  parentId?: string; // ID của organization cha (null nếu là root)
  path: string; // Đường dẫn cây (ví dụ: "/root_group/company1/dept1")
  level: number; // Cấp độ (0 = root, 1, 2, ...)
  isActive: boolean;
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/organization/*` - Full CRUD operations

---

### 3. Agent Module (TÙY CHỌN - Nếu cần tự động hóa)

#### Agent Collection
**Ý nghĩa**: Quản lý các trợ lý tự động (AI Agent) thực hiện các tác vụ tự động
**Tính năng**:
- Tạo, quản lý agent
- Agent được gán thông tin đăng nhập của user để thực hiện hành động
- Check-in/Check-out để cập nhật trạng thái hoạt động
- Quản lý trạng thái (offline/online) và lệnh điều khiển (stop/play)
- Gán users cho agent
- Lưu config data cho agent

**Cần thiết**: ⭐⭐⭐ (TÙY CHỌN - Chỉ cần nếu hệ thống có tính năng tự động hóa)

**Model:**
```typescript
interface Agent {
  id: string;
  name: string;
  describe: string;
  status: number; // 0: offline, 1: online
  command: number; // 0: stop, 1: play
  assignedUsers: string[]; // Array of user IDs
  configData: Record<string, any>; // Cấu hình agent
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/agent/*` - Full CRUD operations
- `/api/v1/agent/check-in/:id` - Check-in agent (cập nhật trạng thái online)
- `/api/v1/agent/check-out/:id` - Check-out agent (cập nhật trạng thái offline)

**Lưu ý**: Agent cần check-in thường xuyên (mỗi 5 phút) để duy trì trạng thái online. Nếu không check-in sau 5 phút, hệ thống tự động chuyển về offline.

---

### 4. Facebook Integration Module (TÙY CHỌN - Nếu cần tích hợp Facebook)

#### AccessToken Collection
**Ý nghĩa**: Quản lý các access tokens để truy cập vào các hệ thống bên ngoài (Facebook, Pancake, etc.)
**Tính năng**:
- Lưu trữ access tokens cho các hệ thống khác
- Gán tokens cho users
- Quản lý trạng thái active/inactive

**Cần thiết**: ⭐⭐⭐ (TÙY CHỌN - Chỉ cần nếu tích hợp với hệ thống bên ngoài)

**Model:**
```typescript
interface AccessToken {
  id: string;
  name: string; // Unique name
  describe: string;
  system: string; // Hệ thống (Facebook, Pancake, etc.)
  value: string; // Token value
  assignedUsers: string[]; // Array of user IDs
  status: number; // 0: active, 1: inactive
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/access-token/*` - Full CRUD operations

---

#### FbPage Collection
**Ý nghĩa**: Quản lý các Facebook Pages được kết nối với hệ thống
**Tính năng**:
- Lưu thông tin Facebook Pages
- Quản lý Page Access Token
- Đồng bộ dữ liệu từ Pancake (panCakeData)
- Quản lý trạng thái đồng bộ (isSync)

**Cần thiết**: ⭐⭐⭐ (TÙY CHỌN - Chỉ cần nếu tích hợp Facebook)

**Model:**
```typescript
interface FbPage {
  id: string;
  pageName: string;
  pageUsername: string;
  pageId: string; // Facebook Page ID (unique)
  isSync: boolean; // Trạng thái đồng bộ
  accessToken: string;
  pageAccessToken: string; // Page Access Token
  panCakeData: Record<string, any>; // Dữ liệu từ Pancake API
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/facebook/page/*` - Full CRUD operations

---

#### FbPost Collection
**Ý nghĩa**: Quản lý các Facebook Posts từ các Pages
**Tính năng**:
- Lưu thông tin các bài viết trên Facebook
- Liên kết với FbPage
- Đồng bộ dữ liệu từ Pancake

**Cần thiết**: ⭐⭐ (TÙY CHỌN - Chỉ cần nếu cần quản lý Facebook Posts)

**Model:**
```typescript
interface FbPost {
  id: string;
  pageId: string; // Reference to FbPage
  postId: string; // Facebook Post ID (unique)
  panCakeData: Record<string, any>; // Dữ liệu từ Pancake API
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/facebook/post/*` - Full CRUD operations

---

#### FbConversation Collection
**Ý nghĩa**: Quản lý các cuộc trò chuyện (conversations) trên Facebook Messenger
**Tính năng**:
- Lưu thông tin conversations từ Facebook Pages
- Theo dõi thời gian cập nhật từ API (panCakeUpdatedAt)
- Liên kết với FbPage và Customer
- Endpoint đặc biệt để lấy conversations sắp xếp theo thời gian cập nhật API

**Cần thiết**: ⭐⭐⭐ (TÙY CHỌN - Chỉ cần nếu cần quản lý Facebook Conversations)

**Model:**
```typescript
interface FbConversation {
  id: string;
  pageId: string; // Reference to FbPage
  pageUsername: string;
  conversationId: string; // Facebook Conversation ID (unique)
  customerId: string; // Facebook Customer ID
  panCakeData: Record<string, any>; // Dữ liệu từ Pancake API
  panCakeUpdatedAt: number; // Thời gian cập nhật từ Pancake API
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/facebook/conversation/*` - Full CRUD operations
- `/api/v1/facebook/conversation/sort-by-api-update` - Lấy conversations sắp xếp theo thời gian cập nhật API

---

#### FbMessage Collection
**Ý nghĩa**: Quản lý các tin nhắn trong conversations trên Facebook Messenger
**Tính năng**:
- Lưu thông tin messages từ Facebook Conversations
- Liên kết với FbPage và FbConversation
- Đồng bộ dữ liệu từ Pancake

**Cần thiết**: ⭐⭐ (TÙY CHỌN - Chỉ cần nếu cần quản lý chi tiết Facebook Messages)

**Model:**
```typescript
interface FbMessage {
  id: string;
  pageId: string; // Reference to FbPage
  pageUsername: string;
  conversationId: string; // Reference to FbConversation
  customerId: string; // Facebook Customer ID
  panCakeData: Record<string, any>; // Dữ liệu từ Pancake API
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/facebook/message/*` - Full CRUD operations

---

### 5. Pancake Integration Module (TÙY CHỌN - Nếu cần tích hợp Pancake)

#### PcOrder Collection
**Ý nghĩa**: Quản lý đơn hàng từ hệ thống Pancake
**Tính năng**:
- Lưu thông tin đơn hàng từ Pancake
- Đồng bộ dữ liệu đầy đủ từ Pancake API (panCakeData)
- Quản lý trạng thái đơn hàng

**Cần thiết**: ⭐⭐ (TÙY CHỌN - Chỉ cần nếu tích hợp với hệ thống Pancake)

**Model:**
```typescript
interface PcOrder {
  id: string;
  pancakeOrderId: string; // Pancake Order ID (unique)
  status: number; // 0: active, 1: inactive
  panCakeData: Record<string, any>; // Full data from Pancake API
  createdAt: number;
  updatedAt: number;
}
```

**Endpoints:**
- `/api/v1/pancake/order/*` - Full CRUD operations

---

## 📡 API Endpoints Chi Tiết

### 1. System Routes

#### Health Check
```
GET /api/v1/system/health
```
**Không cần authentication**

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-12-09T21:02:38Z",
  "services": {
    "api": "ok",
    "database": "ok"
  }
}
```

---

### 2. Authentication Routes (Không cần token)

#### Đăng Nhập Bằng Firebase
```
POST /api/v1/auth/login/firebase
```

**Request Body:**
```json
{
  "idToken": "firebase-id-token-from-client-sdk",
  "hwid": "hardware-id-unique"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Thao tác thành công",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "firebaseUid": "firebase-user-uid",
    "name": "Tên người dùng",
    "email": "user@example.com",
    "emailVerified": true,
    "phone": "+84123456789",
    "phoneVerified": true,
    "avatarUrl": "https://example.com/avatar.jpg",
    "token": "jwt-token-string",
    "createdAt": 1702147200,
    "updatedAt": 1702147200
  },
  "status": "success"
}
```

**Lưu ý:** 
- Lưu `token` để sử dụng cho các request tiếp theo
- User được tạo tự động trong MongoDB nếu chưa tồn tại
- Nếu là user đầu tiên và chưa có admin, user này sẽ tự động trở thành Administrator

#### Đăng Xuất
```
POST /api/v1/auth/logout
```
**Cần authentication**

**Request Body:**
```json
{
  "hwid": "hardware-id-unique"
}
```

#### Lấy Thông Tin Profile
```
GET /api/v1/auth/profile
```
**Cần authentication**

#### Cập Nhật Profile
```
PUT /api/v1/auth/profile
```
**Cần authentication**

**Request Body:**
```json
{
  "name": "Tên mới"
}
```

**Lưu ý:** Email và phone được quản lý bởi Firebase, không thể thay đổi qua API này.

#### Lấy Danh Sách Roles Của User
```
GET /api/v1/auth/roles
```
**Cần authentication**

---

### 3. CRUD Operations Pattern

Hệ thống sử dụng pattern CRUD thống nhất cho tất cả các collections. Các endpoints có format:

#### Create Operations
- `POST /api/v1/{collection}/insert-one` - Tạo một document
- `POST /api/v1/{collection}/insert-many` - Tạo nhiều documents

#### Read Operations
- `GET /api/v1/{collection}/find` - Tìm tất cả (có filter)
- `GET /api/v1/{collection}/find-one` - Tìm một document
- `GET /api/v1/{collection}/find-by-id/:id` - Tìm theo ID
- `POST /api/v1/{collection}/find-by-ids` - Tìm nhiều documents theo IDs
- `GET /api/v1/{collection}/find-with-pagination` - Tìm có phân trang
- `GET /api/v1/{collection}/count` - Đếm số documents
- `GET /api/v1/{collection}/distinct` - Lấy giá trị distinct
- `GET /api/v1/{collection}/exists` - Kiểm tra document tồn tại

#### Update Operations
- `PUT /api/v1/{collection}/update-one` - Cập nhật một document
- `PUT /api/v1/{collection}/update-many` - Cập nhật nhiều documents
- `PUT /api/v1/{collection}/update-by-id/:id` - Cập nhật theo ID
- `PUT /api/v1/{collection}/find-one-and-update` - Tìm và cập nhật
- `POST /api/v1/{collection}/upsert-one` - Upsert một document
- `POST /api/v1/{collection}/upsert-many` - Upsert nhiều documents

#### Delete Operations
- `DELETE /api/v1/{collection}/delete-one` - Xóa một document
- `DELETE /api/v1/{collection}/delete-many` - Xóa nhiều documents
- `DELETE /api/v1/{collection}/delete-by-id/:id` - Xóa theo ID
- `DELETE /api/v1/{collection}/find-one-and-delete` - Tìm và xóa

#### Query Parameters cho Find Operations

**Filter (query string):**
```
GET /api/v1/user/find?filter={"email":"user@example.com"}
```

**Options (query string):**
```
GET /api/v1/user/find?options={"sort":{"createdAt":-1},"limit":10,"skip":0}
```

**Pagination:**
```
GET /api/v1/user/find-with-pagination?page=1&limit=10&filter={"name":"John"}
```

**Response Pagination:**
```json
{
  "code": 200,
  "message": "Thao tác thành công",
  "data": {
    "page": 1,
    "limit": 10,
    "itemCount": 5,
    "items": [ /* danh sách items */ ]
  },
  "status": "success"
}
```

---

### 4. Admin Routes

#### Block User
```
POST /api/v1/admin/user/block
```
**Permission:** `User.Block`

**Request Body:**
```json
{
  "email": "user@example.com",
  "note": "Lý do khóa tài khoản"
}
```

#### Unblock User
```
POST /api/v1/admin/user/unblock
```
**Permission:** `User.Block`

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

#### Set Role for User
```
POST /api/v1/admin/user/role
```
**Permission:** `User.SetRole`

**Request Body:**
```json
{
  "email": "user@example.com",
  "roleID": "role-id-objectid"
}
```

#### Set Administrator (Khi đã có admin)
```
POST /api/v1/admin/user/set-administrator/:id
```
**Permission:** `Init.SetAdmin`

**Path Parameter:**
- `id`: User ID cần set làm administrator

**Lưu ý:** Endpoint này chỉ dùng khi hệ thống đã có admin. Nếu chưa có admin, sử dụng `/init/set-administrator/:id`.

---

### 5. Init Routes (Chỉ hoạt động khi chưa có admin)

**⚠️ QUAN TRỌNG:** Tất cả init endpoints sẽ **tự động bị tắt** (404 Not Found) sau khi hệ thống đã có admin và server restart.

#### Kiểm Tra Trạng Thái Init
```
GET /api/v1/init/status
```
**Không cần authentication**

**Response:**
```json
{
  "code": 200,
  "data": {
    "organization": {
      "initialized": true,
      "error": ""
    },
    "permissions": {
      "initialized": true,
      "count": 50,
      "error": ""
    },
    "roles": {
      "initialized": true,
      "error": ""
    },
    "adminUsers": {
      "count": 1,
      "hasAdmin": true
    }
  },
  "status": "success"
}
```

#### Khởi Tạo Organization Root
```
POST /api/v1/init/organization
```
**Không cần authentication** (chỉ khi chưa có admin)

#### Khởi Tạo Permissions
```
POST /api/v1/init/permissions
```
**Không cần authentication** (chỉ khi chưa có admin)

#### Khởi Tạo Roles
```
POST /api/v1/init/roles
```
**Không cần authentication** (chỉ khi chưa có admin)

#### Khởi Tạo Admin User từ Firebase UID
```
POST /api/v1/init/admin-user
```
**Không cần authentication** (chỉ khi chưa có admin)

**Request Body:**
```json
{
  "firebaseUid": "firebase-user-uid"
}
```

#### Khởi Tạo Tất Cả (One-click Setup)
```
POST /api/v1/init/all
```
**Không cần authentication** (chỉ khi chưa có admin)

Khởi tạo Organization, Permissions, và Roles trong một lần gọi.

#### Set Administrator (Khi chưa có admin)
```
POST /api/v1/init/set-administrator/:id
```
**Không cần authentication** (chỉ khi chưa có admin)

**Path Parameter:**
- `id`: User ID cần set làm administrator

**Lưu ý:** 
- Endpoint này chỉ hoạt động khi hệ thống chưa có admin
- Nếu đã có admin, sẽ trả về 403 và hướng dẫn dùng `/admin/user/set-administrator/:id`

---

## 🔍 Query Examples

### Tìm User Theo Email
```
GET /api/v1/user/find-one?filter={"email":"user@example.com"}
```

### Tìm Users Có Phân Trang
```
GET /api/v1/user/find-with-pagination?page=1&limit=10&filter={"name":{"$regex":"John"}}
```

### Tìm Users Với Sort
```
GET /api/v1/user/find?filter={}&options={"sort":{"createdAt":-1},"limit":20}
```

### Cập Nhật User
```
PUT /api/v1/user/update-by-id/507f1f77bcf86cd799439011
Content-Type: application/json

{
  "name": "Tên mới"
}
```

### Xóa User
```
DELETE /api/v1/user/delete-by-id/507f1f77bcf86cd799439011
```

---

## ⚠️ Error Handling

### Common Error Codes

**Authentication Errors:**
- `AUTH_001` - Lỗi token (thiếu, không hợp lệ, hết hạn)
- `AUTH_002` - Lỗi thông tin đăng nhập
- `AUTH_003` - Lỗi quyền truy cập

**Validation Errors:**
- `VAL_001` - Lỗi dữ liệu đầu vào
- `VAL_002` - Lỗi định dạng dữ liệu

**Database Errors:**
- `DB` - Lỗi database chung
- `DB_001` - Lỗi kết nối database
- `DB_002` - Lỗi truy vấn database

**Business Logic Errors:**
- `BIZ_001` - Lỗi trạng thái nghiệp vụ
- `BIZ_002` - Lỗi thao tác nghiệp vụ

### Error Response Format
```json
{
  "code": "AUTH_001",
  "message": "Token không hợp lệ",
  "details": null,
  "status": "error"
}
```

---

## 📝 Validation Rules

### Firebase Login
- **idToken**: Required, Firebase ID Token từ Firebase Client SDK
- **Hwid**: Required, Hardware ID duy nhất cho mỗi thiết bị

### Common Validation
- Tất cả các trường có tag `validate:"required"` là bắt buộc
- Firebase ID Token phải hợp lệ và chưa hết hạn
- Hwid phải là string không rỗng

---

## 🎯 Frontend Implementation Guide

### 1. API Client Setup

```typescript
// apiClient.ts
const API_BASE_URL = 'http://localhost:8080/api/v1';

class ApiClient {
  private token: string | null = null;
  private hwid: string;

  constructor() {
    // Tạo hoặc lấy HWID từ localStorage
    this.hwid = this.getOrCreateHWID();
  }

  private getOrCreateHWID(): string {
    let hwid = localStorage.getItem('hwid');
    if (!hwid) {
      // Tạo HWID duy nhất (có thể dùng thư viện như device-uuid)
      hwid = this.generateHWID();
      localStorage.setItem('hwid', hwid);
    }
    return hwid;
  }

  private generateHWID(): string {
    // Sử dụng device fingerprint hoặc thư viện device-uuid
    // Ví dụ đơn giản:
    return `hwid_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  setToken(token: string) {
    this.token = token;
    localStorage.setItem('auth_token', token);
  }

  getToken(): string | null {
    return this.token || localStorage.getItem('auth_token');
  }

  getHWID(): string {
    return this.hwid;
  }

  clearToken() {
    this.token = null;
    localStorage.removeItem('auth_token');
  }

  async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const token = this.getToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    const data = await response.json();

    if (!response.ok || data.status === 'error') {
      throw new ApiError(data.message, data.code, response.status);
    }

    return data;
  }

  // CRUD Methods
  async find<T>(collection: string, filter?: any, options?: any): Promise<T[]> {
    const params = new URLSearchParams();
    if (filter) params.append('filter', JSON.stringify(filter));
    if (options) params.append('options', JSON.stringify(options));
    
    const response = await this.request<{ data: T[] }>(
      `/${collection}/find?${params.toString()}`
    );
    return response.data;
  }

  async findOne<T>(collection: string, filter?: any): Promise<T> {
    const params = new URLSearchParams();
    if (filter) params.append('filter', JSON.stringify(filter));
    
    const response = await this.request<{ data: T }>(
      `/${collection}/find-one?${params.toString()}`
    );
    return response.data;
  }

  async findById<T>(collection: string, id: string): Promise<T> {
    const response = await this.request<{ data: T }>(
      `/${collection}/find-by-id/${id}`
    );
    return response.data;
  }

  async insertOne<T>(collection: string, data: any): Promise<T> {
    const response = await this.request<{ data: T }>(
      `/${collection}/insert-one`,
      {
        method: 'POST',
        body: JSON.stringify(data),
      }
    );
    return response.data;
  }

  async updateById<T>(
    collection: string,
    id: string,
    data: any
  ): Promise<T> {
    const response = await this.request<{ data: T }>(
      `/${collection}/update-by-id/${id}`,
      {
        method: 'PUT',
        body: JSON.stringify(data),
      }
    );
    return response.data;
  }

  async deleteById(collection: string, id: string): Promise<void> {
    await this.request(`/${collection}/delete-by-id/${id}`, {
      method: 'DELETE',
    });
  }

  async findWithPagination<T>(
    collection: string,
    page: number = 1,
    limit: number = 10,
    filter?: any
  ): Promise<PaginatedResponse<T>> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    if (filter) params.append('filter', JSON.stringify(filter));

    const response = await this.request<{ data: PaginatedResponse<T> }>(
      `/${collection}/find-with-pagination?${params.toString()}`
    );
    return response.data;
  }
}

// Types
interface ApiResponse<T> {
  code: number | string;
  message: string;
  data: T;
  status: 'success' | 'error';
}

interface PaginatedResponse<T> {
  page: number;
  limit: number;
  itemCount: number;
  items: T[];
}

class ApiError extends Error {
  constructor(
    message: string,
    public code: string,
    public statusCode: number
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export const apiClient = new ApiClient();
```

### 2. Auth Service

```typescript
// authService.ts
import { apiClient } from './apiClient';

export interface FirebaseLoginInput {
  idToken: string; // Firebase ID Token từ Firebase Client SDK
  hwid: string;
}

export interface User {
  id: string;
  firebaseUid: string;
  name: string;
  email?: string;
  emailVerified: boolean;
  phone?: string;
  phoneVerified: boolean;
  avatarUrl?: string;
  token?: string;
  createdAt: number;
  updatedAt: number;
}

class AuthService {
  /**
   * Đăng nhập bằng Firebase ID Token
   * @param idToken Firebase ID Token từ Firebase Client SDK
   */
  async loginWithFirebase(idToken: string): Promise<User> {
    const hwid = apiClient.getHWID();
    const response = await apiClient.request<{ data: User }>(
      '/auth/login/firebase',
      {
        method: 'POST',
        body: JSON.stringify({
          idToken,
          hwid,
        }),
      }
    );

    if (response.data.token) {
      apiClient.setToken(response.data.token);
    }

    return response.data;
  }

  async logout(): Promise<void> {
    const hwid = apiClient.getHWID();
    await apiClient.request('/auth/logout', {
      method: 'POST',
      body: JSON.stringify({ hwid }),
    });
    apiClient.clearToken();
  }

  async getProfile(): Promise<User> {
    const response = await apiClient.request<{ data: User }>(
      '/auth/profile'
    );
    return response.data;
  }

  async updateProfile(name: string): Promise<User> {
    const response = await apiClient.request<{ data: User }>(
      '/auth/profile',
      {
        method: 'PUT',
        body: JSON.stringify({ name }),
      }
    );
    return response.data;
  }

  /**
   * Lưu ý: Email và phone được quản lý bởi Firebase
   * Để thay đổi email/phone, sử dụng Firebase Client SDK
   */

  async getUserRoles() {
    const response = await apiClient.request<{ data: any[] }>(
      '/auth/roles'
    );
    return response.data;
  }
}

export const authService = new AuthService();
```

### 3. User Management Service

```typescript
// userService.ts
import { apiClient } from './apiClient';

export interface User {
  id: string;
  name: string;
  email: string;
  createdAt: number;
  updatedAt: number;
}

class UserService {
  async findAll(filter?: any): Promise<User[]> {
    return apiClient.find<User>('user', filter);
  }

  async findOne(filter: any): Promise<User> {
    return apiClient.findOne<User>('user', filter);
  }

  async findById(id: string): Promise<User> {
    return apiClient.findById<User>('user', id);
  }

  async findWithPagination(
    page: number = 1,
    limit: number = 10,
    filter?: any
  ) {
    return apiClient.findWithPagination<User>('user', page, limit, filter);
  }
}

export const userService = new UserService();
```

### 4. Error Handling

```typescript
// errorHandler.ts
import { ApiError } from './apiClient';

export function handleApiError(error: unknown): string {
  if (error instanceof ApiError) {
    switch (error.code) {
      case 'AUTH_001':
        return 'Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.';
      case 'AUTH_002':
        return 'Thông tin đăng nhập không chính xác.';
      case 'AUTH_003':
        return 'Bạn không có quyền thực hiện thao tác này.';
      case 'VAL_001':
        return 'Dữ liệu không hợp lệ. Vui lòng kiểm tra lại.';
      case 'DB_002':
        return 'Không tìm thấy dữ liệu.';
      default:
        return error.message || 'Đã xảy ra lỗi. Vui lòng thử lại.';
    }
  }

  if (error instanceof Error) {
    return error.message;
  }

  return 'Đã xảy ra lỗi không xác định.';
}
```

---

## 🔑 Important Notes

1. **Firebase Authentication**: 
   - Sử dụng Firebase Client SDK để đăng nhập (Email, Phone OTP, Google, Facebook)
   - Lấy Firebase ID Token từ Firebase sau khi đăng nhập thành công
   - Gửi Firebase ID Token đến `/auth/login/firebase` để nhận JWT token của hệ thống
   - Lưu JWT token vào localStorage hoặc state management
   - Gửi JWT token trong header `Authorization: Bearer <token>` cho mọi request (trừ auth endpoints)

2. **HWID (Hardware ID)**:
   - Cần tạo và lưu trữ một hardware ID duy nhất cho mỗi thiết bị
   - Sử dụng khi login và logout
   - Có thể sử dụng thư viện như `device-uuid` hoặc tạo từ browser fingerprint

3. **Pagination**:
   - Sử dụng `find-with-pagination` cho danh sách lớn
   - Response có format: `{ page, limit, itemCount, items }`

4. **Filter & Options**:
   - Filter và options được truyền qua query string dưới dạng JSON
   - Sử dụng MongoDB query syntax cho filter
   - Options hỗ trợ: `sort`, `limit`, `skip`, `projection`

5. **Error Handling**:
   - Luôn kiểm tra `status === "error"` trong response
   - Hiển thị message từ response cho user
   - Xử lý 401 để redirect về login page

6. **Permissions**:
   - Mỗi endpoint yêu cầu permission cụ thể
   - Format: `<Module>.<Action>`
   - Nếu không có permission, sẽ nhận 403 Forbidden

7. **Organization & Roles**:
   - Roles phải thuộc về một Organization
   - Tên role phải unique trong mỗi Organization
   - Khi tạo role, bắt buộc phải có `organizationId`

8. **Agent Check-in**:
   - Agent cần check-in mỗi 5 phút để duy trì trạng thái online
   - Nếu không check-in sau 5 phút, hệ thống tự động chuyển về offline

---

## 📊 Tóm Tắt Collections Theo Mức Độ Cần Thiết

### ⭐⭐⭐⭐⭐ BẮT BUỘC (Core System)
- **User** - Quản lý người dùng
- **Permission** - Định nghĩa quyền
- **Role** - Định nghĩa vai trò
- **RolePermission** - Liên kết Role-Permission
- **UserRole** - Liên kết User-Role

### ⭐⭐⭐⭐ RẤT QUAN TRỌNG (Nếu cần phân quyền theo tổ chức)
- **Organization** - Cấu trúc tổ chức

### ⭐⭐⭐ TÙY CHỌN (Tích hợp và tự động hóa)
- **Agent** - Trợ lý tự động
- **AccessToken** - Quản lý tokens
- **FbPage** - Facebook Pages
- **FbConversation** - Facebook Conversations

### ⭐⭐ TÙY CHỌN (Chi tiết)
- **FbPost** - Facebook Posts
- **FbMessage** - Facebook Messages
- **PcOrder** - Pancake Orders

---

## 📚 Additional Resources

- Base URL: `http://localhost:8080/api/v1`
- Health Check: `GET /api/v1/system/health`
- All endpoints require authentication except:
  - `/auth/login/firebase`
  - `/init/status` (chỉ khi chưa có admin)
  - `/init/*` (chỉ khi chưa có admin, sẽ bị tắt sau khi có admin)
  - `/system/health`

---

**Tài liệu này cung cấp đầy đủ thông tin về ý nghĩa, tính năng và mức độ cần thiết của từng collection để phát triển frontend tích hợp với API server này.**
