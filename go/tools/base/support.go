package base

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/fdc2hub"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/system"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/teeextensionregistry"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/teemachineregistry"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/teeownerallowlist"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/teeverification"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/teewalletkeymanager"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/teewalletmanager"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/teewalletprojectmanager"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// Support holds an authenticated chain client and all standard TEE contract
// bindings needed by the CLI tools.
type Support struct {
	Prv *ecdsa.PrivateKey

	FlareSystemManager      *system.FlareSystemsManager
	TeeMachineRegistry      *teemachineregistry.TeeMachineRegistry
	TeeWalletProjectManager *teewalletprojectmanager.TeeWalletProjectManager
	TeeWalletManager        *teewalletmanager.TeeWalletManager
	TeeWalletKeyManager     *teewalletkeymanager.TeeWalletKeyManager
	Fdc2Hub                 *fdc2hub.Fdc2Hub
	TeeVerification         *teeverification.TeeVerification
	TeeExtensionRegistry    *teeextensionregistry.TeeExtensionRegistry
	TeeOwnerAllowlist       *teeownerallowlist.TeeOwnerAllowlist

	Addresses *Addresses

	ChainClient *ethclient.Client
	ChainID     *big.Int
}

// Addresses contains the on-chain addresses of the TEE protocol contracts.
type Addresses struct {
	TeeMachineRegistry      common.Address `json:"TeeMachineRegistry"`
	TeeWalletManager        common.Address `json:"TeeWalletManager"`
	TeeWalletKeyManager     common.Address `json:"TeeWalletKeyManager"`
	TeeWalletProjectManager common.Address `json:"TeeWalletProjectManager"`
	FlareSystemManager      common.Address `json:"FlareSystemsManager"`
	Fdc2Hub                 common.Address `json:"Fdc2Hub"`
	TeeVerification         common.Address `json:"TeeVerification"`
	TeeExtensionRegistry    common.Address `json:"TeeExtensionRegistry"`
	TeeOwnerAllowlist       common.Address `json:"TeeOwnerAllowlist"`
}

// DefaultSupport loads .env, reads addresses, connects to the chain, and
// returns a fully initialised Support.
func DefaultSupport(addressesFilePath, chainNodeURL string) (*Support, error) {
	addr := &Addresses{}
	err := ReadAddresses(addressesFilePath, addr)
	if err != nil {
		// Fallback: try array-format JSON
		addr, err = ParseAddresses(addressesFilePath)
		if err != nil {
			return nil, errors.Errorf("read addresses: %s", err)
		}
	}

	cc, err := ethclient.Dial(chainNodeURL)
	if err != nil {
		return nil, errors.Errorf("dial chain: %s", err)
	}

	privKey, err := DefaultPrivateKey()
	if err != nil {
		return nil, err
	}

	return NewSupport(privKey, cc, addr)
}

// DefaultPrivateKey loads the private key from PRIVATE_KEY in the environment,
// optionally loading .env first.
func DefaultPrivateKey() (*ecdsa.PrivateKey, error) {
	// Try loading .env from cwd first, then from project root (for running from go/tools/).
	if err := godotenv.Load(); err != nil {
		if err2 := godotenv.Load("../../.env"); err2 != nil {
			fmt.Printf("Warning: .env not loaded: %v\n", err)
		}
	}

	privKeyString := os.Getenv("PRIVATE_KEY")

	if privKeyString == "" {
		fmt.Println("Warning: PRIVATE_KEY not set, falling back to hardcoded dev key (only works on local devnet)")
		return PrvWithFunds, nil
	}

	privKeyString = strings.TrimPrefix(strings.TrimPrefix(privKeyString, "0x"), "0X")
	privKey, err := crypto.HexToECDSA(privKeyString)
	if err != nil {
		return nil, errors.Errorf("parse private key: %s", err)
	}
	return privKey, nil
}

// NewSupport wires up all contract bindings from the given key, client and
// addresses.
func NewSupport(prv *ecdsa.PrivateKey, chainClient *ethclient.Client, addresses *Addresses) (*Support, error) {
	tr, err := teemachineregistry.NewTeeMachineRegistry(addresses.TeeMachineRegistry, chainClient)
	if err != nil {
		return nil, err
	}
	twm, err := teewalletmanager.NewTeeWalletManager(addresses.TeeWalletManager, chainClient)
	if err != nil {
		return nil, err
	}
	twkm, err := teewalletkeymanager.NewTeeWalletKeyManager(addresses.TeeWalletKeyManager, chainClient)
	if err != nil {
		return nil, err
	}
	twpm, err := teewalletprojectmanager.NewTeeWalletProjectManager(addresses.TeeWalletProjectManager, chainClient)
	if err != nil {
		return nil, err
	}
	sm, err := system.NewFlareSystemsManager(addresses.FlareSystemManager, chainClient)
	if err != nil {
		return nil, err
	}
	ftdc, err := fdc2hub.NewFdc2Hub(addresses.Fdc2Hub, chainClient)
	if err != nil {
		return nil, err
	}
	tv, err := teeverification.NewTeeVerification(addresses.TeeVerification, chainClient)
	if err != nil {
		return nil, err
	}
	ter, err := teeextensionregistry.NewTeeExtensionRegistry(addresses.TeeExtensionRegistry, chainClient)
	if err != nil {
		return nil, err
	}
	toal, err := teeownerallowlist.NewTeeOwnerAllowlist(addresses.TeeOwnerAllowlist, chainClient)
	if err != nil {
		return nil, err
	}

	chainID, err := chainClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	return &Support{
		Prv:                     prv,
		TeeMachineRegistry:      tr,
		TeeWalletManager:        twm,
		TeeWalletKeyManager:     twkm,
		TeeWalletProjectManager: twpm,
		FlareSystemManager:      sm,
		Fdc2Hub:                 ftdc,
		TeeVerification:         tv,
		TeeExtensionRegistry:    ter,
		TeeOwnerAllowlist:       toal,
		ChainClient:             chainClient,
		ChainID:                 chainID,
		Addresses:               addresses,
	}, nil
}

// SendAndCheck sends a transaction produced by a contract binding and waits for
// the receipt. If the binding's estimateGas fails, it performs an explicit
// eth_call dry-run to extract a decoded revert reason with actionable hints.
//
// Usage:
//
//	tx, err := s.TeeMachineRegistry.Register(opts, ...)
//	receipt, err := base.SendAndCheck(tx, err, s)
func SendAndCheck(tx *types.Transaction, bindErr error, s *Support) (*types.Receipt, error) {
	if bindErr != nil {
		// The binding failed (usually estimateGas revert). Try to get a better
		// error by replaying via eth_call using the binding's error data.
		rawHex := extractRevertFromBindError(bindErr)
		if rawHex != "" {
			return nil, revertError(rawHex)
		}
		return nil, errors.Errorf("%s", bindErr)
	}
	return checkTxReceipt(tx, s.ChainClient)
}

// CheckTx waits for a transaction to be mined and returns the receipt.
// It returns an error with the revert reason if the tx failed. When the
// revert matches a known custom error, the error message includes an
// actionable hint for the developer.
func CheckTx(tx *types.Transaction, client *ethclient.Client) (*types.Receipt, error) {
	return checkTxReceipt(tx, client)
}

func checkTxReceipt(tx *types.Transaction, client *ethclient.Client) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return nil, errors.Errorf("wait mined: %s", err)
	}
	if receipt.Status == 0 {
		rawHex := getRevertData(client, tx)
		return nil, revertError(rawHex)
	}
	return receipt, nil
}

func revertError(rawHex string) error {
	if errName := DecodeCustomError(rawHex); errName != "" {
		hint := HintForRevert(rawHex)
		if hint != "" {
			return errors.Errorf("transaction reverted: %s\n  → %s", errName, hint)
		}
		return errors.Errorf("transaction reverted: %s", errName)
	}
	return errors.Errorf("transaction failed: %s", rawHex)
}

// extractRevertFromBindError tries to get the revert selector from a contract
// binding error (estimateGas failure).
func extractRevertFromBindError(err error) string {
	// Try ErrorData() interface (go-ethereum JSON-RPC errors).
	type dataError interface {
		ErrorData() interface{}
	}
	var de dataError
	if stderrors.As(err, &de) {
		if data := de.ErrorData(); data != nil {
			if hexStr, ok := data.(string); ok {
				hexStr = strings.TrimPrefix(hexStr, "0x")
				if len(hexStr) >= 8 {
					return "0x" + hexStr
				}
			}
		}
	}
	return ""
}

// getRevertData replays a failed transaction via eth_call and returns the raw
// hex revert data (e.g. "0x..."). On any error it returns a best-effort message.
func getRevertData(client *ethclient.Client, tx *types.Transaction) string {
	txObj, _, err := client.TransactionByHash(context.Background(), tx.Hash())
	if err != nil {
		return err.Error()
	}
	from, err := types.Sender(types.NewEIP155Signer(txObj.ChainId()), txObj)
	if err != nil {
		return err.Error()
	}

	toAddr := txObj.To()
	msg := ethereum.CallMsg{
		From:     from,
		To:       toAddr,
		Gas:      txObj.Gas(),
		GasPrice: txObj.GasPrice(),
		Value:    txObj.Value(),
		Data:     txObj.Data(),
	}

	// eth_call may return the revert data in the error itself.
	_, callErr := client.CallContract(context.Background(), msg, nil)
	if callErr != nil {
		// Try to extract structured revert data from the error.
		if reason := decodeRevertFromError(callErr); reason != "" {
			return reason
		}
		return callErr.Error()
	}

	return "unknown revert reason"
}

// decodeRevertFromError attempts to extract hex revert data from a JSON-RPC
// error that implements ErrorData().
func decodeRevertFromError(err error) string {
	type dataError interface {
		ErrorData() interface{}
	}

	var de dataError
	if stderrors.As(err, &de) {
		if data := de.ErrorData(); data != nil {
			if hexStr, ok := data.(string); ok {
				hexStr = strings.TrimPrefix(hexStr, "0x")
				if len(hexStr) >= 8 {
					return "0x" + hexStr
				}
			}
		}
	}
	return ""
}

// RawContract is used when parsing the array-format addresses JSON.
type RawContract struct {
	Name         string `json:"name"`
	ContractName string `json:"contractName"`
	Address      string `json:"address"`
}

// ParseAddresses reads addresses from the old array-format JSON as a fallback.
func ParseAddresses(filePath string) (*Addresses, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var raw []RawContract
	if err := json.Unmarshal(file, &raw); err != nil {
		return nil, err
	}

	var addrs Addresses
	targets := fieldMap(&addrs)
	for _, c := range raw {
		if dest, ok := targets[c.Name]; ok {
			*dest = common.HexToAddress(c.Address)
		}
	}
	return &addrs, nil
}

// fieldMap builds "<json tag or field name>" -> pointer to field.
func fieldMap(addrStruct *Addresses) map[string]*common.Address {
	out := make(map[string]*common.Address)
	v := reflect.ValueOf(addrStruct).Elem()
	t := v.Type()
	addrType := reflect.TypeOf(common.Address{})

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)
		if fv.Type() != addrType {
			continue
		}
		tag := ft.Tag.Get("json")
		key := tag
		if key == "" || key == "-" {
			key = ft.Name
		}
		out[key] = fv.Addr().Interface().(*common.Address)
	}
	return out
}
