# Hướng Dẫn Chạy Test

## 🚀 Cách Chạy Test Đơn Giản

### Cách 1: Chạy từ root (KHUYẾN NGHỊ)

```powershell
.\api-tests\test.ps1
```

Script này sẽ tự động:
1. Kiểm tra server có đang chạy không
2. Khởi động server nếu chưa chạy
3. Đợi server sẵn sàng (tối đa 60 giây)
4. Chạy toàn bộ test suite
5. Tự động dừng server sau khi test xong
6. Hiển thị kết quả chi tiết

### Cách 2: Chạy test khi server đã sẵn sàng

Nếu server đã chạy sẵn, bạn có thể bỏ qua bước khởi động:

```powershell
.\api-tests\test.ps1 -SkipServer
```

### Cách 3: Chạy từ thư mục api-tests

```powershell
.\api-tests\test.ps1
```

### Cách 4: Quản lý server thủ công

Nếu muốn quản lý server riêng để debug:

```powershell
# Khởi động server
.\api-tests\scripts\manage_server.ps1 start

# Kiểm tra trạng thái
.\api-tests\scripts\manage_server.ps1 status

# Dừng server
.\api-tests\scripts\manage_server.ps1 stop
```

Sau đó chạy test ở terminal khác:
```powershell
.\api-tests\test.ps1 -SkipServer
```

## 📁 Cấu Trúc Test

```
ff_be_auth/
└── api-tests/                  # Module test chính - TẤT CẢ Ở ĐÂY
    ├── test.ps1                # Script chạy test chính (entry point)
    ├── cases/                  # Test cases (Go)
    │   ├── auth_test.go
    │   ├── admin_test.go
    │   ├── health_test.go
    │   └── ...
    ├── utils/                  # Utilities cho test (Go)
    │   ├── http_client.go
    │   └── test_fixtures.go
    ├── scripts/                # Scripts PowerShell cho test
    │   ├── server.ps1          # Module quản lý server
    │   ├── test_runner.ps1     # Module chạy test suite
    │   ├── utils.ps1           # Utilities PowerShell
    │   └── manage_server.ps1   # Script quản lý server độc lập
    ├── reports/                # Báo cáo test
    ├── templates/              # Templates cho báo cáo
    ├── go.mod                  # Module dependencies
    └── README.md               # Tài liệu chi tiết
```

## ⚙️ Yêu Cầu

- MongoDB phải đang chạy
- Go đã được cài đặt
- File config: `api\config\env\development.env` phải tồn tại

## 📊 Kết Quả

Script sẽ hiển thị:
- Tổng số test cases
- Số test passed
- Số test failed
- Pass rate (%)

## 🔧 Troubleshooting

### Server không khởi động được
- Kiểm tra MongoDB có đang chạy không
- Kiểm tra port 8080 có bị chiếm bởi process khác không
- Xem log trong `logs\app.log`

### Test bị lỗi kết nối
- Đảm bảo server đã sẵn sàng trước khi chạy test
- Kiểm tra health endpoint: `http://localhost:8080/api/v1/system/health`

### Script không tìm thấy file config
- Đảm bảo chạy script từ thư mục gốc của project
- Kiểm tra file `api\config\env\development.env` có tồn tại không

## 📚 Tài Liệu Chi Tiết

Xem `api-tests/README.md` để biết thêm chi tiết về cấu trúc và cách sử dụng module test.
