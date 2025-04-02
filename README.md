# Go Authentication Service

Dịch vụ xác thực và phân quyền được xây dựng bằng Go, sử dụng MongoDB làm database.

## Cấu trúc dự án

```
.
├── app/                    # Source code chính
│   ├── handler/           # Xử lý request/response
│   ├── middleware/        # Middleware (auth, logging, etc.)
│   ├── models/           # Data models
│   ├── router/           # Định nghĩa routes
│   ├── services/         # Business logic
│   └── utility/          # Các utility functions
├── config/               # Cấu hình ứng dụng
├── database/            # Database connections
├── global/              # Global variables và constants
└── tests/               # Test files
```

## Yêu cầu hệ thống

- Go 1.16 trở lên
- MongoDB 4.4 trở lên
- Make (optional, cho build automation)

## Cài đặt

1. Clone repository:
```bash
git clone https://github.com/your-org/ff_be_auth.git
cd ff_be_auth
```

2. Cài đặt dependencies:
```bash
go mod download
```

3. Tạo file .env từ template:
```bash
cp .env.example .env
```

4. Cập nhật các biến môi trường trong file .env

5. Build và chạy:
```bash
go build
./ff_be_auth
```

## API Documentation

### Authentication

- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/register` - Đăng ký
- `POST /api/v1/auth/refresh` - Làm mới token

### Users

- `GET /api/v1/users` - Lấy danh sách users
- `GET /api/v1/users/{id}` - Lấy thông tin user
- `POST /api/v1/users` - Tạo user mới
- `PUT /api/v1/users/{id}` - Cập nhật user
- `DELETE /api/v1/users/{id}` - Xóa user

### Roles & Permissions

- `GET /api/v1/roles` - Lấy danh sách roles
- `POST /api/v1/roles` - Tạo role mới
- `GET /api/v1/permissions` - Lấy danh sách permissions
- `POST /api/v1/permissions` - Tạo permission mới

## Testing

Chạy tests:
```bash
go test ./...
```

## Development

### Code Style

- Sử dụng `gofmt` để format code
- Tuân thủ [Effective Go](https://golang.org/doc/effective_go)
- Comment đầy đủ cho các hàm public

### Git Workflow

1. Tạo branch mới cho feature/bugfix
2. Commit code với message rõ ràng
3. Tạo pull request
4. Code review
5. Merge vào main branch

## Security

- Sử dụng JWT cho authentication
- Mã hóa password với bcrypt
- Rate limiting cho API endpoints
- CORS configuration
- Input validation

## Performance

- Connection pooling cho MongoDB
- Caching khi cần thiết
- Tối ưu database queries
- Sử dụng goroutines và channels hiệu quả

## Contributing

1. Fork repository
2. Tạo feature branch
3. Commit changes
4. Push to branch
5. Tạo Pull Request

## License

MIT License
