# Git shortcuts
s:
	git add .
	commit-emoji
	git push origin main

bootstrap:
	go run cmd/client/admin/main.go bootstrap user
	go run cmd/client/admin/main.go bootstrap group
	go run cmd/client/admin/main.go bootstrap member

# Development environment
.PHONY: dev dev-up dev-down dev-logs

# Start development environment with all services
dev-up:
	@echo "Starting development environment..."
	@docker-compose up -d mysql redis
	@sleep 5
	@echo "Database and Redis are ready!"
	@echo "Starting documentation services..."
	@docker-compose up -d swagger-ui godoc
	@echo "Development environment is ready!"
	@echo "  - Swagger UI: http://localhost:3002"
	@echo "  - Go Docs: http://localhost:3003"
	@echo "  - phpMyAdmin: http://localhost:3001"
	@echo ""
	@echo "To start the API server manually:"
	@echo "  go run cmd/server/main.go"

# Stop development environment
dev-down:
	@echo "Stopping development environment..."
	@docker-compose down

# Full ephemeral mail stack recreate (dns, mysql, mailserver, roundcube)
.PHONY: mail-recreate
mail-recreate:
	@echo "Recreating full mail test environment (ephemeral) ..."
	@bash ./scripts/main.sh env recreate
	@echo "Done. Roundcube: http://localhost:3005"

# Show logs for development services
dev-logs:
	@docker-compose logs -f swagger-ui godoc

# Start local API server
dev-api:
	@echo "Starting local API server..."
	@go run cmd/server/main.go

# Documentation commands
.PHONY: docs docs-swagger docs-godoc docs-open

# Start documentation services only
docs:
	@echo "Starting documentation services..."
	@docker compose up -d swagger-ui godoc
	@echo "Documentation services started!"
	@echo "  - Swagger UI: http://localhost:3002"
	@echo "  - Go Docs: http://localhost:3003"

# Start Swagger UI only
docs-swagger:
	@echo "Starting Swagger UI..."
	@docker compose up -d swagger-ui
	@echo "Swagger UI: http://localhost:3002"

# Start Go documentation server only
docs-godoc:
	@echo "Starting Go documentation server..."
	@docker compose up -d godoc
	@echo "Go Docs: http://localhost:3003"

# Open documentation in browser (macOS)
docs-open:
	@echo "Opening documentation in browser..."
	@open http://localhost:3002 || echo "Please open http://localhost:3002 manually"
	@open http://localhost:3003 || echo "Please open http://localhost:3003 manually"

# Validate Swagger specification
docs-validate:
	@echo "Validating Swagger specification..."
	@docker run --rm -v $(PWD)/docs/swagger:/swagger swaggerapi/swagger-validator:latest /swagger/swagger.yaml

# Generate Go client from Swagger spec (requires swagger-codegen)
docs-generate-client:
	@echo "Generating Go client from Swagger specification..."
	@docker run --rm -v $(PWD):/local swaggerapi/swagger-codegen-cli:latest generate \
		-i /local/docs/swagger/swagger.yaml \
		-l go \
		-o /local/generated/client

# Simplified UML / PlantUML targets
DIAGRAM_DIR := docs/diagrams
PUML ?= $(DIAGRAM_DIR)/classes.puml
PNG ?= $(PUML:.puml=.png)

.PHONY: uml-png uml-png-all uml-clean uml-gen-code

## uml-png: Convert a single .puml (PUML=...) to .png
uml-png:
	@mkdir -p $(DIAGRAM_DIR)
	@if [ ! -f "$(PUML)" ]; then echo "PUML not found: $(PUML)"; exit 1; fi
	@echo "Rendering $(PUML) -> $(PNG)"
	@if which plantuml >/dev/null 2>&1; then \
		plantuml -tpng $(PUML); \
		mv $(DIAGRAM_DIR)/$$(basename $(PUML) .puml).png $(PNG) 2>/dev/null || true; \
		echo "Done (local plantuml)"; \
	else \
		echo "Local plantuml not found. Using Docker."; \
		docker run --rm -v $$(pwd)/$(DIAGRAM_DIR):/workspace plantuml/plantuml -tpng /workspace/$$(basename $(PUML)); \
		echo "Done (docker plantuml)"; \
	fi

## uml-png-all: Convert all .puml under docs/diagrams to .png
uml-png-all:
	@mkdir -p $(DIAGRAM_DIR)
	@set -e; \
	for f in $(DIAGRAM_DIR)/*.puml; do \
		[ -f "$$f" ] || continue; \
		echo "Rendering $$f"; \
		$(MAKE) -s uml-png PUML=$$f; \
	done

## uml-clean: Remove generated PNG files
uml-clean:
	@rm -f $(DIAGRAM_DIR)/*.png
	@echo "Removed PNG diagrams"

## uml-gen-code: (Optional) Generate classes.puml from Go code (kept minimal) // requires goplantuml
uml-gen-code:
	@mkdir -p $(DIAGRAM_DIR)
	@which goplantuml >/dev/null 2>&1 || { echo 'Installing goplantuml ...'; go install github.com/jfeliu007/goplantuml/cmd/goplantuml@latest; }
	@goplantuml -recursive -output $(PUML) ./pkg
	@echo "Generated $(PUML) from code"

# Test targets
.PHONY: test test-server test-repository test-controller test-coverage clean-test

test: test-server

test-server:
	@./test/run_tests.sh server

test-repository:
	@./test/run_tests.sh repository

test-controller:
	@./test/run_tests.sh controller

test-coverage:
	@./test/run_tests.sh coverage

# Clean test artifacts
clean-test:
	@echo "Cleaning test artifacts..."
	@rm -rf test/results/*
	@rm -f coverage.out coverage.html
	@echo "Done."

PKG_STRUCTURE_PUML := $(DIAGRAM_DIR)/pkg_structure.puml
PKG_STRUCTURE_PNG := $(DIAGRAM_DIR)/pkg_structure.png

.PHONY: gen-pkg-structure gen-pkg-structure-png

## gen-pkg-structure: Generate high-level (detailed) pkg/ directory structure PlantUML (revised v3)
## NOTE: .go 拡張子は除去し論理コンポーネント(class)名のみ使用
##       主要インタフェース/エンティティ/依存関係を追加
## Usage: make gen-pkg-structure
gen-pkg-structure:
	@mkdir -p $(DIAGRAM_DIR)
	@echo "@startuml" > $(PKG_STRUCTURE_PUML)
	@echo "title pkg directory structure (detailed)" >> $(PKG_STRUCTURE_PUML)
	@echo "' Generated by make gen-pkg-structure (revised v3)" >> $(PKG_STRUCTURE_PUML)
	@echo "skinparam packageStyle rectangle" >> $(PKG_STRUCTURE_PUML)
	@echo "skinparam defaultFontName Monospaced" >> $(PKG_STRUCTURE_PUML)
	@echo "skinparam shadowing false" >> $(PKG_STRUCTURE_PUML)
	@echo "hide empty members" >> $(PKG_STRUCTURE_PUML)
	@echo "' ==============================================" >> $(PKG_STRUCTURE_PUML)
	@echo "package pkg {" >> $(PKG_STRUCTURE_PUML)
	@echo "  package client {" >> $(PKG_STRUCTURE_PUML)
	@echo "    package controller as client_controller { class AdminCmd; class AppCmd; class AnonymousCmd }" >> $(PKG_STRUCTURE_PUML)
	@echo "    package usecase as client_usecase { class UserUsecase; class GroupUsecase; class MemberUsecase; class RoleUsecase; class CommonUsecase }" >> $(PKG_STRUCTURE_PUML)
	@echo "    package repository as client_repository {" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface CommonRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface UserRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface GroupRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface MemberRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface RoleRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "    }" >> $(PKG_STRUCTURE_PUML)
	@echo "  }" >> $(PKG_STRUCTURE_PUML)
	@echo "  package server {" >> $(PKG_STRUCTURE_PUML)
	@echo "    package controller as server_controller {" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface UserControllerForPublic" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface UserControllerForInternal" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface UserControllerForPrivate" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface GroupControllerForPrivate" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface GroupControllerForInternal" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface MemberControllerForPrivate" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface MemberControllerForInternal" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface RoleControllerForPrivate" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface RoleControllerForInternal" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface CommonControllerForPublic" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface CommonControllerForPrivate" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface CommonControllerForInternal" >> $(PKG_STRUCTURE_PUML)
	@echo "    }" >> $(PKG_STRUCTURE_PUML)
	@echo "    package middleware as server_middleware { class AuthMiddleware }" >> $(PKG_STRUCTURE_PUML)
	@echo "    package repository as server_repository {" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface CommonRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface UserRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface GroupRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface MemberRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "      interface RoleRepository" >> $(PKG_STRUCTURE_PUML)
	@echo "    }" >> $(PKG_STRUCTURE_PUML)
	@echo "  }" >> $(PKG_STRUCTURE_PUML)
	@echo "  package config as config_pkg { class BaseConfig }" >> $(PKG_STRUCTURE_PUML)
	@echo "  package entity {" >> $(PKG_STRUCTURE_PUML)
	@echo "    package model as entity_model { class Users; class Groups; class Members; class TokenPair; class JWTClaims }" >> $(PKG_STRUCTURE_PUML)
	@echo "    package request as entity_request { class UserRequest; class GroupRequest; class MemberRequest; class LoginRequest; class RolePermissionRequest }" >> $(PKG_STRUCTURE_PUML)
	@echo "    package response as entity_response { class UserResponse; class GroupResponse; class MemberResponse; class RoleResponse; class LoginResponse; class CountResponse }" >> $(PKG_STRUCTURE_PUML)
	@echo "  }" >> $(PKG_STRUCTURE_PUML)
	@echo "}" >> $(PKG_STRUCTURE_PUML)
	@echo "' ==============================================" >> $(PKG_STRUCTURE_PUML)
	@echo "' Logical dependency (simplified arrows)" >> $(PKG_STRUCTURE_PUML)
	@echo "client_usecase ..> client_repository : uses" >> $(PKG_STRUCTURE_PUML)
	@echo "client_controller ..> client_usecase : orchestrates" >> $(PKG_STRUCTURE_PUML)
	@echo "client_repository ..> server_controller : HTTP calls" >> $(PKG_STRUCTURE_PUML)
	@echo "server_controller ..> server_repository : invokes" >> $(PKG_STRUCTURE_PUML)
	@echo "server_controller ..> server_middleware : auth" >> $(PKG_STRUCTURE_PUML)
	@echo "server_repository ..> entity_model : persistence" >> $(PKG_STRUCTURE_PUML)
	@echo "server_repository ..> config_pkg : config" >> $(PKG_STRUCTURE_PUML)
	@echo "client_repository ..> config_pkg : config" >> $(PKG_STRUCTURE_PUML)
	@echo "server_controller ..> entity_request : bind" >> $(PKG_STRUCTURE_PUML)
	@echo "server_controller ..> entity_response : returns" >> $(PKG_STRUCTURE_PUML)
	@echo "client_usecase ..> entity_request : build" >> $(PKG_STRUCTURE_PUML)
	@echo "client_usecase ..> entity_response : format" >> $(PKG_STRUCTURE_PUML)
	@echo "@enduml" >> $(PKG_STRUCTURE_PUML)
	@echo "Done: $(PKG_STRUCTURE_PUML)"

## gen-pkg-structure-png: Render PNG for pkg structure
gen-pkg-structure-png: gen-pkg-structure
	@$(MAKE) gen-plantuml-png PUML=$(PKG_STRUCTURE_PUML) PNG=$(PKG_STRUCTURE_PNG)

# Detailed pkg structure with members & methods (leverages qualified full diagram)
PKG_DETAILED_PUML := $(DIAGRAM_DIR)/pkg_structure.puml
PKG_DETAILED_PNG := $(DIAGRAM_DIR)/pkg_structure.png

.PHONY: gen-pkg-structure-detailed gen-pkg-structure-detailed-png
## gen-pkg-structure-detailed: Generate pkg_structure.puml including all members & methods
## Implementation: reuse gen-plantuml-qualified output (classes_qualified.puml)
## Filters: none (full detail)
gen-pkg-structure-detailed: gen-plantuml-qualified
	@cp $(QUALIFIED_PUML) $(PKG_DETAILED_PUML)
	@echo "Copied $(QUALIFIED_PUML) -> $(PKG_DETAILED_PUML) (detailed)"

## gen-pkg-structure-detailed-png: Render PNG for detailed pkg structure
gen-pkg-structure-detailed-png: gen-pkg-structure-detailed
	@$(MAKE) gen-plantuml-png PUML=$(PKG_DETAILED_PUML) PNG=$(PKG_DETAILED_PNG)
