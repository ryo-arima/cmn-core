#!/usr/bin/env bash
# Help function

function help() {
    cat <<EOF
Locky Environment Management CLI

USAGE:
    ./scripts/main.sh <command> [args...]

COMMANDS:
    Environment Management (env):
        env recreate          - Recreate entire environment (ephemeral)
        env up                - Start all services
        env down              - Stop all services
        env status            - Show service status
        env logs [service]    - Show logs (all or specific service)

    Mail Management (mail):
        mail create_accounts  - Create/update mail accounts in running container
        mail generate_accounts- Generate postfix-accounts.cf with hashed passwords

    CI & Testing (ci):
        ci run [suite]        - Run CI orchestration (default: auth)
        ci matrix             - Run all test suites sequentially

    Documentation (docs):
        docs build            - Build GitHub Pages documentation
        docs serve            - Serve documentation locally (port 8000)
        docs architecture     - Generate architecture diagram from mermaid

    DNS Verification (dns):
        dns check             - Verify DNS records configuration

    Secrets Management (secrets):
        secrets init          - Initialize LocalStack Secrets Manager

    Helpers:
        help                  - Show this help message

EXAMPLES:
    # Full environment rebuild
    ./scripts/main.sh env recreate

    # Start services
    ./scripts/main.sh env up

    # Create mail accounts
    ./scripts/main.sh mail create_accounts

    # Run auth test suite
    ./scripts/main.sh ci run auth

    # Build docs
    ./scripts/main.sh docs build

    # Generate architecture diagram
    ./scripts/main.sh docs architecture

    # Check DNS
    ./scripts/main.sh dns check

    # Initialize secrets
    ./scripts/main.sh secrets init

EOF
}
