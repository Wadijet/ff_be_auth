# FolkForm Auth Backend

Hệ thống xác thực và quản lý quyền (RBAC) cho nền tảng FolkForm, được xây dựng bằng Go với Fiber framework.

## 📋 Tổng Quan

FolkForm Auth Backend là một hệ thống backend cung cấp các tính năng:

- 🔐 **Firebase Authentication**: Đăng nhập đa phương thức (Email/Password, Google, Facebook, Phone OTP)
- 👥 **Quản lý Người Dùng**: Tự động tạo user từ Firebase, quản lý profile
- 🔑 **RBAC (Role-Based Access Control)**: Hệ thống phân quyền theo vai trò và tổ chức
- 🏢 **Quản lý Tổ chức**: Cấu trúc tổ chức theo cây (Organization Tree)
- 📱 **Tích hợp Facebook**: Quản lý pages, posts, conversations, messages
- 🛒 **Tích hợp Pancake**: Quản lý đơn hàng
- 🤖 **Quản lý Agent**: Hệ thống trợ lý tự động với check-in/check-out

## 🚀 Bắt Đầu Nhanh

### Yêu Cầu Hệ Thống

- Go 1.23+ 
- MongoDB 4.4+
- Firebase Project (cho Authentication)

### Cài Đặt

1. **Clone repository:**
```bash
git clone <repository-url>
cd ff_be_auth
```

2. **Cài đặt dependencies:**
```bash
cd api
go mod download
```

3. **Cấu hình môi trường:**
```bash
# Copy file cấu hình mẫu
cp config/env/development.env config/env/development.env.local

# Chỉnh sửa các biến môi trường cần thiết
# - MongoDB connection string
# - Firebase credentials
# - JWT secret
```

4. **Chạy server:**
```bash
go run cmd/server/main.go
```

Server sẽ chạy tại `http://localhost:8080`

### Kiểm Tra Sức Khỏe Hệ Thống

```bash
curl http://localhost:8080/api/v1/system/health
```

## 📁 Cấu Trúc Dự Án

```
ff_be_auth/
├── api/                          # Module chính
│   ├── cmd/
│   │   └── server/              # Entry point của ứng dụng
│   ├── core/
│   │   ├── api/                 # API layer
│   │   │   ├── handler/         # HTTP handlers
│   │   │   ├── services/        # Business logic
│   │   │   ├── models/          # Data models
│   │   │   ├── dto/             # Data Transfer Objects
│   │   │   ├── middleware/      # HTTP middleware
│   │   │   └── router/          # Route definitions
│   │   ├── database/            # Database connections
│   │   ├── global/              # Global variables
│   │   ├── logger/              # Logging utilities
│   │   └── utility/              # Utility functions
│   └── config/                  # Configuration files
├── api-tests/                    # Module test
│   ├── cases/                   # Test cases
│   ├── utils/                   # Test utilities
│   └── scripts/                 # Test scripts
├── docs/                        # Tài liệu hệ thống
└── deploy_notes/                # Ghi chú triển khai
```

## 🔧 Cấu Hình

### Biến Môi Trường Quan Trọng

| Biến | Mô Tả | Ví Dụ |
|------|-------|-------|
| `ADDRESS` | Port server | `8080` |
| `JWT_SECRET` | Secret key cho JWT | `your-secret-key` |
| `MONGODB_CONNECTION_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `MONGODB_DBNAME_AUTH` | Database name cho auth | `folkform_auth` |
| `FIREBASE_PROJECT_ID` | Firebase project ID | `your-project-id` |
| `FIREBASE_CREDENTIALS_PATH` | Đường dẫn đến service account JSON | `config/firebase/service-account.json` |

Xem chi tiết tại [docs/01-getting-started/cau-hinh.md](docs/01-getting-started/cau-hinh.md)

## 📚 Tài Liệu

### Tài Liệu Chính

- [📖 Tổng Quan Tài Liệu](docs/README.md) - Index của tất cả tài liệu
- [🚀 Bắt Đầu](docs/01-getting-started/) - Hướng dẫn cài đặt và cấu hình
- [🏗️ Kiến Trúc](docs/02-architecture/) - Kiến trúc và thiết kế hệ thống
- [🔌 API Reference](docs/03-api/) - Tài liệu API endpoints
- [🚢 Triển Khai](docs/04-deployment/) - Hướng dẫn deploy
- [💻 Phát Triển](docs/05-development/) - Hướng dẫn phát triển
- [🧪 Testing](docs/06-testing/) - Hướng dẫn test
- [🔧 Xử Lý Sự Cố](docs/07-troubleshooting/) - Troubleshooting

### Tài Liệu Nhanh

- [📝 Tài Liệu Ngắn Gọn](docs/tai-lieu-he-thong.md) - Tổng quan nhanh về hệ thống
- [🔐 Firebase Authentication](docs/firebase-auth-voi-database.md) - Tích hợp Firebase
- [🔄 Quy Trình Khởi Tạo](docs/quy-trinh-khoi-tao-he-thong.md) - Khởi tạo hệ thống lần đầu

## 🧪 Testing

### Chạy Test Suite

```powershell
# Từ root directory
.\api-tests\test.ps1
```

Script sẽ tự động:
- Kiểm tra server có đang chạy
- Khởi động server nếu cần
- Chạy toàn bộ test suite
- Tạo báo cáo test

Xem chi tiết tại [README_TEST.md](README_TEST.md)

## 🔐 Authentication Flow

1. **Frontend**: User đăng nhập bằng Firebase SDK (Email/Google/Facebook/Phone)
2. **Firebase**: Trả về Firebase ID Token
3. **Frontend**: Gửi ID Token đến `/api/v1/auth/login/firebase`
4. **Backend**: Verify token, tạo/update user trong MongoDB, trả về JWT
5. **Frontend**: Lưu JWT để sử dụng cho các request tiếp theo

Xem chi tiết tại [docs/02-architecture/authentication.md](docs/02-architecture/authentication.md)

## 🛠️ Công Nghệ Sử Dụng

- **Language**: Go 1.23+
- **Framework**: Fiber v3
- **Database**: MongoDB
- **Authentication**: Firebase Authentication
- **Logging**: Logrus
- **Validation**: go-playground/validator

## 📝 License

[Thêm thông tin license nếu có]

## 🤝 Đóng Góp

[Thêm hướng dẫn đóng góp nếu có]

## 📞 Liên Hệ

[Thêm thông tin liên hệ nếu có]

---

**Lưu ý**: Đây là tài liệu tổng quan. Để biết chi tiết, vui lòng xem các tài liệu trong thư mục `docs/`.

