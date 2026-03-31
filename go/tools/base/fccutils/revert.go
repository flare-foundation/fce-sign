package fccutils

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"sign-tools/base"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// DecodeRevertReason attempts to extract and decode the revert reason from
// an error returned by eth_call or eth_estimateGas.
func DecodeRevertReason(err error) string {
	if err == nil {
		return ""
	}

	type dataError interface {
		ErrorData() interface{}
	}

	var de dataError
	if errors.As(err, &de) {
		if data := de.ErrorData(); data != nil {
			if hexStr, ok := data.(string); ok {
				return decodeRevertHex(hexStr)
			}
		}
	}

	return ""
}

// SimulateAndDecodeRevert replays a call via eth_call and attempts to decode
// the revert reason.
func SimulateAndDecodeRevert(
	client *ethclient.Client,
	from common.Address,
	to common.Address,
	value *big.Int,
	data []byte,
) string {
	toAddr := to
	msg := ethereum.CallMsg{
		From:  from,
		To:    &toAddr,
		Value: value,
		Data:  data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		if reason := DecodeRevertReason(err); reason != "" {
			return reason
		}
		return err.Error()
	}

	if len(result) >= 4 {
		return decodeRevertHex(hex.EncodeToString(result))
	}

	return ""
}

func decodeRevertHex(hexStr string) string {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	decoded, err := hex.DecodeString(hexStr)
	if err != nil || len(decoded) < 4 {
		return ""
	}

	// Standard Error(string) revert.
	if reason, unpackErr := abi.UnpackRevert(decoded); unpackErr == nil {
		return reason
	}

	// Try matching a known custom error selector.
	withPrefix := "0x" + hexStr
	if name := base.DecodeCustomError(withPrefix); name != "" {
		return withPrefix
	}

	return "0x" + hexStr
}
