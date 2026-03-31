package base

import (
	"testing"
)

func TestHexToBytes(t *testing.T) {
	b, err := HexToBytes("0xdeadbeef")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 4 || b[0] != 0xde || b[1] != 0xad || b[2] != 0xbe || b[3] != 0xef {
		t.Errorf("unexpected bytes: %x", b)
	}

	b, err = HexToBytes("abcd")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 2 || b[0] != 0xab || b[1] != 0xcd {
		t.Errorf("unexpected bytes: %x", b)
	}

	b, err = HexToBytes("")
	if err != nil {
		t.Fatal(err)
	}
	if b != nil {
		t.Errorf("expected nil for empty hex, got %x", b)
	}

	b, err = HexToBytes("0x")
	if err != nil {
		t.Fatal(err)
	}
	if b != nil {
		t.Errorf("expected nil for 0x, got %x", b)
	}
}

func TestBytesToHex(t *testing.T) {
	got := BytesToHex([]byte{0xde, 0xad})
	if got != "0xdead" {
		t.Errorf("BytesToHex = %s, want 0xdead", got)
	}
}

func TestPadLeft(t *testing.T) {
	b := PadLeft([]byte{0x01}, 4)
	if len(b) != 4 || b[0] != 0 || b[1] != 0 || b[2] != 0 || b[3] != 1 {
		t.Errorf("unexpected PadLeft result: %x", b)
	}

	b = PadLeft([]byte{0x01, 0x02, 0x03, 0x04, 0x05}, 4)
	if len(b) != 4 || b[0] != 0x02 || b[1] != 0x03 || b[2] != 0x04 || b[3] != 0x05 {
		t.Errorf("unexpected PadLeft result for oversized input: %x", b)
	}
}
