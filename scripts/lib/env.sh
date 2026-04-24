#!/usr/bin/env bash
# Environment Management functions

# SCRIPT_DIR and ROOT_DIR are inherited from main.sh

function env_recreate() {
    info "Recreating all containers, networks and volumes (ephemeral test env)"
    cd "$ROOT_DIR"
    docker compose down -v --remove-orphans || true

    # Deep cleanup scoped by compose project label
    PROJECT_LABEL="com.docker.compose.project=locky"
    info "Removing leftover containers (label: ${PROJECT_LABEL})"
    docker ps -a --filter label="$PROJECT_LABEL" -q | xargs -r docker rm -f || true

    info "Removing leftover volumes (label: ${PROJECT_LABEL})"
    docker volume ls --filter label="$PROJECT_LABEL" -q | xargs -r docker volume rm -f || true

    info "Removing leftover networks (label: ${PROJECT_LABEL})"
    docker network ls --filter label="$PROJECT_LABEL" -q | xargs -r docker network rm || true

    docker network rm locky_locky-network 2>/dev/null || true

    info "Validating docker compose config"
    docker compose config >/dev/null

    info "Bringing up dns and mysql (force recreate)"
    docker compose up -d --force-recreate dns mysql

    info "Waiting for MySQL to be ready..."
    ATTEMPTS=60
    until docker compose exec -T mysql sh -lc "mysql -uroot -pmysql -e 'SELECT 1'" >/dev/null 2>&1; do
        ATTEMPTS=$((ATTEMPTS-1)) || true
        if [ "$ATTEMPTS" -le 0 ]; then
            err "MySQL did not become ready in time"
            exit 1
        fi
        sleep 1
    done
    success "MySQL is ready"

    info "Bringing up mailserver and roundcube (force recreate)"
    docker compose up -d --force-recreate mailserver roundcube

    sleep 5 || true

    info "Ensuring mail accounts exist (idempotent)"
    mail_create_accounts || warn "create-accounts returned non-zero; continuing"

    info "Applying Postfix runtime overrides"
    set +e
    docker compose exec -T mailserver postconf -P submission/inet/smtpd_client_restrictions='permit_mynetworks, permit_sasl_authenticated, reject'
    docker compose exec -T mailserver postconf -P submission/inet/smtpd_recipient_restrictions='permit_sasl_authenticated, permit_mynetworks, reject_unauth_destination'
    docker compose exec -T mailserver postconf -e 'smtpd_sasl_auth_enable = yes'
    docker compose exec -T mailserver postconf -e 'smtpd_tls_auth_only = no'
    docker compose exec -T mailserver postconf -e 'mynetworks = 127.0.0.0/8 [::1]/128 10.88.0.0/24'
    set -e
    docker compose exec -T mailserver postfix reload || true
    success "Postfix overrides applied"

    info "Current service status:"
    docker compose ps || true

    success "Full recreate completed."
    echo "- Roundcube: http://localhost:3005"
    echo "- Example login: test1@cmn.local / TestPassword123!"
}

function env_up() {
    info "Starting all services (detached)"
    cd "$ROOT_DIR"
    docker compose up -d
    success "Services started"
    docker compose ps
}

function env_down() {
    info "Stopping all services"
    cd "$ROOT_DIR"
    docker compose down
    success "Services stopped"
}

function env_status() {
    cd "$ROOT_DIR"
    docker compose ps
}

function env_logs() {
    local service="${1:-}"
    cd "$ROOT_DIR"
    if [ -z "$service" ]; then
        docker compose logs -f
    else
        docker compose logs -f "$service"
    fi
}
