package base

import (
	"encoding/hex"
	"encoding/json"
	"strings"
)

// Action is the top-level request received on POST /action.
type Action struct {
	Data                       ActionData `json:"data"`
	AdditionalVariableMessages []string   `json:"additionalVariableMessages,omitempty"`
	Timestamps                 []int64    `json:"timestamps,omitempty"`
	AdditionalActionData       string     `json:"additionalActionData,omitempty"`
	Signatures                 []string   `json:"signatures,omitempty"`
}

// ActionData is the nested "data" field inside an Action.
type ActionData struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	SubmissionTag string `json:"submissionTag"`
	Message       string `json:"message"` // JSON-encoded DataFixed
}

// DataFixed is the decoded content of ActionData.Message.
type DataFixed struct {
	InstructionID          string   `json:"instructionId"`
	TeeID                  string   `json:"teeId,omitempty"`
	Timestamp              int64    `json:"timestamp,omitempty"`
	RewardEpochID          int64    `json:"rewardEpochId,omitempty"`
	OpType                 string   `json:"opType"`
	OpCommand              string   `json:"opCommand"`
	Cosigners              []string `json:"cosigners,omitempty"`
	CosignersThreshold     int64    `json:"cosignersThreshold,omitempty"`
	OriginalMessage        string   `json:"originalMessage,omitempty"`
	AdditionalFixedMessage string   `json:"additionalFixedMessage,omitempty"`
}

// ActionResult is the response returned from POST /action.
type ActionResult struct {
	ID                     string  `json:"id"`
	SubmissionTag          string  `json:"submissionTag"`
	Status                 int     `json:"status"`
	Log                    *string `json:"log,omitempty"`
	OpType                 string  `json:"opType"`
	OpCommand              string  `json:"opCommand"`
	AdditionalResultStatus *string `json:"additionalResultStatus,omitempty"`
	Version                string  `json:"version"`
	Data                   *string `json:"data,omitempty"`
}

// StateResponse is returned from GET /state.
type StateResponse struct {
	StateVersion string          `json:"stateVersion"`
	State        json.RawMessage `json:"state"`
}

// HandlerFunc is the signature for action handlers.
// It receives the original message bytes and returns (data, status, error).
type HandlerFunc func(msg string) (data *string, status int, err error)

// RegisterFunc is called at startup to register handlers and initial state.
type RegisterFunc func(f *Framework)

// ReportStateFunc converts the current state into a JSON-serializable value.
type ReportStateFunc func() json.RawMessage

// Framework provides setState and handle registration to app code.
type Framework struct {
	handlers []handlerEntry
}

type handlerEntry struct {
	opType    string
	opCommand string
	handler   HandlerFunc
}

// Handle registers a handler for an OPType/OPCommand pair.
// Pass "" for opCommand to match any command.
func (f *Framework) Handle(opType, opCommand string, handler HandlerFunc) {
	f.handlers = append(f.handlers, handlerEntry{
		opType:    stringToBytes32Hex(opType),
		opCommand: stringToBytes32Hex(opCommand),
		handler:   handler,
	})
}

// stringToBytes32Hex encodes a UTF-8 string into a 0x-prefixed 32-byte
// zero-right-padded hex string, matching Solidity's bytes32("FOO").
func stringToBytes32Hex(s string) string {
	b := make([]byte, 32)
	copy(b, []byte(s))
	return "0x" + hex.EncodeToString(b)
}

// VersionToHex converts a version string to bytes32 hex.
func VersionToHex(version string) string {
	return stringToBytes32Hex(version)
}

// Lookup finds a handler matching the given opType and opCommand.
func (f *Framework) Lookup(opType, opCommand string) HandlerFunc {
	emptyCmd := stringToBytes32Hex("")
	for _, e := range f.handlers {
		if e.opType != opType {
			continue
		}
		if e.opCommand == emptyCmd || e.opCommand == opCommand {
			return e.handler
		}
	}
	return nil
}

// OpTypeToString converts a bytes32 hex opType back to a trimmed string.
func OpTypeToString(h string) string {
	h = strings.TrimPrefix(h, "0x")
	b, err := hex.DecodeString(h)
	if err != nil {
		return ""
	}
	end := len(b)
	for end > 0 && b[end-1] == 0 {
		end--
	}
	return string(b[:end])
}
