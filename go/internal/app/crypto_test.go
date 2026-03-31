package app

import (
	"math/big"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"sign-extension/internal/base"
)

func TestSignECDSA(t *testing.T) {
	keyBytes := padLeft(big.NewInt(12345).Bytes(), 32)
	privKey, err := parseSecp256k1PrivateKey(keyBytes)
	if err != nil {
		t.Fatalf("parseSecp256k1PrivateKey: %v", err)
	}

	message := []byte("test message")
	sig, err := signECDSA(privKey, message)
	if err != nil {
		t.Fatalf("signECDSA: %v", err)
	}

	if len(sig) != 65 {
		t.Fatalf("expected 65-byte signature, got %d", len(sig))
	}

	// Verify using the secp256k1 library: reconstruct compact sig (v||r||s) and verify.
	hash := base.Keccak256(message)
	compact := make([]byte, 65)
	compact[0] = sig[64] // already 27 or 28, as expected by RecoverCompact
	copy(compact[1:33], sig[0:32])
	copy(compact[33:65], sig[32:64])
	pubKey, _, err := ecdsa.RecoverCompact(compact, hash)
	if err != nil {
		t.Fatalf("RecoverCompact: %v", err)
	}
	if !pubKey.IsEqual(privKey.PubKey()) {
		t.Error("recovered public key does not match")
	}
}

func TestSignECDSADifferentMessages(t *testing.T) {
	keyBytes := padLeft(big.NewInt(999999).Bytes(), 32)
	privKey, err := parseSecp256k1PrivateKey(keyBytes)
	if err != nil {
		t.Fatal(err)
	}

	msg1 := []byte("message one")
	msg2 := []byte("message two")

	sig1, err := signECDSA(privKey, msg1)
	if err != nil {
		t.Fatal(err)
	}
	sig2, err := signECDSA(privKey, msg2)
	if err != nil {
		t.Fatal(err)
	}

	if string(sig1) == string(sig2) {
		t.Error("signatures for different messages should differ")
	}

	// Verify both signatures recover to the same public key.
	hash1 := base.Keccak256(msg1)
	compact1 := make([]byte, 65)
	compact1[0] = sig1[64]
	copy(compact1[1:33], sig1[0:32])
	copy(compact1[33:65], sig1[32:64])
	pub1, _, err := ecdsa.RecoverCompact(compact1, hash1)
	if err != nil || !pub1.IsEqual(privKey.PubKey()) {
		t.Error("sig1 verification failed")
	}

	hash2 := base.Keccak256(msg2)
	compact2 := make([]byte, 65)
	compact2[0] = sig2[64]
	copy(compact2[1:33], sig2[0:32])
	copy(compact2[33:65], sig2[32:64])
	pub2, _, err := ecdsa.RecoverCompact(compact2, hash2)
	if err != nil || !pub2.IsEqual(privKey.PubKey()) {
		t.Error("sig2 verification failed")
	}
}

func TestSignECDSADeterministic(t *testing.T) {
	keyBytes := padLeft(big.NewInt(42).Bytes(), 32)
	privKey, err := parseSecp256k1PrivateKey(keyBytes)
	if err != nil {
		t.Fatal(err)
	}
	message := []byte("deterministic test")

	sig1, err := signECDSA(privKey, message)
	if err != nil {
		t.Fatal(err)
	}
	sig2, err := signECDSA(privKey, message)
	if err != nil {
		t.Fatal(err)
	}
	if string(sig1) != string(sig2) {
		t.Error("signatures for the same message should be deterministic")
	}
}

func TestParseSecp256k1PrivateKey(t *testing.T) {
	keyBytes := padLeft(big.NewInt(1).Bytes(), 32)
	key, err := parseSecp256k1PrivateKey(keyBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the key serializes back to the expected 32-byte scalar.
	serialized := key.Serialize()
	if len(serialized) != 32 {
		t.Errorf("expected 32-byte serialized key, got %d", len(serialized))
	}
	scalar := new(big.Int).SetBytes(serialized)
	if scalar.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("unexpected key scalar: want 1, got %v", scalar)
	}

	_, err = parseSecp256k1PrivateKey(nil)
	if err == nil {
		t.Error("expected error for nil key")
	}

	_, err = parseSecp256k1PrivateKey([]byte{})
	if err == nil {
		t.Error("expected error for empty key")
	}

	_, err = parseSecp256k1PrivateKey(make([]byte, 32))
	if err == nil {
		t.Error("expected error for zero key")
	}

	// secp256k1 curve order N
	order, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	_, err = parseSecp256k1PrivateKey(order.Bytes())
	if err == nil {
		t.Error("expected error for key >= curve order")
	}
}

func TestKeccak256(t *testing.T) {
	hash := base.Keccak256([]byte{})
	expected := "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
	got := base.BytesToHex(hash)
	if got != "0x"+expected {
		t.Errorf("Keccak256('') = %s, want 0x%s", got, expected)
	}
}
