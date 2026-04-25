# Introduction

**cmn-core** is a Go-based API server that delegates all authentication and authorization to external Identity Providers (IdPs).

## Overview

cmn-core provides a multi-tier REST API backed by PostgreSQL and Redis. It does **not** manage users, groups, roles, or passwords internally — all identity management is handled by:

- **Keycloak** (port 8080) — enterprise-grade IdP, OIDC & SAML 2.0
- **Casdoor** (port 9000) — lightweight IdP, OIDC & SAML 2.0 (Apache 2.0)

## Key Features

- **External IdP delegation**: OIDC / SAML 2.0 via Keycloak or Casdoor
- **Multi-tier API**: Share, Internal, and Private endpoints
- **Redis caching**: Token caching and session management
- **PostgreSQL**: Persistent data storage
- **CLI Clients**: Admin, App, and Anonymous command-line interfaces

## Getting Started

See the [Getting Started](./development/getting-started.md) guide for setup instructions.

## Documentation Structure

- **Architecture**: System design and component overview
- **API Reference**: REST API documentation
- **Configuration**: Setup and configuration guides
- **Development**: Build, test, and contributing workflows
- **Appendix**: Swagger, GoDoc
