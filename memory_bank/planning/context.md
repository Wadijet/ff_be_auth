# Current Context

## Status Overview
- Sprint: 1
- Phase: Giai đoạn 1 - Chuẩn bị
- Progress: 40% (8/20 points)

## Active Tasks
1. Framework Migration (In Progress)
   - Chuyển từ FastHTTP sang Fiber
   - Routes migration đang được thực hiện
   - Middleware migration sẽ được thực hiện tiếp theo

2. Clean Architecture Setup (In Progress)
   - Cấu trúc thư mục đã hoàn thành
   - Đang refactor code theo clean architecture
   - Cần setup dependency injection

## Next Focus
1. Metadata System
   - Thiết kế schema cho auth flows, API routes và database config
   - Tạo cấu trúc thư mục metadata với các thành phần auth/, api/, db/
   - Setup metadata parser với YAML và hot reload

2. Middleware Configuration
   - Setup logger middleware
   - Implement error handler
   - Configure CORS

## Technical Context
- Framework: Fiber (migration in progress)
- Architecture: Clean Architecture
- Database: MongoDB (sẽ được cấu hình qua metadata)
- Authentication: JWT (planned)

## Dependencies
- go.mod đã được setup
- Fiber framework đã được thêm
- MongoDB driver đã được cài đặt

## Notes
- Migration được thực hiện theo từng phần để giảm thiểu rủi ro
- Clean Architecture giúp code dễ maintain và test hơn
- Metadata system sẽ giúp cấu hình linh hoạt hơn

## Risks
1. Framework Migration
   - Impact: High
   - Mitigation: Phased approach và extensive testing

2. Performance
   - Impact: Medium
   - Mitigation: Benchmarking và optimization

## Môi trường
- OS: Windows (win32 10.0.19045)
- Shell: PowerShell
- Workspace: D:\Crossborder\ff_be_auth

## Cấu trúc hiện tại
### Framework & Libraries
- Web Framework: FastHTTP + Router
- Database: MongoDB
- Validation: go-playground/validator
- Logging: logrus
- JWT: dgrijalva/jwt-go

### Cấu trúc thư mục
- app/
  + models/: MongoDB models
  + services/: Business logic
  + database/: DB operations
  + router/: API routes
  + handler/: Request handlers
  + middleware/: HTTP middleware
  + utility/: Helper functions
  + global/: Global variables

### Tính năng hiện có
- Authentication & Authorization
- User Management
- Role & Permission Management
- Rate Limiting
- CORS & Timeout
- Panic Recovery

## Lưu ý
### Vấn đề cần giải quyết
- Global variables trong global/
- Thiếu dependency injection
- Chưa tách biệt rõ các layer
- Error handling chưa chuẩn hóa

### Ưu tiên hiện tại
1. Tái cấu trúc theo Clean Architecture
2. Cải thiện dependency injection
3. Chuẩn hóa error handling

### Kế hoạch ngắn hạn
1. Tạo cấu trúc thư mục mới
2. Di chuyển code từng module
3. Tách interface và implementation
4. Tạo service container 