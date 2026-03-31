package base

import (
	"crypto/ecdsa"
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// Coston2 defaults
	DefaultExtensionProxyURL = "http://localhost:6676"
	DefaultChainNodeURL      = "https://coston2-api.flare.network/ext/C/rpc"
	DefaultAddressesFile     = "../../config/coston2/deployed-addresses.json"
)

// ReadAddresses unmarshals a flat JSON addresses file into dest.
func ReadAddresses[T any](filePath string, dest *T) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, dest)
}

// PrvWithFunds is a hardcoded dev-only key for local devnet testing.
// On live networks (Coston2), use the key from .env instead.
var PrvWithFunds *ecdsa.PrivateKey

func init() {
	var err error
	PrvWithFunds, err = crypto.HexToECDSA("804b01a8c27a65cc694a867be76edae3ccce7a7161cda1f67a8349df696d2207")
	if err != nil {
		panic("cannot read privateKey with funds")
	}
}
