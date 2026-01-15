#!/usr/bin/env bash
set -euo pipefail

# Configure Nginx for Study-UPC
# Usage: sudo ./scripts/setup-nginx.sh <server_name> <frontend_dist_path> <backend_host> <backend_port>
# Example: sudo ./scripts/setup-nginx.sh example.com /opt/study-upc/frontend/dist 127.0.0.1 8080

SERVER_NAME="${1:-}"
FRONTEND_DIST="${2:-}"
BACKEND_HOST="${3:-}"
BACKEND_PORT="${4:-}"

if [[ -z "$SERVER_NAME" || -z "$FRONTEND_DIST" || -z "$BACKEND_HOST" || -z "$BACKEND_PORT" ]]; then
  echo "Usage: sudo $0 <server_name> <frontend_dist_path> <backend_host> <backend_port>"
  exit 1
fi

if [[ ! -d "$FRONTEND_DIST" ]]; then
  echo "[ERROR] Frontend dist path not found: $FRONTEND_DIST"
  exit 1
fi

NGINX_CONF="/etc/nginx/conf.d/study-upc.conf"

cat > "$NGINX_CONF" <<EOF
server {
    listen 80;
    server_name ${SERVER_NAME};

    root ${FRONTEND_DIST};
    index index.html;

    location / {
        try_files \$uri \$uri/ /index.html;
    }

    location /api/v1/ {
        proxy_pass http://${BACKEND_HOST}:${BACKEND_PORT}/api/v1/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}
EOF

nginx -t
systemctl reload nginx

echo "[INFO] Nginx configured: $NGINX_CONF"
