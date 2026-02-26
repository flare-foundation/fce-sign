# TEE Extension - Private Key Manager

A TEE extension that accepts private keys and signs messages with them, written in Go.

## Overview

This extension handles `KEY` operations with two commands:
- **UPDATE**: Decrypts an encrypted message using the TEE node's private key, then stores it as an ECDSA private key in memory
- **SIGN**: Signs a message with the stored private key and returns both the message and signature in the action result
