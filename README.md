# TEE Extension - Private Key Manager

A TEE extension that accepts private keys and signs messages with them, written in Go.

## Overview

This extension handles `KEY` operations with two commands:
- **UPDATE**: Decrypts an encrypted message using the TEE node's private key, then stores it as an ECDSA private key in memory
- **SIGN**: Signs a message with the stored private key and returns both the message and signature in the action result

## Skills

This repository ships with skills for AI coding agents (Claude Code, opencode, and others that support the [skills.sh](https://skills.sh) format).

### Install

From within this repository:

```bash
npx skills add .
```

This installs the skills into your agent's configuration directory. To update:

```bash
npx skills update
```

### Available skills

| Skill | Description |
|-------|-------------|
| `create-extension` | Architecture, terminology, protocol spec, and step-by-step implementation guide. Also includes `references/` with JSON schemas, OpenAPI spec, and Solidity interfaces. |

### Using skills

Once installed, prompt your agent naturally:

```
Make a TEE extension that processes instructions with OPType KEY and two OPCommands:
- UPDATE: stores a private key (decrypted with the TEE node's key)
- SIGN: signs a message with the stored private key and returns both message and signature

Write unit tests. Use Go.
```

Or use explicit commands: `create my extension`, `write tests for my extension`.
