# Coding Standards

Tiêu chuẩn code cho dự án.

## 📋 Tổng Quan

Tài liệu này mô tả các tiêu chuẩn code cần tuân thủ trong dự án.

## 📝 Naming Conventions

### Files

- Handler: `handler.<module>.<entity>.go`
- Service: `service.<module>.<entity>.go`
- Model: `model.<module>.<entity>.go`
- DTO: `dto.<module>.<entity>.go`

### Functions

- Handler: `Handle<Action><Entity>`
- Service: `<Action><Entity>`
- Public: PascalCase
- Private: camelCase

### Variables

- Constants: UPPER_SNAKE_CASE
- Variables: camelCase
- Exported: PascalCase

## 🏗️ Code Structure

### Handler

```go
func (h *Handler) HandleAction(c fiber.Ctx) error {
    // 1. Parse input
    var input dto.Input
    if err := h.ParseRequestBody(c, &input); err != nil {
        h.HandleResponse(c, nil, err)
        return nil
    }
    
    // 2. Call service
    result, err := h.service.Action(context.Background(), &input)
    if err != nil {
        h.HandleResponse(c, nil, err)
        return nil
    }
    
    // 3. Return response
    h.HandleResponse(c, result, nil)
    return nil
}
```

### Service

```go
func (s *Service) Action(ctx context.Context, input *dto.Input) (*models.Entity, error) {
    // 1. Validate
    if err := validate(input); err != nil {
        return nil, err
    }
    
    // 2. Business logic
    entity := &models.Entity{
        // ...
    }
    
    // 3. Save
    result, err := s.InsertOne(ctx, entity)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}
```

## ✅ Best Practices

1. **Error Handling**: Luôn xử lý lỗi
2. **Context**: Sử dụng context cho tất cả operations
3. **Validation**: Validate input ở handler
4. **Logging**: Log errors và important events
5. **Comments**: Comment cho public functions

## 📚 Tài Liệu Liên Quan

- [Cấu Trúc Code](cau-truc-code.md)
- [Thêm API Mới](them-api-moi.md)

