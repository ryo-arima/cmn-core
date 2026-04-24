#!/usr/bin/env bash
# CI / Testing functions

# SCRIPT_DIR and ROOT_DIR are inherited from main.sh

function ci_build_images() {
    info "Building images (if changes present)..."
    cd "$ROOT_DIR"
    docker compose -f "${COMPOSE_FILE:-docker-compose.yaml}" build --pull || err "Build step failed"
}

function ci_bring_up() {
    info "Starting services (detached)..."
    cd "$ROOT_DIR"
    docker compose -f "${COMPOSE_FILE:-docker-compose.yaml}" up -d --quiet-pull --remove-orphans
}

function ci_list_containers() {
    cd "$ROOT_DIR"
    docker compose -f "${COMPOSE_FILE:-docker-compose.yaml}" ps -q
}

function ci_wait_health() {
    local timeout=${HEALTH_TIMEOUT:-300}
    info "Waiting for containers to become healthy or running (timeout=${timeout}s)..."
    cd "$ROOT_DIR"
    local end=$((SECONDS+timeout))
    mapfile -t containers < <(ci_list_containers)
    if [[ ${#containers[@]} -eq 0 ]]; then
        err "No containers found after up"; return 1
    fi
    for id in "${containers[@]}"; do
        local name
        name=$(docker inspect -f '{{.Name}}' "$id" | sed 's#^/##')
        info "→ Waiting for $name ($id)"
        while true; do
            local health state
            health=$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "$id" || echo "missing")
            state=$(docker inspect -f '{{.State.Status}}' "$id" || echo "missing")
            if [[ "$health" == "healthy" || ( "$health" == "none" && "$state" == "running" ) ]]; then
                info "   Ready: $name (health=$health state=$state)"
                break
            fi
            if (( SECONDS > end )); then
                err "Timeout: $name (health=$health state=$state)"
                docker logs "$id" || true
                return 1
            fi
            sleep 3
        done
    done
}

function ci_run_tests() {
    local suite="${SUITE:-auth}"
    local test_auth_primary="$ROOT_DIR/test/integration_auth_test.sh"
    local test_auth_fallback="$ROOT_DIR/test/client_auth_integration_test.sh"
    local test_mailflow="$ROOT_DIR/test/mailflow_integration_test.sh"
    
    cd "$ROOT_DIR"
    case "$suite" in
        auth)
            if [[ -f "$test_auth_primary" ]]; then
                info "Running auth test: $test_auth_primary"
                bash "$test_auth_primary"
            elif [[ -f "$test_auth_fallback" ]]; then
                info "Auth primary not found; fallback: $test_auth_fallback"
                bash "$test_auth_fallback"
            else
                warn "No auth test script found; skipping"
            fi
            ;;
        mailflow)
            if [[ -f "$test_mailflow" ]]; then
                info "Running mailflow test: $test_mailflow"
                bash "$test_mailflow"
            else
                warn "Mailflow test script missing; skipping"
            fi
            ;;
        *)
            warn "Unknown SUITE='$suite' — skipping tests"
            ;;
    esac
}

function ci_collect_logs() {
    local log_dir="${LOG_DIR:-artifacts}"
    info "Collecting logs & state..."
    cd "$ROOT_DIR"
    mkdir -p "$log_dir"
    docker compose -f "${COMPOSE_FILE:-docker-compose.yaml}" ps > "$log_dir/ps.txt" 2>&1 || true
    docker compose -f "${COMPOSE_FILE:-docker-compose.yaml}" logs --no-color > "$log_dir/compose.log" 2>&1 || true
    while read -r id; do
        [[ -z "$id" ]] && continue
        local name
        name=$(docker inspect -f '{{.Name}}' "$id" | sed 's#^/##')
        docker logs "$id" --since 1h > "$log_dir/${name}.log" 2>&1 || true
    done < <(ci_list_containers)
}

function ci_teardown() {
    info "Tearing down stack (volumes removed)..."
    cd "$ROOT_DIR"
    docker compose -f "${COMPOSE_FILE:-docker-compose.yaml}" down -v --remove-orphans || err "Down failed"
}

function ci_parallel() {
    local suite="${1:-auth}"
    export COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yaml}"
    export PROJECT_NAME="${COMPOSE_PROJECT_NAME:-locky-local}"
    export HEALTH_TIMEOUT="${HEALTH_TIMEOUT:-300}"
    export LOG_DIR="${LOG_DIR:-artifacts}"
    export SUITE="$suite"

    info "Project: $PROJECT_NAME (suite=$SUITE)"
    
    # Set trap for cleanup
    trap 'warn "Script terminating; capturing logs then teardown"; ci_collect_logs; ci_teardown' EXIT
    
    ci_build_images
    ci_bring_up
    ci_wait_health
    ci_run_tests
}

function ci_run() {
    local suite="${1:-auth}"
    export COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yaml}"
    export PROJECT_NAME="${COMPOSE_PROJECT_NAME:-locky-local}"
    export HEALTH_TIMEOUT="${HEALTH_TIMEOUT:-300}"
    export LOG_DIR="${LOG_DIR:-artifacts}"
    export SUITE="$suite"

    info "Running CI suite: $SUITE (project=$PROJECT_NAME)"
    ci_parallel "$suite"
}

function ci_matrix() {
    info "Running all test suites in sequence (auth + mailflow)"
    ci_run auth
    ci_run mailflow
    success "All test suites completed"
}
