package main

import (
	"crypto/ecdsa"
	"flag"
	"os"

	"sign-tools/base"
	"sign-tools/base/fccutils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/flare-foundation/go-flare-common/pkg/logger"
)

func main() {
	af := flag.String("a", base.DefaultAddressesFile, "file with deployed addresses")
	cf := flag.String("c", base.DefaultChainNodeURL, "chain node url")
	pf := flag.String("p", base.DefaultExtensionProxyURL, "proxy url")
	versionF := flag.String("version", "v0.1.0", "version")
	lf := flag.Bool("l", false, "local/test mode: use hardcoded test code hash instead of the image hash from /info")
	flag.Parse()

	testSupport, err := base.DefaultSupport(*af, *cf)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	// Get teeID from proxy.
	teeInfo, err := fccutils.TeeInfo(*pf)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	var privKey *ecdsa.PrivateKey
	privKeyString := os.Getenv("EXTENSION_OWNER_KEY")
	if privKeyString != "" {
		privKey, err = crypto.HexToECDSA(privKeyString)
		if err != nil {
			fccutils.FatalWithCause(err)
		}
	} else {
		privKey = testSupport.Prv
	}

	teeID, _, err := fccutils.TeeProxyId(teeInfo)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	codeHash := teeInfo.MachineData.CodeHash
	platform := teeInfo.MachineData.Platform
	extensionID := teeInfo.MachineData.ExtensionID.Big()

	// Always register the image hash reported by /info (needed for the
	// signed Register/pre-registration call).
	logger.Infof("Registering version: %s, %s, extension: %v, tee id: %s",
		codeHash, platform, extensionID, teeID)

	err = fccutils.AddTeeVersion(testSupport, privKey,
		extensionID, codeHash, platform,
		common.Hash{}, *versionF)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	if *lf && (codeHash != fccutils.TeeCodeHash || platform != fccutils.TestPlatform) {
		// In test mode, the tee-node attestation proof uses a hardcoded
		// code hash and TEST_PLATFORM regardless of the actual Docker
		// image. Also register these so toProduction proofs match.
		logger.Infof("Also registering hardcoded test version: %s, %s",
			fccutils.TeeCodeHash, fccutils.TestPlatform)

		err = fccutils.AddTeeVersion(testSupport, privKey,
			extensionID, fccutils.TeeCodeHash, fccutils.TestPlatform,
			common.Hash{}, *versionF)
		if err != nil {
			fccutils.FatalWithCause(err)
		}
	}
}
