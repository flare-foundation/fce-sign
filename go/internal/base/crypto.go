package base

import "golang.org/x/crypto/sha3"

// Keccak256 computes the Keccak-256 hash of data.
func Keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}
