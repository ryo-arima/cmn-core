# Getting Started

## Prerequisites

- Go 1.24+
- Docker & Docker Compose

## Setup

```bash
git clone https://github.com/ryo-arima/cmn-core.git
cd cmn-core

# Platform setting (Apple Silicon)
echo "DOCKER_PLATFORM=linux/arm64" > .env.local

# Copy config
cp etc/app.yaml.example etc/app.yaml

# Start all dev services
make dev-up
```

Services started:

| Service | URL | Credentials |
|---|---|---|
| cmn-core server | http://localhost:8000 | — |
| Keycloak admin | http://localhost:8080/admin | admin / admin |
| Keycloak user portal | http://localhost:8080/realms/cmn/account | user01-10 / Password123! |
| Casdoor admin | http://localhost:9000 | admin / 123 |
| Casdoor user portal | http://localhost:9000/login/cmn | user01-10 / Password123! |
| PostgreSQL | localhost:5432 | user / password |
| Redis | localhost:6379 | — |
| pgAdmin | http://localhost:5050 | — |
| Roundcube | http://localhost:3005 | — |

## Build & Run

```bash
make build
./.bin/cmn-server
```

## Stop

```bash
make dev-down
```

## Makefile Reference

```bash
make build      # Build all binaries
make test       # Run unit tests
make dev-up     # Start dev environment
make dev-down   # Stop and remove volumes
make docs       # Build documentation
```


## Manual Setup

### 1. Install Dependencies

```bash
# Install Go dependencies
go mod download
go mod vendor
```

### 2. Setup Database

```bash
# Create MySQL database
mysql -u root -p
> CREATE DATABASE cmn_core;
> CREATE USER 'cmn_core'@'localhost' IDENTIFIED BY 'password';
> GRANT ALL PRIVILEGES ON cmn_core.* TO 'cmn_core'@'localhost';
> FLUSH PRIVILEGES;
```

### 3. Setup Redis

```bash
# Start Redis (macOS with Homebrew)
brew services start redis

# Or with Docker
docker run -d -p 6379:6379 redis:alpine
```

### 4. Configure Application

```bash
# Copy configuration template
cp etc/app.yaml.example etc/app.yaml

# Edit configuration
vi etc/app.yaml
```

Update the following sections:
- MySQL connection details
- Redis connection details
- JWT secrets

### 5. Build the Application

```bash
# Build all binaries
make build

# Or build individually
go build -o .bin/cmn-server ./cmd/server/main.go
go build -o .bin/cmn-client-admin ./cmd/client/admin/main.go
go build -o .bin/cmn-client-app ./cmd/client/app/main.go
go build -o .bin/cmn-client-anonymous ./cmd/client/anonymous/main.go
```

### 6. Run the Server

```bash
./.bin/cmn-server
```

## Verify Installation

### Check Server Status

```bash
# Health check
curl http://localhost:8080/v1/public/health

# Expected response:
# {"status": "ok"}
```

### Register a User

```bash
curl -X POST http://localhost:8080/v1/public/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "securepassword"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/v1/public/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "securepassword"
  }'

# Save the JWT token from response
export TOKEN="<jwt_token_here>"
```

### Access Protected Endpoint

```bash
curl http://localhost:8080/v1/internal/users \
  -H "Authorization: Bearer $TOKEN"
```

## Using the CLI Clients

### Admin Client

```bash
./.bin/cmn-client-admin --help
```

### App Client

```bash
./.bin/cmn-client-app --help
```

### Anonymous Client

```bash
./.bin/cmn-client-anonymous --help
```

## Makefile Commands

cmn-core includes a Makefile for common tasks:

```bash
# Build all binaries
make build

# Run tests
make test

# Clean build artifacts
make clean

# Generate documentation
make docs

# Start development environment
make dev-up

# Stop development environment
make dev-down
```

## Development Workflow

1. **Make Changes**: Edit source code
2. **Build**: `make build`
3. **Test**: `make test`
4. **Run**: `./.bin/cmn-server`
5. **Verify**: Test with curl or CLI clients

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Database Connection Failed

- Verify MySQL is running: `mysql -u root -p`
- Check credentials in `etc/app.yaml`
- Ensure database exists: `SHOW DATABASES;`

### Redis Connection Failed

- Verify Redis is running: `redis-cli ping`
- Check Redis host/port in configuration
- Test connection: `redis-cli -h localhost -p 6379`

### Casbin Policy Errors

- Verify policy files exist in `etc/casbin/`
- Check CSV format (no trailing commas)
- Validate model syntax

## Next Steps

- [Configuration Guide](../configuration/guide.md)
- [Building](./building.md)
- [Testing](./testing.md)
- [API Overview](../api/overview.md)
