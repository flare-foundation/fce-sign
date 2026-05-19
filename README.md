# TEE Extension Example — Private Key Manager (sign)

An example TEE extension that stores a private key and signs messages with it.

> **Warning**: This repo is for demonstration purposes only. Storing encrypted
> secrets on-chain is not advisable in production — on-chain data is public
> and encryption can be broken over time. A production extension should use
> off-chain channels for secret delivery.

## Layout & deployable surface

This repo contains three implementations of the same extension:

| Language   | Directory                    | Deployable on Coston/Coston2? |
| ---------- | ---------------------------- | ----------------------------- |
| Go         | [`go/`](go/)                 | **Yes** — production target   |
| Python     | [`python/`](python/)         | No — reference implementation |
| TypeScript | [`typescript/`](typescript/) | No — reference implementation |

The **Go** implementation is the one wired into `Dockerfile`, `docker-compose.yaml`,
and `scripts/*` for the testnet deploy flow described in
[`testnet-deployment.md`](testnet-deployment.md).

Python and TypeScript implementations are kept as side-by-side references for
hackathon participants who want to study the protocol in a different language.
They have their own `Dockerfile` and tests but do **not** deploy via the
scaffold-aligned flow — they're code-reading material, not deploy artifacts.

### Running the hackathon-style language tests

| Language   | Test command                                                        |
| ---------- | ------------------------------------------------------------------- |
| Go         | `cd go && go test ./...`                                            |
| Python     | `cd python && python3 -m unittest discover -s tests -p 'test_*.py'` |
| TypeScript | `cd typescript && npm ci && npm test`                               |

## Shared contract

`contracts/InstructionSender.sol` is shared across all implementations. It
declares `OP_TYPE_KEY = bytes32("KEY")`, `OP_COMMAND_UPDATE = bytes32("UPDATE")`
and `OP_COMMAND_SIGN = bytes32("SIGN")` and exposes `updateKey(bytes)` and
`sign(bytes)` entry points that route through the Flare TEE Manager diamond.

## Deploying and Testing

The full testnet flow (Coston/Coston2 with a devops-hosted Confidential Space
VM) is documented in [`testnet-deployment.md`](testnet-deployment.md). The
short version:

```bash
bash ./scripts/use-chain.sh coston2       # or coston
bash ./scripts/full-setup.sh --chain coston2 --test
```

For local development against a Hardhat devnet + locally-built Docker stack:

```bash
bash ./scripts/use-chain.sh local
bash ./scripts/full-setup.sh --test       # defaults to --chain local
```

Each phase can also be run individually:

```bash
./scripts/pre-build.sh         # 1. Deploy contract + register extension → config/extension.env
./scripts/start-services.sh    # 2. Docker compose up (redis + ext-proxy + extension-tee)
./scripts/post-build.sh        # 3. Allow TEE version + register TEE machine on-chain
./scripts/test.sh              # 4. End-to-end UPDATE/SIGN test against the running TEE
./scripts/stop-services.sh     # Tear down
```

### Prerequisites

- **Docker Desktop** (Linux containers) — for the local stack
- **Go 1.25.1+** — for the deploy + registration tools in `go/tools/`
- **Foundry** (`forge`, `cast`, `jq`) — for Solidity compilation and contract bindings
- **Bash** — Git Bash works on Windows
- **Sibling repos** at `../../tee-node/` and `../../tee-proxy/` (cloned next to
  `extensions/` under a shared `tee/` parent — see
  [`REPRODUCIBILITY.md`](REPRODUCIBILITY.md) for layout)
- **A funded private key** for the target chain. Set as `DEPLOYMENT_PRIVATE_KEY`
  in `.env.<chain>` (no `0x` prefix). Fund at
  [`faucet.flare.network`](https://faucet.flare.network/).
- For Coston/Coston2 deploys: a devops contact who'll run the TEE on a real
  GCP Confidential Space VM. See `testnet-deployment.md` for the full handoff.

### Chain selection

`.env` is a per-chain file. `scripts/use-chain.sh <chain>` copies the active
chain's template (`.env.coston`, `.env.coston2`, or `.env.example` for local)
over `.env`. All scripts then source `.env` automatically.

| Chain     | `.env.<chain>`  | Addresses file                          | RPC                                              |
| --------- | --------------- | --------------------------------------- | ------------------------------------------------ |
| local     | `.env.example`  | `e2e/docker/sim_dump/deployed-addresses.json` (auto-detected) | `http://127.0.0.1:8545`                          |
| coston    | `.env.coston`   | `config/coston/deployed-addresses.json` | `https://coston-api.flare.network/ext/C/rpc`     |
| coston2   | `.env.coston2`  | `config/coston2/deployed-addresses.json`| `https://coston2-api.flare.network/ext/C/rpc`    |

### Generated artifacts

`pre-build.sh` writes the new `EXTENSION_ID` and `INSTRUCTION_SENDER` to
`config/extension.env`. Every subsequent script (`start-services`, `post-build`,
`test`) reads this file automatically — no manual `.env` edits required.

## Reproducible builds

The Dockerfile is reproducible: same source + same `SOURCE_DATE_EPOCH` yields a
bit-for-bit identical image. See [`REPRODUCIBILITY.md`](REPRODUCIBILITY.md).

## Troubleshooting

See `testnet-deployment.md` § Troubleshooting for the full catalogue. Common
issues:

- **`connect: connection refused` from ext-proxy** — VPN to Flare's indexer DB
  (`35.241.249.150:3306`) is down on Coston/Coston2 deploys.
- **`Verification.TeeNotFound`** — `NORMAL_PROXY_URL` is pointed at the wrong
  chain's FTDC proxy.
- **`Verification.ChallengeExpired`** — re-run `post-build.sh`; the patched
  `register-tee` already passes `-command rRap` for fresh attestation.
- **`code hashes do not match`** — `SIMULATED_TEE` and the TEE's `MODE` env
  disagree. Both must point at "real" for testnet (`SIMULATED_TEE=false`,
  `MODE=0`).

## Related docs

| Doc                                            | What it covers                                                   |
| ---------------------------------------------- | ---------------------------------------------------------------- |
| [`testnet-deployment.md`](testnet-deployment.md) | End-to-end Coston/Coston2 deploy with devops handoff             |
| [`REPRODUCIBILITY.md`](REPRODUCIBILITY.md)       | `SOURCE_DATE_EPOCH` and reproducible image builds                |
| [`go/`](go/)                                    | Go extension binary + tooling (production-deployable)            |
| [`python/`](python/) / [`typescript/`](typescript/) | Reference implementations of the same protocol (no deploy) |
