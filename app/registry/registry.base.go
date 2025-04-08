// Package registry cung cấp implementation của registry pattern với generic type.
// Package này cho phép quản lý các singleton instances trong ứng dụng một cách thread-safe.
// Sử dụng generic type để có thể tái sử dụng cho nhiều loại đối tượng khác nhau.
package registry

import (
	"fmt"
	"sync"
)

// globalRegistry lưu trữ instance của registry
var globalRegistry = NewRegistry[any]()

// GetRegistry trả về instance của global registry.
//
// Returns:
//   - *Registry[any]: Instance của global registry
//
// Thread-safety: Safe for concurrent use
//
// Example:
//
//	registry := GetRegistry()
//	registry.Register("myService", service)
func GetRegistry() *Registry[any] {
	return globalRegistry
}

// RegistryError đại diện cho các lỗi từ registry operations
type RegistryError struct {
	Message string
}

func (e *RegistryError) Error() string {
	return e.Message
}

// Registry là một thread-safe generic registry pattern implementation.
// Type parameter T cho phép registry quản lý bất kỳ loại object nào.
// Thread-safety được đảm bảo thông qua sync.RWMutex.
//
// Example:
//
//	// Tạo registry cho kiểu string
//	strRegistry := NewRegistry[string]()
//
//	// Đăng ký một item
//	strRegistry.Register("key", "value")
//
//	// Lấy item
//	if value, exists := strRegistry.Get("key"); exists {
//	    fmt.Println(value)
//	}
type Registry[T any] struct {
	items map[string]T // Map lưu trữ các items theo key
	mu    sync.RWMutex // Mutex để đảm bảo thread-safety
}

// NewRegistry tạo và trả về một registry mới.
// Generic type T xác định loại items mà registry sẽ quản lý.
//
// Returns:
//   - *Registry[T]: Registry instance mới, đã được khởi tạo
//
// Example:
//
//	registry := NewRegistry[int]()
func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		items: make(map[string]T),
	}
}

// Register đăng ký một item mới vào registry.
// Nếu item với name đã tồn tại, nó sẽ bị ghi đè.
//
// Parameters:
//   - name: Định danh duy nhất cho item
//   - item: Item cần đăng ký
//
// Returns:
//   - error: Trả về lỗi nếu name rỗng
//
// Thread-safety: Safe for concurrent use
//
// Example:
//
//	err := registry.Register("counter", 42)
func (r *Registry[T]) Register(name string, item T) error {
	if name == "" {
		return &RegistryError{Message: "name cannot be empty"}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[name] = item
	return nil
}

// Get lấy item theo tên.
// Trả về item và một boolean cho biết item có tồn tại hay không.
//
// Parameters:
//   - name: Tên của item cần lấy
//
// Returns:
//   - T: Item nếu tìm thấy, zero value của T nếu không tìm thấy
//   - bool: true nếu item tồn tại, false nếu không
//
// Thread-safety: Safe for concurrent use
//
// Example:
//
//	if value, exists := registry.Get("counter"); exists {
//	    fmt.Printf("Counter value: %d\n", value)
//	}
func (r *Registry[T]) Get(name string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.items[name]
	return item, exists
}

// MustGet lấy item theo tên, trả về error nếu không tìm thấy.
//
// Parameters:
//   - name: Tên của item cần lấy
//
// Returns:
//   - T: Item nếu tìm thấy
//   - error: Error nếu không tìm thấy item
//
// Thread-safety: Safe for concurrent use
//
// Example:
//
//	value, err := registry.MustGet("counter")
//	if err != nil {
//	    log.Printf("Error: %v", err)
//	    return
//	}
func (r *Registry[T]) MustGet(name string) (T, error) {
	if item, exists := r.Get(name); exists {
		return item, nil
	}
	var zero T
	return zero, &RegistryError{Message: fmt.Sprintf("item not found: %s", name)}
}

// Delete xóa một item khỏi registry.
//
// Parameters:
//   - name: Tên của item cần xóa
//
// Returns:
//   - bool: true nếu item đã được xóa, false nếu item không tồn tại
//
// Thread-safety: Safe for concurrent use
//
// Example:
//
//	if deleted := registry.Delete("counter"); deleted {
//	    fmt.Println("Counter was deleted")
//	}
func (r *Registry[T]) Delete(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[name]; exists {
		delete(r.items, name)
		return true
	}
	return false
}

// GetAll trả về danh sách tên của tất cả items trong registry.
//
// Returns:
//   - []string: Danh sách tên các items
//
// Thread-safety: Safe for concurrent use
//
// Example:
//
//	names := registry.GetAll()
//	for _, name := range names {
//	    fmt.Printf("Item: %s\n", name)
//	}
func (r *Registry[T]) GetAll() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.items))
	for name := range r.items {
		names = append(names, name)
	}
	return names
}

// Clear xóa tất cả items trong registry.
//
// Thread-safety: Safe for concurrent use
//
// Example:
//
//	registry.Clear()
func (r *Registry[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]T)
}
