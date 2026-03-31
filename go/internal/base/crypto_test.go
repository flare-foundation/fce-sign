package base

import "testing"

func TestKeccak256(t *testing.T) {
	hash := Keccak256([]byte{})
	expected := "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
	got := BytesToHex(hash)
	if got != "0x"+expected {
		t.Errorf("Keccak256('') = %s, want 0x%s", got, expected)
	}
}
