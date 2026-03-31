package base

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestDecodeCustomError_KnownErrors(t *testing.T) {
	tests := []struct {
		soliditySig  string
		wantName     string
		wantHintSub  string // substring expected in the hint
	}{
		{"FeeTooLow()", "FeeTooLow", "FEE_WEI"},
		{"AlreadyRegistered()", "AlreadyRegistered", "already registered on-chain"},
		{"VersionAlreadyExists()", "VersionAlreadyExists", "skip this step"},
		{"PlatformAlreadyExists(bytes32)", "PlatformAlreadyExists", "skip this step"},
		{"OwnerAlreadyAllowed(address)", "OwnerAlreadyAllowed", "Safe to continue"},
		{"KeyTypeAlreadyExists(bytes32)", "KeyTypeAlreadyExists", "Safe to continue"},
		{"TeeMachineNotAvailable()", "TeeMachineNotAvailable", "register-tee"},
		{"OnlyExtensionOwner()", "OnlyExtensionOwner", "PRIVATE_KEY"},
		{"OnlyOwner()", "OnlyOwner", "PRIVATE_KEY"},
		{"InvalidUrl()", "InvalidUrl", "TUNNEL_URL"},
		{"InvalidTeeStatus()", "InvalidTeeStatus", "PRODUCTION"},
		{"VersionNotSupported()", "VersionNotSupported", "allow-tee-version"},
		{"InvalidAvailabilityCheckStatus()", "InvalidAvailabilityCheckStatus", "Try again"},
		{"AvailabilityCheckTimestampInvalid()", "AvailabilityCheckTimestampInvalid", "Re-run register-tee"},
		{"OnlyInstructionsSender()", "OnlyInstructionsSender", "INSTRUCTION_SENDER"},
	}

	for _, tt := range tests {
		selector := "0x" + hex.EncodeToString(crypto.Keccak256([]byte(tt.soliditySig))[:4])

		name := DecodeCustomError(selector)
		if name != tt.wantName {
			t.Errorf("DecodeCustomError(%s [%s]): got %q, want %q", tt.soliditySig, selector, name, tt.wantName)
		}

		hint := HintForRevert(selector)
		if hint == "" {
			t.Errorf("HintForRevert(%s [%s]): got empty hint", tt.soliditySig, selector)
		} else if !contains(hint, tt.wantHintSub) {
			t.Errorf("HintForRevert(%s): hint %q does not contain %q", tt.soliditySig, hint, tt.wantHintSub)
		}
	}
}

func TestDecodeCustomError_UnknownSelector(t *testing.T) {
	name := DecodeCustomError("0xdeadbeef")
	if name != "" {
		t.Errorf("expected empty for unknown selector, got %q", name)
	}

	hint := HintForRevert("0xdeadbeef")
	if hint != "" {
		t.Errorf("expected empty hint for unknown selector, got %q", hint)
	}
}

func TestDecodeCustomError_WithTrailingData(t *testing.T) {
	// Custom errors with parameters have trailing ABI-encoded data after the selector.
	// e.g. PlatformAlreadyExists(bytes32) has 4 + 32 bytes.
	selector := hex.EncodeToString(crypto.Keccak256([]byte("PlatformAlreadyExists(bytes32)"))[:4])
	rawHex := "0x" + selector + "0000000000000000000000000000000000000000000000000000000000000001"

	name := DecodeCustomError(rawHex)
	if name != "PlatformAlreadyExists" {
		t.Errorf("expected PlatformAlreadyExists, got %q", name)
	}
}

func TestExtractSelector_EdgeCases(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"0x", ""},
		{"0xaabb", ""},       // too short
		{"0xaabbccdd", "aabbccdd"},
		{"0XAABBCCDD", "aabbccdd"}, // uppercase prefix
		{"aabbccdd", "aabbccdd"},   // no prefix
	}

	for _, tt := range tests {
		got := extractSelector(tt.input)
		if got != tt.want {
			t.Errorf("extractSelector(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestExtractSelector_EmbeddedInErrorMessage(t *testing.T) {
	alreadyRegSel := "0x" + hex.EncodeToString(crypto.Keccak256([]byte("AlreadyRegistered()"))[:4])
	feeSelHex := hex.EncodeToString(crypto.Keccak256([]byte("FeeTooLow()"))[:4])

	tests := []struct {
		name     string
		input    string
		wantName string
		wantHint bool
	}{
		{"raw hex", alreadyRegSel, "AlreadyRegistered", true},
		{"embedded after execution reverted", "execution reverted: " + alreadyRegSel, "AlreadyRegistered", true},
		{"web3 ContractCustomError", "('" + alreadyRegSel + "', '" + alreadyRegSel + "')", "AlreadyRegistered", true},
		{"nested error", "error: execution reverted: 0x" + feeSelHex, "FeeTooLow", true},
		{"no selector", "execution reverted", "", false},
		{"empty", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName := DecodeCustomError(tt.input)
			if gotName != tt.wantName {
				t.Errorf("DecodeCustomError(%q) = %q, want %q", tt.input, gotName, tt.wantName)
			}
			gotHint := HintForRevert(tt.input)
			if tt.wantHint && gotHint == "" {
				t.Errorf("HintForRevert(%q) returned empty, want non-empty", tt.input)
			}
			if !tt.wantHint && gotHint != "" {
				t.Errorf("HintForRevert(%q) = %q, want empty", tt.input, gotHint)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
