# 🚀 TEE Extension Deployment — Step by Step

Linear recipe to deploy a TEE extension to Flare Coston or Coston2. Run the steps top to bottom.

## Prerequisites

- 🐳 Docker Desktop (Linux containers)
- 🐹 Go 1.25.1+
- 🔨 Foundry (`forge`, `cast`)
- `jq`
- Bash (Git Bash on Windows works)
- VPN access to Flare's indexer DB (`35.241.249.150:3306`) — **only required if you run your own `ext-proxy` locally**. If you're using a devops-hosted proxy (the normal Coston/Coston2 path), devops's proxy queries the indexer for you and you never touch it.

## 1. Clone the extension repo

No sibling repos are required for the Go path — not for building the image, and
not for the deploy/registration/test tooling. Everything resolves `tee-node`
and `tee-proxy` from the public `github.com/flare-foundation` repos:

- **Extension image** (`Dockerfile`) and **deploy tooling** (`go/tools/`, used by
  `pre-build.sh` / `post-build.sh` / `test.sh` / `check-tee-state`) pull them as
  **Go modules** through `GOPROXY` (`proxy.golang.org`, falling back to `direct`),
  pinned by hash in `go/go.sum` (`tee-node v0.0.20`) and `go/tools/go.sum`
  (`tee-node v0.0.20` + `tee-proxy v0.0.17`). No `replace` directives, nothing
  copied from disk.
- **Proxy image** (`proxy/Dockerfile`) `git clone`s `tee-proxy` + `tee-node`
  straight from GitHub at build time (override the refs with the
  `TEE_PROXY_VERSION` / `TEE_NODE_VERSION` build args).

So a flat checkout of just your extension is enough:

```text
<workspace>/
└── extensions/
    └── <your-extension>/
```

> [!NOTE]
> All three languages are sibling-free. The **Python** and **TypeScript** image
> builds also `git clone` `tee-node` from GitHub at build time (no local
> checkout) and build from this same dir — `start-services.sh` just points
> `EXTENSION_DOCKERFILE` at the right one per `LANGUAGE`. The only caveat is
> reproducibility: the Python/TS images reach same-machine determinism but not
> cross-machine bit-for-bit parity — see [`REPRODUCIBILITY.md`](REPRODUCIBILITY.md).

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

LANGUAGE=go                                                          # go | python | typescript (which Dockerfile build-image.sh builds)
```

Activate it (optionally selecting the extension language):

```bash
bash ./scripts/use-chain.sh <chain> [go|python|typescript]   # language defaults to go
bash ./scripts/use-chain.sh --list                           # list available chains + languages
bash ./scripts/use-chain.sh --help                           # usage
```

Copies `.env.<chain>` → `.env` (and sets `LANGUAGE` if you passed one), which all
scripts auto-load. The language you pick here is what `build-image.sh` builds in Step 5.

## 4. Register the extension on-chain

```bash
bash ./scripts/pre-build.sh
```

Generates Go bindings, compiles Solidity, deploys `InstructionSender`, and registers the extension on-chain. Writes `EXTENSION_ID` and `INSTRUCTION_SENDER` to `config/extension.env`.

> [!IMPORTANT]
> Each run mints a **new** extension + InstructionSender and overwrites `config/extension.env`, so the script refuses to run if that file already exists.
>
> - **First deploy:** run the command above as-is.
> - **Want a new extension + InstructionSender** (e.g. after a diamond redeploy): opt in with `--force`:
>   ```bash
>   bash ./scripts/pre-build.sh --force   # or: PRE_BUILD_FORCE=1 bash ./scripts/pre-build.sh
>   ```
>
> Forcing a re-run against an existing TEE on a shared proxy orphans it, and `test.sh` later reverts with `MachineManager.TooMany()`. See [Troubleshooting](#troubleshooting) to recover.

Read the new values — `EXTENSION_ID` is part of the hand-off in Step 6:

```bash
cat config/extension.env
```

## 5. Build the Docker image

`build-image.sh` builds the image for the language in your `.env` (set by
`use-chain.sh` in Step 3), pins `SOURCE_DATE_EPOCH` for a reproducible
`codeHash`, verifies `MODE=0` is baked in (production attestation — FTDC rejects
`MODE=1`), and saves a tar for the hand-off:

```bash
bash ./scripts/build-image.sh                 # build .env's LANGUAGE, tag v0.1.0, save tar
bash ./scripts/build-image.sh -l typescript   # override the language
bash ./scripts/build-image.sh -v v0.1.1       # set the version/tag
```

It writes `sign-extension-<language>-<version>.tar` and prints the image ID to hand to devops.

<details>
<summary>Equivalent manual commands</summary>

```bash
export SOURCE_DATE_EPOCH=$(git log -1 --format=%ct)
docker build -f typescript/Dockerfile -t sign-extension-ts:v0.1.0 .   # or Dockerfile / python/Dockerfile
docker save sign-extension-ts:v0.1.0 -o sign-extension-ts-v0.1.0.tar
```

`docker build -t` names + tags in one step (no separate `docker tag`), and
BuildKit auto-reads `SOURCE_DATE_EPOCH` from the env. Confirm `MODE=0` with:

```bash
docker inspect sign-extension-ts:v0.1.0 --format '{{range .Config.Env}}{{println .}}{{end}}' | grep MODE
# expected: MODE=0
```
</details>

## 6. Deploy the image on a Confidential Space VM

Hand off (or deploy yourself) to a GCP Confidential Space VM with:

- The image (tar or registry URL+tag)
- Workload-launch env: `INITIAL_OWNER`, `CHAIN_URL`, `EXTENSION_ID` (from Step 4), `PROXY_URL` (proxy URL reachable from the TEE)
- Public HTTPS routed to port `6664` of the proxy container

You receive back the **public proxy URL** — which you only learn now, after the
hand-off. Put it in your `.env.<chain>` template, then re-run `use-chain.sh` so
the active `.env` picks it up. (`post-build.sh` / `test.sh` read `EXT_PROXY_URL`
from `.env`; the convention is to edit the chain template, never `.env`
directly.)

```bash
# in .env.<chain>
EXT_PROXY_URL=<public proxy URL>
```

```bash
bash ./scripts/use-chain.sh <chain>   # re-copies .env.<chain> → .env, now with EXT_PROXY_URL set
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

1. `bash ./scripts/pre-build.sh --force` — mints a fresh `EXTENSION_ID`. The `--force` opt-in is required because `config/extension.env` still has the now-invalid values from the previous deploy (see Step 4).
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
