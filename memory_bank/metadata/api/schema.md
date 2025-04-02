# API Metadata Schema

## Endpoint Schema
```yaml
endpoint:
  type: object
  required: [path, method]
  properties:
    path:
      type: string
    method:
      type: string
      enum: [GET, POST, PUT, DELETE, PATCH]
    tags:
      type: array
      items:
        type: string
    summary:
      type: string
    description:
      type: string
    parameters:
      type: array
      items:
        type: object
        required: [name, in]
        properties:
          name:
            type: string
          in:
            type: string
            enum: [path, query, header, cookie]
          required:
            type: boolean
          schema:
            type: object
    requestBody:
      type: object
      properties:
        required:
          type: boolean
        content:
          type: object
    responses:
      type: object
      additionalProperties:
        type: object
        properties:
          description:
            type: string
          content:
            type: object
    security:
      type: array
      items:
        type: object
```

## Security Schema
```yaml
security:
  type: object
  properties:
    schemes:
      type: object
      properties:
        jwt:
          type: object
          properties:
            type:
              type: string
              enum: [http, apiKey]
            scheme:
              type: string
            bearerFormat:
              type: string
        oauth2:
          type: object
          properties:
            type:
              type: string
            flows:
              type: object
```

## Validation Schema
```yaml
validation:
  type: object
  properties:
    types:
      type: object
      additionalProperties:
        type: object
        properties:
          type:
            type: string
          properties:
            type: object
          required:
            type: array
            items:
              type: string
    formats:
      type: object
      additionalProperties:
        type: object
        properties:
          pattern:
            type: string
          validate:
            type: string
```

## Rate Limit Schema
```yaml
rate_limit:
  type: object
  properties:
    default:
      type: object
      properties:
        requests:
          type: integer
        window:
          type: string
    paths:
      type: object
      additionalProperties:
        type: object
        properties:
          requests:
            type: integer
          window:
            type: string
    ips:
      type: object
      additionalProperties:
        type: object
        properties:
          requests:
            type: integer
          window:
            type: string
```

## Documentation Schema
```yaml
documentation:
  type: object
  properties:
    info:
      type: object
      properties:
        title:
          type: string
        version:
          type: string
        description:
          type: string
    servers:
      type: array
      items:
        type: object
        properties:
          url:
            type: string
          description:
            type: string
    tags:
      type: array
      items:
        type: object
        properties:
          name:
            type: string
          description:
            type: string
    externalDocs:
      type: object
      properties:
        url:
          type: string
        description:
          type: string
```

## Monitoring Schema
```yaml
monitoring:
  type: object
  properties:
    metrics:
      type: array
      items:
        type: object
        properties:
          name:
            type: string
          type:
            type: string
          labels:
            type: array
            items:
              type: string
    tracing:
      type: object
      properties:
        enabled:
          type: boolean
        sampler:
          type: object
    logging:
      type: object
      properties:
        level:
          type: string
        format:
          type: string
``` 