# TEE Extension Example - Private Key Manager (Go)

An example TEE extension that stores a private key and signs messages with it.

## For Hackathon Participants

This is a **working example** to use as a starting point. You should modify the
files in `internal/app/` and the shared `contract/InstructionSender.sol` to build
your own extension. The files in `internal/base/` are framework infrastructure --
you should not need to modify them.

### What to change

| File | Purpose |
|------|---------|
| `internal/app/handlers.go` | Your business logic -- register handlers, process messages |
| `internal/app/config.go` | Version constant |
| `internal/app/types.go` | Request/response types for external calls |
| `internal/app/abi.go` | ABI encoding for your specific data types |
| `internal/app/crypto.go` | Cryptographic operations (only if your extension needs them) |
| `contract/InstructionSender.sol` | On-chain contract that sends instructions to your extension |

### What's provided by `base/`

| Package | Functions |
|---------|-----------|
| `base` (encoding) | `HexToBytes(hex)`, `BytesToHex(bytes)` |
| `base` (crypto) | `Keccak256(data)` |
| `base` (types) | `Framework` (handler registration), protocol types |
| `base` (server) | HTTP server (you never call this directly) |

### Handler signature

```go
func myHandler(msg string) (data *string, status int, err error) {
    // msg is the hex-encoded originalMessage from the on-chain instruction
    // Return: data, status, error
    //   status: 0 = error, 1 = success, >=2 = pending
    return &dataHex, 1, nil
}
```

## Tools (`go/tools/`)

The `tools/` directory contains Go programs for deploying, registering, and
testing the extension on Coston2. It is a separate Go module (`sign-tools`)
from the extension runtime.

> **Note**: These tools work for **all extension languages** (Go, Python,
> TypeScript). The scripts interact with smart contracts and the TEE proxy —
> they don't depend on the extension's implementation language. Set `LANGUAGE`
> in `.env` to choose which Docker image to build.

### Structure

```
tools/
  base/              # Generic -- copy to other extension repos as-is
    configs.go       # Constants, ReadAddresses, dev key
    support.go       # Support struct, .env loading, chain client
    fccutils/        # TEE contract helpers (registration, versioning, etc.)
  app/               # Extension-specific
    contract/        # Generated Go bindings (autogen.go, ABI, BIN)
    deploy.go        # Deploy contract, setExtensionId, send instructions
    test.go          # End-to-end test (ECIES, updateKey, sign, verify)
    generate.go      # go:generate directive for abigen
  cmd/
    deploy-contract/    # Deploy InstructionSender
    register-extension/ # Register extension + allowlists
    allow-tee-version/  # Add TEE version (code hash + platform)
    register-tee/       # Register TEE machine (pre-reg -> attest -> produce)
    run-test/           # Run end-to-end test
```

### Prerequisites

- Go >= 1.23
- `abigen` (from go-ethereum, for regenerating bindings)
- `.env` file at the repo root (copy from `.env.example`)

### Running the tools

All commands are run from the `go/tools/` directory. They read configuration
from `.env` at the repo root and `config/coston2/deployed-addresses.json`.

```bash
# Deploy the InstructionSender contract
go run ./cmd/deploy-contract

# Register the extension (after deploying)
go run ./cmd/register-extension --instructionSender 0x<address>

# Add TEE version (after extension stack is running)
go run ./cmd/allow-tee-version -p http://localhost:6676

# Register TEE machine (pre-reg, attest, to-production)
go run ./cmd/register-tee -p http://localhost:6676 -l

# Run the end-to-end test
go run ./cmd/run-test --instructionSender 0x<address> -p http://localhost:6676
```

### Environment variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PRIV_KEY` or `PRIVATE_KEY` | Wallet private key | Dev key (local only) |
| `CHAIN_URL` | Coston2 RPC URL | `https://coston2-api.flare.network/ext/C/rpc` |
| `ADDRESSES_FILE` | Deployed addresses JSON | `../../config/coston2/deployed-addresses.json` |

### Command-line flags

| Flag | Description | Used by |
|------|-------------|---------|
| `--instructionSender` | InstructionSender contract address | `register-extension`, `run-test` |
| `-p, --proxy` | Extension proxy URL | `allow-tee-version`, `register-tee`, `run-test` |
| `-l, --local` | Use test attestation (no real GCP JWT) | `register-tee` |
| `-ep, --ext-proxy` | Existing production TEE proxy for attestation | `register-tee` |

### Regenerating contract bindings

If you modify `contract/InstructionSender.sol`:

1. Compile with Foundry:
   ```bash
   cd ../../contract && forge build --root . --contracts . --out out
   ```

2. Extract ABI and BIN:
   ```bash
   jq -r '.abi' contract/out/InstructionSender.sol/InstructionSender.json > go/tools/app/contract/InstructionSender.abi
   jq -r '.bytecode.object' contract/out/InstructionSender.sol/InstructionSender.json > go/tools/app/contract/InstructionSender.bin
   ```

3. Regenerate Go bindings:
   ```bash
   cd go/tools && go generate ./...
   ```

## Testing

```bash
go test ./...
```
