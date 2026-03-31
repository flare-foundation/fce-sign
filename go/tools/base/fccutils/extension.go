package fccutils

import (
	"context"
	"math/big"

	"sign-tools/base"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/teeextensionregistry"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/teeownerallowlist"
	"github.com/flare-foundation/go-flare-common/pkg/logger"
	"github.com/flare-foundation/tee-node/pkg/wallets"
	"github.com/pkg/errors"
)

var DefaultExtensionId = big.NewInt(0)

// SetupExtension registers an extension, allows the caller as TEE machine
// owner and wallet project manager, and adds EVM key type support.
// Steps that are already complete are skipped gracefully.
func SetupExtension(
	s *base.Support,
	governanceHash common.Hash,
	instructionsSenderAddress, stateVerifierAddress common.Address,
) (*big.Int, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(s.Prv, s.ChainID)
	if err != nil {
		return nil, err
	}

	extRegistered, _, err := registerExtension(s, opts, instructionsSenderAddress, stateVerifierAddress)
	if err != nil {
		return nil, err
	}
	extensionID := extRegistered.ExtensionId
	logger.Infof("Extension registered with ID: %s", extensionID.String())

	callerAddr := crypto.PubkeyToAddress(s.Prv.PublicKey)

	// Allow TEE machine owner (skip if already allowed).
	alreadyMachineOwner, err := s.TeeOwnerAllowlist.IsAllowedTeeMachineOwner(nil, extensionID, callerAddr)
	if err != nil {
		return nil, errors.Errorf("check TEE machine owner allowlist: %s", err)
	}
	if alreadyMachineOwner {
		logger.Infof("TEE machine owner already allowed, skipping.")
	} else {
		_, err = allowTeeMachineOwners(s, opts, extensionID, []common.Address{callerAddr})
		if err != nil {
			return nil, err
		}
		logger.Infof("TEE machine owner allowed: %s", callerAddr.Hex())
	}

	// Allow wallet project manager owner (skip if already allowed).
	alreadyProjectOwner, err := s.TeeOwnerAllowlist.IsAllowedTeeWalletProjectOwner(nil, extensionID, callerAddr)
	if err != nil {
		return nil, errors.Errorf("check wallet project owner allowlist: %s", err)
	}
	if alreadyProjectOwner {
		logger.Infof("Wallet project manager owner already allowed, skipping.")
	} else {
		_, err = allowTeeProjectManagerOwners(s, opts, extensionID, []common.Address{callerAddr})
		if err != nil {
			return nil, err
		}
		logger.Infof("Wallet project manager owner allowed: %s", callerAddr.Hex())
	}

	// Add EVM key type (skip if already supported).
	isKeyTypeSupported, err := IsKeyTypeSupported(s, extensionID, wallets.EVMType)
	if err != nil {
		return nil, err
	}
	if isKeyTypeSupported {
		logger.Infof("Key type %s already supported, skipping.", wallets.EVMType)
	} else {
		logger.Infof("Adding key type %s to extension %s", wallets.EVMType, extensionID)
		err = AddSupportedKeyTypes(s, extensionID, []common.Hash{wallets.EVMType})
		if err != nil {
			return nil, err
		}
	}

	return extensionID, nil
}

func AddSupportedKeyTypes(s *base.Support, extensionId *big.Int, keyTypes []common.Hash) error {
	opts, err := bind.NewKeyedTransactorWithChainID(s.Prv, s.ChainID)
	if err != nil {
		return errors.Errorf("%s", err)
	}

	keyTypesBytes32 := HashArrayToBytes32Array(keyTypes)

	tx, err := s.TeeExtensionRegistry.AddSupportedKeyTypes(opts, extensionId, keyTypesBytes32)
	_, err = base.SendAndCheck(tx, err, s)
	if err != nil {
		return err
	}
	return nil
}

func IsKeyTypeSupported(s *base.Support, extensionId *big.Int, keyType common.Hash) (bool, error) {
	callOpts := &bind.CallOpts{
		From:    crypto.PubkeyToAddress(s.Prv.PublicKey),
		Context: context.Background(),
	}
	return s.TeeExtensionRegistry.IsKeyTypeSupported(callOpts, extensionId, keyType)
}

func registerExtension(
	s *base.Support, opts *bind.TransactOpts, instructionsSenderAddress, stateVerifierAddress common.Address,
) (
	*teeextensionregistry.TeeExtensionRegistryTeeExtensionRegistered, *teeextensionregistry.TeeExtensionRegistryTeeExtensionContractsSet, error,
) {
	tx, err := s.TeeExtensionRegistry.Register(opts, stateVerifierAddress, instructionsSenderAddress)
	receipt, err := base.SendAndCheck(tx, err, s)
	if err != nil {
		return nil, nil, err
	}

	extensionRegistered, err := s.TeeExtensionRegistry.ParseTeeExtensionRegistered(*receipt.Logs[0])
	if err != nil {
		return nil, nil, errors.Errorf("failed to parse TeeExtensionRegistered event: %s", err)
	}

	extensionContractsSet, err := s.TeeExtensionRegistry.ParseTeeExtensionContractsSet(*receipt.Logs[1])
	if err != nil {
		return nil, nil, errors.Errorf("failed to parse TeeExtensionContractsSet event: %s", err)
	}

	return extensionRegistered, extensionContractsSet, nil
}

func allowTeeMachineOwners(s *base.Support, opts *bind.TransactOpts, extensionId *big.Int, owners []common.Address) (*teeownerallowlist.TeeOwnerAllowlistAllowedTeeMachineOwnersAdded, error) {
	tx, err := s.TeeOwnerAllowlist.AddAllowedTeeMachineOwners(opts, extensionId, owners)
	receipt, err := base.SendAndCheck(tx, err, s)
	if err != nil {
		return nil, err
	}

	ownersAdded, err := s.TeeOwnerAllowlist.ParseAllowedTeeMachineOwnersAdded(*receipt.Logs[0])
	if err != nil {
		return nil, errors.Errorf("failed to parse AllowedTeeMachineOwnersAdded event: %s", err)
	}

	return ownersAdded, nil
}

func allowTeeProjectManagerOwners(s *base.Support, opts *bind.TransactOpts, extensionId *big.Int, owners []common.Address) (*teeownerallowlist.TeeOwnerAllowlistAllowedTeeWalletProjectOwnersAdded, error) {
	tx, err := s.TeeOwnerAllowlist.AddAllowedTeeWalletProjectOwners(opts, extensionId, owners)
	receipt, err := base.SendAndCheck(tx, err, s)
	if err != nil {
		return nil, err
	}

	ownersAdded, err := s.TeeOwnerAllowlist.ParseAllowedTeeWalletProjectOwnersAdded(*receipt.Logs[0])
	if err != nil {
		return nil, errors.Errorf("failed to parse AllowedTeeWalletProjectOwnersAdded event: %s", err)
	}

	return ownersAdded, nil
}
