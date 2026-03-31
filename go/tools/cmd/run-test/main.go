package main

import (
	"flag"
	"os"
	"time"

	"sign-tools/app"
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
	pf := flag.String("p", base.DefaultExtensionProxyURL, "extension proxy url")
	instructionSenderF := flag.String("instructionSender", defaultInstructionSender, "InstructionSender contract address")
	timeoutF := flag.Duration("timeout", 120*time.Second, "poll timeout for instruction results")
	flag.Parse()

	if *instructionSenderF == "" {
		logger.Fatal("--instructionSender flag is required (or set INSTRUCTION_SENDER in .env)")
	}

	instructionSenderAddress := common.HexToAddress(*instructionSenderF)

	testSupport, err := base.DefaultSupport(*af, *cf)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	err = app.RunE2ETest(testSupport, instructionSenderAddress, *pf, *timeoutF)
	if err != nil {
		fccutils.FatalWithCause(err)
	}
}
