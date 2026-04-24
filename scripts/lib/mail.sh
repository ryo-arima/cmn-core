#!/usr/bin/env bash
# Mail Management functions

# SCRIPT_DIR and ROOT_DIR are inherited from main.sh

function mail_create_accounts() {
    local CONTAINER_NAME="locky-mailserver"

    if ! docker ps | grep -q "$CONTAINER_NAME"; then
        err "$CONTAINER_NAME is not running"
        return 1
    fi

    info "Creating mail accounts in $CONTAINER_NAME..."

    local ACCOUNTS=(
        "admin@cmn.local:AdminPassword123!"
        "test1@cmn.local:TestPassword123!"
        "test2@cmn.local:TestPassword123!"
        "test3@cmn.local:TestPassword123!"
        "test4@cmn.local:TestPassword123!"
        "test5@cmn.local:TestPassword123!"
        "user1@cmn.local:User1Password123!"
        "user2@cmn.local:User2Password123!"
        "developer@cmn.local:DevPassword123!"
        "noreply@cmn.local:NoReplyPassword123!"
        "support@cmn.local:SupportPassword123!"
    )

    info "Retrieving existing accounts..."
    local EXISTING
    EXISTING=$(docker exec "$CONTAINER_NAME" setup email list || true)
    echo "$EXISTING" | sed 's/^/  - /'

    for account in "${ACCOUNTS[@]}"; do
        local email="${account%%:*}"
        local password="${account#*:}"
        info "Ensuring account exists: $email"
        if docker exec "$CONTAINER_NAME" setup email update "$email" "$password" 2>&1 | grep -qi "does not exist"; then
            info "  -> creating (not found)"
            if ! docker exec "$CONTAINER_NAME" setup email add "$email" "$password"; then
                info "  -> add failed, retrying update"
                docker exec "$CONTAINER_NAME" setup email update "$email" "$password" || true
            fi
        else
            info "  -> password updated"
        fi
    done

    info "Final account list:"
    docker exec "$CONTAINER_NAME" setup email list | sed 's/^/  * /'

    success "All mail accounts created successfully!"
}

function mail_generate_accounts() {
    info "Generating postfix-accounts.cf with hashed passwords"
    local CONFIG_DIR="$ROOT_DIR/scripts/data/mailserver"
    local ACCOUNTS_FILE="$CONFIG_DIR/postfix-accounts.cf"

    python3 - <<'PYTHON'
import crypt
import os

accounts = {
    "admin@cmn.local": "AdminPassword123!",
    "test1@cmn.local": "TestPassword123!",
    "test2@cmn.local": "TestPassword123!",
    "test3@cmn.local": "TestPassword123!",
    "test4@cmn.local": "TestPassword123!",
    "test5@cmn.local": "TestPassword123!",
    "user1@cmn.local": "User1Password123!",
    "user2@cmn.local": "User2Password123!",
    "developer@cmn.local": "DevPassword123!",
    "noreply@cmn.local": "NoReplyPassword123!",
    "support@cmn.local": "SupportPassword123!"
}

config_dir = "scripts/data/mailserver"
accounts_file = os.path.join(config_dir, "postfix-accounts.cf")

os.makedirs(config_dir, exist_ok=True)

print("Generating password hashes for mail accounts...")

with open(accounts_file, 'w') as f:
    for email, password in sorted(accounts.items()):
        print(f"Processing: {email}")
        password_hash = crypt.crypt(password, crypt.mksalt(crypt.METHOD_SHA512))
        f.write(f"{email}|{password_hash}\n")

print()
print("✓ Password hashes generated successfully!")
print(f"✓ Configuration saved to: {accounts_file}")
PYTHON

    success "Restart mailserver to apply: docker compose restart mailserver"
}
