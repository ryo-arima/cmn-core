# GoDoc

cmn-core's Go package documentation is generated using standard Go documentation tools.

## Viewing GoDoc

### GitHub Pages

**[📖 GoDoc Overview](https://ryo-arima.github.io/cmn-core/godoc/index.html)** (Links to pkg.go.dev)

Or use the relative link when browsing documentation:
**[📖 GoDoc Overview](../../godoc/index.html)**

### Online - pkg.go.dev (When Published)

Once the package is published, view the official Go documentation online:

**[📖 View on pkg.go.dev](https://pkg.go.dev/github.com/ryo-arima/cmn-core)**

### Local Documentation

Generate and view documentation locally:

```bash
# Install godoc tool (if not already installed)
go install golang.org/x/tools/cmd/godoc@latest

# Start local documentation server
godoc -http=:6060

# Open in browser
open http://localhost:6060/pkg/github.com/ryo-arima/cmn-core/
```

### Generated HTML

Pre-generated HTML documentation is available in:

```
docs/godoc/
```

## Package Structure

### Main Packages

| Package | Description |
|---------|-------------|
| `pkg/server` | HTTP server implementation |
| `pkg/client` | CLI client implementations |
| `pkg/config` | Configuration management |
| `pkg/entity` | Data models and DTOs |

### Server Packages

| Package | Description |
|---------|-------------|
| `pkg/server/controller` | Request handlers |
| `pkg/server/middleware` | Middleware components |
| `pkg/server/repository` | Data access layer |

### Entity Packages

| Package | Description |
|---------|-------------|
| `pkg/entity/model` | Database models |
| `pkg/entity/request` | Request DTOs |
| `pkg/entity/response` | Response DTOs |

### Client Packages

| Package | Description |
|---------|-------------|
| `pkg/client/controller` | CLI command handlers |
| `pkg/client/repository` | API client layer |
| `pkg/client/usecase` | Client business logic |

## Key Types and Functions

### Configuration

```go
// BaseConfig holds application configuration
type BaseConfig struct {
    YamlConfig YamlConfig
    DB         *gorm.DB
}

// LoadConfig loads configuration from YAML file
func LoadConfig(path string) (*BaseConfig, error)
```

### Server

```go
// InitRouter initializes the Gin router with all routes and middleware
func InitRouter(conf config.BaseConfig) *gin.Engine

// Start starts the HTTP server
func Start(router *gin.Engine, port int) error
```

### Repository

```go
// UserRepository provides user data access
type UserRepository interface {
    ListUsers(filters map[string]interface{}) ([]model.User, error)
    GetUser(id uint) (*model.User, error)
    CreateUser(user *model.User) error
    UpdateUser(user *model.User) error
    DeleteUser(id uint) error
}
```

### Entity Models

```go
// User represents a user account
type User struct {
    ID        uint      `gorm:"primaryKey"`
    UUID      string    `gorm:"uniqueIndex"`
    Name      string
    Email     string    `gorm:"uniqueIndex"`
    Password  string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `gorm:"index"`
}
```

## Documentation Standards

### Package Comments

Every package should have a package-level comment:

```go
// Package server implements the HTTP server and routing logic.
//
// The server package provides:
//   - Gin-based HTTP routing
//   - Middleware integration (JWT, Casbin, logging)
//   - Controller initialization
//   - Repository wiring
//
// Example usage:
//
//   config := config.LoadConfig("etc/app.yaml")
//   router := server.InitRouter(config)
//   server.Start(router, 8080)
package server
```

### Function Comments

All exported functions must have documentation:

```go
// ListUsers retrieves a list of users based on the provided filters.
//
// Filters can include:
//   - name: Filter by name (partial match)
//   - email: Filter by email (exact match)
//   - limit: Maximum number of results
//   - offset: Number of results to skip
//
// Returns a slice of User models and an error if the operation fails.
func (r *userRepository) ListUsers(filters map[string]interface{}) ([]model.User, error)
```

### Type Comments

All exported types must be documented:

```go
// User represents a user account in the system.
//
// Users can belong to multiple groups and have various roles
// that determine their permissions within the application.
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UUID      string    `gorm:"uniqueIndex" json:"uuid"`
    Name      string    `json:"name"`
    Email     string    `gorm:"uniqueIndex" json:"email"`
    Password  string    `json:"-"` // Never serialize password
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
```

## Generating Documentation

### Standard godoc

```bash
# Generate HTML documentation
godoc -html github.com/ryo-arima/cmn-core > docs/godoc/index.html

# Generate documentation for specific package
godoc -html github.com/ryo-arima/cmn-core/pkg/server > docs/godoc/server.html
```

### pkgsite (Modern godoc)

```bash
# Install pkgsite
go install golang.org/x/pkgsite/cmd/pkgsite@latest

# Run local server
pkgsite -http=:6060

# View documentation
open http://localhost:6060/github.com/ryo-arima/cmn-core
```

## Documentation Best Practices

### Writing Good Documentation

1. **Start with a summary**: One-line description of what it does
2. **Explain the why**: Not just what, but why it exists
3. **Provide examples**: Show common usage patterns
4. **Document parameters**: Explain each parameter's purpose
5. **Document return values**: Explain what's returned and when
6. **Document errors**: List possible error conditions
7. **Link related items**: Reference related types/functions

### Example Structure

```go
// FunctionName does something specific.
//
// More detailed explanation of what the function does,
// why it exists, and how it should be used.
//
// Parameters:
//   - param1: Description of first parameter
//   - param2: Description of second parameter
//
// Returns:
//   - result: Description of return value
//   - error: Possible error conditions
//
// Example:
//
//   result, err := FunctionName(arg1, arg2)
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(result)
func FunctionName(param1 string, param2 int) (result string, err error)
```

## Package Dependencies

View package dependencies:

```bash
# Generate dependency graph
go mod graph

# Visualize with graphviz
go mod graph | dot -Tsvg -o docs/dependencies.svg
```

## Code Examples in Documentation

GoDoc supports executable examples:

```go
// Example_listUsers demonstrates how to list users
func Example_listUsers() {
    config := config.LoadConfig("etc/app.yaml")
    repo := repository.NewUserRepository(config)
    
    users, err := repo.ListUsers(map[string]interface{}{
        "limit": 10,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    for _, user := range users {
        fmt.Println(user.Name)
    }
}
```

## Testing Documentation

Verify documentation builds correctly:

```bash
# Check for documentation issues
go vet ./...

# Run documentation tests
go test -v ./...
```

## Next Steps

- [Swagger Documentation](./swagger.md)
- [API Reference](../api/overview.md)
- [Contributing](../development/contributing.md)
