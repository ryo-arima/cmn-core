# Authentication

cmn-core delegates all authentication to external IdPs via **OIDC (OpenID Connect)**.

## Supported IdPs

| IdP | Port | Login URL |
|---|---|---|
| Keycloak | 8080 | http://localhost:8080/realms/cmn/account |
| Casdoor | 9000 | http://localhost:9000/login/cmn |

## OIDC Flow

1. Client redirects user to IdP login page
2. User authenticates with IdP
3. IdP redirects to `GET /v1/share/auth/oidc/callback` with authorization code
4. Server exchanges code for JWT token
5. JWT token is used for subsequent API requests

## Using the JWT Token

Include the token in the `Authorization` header:

```http
Authorization: Bearer <jwt_token>
```

## Development Users (cmn realm / cmn org)

| Username | Password | Role |
|---|---|---|
| user01 – user10 | `Password123!` | app |

## OIDC Client Settings

| Field | Value |
|---|---|
| Client ID | `cmn-core-client-id` |
| Client Secret | `cmn-core-client-secret` |
| Redirect URI | `http://localhost:8000/v1/share/auth/oidc/callback` |
