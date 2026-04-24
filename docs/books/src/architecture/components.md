# Components Overview

This section provides detailed information about each layer of the cmn-core architecture.

## Layer Summary

| Layer | Purpose | Key Components |
|-------|---------|----------------|
| [Client Layer](./client-layer.md) | User interfaces and CLI tools | Admin CLI, App CLI, Anonymous CLI |
| [API Layer](./api-layer.md) | HTTP routing and middleware | Gin Router, JWT Auth, Casbin RBAC, Logger |
| [Controller Layer](./controller-layer.md) | Request handling and validation | Public, Internal, Private Controllers |
| [Business Logic Layer](./business-layer.md) | Use cases and business rules | User, Group, Member, Role Usecases |
| [Repository Layer](./repository-layer.md) | Data access abstraction | User, Group, Member, Role Repositories |
| [Data Layer](./data-layer.md) | Persistent storage | MySQL/TiDB, Redis, Casbin Policies |

## Component Interaction

Each layer communicates with adjacent layers through well-defined interfaces:

```
Client → API → Controller → Business Logic → Repository → Data
```

### Dependency Flow

- **Downward Dependencies**: Each layer depends only on the layer directly below it
- **Interface-Based**: Layers communicate through interfaces, not concrete implementations
- **Testability**: Each layer can be tested independently with mocks

## Common Patterns

### Repository Pattern

All data access goes through repositories, providing:
- Abstraction from underlying storage
- Consistent error handling
- Transaction management
- Query optimization

### Use Case Pattern

Business logic is encapsulated in use cases:
- Single responsibility per use case
- Orchestration of multiple repositories
- Business rule validation
- Domain logic isolation

### Middleware Pattern

Cross-cutting concerns are handled by middleware:
- Authentication
- Authorization
- Logging
- Error handling

## Next Steps

Explore each layer in detail:
- [Client Layer](./client-layer.md)
- [API Layer](./api-layer.md)
- [Controller Layer](./controller-layer.md)
- [Business Logic Layer](./business-layer.md)
- [Repository Layer](./repository-layer.md)
- [Data Layer](./data-layer.md)
