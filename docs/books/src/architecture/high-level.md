# High-Level Architecture

## Overview

cmn-core delegates all identity management to external IdPs. The server validates JWT tokens issued by Keycloak or Casdoor and enforces role-based access at the API tier.

```
Browser / CLI
    │
    ▼
[Keycloak :8080] or [Casdoor :9000]   ← authentication & user management
    │  JWT token
    ▼
cmn-core server :8000
    ├── /v1/share/     (JWT required)
    ├── /v1/internal/  (JWT + role check)
    └── /v1/private/   (JWT required)
    │
    ├── PostgreSQL :5432
    └── Redis :6379
```

## System Layers

### 1. API Layer
Gin-based HTTP router. Validates JWT tokens issued by the external IdP.

### 2. Controller Layer
Handlers organized by access tier (Share / Internal / Private).

### 3. Business Logic Layer
Use-case implementations. No user/group/role CRUD — delegated to IdPs.

### 4. Repository Layer
Abstracts PostgreSQL and Redis access.

### 5. Data Layer
- **PostgreSQL**: Application data
- **Redis**: Token cache, session management

## Authentication Flow

1. User authenticates with Keycloak or Casdoor
2. IdP redirects to `GET /v1/share/auth/oidc/callback` with authorization code
3. Server exchanges code for JWT, validates signature against IdP's JWKS
4. JWT claims (including role) are used for authorization on subsequent requests

6. Repository performs database operations
7. Response flows back through the layers

## Configuration

The system is configured through:

- **app.yaml**: Main configuration file (database, Redis, JWT settings)
- **Casbin Model Files**: Define RBAC model structure
- **Casbin Policy Files**: Define actual permissions

## Scalability Considerations

- **Stateless API**: JWT tokens enable horizontal scaling
- **Redis Caching**: Reduces database load
- **Connection Pooling**: Efficient database connection management
- **Casbin Policy Caching**: In-memory policy evaluation for fast authorization

## Security Features

- **JWT with HS256**: Secure token-based authentication
- **Token Denylist**: Revoked tokens stored in Redis
- **Casbin RBAC**: Policy-based authorization
- **Multi-tier Access**: Public, Internal, and Private endpoint separation
- **Admin Email Check**: Additional verification for administrative operations

## Next Steps

- Learn more about individual [Components](./components.md)
- Explore the [API Reference](../api/overview.md)
- Review [Configuration Guide](../configuration/guide.md)
