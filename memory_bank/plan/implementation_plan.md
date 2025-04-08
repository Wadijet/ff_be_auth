# Kế hoạch triển khai ETL Pipeline

## Phase 0: Registry Refactoring (1 ngày)

### 0.1 Generic Registry Implementation
```go
// app/registry/registry.go
type Registry[T any] struct {
    items map[string]T
    mu sync.RWMutex
}

func NewRegistry[T any]() *Registry[T] {
    return &Registry[T]{
        items: make(map[string]T),
    }
}
```

### 0.2 Chuyển đổi Collection Registry
```go
// app/registry/init.go
var (
    Collections = NewRegistry[*mongo.Collection]()
)

// Cập nhật các hàm hiện có để sử dụng generic registry
func InitCollections(db *mongo.Client, cfg *config.Configuration) error {
    database := db.Database(cfg.MongoDB_DBNameAuth)
    for _, name := range GetCollectionNames() {
        Collections.Register(name, database.Collection(name))
    }
    return nil
}
```

### 0.3 Unit Tests cho Generic Registry
```go
// app/registry/registry_test.go
func TestGenericRegistry(t *testing.T) {
    registry := NewRegistry[string]()
    registry.Register("test", "value")
    // ... more tests
}
```

## Phase 1: ETL Components (2-3 ngày)

### 1.1 ETL Registry Setup
```go
// app/registry/init.go
var (
    // Tận dụng generic registry
    DataSources = NewRegistry[types.DataSource]()
    Transformers = NewRegistry[types.Transformer]()
    Destinations = NewRegistry[types.Destination]()
    Pipelines = NewRegistry[types.Pipeline]()
)
```

### 1.2 Cấu trúc thư mục
```
app/
├── registry/
│   ├── collection.go    (hiện có)
│   ├── init.go         (cập nhật)
│   └── etl.go          (mới)
├── etl/
│   ├── types/          (interfaces chung)
│   ├── datasource/     (nguồn dữ liệu)
│   ├── transformer/    (xử lý dữ liệu)
│   ├── destination/    (đích đến)
│   └── pipeline/       (workflow)
└── config/
    └── etl/
        ├── datasources/
        ├── transformers/
        ├── destinations/
        └── pipelines/
```

### 1.3 Registry Implementation
```go
// app/registry/etl.go
type ETLComponentRegistry struct {
    dataSources  map[string]types.DataSource
    transformers map[string]types.Transformer
    destinations map[string]types.Destination
    pipelines    map[string]types.Pipeline
    mu sync.RWMutex
}

// Singleton pattern
var (
    etlRegistry *ETLComponentRegistry
    etlOnce sync.Once
)
```

### 1.4 Interface Definitions
```go
// app/etl/types/interfaces.go
type DataSource interface {
    Extract(ctx context.Context) ([]byte, error)
    GetConfig() interface{}
}

type Transformer interface {
    Transform(ctx context.Context, input []byte) ([]byte, error)
    GetConfig() interface{}
}

type Destination interface {
    Load(ctx context.Context, data []byte) error
    GetConfig() interface{}
}

type Pipeline interface {
    Execute(ctx context.Context) error
    GetConfig() interface{}
}
```

## Phase 2: Components cơ bản (3-4 ngày)

### 2.1 DataSource Implementation
```go
// app/etl/datasource/rest_api.go
type RestAPISource struct {
    config RestAPIConfig
    client *http.Client
}

type RestAPIConfig struct {
    URL     string            `yaml:"url"`
    Method  string            `yaml:"method"`
    Headers map[string]string `yaml:"headers"`
}
```

### 2.2 Transformer Implementation
```go
// app/etl/transformer/field_mapper.go
type FieldMapper struct {
    config FieldMapConfig
}

type FieldMapConfig struct {
    Mappings []FieldMapping `yaml:"mappings"`
}

type FieldMapping struct {
    Source string `yaml:"source"`
    Target string `yaml:"target"`
}
```

### 2.3 Destination Implementation
```go
// app/etl/destination/internal_api.go
type InternalAPI struct {
    config InternalAPIConfig
    client *http.Client
}

type InternalAPIConfig struct {
    Endpoint string            `yaml:"endpoint"`
    Method   string            `yaml:"method"`
    Headers  map[string]string `yaml:"headers"`
}
```

## Phase 3: Configuration (2-3 ngày)

### 3.1 Config Structure
```yaml
# config/etl/datasources/rest_api.yaml
datasources:
  - id: "user_api"
    type: "rest-api"
    config:
      url: "https://api.example.com/users"
      method: "GET"
      headers:
        Authorization: "${AUTH_TOKEN}"

# config/etl/transformers/user_mapping.yaml
transformers:
  - id: "user_transform"
    type: "field-mapping"
    config:
      mappings:
        - source: "data.id"
          target: "userId"
        - source: "data.attributes.name"
          target: "fullName"

# config/etl/pipelines/user_sync.yaml
pipelines:
  - id: "user_sync"
    steps:
      - type: "extract"
        source: "user_api"
      - type: "transform"
        transformer: "user_transform"
      - type: "load"
        destination: "internal_users"
```

## Phase 4: Testing & Documentation (2-3 ngày)

### 4.1 Unit Tests
- Registry tests
- Component tests
- Config loading tests

### 4.2 Integration Tests
- End-to-end pipeline tests
- Config-driven tests

### 4.3 Documentation
- API documentation
- Configuration guide
- Example pipelines

## Timeline
1. Phase 1: Ngày 1-3
2. Phase 2: Ngày 4-7
3. Phase 3: Ngày 8-10
4. Phase 4: Ngày 11-13

## Dependencies
1. Registry pattern hiện có
2. HTTP client
3. YAML parser
4. JSON parser

## Milestones
1. Milestone 1: Cấu trúc cơ bản hoàn thành (Cuối Phase 1)
2. Milestone 2: Components cơ bản hoạt động (Cuối Phase 2)
3. Milestone 3: Config system hoạt động (Cuối Phase 3)
4. Milestone 4: Hệ thống hoàn chỉnh với tests và docs (Cuối Phase 4) 