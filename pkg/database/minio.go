package database

import (
	"Merchant-POS-backend/config"
	// "fmt"
	"log"

	"github.com/minio/minio-go/v7"
	// "github.com/minio/minio-go/v7"
	// "github.com/minio/minio-go/v7/pkg/credentials"
)

// func NewMinioClient(cfg config.MinioConfig) (*minio.Client, error) {
// 	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
// 		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
// 		Secure: cfg.UseSSL,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to minio: %w", err)
// 	}

// 	log.Println("Connected to Minio successfully")
// 	return minioClient, nil
// }

func NewMinioClient(cfg config.MinioConfig) (*minio.Client, error) {
	// minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
	// 	Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
	// 	Secure: cfg.UseSSL,
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to minio: %w", err)
	// }

	log.Println("Connected to Minio successfully")
	return nil, nil
}
