// Package main provides a combined entry point for Docker that starts both the
// tee-node server and the sign extension in a single process. Unlike the
// two-binary approach, this uses teeServer.StartServerExtension() which does
// NOT override TestCodeHash, ensuring the TEE reports the correct code hash
// (194844...) that the Coston2 production verifier expects.
package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/flare-foundation/go-flare-common/pkg/logger"
	teeServer "github.com/flare-foundation/tee-node/pkg/server"

	"sign-extension/internal/app"
	"sign-extension/internal/base"
)

func main() {
	configPort := intEnv("CONFIG_PORT", 5502)
	signPort := intEnv("SIGN_PORT", 8882)
	extensionPort := intEnv("EXTENSION_PORT", 8883)

	// Start tee-node in extension mode (config server, sign server, forward router).
	go teeServer.StartServerExtension(configPort, signPort, extensionPort)

	// Start the sign extension server on the extension port.
	app.SetSignPort(strconv.Itoa(signPort))
	srv := base.New(strconv.Itoa(extensionPort), strconv.Itoa(signPort), app.Version, app.Register, app.ReportState)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Errorf("extension server error: %v", err)
		}
	}()

	logger.Infof("sign extension TEE running (config=%d, sign=%d, ext=%d)", configPort, signPort, extensionPort)

	// Wait for signal.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Info("shutting down")
}

func intEnv(key string, fallback int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return v
	}
	return fallback
}
