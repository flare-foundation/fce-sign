# TEE Extension Example - Private Key Manager

An example TEE extension that stores a private key and signs messages with it.
Use this as a **hackathon starter template**: clone it, modify the code to create
your own extension, then deploy/register/test it on Coston2.

> **Warning**: This repo is for demonstration purposes only. Storing encrypted
> secrets on-chain is not advisable in production — on-chain data is public
> and encryption can be broken over time. A production extension should use
> off-chain channels for secret delivery.

## For Hackathon Participants

Pick the language you're most comfortable with and work inside its directory.
You should modify the files in `app/` and the shared
`contract/InstructionSender.sol`. The files in `base/` are framework
infrastructure -- you should not need to modify them.

| Language   | Directory                    | Test command                                                        |
| ---------- | ---------------------------- | ------------------------------------------------------------------- |
| Go         | [`go/`](go/)                 | `cd go && go test ./...`                                            |
| Python     | [`python/`](python/)         | `cd python && python3 -m unittest discover -s tests -p 'test_*.py'` |
| TypeScript | [`typescript/`](typescript/) | `cd typescript && npm ci && npm test`                               |

See each directory's `README.md` for details on the handler signature, what's
provided by `base/`, and what files to change.

### Agent skills (optional)

Install the skill for agentic coding:

```bash
npx skills add .
```

## Shared contract

`contract/InstructionSender.sol` is shared across all implementations. Update it
to match your extension's OPType/OPCommand constants.

## Deploying and Testing on Coston2

Run the sign extension locally, expose it to the internet via a tunnel
(cloudflared, ngrok, etc.), and register + test it on the Coston2 testnet
(chain ID 114).

All deployment, registration, and testing tools are in `go/tools/` and work
for **all extension languages**. The scripts interact with smart contracts and
the TEE proxy — they don't depend on the extension's implementation language.
Set `LANGUAGE` in `.env` to choose which Docker image to build.

See [`go/README.md`](go/README.md#tools-gotools) for tool details.

### Prerequisites

- Docker
- [cloudflared](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/)
  to expose a local port to the internet (no account required; [ngrok](https://ngrok.com/) also works but needs sign-up)
- [Foundry](https://book.getfoundry.sh/getting-started/installation) (`forge`, `cast`)
  for contract compilation and verification
- A funded Coston2 wallet (needs C2FLR for gas + TEE registration fees)
- Go >= 1.23 (for the deployment/registration tools in `go/tools/`)

### Step 0: Configure environment

```bash
cp .env.example .env
# Edit .env and fill in PRIVATE_KEY and INITIAL_OWNER.
# Set LANGUAGE=go, LANGUAGE=python, or LANGUAGE=typescript
# to choose which implementation the Docker stack runs.

cp config/proxy/extension_proxy.toml.example config/proxy/extension_proxy.toml
# Edit config/proxy/extension_proxy.toml and fill in the DB credentials
# for the Coston2 C-chain indexer (username and password in the [db] section).
```

### Step 1: Deploy InstructionSender

```bash
cd go/tools
go run ./cmd/deploy-contract
```

The deploy tool automatically verifies the contract source on the
[Coston2 block explorer](https://coston2-explorer.flare.network/).
Pass `--no-verify` to skip.

Save the printed address in `.env`:

```bash
# Add to .env
INSTRUCTION_SENDER="<deployed-address>"
```

### Step 2: Register the extension

```bash
cd go/tools
go run ./cmd/register-extension
```

Save the printed extension ID in `.env`:

```bash
# Add to .env
EXTENSION_ID="0x<64-hex-chars>"
```

### Step 3: Start the extension stack

Build and start the stack (rebuild when switching `LANGUAGE` in `.env`):

```bash
docker compose build
docker compose up -d
```

Wait for the proxy to become healthy:

```bash
until curl -sf http://localhost:6676/info >/dev/null 2>&1; do sleep 2; done
echo "Extension proxy is ready"
```

### Step 4: Start tunnel

In a separate terminal, expose the extension proxy port (6676) to the internet:

```bash
# Using cloudflared (no account required):
cloudflared tunnel --url http://localhost:6676

# Or using ngrok:
ngrok http 6676
```

Note the public HTTPS URL and add it to `.env`:

```bash
# Add to .env
TUNNEL_URL="https://<your-tunnel-url>"
```

> **Note**: The tunnel must stay running for the entire session. If your
> computer sleeps or restarts, restart the tunnel and update `TUNNEL_URL`
> in `.env` with the new URL.

### Step 5: Add TEE version

```bash
cd go/tools
go run ./cmd/allow-tee-version -p http://localhost:6676
```

### Step 6: Register the TEE machine

Make sure `TUNNEL_URL` is set correctly in `.env`.

```bash
cd go/tools
go run ./cmd/register-tee -p http://localhost:6676 -l
```

The `-l` flag enables local/test mode (required when the TEE returns a test
attestation token instead of a real GCP JWT).

The `-p` flag specifies an existing production TEE on extension 0 that
performs the FTDC availability check on your TEE. It defaults to
`https://tee-proxy-coston2-1.flare.rocks` (the Coston2 public TEE proxy).

### Step 7: Run the end-to-end test

Make sure `INSTRUCTION_SENDER` and `TUNNEL_URL` are set correctly in `.env`.

```bash
cd go/tools
go run ./cmd/run-test -p http://localhost:6676
```

The test will:

1. Call `setExtensionId()` on the InstructionSender
2. Fetch the TEE's public key from the proxy
3. ECIES-encrypt a test private key and send `updateKey` on-chain
4. Wait for the TEE to process the instruction
5. Send a `sign` instruction on-chain
6. Verify the returned signature matches the test private key

---

## Port reference

| Service            | Container port | Host port |
| ------------------ | -------------- | --------- |
| ext-proxy internal | 6663           | 6675      |
| ext-proxy external | 6664           | 6676      |
| redis              | 6379           | 6383      |

The tunnel exposes host port 6676 (ext-proxy external) to the internet.

## Troubleshooting

### Proxy won't start / DB sync error

The proxy needs a synced C-chain indexer DB. Check the proxy logs and verify
the DB credentials in `config/proxy/extension_proxy.toml`:

```bash
docker compose logs ext-proxy
```

### Transaction reverts

Ensure your wallet has enough C2FLR for gas + fees. The TEE fee calculator
determines the required fee for each operation.

### to-production times out

Try restarting the proxy — it may have missed a signing policy round:

```bash
docker compose down
docker compose up -d
```

If that doesn't help, the FDC attestation flow requires active relay providers
on Coston2. If no relay infrastructure is running, the availability check won't
complete.

### Tunnel URL changed

If your tunnel restarts and the URL changes, update `TUNNEL_URL` in `.env`
and restart the Docker stack (`docker compose down && docker compose up -d`),
then re-run steps 5-6 (allow-tee-version + register-tee) to register a new
TEE machine with the new URL.

## Cleanup

To shut down all local services and prepare for a fresh start:

### Stop the Docker stack

```bash
docker compose down
```

This stops and removes all containers (redis, ext-proxy, extension-tee).

### Full reset (start from scratch)

If you want to completely reset and follow the README from the beginning:

```bash
# Remove built images (forces rebuild)
docker compose down --rmi local

# Clear environment state
rm -f .env config/proxy/extension_proxy.toml
```

After a full reset, start again from [Step 0](#step-0-configure-environment).

> **Note**: On-chain state (deployed contracts, registered extensions, registered
> TEEs) cannot be reset. Each fresh start will deploy a new InstructionSender
> contract and register a new extension. This is fine for testing — Coston2 is
> a testnet.
