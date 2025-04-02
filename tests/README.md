# Hướng dẫn Test API

## Yêu cầu hệ thống
- Go 1.16 trở lên
- MongoDB 4.4 trở lên
- Git

## Cấu trúc thư mục test
```
tests/
  ├── auth_test.go        # Test các API xác thực (login, register, logout)
  ├── user_test.go        # Test các API người dùng (CRUD, block/unblock)
  ├── permission_test.go  # Test các API phân quyền
  ├── role_test.go        # Test các API vai trò
  ├── agent_test.go       # Test các API agent
  ├── fb_test.go          # Test các API Facebook (pages, conversations, messages, posts)
  ├── order_test.go       # Test các API orders
  └── utils/
      ├── test_utils.go   # Các hàm tiện ích cho test (setup, teardown, helpers)
      └── test_data.go    # Dữ liệu test mẫu
```

## Các bước thực hiện test

### 1. Khởi tạo môi trường test

#### 1.1. Cấu hình môi trường test
Tạo file `.env.test` trong thư mục gốc:
```env
# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB_NAME_AUTH=ff_be_auth_test

# JWT Configuration
JWT_SECRET=your_test_jwt_secret
JWT_EXPIRATION=24h

# Server Configuration
SERVER_PORT=8080
INIT_MODE=true
```

#### 1.2. Khởi tạo database test
```bash
# Kết nối MongoDB
mongosh

# Tạo database mới
use ff_be_auth_test

# Database sẽ được tự động khởi tạo khi chạy server với INIT_MODE=true
```

### 2. Khởi động server
```bash
# Terminal 1
go run main.go
```

### 3. Khởi tạo user admin

#### 3.1. Đăng ký user mới
```bash
# Gọi API đăng ký user
curl -X POST http://localhost:8080/api/v1/users/register \
-H "Content-Type: application/json" \
-d '{
    "email": "admin@example.com",
    "password": "admin123",
    "fullName": "Admin User"
}'
```

#### 3.2. Set quyền Administrator
```bash
# Gọi API set quyền Administrator
curl -X POST http://localhost:8080/api/v1/init/setadmin/{user_id} \
-H "Content-Type: application/json" \
-H "Authorization: Bearer {token}"
```

### 4. Chạy test
```bash
# Terminal 2
# Chạy tất cả test
go test ./tests -v > test_results.log

# Chạy test theo module
go test ./tests/auth_test.go -v
go test ./tests/user_test.go -v
```

## Quy trình test

### 1. Test Authentication (auth_test.go)
- Đăng ký user mới
- Đăng nhập và lấy token
- Test logout
- Test lấy thông tin user hiện tại
- Test đổi mật khẩu
- Test cập nhật thông tin

### 2. Test User Management (user_test.go)
- CRUD user
- Block/Unblock user
- Phân quyền user
- Lấy danh sách user
- Tìm kiếm user

### 3. Test Role & Permission (role_test.go, permission_test.go)
- CRUD roles
- CRUD permissions
- Phân quyền cho role
- Gán role cho user
- Kiểm tra quyền của user

### 4. Test Agent (agent_test.go)
- CRUD agent
- Check-in/Check-out agent
- Quản lý trạng thái agent

### 5. Test Facebook Integration (fb_test.go)
- CRUD Facebook pages
- Cập nhật token Facebook
- CRUD conversations
- CRUD messages
- CRUD posts

### 6. Test Orders (order_test.go)
- CRUD orders
- Quản lý trạng thái đơn hàng
- Tìm kiếm đơn hàng

## Kiểm tra kết quả
- Xem file `test_results.log` để kiểm tra kết quả test
- Format kết quả test:
```
=== RUN   TestUserRegistration
--- PASS: TestUserRegistration (0.52s)
=== RUN   TestUserLogin
--- PASS: TestUserLogin (0.45s)
...
```

## Lưu ý quan trọng
1. Luôn sử dụng database riêng cho test (ff_be_auth_test)
2. Không chạy test trên môi trường production
3. Đảm bảo server đang chạy trước khi thực hiện test
4. Backup dữ liệu test nếu cần
5. Kiểm tra logs của server khi test fail
6. Đảm bảo đã set quyền Administrator cho user test trước khi chạy các test khác

## Troubleshooting
1. Nếu test fail, kiểm tra:
   - Server có đang chạy không
   - Kết nối database có thành công không
   - Token authentication có hợp lệ không
   - Input data có đúng format không
   - Server logs có lỗi gì không
   - User test có quyền Administrator không

2. Các lỗi thường gặp:
   - Connection refused (8080)
   - Authentication failed
   - Invalid input data
   - Missing required fields
   - Database connection error
   - Permission denied

3. Cách khắc phục:
   - Kiểm tra port server
   - Kiểm tra kết nối MongoDB
   - Kiểm tra format request/response
   - Xem logs server
   - Xem logs test
   - Kiểm tra quyền của user test