#!/bin/bash
# Tambahkan atau update akun admin via docker compose exec
# Usage  : bash scripts/add-admin.sh [email] [password] [username]
# Example: bash scripts/add-admin.sh admin@example.com secret123 adminuser

set -euo pipefail

EMAIL="${1:-mrg@raffimrg.my.id}"
PASSWORD="${2:-123123123}"
USERNAME="${3:-mrg}"

# Lokasi project (direktori di atas scripts/)
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
ENV_FILE="$PROJECT_DIR/.env"

if [ ! -f "$ENV_FILE" ]; then
  echo "ERROR: file .env tidak ditemukan di $ENV_FILE"
  exit 1
fi

DB_NAME=$(grep '^DB_NAME=' "$ENV_FILE" | cut -d= -f2)
DB_ROOT_PASSWORD=$(grep '^DB_ROOT_PASSWORD=' "$ENV_FILE" | cut -d= -f2)
EXTERNAL_API_PORT=$(grep '^EXTERNAL_API_PORT=' "$ENV_FILE" | cut -d= -f2)
API="http://localhost:${EXTERNAL_API_PORT}/api"

echo "=== Add Admin Script ==="
echo "Email    : $EMAIL"
echo "Username : $USERNAME"
echo "Database : $DB_NAME"
echo ""

# Step 1: Cek apakah user sudah ada
EXISTING=$(docker compose -f "$PROJECT_DIR/docker-compose.yml" exec -T mysql \
  mysql -u root -p"${DB_ROOT_PASSWORD}" -N -s \
  -e "SELECT COUNT(*) FROM ${DB_NAME}.users WHERE email='${EMAIL}';" 2>/dev/null)

if [ "$EXISTING" = "0" ]; then
  echo "→ Akun belum ada, register via API..."
  HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$USERNAME\",\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

  if [ "$HTTP_CODE" = "201" ]; then
    echo "✓ Akun berhasil dibuat"
  else
    echo "✗ Register gagal (HTTP $HTTP_CODE) — pastikan backend berjalan"
    exit 1
  fi
else
  echo "→ Akun sudah ada (skip register)"
fi

# Step 2: Set role admin via docker compose exec mysql
echo "→ Update role ke 'admin'..."
docker compose -f "$PROJECT_DIR/docker-compose.yml" exec -T mysql \
  mysql -u root -p"${DB_ROOT_PASSWORD}" \
  -e "UPDATE ${DB_NAME}.users SET role='admin' WHERE email='${EMAIL}';" 2>/dev/null

# Verifikasi
RESULT=$(docker compose -f "$PROJECT_DIR/docker-compose.yml" exec -T mysql \
  mysql -u root -p"${DB_ROOT_PASSWORD}" -N -s \
  -e "SELECT CONCAT(username, ' | ', email, ' | role:', role) FROM ${DB_NAME}.users WHERE email='${EMAIL}';" 2>/dev/null)

echo ""
echo "✓ Selesai: $RESULT"
