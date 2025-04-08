# ETL Pipeline Development Tasks

## Phase 1: Cấu trúc cơ bản
- [ ] Tạo cấu trúc thư mục
  - [ ] Tạo thư mục app/etl và các subfolders
  - [ ] Tạo thư mục config/etl

- [ ] Registry Implementation
  - [ ] Tạo file app/registry/etl.go
  - [ ] Implement ETLComponentRegistry
  - [ ] Bổ sung InitETLComponents vào app/registry/init.go

- [ ] Interface Definitions
  - [ ] DataSource interface và implementations
  - [ ] Transformer interface và implementations
  - [ ] Destination interface và implementations
  - [ ] Pipeline interface và implementations

## Phase 2: Components cơ bản
- [ ] DataSource
  - [ ] REST API connector
  - [ ] Factory pattern cho DataSource
  - [ ] Config parser

- [ ] Transformer
  - [ ] Field mapping transformer
  - [ ] Factory pattern cho Transformer
  - [ ] Config parser

- [ ] Destination
  - [ ] Internal API connector
  - [ ] Factory pattern cho Destination
  - [ ] Config parser

- [ ] Pipeline
  - [ ] Basic workflow engine
  - [ ] Factory pattern cho Pipeline
  - [ ] Config parser

## Phase 3: Configuration
- [ ] Config Files
  - [ ] datasources.yaml template
  - [ ] transformers.yaml template
  - [ ] destinations.yaml template
  - [ ] pipelines.yaml template

- [ ] Config Loading
  - [ ] YAML parser utilities
  - [ ] Environment variable substitution
  - [ ] Config validation

## Phase 4: Testing & Documentation
- [ ] Unit Tests
  - [ ] Registry tests
  - [ ] Component tests
  - [ ] Config loading tests

- [ ] Integration Tests
  - [ ] End-to-end pipeline tests
  - [ ] Config-driven tests

- [ ] Documentation
  - [ ] API documentation
  - [ ] Configuration guide
  - [ ] Example pipelines 