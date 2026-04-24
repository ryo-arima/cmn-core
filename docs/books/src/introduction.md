# Introduction

**cmn-core** is a robust Role-Based Access Control (RBAC) service built with Go, designed to provide comprehensive user, group, member, and role management with fine-grained permissions.

## Overview

cmn-core implements a multi-layered architecture that combines:
- **JWT-based authentication** for secure user identification
- **Casbin RBAC** for flexible authorization policies
- **Redis caching** for performance optimization
- **MySQL/TiDB** for persistent data storage

## Key Features

- **User Management**: Complete CRUD operations for user accounts
- **Group Management**: Organize users into logical groups
- **Member Management**: Control group membership and relationships
- **Role Management**: Define and assign fine-grained permissions
- **Multi-tier API**: Public, Internal, and Private endpoints with different access levels
- **Casbin Integration**: Policy-based authorization with dual enforcer setup
  - App-wide permissions (cmn-core policies)
  - Resource-specific permissions (resource policies)

## Architecture Highlights

cmn-core follows a clean architecture pattern with clear separation of concerns:

1. **Client Layer**: CLI tools for different user roles (Admin, App, Anonymous)
2. **API Layer**: Gin-based HTTP router with middleware stack
3. **Controller Layer**: Request handlers organized by access level
4. **Business Logic Layer**: Use case implementations
5. **Repository Layer**: Data access abstraction
6. **Data Layer**: MySQL/TiDB, Redis, and Casbin policy storage

## Use Cases

- **Multi-tenant applications** requiring organization-level access control
- **Enterprise systems** with complex permission requirements
- **Microservices** needing centralized authentication and authorization
- **APIs** requiring different access levels (public, internal, private)

## Getting Started

See the [Getting Started](./development/getting-started.md) guide for installation and setup instructions.

## Documentation Structure

This documentation is organized into the following sections:

- **Architecture**: Deep dive into system design and components
- **API Reference**: Complete API documentation with examples
- **Configuration**: Setup and configuration guides
- **Development**: Contributing and development workflows
- **Appendix**: Additional resources (Swagger, GoDoc)
