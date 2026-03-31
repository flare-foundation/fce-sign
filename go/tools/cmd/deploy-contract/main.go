package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"sign-tools/app"
	"sign-tools/base"
	"sign-tools/base/fccutils"

	"github.com/flare-foundation/go-flare-common/pkg/logger"
)

func main() {
	af := flag.String("a", base.DefaultAddressesFile, "file with deployed addresses")
	cf := flag.String("c", base.DefaultChainNodeURL, "chain node url")
	outFile := flag.String("o", "", "write deployed address to this file (optional)")
	verify := flag.Bool("verify", true, "verify contract on block explorer after deployment")
	explorerURL := flag.String("explorer-url", "https://coston2-explorer.flare.network/api", "block explorer API URL for verification")
	flag.Parse()

	testSupport, err := base.DefaultSupport(*af, *cf)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	logger.Infof("Deploying InstructionSender contract...")
	address, _, err := app.DeployInstructionSender(testSupport)
	if err != nil {
		fccutils.FatalWithCause(err)
	}

	logger.Infof("InstructionSender deployed at: %s", address.Hex())

	// Optionally write address to file for script consumption.
	if *outFile != "" {
		os.MkdirAll(filepath.Dir(*outFile), 0755)
		os.WriteFile(*outFile, []byte(address.Hex()), 0644)
	}

	// Verify contract on block explorer.
	if *verify {
		verifyContract(address.Hex(), testSupport.Addresses, *explorerURL)
	}

	// Machine-readable output on stdout.
	fmt.Println(address.Hex())
}

func verifyContract(address string, addresses *base.Addresses, explorerURL string) {
	// Check if forge and cast are available.
	if _, err := exec.LookPath("forge"); err != nil {
		logger.Warnf("forge not found, skipping contract verification (install Foundry to enable)")
		return
	}
	if _, err := exec.LookPath("cast"); err != nil {
		logger.Warnf("cast not found, skipping contract verification (install Foundry to enable)")
		return
	}

	// Encode constructor args.
	castArgs := exec.Command("cast", "abi-encode",
		"constructor(address,address)",
		addresses.TeeExtensionRegistry.Hex(),
		addresses.TeeMachineRegistry.Hex(),
	)
	constructorArgs, err := castArgs.Output()
	if err != nil {
		logger.Warnf("Failed to encode constructor args: %v", err)
		return
	}

	// Find the contract directory (relative to go/tools/).
	contractDir := "../../contract"
	if _, err := os.Stat(contractDir); err != nil {
		logger.Warnf("Contract directory not found at %s, skipping verification", contractDir)
		return
	}

	logger.Infof("Verifying contract on block explorer...")
	cmd := exec.Command("forge", "verify-contract",
		"--verifier", "etherscan",
		"--verifier-url", explorerURL,
		"--etherscan-api-key", "placeholder",
		"--compiler-version", "0.8.31",
		"--evm-version", "prague",
		"--constructor-args", string(constructorArgs[:len(constructorArgs)-1]), // trim newline
		"--root", contractDir,
		address,
		"InstructionSender.sol:InstructionSender",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Warnf("Contract verification failed: %v (contract is deployed but not verified on explorer)", err)
		return
	}
	logger.Infof("Contract verified on block explorer")
}
