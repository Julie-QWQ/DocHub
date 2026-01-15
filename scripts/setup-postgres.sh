#!/usr/bin/env bash
set -euo pipefail

# Configure PostgreSQL on Alibaba Cloud Linux (dnf/yum install)
# Usage: sudo ./scripts/setup-postgres.sh <db_user> <db_password> <db_name>

DB_USER="${1:-}"
DB_PASSWORD="${2:-}"
DB_NAME="${3:-}"

if [[ -z "$DB_USER" || -z "$DB_PASSWORD" || -z "$DB_NAME" ]]; then
  echo "Usage: sudo $0 <db_user> <db_password> <db_name>"
  exit 1
fi

if ! command -v psql >/dev/null 2>&1; then
  echo "[ERROR] psql not found. Install PostgreSQL first."
  exit 1
fi

if ! systemctl is-active --quiet postgresql; then
  echo "[INFO] PostgreSQL not running. Attempting init and start..."
  if command -v postgresql-setup >/dev/null 2>&1; then
    postgresql-setup --initdb || true
  fi
  systemctl start postgresql
fi

systemctl enable postgresql

echo "[INFO] Creating user/database (if not exists)..."
sudo -u postgres psql -v ON_ERROR_STOP=1 <<SQL
DO \$\$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '${DB_USER}') THEN
    CREATE ROLE ${DB_USER} LOGIN PASSWORD '${DB_PASSWORD}';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}') THEN
    CREATE DATABASE ${DB_NAME} OWNER ${DB_USER};
  END IF;
END
\$\$;
GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};
SQL

echo "[INFO] PostgreSQL initialized."
echo "Next: edit pg_hba.conf/postgresql.conf if you need remote access."
