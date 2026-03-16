# TEE Extension - Private Key Manager

A TEE extension that accepts private keys and signs messages with them.

## Overview

This extension handles `KEY` operations with two commands:
- **UPDATE**: Decrypts an encrypted message using the TEE node's private key, then stores it as an ECDSA private key in memory
- **SIGN**: Signs a message with the stored private key and returns both the message and signature in the action result

## Implementations

This repository contains three implementations on separate branches:

| Branch | Language | Crypto |
|--------|----------|--------|
| `go` | Go | `golang.org/x/crypto/sha3` (Keccak-256), stdlib `crypto/ecdsa` |
| `python` | Python | Pure Python secp256k1 ECDSA + Keccak-256 (zero dependencies) |
| `typescript` | TypeScript | `@noble/hashes` (Keccak-256), `@noble/secp256k1` (ECDSA) |

Each branch contains:
- `extension/` — the extension server (with `base/` infrastructure and `app/` custom code)
- `contracts/InstructionSender.sol` — the on-chain instruction sender contract
- `Dockerfile` — multi-stage build for tee-node + extension
- Unit tests
