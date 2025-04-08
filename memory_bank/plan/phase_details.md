# Chi tiết các Phase phát triển

## Phase 0: Registry Refactoring

### Tasks
1. Generic Registry Implementation
   ```go
   // app/registry/registry.go
   type Registry[T any] struct {
       items map[string]T
       mu sync.RWMutex
   }
   ```
   - Implement các methods cơ bản
   - Viết unit tests
   - Đảm bảo thread safety

2. Collection Registry Migration
   - Chuyển đổi code hiện tại sang generic registry
   - Cập nhật các hàm helper
   - Đảm bảo backward compatibility

3. Testing & Validation
   - Unit tests cho generic registry
   - Integration tests với code hiện có
   - Performance testing
   - Validation các use cases hiện tại

### Dependencies
- Go version >= 1.18 (for generics)
- Testing framework
- Existing codebase

### Deliverables
- [x] Generic registry implementation
- [ ] Migrated collection registry
- [ ] Test coverage
- [ ] Documentation update

## Phase 1: Cấu trúc cơ bản

### Ngày 1: Setup & Registry
1. Tạo cấu trúc thư mục
   ```
   mkdir -p app/etl/{types,datasource,transformer,destination,pipeline}
   mkdir -p config/etl/{datasources,transformers,destinations,pipelines}
   ```

2. Tạo file registry
   ```go
   // app/registry/etl.go
   func init() {
       // Initialize ETL registry
   }
   ```

3. Định nghĩa interfaces
   ```go
   // app/etl/types/interfaces.go
   type DataSource interface { ... }
   ```

### Ngày 2-3: Components & Tests
1. Implement base types
2. Viết unit tests
3. Setup CI pipeline

## Phase 2: Components cơ bản

### Ngày 4-5: DataSource & Transformer
1. REST API source
   - HTTP client wrapper
   - Config parsing
   - Error handling

2. Field mapper
   - JSON path parsing
   - Mapping logic
   - Validation

### Ngày 6-7: Destination & Pipeline
1. Internal API destination
   - HTTP client
   - Retry logic
   - Error handling

2. Basic pipeline
   - Step execution
   - Error propagation
   - State management

## Phase 3: Configuration System

### Ngày 8-9: Config Loading
1. YAML parser
   - File loading
   - Environment variables
   - Validation

2. Config management
   - Hot reload
   - Caching
   - Versioning

### Ngày 9-10: Integration
1. Component factory
   - Dynamic creation
   - Type registry
   - Validation

2. Pipeline builder
   - Config parsing
   - Step validation
   - Error handling

## Phase 4: Testing & Documentation

### Ngày 11-12: Testing
1. Unit tests
   - Component tests
   - Config tests
   - Error cases

2. Integration tests
   - End-to-end flows
   - Performance tests
   - Edge cases

### Ngày 13: Documentation
1. API docs
   - Interface documentation
   - Example usage
   - Best practices

2. Config guide
   - Format specification
   - Examples
   - Troubleshooting

## Dependencies cho mỗi Phase

### Phase 1
- Go modules setup
- Registry pattern
- Testing framework

### Phase 2
- HTTP client
- JSON parser
- Context package

### Phase 3
- YAML parser
- Environment config
- Validation library

### Phase 4
- Testing tools
- Documentation generator
- CI/CD pipeline

## Milestones & Deliverables

### Milestone 1: Cấu trúc cơ bản
- [x] Cấu trúc thư mục
- [ ] Registry implementation
- [ ] Interface definitions
- [ ] Basic tests

### Milestone 2: Components
- [ ] REST API source
- [ ] Field mapper
- [ ] Internal API destination
- [ ] Basic pipeline

### Milestone 3: Configuration
- [ ] YAML parsing
- [ ] Config management
- [ ] Component factory
- [ ] Pipeline builder

### Milestone 4: Quality
- [ ] Unit tests
- [ ] Integration tests
- [ ] Documentation
- [ ] Examples 