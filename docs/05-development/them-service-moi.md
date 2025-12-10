# Thêm Service Mới

Hướng dẫn thêm service mới vào hệ thống.

## 📋 Tổng Quan

Service chứa business logic của ứng dụng. Tài liệu này hướng dẫn cách tạo service mới.

## 🚀 Các Bước

### 1. Tạo Service Struct

**File:** `api/core/api/services/service.<module>.<entity>.go`

```go
package services

import (
    "context"
    "meta_commerce/core/api/models/mongodb"
    "go.mongodb.org/mongo-driver/mongo"
)

type EntityService struct {
    *BaseService[mongodb.Entity]
    // Additional dependencies
}

func NewEntityService() (*EntityService, error) {
    baseService, err := NewBaseService[mongodb.Entity]("entities")
    if err != nil {
        return nil, err
    }
    
    return &EntityService{
        BaseService: baseService,
    }, nil
}
```

### 2. Thêm Business Logic

```go
func (s *EntityService) CustomBusinessLogic(ctx context.Context, input *CustomInput) (*mongodb.Entity, error) {
    // Business logic here
    entity := &mongodb.Entity{
        Name: input.Name,
    }
    
    result, err := s.InsertOne(ctx, entity)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}
```

### 3. Sử Dụng Service trong Handler

```go
func (h *EntityHandler) HandleCustomAction(c fiber.Ctx) error {
    var input dto.CustomInput
    if err := h.ParseRequestBody(c, &input); err != nil {
        h.HandleResponse(c, nil, err)
        return nil
    }
    
    result, err := h.entityService.CustomBusinessLogic(context.Background(), &input)
    if err != nil {
        h.HandleResponse(c, nil, err)
        return nil
    }
    
    h.HandleResponse(c, result, nil)
    return nil
}
```

## 📝 Best Practices

1. **Separation of Concerns**: Service chỉ chứa business logic
2. **Error Handling**: Xử lý lỗi đúng cách
3. **Validation**: Validate input trước khi xử lý
4. **Transactions**: Sử dụng transactions cho operations phức tạp

## 📚 Tài Liệu Liên Quan

- [Cấu Trúc Code](cau-truc-code.md)
- [Thêm API Mới](them-api-moi.md)
- [Coding Standards](coding-standards.md)

