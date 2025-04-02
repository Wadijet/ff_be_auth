# ETL Pipeline

## Cấu trúc Metadata

### 1. Data Sources (metadata/etl/sources.yaml)
```yaml
sources:
  - name: external_users
    type: api
    config:
      url: ${EXTERNAL_API_URL}
      method: GET
      headers:
        Authorization: ${API_KEY}
    schedule: "0 */6 * * *"  # 6 tiếng/lần

  - name: legacy_auth_db
    type: database
    config:
      type: mysql
      connection_string: ${LEGACY_DB_URL}
    tables: [users, roles, permissions]

  - name: user_activity_logs
    type: file
    config:
      path: /var/log/auth/*.log
      format: json
    watch: true  # Real-time monitoring
```

### 2. Transformations (metadata/etl/transforms.yaml)
```yaml
transformations:
  - name: user_mapping
    source: external_users
    rules:
      - field: email
        action: lowercase
      - field: role
        action: map
        mapping:
          admin: system_admin
          user: basic_user
      - field: status
        action: default
        value: active

  - name: legacy_user_transform
    source: legacy_auth_db
    rules:
      - field: password
        action: hash
        algorithm: bcrypt
      - field: permissions
        action: convert_to_roles
        mapping_file: role_mappings.yaml
```

### 3. Loading Rules (metadata/etl/loading.yaml)
```yaml
loading:
  - target: users
    source: external_users
    mode: upsert
    key_fields: [email]
    conflict_resolution:
      - field: last_updated
        strategy: newest_wins
      - field: status
        strategy: preserve_active

  - target: audit_logs
    source: user_activity_logs
    mode: append
    batch_size: 1000
    retention:
      duration: 90d
      archive: true
```

### 4. Pipeline Workflows (metadata/etl/workflows.yaml)
```yaml
workflows:
  - name: user_sync
    schedule: "0 */6 * * *"
    steps:
      - extract:
          source: external_users
      - transform:
          use: user_mapping
      - validate:
          schema: user_schema
      - load:
          target: users
    error_handling:
      retry_count: 3
      notify: ["slack", "email"]

  - name: legacy_migration
    trigger: manual
    steps:
      - extract:
          source: legacy_auth_db
      - transform:
          use: legacy_user_transform
      - validate:
          schema: user_schema
      - load:
          target: users
    monitoring:
      metrics: ["processed_records", "error_rate"]
      alerts:
        - condition: "error_rate > 5%"
          severity: high
```

## Tính năng

### 1. Data Collection
- API integration
- Database connections
- File monitoring
- Real-time streaming
- Scheduled fetching
- Incremental loading

### 2. Data Processing
- Field mapping
- Data transformation
- Validation rules
- Error handling
- Data cleansing
- Format conversion

### 3. Data Loading
- Upsert/Append strategies
- Conflict resolution
- Batch processing
- Transaction management
- Data versioning
- Rollback support

### 4. Workflow Management
- Scheduled execution
- Manual triggers
- Dependencies handling
- Error recovery
- Progress tracking
- Resource management

### 5. Monitoring & Logging
- Progress tracking
- Error reporting
- Performance metrics
- Audit trails
- Health checks
- Alert notifications

### 6. Security
- Credential management
- Data encryption
- Access control
- Compliance tracking
- Data masking
- Audit logging

## Use Cases

### 1. User Synchronization
- Sync với external systems
- Migrate từ legacy systems
- Consolidate user data
- Clean up duplicates
- Update user metadata
- Sync permissions

### 2. Audit & Compliance
- Collect activity logs
- Generate compliance reports
- Track security events
- Archive historical data
- Monitor access patterns
- Detect anomalies

### 3. Data Enrichment
- Add user metadata
- Update role mappings
- Enhance security info
- Maintain data quality
- Add derived fields
- Calculate metrics

## Best Practices

### 1. Data Quality
- Validate input data
- Handle missing values
- Check data types
- Enforce constraints
- Monitor quality metrics
- Regular cleanup

### 2. Performance
- Batch processing
- Incremental updates
- Resource management
- Connection pooling
- Query optimization
- Caching strategies

### 3. Error Handling
- Retry mechanisms
- Error logging
- Alert notifications
- Fallback strategies
- Data recovery
- Transaction management

### 4. Security
- Encrypt sensitive data
- Secure connections
- Access control
- Audit logging
- Data masking
- Compliance checks 