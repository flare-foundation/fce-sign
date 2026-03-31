package main

import (
	"encoding/hex"
	"flag"
	"os"

	"sign-tools/base"
	"sign-tools/base/fccutils"

	"github.com/flare-foundation/go-flare-common/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env for TUNNEL_URL
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")

	defaultHostURL := os.Getenv("TUNNEL_URL")

	af := flag.String("a", base.DefaultAddressesFile, "file with deployed addresses")
	cf := flag.String("c", base.DefaultChainNodeURL, "chain node url")
	pf := flag.String("p", base.DefaultExtensionProxyURL, "extension proxy url (local)")
	hostF := flag.String("host", defaultHostURL, "on-chain host URL reachable by external TEEs (e.g. tunnel URL)")
	defaultExternalProxy := os.Getenv("NORMAL_PROXY_URL")
	if defaultExternalProxy == "" {
		defaultExternalProxy = "https://tee-proxy-coston2-1.flare.rocks"
	}
	epf := flag.String("ep", defaultExternalProxy, "external proxy url for FTDC availability check (must be a production TEE on extension 0)")
	lf := flag.Bool("l", false, "local")
	instructionF := flag.String("i", "", "instructionID")
	command := flag.String("command", "rap", "command (rap)")

	flag.Parse()

	hostURL := *hostF
	if hostURL == "" {
		hostURL = *pf
	}

	testSupport, err := base.DefaultSupport(*af, *cf)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	// Get teeID from proxy.
	teeInfo, err := fccutils.TeeInfo(*pf)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	teeID, _, err := fccutils.TeeProxyId(teeInfo)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	ftdcTeeID, _, err := fccutils.GetTeeProxyID(*epf)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	// Verify code hash / platform.
	_, _, err = fccutils.GetCodeHashAndPlatform(teeInfo, *lf)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	logger.Infof("Registration of TEE with ID %s", hex.EncodeToString(teeID[:]))
	err = fccutils.RegisterNode(testSupport, teeInfo, hostURL, *epf, ftdcTeeID, *command, *instructionF)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	logger.Infof("Registered TEE node with id %s", teeID)
}
