#!/bin/bash
# setup-ssl.sh - SSL certificate setup using Let's Encrypt
# Usage: ./setup-ssl.sh <domain> <email>

set -e

DOMAIN=${1:-erp.company.vn}
EMAIL=${2:-admin@company.vn}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[SSL Setup]${NC} $1"; }
warn() { echo -e "${YELLOW}[SSL Setup] WARNING:${NC} $1"; }
error() { echo -e "${RED}[SSL Setup] ERROR:${NC} $1"; exit 1; }

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    error "Please run as root (sudo ./setup-ssl.sh)"
fi

log "Setting up SSL for domain: $DOMAIN"
log "Admin email: $EMAIL"

# Install certbot if not present
if ! command -v certbot &> /dev/null; then
    log "Installing certbot..."
    apt update
    apt install -y certbot python3-certbot-nginx
fi

# Stop nginx temporarily if running
if systemctl is-active --quiet nginx; then
    log "Stopping nginx temporarily..."
    systemctl stop nginx
fi

# Obtain certificate
log "Obtaining SSL certificate..."
certbot certonly --standalone \
    -d "$DOMAIN" \
    --non-interactive \
    --agree-tos \
    --email "$EMAIL" \
    --no-eff-email

# Check if certificate was obtained
CERT_PATH="/etc/letsencrypt/live/$DOMAIN"
if [ ! -d "$CERT_PATH" ]; then
    error "Certificate was not obtained. Check certbot logs."
fi

log "Certificate obtained successfully!"
log "Certificate path: $CERT_PATH"

# Create nginx SSL config
log "Creating nginx SSL configuration..."
cat > /etc/nginx/snippets/ssl-$DOMAIN.conf << EOF
ssl_certificate $CERT_PATH/fullchain.pem;
ssl_certificate_key $CERT_PATH/privkey.pem;
ssl_session_timeout 1d;
ssl_session_cache shared:SSL:50m;
ssl_session_tickets off;

# Modern configuration
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
ssl_prefer_server_ciphers off;

# HSTS
add_header Strict-Transport-Security "max-age=63072000; includeSubDomains" always;
EOF

# Setup auto-renewal cron job
log "Setting up auto-renewal..."
if ! crontab -l 2>/dev/null | grep -q "certbot renew"; then
    (crontab -l 2>/dev/null; echo "0 0 1 * * certbot renew --quiet --post-hook 'systemctl reload nginx'") | crontab -
    log "Auto-renewal cron job added (monthly)"
fi

# Restart nginx
if [ -f /etc/nginx/nginx.conf ]; then
    log "Starting nginx..."
    systemctl start nginx
    systemctl enable nginx
fi

log "========================================="
log "SSL Setup Complete!"
log "========================================="
log "Certificate: $CERT_PATH/fullchain.pem"
log "Private Key: $CERT_PATH/privkey.pem"
log "Nginx snippet: /etc/nginx/snippets/ssl-$DOMAIN.conf"
log ""
log "Add to your nginx server block:"
log "  listen 443 ssl http2;"
log "  include snippets/ssl-$DOMAIN.conf;"
