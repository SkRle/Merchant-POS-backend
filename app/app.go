package app

import (
	"Merchant-POS-backend/config"
	"Merchant-POS-backend/pkg/database"

	"github.com/centrifugal/gocent/v3"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// App holds all the dependencies required for the application
type App struct {
	Config     *config.Config
	DB         *gorm.DB
	Redis      *redis.Client
	Minio      *minio.Client
	Centrifugo *gocent.Client
}

// NewApp initializes and returns a new App instance with all dependencies
func NewApp(cfg *config.Config) (*App, error) {
	// Initialize Postgres
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	// Initialize Redis
	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, err
	}

	// Initialize Minio
	// minioClient, err := database.NewMinioClient(cfg.Minio)
	// if err != nil {
	// 	return nil, err
	// }

	// Initialize Centrifugo
	centrifugoClient := database.NewCentrifugoClient(cfg.Centrifugo)

	app := &App{
		Config: cfg,
		DB:     db,
		Redis:  redisClient,
		// Minio:      minioClient,
		Centrifugo: centrifugoClient,
	}

	return app, nil
}
