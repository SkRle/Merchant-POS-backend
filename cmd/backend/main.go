package main

import (
	"Merchant-POS-backend/app"
	"Merchant-POS-backend/config"
	"log"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Application (DB, Redis, Minio, Centrifugo, etc.)
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application dependencies: %v", err)
	}

	// 3. Start Fiber Server
	if err := application.StartServer(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
