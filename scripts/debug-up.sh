#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if command -v docker-compose >/dev/null 2>&1; then
  COMPOSE_CMD="docker-compose"
elif docker compose version >/dev/null 2>&1; then
  COMPOSE_CMD="docker compose"
else
  echo "[ERRO] Docker Compose não encontrado (docker-compose ou docker compose)."
  exit 1
fi

if [ ! -f ".env" ]; then
  echo "[INFO] .env não encontrado. Copiando de .env-example..."
  cp .env-example .env
fi

echo "[INFO] Subindo dependências para debug local..."
$COMPOSE_CMD up -d db dev

echo "[INFO] Aguardando banco ficar pronto..."
$COMPOSE_CMD exec -T db sh -c 'until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do sleep 1; done'

echo "[INFO] Executando migrations..."
if command -v go >/dev/null 2>&1; then
  go run ./cmd/migrate
else
  echo "[INFO] Go não encontrado no host. Rodando migration no container dev..."
  POSTGRES_USER_VALUE="$(grep '^POSTGRES_USER=' .env | cut -d'=' -f2- | tr -d '\r')"
  POSTGRES_PASSWORD_VALUE="$(grep '^POSTGRES_PASSWORD=' .env | cut -d'=' -f2- | tr -d '\r')"
  POSTGRES_DB_VALUE="$(grep '^POSTGRES_DB=' .env | cut -d'=' -f2- | tr -d '\r')"
  POSTGRES_HOST_VALUE="$(grep '^POSTGRES_HOST=' .env | cut -d'=' -f2- | tr -d '\r')"
  POSTGRES_PORT_VALUE="$(grep '^POSTGRES_PORT=' .env | cut -d'=' -f2- | tr -d '\r')"

  $COMPOSE_CMD exec -T \
    -e POSTGRES_USER="$POSTGRES_USER_VALUE" \
    -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD_VALUE" \
    -e POSTGRES_DB="$POSTGRES_DB_VALUE" \
    -e POSTGRES_HOST="$POSTGRES_HOST_VALUE" \
    -e POSTGRES_PORT="$POSTGRES_PORT_VALUE" \
    dev go run ./cmd/migrate
fi

echo ""
echo "[OK] Dependências prontas para debug local."
echo "[INFO] Agora inicie a API em modo debug no VS Code (cmd/api/main.go)."
echo "[INFO] API esperada em: http://localhost:8080"
