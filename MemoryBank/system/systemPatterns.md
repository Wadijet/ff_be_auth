# System Patterns

## Nguyên tắc phát triển

### Tối giản và hiệu quả
- Luôn lựa chọn giải pháp tinh, gọn, phù hợp để nhanh chóng đưa dự án vào hoạt động
- Ưu tiên các giải pháp đơn giản nhất có thể giải quyết vấn đề
- Tránh over-engineering và phức tạp hóa không cần thiết

### Tái sử dụng code
- Luôn đề cao việc tái sử dụng code, làm sao để code ít nhất có thể
- Tạo các utility và helper function cho các tác vụ lặp lại
- Trừu tượng hóa các logic chung thành các module có thể tái sử dụng

## Design Patterns

### Middleware Chain
Ứng dụng sử dụng mẫu Middleware Chain để xử lý các request HTTP. Các middleware được áp dụng theo thứ tự cụ thể:
1. Recovery - Xử lý panic
2. Timeout - Giới hạn thời gian xử lý request
3. RateLimit - Giới hạn số lượng request từ một IP
4. Measure - Đo lường hiệu suất
5. CORS - Xử lý Cross-Origin Resource Sharing

### Repository Pattern
Tương tác với MongoDB được trừu tượng hóa qua Repository Pattern, giúp tách biệt logic nghiệp vụ và truy cập dữ liệu.

### Factory Pattern
Sử dụng Factory Pattern để khởi tạo các đối tượng dịch vụ và repository.

## Coding Patterns

### Error Handling
Xử lý lỗi thông qua logging và trả về response JSON với thông tin lỗi:
```go
utility.JSON(ctx, utility.Payload(false, data, "Thông báo lỗi"))
```

### Dependency Injection
Các dependency được inject qua constructor:
```go
func NewXXXService(dependency *SomeDependency) *XXXService {
    return &XXXService{dependency: dependency}
}
```

### Context Propagation
Sử dụng context.Context để truyền thông tin giữa các layer:
```go
func SomeFunction(ctx context.Context, param string) error {
    // ...
}
``` 