package database

import (
	"Merchant-POS-backend/config"
	"log"

	"github.com/centrifugal/gocent/v3"
)

func NewCentrifugoClient(cfg config.CentrifugoConfig) *gocent.Client {
	client := gocent.New(gocent.Config{
		Addr: cfg.Endpoint,
		Key:  cfg.APIKey,
	})

	log.Println("Centrifugo client initialized")
	return client
}
