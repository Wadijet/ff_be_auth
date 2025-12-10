# Tài liệu ngắn gọn cho FolkForm Auth Backend

Tài liệu này giúp bạn hiểu hệ thống trong vài phút. Toàn bộ nội dung đều viết bằng tiếng Việt, dùng từ ngữ đơn giản để bất kỳ thành viên mới nào cũng có thể đọc và làm theo.

## 1. FolkForm Auth làm gì?
- Quản lý đăng nhập, đăng ký và đổi mật khẩu cho người dùng FolkForm.
- Cấp quyền theo vai trò (RBAC) để bảo vệ các API quản trị.
- Cung cấp nền tảng mở rộng cho các dịch vụ tích hợp (Facebook, đối tác khác).

## 2. Hệ thống vận hành ra sao?
```
Client → Fiber API → Handler → Service → MongoDB
                      ↓
                 Middleware
```
- **Fiber**: máy chủ HTTP.
- **Middleware**: đọc JWT, kiểm tra quyền, log lỗi.
- **Handler**: nhận request, gọi đúng service, trả response chuẩn.
- **Service**: chứa toàn bộ nghiệp vụ.
- **MongoDB**: lưu người dùng, vai trò, token...

## 3. Các thư mục quan trọng
- `cmd/server`: file main, khởi tạo logger, cache, registry rồi bật Fiber.
- `core/api/handler`: chứa handler CRUD và handler đặc thù (auth, admin, facebook...).
- `core/api/services`: logic nghiệp vụ (đăng nhập, quản lý role, webhook...).
- `core/database`: tạo kết nối MongoDB.
- `config`: đọc biến môi trường, cổng chạy, chuỗi kết nối.
- `tests`: bộ test Go cùng template báo cáo.
- `postman/`: collection dùng thử API.

## 4. Một request đi qua những bước nào?
1. Client gọi `/api/v1/...` và gửi JWT ở header.
2. Middleware đọc JWT, gắn thông tin quyền vào context.
3. Handler kiểm tra input, gọi service tương ứng.
4. Service xử lý nghiệp vụ, đọc/ghi MongoDB qua repository.
5. Handler trả JSON chung định dạng trong `handler.base.response.go`.

## 5. Cấu hình và cách chạy local
- Biến mẫu nằm ở `config/env/development.env`.
- Các biến chính:
  - `APP_PORT`: cổng Fiber.
  - `MONGO_URI`: kết nối MongoDB.
  - `JWT_SECRET`: khóa ký token.
  - `CACHE_TTL`: thời gian sống cache.
- Chạy dự án:
```
go mod tidy      # lần đầu
go run ./cmd/server
```
- Nhật ký ghi tại `logs/app.log`, lỗi nền tảng ghi `logs/jobs.log`.

## 6. Nhóm API chính
- **Auth**: đăng ký, đăng nhập, refresh token, đổi mật khẩu.
- **RBAC**: CRUD Role, Permission, RolePermission, UserRole.
- **Admin**: quản lý user, page, post, order.
- **Facebook/Partner**: webhook và đồng bộ dữ liệu bên ngoài.

## 7. Kiểm thử và quan sát
- Unit test nhanh: `go test ./tests/cases -v`.
- Kiểm thử thủ công: dùng các collection trong `postman/`.
- Theo dõi hệ thống sau khi deploy bằng cách đọc file log hoặc `tail -f logs/app.log`.

## 8. Mở rộng và bảo mật
- Quy tắc: handler mỏng, service chứa logic, tận dụng middleware chung.
- Thêm API mới:
  1. Tạo handler (có thể reuse base CRUD).
  2. Khai báo route trong `core/api/router/routes.go`.
  3. Thêm service + model nếu cần.
- Mật khẩu luôn được băm (`core/utility/cipher.go`).
- JWT kiểm hạn, phải giữ kín `JWT_SECRET`.
- Route quản trị luôn cần role phù hợp, cấu hình trong `handler.auth.*`.

---
Nếu cần thêm chi tiết, hãy xem `memory_bank/techContext.md`, phần code tương ứng trong `core/api` hoặc các ghi chú trong thư mục `deploy_notes/`.

