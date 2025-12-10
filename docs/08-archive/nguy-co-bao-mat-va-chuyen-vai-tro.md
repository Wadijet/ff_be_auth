# NGUY CƠ BẢO MẬT VÀ CHUYỂN VAI TRÒ

Tài liệu này phân tích các nguy cơ bảo mật và khả năng chuyển vai trò trong hệ thống.

---

## 1. CHUYỂN VAI TRÒ

### 1.1. Có thể chuyển vai trò

**Có**, hệ thống hỗ trợ chuyển vai trò qua các endpoint:

#### A. Set Role cho User
**Endpoint:** `POST /api/v1/admin/user/role`
**Permission:** `User.SetRole`

```json
{
  "email": "user@example.com",
  "roleID": "role_id_here"
}
```

**Lưu ý:** Endpoint này chỉ **thay thế** role, không thêm role mới.

#### B. Update User Roles (Thêm/Xóa nhiều roles)
**Endpoint:** `PUT /api/v1/user-role/update-user`
**Permission:** `UserRole.Update`

```json
{
  "userID": "user_id_here",
  "roleIDs": ["role_id_1", "role_id_2"]
}
```

**Cách hoạt động:**
- Xóa tất cả roles cũ của user
- Thêm các roles mới từ `roleIDs`

#### C. CRUD UserRole
**Endpoints:**
- `POST /api/v1/user-role` - Thêm role cho user
- `DELETE /api/v1/user-role/:id` - Xóa role khỏi user
- `GET /api/v1/user-role` - Lấy danh sách user roles

**Permission:** `UserRole.*` (Insert, Update, Delete, Read)

---

### 1.2. Chuyển vai trò Administrator

**Có thể**, nhưng cần cẩn thận:

#### Cách 1: Thêm admin mới (giữ admin cũ)
```bash
POST /api/v1/init/set-administrator/:newUserID
Authorization: Bearer <admin_token>
```

#### Cách 2: Xóa admin cũ, thêm admin mới
```bash
# 1. Xóa admin role khỏi user cũ
DELETE /api/v1/user-role/:oldUserRoleID
Authorization: Bearer <admin_token>

# 2. Thêm admin role cho user mới
POST /api/v1/init/set-administrator/:newUserID
Authorization: Bearer <admin_token>
```

#### Cách 3: Update roles (thay thế tất cả)
```bash
PUT /api/v1/user-role/update-user
Authorization: Bearer <admin_token>
{
  "userID": "new_user_id",
  "roleIDs": ["administrator_role_id"]
}
```

**⚠️ LƯU Ý:** Phải đảm bảo luôn có ít nhất 1 admin trong hệ thống!

---

## 2. NGUY CƠ BẢO MẬT

### 2.1. Phân tích các endpoint init

| Endpoint | Auth | Guard | Nguy cơ | Mô tả |
|----------|------|-------|---------|-------|
| `/init/status` | ❌ | ❌ | ⚠️ Thấp | Chỉ đọc, không thay đổi dữ liệu |
| `/init/organization` | ❌ | ✅ | ✅ **Thấp** | **Tự động disable khi có admin** |
| `/init/permissions` | ❌ | ✅ | ✅ **Thấp** | **Tự động disable khi có admin** |
| `/init/roles` | ❌ | ✅ | ✅ **Thấp** | **Tự động disable khi có admin** |
| `/init/admin-user` | ❌ | ✅ | ✅ **Thấp** | **Tự động disable khi có admin** |
| `/init/all` | ❌ | ✅ | ✅ **Thấp** | **Tự động disable khi có admin** |
| `/init/set-administrator/:id` | ✅* | ❌ | ⚠️ Trung bình | Bỏ qua auth nếu chưa có admin |

*Bỏ qua permission check nếu chưa có admin
✅ Guard = InitGuardMiddleware - tự động disable khi có admin

---

### 2.2. Nguy cơ chiếm quyền

#### ✅ Nguy cơ ĐÃ ĐƯỢC GIẢM THIỂU

**1. Endpoint `/init/admin-user` - ĐÃ ĐƯỢC BẢO VỆ**
```bash
POST /api/v1/init/admin-user
{
  "firebaseUid": "attacker_firebase_uid"
}
```

**Trước đây:**
- 🔴 Attacker có thể tạo admin user bất cứ lúc nào
- 🔴 Không cần authentication
- 🔴 Nguy cơ cao

**Hiện tại:**
- ✅ **Tự động disable khi có admin** (InitGuardMiddleware)
- ✅ Attacker không thể tạo admin sau khi hệ thống đã setup
- ✅ Nguy cơ thấp

**Giải pháp bổ sung (nếu cần):**
- ✅ **Bảo vệ bằng IP whitelist** (chỉ cho phép từ IP cụ thể)
- ✅ **Thêm rate limiting**
- ✅ **Thêm secret key** (phải có key mới gọi được)

---

**2. Endpoint `/init/all` - ĐÃ ĐƯỢC BẢO VỆ**
```bash
POST /api/v1/init/all
```

**Trước đây:**
- 🔴 Attacker có thể setup lại toàn bộ hệ thống bất cứ lúc nào
- 🔴 Nguy cơ cao

**Hiện tại:**
- ✅ **Tự động disable khi có admin** (InitGuardMiddleware)
- ✅ Attacker không thể setup lại sau khi hệ thống đã có admin
- ✅ Nguy cơ thấp

**Giải pháp bổ sung (nếu cần):**
- ✅ **Bảo vệ bằng IP whitelist**
- ✅ **Thêm secret key**

---

**3. Endpoint `/init/set-administrator/:id` bỏ qua auth nếu chưa có admin**
```bash
POST /api/v1/init/set-administrator/:userID
# Không cần token nếu chưa có admin
```

**Nguy cơ:**
- Nếu hệ thống chưa có admin, attacker có thể set admin cho bất kỳ user nào
- Cần user đã login trước (có user ID)

**Giải pháp:**
- ✅ **Đã có:** Chỉ bỏ qua auth khi chưa có admin (cần thiết cho setup lần đầu)
- ✅ **Thêm:** Kiểm tra IP hoặc secret key
- ✅ **Thêm:** Rate limiting

---

#### ⚠️ Nguy cơ TRUNG BÌNH

**1. User đầu tiên tự động trở thành admin**

**Nguy cơ:**
- Nếu attacker là người đầu tiên login, sẽ tự động trở thành admin
- Phụ thuộc vào thời điểm

**Giải pháp:**
- ✅ **Đã có:** Đây là phương án phổ biến, chấp nhận rủi ro
- ✅ **Thêm:** Giới hạn IP cho phép login lần đầu
- ✅ **Thêm:** Set `FIREBASE_ADMIN_UID` để tạo admin trước

---

**2. Endpoint `/init/organization`, `/init/permissions`, `/init/roles`**

**Nguy cơ:**
- Có thể tạo lại các đơn vị cơ bản
- Không trực tiếp tạo admin, nhưng có thể ảnh hưởng đến hệ thống

**Giải pháp:**
- ✅ **Bảo vệ bằng IP whitelist**
- ✅ **Chỉ enable trong development/staging**

---

### 2.3. Nguy cơ mất quyền admin

#### Nguy cơ: Xóa nhầm admin

**Tình huống:**
- Admin xóa nhầm admin role của chính mình
- Hoặc xóa admin role của admin duy nhất

**Hậu quả:**
- Không còn admin nào trong hệ thống
- Không thể quản lý hệ thống

**Giải pháp:**
- ✅ **Thêm validation:** Không cho phép xóa admin role nếu chỉ còn 1 admin
- ✅ **Thêm:** Yêu cầu xác nhận khi xóa admin role
- ✅ **Thêm:** Log tất cả thao tác xóa admin role

---

## 3. BIỆN PHÁP BẢO MẬT ĐỀ XUẤT

### 3.1. Bảo vệ các endpoint init

**✅ ĐÃ TRIỂN KHAI:** InitGuardMiddleware tự động disable các init endpoints sau khi có admin.

**Các endpoint đã được bảo vệ:**
- `/init/organization`
- `/init/permissions`
- `/init/roles`
- `/init/admin-user`
- `/init/all`

**Kết quả:** Nguy cơ đã được giảm thiểu đáng kể. Các endpoint chỉ hoạt động khi chưa có admin.

**Các biện pháp bổ sung (nếu cần):**

#### Phương án 1: IP Whitelist (Tùy chọn)

```go
// Middleware kiểm tra IP
func InitIPWhitelistMiddleware() fiber.Handler {
    allowedIPs := []string{"127.0.0.1", "::1", "192.168.1.100"}
    
    return func(c fiber.Ctx) error {
        clientIP := c.IP()
        if !contains(allowedIPs, clientIP) {
            return c.Status(403).JSON(fiber.Map{
                "error": "Access denied",
            })
        }
        return c.Next()
    }
}
```

**Áp dụng:**
- `/init/admin-user`
- `/init/all`
- `/init/organization`
- `/init/permissions`
- `/init/roles`

---

#### Phương án 2: Secret Key

```go
// Middleware kiểm tra secret key
func InitSecretKeyMiddleware() fiber.Handler {
    secretKey := os.Getenv("INIT_SECRET_KEY")
    
    return func(c fiber.Ctx) error {
        providedKey := c.Get("X-Init-Secret-Key")
        if providedKey != secretKey {
            return c.Status(403).JSON(fiber.Map{
                "error": "Invalid secret key",
            })
        }
        return c.Next()
    }
}
```

**Sử dụng:**
```bash
POST /api/v1/init/admin-user
X-Init-Secret-Key: your_secret_key_here
{
  "firebaseUid": "..."
}
```

---

#### Phương án 3: Chỉ enable trong Development/Staging

```go
func InitEnvironmentMiddleware() fiber.Handler {
    env := os.Getenv("GO_ENV")
    
    return func(c fiber.Ctx) error {
        if env == "production" {
            return c.Status(403).JSON(fiber.Map{
                "error": "Init endpoints disabled in production",
            })
        }
        return c.Next()
    }
}
```

---

### 3.2. Bảo vệ chống xóa admin cuối cùng

```go
// Validation khi xóa admin role
func (s *UserRoleService) DeleteUserRole(ctx context.Context, userRoleID primitive.ObjectID) error {
    // Lấy userRole
    userRole, err := s.FindOneById(ctx, userRoleID)
    if err != nil {
        return err
    }
    
    // Kiểm tra nếu là admin role
    role, err := roleService.FindOneById(ctx, userRole.RoleID)
    if err == nil && role.Name == "Administrator" {
        // Đếm số admin còn lại
        adminCount, err := s.CountAdmins(ctx)
        if err == nil && adminCount <= 1 {
            return errors.New("cannot remove last administrator")
        }
    }
    
    // Xóa userRole
    return s.DeleteById(ctx, userRoleID)
}
```

---

### 3.3. Rate Limiting

```go
// Rate limiting cho init endpoints
func InitRateLimitMiddleware() fiber.Handler {
    limiter := rate.NewLimiter(rate.Every(time.Minute), 5) // 5 requests per minute
    
    return func(c fiber.Ctx) error {
        if !limiter.Allow() {
            return c.Status(429).JSON(fiber.Map{
                "error": "Too many requests",
            })
        }
        return c.Next()
    }
}
```

---

## 4. KHUYẾN NGHỊ

### Development/Staging:
- ✅ **ĐÃ ĐỦ:** InitGuardMiddleware tự động bảo vệ
- ✅ Có thể thêm IP whitelist nếu cần (tùy chọn)

### Production:
- ✅ **ĐÃ ĐỦ:** InitGuardMiddleware tự động disable khi có admin
- ✅ **KHUYẾN NGHỊ:** Set `FIREBASE_ADMIN_UID` để tạo admin trước (an toàn hơn)
- ✅ **KHUYẾN NGHỊ:** Có thể thêm IP whitelist nếu muốn bảo vệ thêm
- 🔴 **BẮT BUỘC:** Validation không cho xóa admin cuối cùng (cần triển khai)
- ✅ **KHUYẾN NGHỊ:** Log tất cả thao tác admin

---

## 5. TÓM TẮT

### Chuyển vai trò:
- ✅ **Có thể** chuyển vai trò qua các endpoint
- ✅ **Có thể** thêm/xóa admin
- ⚠️ **Phải đảm bảo** luôn có ít nhất 1 admin

### Nguy cơ bảo mật:
- ✅ **ĐÃ GIẢM:** Endpoint `/init/admin-user` và `/init/all` tự động disable khi có admin
- ⚠️ **TRUNG BÌNH:** User đầu tiên tự động trở thành admin (phổ biến, chấp nhận được)
- ⚠️ **TRUNG BÌNH:** Có thể xóa nhầm admin cuối cùng (cần validation)

### Biện pháp:
- ✅ IP whitelist
- ✅ Secret key
- ✅ Chỉ enable trong development/staging
- ✅ Validation không cho xóa admin cuối cùng
- ✅ Rate limiting
- ✅ Logging

---

**Cần triển khai các biện pháp bảo mật ngay! 🔒**

