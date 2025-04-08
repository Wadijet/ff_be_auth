# Phase 0: Registry Refactoring

## Quy ước chung

### 1. Quy ước đặt tên file
Tất cả các file trong thư mục registry sẽ theo format: `registry.{name}.go`

Ví dụ:
- `registry.base.go`: Generic registry implementation
- `registry.collection.go`: Collection registry
- `registry.init.go`: Initialization code
- `registry.etl.go`: ETL registry (sẽ thêm sau)
- `registry.types.go`: Type definitions (nếu cần)
- `registry.test.go`: Test utilities (nếu cần)

### 2. Quy ước Documentation
- Mỗi file phải có file header mô tả mục đích và cách sử dụng
- Mỗi type/interface phải có comment mô tả đầy đủ
- Mỗi function phải có comment theo format:
  ```go
  // FunctionName mô tả ngắn gọn chức năng
  // 
  // Chi tiết về cách hoạt động và use cases
  //
  // Parameters:
  //   - param1: mô tả parameter
  //   - param2: mô tả parameter
  //
  // Returns:
  //   - returnType: mô tả giá trị trả về
  //
  // Errors:
  //   - error1: trong trường hợp nào
  //   - error2: trong trường hợp nào
  //
  // Thread-safety: mô tả về thread-safety
  ```

## Implementation Details

### 1. Generic Registry Base
```go
// registry.base.go

// Package registry cung cấp implementation của registry pattern với generic type
// Sử dụng để quản lý các singleton instances trong ứng dụng
package registry

import (
    "fmt"
    "sync"
)

// Registry là một thread-safe generic registry pattern implementation
// Generic type T cho phép registry quản lý bất kỳ loại object nào
// Thread-safety được đảm bảo thông qua sync.RWMutex
type Registry[T any] struct {
    items map[string]T    // Map lưu trữ các items theo key
    mu    sync.RWMutex    // Mutex để đảm bảo thread-safety
}

// NewRegistry tạo và trả về một registry mới
// Generic type T xác định loại items mà registry sẽ quản lý
//
// Returns:
//   - *Registry[T]: Registry instance mới, đã được khởi tạo
func NewRegistry[T any]() *Registry[T] {
    return &Registry[T]{
        items: make(map[string]T),
    }
}

// Register đăng ký một item mới vào registry
//
// Parameters:
//   - name: Định danh duy nhất cho item
//   - item: Item cần đăng ký
//
// Thread-safety: Safe for concurrent use
func (r *Registry[T]) Register(name string, item T) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.items[name] = item
}

// Get lấy item theo tên
func (r *Registry[T]) Get(name string) T {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.items[name]
}

// MustGet lấy item, panic nếu không tồn tại
func (r *Registry[T]) MustGet(name string) T {
    item := r.Get(name)
    var zero T
    if item == zero {
        panic("item not found: " + name)
    }
    return item
}

// GetAll trả về danh sách tên các items
func (r *Registry[T]) GetAll() []string {
    r.mu.RLock()
    defer r.mu.RUnlock()
    names := make([]string, 0, len(r.items))
    for name := range r.items {
        names = append(names, name)
    }
    return names
}

// Clear xóa tất cả items
func (r *Registry[T]) Clear() {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.items = make(map[string]T)
}
```

### 2. Collection Registry
```go
// registry.collection.go

// Package registry cung cấp implementation của MongoDB collection registry
// Đảm bảo mỗi collection chỉ được khởi tạo một lần và có thể truy cập từ bất kỳ đâu
package registry

import (
    "go.mongodb.org/mongo-driver/mongo"
)

// Các biến global cho collection registry
var (
    collectionRegistry *Registry[*mongo.Collection]  // Registry instance cho MongoDB collections
    collectionOnce    sync.Once                     // Đảm bảo singleton pattern
)

// GetCollectionRegistry trả về instance của collection registry
func GetCollectionRegistry() *Registry[*mongo.Collection] {
    collectionOnce.Do(func() {
        collectionRegistry = NewRegistry[*mongo.Collection]()
    })
    return collectionRegistry
}

// Các hàm wrapper để giữ backward compatibility
func RegisterCollection(name string, collection *mongo.Collection) {
    GetCollectionRegistry().Register(name, collection)
}

func GetCollection(name string) *mongo.Collection {
    return GetCollectionRegistry().Get(name)
}

func MustGetCollection(name string) *mongo.Collection {
    return GetCollectionRegistry().MustGet(name)
}

func CollectionNames() []string {
    return GetCollectionRegistry().GetAll()
}
```

### 3. Init Registry
```go
// registry.init.go

// Package registry cung cấp các functions khởi tạo cho registry system
package registry

import (
    "go.mongodb.org/mongo-driver/mongo"
)

func InitCollections(db *mongo.Client, cfg *config.Configuration) error {
    registry := GetCollectionRegistry()
    database := db.Database(cfg.MongoDB_DBNameAuth)
    
    for _, name := range GetCollectionNames() {
        registry.Register(name, database.Collection(name))
    }
    return nil
}
```

## Các bước thực hiện:

1. Documentation & Comments
   - Viết file headers
   - Comments cho types/interfaces
   - Comments cho functions
   - Examples trong comments

2. Implementation
   - Tạo `registry.base.go` với generic registry
   - Rename và update các file hiện có
   - Cập nhật code với comments đầy đủ

3. Testing
   - Unit tests với comments đầy đủ
   - Test cases cho zero values
   - Test concurrent access
   - Test backward compatibility

4. Verification
   - Code review
   - Documentation review
   - Test coverage review

## Lợi ích:
1. Code gọn hơn, tái sử dụng tốt hơn
2. Type safety với generics
3. Documentation đầy đủ, dễ hiểu
4. Dễ dàng maintain và mở rộng
5. Cấu trúc file rõ ràng, nhất quán

## Migration Guide:
1. Tạo files mới với cấu trúc mới
2. Chuyển đổi code dần dần
3. Giữ backward compatibility
4. Xóa code cũ khi đã stable

Bạn có muốn chuyển sang CREATIVE mode để bắt đầu implementation không? 