package app

import (
	"bytes"
	"testing"
)

func TestAbiEncodeTwoRoundTrip(t *testing.T) {
	a := []byte("hello world")
	b := []byte{0xde, 0xad, 0xbe, 0xef}

	encoded, err := abiEncodeTwo(a, b)
	if err != nil {
		t.Fatalf("abiEncodeTwo: %v", err)
	}

	decodedA, decodedB, err := abiDecodeTwo(encoded)
	if err != nil {
		t.Fatalf("abiDecodeTwo: %v", err)
	}

	if !bytes.Equal(a, decodedA) {
		t.Errorf("first arg mismatch: got %x, want %x", decodedA, a)
	}
	if !bytes.Equal(b, decodedB) {
		t.Errorf("second arg mismatch: got %x, want %x", decodedB, b)
	}
}

func TestAbiEncodeTwoEmpty(t *testing.T) {
	encoded, err := abiEncodeTwo(nil, nil)
	if err != nil {
		t.Fatalf("abiEncodeTwo: %v", err)
	}

	decodedA, decodedB, err := abiDecodeTwo(encoded)
	if err != nil {
		t.Fatalf("abiDecodeTwo: %v", err)
	}

	if len(decodedA) != 0 {
		t.Errorf("expected empty first arg, got %x", decodedA)
	}
	if len(decodedB) != 0 {
		t.Errorf("expected empty second arg, got %x", decodedB)
	}
}

func TestAbiEncodeTwoExact32(t *testing.T) {
	a := make([]byte, 32)
	for i := range a {
		a[i] = byte(i)
	}
	b := []byte{0xff}

	encoded, err := abiEncodeTwo(a, b)
	if err != nil {
		t.Fatalf("abiEncodeTwo: %v", err)
	}

	decodedA, decodedB, err := abiDecodeTwo(encoded)
	if err != nil {
		t.Fatalf("abiDecodeTwo: %v", err)
	}

	if !bytes.Equal(a, decodedA) {
		t.Errorf("first arg mismatch")
	}
	if !bytes.Equal(b, decodedB) {
		t.Errorf("second arg mismatch")
	}
}
