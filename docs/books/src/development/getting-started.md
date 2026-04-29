# Getting Started

## Prerequisites

- Go 1.25+
- Docker & Docker Compose

## Setup

```bash
git clone https://github.com/ryo-arima/cmn-core.git
cd cmn-core

# Platform setting (Apple Silicon)
echo "DOCKER_PLATFORM=linux/arm64" > .env.local

# Start all dev services (PostgreSQL, Redis, Keycloak, Casdoor, ...)
make dev-up

# Build and start the application server container
make svr-up
```

Services started by `make dev-up`:

| Service | URL | Credentials |
|---|---|---|
| Keycloak admin | http://localhost:8080/admin | admin / admin |
| Keycloak user portal | http://localhost:8080/realms/cmn/account | user01-50 / Password123! |
| Casdoor admin | http://localhost:9000 | admin / 123 |
| Casdoor user portal | http://localhost:9000/login/cmn | user01-50 / Password123! |
| Casdoor admin user | http://localhost:9000/login/cmn | admin@cmn.local / Admin123! |
| PostgreSQL | localhost:5432 | user / password |
| Redis | localhost:6379 | — |
| pgAdmin | http://localhost:5050 | — |
| Roundcube | http://localhost:3005 | — |

The application server (`cmn-server`) starts on `http://localhost:8000` after `make svr-up`.

## Build Binaries Locally

```bash
make build
# Binaries are placed in .bin/
```

## Stop

```bash
make dev-down   # stop all infra containers
make svr-down   # stop only the server container
```

## Makefile Reference

```bash
make build      # Build all binaries (.bin/)
make test       # Run all tests
make test-unit  # Run unit tests only
make dev-up     # Start infra (PostgreSQL, Redis, IdPs, ...)
make dev-down   # Stop and remove volumes
make svr-up     # Build Docker image and start server container
make svr-down   # Stop server container
make docs       # Start Swagger UI + GoDoc containers
```

## CLI Usage

Admin client example:

```bash
.bin/admin-client -c etc/.cmn/client/credentials/admin.yaml user list
```

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
