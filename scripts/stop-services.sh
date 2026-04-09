#!/usr/bin/env bash
# stop-services.sh — Stop the sign extension Docker Compose stack.
#
# Usage: ./scripts/stop-services.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
log() { echo -e "${GREEN}[stop-services]${NC} $*"; }

log "Stopping Docker Compose services..."
docker compose -f "$PROJECT_DIR/docker-compose.yaml" down
log "Done."
