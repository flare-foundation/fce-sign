package app

import (
	"fmt"
	"math/big"
)

// abiEncodeTwo ABI-encodes two dynamic byte arrays: (bytes, bytes).
// Layout:
//
//	offset of first bytes  (32 bytes) = 64
//	offset of second bytes (32 bytes) = 64 + 32 + padded(len(a))
//	length of a            (32 bytes)
//	a data                 (ceil(len(a)/32)*32 bytes)
//	length of b            (32 bytes)
//	b data                 (ceil(len(b)/32)*32 bytes)
func abiEncodeTwo(a, b []byte) ([]byte, error) {
	aPadded := padToMultipleOf32(a)
	bPadded := padToMultipleOf32(b)

	offsetA := big.NewInt(64)
	offsetB := big.NewInt(int64(64 + 32 + len(aPadded)))

	buf := make([]byte, 0, 64+32+len(aPadded)+32+len(bPadded))

	buf = append(buf, padLeft(offsetA.Bytes(), 32)...)
	buf = append(buf, padLeft(offsetB.Bytes(), 32)...)

	buf = append(buf, padLeft(big.NewInt(int64(len(a))).Bytes(), 32)...)
	buf = append(buf, aPadded...)

	buf = append(buf, padLeft(big.NewInt(int64(len(b))).Bytes(), 32)...)
	buf = append(buf, bPadded...)

	return buf, nil
}

// abiDecodeTwo decodes ABI-encoded (bytes, bytes) back into two byte slices.
func abiDecodeTwo(data []byte) ([]byte, []byte, error) {
	if len(data) < 128 {
		return nil, nil, fmt.Errorf("data too short for ABI-encoded (bytes, bytes)")
	}

	offsetA := new(big.Int).SetBytes(data[0:32]).Int64()
	offsetB := new(big.Int).SetBytes(data[32:64]).Int64()

	a, err := abiReadBytes(data, offsetA)
	if err != nil {
		return nil, nil, err
	}
	b, err := abiReadBytes(data, offsetB)
	if err != nil {
		return nil, nil, err
	}
	return a, b, nil
}

func abiReadBytes(data []byte, offset int64) ([]byte, error) {
	if offset+32 > int64(len(data)) {
		return nil, fmt.Errorf("offset %d out of range", offset)
	}
	length := new(big.Int).SetBytes(data[offset : offset+32]).Int64()
	start := offset + 32
	if start+length > int64(len(data)) {
		return nil, fmt.Errorf("bytes data out of range")
	}
	return data[start : start+length], nil
}

// padToMultipleOf32 pads data to a multiple of 32 bytes with trailing zeros.
func padToMultipleOf32(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	remainder := len(data) % 32
	if remainder == 0 {
		result := make([]byte, len(data))
		copy(result, data)
		return result
	}
	padded := make([]byte, len(data)+32-remainder)
	copy(padded, data)
	return padded
}
