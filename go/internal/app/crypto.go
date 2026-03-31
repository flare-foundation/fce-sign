package app

import (
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"sign-extension/internal/base"
)

func padLeft(b []byte, size int) []byte {
	if len(b) >= size {
		return b[len(b)-size:]
	}
	result := make([]byte, size)
	copy(result[size-len(b):], b)
	return result
}

func parseSecp256k1PrivateKey(b []byte) (*secp256k1.PrivateKey, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("key bytes are empty")
	}
	if len(b) > 32 {
		return nil, fmt.Errorf("key too long: %d bytes", len(b))
	}
	key := secp256k1.PrivKeyFromBytes(padLeft(b, 32))
	if key.Key.IsZero() {
		return nil, fmt.Errorf("key is zero")
	}
	return key, nil
}

func signECDSA(key *secp256k1.PrivateKey, message []byte) ([]byte, error) {
	hash := base.Keccak256(message)
	// SignCompact returns [recoveryFlag, r(32), s(32)] where recoveryFlag is 27 (even y) or 28 (odd y).
	sig := ecdsa.SignCompact(key, hash, false)
	// Build Ethereum-style [r(32), s(32), v] where v is 27 or 28.
	result := make([]byte, 65)
	copy(result[0:32], sig[1:33])   // r
	copy(result[32:64], sig[33:65]) // s
	result[64] = sig[0]             // v (already 27 or 28)
	return result, nil
}
