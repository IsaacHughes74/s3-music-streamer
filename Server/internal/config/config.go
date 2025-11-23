package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabasePath string
	S3Bucket     string
	S3Region     string
	ServerPort   string
	AWSAccessKey string
	AWSSecretKey string
}

func Load() *Config {
	// Try loading .env from current directory, then parent directories
	envPaths := []string{".env", "../.env", "../../.env"}
	loaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded environment from %s", path)
			loaded = true
			break
		}
	}
	if !loaded {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DatabasePath: getEnv("DATABASE_PATH", "./music.db"),
		S3Bucket:     getEnv("S3_BUCKET", ""),
		S3Region:     getEnv("S3_REGION", "us-east-1"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		AWSAccessKey: getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
