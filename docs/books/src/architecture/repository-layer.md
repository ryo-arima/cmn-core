# Repository Layer

The repository layer abstracts all external data sources: the configured IdP (Keycloak or Casdoor) and PostgreSQL/Redis.

## IdPManager

`repository.IdPManager` is the single interface to the identity provider.

```go
type IdPManager interface {
    // Users
    GetUser(ctx, id)    (*model.LoUser, error)
    ListUsers(ctx)      ([]model.LoUser, error)
    CreateUser(ctx, input) (*model.LoUser, error)
    UpdateUser(ctx, id, input) error
    DeleteUser(ctx, id) error

    // Groups
    GetGroup(ctx, id)   (*model.LoGroup, error)
    ListGroups(ctx)     ([]model.LoGroup, error)
    CreateGroup(ctx, input) (*model.LoGroup, error)
    UpdateGroup(ctx, id, input) error
    DeleteGroup(ctx, id) error

    // Members & Roles (similar pattern)
    ...
}
```

Two implementations exist:

| File | IdP |
|---|---|
| `pkg/server/repository/casdoor.go` | Casdoor REST API |
| `pkg/server/repository/keycloak.go` | Keycloak Admin REST API |

The implementation is selected at startup based on `Application.Server.IdP.provider`.

## Group (pg_groups) — Casdoor only

`repository.Group` persists the UUID ↔ display name mapping when using Casdoor.

```go
type Group interface {
    Upsert(ctx, uuid, name string) error
    LookupName(ctx, uuid string) string
    LookupNames(ctx, uuids []string) map[string]string
    SoftDelete(ctx, uuid string) error
}
```

Implemented in `pkg/server/repository/group_psql.go` using GORM against the `pg_groups` table.

The instance is `nil` when using Keycloak (groups have native UUIDs).

## Common Repository

`repository.Common` handles token operations (validation, caching, denylist) using Redis.

## Initialization

All repositories are wired in `pkg/server/router.go`:

```go
idpManager, _ := repository.NewIdPManager(conf)   // Keycloak or Casdoor
var groupStore repository.Group
if provider == "casdoor" && dbConnection != nil {
    groupStore = repository.NewGroup(dbConnection)
}
```
