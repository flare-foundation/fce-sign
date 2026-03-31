package base

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

type revertHint struct {
	errorName string
	hint      string
}

// revertHints maps 4-byte custom error selectors (hex-encoded, no 0x) to
// human-readable error names and actionable next-step guidance.
var revertHints map[string]revertHint

func init() {
	type entry struct {
		sig  string // Solidity error signature, e.g. "FeeTooLow()"
		name string
		hint string
	}

	entries := []entry{
		{"VersionAlreadyExists()", "VersionAlreadyExists",
			"This code hash is already registered for this extension. You can skip this step."},
		{"PlatformAlreadyExists(bytes32)", "PlatformAlreadyExists",
			"Platform already registered for this code hash. You can skip this step."},
		{"AlreadyRegistered()", "AlreadyRegistered",
			"This TEE machine is already registered on-chain. To register a new TEE, restart the Docker stack (docker compose down && docker compose up -d) to generate a new TEE identity, then re-run allow-tee-version and register-tee."},
		{"OwnerAlreadyAllowed(address)", "OwnerAlreadyAllowed",
			"This owner is already in the allowlist. Safe to continue."},
		{"KeyTypeAlreadyExists(bytes32)", "KeyTypeAlreadyExists",
			"Key type already supported on this extension. Safe to continue."},
		{"FeeTooLow()", "FeeTooLow",
			"The fee sent with this instruction is below the on-chain minimum. Increase FEE_WEI in .env and re-run."},
		{"TeeMachineNotAvailable()", "TeeMachineNotAvailable",
			"TEE machine is not in PRODUCTION status. Complete the register-tee step first."},
		{"OnlyExtensionOwner()", "OnlyExtensionOwner",
			"Only the extension owner can perform this action. Check that PRIVATE_KEY in .env matches the key used to register the extension."},
		{"OnlyOwner()", "OnlyOwner",
			"Only the TEE machine owner can perform this action. Check that PRIVATE_KEY and INITIAL_OWNER in .env are correct."},
		{"InvalidUrl()", "InvalidUrl",
			"Invalid or empty URL. Set TUNNEL_URL in .env to your public tunnel URL."},
		{"InvalidTeeStatus()", "InvalidTeeStatus",
			"TEE is in an unexpected status for this operation. If the TEE is already in PRODUCTION, you can skip this step."},
		{"VersionNotSupported()", "VersionNotSupported",
			"Code hash or platform is not registered for this extension. Run allow-tee-version first."},
		{"InvalidAvailabilityCheckStatus()", "InvalidAvailabilityCheckStatus",
			"The availability check did not return OK. The FTDC TEE may not have processed the request yet. Try again in a moment."},
		{"AvailabilityCheckTimestampInvalid()", "AvailabilityCheckTimestampInvalid",
			"Availability check timestamp is before the TEE's last status change. Re-run register-tee to get a fresh availability check."},
		{"OnlyInstructionsSender()", "OnlyInstructionsSender",
			"The calling contract is not the registered instructions sender for this extension. Redeploy the InstructionSender contract and update INSTRUCTION_SENDER in .env."},
	}

	revertHints = make(map[string]revertHint, len(entries))
	for _, e := range entries {
		selector := hex.EncodeToString(crypto.Keccak256([]byte(e.sig))[:4])
		revertHints[selector] = revertHint{errorName: e.name, hint: e.hint}
	}
}

// HintForRevert returns an actionable hint for the given raw revert hex or
// message string. Returns empty string if no hint matches.
func HintForRevert(rawHexOrMessage string) string {
	sel := extractSelector(rawHexOrMessage)
	if sel == "" {
		return ""
	}
	if h, ok := revertHints[sel]; ok {
		return h.hint
	}
	return ""
}

// DecodeCustomError returns the human-readable error name (e.g.
// "AlreadyRegistered") for the given raw hex revert data, or empty string.
func DecodeCustomError(rawHex string) string {
	sel := extractSelector(rawHex)
	if sel == "" {
		return ""
	}
	if h, ok := revertHints[sel]; ok {
		return h.errorName
	}
	return ""
}

func extractSelector(s string) string {
	// First try: the string itself is hex data (e.g. "0xceb05b68..." or "ceb05b68...").
	stripped := strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(stripped) >= 8 && isHex(stripped[:8]) {
		return strings.ToLower(stripped[:8])
	}

	// Second try: find a 0x-prefixed hex selector embedded in the string
	// (e.g. "execution reverted: 0xceb05b68" or error messages containing selectors).
	for i := 0; i < len(s)-9; i++ {
		if s[i] == '0' && (s[i+1] == 'x' || s[i+1] == 'X') {
			candidate := s[i+2:]
			if len(candidate) >= 8 && isHex(candidate[:8]) {
				return strings.ToLower(candidate[:8])
			}
		}
	}

	return ""
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
