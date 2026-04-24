# API Overview

cmn-core provides a RESTful API built with the Gin framework. The API is organized into three access levels to enforce security boundaries.

## Base URL

```
http://localhost:8080/v1
```

## Access Levels

### Public Endpoints (`/v1/public`)

Public endpoints are accessible without authentication:

- **User Registration**: Create new user accounts
- **User Login**: Authenticate and receive JWT token
- **Health Check**: Service status

### Internal Endpoints (`/v1/internal`)

Internal endpoints require JWT authentication and are intended for internal services:

- **User Management**: List, count, get user details
- **Group Management**: CRUD operations for groups
- **Member Management**: Manage group memberships
- **Role Management**: Query roles and permissions

### Private Endpoints (`/v1/private`)

Private endpoints require both JWT authentication and specific Casbin permissions:

- **User Administration**: Update, delete users
- **Group Administration**: Full group management
- **Member Administration**: Full membership control
- **Role Administration**: Create, update, delete roles

## Authentication

### JWT Token

All authenticated endpoints require a JWT token in the Authorization header:

```http
Authorization: Bearer <jwt_token>
```

### Token Lifecycle

1. **Obtain Token**: POST to `/v1/public/users/login` with credentials
2. **Use Token**: Include in Authorization header for subsequent requests
3. **Token Expiry**: Tokens expire after configured duration (default: 24 hours)
4. **Token Revocation**: Logout endpoint adds token to Redis denylist

## Authorization

Authorization is handled by Casbin with two policy sets:

### App Policies (`etc/casbin/cmn/`)

Controls access to API endpoints based on user roles.

Example policy:
```csv
p, admin, /v1/private/*, *
p, user, /v1/internal/users, GET
```

### Resource Policies (`etc/casbin/resources/`)

Controls access to specific resources (groups, members).

Example policy:
```csv
p, user:123, group:456, read
p, user:123, group:456, write
```

## Request/Response Format

### Request Format

```json
{
  "name": "example",
  "email": "user@example.com"
}
```

### Response Format

```json
{
  "data": {
    "id": 1,
    "uuid": "f3b3b3b3-3b3b-3b3b-3b3b-3b3b3b3b3b3b",
    "name": "example"
  },
  "message": "Success"
}
```

### Error Format

```json
{
  "error": "Invalid credentials",
  "code": 401
}
```

## Rate Limiting

Currently, rate limiting is not implemented but can be added via middleware.

## API Versioning

The API is versioned in the URL path (`/v1/`). Future versions will be released as `/v2/`, etc.

## Next Steps

- [Authentication Details](./authentication.md)
- [Endpoint Reference](./endpoints.md)
- [Swagger Documentation](../appendix/swagger.md)
