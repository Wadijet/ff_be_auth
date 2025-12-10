# Go Workspace Configuration

## 📋 Tổng Quan

Dự án sử dụng **Go Workspace** để quản lý nhiều module trong cùng một workspace. Điều này cho phép tách biệt module test khỏi module chính một cách rõ ràng nhưng vẫn dễ quản lý.

## 🏗️ Cấu Trúc

```
ff_be_auth/                    # Root workspace
├── go.work                     # Workspace configuration
├── api/                        # Module chính (meta_commerce)
│   ├── go.mod
│   ├── cmd/                    # Application entry points
│   ├── core/                   # Core business logic
│   └── config/                 # Configuration files
└── api-tests/                  # Module test (ff_be_auth_tests)
    ├── go.mod
    ├── cases/                  # Test cases
    ├── utils/                  # Test utilities
    ├── reports/                # Test reports
    └── run_tests.ps1           # Test runner script
```

## 📝 File go.work

File `go.work` định nghĩa các module trong workspace:

```go
go 1.23

use (
	./api        // Module chính (meta_commerce)
	./api-tests  // Module test (ff_be_auth_tests)
)
```

## 🔧 Sử Dụng

### Chạy từ Root (Khuyến nghị)
```powershell
# Go tự động nhận diện workspace
go test -v ./api-tests/cases/...

# Build module chính
go build ./api/cmd/server

# Chạy server
go run ./api/cmd/server
```

### Chạy từ Module Test
```powershell
cd api-tests
go test -v ./cases/...
```

### Cập nhật Dependencies
```powershell
# Cập nhật dependencies cho module chính
cd api
go mod tidy

# Cập nhật dependencies cho module test
cd api-tests
go mod tidy

# Hoặc từ root với workspace
go work sync
```

## ✅ Lợi Ích

1. **Tách biệt rõ ràng**: Test là module độc lập
2. **Dependencies riêng**: Mỗi module quản lý dependencies của mình
3. **Dễ maintain**: Có thể versioning riêng nếu cần
4. **Đơn giản**: Không cần quản lý nhiều repo

## 📚 Tài Liệu Tham Khảo

- [Go Workspaces Documentation](https://go.dev/doc/tutorial/workspaces)
- [Go Modules Documentation](https://go.dev/doc/modules/managing-dependencies)

