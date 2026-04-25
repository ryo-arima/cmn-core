# Components Overview

| Layer | Purpose | Key Components |
|---|---|---|
| [API Layer](./api-layer.md) | HTTP routing and JWT validation | Gin Router, JWT middleware, Logger |
| [Controller Layer](./controller-layer.md) | Request handling | Share, Internal, Private controllers |
| [Business Logic Layer](./business-layer.md) | Use-case implementations | Application business rules |
| [Repository Layer](./repository-layer.md) | Data access abstraction | PostgreSQL, Redis repositories |
| [Data Layer](./data-layer.md) | Persistent storage | PostgreSQL, Redis |

Identity management (users, roles, groups) is fully delegated to **Keycloak** and **Casdoor**.

- [Controller Layer](./controller-layer.md)
- [Business Logic Layer](./business-layer.md)
- [Repository Layer](./repository-layer.md)
- [Data Layer](./data-layer.md)
