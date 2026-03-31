package main

import (
	"flag"
	"fmt"
	"os"

	"sign-tools/base"
	"sign-tools/base/fccutils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/flare-foundation/go-flare-common/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env for INSTRUCTION_SENDER
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")

	defaultInstructionSender := os.Getenv("INSTRUCTION_SENDER")

	af := flag.String("a", base.DefaultAddressesFile, "file with deployed addresses")
	cf := flag.String("c", base.DefaultChainNodeURL, "chain node url")
	instructionSenderF := flag.String("instructionSender", defaultInstructionSender, "InstructionSender contract address")
	governanceHashF := flag.String("governanceHash", "", "governance hash (optional)")
	flag.Parse()

	if *instructionSenderF == "" {
		logger.Fatal("--instructionSender flag is required (or set INSTRUCTION_SENDER in .env)")
	}

	testSupport, err := base.DefaultSupport(*af, *cf)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	governanceHash := common.HexToHash(*governanceHashF)
	instructionSenderAddress := common.HexToAddress(*instructionSenderF)

	logger.Infof("Registering extension with InstructionSender %s...", instructionSenderAddress.Hex())
	extensionID, err := fccutils.SetupExtension(testSupport, governanceHash, instructionSenderAddress, common.Address{})
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	extensionIDHex := fmt.Sprintf("0x%064x", extensionID)
	logger.Infof("Extension registered with ID: %s", extensionIDHex)

	// Machine-readable output on stdout.
	fmt.Println(extensionIDHex)
}
