# Technical Context - ETL Pipeline

## Công nghệ hiện có
1. **Backend Framework**
   - Golang
   - MongoDB
   - Registry Pattern đã implement

2. **Cấu trúc hiện tại**
   ```
   app/
   ├── registry/
   │   ├── collection.go
   │   └── init.go
   ```

## Kiến trúc ETL đề xuất
1. **Cấu trúc thư mục**
   ```
   app/
   ├── registry/
   │   ├── collection.go
   │   ├── init.go
   │   └── etl.go
   ├── etl/
   │   ├── datasource/
   │   ├── transformer/
   │   ├── destination/
   │   └── pipeline/
   └── config/
       └── etl/
   ```

2. **Components**
   - **DataSource**
     ```go
     type DataSource interface {
         Extract(ctx context.Context) ([]byte, error)
         GetConfig() interface{}
     }
     ```
   
   - **Transformer**
     ```go
     type Transformer interface {
         Transform(ctx context.Context, input []byte) ([]byte, error)
         GetConfig() interface{}
     }
     ```
   
   - **Destination**
     ```go
     type Destination interface {
         Load(ctx context.Context, data []byte) error
         GetConfig() interface{}
     }
     ```

   - **Pipeline**
     ```go
     type Pipeline interface {
         Execute(ctx context.Context) error
         GetConfig() interface{}
     }
     ```

3. **Registry Pattern**
   ```go
   type ETLComponentRegistry struct {
       dataSources  map[string]DataSource
       transformers map[string]Transformer
       destinations map[string]Destination
       pipelines    map[string]Pipeline
   }
   ```

## Cấu trúc Config
1. **DataSource Config**
   ```yaml
   datasources:
     - id: "source1"
       type: "rest-api"
       config:
         url: "https://api.example.com/data"
         method: "GET"
         headers:
           Authorization: "${AUTH_TOKEN}"
   ```

2. **Transformer Config**
   ```yaml
   transformers:
     - id: "transform1"
       type: "field-mapping"
       config:
         mappings:
           - source: "data.id"
             target: "userId"
   ```

3. **Pipeline Config**
   ```yaml
   pipelines:
     - id: "sync-pipeline"
       steps:
         - type: "extract"
           source: "source1"
         - type: "transform"
           transformer: "transform1"
         - type: "load"
           destination: "dest1"
   ``` 