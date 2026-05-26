# 🚀 TEE Extension Deployment — Step by Step

Linear recipe to deploy a TEE extension to Flare Coston or Coston2. Run the steps top to bottom.

## Prerequisites

- 🐳 Docker Desktop (Linux containers)
- 🐹 Go 1.25.1+
- 🔨 Foundry (`forge`, `cast`)
- `jq`
- Bash (Git Bash on Windows works)
- VPN access to Flare's indexer DB (`35.241.249.150:3306`) — **only required if you run your own `ext-proxy` locally**. If you're using a devops-hosted proxy (the normal Coston/Coston2 path), devops's proxy queries the indexer for you and you never touch it.

## 1. Clone sibling repos

The extension's Dockerfiles consume both repos from `../../tee-node/`.

```text
<workspace>/tee/
├── tee-node/         # gitlab.com/flarenetwork/tee/tee-node, tag v0.0.20
├── tee-proxy/        # gitlab.com/flarenetwork/tee/tee-proxy, main @ a3adb51 or newer
└── extensions/
    └── <your-extension>/
```

> [!IMPORTANT]
> Use `tee-proxy` at commit `a3adb51` (or any commit after tag `v0.0.17`) — not the
> `v0.0.17` tag itself. `v0.0.17`'s Dockerfile has the build context wrong
> (`WORKDIR /app/tee-proxy` before `COPY . .`) and `COPY config.example.toml`
> in the final stage points at a path that doesn't exist outside the builder.
> Commit `a3adb51` ("chore: use local tee-node via replace directive in build")
> fixes both. A `v0.0.18` tag has not been cut yet — once it is, update this line.

## 2. Generate a funded deployer key

```bash
cast wallet new
cast wallet address --private-key 0x<private-key>
```

The derived address becomes your `INITIAL_OWNER`. Fund it from the target chain's faucet.

| Chain   | Faucet                                 |
| ------- | -------------------------------------- |
| Coston  | `https://faucet.flare.network/coston`  |
| Coston2 | `https://faucet.flare.network/coston2` |

## 3. Create `.env.<chain>`

Copy `.env.example` to `.env.coston` or `.env.coston2`. Fill in:

```bash
CHAIN=coston2                                                         # or coston
CHAIN_URL=https://coston2-api.flare.network/ext/C/rpc                 # chain RPC
ADDRESSES_FILE=./config/coston2/deployed-addresses.json
NORMAL_PROXY_URL=https://tee-proxy-coston2-1.flare.rocks              # FTDC proxy
EXT_PROXY_URL=                                                        # leave empty — set in Step 6

LOCAL_MODE=false
SIMULATED_TEE=false
DEPLOYMENT_PRIVATE_KEY=<private key, no 0x prefix>
INITIAL_OWNER=0x<derived address from Step 2>
```

Activate it:

```bash
bash ./scripts/use-chain.sh <chain>
```

Copies `.env.<chain>` → `.env`, which all scripts auto-load.

## 4. Register the extension on-chain

```bash
bash ./scripts/pre-build.sh
```

Generates Go bindings, compiles Solidity, deploys `InstructionSender`, and registers the extension on-chain. Writes `EXTENSION_ID` and `INSTRUCTION_SENDER` to `config/extension.env`.

> [!WARNING]
> **`pre-build.sh` is destructive — it mints a NEW extension and a NEW InstructionSender on every run, and overwrites `config/extension.env`.**
>
> Re-running it against an existing TEE (the normal case on a shared Coston/Coston2 proxy, where the proxy's signing key is persistent) orphans the TEE: the TEE record on-chain stays bound to the previous extension, while the new InstructionSender points at the empty new extension. `post-build.sh` will then skip `ToProduction` (the TEE is already `PRODUCTION` on the old extension) and `test.sh` will revert with `MachineManager.TooMany()` because `getRandomTeeIds(newExtId, 1)` finds zero active TEEs.
>
> The script now refuses to clobber an existing `config/extension.env` unless you explicitly opt in:
>
> ```bash
> PRE_BUILD_FORCE=1 ./scripts/pre-build.sh
> # or:   ./scripts/pre-build.sh --force
> ```
>
> If you hit `TooMany()` from a previous run, run the diagnostic (from `go/tools`):
>
> ```bash
> go run ./cmd/check-tee-state \
>     -a ../../config/<chain>/deployed-addresses.json \
>     -c "$CHAIN_URL" \
>     -p "$EXT_PROXY_URL" \
>     -instructionSender "$INSTRUCTION_SENDER"
> ```
>
> It prints the InstructionSender registered against the TEE's actual extension — paste that address into `config/extension.env` and skip back to `test.sh`.

Read the new values — `EXTENSION_ID` is part of the hand-off in Step 6:

```bash
cat config/extension.env
```

## 5. Build the Docker image

Confirm `MODE=0` is the default in your extension's `Dockerfile` (`MODE=0` is the production attestation backend; `MODE=1` produces simulated attestation that FTDC rejects):

```dockerfile
ENV MODE=0 CONFIG_PORT=5501 SIGN_PORT=7701 EXTENSION_PORT=7702
```

Then build:

```powershell
$env:SOURCE_DATE_EPOCH = (git log -1 --format=%ct)
docker compose -f docker-compose.yaml build --no-cache extension-tee
docker tag <your-extension>-extension-tee:latest <your-extension>:v0.1.0
docker save <your-extension>:v0.1.0 -o <your-extension>-v0.1.0.tar
```

Setting `SOURCE_DATE_EPOCH` makes the build reproducible (same source → same `codeHash`).

Verify `MODE=0` is baked into the image:

```powershell
docker inspect <your-extension>:v0.1.0 --format '{{range .Config.Env}}{{println .}}{{end}}' | Select-String MODE
# expected: MODE=0
```

## 6. Deploy the image on a Confidential Space VM

Hand off (or deploy yourself) to a GCP Confidential Space VM with:

- The image (tar or registry URL+tag)
- Workload-launch env: `INITIAL_OWNER`, `CHAIN_URL`, `EXTENSION_ID` (from Step 4), `PROXY_URL` (proxy URL reachable from the TEE)
- Public HTTPS routed to port `6664` of the proxy container

You receive back the **public proxy URL**. Add it to `.env.<chain>` and re-activate:

```bash
# in .env.<chain>
EXT_PROXY_URL=<public proxy URL>
```

```bash
bash ./scripts/use-chain.sh <chain>
```

## 7. Verify the proxy `/info`

```powershell
curl -s $env:EXT_PROXY_URL/info | jq '.machineData'
```

Required values:

| Field          | Expected                                                      |
| -------------- | ------------------------------------------------------------- |
| `platform`     | starts with `0x4743505f414d445f534556…` (GCP_AMD_SEV)         |
| `codeHash`     | real measured hash (**not** `0x194844cf…` — that's simulated) |
| `extensionId`  | matches your `config/extension.env` `EXTENSION_ID`            |
| `initialOwner` | matches your `INITIAL_OWNER`                                  |

If `extensionId` is wrong, ask the VM operator to restart the container with the correct `EXTENSION_ID` env override (no image rebuild needed — it's a launch-policy override).

## 8. Register the TEE machine

> [!NOTE]
> `post-build.sh` already invokes `register-tee` with `-command rRap` (not the default `rap`) — this is a load-bearing detail; don't revert it. Step `a` (availability check) needs a one-time **challenge** — a random number from the contract that the TEE signs to prove it's alive. By default only `r` issues it, but `r` skips itself once the TEE is registered on-chain. So re-runs (image changes, diamond cuts, retries) would revert with `Verification.ChallengeExpired`. Capital `R` issues the challenge directly — decoupled from `r` — so re-runs work.

Run:

```bash
bash ./scripts/post-build.sh
```

- `allow-tee-version` whitelists the codeHash for your extension.
- `register-tee -command rRap` pre-registers the TEE, requests fresh attestation, runs the FTDC availability check, promotes to production.

## 9. End-to-end test

```bash
bash ./scripts/test.sh
```

Sends test instructions through the deployed TEE and verifies the round-trip.

---

## When the extension image changes

1. Rebuild and hand off the new image.
2. The VM is re-deployed → `codeHash` changes.
3. `bash ./scripts/post-build.sh` whitelists the new codeHash.
4. `bash ./scripts/test.sh`.

## When the `FlareTeeManager` diamond is re-deployed

All extension registrations on that chain are wiped:

1. `PRE_BUILD_FORCE=1 bash ./scripts/pre-build.sh` — mints a fresh `EXTENSION_ID`. The `--force` opt-in is required because `config/extension.env` still has the now-invalid values from the previous deploy (see Step 4).
2. Send the new `EXTENSION_ID` to the VM operator. They restart the container with `EXTENSION_ID=<new value>` as a launch-policy env override — no image rebuild needed.
3. Re-curl `/info` and confirm `extensionId` matches.
4. `bash ./scripts/post-build.sh`.
5. `bash ./scripts/test.sh`.

## Troubleshooting

If `test.sh` reverts or `post-build.sh` fails, run the diagnostic before changing anything:

```bash
cd go/tools
go run ./cmd/check-tee-state \
    -a ../../config/<chain>/deployed-addresses.json \
    -c "$CHAIN_URL" \
    -p "$EXT_PROXY_URL" \
    -instructionSender "$INSTRUCTION_SENDER"
```

It reads (read-only — no transactions) the TEE proxy's `/info`, the on-chain TEE machine record, `InstructionSender._extensionId` (via storage slot 0), and `getActiveTeeMachines` for each extension involved. The output ends with a verdict that maps directly to a fix:

| Symptom | Verdict line | Fix |
| --- | --- | --- |
| `test.sh` reverts with `0xd65ac61e` (`MachineManager.TooMany()`) | `MISMATCH: InstructionSender ext=X ≠ TEE on-chain ext=Y` | Set `INSTRUCTION_SENDER` in `config/extension.env` to the address the diag prints under `[TEE record ext=Y]`. |
| `post-build.sh` reverts with `MachineManager.InvalidTeeStatus()` | `toProduction will revert: status=1 (PRODUCTION)` | TEE is already promoted. Skip post-build entirely, or rely on the idempotency guard in `registration.go` (which short-circuits when status=PRODUCTION). |
| `check-tee-state` says active set is empty *and* IDs all agree | (no MISMATCH line; "active set was emptied for a non-status reason") | TEE was banned or its version disabled. Investigate via on-chain events; `pause()` → re-promote is the recovery path, but only the TEE machine owner can do it. |
