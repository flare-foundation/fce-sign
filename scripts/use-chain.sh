#!/usr/bin/env bash
# use-chain.sh — Swap the active .env for a given chain, optionally setting the
# extension language.
#
# Usage:
#   ./scripts/use-chain.sh <chain> [go|python|typescript]
#   ./scripts/use-chain.sh --list
#   ./scripts/use-chain.sh --help
#
# Reads .env.<chain> from the project root and copies it to .env. If a language
# is given, the LANGUAGE= line in the new .env is set (or appended) to it.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

LANGUAGES="go python typescript"

usage() { echo "usage: $0 <chain> [go|python|typescript] | --list | --help" >&2; }

print_help() {
    cat <<EOF
use-chain.sh — activate a chain's .env (and optionally set the extension language)

Usage:
  $0 <chain> [language]   Copy .env.<chain> → .env; optionally set LANGUAGE
  $0 --list               List available chains and languages
  $0 --help, -h           Show this help

Arguments:
  <chain>      A chain with a .env.<chain> template (see --list).
  [language]   Optional: go | python | typescript. Rewrites LANGUAGE= in .env.

Examples:
  $0 coston2                 Activate Coston2, keep the template's LANGUAGE
  $0 coston2 typescript      Activate Coston2 and set LANGUAGE=typescript
EOF
}

list_options() {
    echo "Available chains (from .env.* templates in $PROJECT_DIR):"
    local found=0 f name
    for f in "$PROJECT_DIR"/.env.*; do
        [[ -e "$f" ]] || continue
        name="${f##*/.env.}"
        [[ "$name" == "example" ]] && continue
        echo "  - $name"
        found=1
    done
    [[ "$found" -eq 1 ]] || echo "  (none found — create one by copying .env.example)"
    echo ""
    echo "Available languages:"
    local l
    for l in $LANGUAGES; do
        if [[ "$l" == "go" ]]; then
            echo "  - $l (default)"
        else
            echo "  - $l"
        fi
    done
}

# --- arg parsing ---
[[ $# -ge 1 ]] || { usage; exit 1; }

case "$1" in
    -h|--help) print_help; exit 0 ;;
    --list)    list_options; exit 0 ;;
    -*)        echo "unknown flag: $1" >&2; usage; exit 1 ;;
esac

[[ $# -le 2 ]] || { usage; exit 1; }
CHAIN="$1"
LANGUAGE="${2:-}"

SRC="$PROJECT_DIR/.env.$CHAIN"
DST="$PROJECT_DIR/.env"

[[ -f "$SRC" ]] || { echo "no such file: $SRC" >&2; echo "run '$0 --list' to see available chains" >&2; exit 1; }

cp "$SRC" "$DST"
echo "[use-chain] copied .env.$CHAIN → .env"

if [[ -n "$LANGUAGE" ]]; then
    case " $LANGUAGES " in
        *" $LANGUAGE "*) ;;
        *) echo "unknown language: $LANGUAGE (valid: $LANGUAGES)" >&2; exit 1 ;;
    esac
    if grep -qE '^LANGUAGE=' "$DST"; then
        tmp="$(mktemp)"
        sed -E "s|^LANGUAGE=.*|LANGUAGE=$LANGUAGE|" "$DST" > "$tmp" && mv "$tmp" "$DST"
    else
        printf '\n# Extension language: go|python|typescript.\nLANGUAGE=%s\n' "$LANGUAGE" >> "$DST"
    fi
fi

echo "[use-chain] EXT_PROXY_URL=$(grep -E '^EXT_PROXY_URL=' "$DST" | head -1 | cut -d= -f2-)"
echo "[use-chain] CHAIN_URL=$(grep -E '^CHAIN_URL=' "$DST" | head -1 | cut -d= -f2-)"
echo "[use-chain] LANGUAGE=$(grep -E '^LANGUAGE=' "$DST" | head -1 | cut -d= -f2-)"
