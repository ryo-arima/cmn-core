# cmn-core

A robust Role-Based Access Control (RBAC) service built with Go, providing comprehensive user, group, member, and role management with fine-grained permissions.

[![Go Version](https://img.shields.io/badge/Go-1.25.5+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![E2E Tests](https://github.com/ryo-arima/cmn-core/actions/workflows/e2e-test.yml/badge.svg)](https://github.com/ryo-arima/cmn-core/actions/workflows/e2e-test.yml)
[![Documentation](https://img.shields.io/badge/docs-GitHub%20Pages-success)](https://ryo-arima.github.io/cmn-core/)

## Features

- 🔐 **JWT Authentication**: Secure token-based authentication with HS256
- 🛡️ **Casbin RBAC**: Flexible policy-based authorization
- 👥 **User Management**: Complete CRUD operations for user accounts
- 📦 **Group Management**: Organize users into logical groups
- 👤 **Member Management**: Control group membership and relationships
- 🎭 **Role Management**: Define and assign fine-grained permissions
- ⚡ **Redis Caching**: High-performance caching with token denylist
- 🗄️ **MySQL/TiDB**: Reliable persistent data storage
- 🌐 **Multi-tier API**: Public, Internal, and Private endpoints
- 📱 **CLI Clients**: Admin, App, and Anonymous command-line interfaces

## Quick Start

### Prerequisites

- Go 1.25.5 or higher
- MySQL 8.0+ or TiDB
- Redis 6.0+
- Docker & Docker Compose (optional)

### Installation

```bash
# Clone the repository
git clone https://github.com/ryo-arima/cmn-core.git
cd cmn-core

# Start dependencies with Docker Compose
docker compose up -d

# Copy configuration
cp etc/app.dev.yaml etc/app.yaml

# Build the server
go build -o .bin/cmn-server ./cmd/server/main.go

# Start the server
./.bin/cmn-server
```

The server will start on `http://localhost:8080`.

### Quick Test

```bash
# Register a user
curl -X POST http://localhost:8080/v1/public/user \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/v1/public/token \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Use the JWT token from the response for authenticated requests
```

## Architecture

cmn-core follows a clean, layered architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                       Client Layer                          │
│              (Admin CLI, App CLI, Anonymous CLI)            │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                        API Layer                            │
│        (Gin Router, JWT Auth, Casbin RBAC, Logger)         │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Controller Layer                         │
│          (Public, Internal, Private Controllers)            │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                 Business Logic Layer                        │
│              (User, Group, Member, Role Usecases)           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                          │
│         (User, Group, Member, Role Repositories)            │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      Data Layer                             │
│              (MySQL/TiDB, Redis, Casbin Policies)           │
└─────────────────────────────────────────────────────────────┘
```

[View detailed architecture diagram →](https://ryo-arima.github.io/cmn-core/architecture/high-level.html)

## Documentation

Comprehensive documentation is available at **[https://ryo-arima.github.io/cmn-core/](https://ryo-arima.github.io/cmn-core/)**

- **[Getting Started](https://ryo-arima.github.io/cmn-core/development/getting-started.html)** - Installation and setup
- **[Architecture](https://ryo-arima.github.io/cmn-core/architecture/high-level.html)** - System design and components
- **[API Reference](https://ryo-arima.github.io/cmn-core/api/overview.html)** - REST API documentation
- **[Configuration](https://ryo-arima.github.io/cmn-core/configuration/guide.html)** - Configuration guide
- **[Swagger UI](https://ryo-arima.github.io/cmn-core/swagger/index.html)** - Interactive API documentation
- **[GoDoc](https://ryo-arima.github.io/cmn-core/godoc/index.html)** - Go package documentation

## API Endpoints

### Public Endpoints (No Authentication Required)

- `POST /v1/public/user` - Register a new user
- `POST /v1/public/token` - Authenticate and get JWT token
- `GET /health` - Health check

### Share Endpoints (JWT Required)

- `POST /v1/share/token/refresh` - Refresh JWT token
- `DELETE /v1/share/token` - Logout and invalidate token
- `GET /v1/share/token/validate` - Validate JWT token
- `GET /v1/share/token/user` - Get user info from token

### Internal Endpoints (JWT + Casbin Authorization Required)

- `GET /v1/internal/users` - List users
- `GET /v1/internal/groups` - List groups
- `GET /v1/internal/members` - List members
- `GET /v1/internal/roles` - List roles
- `GET /v1/internal/resources` - List resources

### Private Endpoints (JWT Required, No Casbin)

- `PUT /v1/private/user/{id}` - Update user
- `DELETE /v1/private/user/{id}` - Delete user
- `POST /v1/private/group` - Create group
- `POST /v1/private/role` - Create role
- `POST /v1/private/resource` - Create resource

[Full API documentation →](https://ryo-arima.github.io/cmn-core/swagger/index.html)

## Configuration

Create `etc/app.yaml` from the template:

```bash
cp etc/app.yaml.example etc/app.yaml
```

Edit the configuration with your settings:

```yaml
Server:
  host: "0.0.0.0"
  port: 8080
  jwt_secret: "your-secure-secret-256-bits"
  Redis:
    JWTCache: true      # Enable JWT token caching
    CacheTTL: 3600      # Cache TTL in seconds

MySQL:
  host: "localhost"
  user: "cmn-core"
  pass: "password"
  db: "cmn-core"

Redis:
  host: "localhost"
  port: 6379
  db: 0

Casbin:
  app_model: "etc/casbin/cmn/model.conf"
  app_policy: "etc/casbin/cmn/policy.csv"
```

[Configuration guide →](https://ryo-arima.github.io/cmn-core/configuration/guide.html)

## Development

### Building

```bash
# Build all binaries
make build

# Or build individually
go build -o .bin/cmn-server ./cmd/server/main.go
go build -o .bin/cmn-client-admin ./cmd/client/admin/main.go
go build -o .bin/cmn-client-app ./cmd/client/app/main.go
go build -o .bin/cmn-client-anonymous ./cmd/client/anonymous/main.go
```

### Testing

#### Unit Tests

```bash
# Run tests
make test

# Run with coverage
go test -v -cover ./...
```

#### E2E Tests

E2E tests verify the entire system including server, database, Redis, and CLI clients.

**Prerequisites:**
- Docker & Docker Compose
- Go 1.25.5+

**Platform Configuration:**

For Apple Silicon (M1/M2/M3) Macs:
```bash
# Create .env.local (ignored by git)
echo "DOCKER_PLATFORM=linux/arm64" > .env.local
```

For Intel Macs and Linux x86_64:
```bash
# Create .env.local (ignored by git)
echo "DOCKER_PLATFORM=linux/amd64" > .env.local
```

**Run E2E Tests:**

```bash
# 1. Start dependencies (MySQL and Redis)
docker compose up -d postgres redis

# 2. Wait for services to be ready (about 10 seconds)
sleep 10

# 3. Build CLI binaries
mkdir -p bin
go build -o bin/cmn-admin ./cmd/client/admin
go build -o bin/cmn-app ./cmd/client/app
go build -o bin/cmn-anonymous ./cmd/client/anonymous

# 4. Run E2E tests
go test -v -timeout 15m ./test/e2e/testcase/

# 5. Cleanup
docker compose down -v
```

**Test Coverage:**
- ✅ Anonymous User Registration
- ⏭️ Authentication Flow (Login, Token validation, Refresh, Logout) - Endpoints need implementation
- ⏭️ App User Group CRUD - Requires authentication implementation
- ⏭️ User/Role Read Operations - Requires authentication implementation
- ⏭️ Admin operations - Requires admin role assignment

**Note:** Most E2E tests are currently skipped pending full authentication and authorization implementation.

**GitHub Actions:**
E2E tests run automatically on pull requests and pushes to `dev` branch using `linux/amd64` platform.

## Ephemeral Mail/Test Environment (Experimental)

An internal mail sandbox (Postfix/Dovecot via docker-mailserver + dnsmasq + Roundcube) can be fully recreated for browser‑based tests. All data is ephemeral.

```bash
# Full teardown & rebuild (containers, network, volumes) + account provisioning
./scripts/main.sh env recreate

# Access Roundcube (webmail)
open http://localhost:3005  # or manually open in browser

# Example login
#   user: test1@cmn.local
#   pass: TestPassword123!
```

Send a test message from test1 to test2 and verify it appears in test2's inbox after switching accounts. Logs:

```bash
# Postfix / Dovecot logs (mailserver container)
docker compose logs -f mailserver
```

To iterate after config changes always use force recreate:

```bash
docker compose up -d --force-recreate mailserver roundcube
```

If authentication fails, rerun the full recreate script to ensure accounts are re-applied cleanly.

### Documentation

```bash
# Build documentation
./scripts/main.sh docs build

# Serve documentation locally
cd docs/dist && python3 -m http.server 8000
```

### Publishing to pkg.go.dev

To make your package available on pkg.go.dev, you have two options:

#### Option 1: GitHub Actions (Recommended)

**Manual trigger:**
1. Go to **Actions** tab in GitHub
2. Select "Create Release Tag" workflow
3. Click "Run workflow"
4. Enter version (e.g., `v0.1.0`)
5. Add release notes (optional)
6. Click "Run workflow"

**Automatic trigger:**
1. Update the `VERSION` file with new version (e.g., `0.2.0`)
2. Commit and push to main branch
3. GitHub Actions will automatically create the tag

#### Option 2: Manual Script

```bash
# Publish with default version (v0.1.0)
./scripts/publish-to-pkggodev.sh

# Or specify a version
./scripts/publish-to-pkggodev.sh v0.2.0
```

After publishing (either method), your package will be available at:
- https://pkg.go.dev/github.com/ryo-arima/cmn-core
- https://pkg.go.dev/github.com/ryo-arima/cmn-core/pkg/server
- https://pkg.go.dev/github.com/ryo-arima/cmn-core/pkg/client

**Note:** It may take 5-10 minutes for pkg.go.dev to index the package after tag creation.

## CLI Clients

### Admin Client

Administrative operations with elevated privileges:

```bash
./.bin/cmn-client-admin --help
```

### App Client

Application-level operations for authenticated users:

```bash
./.bin/cmn-client-app --help
```

### Anonymous Client

Public operations without authentication:

```bash
./.bin/cmn-client-anonymous --help
```

## Project Structure

```
cmn-core/
├── cmd/                    # Command-line applications
│   ├── server/            # HTTP server
│   └── client/            # CLI clients (admin, app, anonymous)
├── pkg/                   # Shared packages
│   ├── server/           # Server implementation
│   │   ├── controller/   # HTTP handlers
│   │   ├── middleware/   # Middleware components
│   │   └── repository/   # Data access layer
│   ├── client/           # Client implementation
│   ├── config/           # Configuration management
│   └── entity/           # Data models and DTOs
├── etc/                   # Configuration files
│   ├── casbin/           # Casbin policy files
│   └── app.yaml.example  # Configuration template
├── docs/                  # Documentation
│   ├── books/            # mdBook source
│   ├── dist/             # Built documentation (GitHub Pages)
│   ├── architecture/     # Architecture diagrams
│   └── swagger/          # OpenAPI specification
├── scripts/              # Build and utility scripts
└── test/                 # Test files
```

## Technology Stack

- **Language**: Go 1.25.5+
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Authentication**: JWT with [golang-jwt](https://github.com/golang-jwt/jwt)
- **Authorization**: [Casbin](https://casbin.org/) (Group-based RBAC)
- **Database**: MySQL 8.0+ / TiDB
- **Cache**: Redis 6.0+ (JWT token caching, session management)
- **Documentation**: mdBook, Swagger/OpenAPI

## Contributing

Contributions are welcome! Please read our [Contributing Guide](https://ryo-arima.github.io/cmn-core/development/contributing.html) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Testing

### E2E Tests

End-to-end tests are located in `test/e2e/testcase/` and test the complete application stack.

#### Running E2E Tests

**1. Start required services (MySQL and Redis only):**
```bash
# For ARM64 (Mac M1/M2/M3, etc.)
docker compose up -d postgres redis

# For AMD64 (GitHub Actions, x86_64 Linux, etc.)
DOCKER_PLATFORM=linux/amd64 docker compose up -d postgres redis
```

**2. Build CLI binaries:**
```bash
make build
```

**3. Run E2E tests:**
```bash
# Run all E2E tests
go test -v ./test/e2e/testcase/

# Run specific test
go test -v ./test/e2e/testcase/ -run TestAuthenticationFlow
```

**4. Stop services:**
```bash
docker compose down
```

#### Test Coverage

Currently implemented E2E tests:
- ✅ **TestMain** - Test environment initialization
- ⏭️ All other tests - Pending authentication and authorization implementation

**Note:** E2E test implementation is in progress. Most test cases are currently disabled pending full system integration.

#### CI/CD

E2E tests run automatically on GitHub Actions for every pull request and push to main branch. See [`.github/workflows/e2e-test.yml`](.github/workflows/e2e-test.yml) for details.

## Links

- **Documentation**: https://ryo-arima.github.io/cmn-core/
- **API Documentation**: https://ryo-arima.github.io/cmn-core/swagger/index.html
- **GoDoc**: https://ryo-arima.github.io/cmn-core/godoc/index.html (or [pkg.go.dev](https://pkg.go.dev/github.com/ryo-arima/cmn-core) when published)
- **Issue Tracker**: https://github.com/ryo-arima/cmn-core/issues
- **Discussions**: https://github.com/ryo-arima/cmn-core/discussions

## Support

- 📖 [Documentation](https://ryo-arima.github.io/cmn-core/)
- 💬 [GitHub Discussions](https://github.com/ryo-arima/cmn-core/discussions)
- 🐛 [Issue Tracker](https://github.com/ryo-arima/cmn-core/issues)

---

Made with ❤️ by [Ryo ARIMA](https://github.com/ryo-arima)
