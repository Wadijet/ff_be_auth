# Auth Metadata Schema

## User Schema
```yaml
user:
  type: object
  required: [email, password]
  properties:
    email:
      type: string
      format: email
    password:
      type: string
      minLength: 8
    status:
      type: string
      enum: [active, inactive, pending]
    roles:
      type: array
      items:
        type: string
    metadata:
      type: object
      additionalProperties: true
```

## Role Schema
```yaml
role:
  type: object
  required: [name]
  properties:
    name:
      type: string
    permissions:
      type: array
      items:
        type: string
    metadata:
      type: object
      additionalProperties: true
```

## Permission Schema
```yaml
permission:
  type: object
  required: [resource, action]
  properties:
    resource:
      type: string
    action:
      type: string
      enum: [create, read, update, delete]
    conditions:
      type: object
      additionalProperties: true
```

## Authentication Flow
```yaml
auth_flow:
  type: object
  required: [type]
  properties:
    type:
      type: string
      enum: [basic, oauth, mfa]
    config:
      type: object
      additionalProperties: true
    steps:
      type: array
      items:
        type: object
        required: [action]
        properties:
          action:
            type: string
          config:
            type: object
```

## Security Policy
```yaml
security_policy:
  type: object
  properties:
    password_policy:
      type: object
      properties:
        min_length:
          type: integer
          minimum: 8
        require_uppercase:
          type: boolean
        require_numbers:
          type: boolean
        require_special:
          type: boolean
    rate_limit:
      type: object
      properties:
        requests_per_minute:
          type: integer
        burst:
          type: integer
    session:
      type: object
      properties:
        timeout:
          type: integer
        max_sessions:
          type: integer
```

## OAuth Provider
```yaml
oauth_provider:
  type: object
  required: [name, client_id]
  properties:
    name:
      type: string
    client_id:
      type: string
    client_secret:
      type: string
    redirect_uri:
      type: string
    scopes:
      type: array
      items:
        type: string
```

## MFA Configuration
```yaml
mfa_config:
  type: object
  properties:
    providers:
      type: array
      items:
        type: object
        required: [type]
        properties:
          type:
            type: string
            enum: [totp, sms, email]
          config:
            type: object
    required_roles:
      type: array
      items:
        type: string
```

## Audit Configuration
```yaml
audit_config:
  type: object
  properties:
    events:
      type: array
      items:
        type: object
        required: [type]
        properties:
          type:
            type: string
          level:
            type: string
            enum: [info, warn, error]
    retention:
      type: object
      properties:
        duration:
          type: string
        max_size:
          type: string
``` 