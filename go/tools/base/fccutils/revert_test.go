package fccutils

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestDecodeRevertHex_StandardError(t *testing.T) {
	// Error(string) with message "test error"
	// Selector: 08c379a0
	// Then ABI-encoded string "test error"
	reason := decodeRevertHex("08c379a0" +
		"0000000000000000000000000000000000000000000000000000000000000020" +
		"000000000000000000000000000000000000000000000000000000000000000a" +
		"74657374206572726f7200000000000000000000000000000000000000000000")
	if reason != "test error" {
		t.Errorf("expected 'test error', got %q", reason)
	}
}

func TestDecodeRevertHex_CustomError(t *testing.T) {
	// AlreadyRegistered() should return the raw hex (prefixed) since it's a
	// known custom error.
	selector := hex.EncodeToString(crypto.Keccak256([]byte("AlreadyRegistered()"))[:4])
	result := decodeRevertHex(selector)
	if result != "0x"+selector {
		t.Errorf("expected 0x%s, got %q", selector, result)
	}
}

func TestDecodeRevertHex_UnknownSelector(t *testing.T) {
	// Unknown 4-byte selector should still return 0x-prefixed hex.
	result := decodeRevertHex("deadbeef")
	if result != "0xdeadbeef" {
		t.Errorf("expected '0xdeadbeef', got %q", result)
	}
}

func TestDecodeRevertHex_TooShort(t *testing.T) {
	result := decodeRevertHex("abcd")
	if result != "" {
		t.Errorf("expected empty for short input, got %q", result)
	}
}

func TestDecodeRevertHex_WithPrefix(t *testing.T) {
	selector := hex.EncodeToString(crypto.Keccak256([]byte("FeeTooLow()"))[:4])
	result := decodeRevertHex("0x" + selector)
	if result != "0x"+selector {
		t.Errorf("expected 0x%s, got %q", selector, result)
	}
}
