package app

import (
	"context"
	"math/big"

	"sign-tools/app/contract"
	"sign-tools/base"
	"sign-tools/base/fccutils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// DeployInstructionSender deploys the InstructionSender contract and returns
// its address.
func DeployInstructionSender(s *base.Support) (common.Address, *contract.InstructionSender, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(s.Prv, s.ChainID)
	if err != nil {
		return common.Address{}, nil, errors.Errorf("failed to create transactor: %s", err)
	}

	address, tx, inst, err := contract.DeployInstructionSender(
		opts, s.ChainClient, s.Addresses.TeeExtensionRegistry, s.Addresses.TeeMachineRegistry,
	)
	if err != nil {
		return common.Address{}, nil, errors.Errorf("failed to deploy contract: %s", err)
	}

	receipt, err := bind.WaitMined(context.Background(), s.ChainClient, tx)
	if err != nil {
		return common.Address{}, nil, errors.Errorf("failed waiting for deployment: %s", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Address{}, nil, errors.New("contract deployment failed")
	}

	return address, inst, nil
}

// SetExtensionId calls setExtensionId on the InstructionSender contract.
// This is idempotent — it skips if the extension ID is already set.
func SetExtensionId(s *base.Support, instructionSenderAddress common.Address) error {
	sender, err := contract.NewInstructionSender(instructionSenderAddress, s.ChainClient)
	if err != nil {
		return errors.Errorf("failed to bind contract: %s", err)
	}

	// Check if already set.
	extID, err := sender.ExtensionId(&bind.CallOpts{})
	if err != nil {
		return errors.Errorf("failed to read extensionId: %s", err)
	}
	if extID.Sign() != 0 {
		return nil // already set
	}

	opts, err := bind.NewKeyedTransactorWithChainID(s.Prv, s.ChainID)
	if err != nil {
		return errors.Errorf("failed to create transactor: %s", err)
	}

	tx, err := sender.SetExtensionId(opts)
	if err != nil {
		reason := fccutils.DecodeRevertReason(err)
		if reason == "" {
			parsed, _ := contract.InstructionSenderMetaData.GetAbi()
			if parsed != nil {
				callData, packErr := parsed.Pack("setExtensionId")
				if packErr == nil {
					from := crypto.PubkeyToAddress(s.Prv.PublicKey)
					reason = fccutils.SimulateAndDecodeRevert(
						s.ChainClient, from, instructionSenderAddress, nil, callData,
					)
				}
			}
		}
		if reason != "" {
			return errors.Errorf("failed to call setExtensionId: %s (revert reason: %s)", err, reason)
		}
		return errors.Errorf("failed to call setExtensionId: %s", err)
	}

	receipt, err := bind.WaitMined(context.Background(), s.ChainClient, tx)
	if err != nil {
		return errors.Errorf("failed waiting for transaction: %s", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		parsed, _ := contract.InstructionSenderMetaData.GetAbi()
		if parsed != nil {
			callData, packErr := parsed.Pack("setExtensionId")
			if packErr == nil {
				from := crypto.PubkeyToAddress(s.Prv.PublicKey)
				reason := fccutils.SimulateAndDecodeRevert(
					s.ChainClient, from, instructionSenderAddress, nil, callData,
				)
				if reason != "" {
					return errors.Errorf("setExtensionId transaction failed (revert reason: %s)", reason)
				}
			}
		}
		return errors.New("setExtensionId transaction failed")
	}

	return nil
}

// SendUpdateKey sends an updateKey instruction via the InstructionSender.
// Returns the instruction ID parsed from the TeeInstructionsSent event.
func SendUpdateKey(s *base.Support, instructionSenderAddress common.Address, encryptedKey []byte, fee *big.Int) (common.Hash, error) {
	sender, err := contract.NewInstructionSender(instructionSenderAddress, s.ChainClient)
	if err != nil {
		return common.Hash{}, errors.Errorf("failed to bind contract: %s", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(s.Prv, s.ChainID)
	if err != nil {
		return common.Hash{}, errors.Errorf("failed to create transactor: %s", err)
	}
	opts.Value = fee

	tx, err := sender.UpdateKey(opts, encryptedKey)
	if err != nil {
		return common.Hash{}, errors.Errorf("updateKey: %s", err)
	}

	receipt, err := bind.WaitMined(context.Background(), s.ChainClient, tx)
	if err != nil {
		return common.Hash{}, errors.Errorf("wait updateKey: %s", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Hash{}, errors.New("updateKey transaction failed")
	}

	return parseInstructionID(receipt, s)
}

// SendSign sends a sign instruction via the InstructionSender.
// Returns the instruction ID parsed from the TeeInstructionsSent event.
func SendSign(s *base.Support, instructionSenderAddress common.Address, message []byte, fee *big.Int) (common.Hash, error) {
	sender, err := contract.NewInstructionSender(instructionSenderAddress, s.ChainClient)
	if err != nil {
		return common.Hash{}, errors.Errorf("failed to bind contract: %s", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(s.Prv, s.ChainID)
	if err != nil {
		return common.Hash{}, errors.Errorf("failed to create transactor: %s", err)
	}
	opts.Value = fee

	tx, err := sender.Sign(opts, message)
	if err != nil {
		return common.Hash{}, errors.Errorf("sign: %s", err)
	}

	receipt, err := bind.WaitMined(context.Background(), s.ChainClient, tx)
	if err != nil {
		return common.Hash{}, errors.Errorf("wait sign: %s", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Hash{}, errors.New("sign transaction failed")
	}

	return parseInstructionID(receipt, s)
}

// parseInstructionID extracts the instruction ID from the TeeInstructionsSent event.
func parseInstructionID(receipt *types.Receipt, s *base.Support) (common.Hash, error) {
	if s.TeeExtensionRegistry != nil && len(receipt.Logs) > 0 {
		event, err := s.TeeExtensionRegistry.ParseTeeInstructionsSent(*receipt.Logs[0])
		if err == nil {
			return common.Hash(event.InstructionId), nil
		}
	}
	// Fallback: read Topics[2] from the first log
	if len(receipt.Logs) > 0 && len(receipt.Logs[0].Topics) >= 3 {
		return receipt.Logs[0].Topics[2], nil
	}
	return common.Hash{}, errors.New("could not extract instruction ID from receipt")
}
