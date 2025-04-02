# Cấu trúc Metadata

## Tổng quan
Metadata được tổ chức thành các module riêng biệt, mỗi module quản lý một khía cạnh của hệ thống authentication. Tất cả được định nghĩa bằng YAML để dễ đọc và maintain.

## Cấu trúc thư mục
```
metadata/
├── auth/                  # Authentication configs
│   ├── methods/          # Auth methods
│   ├── providers/        # Identity providers
│   └── flows/           # Auth flows
├── users/                # User management
│   ├── schemas/         # User schemas
│   ├── roles/           # Role definitions
│   └── policies/        # Access policies
├── api/                  # API definitions
│   ├── endpoints/       # API endpoints
│   ├── validations/     # Request validation
│   └── responses/       # Response templates
└── system/              # System configs
    ├── database/        # DB connections
    ├── security/        # Security settings
    └── monitoring/      # Monitoring configs
```

## Chi tiết metadata

### 1. Authentication Methods
```yaml
auth_method:
  name: password
  type: basic
  enabled: true
  config:
    hash_algorithm: bcrypt
    hash_rounds: 12
    min_length: 8
    require_special: true
    require_number: true
  validation:
    - rule: length
      min: 8
      max: 32
    - rule: complexity
      special: true
      number: true
  error_messages:
    too_short: "Mật khẩu phải có ít nhất 8 ký tự"
    too_simple: "Mật khẩu phải chứa ký tự đặc biệt và số"
```

### 2. Role Definitions
```yaml
role:
  name: admin
  description: "Quản trị viên hệ thống"
  permissions:
    - resource: users
      actions: [create, read, update, delete]
      conditions:
        organization_id: ${user.org_id}
    - resource: roles
      actions: [read, assign]
  inheritance:
    - role: user
    - role: manager
  metadata:
    created_at: ${timestamp}
    updated_at: ${timestamp}
```

### 3. API Endpoints
```yaml
endpoint:
  path: /auth/login
  method: POST
  auth_required: false
  rate_limit:
    requests: 5
    window: 60
  validation:
    body:
      type: object
      properties:
        username:
          type: string
          format: email
        password:
          type: string
          min_length: 8
  workflow:
    - validate_input
    - check_rate_limit
    - authenticate_user
    - generate_tokens
    - audit_log
  responses:
    200:
      description: "Đăng nhập thành công"
      schema: LoginResponse
    401:
      description: "Thông tin không hợp lệ"
      schema: ErrorResponse
```

### 4. Security Policies
```yaml
security_policy:
  name: password_policy
  rules:
    - type: complexity
      min_length: 8
      require_uppercase: true
      require_lowercase: true
      require_number: true
      require_special: true
    - type: history
      remember_last: 5
      no_reuse_within_days: 90
    - type: expiry
      max_age_days: 90
      remind_before_days: 7
  lockout:
    max_attempts: 5
    window_minutes: 30
    lockout_minutes: 60
```

### 5. Monitoring Config
```yaml
monitoring:
  metrics:
    - name: login_attempts
      type: counter
      labels: [status, method, ip]
    - name: active_sessions
      type: gauge
      labels: [user_type, device]
  alerts:
    - name: high_failure_rate
      condition: "rate(login_attempts{status='failed'}[5m]) > 10"
      severity: warning
      channels: [slack, email]
  dashboards:
    - name: auth_overview
      panels:
        - title: "Login Success Rate"
          type: graph
          metric: login_attempts
        - title: "Active Sessions"
          type: gauge
          metric: active_sessions
```

### 6. Database Schema
```yaml
collection:
  name: users
  indexes:
    - fields: [email]
      unique: true
    - fields: [username]
      unique: true
    - fields: [created_at]
  schema:
    required: [email, username, password_hash]
    properties:
      email:
        type: string
        format: email
      username:
        type: string
        min_length: 3
      password_hash:
        type: string
      roles:
        type: array
        items:
          type: string
      metadata:
        type: object
        properties:
          last_login:
            type: datetime
          failed_attempts:
            type: integer
```

## Validation Rules

### 1. Schema Validation
- Kiểm tra cấu trúc metadata
- Validate required fields
- Type checking
- Format validation

### 2. Business Rules
- Dependency checking
- Conflict detection
- Security validation
- Performance impact

### 3. Version Control
- Semantic versioning
- Change tracking
- Rollback support
- Migration paths

## Best Practices

### 1. Tổ chức file
- Chia nhỏ theo module
- Sử dụng meaningful names
- Comment đầy đủ
- Version control

### 2. Validation
- Schema validation
- Business logic validation
- Security checks
- Performance impact

### 3. Maintenance
- Regular reviews
- Clean up unused configs
- Update documentation
- Monitor usage

### 4. Security
- Encrypt sensitive data
- Access control
- Audit logging
- Regular backups 