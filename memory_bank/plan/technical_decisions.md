# Quyết định kỹ thuật ETL Pipeline

## 1. Kiến trúc tổng thể

### 1.1 Generic Registry Pattern
- **Quyết định**: Sử dụng generic registry pattern
- **Lý do**:
  - Giảm thiểu code trùng lặp
  - Tái sử dụng logic cho mọi loại component
  - Type safety với Go generics
- **Implementation**:
  ```go
  type Registry[T any] struct {
      items map[string]T
      mu sync.RWMutex
  }
  ```
- **Trade-offs**:
  + Ưu điểm:
    * Code ít và tập trung
    * Type safety tại compile-time
    * Dễ mở rộng cho components mới
    * Tái sử dụng 100% logic registry
  - Nhược điểm:
    * Cần refactor code hiện có
    * Yêu cầu Go version hỗ trợ generics

### 1.2 Component-based Architecture
- **Quyết định**: Chia nhỏ thành các components độc lập
- **Lý do**:
  - Dễ dàng mở rộng và bảo trì
  - Có thể thay thế từng component
  - Tái sử dụng components
- **Components**:
  1. DataSource
  2. Transformer
  3. Destination
  4. Pipeline

## 2. Thiết kế Interface

### 2.1 Generic Interfaces
- **Quyết định**: Sử dụng []byte cho data transfer
- **Lý do**:
  - Linh hoạt với nhiều loại dữ liệu
  - Dễ dàng chuyển đổi giữa JSON/YAML
  - Hiệu năng tốt với dữ liệu lớn

### 2.2 Configuration Interface
- **Quyết định**: Mỗi component có GetConfig()
- **Lý do**:
  - Dễ dàng serialize/deserialize
  - Thuận tiện cho việc lưu trữ
  - Hỗ trợ validation

## 3. Configuration Management

### 3.1 YAML Format
- **Quyết định**: Sử dụng YAML cho config
- **Lý do**:
  - Dễ đọc và viết
  - Hỗ trợ comments
  - Phù hợp cho cấu trúc phức tạp
- **Structure**:
  ```yaml
  component:
    id: string
    type: string
    config: object
  ```

### 3.2 Config Organization
- **Quyết định**: Tách config theo loại component
- **Lý do**:
  - Dễ quản lý
  - Có thể reuse configs
  - Dễ version control

## 4. Error Handling & Logging

### 4.1 Error Propagation
- **Quyết định**: Sử dụng error wrapping
- **Lý do**:
  - Dễ debug
  - Giữ context của lỗi
  - Phù hợp với Go's error handling

### 4.2 Logging Strategy
- **Quyết định**: Structured logging
- **Lý do**:
  - Dễ parse và analyze
  - Hỗ trợ monitoring
  - Consistent format

## 5. Testing Strategy

### 5.1 Test Levels
- **Quyết định**: Unit + Integration tests
- **Lý do**:
  - Coverage đầy đủ
  - Phát hiện lỗi sớm
  - Đảm bảo chất lượng

### 5.2 Test Data
- **Quyết định**: Sử dụng test fixtures
- **Lý do**:
  - Tái sử dụng test data
  - Dễ maintain
  - Consistent test cases 