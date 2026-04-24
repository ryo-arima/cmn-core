# Swagger Documentation

cmn-core provides comprehensive API documentation via Swagger/OpenAPI specification.

## Accessing Swagger UI

### GitHub Pages (Recommended)

View the interactive Swagger UI online:

**[📖 Open Swagger UI](https://ryo-arima.github.io/cmn-core/swagger/index.html)** (Opens in new window)

Or use the relative link when browsing documentation:
**[📖 Open Swagger UI](../../swagger/index.html)**

### Local Development

When running cmn-core locally, Swagger UI is available at:

```
http://localhost:8080/swagger/index.html
```

### Swagger Files

The Swagger specification is maintained in:
- **Source**: `docs/swagger/swagger.yaml`
- **Online**: [swagger.yaml](../../swagger/swagger.yaml)
- **Generated**: Auto-generated from code annotations

## Swagger Specification Overview

The Swagger specification includes:

### API Information

- **Title**: cmn-core API
- **Version**: 1.0.0
- **Base Path**: `/v1`
- **Schemes**: HTTP, HTTPS
- **Consumes**: `application/json`
- **Produces**: `application/json`

### Endpoints Documented

#### Public Endpoints

- `POST /v1/public/users/register` - User registration
- `POST /v1/public/users/login` - User authentication
- `GET /v1/public/health` - Health check

#### Internal Endpoints

- `GET /v1/internal/users` - List users
- `GET /v1/internal/users/count` - Count users
- `GET /v1/internal/users/{id}` - Get user details
- `GET /v1/internal/groups` - List groups
- `GET /v1/internal/groups/count` - Count groups
- `GET /v1/internal/members` - List members
- `GET /v1/internal/roles` - List roles

#### Private Endpoints

- `PUT /v1/private/users/{id}` - Update user
- `DELETE /v1/private/users/{id}` - Delete user
- `POST /v1/private/groups` - Create group
- `PUT /v1/private/groups/{id}` - Update group
- `DELETE /v1/private/groups/{id}` - Delete group
- `POST /v1/private/members` - Add member
- `DELETE /v1/private/members/{id}` - Remove member
- `POST /v1/private/roles` - Create role
- `PUT /v1/private/roles/{id}` - Update role
- `DELETE /v1/private/roles/{id}` - Delete role

### Data Models

The specification includes complete schemas for:

#### Request Models

- `UserRequest` - User registration/update
- `GroupRequest` - Group creation/update
- `MemberRequest` - Member addition
- `RoleRequest` - Role creation/update

#### Response Models

- `User` - User object
- `Group` - Group object
- `Member` - Member object
- `Role` - Role object
- `Count` - Count response

#### Common Properties

All models include:
- `id` - Unique identifier (uint64)
- `uuid` - UUID string
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp
- `deleted_at` - Soft delete timestamp (nullable)

## Using Swagger UI

### Try Out API Endpoints

1. Open Swagger UI in browser
2. Click on an endpoint to expand
3. Click "Try it out"
4. Fill in parameters
5. Click "Execute"
6. View response

### Authentication

For protected endpoints:

1. Get JWT token via `/v1/public/users/login`
2. Click "Authorize" button in Swagger UI
3. Enter token in format: `Bearer <token>`
4. Click "Authorize"
5. Try protected endpoints

## Generating Swagger Documentation

### From Code Annotations

cmn-core uses `swaggo/swag` for generating Swagger docs from code comments:

```bash
# Install swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g cmd/server/main.go -o docs/swagger

# Documentation will be generated in docs/swagger/
```

### Swagger Annotations Example

```go
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body entity.UserRequest true "User registration request"
// @Success 201 {object} entity.User
// @Failure 400 {object} map[string]string
// @Router /v1/public/users/register [post]
func (c *userController) Register(ctx *gin.Context) {
    // Implementation
}
```

## Swagger vs OpenAPI

- **OpenAPI 3.0**: Modern specification format
- **Swagger 2.0**: Legacy format (currently used by cmn-core)
- **Migration**: Consider upgrading to OpenAPI 3.0

## Interactive Documentation

### Swagger Editor

Edit the specification online:
1. Visit [Swagger Editor](https://editor.swagger.io/)
2. Paste contents of `docs/swagger/swagger.yaml`
3. View formatted documentation
4. Validate specification

### Postman Integration

Import Swagger spec into Postman:
1. Open Postman
2. Click "Import"
3. Select `docs/swagger/swagger.yaml`
4. Use generated collection for testing

## API Testing with Swagger

Swagger UI provides built-in testing:

1. **Authentication**: Test login flow
2. **CRUD Operations**: Test all endpoints
3. **Validation**: Verify request/response schemas
4. **Error Handling**: Test error scenarios

## Best Practices

### Maintaining Documentation

- Update Swagger annotations when changing APIs
- Regenerate docs after code changes
- Validate specification regularly
- Keep examples up-to-date

### Documentation Quality

- Clear descriptions for all endpoints
- Comprehensive parameter documentation
- Example request/response bodies
- Document all error codes

## Swagger Tools

### CLI Tools

```bash
# Validate Swagger spec
swagger validate docs/swagger/swagger.yaml

# Generate client SDKs
swagger generate client -f docs/swagger/swagger.yaml

# Generate server code
swagger generate server -f docs/swagger/swagger.yaml
```

### Online Validators

- [Swagger Validator](https://validator.swagger.io/)
- [OpenAPI Validator](https://apitools.dev/swagger-parser/online/)

## Next Steps

- [API Overview](../api/overview.md)
- [Authentication](../api/authentication.md)
- [Endpoint Reference](../api/endpoints.md)
- [GoDoc](./godoc.md)
