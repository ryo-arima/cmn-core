# High-Level Architecture

## Overview

cmn-core delegates all identity management to external IdPs. The server validates JWT tokens issued by Keycloak or Casdoor and enforces role-based access at the API tier.

```
Browser / CLI
    │
    ▼
[Keycloak :8080] or [Casdoor :9000]   ← authentication & user/group management
    │  JWT token
    ▼
cmn-core server :8000
    ├── /v1/share/     (JWT required)
    ├── /v1/internal/  (JWT + role check)
    └── /v1/private/   (JWT + admin role required)
    │
    ├── PostgreSQL :5432  (pg_groups: Casdoor UUID↔name mapping)
    └── Redis :6379       (token cache / denylist)
```

## System Layers

### 1. API Layer
Gin-based HTTP router. Validates JWT tokens issued by the external IdP.

### 2. Controller Layer
Handlers organized by access tier (Share / Internal / Private).

### 3. Business Logic Layer
Use-case implementations. User/group operations are proxied to the configured IdP.

### 4. Repository Layer
- **IdPManager**: Abstracts Keycloak and Casdoor APIs (users, groups, members, roles).
- **Group (pg_groups)**: Casdoor-only store mapping generated UUIDs to display names.

### 5. Data Layer
- **PostgreSQL**: `pg_groups` table (Casdoor group UUID ↔ display name).
- **Redis**: Token cache, session management.

## Authentication Flow

1. User authenticates with Keycloak or Casdoor
2. IdP redirects to `GET /v1/share/auth/oidc/callback` with authorization code
3. Server exchanges code for JWT, validates signature against IdP's JWKS
4. JWT claims (including role) are used for authorization on subsequent requests
5. Response flows back through the layers

## Casdoor Group UUID Management

Casdoor does not assign native UUIDs to groups. cmn-core generates a UUID at creation time and stores it as the Casdoor group `name`. The human-readable display name is stored separately in the `pg_groups` PostgreSQL table.

```
Casdoor group.name  = UUID (e.g. e4e83e99-0e8a-5d7b-9b0b-1fd02db4dad2)
pg_groups.uuid      = same UUID
pg_groups.name      = display name (e.g. group001)
```

The usecase layer (`pkg/server/usecase/group.go`) translates between UUIDs and display names transparently.

## Configuration

The system is configured through `server.yaml` (or `app.yaml`). Key sections:

- **PostgreSQL** / **Redis**: Connection settings
- **Application.Server.IdP**: `provider: "casdoor"` or `"keycloak"` and credentials
- **Application.Server.Auth.OIDC**: Issuer URL and JWKS settings for JWT validation

## Security Features

- **JWT validation**: Signature verified against IdP's JWKS endpoint; `iss`, `aud`, `exp` checked
- **Token denylist**: Revoked tokens stored in Redis
- **Role enforcement**: Admin email list checked for `/v1/private/` endpoints
- **Multi-tier access**: Public / Share / Internal / Private endpoint separation

## Next Steps

- Learn more about individual [Components](./components.md)
- Explore the [API Reference](../api/overview.md)
- Review [Configuration Guide](../configuration/guide.md)
