# Kiến trúc hệ thống

## Kiến trúc hiện tại
### Cấu trúc thư mục
```
app/
├── models/      # MongoDB models
├── services/    # Business logic
├── database/    # DB operations
├── router/      # API routes
├── handler/     # Request handlers
├── middleware/  # HTTP middleware
├── utility/     # Helper functions
└── global/      # Global variables
```

### Các thành phần chính
1. **Models Layer**
   - Định nghĩa cấu trúc dữ liệu
   - MongoDB schema
   - Validation rules

2. **Services Layer**
   - Business logic
   - Data processing
   - External service integration

3. **Database Layer**
   - MongoDB operations
   - Data access methods
   - Transaction handling

4. **Router Layer**
   - API route definitions
   - Route grouping
   - Middleware binding

5. **Handler Layer**
   - Request handling
   - Response formatting
   - Input validation

6. **Middleware Layer**
   - Authentication
   - Authorization
   - Rate limiting
   - CORS
   - Recovery

7. **Utility Layer**
   - Helper functions
   - Common utilities
   - Shared constants

8. **Global Layer**
   - Global variables
   - Shared instances
   - Configuration

### Vấn đề hiện tại
1. **Tight Coupling**
   - Global variables gây khó kiểm soát
   - Các layer phụ thuộc chặt chẽ
   - Khó thay đổi implementation

2. **Dependency Management**
   - Thiếu dependency injection
   - Sử dụng global state
   - Khó unit test

3. **Error Handling**
   - Chưa chuẩn hóa error types
   - Thiếu error middleware
   - Response format không nhất quán

## Kế hoạch cải tiến
### Clean Architecture
```
app/
├── domain/           # Business entities & interfaces
│   ├── entity/      # Domain models
│   ├── repository/  # Repository interfaces
│   └── service/     # Service interfaces
├── infrastructure/  # External implementations
│   ├── persistence/ # Database implementations
│   └── auth/        # Auth implementations
├── interfaces/      # Interface adapters
│   ├── http/        # HTTP handlers
│   └── middleware/  # HTTP middleware
└── application/     # Use cases & business logic
    ├── service/     # Service implementations
    └── dto/         # Data transfer objects
```

### Các layer mới
1. **Domain Layer**
   - Business entities
   - Repository interfaces
   - Service interfaces
   - Domain errors
   - Value objects

2. **Infrastructure Layer**
   - Database implementations
   - External service adapters
   - Authentication providers
   - Caching implementations

3. **Interface Layer**
   - HTTP handlers
   - Middleware
   - Request/Response DTOs
   - API versioning

4. **Application Layer**
   - Use cases
   - Service implementations
   - Business rules
   - Transaction scripts

### Dependency Injection
1. **Service Container**
   - Centralized dependency management
   - Interface-based injection
   - Scoped instances
   - Testing support

2. **Interface Segregation**
   - Minimal interfaces
   - Clear dependencies
   - Loose coupling

### Error Handling
1. **Domain Errors**
   - Business error types
   - Error hierarchies
   - Error codes

2. **Error Middleware**
   - Centralized error handling
   - Error logging
   - Response formatting

3. **Response Standards**
   - Consistent error format
   - HTTP status codes
   - Error messages

## Quy trình triển khai
1. **Chuẩn bị**
   - Tạo cấu trúc thư mục mới
   - Định nghĩa interfaces
   - Setup dependency injection

2. **Di chuyển code**
   - Chuyển entities vào domain
   - Implement repositories
   - Tạo use cases
   - Cập nhật handlers

3. **Tối ưu hóa**
   - Loại bỏ global state
   - Chuẩn hóa error handling
   - Thêm unit tests

4. **Kiểm thử**
   - Unit testing
   - Integration testing
   - Performance testing 