package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"depop-alerter/alerter"
	"depop-alerter/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	app, err := alerter.NewAlerter(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize alerter: %v", err)
	}

	app.Start()

	// Wait for interrupt signal to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received termination signal. Shutting down...")
	app.Stop()
}
