package main

import (
	"log"
	"os"

	"sign-extension/internal/app"
	"sign-extension/internal/base"
)

func main() {
	extPort := os.Getenv("EXTENSION_PORT")
	if extPort == "" {
		extPort = "8080"
	}

	signPort := os.Getenv("SIGN_PORT")
	if signPort == "" {
		signPort = "9090"
	}

	app.SetSignPort(signPort)
	srv := base.New(extPort, signPort, app.Version, app.Register, app.ReportState)

	log.Printf("extension listening on port %s", extPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
