#!/usr/bin/env bash
# Unified CLI for Locky environment management
# Usage: ./scripts/main.sh <command> [args...]
# Example: ./scripts/main.sh env recreate
#          ./scripts/main.sh mail create_accounts
#          ./scripts/main.sh ci run auth

set -euo pipefail

COMMAND="${1:-help}"
shift || true
ARGS="$@"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

# Load library functions
source "$SCRIPT_DIR/lib/common.sh"
source "$SCRIPT_DIR/lib/env.sh"
source "$SCRIPT_DIR/lib/mail.sh"
source "$SCRIPT_DIR/lib/ci.sh"
source "$SCRIPT_DIR/lib/docs.sh"
source "$SCRIPT_DIR/lib/dns.sh"
source "$SCRIPT_DIR/lib/secrets.sh"
source "$SCRIPT_DIR/lib/help.sh"

# ============================================================================
# Command Router
# ============================================================================
case "$COMMAND" in
    # Environment commands
    env)
        SUBCOMMAND="${1:-help}"
        shift || true
        case "$SUBCOMMAND" in
            recreate) env_recreate ;;
            up) env_up ;;
            down) env_down ;;
            status) env_status ;;
            logs) env_logs "$@" ;;
            *) err "Unknown env subcommand: $SUBCOMMAND"; help; exit 1 ;;
        esac
        ;;
    
    # Mail commands
    mail)
        SUBCOMMAND="${1:-help}"
        shift || true
        case "$SUBCOMMAND" in
            create_accounts) mail_create_accounts ;;
            generate_accounts) mail_generate_accounts ;;
            *) err "Unknown mail subcommand: $SUBCOMMAND"; help; exit 1 ;;
        esac
        ;;
    
    # CI commands
    ci)
        SUBCOMMAND="${1:-help}"
        shift || true
        case "$SUBCOMMAND" in
            run) ci_run "$@" ;;
            matrix) ci_matrix ;;
            *) err "Unknown ci subcommand: $SUBCOMMAND"; help; exit 1 ;;
        esac
        ;;
    
    # Docs commands
    docs)
        SUBCOMMAND="${1:-help}"
        shift || true
        case "$SUBCOMMAND" in
            build) docs_build ;;
            serve) docs_serve ;;
            architecture) docs_architecture ;;
            *) err "Unknown docs subcommand: $SUBCOMMAND"; help; exit 1 ;;
        esac
        ;;
    
    # DNS commands
    dns)
        SUBCOMMAND="${1:-help}"
        shift || true
        case "$SUBCOMMAND" in
            check) dns_check ;;
            *) err "Unknown dns subcommand: $SUBCOMMAND"; help; exit 1 ;;
        esac
        ;;
    
    # Secrets commands
    secrets)
        SUBCOMMAND="${1:-help}"
        shift || true
        case "$SUBCOMMAND" in
            init) secrets_init ;;
            *) err "Unknown secrets subcommand: $SUBCOMMAND"; help; exit 1 ;;
        esac
        ;;
    
    # Help
    help|--help|-h|"")
        help
        ;;
    
    *)
        err "Unknown command: $COMMAND"
        help
        exit 1
        ;;
esac
