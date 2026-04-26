# cmn-core

<p align="center">
  <img src="docs/images/image01.png" alt="cmn-core" width="600">
  <br>
  <sub>The Go gopher was designed by <a href="https://reneefrench.blogspot.com/">Renée French</a>, licensed under <a href="https://creativecommons.org/licenses/by/4.0/">CC BY 4.0</a>.</sub>
</p>

Go-based API server that delegates authentication and authorization entirely to external IdPs.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![E2E Tests](https://github.com/ryo-arima/cmn-core/actions/workflows/e2e-test.yml/badge.svg)](https://github.com/ryo-arima/cmn-core/actions/workflows/e2e-test.yml)
[![Documentation](https://img.shields.io/badge/docs-GitHub%20Pages-success)](https://ryo-arima.github.io/cmn-core/)

## Overview

cmn-core provides a multi-tier REST API backed by PostgreSQL and Redis. All user management, authentication, and authorization are fully delegated to external IdPs — **Keycloak** (port 8080) and **Casdoor** (port 9000). The application itself performs no internal user CRUD.

CLI clients authenticate by posting credentials to the server (`POST /v1/public/login`); the server then obtains a JWT from the configured IdP and returns it to the client. Clients never contact the IdP directly.

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose

### Start Dev Environment

```bash
cp etc/app.yaml.example etc/app.yaml
make dev-up
```

Services started by `make dev-up`:

| Service | URL | Credentials |
|---|---|---|
| PostgreSQL | localhost:5432 | user/password |
| Redis | localhost:6379 | — |
| Keycloak (admin) | http://localhost:8080/admin | admin / admin |
| Keycloak (cmn realm users) | http://localhost:8080/realms/cmn/account | user01-50 / Password123! |
| Casdoor (admin) | http://localhost:9000 | admin / 123 |
| Casdoor (cmn org users) | http://localhost:9000/login/cmn | user01-50 / Password123! |
| Casdoor (cmn admin user) | http://localhost:9000/login/cmn | admin@example.com / Admin123! |
| pgAdmin | http://localhost:5050 | — |
| Roundcube | http://localhost:3005 | — |

### Build & Run Server

```bash
make build
./.bin/cmn-server
```

Server starts on `http://localhost:8000`.

## Architecture

Authentication is handled by external IdPs via OIDC / SAML 2.0. The server validates tokens issued by Keycloak or Casdoor.

```
Client → [Keycloak / Casdoor] → JWT Token → cmn-core Server → PostgreSQL / Redis
```

### API Tiers

| Tier | Path prefix | Auth |
|---|---|---|
| Public | `/v1/public/` | None (login endpoint) |
| Share | `/v1/share/` | JWT required |
| Internal | `/v1/internal/` | JWT required |
| Private | `/v1/private/` | JWT + admin role required |

OIDC callback: `GET /v1/share/auth/oidc/callback`

## Project Structure

```
cmd/                  # Entrypoints (server, client/admin, client/app, client/anonymous)
pkg/
  auth/               # OIDC / SAML provider interface and implementations
  config/             # Configuration loading (YAML + AWS Secrets Manager)
  entity/             # Models, request/response structs
  server/             # Gin router, controllers, repositories
  client/             # CLI client implementations
etc/
  app.yaml.example    # Configuration template
  keycloak/           # Keycloak realm import (cmn-realm.json)
  casdoor/            # Casdoor config (app.conf) and seed data (init_data.json)
scripts/
  data/postgres/      # PostgreSQL init SQL
test/
  unit/               # Unit tests
  e2e/                # End-to-end tests
```

## Configuration

```bash
cp etc/app.yaml.example etc/app.yaml
```

Key settings in `etc/app.yaml`:

```yaml
Application:
  Server:
    port: 8000
    admin:
      emails:
        - "admin@example.com"
    jwt_secret: "CHANGE_THIS_JWT_SECRET_IN_PRODUCTION"
    log_level: "info"
    redis:
      jwt_cache: true
      cache_ttl: 1800
    idp:
      provider: "casdoor"  # or "keycloak"
      casdoor:
        base_url: "http://localhost:9000"
        client_id: "cmn-core-client-id"
        client_secret: "cmn-core-client-secret"
        organization: "cmn"
    auth:
      oidc:
        issuer_url: "http://localhost:9000"   # JWT iss claim value
        provider_url: "http://localhost:9000" # OIDC discovery URL (set to internal URL in Docker)
        client_id: "cmn-core-client-id"

PostgreSQL:
  host: "localhost"
  user: "user"
  pass: "password"
  port: "5432"
  db: "cmn_core"
  sslmode: "disable"

Redis:
  host: "localhost"
  port: 6379
  user: "default"
  pass: ""
  db: 0
```

## Development

```bash
# Run unit tests
make test

# Run with coverage
go test -v -cover ./...

# E2E tests (requires dev environment running)
go test -v -timeout 15m ./test/e2e/testcase/
```

### Platform (Apple Silicon)

```bash
echo "DOCKER_PLATFORM=linux/arm64" > .env.local
```

## Mail Sandbox (Experimental)

Ephemeral Postfix/Dovecot + Roundcube environment for email testing:

```bash
./scripts/main.sh env recreate
open http://localhost:3005   # Roundcube — user: test1@cmn.local / pass: TestPassword123!
```

## Documentation

Full documentation: **[https://ryo-arima.github.io/cmn-core/](https://ryo-arima.github.io/cmn-core/)**

```bash
# Build docs locally
make docs

# Serve docs locally
make docs-serve
```

## License

MIT — see [LICENSE](LICENSE).
