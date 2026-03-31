package base

import (
	"encoding/hex"
	"strings"
)

// HexToBytes decodes a hex string (optional 0x prefix) to bytes.
func HexToBytes(h string) ([]byte, error) {
	h = strings.TrimPrefix(h, "0x")
	if len(h) == 0 {
		return nil, nil
	}
	return hex.DecodeString(h)
}

// BytesToHex encodes bytes to a 0x-prefixed hex string.
func BytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// PadLeft pads a byte slice to the specified length with leading zeros.
func PadLeft(b []byte, size int) []byte {
	if len(b) >= size {
		return b[len(b)-size:]
	}
	result := make([]byte, size)
	copy(result[size-len(b):], b)
	return result
}
