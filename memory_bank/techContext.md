# Technical Context - FolkForm Authentication Backend

## Kiến Trúc Hệ Thống

### API Layer (`core/api/`)
- **Handler**: Xử lý requests và responses
- **Middleware**: 
  - Authentication middleware với permission-based access control
  - Logging và error handling
- **Router**: 
  - Sử dụng Fiber framework
  - Cấu trúc RESTful với versioning (/api/v1)
  - CRUD operations tự động cho các resource
  - Custom routes cho authentication và specific features
- **Models**: 
  - MongoDB models với BSON/JSON mapping
  - Validation rules cho input data
  - Relationship handling (User-Role-Permission)

### Database Layer (`core/database/`)
- MongoDB integration
- Repository pattern implementation
- Data access interfaces

### Business Logic (`core/services/`)
- Authentication service
- Authorization service
- User management service

### Configuration (`config/`)
- Environment-based configuration
- Application settings
- Database configuration

## Dependencies
```go
// Từ go.mod
module ff_be_auth

go 1.21

// Các dependencies chính cần được liệt kê từ go.mod
```

## Development Tools
- Go 1.21
- MongoDB
- Custom development modes (VAN, PLAN, CREATIVE, IMPLEMENT, QA)

## Cấu Trúc Thư Mục Chi Tiết
```
ff_be_auth/
├── cmd/
│   ├── server/     # HTTP server entry point
│   └── worker/     # Background worker entry point
├── core/
│   ├── api/        # API layer
│   ├── database/   # Database access
│   ├── services/   # Business logic
│   └── models/     # Domain models
├── config/         # Configuration
└── custom_modes/   # Development modes
```

## Quy Ước Phát Triển
1. Clean Architecture principles
2. RESTful API design
3. Error handling standards
4. Logging conventions
5. Testing requirements

## Security Considerations
- Authentication mechanisms
- Authorization flows
- Token management
- Password hashing
- Session handling

## Performance Goals
- Response time targets
- Throughput requirements
- Scalability considerations
- Resource utilization limits

## Monitoring & Logging
- Log levels và formats
- Metrics collection
- Performance monitoring
- Error tracking

## Testing Strategy
- Unit tests
- Integration tests
- API tests
- Performance tests
- Security tests 

### Database Schema
#### User Model
```go
type User struct {
    ID        primitive.ObjectID
    Name      string
    Email     string             // unique index
    Password  string
    Salt      string
    Token     string
    Tokens    []Token           // Multi-device support
    IsBlock   bool
    BlockNote string
    CreatedAt int64
    UpdatedAt int64
}
```

#### Role & Permission System
- **Role**: Định nghĩa vai trò người dùng
- **Permission**: Các quyền trong hệ thống
- **RolePermission**: Mapping giữa Role và Permission
- **UserRole**: Mapping giữa User và Role

### API Endpoints
#### Authentication
- POST `/api/v1/users/login`: Đăng nhập
- POST `/api/v1/users/register`: Đăng ký
- POST `/api/v1/users/logout`: Đăng xuất
- GET `/api/v1/users/profile`: Lấy thông tin profile
- PUT `/api/v1/users/profile`: Cập nhật profile
- PUT `/api/v1/users/change-password`: Đổi mật khẩu

#### CRUD Operations
Mỗi resource (User, Role, Permission, etc.) có các endpoints:
- POST `/`: Create one
- POST `/batch`: Create many
- GET `/`: Find all
- GET `/:id`: Find by ID
- PUT `/:id`: Update by ID
- DELETE `/:id`: Delete by ID 