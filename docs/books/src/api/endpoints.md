# Endpoints

## Share Endpoints (JWT required)

| Method | Path | Description |
|---|---|---|
| GET | `/v1/share/auth/oidc/callback` | OIDC callback |

## Internal Endpoints (JWT + role check)

Role-gated endpoints. Roles are provided by the external IdP in the JWT claims.

## Private Endpoints (JWT required)

Authenticated endpoints without additional role enforcement.

## Health Check

```
GET /health
```

For the full interactive reference, see the [Swagger Documentation](../appendix/swagger.md).
