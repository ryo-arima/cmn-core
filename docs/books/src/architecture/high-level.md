# High-Level Architecture

This page provides an overview of cmn-core's high-level architecture and how its components interact.

## Architecture Diagram

![High-Level Architecture](../../architecture/high-level-architecture.svg)

## System Layers

cmn-core is built on a multi-layered architecture that promotes separation of concerns and maintainability:

### 1. Client Layer

The client layer provides multiple CLI tools for different use cases:

- **Admin CLI** (`cmn-client-admin`): Administrative operations with elevated privileges
- **App CLI** (`cmn-client-app`): Application-level operations for authenticated users
- **Anonymous CLI** (`cmn-client-anonymous`): Public operations that don't require authentication

### 2. API Layer

The API layer handles HTTP requests and applies cross-cutting concerns:

- **Gin Router**: High-performance HTTP router and middleware framework
- **Middleware Stack**:
  - **JWT Authentication**: Validates and decodes JWT tokens
  - **Casbin Authorization**: Enforces RBAC policies
  - **Request Logger**: Logs all incoming requests for auditing

### 3. Controller Layer

Controllers are organized by access level to enforce security boundaries:

- **Public Controllers**: Endpoints accessible without authentication (e.g., user registration, login)
- **Internal Controllers**: Endpoints for authenticated internal services
- **Private Controllers**: Endpoints requiring specific permissions

Each controller handles:
- **User Controller**: User account operations
- **Group Controller**: Group management
- **Member Controller**: Group membership operations
- **Role Controller**: Permission management

### 4. Business Logic Layer

Use cases implement business rules and orchestrate repository operations:

- **User Usecase**: User business logic
- **Group Usecase**: Group business logic
- **Member Usecase**: Membership business logic
- **Role Usecase**: Permission business logic

### 5. Repository Layer

Repositories abstract data access and provide a clean interface to the data layer:

- **User Repository**: User data operations
- **Group Repository**: Group data operations
- **Member Repository**: Membership data operations
- **Role Repository**: Role and permission operations (Casbin-backed)
- **Common Repository**: Shared operations (JWT token management)
- **Redis Repository**: Cache operations

### 6. Data Layer

The data layer consists of multiple storage systems:

- **MySQL/TiDB**: Primary relational database for users, groups, and members
- **Redis**: Caching layer for JWT token denylist and session management
- **Casbin Policies**: Authorization policy storage
  - **Auth Policy**: Application-wide permissions
  - **Resource Policy**: Resource-specific permissions

## Data Flow

### Authentication Flow

1. User sends credentials to Public Controller (Login endpoint)
2. Controller validates credentials via User Repository
3. JWT token is generated and returned to client
4. Subsequent requests include JWT token in Authorization header
5. Middleware validates token and extracts user information
6. Casbin enforcer checks user permissions
7. Request proceeds to appropriate controller if authorized

### Authorization Flow

cmn-core uses a dual-enforcer approach for fine-grained access control:

1. **App Enforcer** (`etc/casbin/cmn/`): Controls access to API endpoints
2. **Resource Enforcer** (`etc/casbin/resources/`): Controls access to specific resources

### CRUD Operation Flow

1. Client sends request through CLI
2. Gin router routes to appropriate controller
3. Middleware validates authentication and authorization
4. Controller parses and validates request
5. Use case applies business logic
6. Repository performs database operations
7. Response flows back through the layers

## Configuration

The system is configured through:

- **app.yaml**: Main configuration file (database, Redis, JWT settings)
- **Casbin Model Files**: Define RBAC model structure
- **Casbin Policy Files**: Define actual permissions

## Scalability Considerations

- **Stateless API**: JWT tokens enable horizontal scaling
- **Redis Caching**: Reduces database load
- **Connection Pooling**: Efficient database connection management
- **Casbin Policy Caching**: In-memory policy evaluation for fast authorization

## Security Features

- **JWT with HS256**: Secure token-based authentication
- **Token Denylist**: Revoked tokens stored in Redis
- **Casbin RBAC**: Policy-based authorization
- **Multi-tier Access**: Public, Internal, and Private endpoint separation
- **Admin Email Check**: Additional verification for administrative operations

## Next Steps

- Learn more about individual [Components](./components.md)
- Explore the [API Reference](../api/overview.md)
- Review [Configuration Guide](../configuration/guide.md)
