# Configuration Guide

cmn-core uses YAML-based configuration files for all runtime settings. This guide explains how to configure cmn-core for different environments.

## Configuration Files

| File | Purpose | Status |
|------|---------|--------|
| `etc/app.yaml` | Runtime configuration | **Not in Git** (local only) |
| `etc/app.yaml.example` | Production template | In Git (template) |
| `etc/app.dev.yaml` | Development defaults | In Git (for development) |

## Quick Start

### Development Setup

```bash
# Copy development configuration
cp etc/app.dev.yaml etc/app.yaml

# Or start with Docker Compose (auto-configured)
docker-compose up -d
```

### Production Setup

```bash
# Create configuration from template
cp etc/app.yaml.example etc/app.yaml

# Edit with your actual credentials
vi etc/app.yaml
```

## Configuration Sections

### Server Configuration

```yaml
Server:
  host: "0.0.0.0"
  port: 8080
  jwt_secret: "your-secure-secret-256-bits"
  jwt:
    key: "your-jwt-key"
```

**Important**:
- `jwt_secret`: Must be at least 256 bits (32 characters)
- `jwt.key`: Used for signing JWT tokens
- Generate secure random values for production

### Database Configuration

```yaml
MySQL:
  host: "mysql-host"
  user: "mysql-user"
  pass: "mysql-password"
  port: 3306
  db: "database-name"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
```

**Connection Pool Settings**:
- `max_idle_conns`: Maximum idle connections (default: 10)
- `max_open_conns`: Maximum open connections (default: 100)
- `conn_max_lifetime`: Connection lifetime in seconds (default: 3600)

### Redis Configuration

```yaml
Redis:
  host: "redis-host"
  port: 6379
  pass: "redis-password"
  db: 0
```

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
