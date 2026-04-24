# Configuration Files Guide

## File Structure

| File | Purpose | Description |
|------|---------|-------------|
| `app.yaml.example` | Template | Configuration template for production environment |
| `app.dev.yaml` | Development | Development environment settings (works with docker-compose.yaml) |
| `app.yaml` | Runtime config | **Not included in Git** - Actual configuration file |

## Setup Instructions

### 1. Development Environment
```bash
# Use development configuration
cp etc/app.dev.yaml etc/app.yaml

# Or use Docker Compose environment directly
make dev-up
```

### 2. Production Environment
```bash
# Create production config from template
cp etc/app.yaml.example etc/app.yaml

# Edit configuration file with actual credentials
vi etc/app.yaml
```

## Required Configuration Changes

For production environment, you must change the following items:

### JWT Configuration
```yaml
Server:
  jwt_secret: "YOUR_SECURE_JWT_SECRET_HERE"  # Random string of 256+ bits
  jwt:
    key: "YOUR_JWT_KEY_HERE"
```

### Database Configuration
```yaml
MySQL:
  host: "your-mysql-host"
  user: "your-mysql-user"
  pass: "your-mysql-password"
  port: 3306
  db: "your-database-name"
```

### Redis Configuration
```yaml
Redis:
  host: "your-redis-host"
  port: 6379
  user: "default"
  pass: "your-redis-password"
  db: "your-redis-database"
```

### Client Configuration
```yaml
Client:
  ServerEndpoint: "https://your-production-domain.com"
  UserEmail: "your-default-user@domain.com"
  UserPassword: "your-secure-password"
```

### Administrator Configuration
```yaml
Server:
  admin:
    emails:
      - "admin@your-domain.com"
```

## Security Considerations

- `etc/app.yaml` is included in `.gitignore` and will not be committed to Git
- Never include production credentials in code repositories
- Use sufficiently long (256+ bits) random strings for JWT secrets
- Regularly rotate passwords and JWT secrets

## Environment Variable Configuration (Recommended)

For production environments, using environment variables instead of configuration files is also recommended:

```bash
export JWT_SECRET="your-jwt-secret"
export MYSQL_HOST="your-mysql-host"
export MYSQL_USER="your-mysql-user"
export MYSQL_PASS="your-mysql-password"
export REDIS_HOST="your-redis-host"
export REDIS_PASS="your-redis-password"
```

Support for environment variables in application code is planned for future versions.