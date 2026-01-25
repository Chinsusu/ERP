---
description: Generate OpenAPI 3.0 specification for a service
---

# Generate OpenAPI Specification

## Usage
```
/generate-openapi [SERVICE_NAME]
```

## Steps

1. Read service documentation from `/opt/ERP/docs/[##-SERVICE-NAME].md`
2. Identify all API endpoints from the handlers
3. Create `openapi.yaml` in `services/[service-name]/docs/`

## OpenAPI Requirements

- Use OpenAPI 3.0.3 specification
- Include Bearer JWT authentication
- Include full request/response schemas
- Include error responses (400, 401, 403, 404, 500)
- Include examples for each endpoint

## Example Structure

```yaml
openapi: 3.0.3
info:
  title: Service Name API
  version: 1.0.0
  description: Description of the service

servers:
  - url: http://localhost:8080/api/v1
    description: Local development

security:
  - bearerAuth: []

paths:
  /resource:
    get:
      summary: List resources
      tags: [Resource]
      parameters:
        - $ref: '#/components/parameters/PageParam'
        - $ref: '#/components/parameters/LimitParam'
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ResourceListResponse'
        '401':
          $ref: '#/components/responses/Unauthorized'

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    Resource:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        created_at:
          type: string
          format: date-time

  responses:
    Unauthorized:
      description: Authentication required
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'

  parameters:
    PageParam:
      name: page
      in: query
      schema:
        type: integer
        default: 1
    LimitParam:
      name: limit
      in: query
      schema:
        type: integer
        default: 20
```

## Common Schemas
Include standard schemas: Error, Pagination, UUID, Timestamp formats.
