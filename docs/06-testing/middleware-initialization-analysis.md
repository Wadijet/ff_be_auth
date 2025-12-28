# Phân Tích Cách Khởi Tạo Middleware

## 📋 Tổng Quan

Tài liệu này phân tích cách middleware được khởi tạo và đăng ký trong server để đảm bảo chúng được gọi đúng thứ tự.

## 🔄 Luồng Khởi Tạo Server

### 1. Entry Point (`main.go`)

```go
func main() {
    initLogger()           // 1. Khởi tạo logger
    InitGlobal()          // 2. Khởi tạo global variables
    InitRegistry()        // 3. Khởi tạo registry
    InitDefaultData()     // 4. Khởi tạo dữ liệu mặc định
    // ... notification processor
    main_thread()         // 5. Chạy server
}

func main_thread() {
    app := InitFiberApp() // Khởi tạo Fiber app với middleware
    // ... start server
}
```

### 2. Khởi Tạo Fiber App (`init.fiber.go`)

Thứ tự đăng ký middleware trong `InitFiberApp()`:

#### **Thứ Tự Middleware Toàn Cục (Global Middleware)**

1. **Request ID Middleware** (dòng 130-135)
   - Tạo ID duy nhất cho mỗi request
   - Header: `X-Request-ID`
   - ✅ **Được đăng ký đúng**

2. **Debug Middleware** (dòng 138-165)
   - Log tất cả requests và responses
   - ✅ **Được đăng ký đúng**

3. **CORS Middleware** (dòng 182-198)
   - Xử lý cross-origin requests
   - **QUAN TRỌNG**: Phải đặt ở đầu để xử lý preflight requests
   - ✅ **Được đăng ký đúng**

4. **Security Headers Middleware** (dòng 201-209)
   - Thêm các security headers
   - ✅ **Được đăng ký đúng**

5. **Rate Limiting Middleware** (dòng 213-241)
   - Giới hạn số request (nếu enabled)
   - ✅ **Được đăng ký đúng**

6. **Recover Middleware** (dòng 244-273)
   - Xử lý panic và trả về error response
   - ✅ **Được đăng ký đúng**

7. **Logger Middleware** (dòng 276-284)
   - Log requests với format chuẩn
   - ✅ **Được đăng ký đúng**

8. **Setup Routes** (dòng 287)
   - Đăng ký tất cả routes
   - ✅ **Được gọi sau khi đăng ký middleware toàn cục**

## 🔐 Middleware Theo Route

### Trong `routes.go` - `registerCRUDRoutes()`

Thứ tự middleware cho mỗi route CRUD:

```go
router.Post(path, 
    authMiddleware,           // 1. AuthMiddleware - Xác thực và kiểm tra quyền
    orgContextMiddleware,     // 2. OrganizationContextMiddleware - Set organization context
    handler                   // 3. Handler - Xử lý request
)
```

**Thứ tự thực thi:**
1. `AuthMiddleware` - Kiểm tra token, permission, role
2. `OrganizationContextMiddleware` - Set `active_role_id` và `active_organization_id`
3. Handler - Xử lý business logic

### Chi Tiết Middleware Theo Route

#### **AuthMiddleware** (`middleware.auth.go`)

**Chức năng:**
- ✅ Lấy token từ header `Authorization`
- ✅ Validate token và tìm user
- ✅ Kiểm tra user có bị block không
- ✅ Lưu `user_id` và `user` vào context
- ✅ Nếu có `requirePermission`, kiểm tra:
  - Header `X-Active-Role-ID` (BẮT BUỘC)
  - User có role này không
  - User có permission trong role context không
- ✅ Lưu `minScope` vào context

**Logging:**
- Có log chi tiết ở mức Debug và Info
- Log khi middleware được tạo và khi được gọi

#### **OrganizationContextMiddleware** (`middleware.organization_context.go`)

**Chức năng:**
- ✅ Lấy `user_id` từ context (đã được set bởi AuthMiddleware)
- ✅ Lấy `X-Active-Role-ID` từ header
- ✅ Validate user có role này không
- ✅ Từ role, suy ra `organization_id`
- ✅ Lưu `active_role_id` và `active_organization_id` vào context

**Lưu ý:**
- Context làm việc là **ROLE**, không phải organization
- Organization được tự động suy ra từ role

## ✅ Kiểm Tra Thứ Tự Middleware

### Thứ Tự Thực Thi Cho Một Request

```
1. Request ID Middleware
   ↓
2. Debug Middleware
   ↓
3. CORS Middleware
   ↓
4. Security Headers Middleware
   ↓
5. Rate Limiting Middleware (nếu enabled)
   ↓
6. Recover Middleware
   ↓
7. Logger Middleware
   ↓
8. Route Matching
   ↓
9. AuthMiddleware (nếu route yêu cầu)
   ↓
10. OrganizationContextMiddleware (nếu route yêu cầu)
   ↓
11. Handler
   ↓
12. Response Middleware (nếu có)
```

## 🔍 Phát Hiện Vấn Đề

### ✅ Không Có Vấn Đề Nghiêm Trọng

Tất cả middleware được đăng ký đúng thứ tự:

1. ✅ **CORS Middleware** được đặt đúng vị trí (sau Request ID, trước các middleware khác)
2. ✅ **AuthMiddleware** được gọi trước **OrganizationContextMiddleware** (đúng vì cần `user_id`)
3. ✅ **Recover Middleware** được đặt đúng vị trí (sau các middleware khác, trước handler)
4. ✅ **Routes được đăng ký sau middleware toàn cục** (đúng)

### ⚠️ Lưu Ý

1. **Debug Middleware** có thể tạo nhiều log, nên tắt trong production
2. **Rate Limiting** chỉ hoạt động nếu được enable trong config
3. **OrganizationContextMiddleware** phụ thuộc vào `user_id` từ AuthMiddleware, nên phải đặt sau AuthMiddleware

## 📝 Khuyến Nghị

1. ✅ **Giữ nguyên thứ tự hiện tại** - Đã đúng
2. ✅ **CORS Middleware ở đầu** - Đúng vị trí
3. ✅ **AuthMiddleware trước OrganizationContextMiddleware** - Đúng thứ tự
4. ⚠️ **Cân nhắc tắt Debug Middleware trong production**

## 🧪 Cách Kiểm Tra Middleware Có Được Gọi

### 1. Kiểm Tra Logs

Middleware có logging chi tiết:
- `AuthMiddleware`: Log khi được tạo và khi được gọi
- `Debug Middleware`: Log tất cả requests/responses
- `Logger Middleware`: Log với format chuẩn

### 2. Kiểm Tra Headers

- `X-Request-ID`: Được tạo bởi Request ID Middleware
- `X-Content-Type-Options`: Được set bởi Security Headers Middleware
- CORS headers: Được set bởi CORS Middleware

### 3. Kiểm Tra Context

- `user_id`: Được set bởi AuthMiddleware
- `active_role_id`: Được set bởi OrganizationContextMiddleware
- `active_organization_id`: Được set bởi OrganizationContextMiddleware

## 📊 Kết Luận

✅ **Middleware được khởi tạo và đăng ký đúng thứ tự**
✅ **Không có vấn đề về thứ tự thực thi**
✅ **Các middleware phụ thuộc được đặt đúng vị trí**

**Không cần thay đổi gì về cách khởi tạo middleware.**
