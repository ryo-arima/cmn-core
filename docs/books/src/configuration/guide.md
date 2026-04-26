# Configuration Guide

cmn-core uses a YAML configuration file (`etc/server.yaml`).

## Quick Start

```bash
cp etc/server.yaml.example etc/server.yaml
```

## Configuration Sections

### Server

```yaml
Application:
  Server:
    port: 8000
    admin:
      emails:
        - "admin@cmn.local"   # Emails with the admin role
    jwt_secret: "CHANGE_THIS_JWT_SECRET_IN_PRODUCTION"
    log_level: "info"            # debug / info / warn / error
    redis:
      jwt_cache: true
      cache_ttl: 1800            # seconds (0 = use token expiry)
```

### Identity Provider (user/group management)

```yaml
    idp:
      provider: "casdoor"        # or "keycloak"
      casdoor:
        base_url: "http://localhost:9000"
        client_id: "cmn-core-client-id"
        client_secret: "cmn-core-client-secret"
        organization: "cmn"
      # keycloak:
      #   base_url: "http://localhost:8080"
      #   realm: "cmn"
      #   admin_client_id: "admin-cli"
      #   admin_client_secret: "CHANGE_THIS_SECRET"
```

### OIDC (JWT validation)

```yaml
    auth:
      oidc:
        issuer_url: "http://localhost:9000"   # Must match the `iss` claim in issued JWTs
        provider_url: "http://localhost:9000" # OIDC discovery URL.
                                               # Set to the internal service URL when running
                                               # in Docker (e.g. http://casdoor:8000) if the
                                               # public issuer URL is not reachable from within
                                               # the container network.
        client_id: "cmn-core-client-id"
```

### PostgreSQL

```yaml
PostgreSQL:
  host: "localhost"
  user: "user"
  pass: "password"
  port: "5432"
  db: "cmn_core"
  sslmode: "disable"
```

### Redis

```yaml
Redis:
  host: "localhost"
  port: 6379
  user: "default"
  pass: ""
  db: 0
```

## Security Best Practices

- **Never commit `etc/server.yaml`** (it is in `.gitignore`)
- Use strong random values for `jwt_secret`:
  ```bash
  openssl rand -base64 32
  ```
- Use separate database credentials per environment

## Client Credential Files

CLI clients read credentials from separate YAML files:

```yaml
# etc/.cmn/client/credentials/<name>.yaml
Application:
  Client:
    ServerEndpoint: "http://localhost:8000"
    credentials:
      email: "admin@cmn.local"
      password: "Admin123!"
```

Clients send credentials to `POST /v1/public/login` on the server; the server obtains a JWT from the IdP and returns it. The token is cached locally and reused until it expires.

## Next Steps

- [Environment Variables](./environment.md)
- [Getting Started](../development/getting-started.md)

## Validation

Validate your configuration before starting:

```bash
# Test database connection
go run cmd/server/main.go --config etc/app.yaml --validate

# Check Casbin policies
casbin-cli check etc/casbin/cmn/model.conf etc/casbin/cmn/policy.csv
```

## Troubleshooting

### Connection Errors

```
Error: dial tcp: lookup mysql-host: no such host
```

**Solution**: Verify hostname and network connectivity

### JWT Errors

```
Error: jwt_secret must be at least 256 bits
```

**Solution**: Generate a longer secret:
```bash
openssl rand -base64 32
```

### Casbin Errors

```
Error: failed to load casbin policy
```

**Solution**: Check file paths and CSV format in policy files

## Next Steps

- [Environment Variables](./environment.md)
- [Casbin Policies](./casbin.md)
- [Getting Started](../development/getting-started.md)
