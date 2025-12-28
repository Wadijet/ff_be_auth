# Kịch Bản Test Organization Ownership

Tài liệu này mô tả các kịch bản test cho tính năng phân quyền dữ liệu theo organization.

## 📋 Tổng Quan

Sau khi triển khai organization ownership, hệ thống sẽ:
- Tự động gán `organizationId` khi tạo dữ liệu mới
- Tự động filter dữ liệu theo quyền của user (bao gồm cả parent organizations)
- Validate quyền truy cập khi update/delete

## 🔑 Headers Cần Thiết

### 1. Authorization Header
```
Authorization: Bearer <JWT_TOKEN>
```

### 2. X-Active-Role-ID Header (Mới)
```
X-Active-Role-ID: <ROLE_ID>
```
Header này xác định organization context mà user đang làm việc. Nếu không có, hệ thống sẽ tự động chọn role đầu tiên của user.

## 📝 Test Scenarios

### Scenario 1: Lấy Danh Sách Roles Của User

**Endpoint:** `GET /api/v1/auth/roles`

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Expected Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    {
      "roleId": "507f1f77bcf86cd799439011",
      "roleName": "Manager",
      "organizationId": "507f1f77bcf86cd799439012",
      "organizationName": "Company A",
      "organizationCode": "COMPANY_A",
      "organizationType": 2,
      "organizationLevel": 1
    }
  ],
  "status": "success"
}
```

**Test Steps:**
1. Login để lấy JWT token
2. Gọi API `/auth/roles`
3. Verify response có đầy đủ thông tin role và organization
4. Lưu `roleId` để dùng cho các test tiếp theo

---

### Scenario 2: Tạo Dữ Liệu Với Organization Context

**Endpoint:** `POST /api/v1/fb-customer/insert-one`

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
X-Active-Role-ID: <ROLE_ID>
```

**Request Body:**
```json
{
  "customerId": "test_customer_123",
  "name": "Test Customer",
  "email": "test@example.com"
}
```

**Expected Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "customerId": "test_customer_123",
    "name": "Test Customer",
    "email": "test@example.com",
    "organizationId": "507f1f77bcf86cd799439012",  // ✅ Tự động gán
    "createdAt": 1234567890,
    "updatedAt": 1234567890
  },
  "status": "success"
}
```

**Test Steps:**
1. Set `X-Active-Role-ID` header với role ID từ Scenario 1
2. Tạo customer mới
3. Verify `organizationId` đã được tự động gán đúng với organization của role

---

### Scenario 3: Filter Dữ Liệu Theo Organization

**Endpoint:** `GET /api/v1/fb-customer/find`

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
X-Active-Role-ID: <ROLE_ID>
```

**Expected Behavior:**
- Chỉ trả về customers thuộc organization của role
- Bao gồm cả customers của parent organizations (inverse lookup)

**Test Steps:**
1. Tạo customer ở organization cha (Company)
2. Tạo customer ở organization con (Department)
3. Set `X-Active-Role-ID` với role ở organization con
4. Gọi API `/fb-customer/find`
5. Verify response chỉ chứa customers mà user có quyền xem

---

### Scenario 4: Scope = 0 (Self) - Chỉ Xem Dữ Liệu Của Organization Mình

**Setup:**
- Tạo role với permission `FbCustomer.Read` và `Scope = 0`
- Role thuộc Department A

**Test Steps:**
1. Tạo customer ở Department A
2. Tạo customer ở Department B (cùng Company)
3. Set `X-Active-Role-ID` với role ở Department A
4. Gọi API `/fb-customer/find`
5. **Expected:** Chỉ thấy customer ở Department A

---

### Scenario 5: Scope = 1 (Children) - Xem Dữ Liệu Của Organization Mình Và Con

**Setup:**
- Tạo role với permission `FbCustomer.Read` và `Scope = 1`
- Role thuộc Company A
- Company A có Department B và Department C

**Test Steps:**
1. Tạo customer ở Company A
2. Tạo customer ở Department B
3. Tạo customer ở Department C
4. Set `X-Active-Role-ID` với role ở Company A
5. Gọi API `/fb-customer/find`
6. **Expected:** Thấy tất cả customers ở Company A, Department B, và Department C

---

### Scenario 6: Inverse Parent Lookup - Xem Dữ Liệu Cấp Trên

**Setup:**
- User có role ở Department B
- Department B thuộc Company A

**Test Steps:**
1. Tạo customer ở Company A (organization cha)
2. Set `X-Active-Role-ID` với role ở Department B
3. Gọi API `/fb-customer/find`
4. **Expected:** Thấy customer ở Company A (tự động thông qua inverse parent lookup)

---

### Scenario 7: Update/Delete Với Organization Filter

**Endpoint:** `PUT /api/v1/fb-customer/update-by-id/:id`

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
X-Active-Role-ID: <ROLE_ID>
```

**Test Steps:**
1. Tạo customer với role A
2. Set `X-Active-Role-ID` với role B (khác organization)
3. Thử update customer từ role A
4. **Expected:** Lỗi 403 Forbidden - Không có quyền truy cập

---

### Scenario 8: Không Cho Phép Update organizationId

**Endpoint:** `PUT /api/v1/fb-customer/update-by-id/:id`

**Request Body:**
```json
{
  "name": "Updated Name",
  "organizationId": "507f1f77bcf86cd799439999"  // ❌ Không được phép
}
```

**Expected Behavior:**
- Field `organizationId` sẽ bị bỏ qua trong update
- Chỉ update các field khác

---

### Scenario 9: Multi-Client Support

**Test Steps:**
1. Client 1: Set `X-Active-Role-ID: ROLE_A`, tạo customer A
2. Client 2: Set `X-Active-Role-ID: ROLE_B`, tạo customer B
3. Verify mỗi client chỉ thấy dữ liệu của organization tương ứng

---

### Scenario 10: Collections Không Có OrganizationID

**Test Collections:**
- Users
- Permissions
- Organizations
- UserRoles
- RolePermissions

**Expected Behavior:**
- CRUD operations hoạt động bình thường
- Không có filter theo organizationId
- Không tự động gán organizationId

---

## 🧪 Test Cases Chi Tiết

### Test Case 1: Tạo Customer Với Organization Context

```bash
# 1. Login
POST /api/v1/auth/login/firebase
{
  "idToken": "<FIREBASE_TOKEN>",
  "hwid": "test_device"
}

# 2. Lấy roles
GET /api/v1/auth/roles
Authorization: Bearer <JWT_TOKEN>

# 3. Tạo customer với active role
POST /api/v1/fb-customer/insert-one
Authorization: Bearer <JWT_TOKEN>
X-Active-Role-ID: <ROLE_ID>
{
  "customerId": "test_001",
  "name": "Test Customer",
  "email": "test@example.com"
}

# Verify: Response có organizationId = organization của role
```

### Test Case 2: Filter Customers Theo Organization

```bash
# 1. Tạo customer ở organization A
POST /api/v1/fb-customer/insert-one
X-Active-Role-ID: <ROLE_ORG_A>
{
  "customerId": "customer_a",
  "name": "Customer A"
}

# 2. Tạo customer ở organization B
POST /api/v1/fb-customer/insert-one
X-Active-Role-ID: <ROLE_ORG_B>
{
  "customerId": "customer_b",
  "name": "Customer B"
}

# 3. Query với role A
GET /api/v1/fb-customer/find
X-Active-Role-ID: <ROLE_ORG_A>

# Verify: Chỉ thấy customer_a, không thấy customer_b
```

### Test Case 3: Inverse Parent Lookup

```bash
# 1. Tạo customer ở Company (parent)
POST /api/v1/fb-customer/insert-one
X-Active-Role-ID: <ROLE_COMPANY>
{
  "customerId": "parent_customer",
  "name": "Parent Customer"
}

# 2. Query với role Department (child)
GET /api/v1/fb-customer/find
X-Active-Role-ID: <ROLE_DEPARTMENT>

# Verify: Thấy parent_customer (tự động thông qua inverse lookup)
```

---

## ✅ Checklist Test

- [ ] Test lấy danh sách roles với thông tin organization
- [ ] Test tạo dữ liệu với organization context (tự động gán organizationId)
- [ ] Test filter dữ liệu theo organization
- [ ] Test scope = 0 (chỉ xem dữ liệu của organization mình)
- [ ] Test scope = 1 (xem dữ liệu của organization mình và con)
- [ ] Test inverse parent lookup (xem dữ liệu cấp trên)
- [ ] Test update/delete với organization filter
- [ ] Test không cho phép update organizationId
- [ ] Test multi-client support
- [ ] Test collections không có organizationId hoạt động bình thường

---

## 📚 Tài Liệu Liên Quan

- [Organization Ownership Analysis](../../02-architecture/organization-ownership-analysis.md)
- [Implementation Plan](../../02-architecture/organization-ownership-implementation-plan.md)
- [Collections Without Organization](../../02-architecture/collections-without-organization.md)

