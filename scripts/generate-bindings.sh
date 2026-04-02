#!/usr/bin/env bash
# generate-bindings.sh — Compile Solidity contracts and generate Go bindings.
#
# Prerequisites: forge (Foundry), jq
#
# Usage: ./scripts/generate-bindings.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# --- Contract name and Go output ---
CONTRACT_NAME="InstructionSender"
BINDINGS_DIR="$PROJECT_DIR/go/tools/app/contract"

echo "=== Step 1: Compile Solidity contracts ==="
cd "$PROJECT_DIR/contract"
forge build --root . --contracts . --out out

FORGE_OUT="$PROJECT_DIR/contract/out/InstructionSender.sol/${CONTRACT_NAME}.json"
if [[ ! -f "$FORGE_OUT" ]]; then
    echo "ERROR: forge output not found at $FORGE_OUT"
    exit 1
fi

echo "=== Step 2: Extract ABI and BIN ==="
mkdir -p "$BINDINGS_DIR"

jq '.abi' "$FORGE_OUT" > "$BINDINGS_DIR/${CONTRACT_NAME}.abi"
jq -r '.bytecode.object' "$FORGE_OUT" | sed 's/^0x//' > "$BINDINGS_DIR/${CONTRACT_NAME}.bin"

echo "  ABI → $BINDINGS_DIR/${CONTRACT_NAME}.abi"
echo "  BIN → $BINDINGS_DIR/${CONTRACT_NAME}.bin"

echo "=== Step 3: Generate Go bindings ==="
cd "$PROJECT_DIR/go/tools"
go generate ./app/...

echo "=== Done ==="
echo "Generated: $BINDINGS_DIR/autogen.go"
