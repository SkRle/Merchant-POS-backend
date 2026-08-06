package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	Minio      MinioConfig      `yaml:"minio"`
	Centrifugo CentrifugoConfig `yaml:"centrifugo"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	UseSSL    bool   `yaml:"use_ssl"`
}

type CentrifugoConfig struct {
	Endpoint string `yaml:"endpoint"`
	APIKey   string `yaml:"api_key"`
}

func LoadConfig() (*Config, error) {
	// Try loading .env file, but don't fail if it doesn't exist (e.g., in production)
	_ = godotenv.Load()

	var cfg Config

	// Load from config.yml if it exists
	data, err := os.ReadFile("config.yml")
	if err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("error parsing config.yml: %w", err)
		}
	} else {
		// Fallback to default values if config.yml is missing
		cfg = Config{
			Server: ServerConfig{Port: 4500, Env: "local"},
		}
	}

	// Override with Environment Variables if present
	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Server.Port = port
		}
	}
	if env := os.Getenv("ENV"); env != "" {
		cfg.Server.Env = env
	}

	// Database
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}
	if dbPortStr := os.Getenv("DB_PORT"); dbPortStr != "" {
		if dbPort, err := strconv.Atoi(dbPortStr); err == nil {
			cfg.Database.Port = dbPort
		}
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.Database.User = dbUser
	}
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		cfg.Database.Password = dbPass
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.Database.Name = dbName
	}

	// Redis
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		cfg.Redis.Addr = redisAddr
	}
	if redisPass := os.Getenv("REDIS_PASSWORD"); redisPass != "" {
		cfg.Redis.Password = redisPass
	}

	// Minio
	if minioEndpoint := os.Getenv("MINIO_ENDPOINT"); minioEndpoint != "" {
		cfg.Minio.Endpoint = minioEndpoint
	}
	if minioAccess := os.Getenv("MINIO_ACCESS_KEY"); minioAccess != "" {
		cfg.Minio.AccessKey = minioAccess
	}
	if minioSecret := os.Getenv("MINIO_SECRET_KEY"); minioSecret != "" {
		cfg.Minio.SecretKey = minioSecret
	}
	if minioSSLStr := os.Getenv("MINIO_USE_SSL"); minioSSLStr != "" {
		cfg.Minio.UseSSL = minioSSLStr == "true"
	}

	// Centrifugo
	if centEndpoint := os.Getenv("CENTRIFUGO_ENDPOINT"); centEndpoint != "" {
		cfg.Centrifugo.Endpoint = centEndpoint
	}
	if centAPIKey := os.Getenv("CENTRIFUGO_API_KEY"); centAPIKey != "" {
		cfg.Centrifugo.APIKey = centAPIKey
	}

	return &cfg, nil
}
