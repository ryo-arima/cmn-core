# Configuration Guide

cmn-core uses a YAML configuration file (`etc/app.yaml`).

## Quick Start

```bash
cp etc/app.yaml.example etc/app.yaml
```

## Configuration Sections

### Server

```yaml
Application:
  Server:
    port: 8000
    jwt_secret: "CHANGE_THIS_JWT_SECRET_IN_PRODUCTION"
    log_level: "debug"
    redis:
      jwt_cache: true
      cache_ttl: 1800
    jwt:
      key: "CHANGE_THIS_JWT_KEY"
    auth:
      provider: "oidc"   # "oidc" or "saml"
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
  db: "0"
```

## Security Best Practices

- **Never commit `etc/app.yaml`** (it is in `.gitignore`)
- Use strong random values for `jwt_secret` and `jwt.key`:
  ```bash
  openssl rand -base64 32
  ```
- Use separate database credentials per environment

## Next Steps

- [Environment Variables](./environment.md)
- [Getting Started](../development/getting-started.md)


**Note**: The `db` field accepts either integer or string format.

### Casbin Configuration

```yaml
Casbin:
  app_model: "etc/casbin/cmn/model.conf"
  app_policy: "etc/casbin/cmn/policy.csv"
  resource_model: "etc/casbin/resources/model.conf"
  resource_policy: "etc/casbin/resources/policy.csv"
```

**Dual Enforcer Setup**:
- **App Enforcer**: Controls API endpoint access
- **Resource Enforcer**: Controls resource-level permissions

## Environment-Specific Configuration

### Development

Development configuration (`app.dev.yaml`) includes:
- Localhost database connections
- Default credentials for Docker Compose
- Debug-friendly settings

### Production

Production configuration must include:
- Secure JWT secrets (256+ bits)
- Production database credentials
- Redis credentials
- Appropriate connection pool sizes
- Production-grade Casbin policies

## Security Best Practices

1. **Never commit `etc/app.yaml`** - It's in `.gitignore` for a reason
2. **Use environment variables** for sensitive data (optional approach)
3. **Rotate JWT secrets** regularly
4. **Use strong passwords** for MySQL and Redis
5. **Limit connection pool sizes** based on your infrastructure
6. **Review Casbin policies** before deployment

## Environment Variables (Alternative)

While cmn-core primarily uses YAML configuration, you can also use environment variables:

```bash
export CMN_JWT_SECRET="your-secret"
export CMN_MYSQL_PASSWORD="your-db-password"
export CMN_REDIS_PASSWORD="your-redis-password"
```

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
