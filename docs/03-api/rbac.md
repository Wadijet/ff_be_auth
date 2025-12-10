# RBAC APIs

Tài liệu về các API endpoints quản lý Role-Based Access Control (Role, Permission, RolePermission, UserRole).

## 📋 Tổng Quan

Hệ thống RBAC bao gồm:
- **Permission**: Quyền cụ thể (ví dụ: `User.Read`, `Role.Update`)
- **Role**: Vai trò chứa nhiều permissions
- **RolePermission**: Mapping giữa Role và Permission
- **UserRole**: Mapping giữa User và Role

## 🔐 Permission APIs

Tất cả endpoints nằm dưới `/api/v1/permission/` (Read-only).

### Endpoints
- `GET /api/v1/permission/find` - Tìm tất cả permissions
- `GET /api/v1/permission/find-one` - Tìm một permission
- `GET /api/v1/permission/find-by-id/:id` - Tìm permission theo ID
- `GET /api/v1/permission/find-by-ids` - Tìm nhiều permissions theo IDs
- `GET /api/v1/permission/find-with-pagination` - Tìm với phân trang
- `GET /api/v1/permission/count` - Đếm permissions

**Authentication:** Cần (Permission: `Permission.Read`)

## 🔐 Role APIs

Tất cả endpoints nằm dưới `/api/v1/role/` (Full CRUD).

### Endpoints
- `POST /api/v1/role/insert-one` - Tạo role mới (Permission: `Role.Insert`)
- `GET /api/v1/role/find` - Tìm tất cả roles (Permission: `Role.Read`)
- `GET /api/v1/role/find-by-id/:id` - Tìm role theo ID (Permission: `Role.Read`)
- `PUT /api/v1/role/update-by-id/:id` - Cập nhật role (Permission: `Role.Update`)
- `DELETE /api/v1/role/delete-by-id/:id` - Xóa role (Permission: `Role.Delete`)

### Ví Dụ: Tạo Role

**Request:**
```json
POST /api/v1/role/insert-one
{
  "name": "Manager",
  "code": "MANAGER",
  "organizationId": "507f1f77bcf86cd799439012",
  "description": "Manager role"
}
```

**Response:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "name": "Manager",
    "code": "MANAGER",
    "organizationId": "507f1f77bcf86cd799439012",
    "description": "Manager role"
  }
}
```

## 🔐 RolePermission APIs

Tất cả endpoints nằm dưới `/api/v1/role-permission/` (Full CRUD).

### Endpoints CRUD
- `POST /api/v1/role-permission/insert-one` - Tạo mapping (Permission: `RolePermission.Insert`)
- `GET /api/v1/role-permission/find` - Tìm mappings (Permission: `RolePermission.Read`)
- `PUT /api/v1/role-permission/update-by-id/:id` - Cập nhật mapping (Permission: `RolePermission.Update`)
- `DELETE /api/v1/role-permission/delete-by-id/:id` - Xóa mapping (Permission: `RolePermission.Delete`)

### Endpoint Đặc Biệt: Update Role Permissions

Cập nhật tất cả permissions của một role.

**Endpoint:** `PUT /api/v1/role-permission/update-role`

**Authentication:** Cần (Permission: `RolePermission.Update`)

**Request Body:**
```json
{
  "roleId": "507f1f77bcf86cd799439011",
  "permissionIds": ["507f1f77bcf86cd799439012", "507f1f77bcf86cd799439013"]
}
```

**Response:**
```json
{
  "data": {
    "message": "Role permissions updated successfully"
  }
}
```

## 🔐 UserRole APIs

Tất cả endpoints nằm dưới `/api/v1/user-role/` (Full CRUD).

### Endpoints
- `POST /api/v1/user-role/insert-one` - Gán role cho user (Permission: `UserRole.Insert`)
- `GET /api/v1/user-role/find` - Tìm mappings (Permission: `UserRole.Read`)
- `PUT /api/v1/user-role/update-by-id/:id` - Cập nhật mapping (Permission: `UserRole.Update`)
- `DELETE /api/v1/user-role/delete-by-id/:id` - Xóa mapping (Permission: `UserRole.Delete`)

### Ví Dụ: Gán Role cho User

**Request:**
```json
POST /api/v1/user-role/insert-one
{
  "userId": "507f1f77bcf86cd799439011",
  "roleId": "507f1f77bcf86cd799439012"
}
```

**Response:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439013",
    "userId": "507f1f77bcf86cd799439011",
    "roleId": "507f1f77bcf86cd799439012"
  }
}
```

## 📝 Lưu Ý

- Tất cả endpoints đều yêu cầu authentication
- Mỗi endpoint yêu cầu permission tương ứng
- Permission collection là read-only (chỉ có thể đọc)
- Role, RolePermission, UserRole có full CRUD operations

## 📚 Tài Liệu Liên Quan

- [RBAC System](../02-architecture/rbac.md)
- [Admin APIs](admin.md)
- [User Management APIs](user-management.md)

