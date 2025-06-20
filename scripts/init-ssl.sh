#!/bin/sh

# SSL Certificate Initialization Script
# This script ensures SSL certificates exist before nginx starts

set -e

SSL_DIR="/ssl"
CERT_FILE="$SSL_DIR/nginx.crt"
KEY_FILE="$SSL_DIR/nginx.key"

echo "=== SSL Certificate Initialization ==="

# Create SSL directory if it doesn't exist
mkdir -p "$SSL_DIR"

# Check if certificates already exist
if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    echo "SSL certificates already exist"
    echo "   Certificate: $CERT_FILE"
    echo "   Private Key: $KEY_FILE"
    
    # Verify certificate is valid
    if openssl x509 -in "$CERT_FILE" -noout -text >/dev/null 2>&1; then
        echo "Certificate is valid"
        exit 0
    else
        echo "Certificate is invalid, regenerating..."
        rm -f "$CERT_FILE" "$KEY_FILE"
    fi
fi

echo "Generating new SSL certificates..."

# Install OpenSSL if not available
if ! command -v openssl >/dev/null 2>&1; then
    echo "Installing OpenSSL..."
    apk add --no-cache openssl
fi

# Generate self-signed certificate
echo "Creating self-signed certificate..."
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$KEY_FILE" \
    -out "$CERT_FILE" \
    -subj "/C=UA/ST=Kyiv/L=Kyiv/O=DevelopsToday/OU=SpyCats/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,DNS:*.localhost,IP:127.0.0.1,IP:0.0.0.0"

# Set proper permissions
chmod 644 "$CERT_FILE"
chmod 600 "$KEY_FILE"

# Verify generated certificate
if openssl x509 -in "$CERT_FILE" -noout -text >/dev/null 2>&1; then
    echo "   SSL certificates generated successfully!"
    echo "   Certificate: $CERT_FILE"
    echo "   Private Key: $KEY_FILE"
    
    # Show certificate info
    echo " Certificate Information:"
    openssl x509 -in "$CERT_FILE" -noout -subject -dates
else
    echo " Failed to generate valid SSL certificate"
    exit 1
fi

echo "=== SSL Initialization Complete ==="
