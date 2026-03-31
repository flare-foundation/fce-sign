package fccutils

import (
	stderrors "errors"

	"sign-tools/base"

	"github.com/ethereum/go-ethereum/common"
	"github.com/flare-foundation/go-flare-common/pkg/logger"
	"github.com/flare-foundation/go-flare-common/pkg/tee/structs"
	"github.com/flare-foundation/go-flare-common/pkg/tee/structs/connector"
	"github.com/pkg/errors"
)

var (
	PlatformIntel   common.Hash = common.HexToHash("4743505f494e54454c5f54445800000000000000000000000000000000000000") // GCP_INTEL_TDX
	PlatformAMD     common.Hash = common.HexToHash("4743505f414d445f534556000000000000000000000000000000000000000000") // GCP_AMD_SEV
	PlatformAMDESEV common.Hash = common.HexToHash("4743505f414d445f5345565f4553000000000000000000000000000000000000") // GCP_AMD_SEV_ES
	TestPlatform    common.Hash = common.HexToHash("544553545f504c4154464f524d00000000000000000000000000000000000000") // TEST_PLATFORM
	TeeCodeHash     common.Hash = common.HexToHash("194844cf417dde867073e5ab7199fa4d21fd82b5dbe2bdea8b3d7fc18d10fdc2")
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func FatalWithCause(err error) {
	msg := err.Error()

	// Try to extract structured revert data from the error.
	revertHex := extractRevertHex(err)
	if revertHex == "" {
		revertHex = msg
	}

	if hint := base.HintForRevert(revertHex); hint != "" {
		if name := base.DecodeCustomError(revertHex); name != "" {
			msg = msg + " (" + name + ")\n  → " + hint
		} else {
			msg = msg + "\n  → " + hint
		}
	}

	errCause, ok := err.(StackTracer)
	if ok {
		st := errCause.StackTrace()
		logger.Fatalf("Error: %s %+v", msg, st)
	} else {
		logger.Fatalf("Error: %s", msg)
	}
}

// extractRevertHex tries to get raw hex revert data from an error that
// implements the ErrorData() interface (as go-ethereum JSON-RPC errors do).
func extractRevertHex(err error) string {
	type dataError interface {
		ErrorData() interface{}
	}
	var de dataError
	if stderrors.As(err, &de) {
		if data := de.ErrorData(); data != nil {
			if hexStr, ok := data.(string); ok {
				return hexStr
			}
		}
	}
	return ""
}

func EncodeFTDCTeeAvailabilityCheckRequest(data connector.ITeeAvailabilityCheckRequestBody) ([]byte, error) {
	return structs.Encode(connector.AttestationTypeArguments[connector.AvailabilityCheck].Request, data)
}

func DecodeFTDCTeeAvailabilityCheckRequest(data []byte) (connector.ITeeAvailabilityCheckRequestBody, error) {
	var request connector.ITeeAvailabilityCheckRequestBody
	err := structs.DecodeTo(connector.AttestationTypeArguments[connector.AvailabilityCheck].Request, data, &request)
	if err != nil {
		return connector.ITeeAvailabilityCheckRequestBody{}, errors.Errorf("%s", err)
	}
	return request, nil
}

func EncodeFTDCTeeAvailabilityCheckResponse(data connector.ITeeAvailabilityCheckResponseBody) ([]byte, error) {
	return structs.Encode(connector.AttestationTypeArguments[connector.AvailabilityCheck].Response, data)
}

func DecodeFTDCTeeAvailabilityCheckResponse(data []byte) (connector.ITeeAvailabilityCheckResponseBody, error) {
	var request connector.ITeeAvailabilityCheckResponseBody
	err := structs.DecodeTo(connector.AttestationTypeArguments[connector.AvailabilityCheck].Response, data, &request)
	if err != nil {
		return connector.ITeeAvailabilityCheckResponseBody{}, errors.Errorf("%s", err)
	}
	return request, nil
}
