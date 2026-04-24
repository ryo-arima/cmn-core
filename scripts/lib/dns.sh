#!/usr/bin/env bash
# DNS Verification functions

# SCRIPT_DIR and ROOT_DIR are inherited from main.sh

function dns_check() {
    info "Checking DNS records for Locky mail server"
    
    local DNS_SERVER="10.88.0.53"
    local DOMAIN="cmn.local"
    local MAIL_SERVER="mail.cmn.local"
    
    local GREEN='\033[0;32m'
    local RED='\033[0;31m'
    local YELLOW='\033[1;33m'
    local NC='\033[0m'
    
    if ! docker ps | grep -q "locky-dns"; then
        err "DNS container is not running"
        info "Start with: ./scripts/main.sh env up"
        return 1
    fi
    
    success "DNS container is running"
    
    echo ""
    echo "======================================"
    echo "1. A Record (Address Record)"
    echo "======================================"
    echo "Purpose: Maps domain name to IP address"
    echo ""
    
    echo "Checking A record for $DOMAIN..."
    local RESULT
    RESULT=$(docker run --rm --network locky_locky-network alpine nslookup "$DOMAIN" "$DNS_SERVER" 2>/dev/null | grep "Address" | tail -n 1 | awk '{print $2}')
    if [ "$RESULT" = "10.88.0.10" ]; then
        echo -e "${GREEN}✓${NC} A record for $DOMAIN -> 10.88.0.10"
    else
        echo -e "${RED}✗${NC} A record incorrect (expected 10.88.0.10, got $RESULT)"
    fi
    
    echo ""
    echo "Checking A record for $MAIL_SERVER..."
    RESULT=$(docker run --rm --network locky_locky-network alpine nslookup "$MAIL_SERVER" "$DNS_SERVER" 2>/dev/null | grep "Address" | tail -n 1 | awk '{print $2}')
    if [ "$RESULT" = "10.88.0.10" ]; then
        echo -e "${GREEN}✓${NC} A record for $MAIL_SERVER -> 10.88.0.10"
    else
        echo -e "${RED}✗${NC} A record incorrect (expected 10.88.0.10, got $RESULT)"
    fi
    
    echo ""
    echo "======================================"
    echo "2. MX Record (Mail Exchange)"
    echo "======================================"
    echo "Purpose: Specifies mail server for domain"
    echo ""
    
    echo "Checking MX record for $DOMAIN..."
    local MX_RESULT
    MX_RESULT=$(docker run --rm --network locky_locky-network alpine sh -c "apk add --no-cache bind-tools >/dev/null 2>&1 && dig @$DNS_SERVER $DOMAIN MX +short" 2>/dev/null)
    if echo "$MX_RESULT" | grep -q "mail.cmn.local"; then
        echo -e "${GREEN}✓${NC} MX record found:"
        echo "   $MX_RESULT"
    else
        echo -e "${RED}✗${NC} MX record not found"
        echo "   Expected: 10 mail.cmn.local."
    fi
    
    echo ""
    echo "======================================"
    echo "3. SPF Record (TXT)"
    echo "======================================"
    echo "Purpose: Prevents email spoofing"
    echo ""
    
    echo "Checking SPF record for $DOMAIN..."
    local SPF_RESULT
    SPF_RESULT=$(docker run --rm --network locky_locky-network alpine sh -c "apk add --no-cache bind-tools >/dev/null 2>&1 && dig @$DNS_SERVER $DOMAIN TXT +short" 2>/dev/null | grep "spf1")
    if [ -n "$SPF_RESULT" ]; then
        echo -e "${GREEN}✓${NC} SPF record found:"
        echo "   $SPF_RESULT"
    else
        echo -e "${RED}✗${NC} SPF record not found"
        echo "   Expected: \"v=spf1 mx a ip4:10.88.0.10 ~all\""
    fi
    
    echo ""
    echo "======================================"
    echo "4. DMARC Record (TXT)"
    echo "======================================"
    echo "Purpose: Email authentication policy"
    echo ""
    
    echo "Checking DMARC record for _dmarc.$DOMAIN..."
    local DMARC_RESULT
    DMARC_RESULT=$(docker run --rm --network locky_locky-network alpine sh -c "apk add --no-cache bind-tools >/dev/null 2>&1 && dig @$DNS_SERVER _dmarc.$DOMAIN TXT +short" 2>/dev/null)
    if [ -n "$DMARC_RESULT" ]; then
        echo -e "${GREEN}✓${NC} DMARC record found:"
        echo "   $DMARC_RESULT"
    else
        echo -e "${YELLOW}⚠${NC} DMARC record not found (optional but recommended)"
        echo "   Expected: \"v=DMARC1; p=none; rua=mailto:dmarc@cmn.local\""
    fi
    
    echo ""
    echo "======================================"
    echo "Summary"
    echo "======================================"
    echo ""
    echo "DNS Server: $DNS_SERVER (locky-dns container)"
    echo "Domain: $DOMAIN"
    echo "Mail Server: $MAIL_SERVER"
    echo "IP Address: 10.88.0.10"
    echo ""
    success "DNS check complete"
}
