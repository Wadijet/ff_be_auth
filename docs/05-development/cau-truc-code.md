# Cấu Trúc Code

Tài liệu về cấu trúc và tổ chức code trong dự án.

## 📋 Tổng Quan

Dự án được tổ chức theo kiến trúc layered với các layer rõ ràng.

## 🏗️ Cấu Trúc Thư Mục

```
api/
├── cmd/
│   └── server/          # Entry point
├── core/
│   ├── api/            # API layer
│   │   ├── handler/    # HTTP handlers
│   │   ├── services/   # Business logic
│   │   ├── models/     # Data models
│   │   ├── dto/        # Data Transfer Objects
│   │   ├── middleware/ # HTTP middleware
│   │   └── router/     # Route definitions
│   ├── database/       # Database connections
│   ├── global/         # Global variables
│   ├── logger/         # Logging utilities
│   └── utility/        # Utility functions
└── config/             # Configuration files
```

## 📝 Naming Conventions

### Files

- Handler: `handler.<module>.<entity>.go`
- Service: `service.<module>.<entity>.go`
- Model: `model.<module>.<entity>.go`
- DTO: `dto.<module>.<entity>.go`

### Functions

- Handler: `Handle<Action><Entity>`
- Service: `<Action><Entity>`
- Utility: `<Action><Entity>`

## 🔄 Flow

```
Request → Router → Middleware → Handler → Service → Repository → Database
```

## 📚 Tài Liệu Liên Quan

- [Thêm API Mới](them-api-moi.md)
- [Thêm Service Mới](them-service-moi.md)
- [Coding Standards](coding-standards.md)

