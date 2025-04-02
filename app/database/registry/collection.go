package registry

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

// CollectionRegistry là một singleton quản lý và cung cấp truy cập đến các MongoDB collections
// Đảm bảo rằng mỗi collection chỉ được khởi tạo một lần và có thể truy cập từ bất kỳ đâu trong ứng dụng
type CollectionRegistry struct {
	// collections lưu trữ map của tên collection và instance của collection
	// Key: tên collection (string)
	// Value: con trỏ đến mongo.Collection
	collections map[string]*mongo.Collection

	// mu là mutex để đảm bảo thread-safe khi truy cập collections
	mu sync.RWMutex
}

var (
	// registry là instance duy nhất của CollectionRegistry
	registry *CollectionRegistry
	// once đảm bảo registry chỉ được khởi tạo một lần
	once sync.Once
)

// GetRegistry trả về instance duy nhất của CollectionRegistry
// Sử dụng sync.Once để đảm bảo registry chỉ được khởi tạo một lần
// Returns:
//   - *CollectionRegistry: Instance duy nhất của registry
func GetRegistry() *CollectionRegistry {
	once.Do(func() {
		registry = &CollectionRegistry{
			collections: make(map[string]*mongo.Collection),
		}
	})
	return registry
}

// RegisterCollection đăng ký một collection mới vào registry
// Parameters:
//   - name: Tên của collection (string)
//   - collection: Instance của mongo.Collection cần đăng ký
//
// Thread-safe: Có sử dụng mutex để đảm bảo an toàn khi ghi
func (r *CollectionRegistry) RegisterCollection(name string, collection *mongo.Collection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.collections[name] = collection
}

// GetCollection lấy collection theo tên
// Parameters:
//   - name: Tên của collection cần lấy
//
// Returns:
//   - *mongo.Collection: Collection tìm thấy hoặc nil nếu không tồn tại
//
// Thread-safe: Có sử dụng mutex để đảm bảo an toàn khi đọc
func (r *CollectionRegistry) GetCollection(name string) *mongo.Collection {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.collections[name]
}

// MustGetCollection lấy collection theo tên, panic nếu không tìm thấy
// Parameters:
//   - name: Tên của collection cần lấy
//
// Returns:
//   - *mongo.Collection: Collection tìm thấy
//
// Panics:
//   - Nếu collection không tồn tại
//
// Thread-safe: Có sử dụng mutex để đảm bảo an toàn khi đọc
func (r *CollectionRegistry) MustGetCollection(name string) *mongo.Collection {
	collection := r.GetCollection(name)
	if collection == nil {
		panic("collection not found: " + name)
	}
	return collection
}

// CollectionNames trả về danh sách tên các collection đã đăng ký
// Returns:
//   - []string: Danh sách tên các collection
//
// Thread-safe: Có sử dụng mutex để đảm bảo an toàn khi đọc
func (r *CollectionRegistry) CollectionNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.collections))
	for name := range r.collections {
		names = append(names, name)
	}
	return names
}

// Clear xóa tất cả collections đã đăng ký
// Sử dụng trong trường hợp cần reset registry hoặc cleanup
// Thread-safe: Có sử dụng mutex để đảm bảo an toàn khi ghi
func (r *CollectionRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.collections = make(map[string]*mongo.Collection)
}
