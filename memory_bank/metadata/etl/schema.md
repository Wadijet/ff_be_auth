# ETL Metadata Schema

## Data Source Schema
```yaml
data_source:
  type: object
  required: [name, type]
  properties:
    name:
      type: string
    type:
      type: string
      enum: [api, database, file]
    config:
      type: object
      properties:
        api:
          type: object
          properties:
            url:
              type: string
            method:
              type: string
            headers:
              type: object
        database:
          type: object
          properties:
            type:
              type: string
            connection_string:
              type: string
            tables:
              type: array
              items:
                type: string
        file:
          type: object
          properties:
            path:
              type: string
            format:
              type: string
            watch:
              type: boolean
    schedule:
      type: string
```

## Transformation Schema
```yaml
transformation:
  type: object
  required: [name, source]
  properties:
    name:
      type: string
    source:
      type: string
    rules:
      type: array
      items:
        type: object
        required: [field, action]
        properties:
          field:
            type: string
          action:
            type: string
            enum: [lowercase, uppercase, map, default, hash]
          mapping:
            type: object
          value:
            type: string
```

## Loading Schema
```yaml
loading:
  type: object
  required: [target, source]
  properties:
    target:
      type: string
    source:
      type: string
    mode:
      type: string
      enum: [insert, update, upsert, append]
    key_fields:
      type: array
      items:
        type: string
    conflict_resolution:
      type: array
      items:
        type: object
        properties:
          field:
            type: string
          strategy:
            type: string
    batch_size:
      type: integer
    retention:
      type: object
      properties:
        duration:
          type: string
        archive:
          type: boolean
```

## Workflow Schema
```yaml
workflow:
  type: object
  required: [name]
  properties:
    name:
      type: string
    schedule:
      type: string
    trigger:
      type: string
      enum: [scheduled, manual, event]
    steps:
      type: array
      items:
        type: object
        properties:
          extract:
            type: object
            properties:
              source:
                type: string
          transform:
            type: object
            properties:
              use:
                type: string
          validate:
            type: object
            properties:
              schema:
                type: string
          load:
            type: object
            properties:
              target:
                type: string
    error_handling:
      type: object
      properties:
        retry_count:
          type: integer
        notify:
          type: array
          items:
            type: string
    monitoring:
      type: object
      properties:
        metrics:
          type: array
          items:
            type: string
        alerts:
          type: array
          items:
            type: object
            properties:
              condition:
                type: string
              severity:
                type: string
```

## Validation Schema
```yaml
validation:
  type: object
  required: [name]
  properties:
    name:
      type: string
    rules:
      type: array
      items:
        type: object
        required: [field, type]
        properties:
          field:
            type: string
          type:
            type: string
          required:
            type: boolean
          pattern:
            type: string
          min:
            type: number
          max:
            type: number
          enum:
            type: array
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
          source:
            type: string
    alerts:
      type: array
      items:
        type: object
        properties:
          name:
            type: string
          condition:
            type: string
          channels:
            type: array
            items:
              type: string
    logging:
      type: object
      properties:
        level:
          type: string
        retention:
          type: string
``` 