package app

// Version is the SemVer version of this extension.
const Version = "0.1.0"

// OPType and OPCommand constants — must match the bytes32 constants in InstructionSender.sol.
const (
	OpTypeKey       = "KEY"
	OpCommandUpdate = "UPDATE"
	OpCommandSign   = "SIGN"
)
