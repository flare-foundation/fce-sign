#!/usr/bin/env bash
# use-chain.sh — Swap the active .env for a given chain.
#
# Usage:
#   ./scripts/use-chain.sh coston
#   ./scripts/use-chain.sh coston2
#   ./scripts/use-chain.sh local
#
# Reads .env.<chain> from the project root and copies it to .env.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

[[ $# -eq 1 ]] || { echo "usage: $0 <local|coston|coston2>" >&2; exit 1; }
CHAIN="$1"

SRC="$PROJECT_DIR/.env.$CHAIN"
DST="$PROJECT_DIR/.env"

[[ -f "$SRC" ]] || { echo "no such file: $SRC" >&2; exit 1; }

cp "$SRC" "$DST"
echo "[use-chain] copied .env.$CHAIN → .env"
echo "[use-chain] EXT_PROXY_URL=$(grep -E '^EXT_PROXY_URL=' "$DST" | head -1)"
echo "[use-chain] CHAIN_URL=$(grep -E '^CHAIN_URL=' "$DST" | head -1)"
