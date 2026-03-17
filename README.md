# TEE Extension Example - Private Key Manager

An example TEE extension that stores a private key and signs messages with it.

## For Hackathon Participants

Pick the branch for your preferred language and use it as a starting point.
Modify the files in `app/` and `contract/InstructionSender.sol` to build your
own extension. The files in `base/` are framework infrastructure — you should
not need to modify them. See the branch README for language-specific details.

## Branches

| Branch | Language |
|--------|----------|
| `go` | Go |
| `python` | Python |
| `typescript` | TypeScript |

Each branch contains:
- `extension/` — the extension server (`base/` = infrastructure, `app/` = your code)
- `contract/InstructionSender.sol` — the on-chain instruction sender contract
- `Dockerfile` — multi-stage build for tee-node + extension
- Unit tests
