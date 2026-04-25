# API Overview

cmn-core provides a RESTful API built with Gin. Authentication is fully handled by external IdPs (Keycloak / Casdoor) via OIDC.

## Base URL

```
http://localhost:8000/v1
```

## Access Levels

| Tier | Path prefix | Auth required |
|---|---|---|
| Share | `/v1/share/` | JWT (from IdP) |
| Internal | `/v1/internal/` | JWT + role check |
| Private | `/v1/private/` | JWT |

## Authentication

All authenticated endpoints require the JWT token issued by the external IdP:

```http
Authorization: Bearer <jwt_token>
```

OIDC callback endpoint (used during login flow):

```
GET /v1/share/auth/oidc/callback
```

## Request / Response Format

```json
// Success
{ "data": { ... }, "message": "Success" }

// Error
{ "error": "Unauthorized", "code": 401 }
```

## API Versioning

The API is versioned in the URL path (`/v1/`). Future versions will be `/v2/`, etc.

## Next Steps

- [Authentication Details](./authentication.md)
- [Endpoint Reference](./endpoints.md)
- [Swagger Documentation](../appendix/swagger.md)
