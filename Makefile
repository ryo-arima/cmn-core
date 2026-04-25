.PHONY: s bootstrap build dev-up dev-down dev-api docs localstack logs test clean

# ---------------------------------------------------------------------------
# Docker Compose — all service files merged via -f flags
# ---------------------------------------------------------------------------
COMPOSE = docker compose --project-directory . \
	-f docker/network.yaml \
	-f docker/dns.yaml \
	-f docker/postgres.yaml \
	-f docker/pgadmin.yaml \
	-f docker/redis.yaml \
	-f docker/keycloak.yaml \
	-f docker/casdoor.yaml \
	-f docker/mailserver.yaml \
	-f docker/roundcube.yaml \
	-f docker/server.yaml \
	-f docker/swagger.yaml \
	-f docker/godoc.yaml \
	-f docker/localstack.yaml

# Git shortcut
s:
	git add .
	commit-emoji
	git push origin main

# Build binaries into .bin/
build:
	mkdir -p .bin
	go build -o .bin/admin-client    ./cmd/client/admin
	go build -o .bin/app-client      ./cmd/client/app
	go build -o .bin/anonymous-client ./cmd/client/anonymous
	go build -o .bin/server   ./cmd/server

# Development environment (core infra + IdP + mail + app server)
dev-up:
	$(COMPOSE) up -d postgres redis pgadmin \
		keycloak \
		casdoor \
		dns mailserver roundcube \
		server

dev-down:
	$(COMPOSE) down -v --remove-orphans

# Documentation (swagger-ui + godoc)
docs:
	$(COMPOSE) up -d swagger-ui godoc
	@echo "Swagger UI: http://localhost:3002"
	@echo "Go Docs:    http://localhost:3003"

# AWS mock (localstack)
localstack:
	$(COMPOSE) up -d localstack

# Logs (pass SERVICE= to filter, e.g. make logs SERVICE=postgres)
logs:
	$(COMPOSE) logs -f $(SERVICE)

# Tests
test:
	go test ./...

test-unit:
	go test ./test/unit/...

clean:
	rm -f coverage.out coverage.html

