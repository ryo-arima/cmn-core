.PHONY: s bootstrap build dev-up dev-down dev-api docs localstack logs test clean

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

# Development environment (core + IdP + mail; docs/localstack excluded)
dev-up:
	docker compose up -d postgres redis pgadmin \
		keycloak \
		casdoor \
		dns mailserver roundcube

dev-down:
	docker compose down -v --remove-orphans

# Documentation (swagger-ui + godoc)
docs:
	docker compose up -d swagger-ui godoc
	@echo "Swagger UI: http://localhost:3002"
	@echo "Go Docs:    http://localhost:3003"

# AWS mock (localstack)
localstack:
	docker compose up -d localstack

# Logs (pass SERVICE= to filter, e.g. make logs SERVICE=postgres)
logs:
	docker compose logs -f $(SERVICE)

# Tests
test:
	go test ./...

test-unit:
	go test ./test/unit/...

clean:
	rm -f coverage.out coverage.html



