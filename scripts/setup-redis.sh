#!/usr/bin/env bash
set -euo pipefail

# Configure Redis on Alibaba Cloud Linux (dnf/yum install)
# Usage: sudo ./scripts/setup-redis.sh <redis_password>

REDIS_PASSWORD="${1:-}"
if [[ -z "$REDIS_PASSWORD" ]]; then
  echo "Usage: sudo $0 <redis_password>"
  exit 1
fi

REDIS_CONF="/etc/redis.conf"
if [[ ! -f "$REDIS_CONF" ]]; then
  echo "[ERROR] Redis config not found at $REDIS_CONF"
  exit 1
fi

echo "[INFO] Updating $REDIS_CONF ..."
cp "$REDIS_CONF" "${REDIS_CONF}.bak"

if grep -qE '^\s*requirepass\s+' "$REDIS_CONF"; then
  sed -i "s|^\s*requirepass\s\+.*|requirepass ${REDIS_PASSWORD}|" "$REDIS_CONF"
else
  echo "requirepass ${REDIS_PASSWORD}" >> "$REDIS_CONF"
fi

if grep -qE '^\s*bind\s+' "$REDIS_CONF"; then
  sed -i "s|^\s*bind\s\+.*|bind 127.0.0.1|" "$REDIS_CONF"
else
  echo "bind 127.0.0.1" >> "$REDIS_CONF"
fi

systemctl restart redis
systemctl enable redis

echo "[INFO] Redis configured and restarted."
