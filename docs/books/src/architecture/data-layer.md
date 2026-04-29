# Data Layer

## PostgreSQL

The application database is `cmn_core`. It is initialized by `scripts/data/postgres/init.sql` on first start.

### Tables

#### `pg_groups`

Used only when `idp.provider = "casdoor"`. Stores the mapping between the UUID used as the Casdoor group `name` and the human-readable display name.

```sql
CREATE TABLE pg_groups (
    id         BIGSERIAL    PRIMARY KEY,
    uuid       TEXT         NOT NULL UNIQUE,
    name       TEXT         NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ  -- soft delete
);
```

Seed data (100 groups) is embedded in `init.sql` and can be regenerated with:

```bash
python3 scripts/gen_group_seed.py
```

This script also updates `etc/casdoor/init_data.json` so Casdoor group names match the UUIDs in `pg_groups`.

### Connection

Configured under `PostgreSQL:` in `server.yaml`:

```yaml
PostgreSQL:
  host: "localhost"
  user: "user"
  pass: "password"
  port: "5432"
  db: "cmn_core"
  sslmode: "disable"
```

## Redis

Used for JWT token caching and the token denylist (logout).

```yaml
Redis:
  host: "localhost"
  port: 6379
  user: "default"
  pass: ""
  db: 0
```

Token caching is controlled by `Application.Server.redis.jwt_cache` and `cache_ttl` (seconds).

## Dev Tools

| Tool | URL | Credentials |
|---|---|---|
| pgAdmin | http://localhost:3001 | admin@cmn.local / admin |
| Roundcube (mail) | http://localhost:3005 | test1@cmn.local / TestPassword123! |
