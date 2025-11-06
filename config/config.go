package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func Load() Config {
	// Only auto-load if DATABASE_URL is not already set (might be set via --env-path flag)
	if os.Getenv("DATABASE_URL") == "" {
		// Load from config/.env first, then fallback to .env in root
		if err := godotenv.Load("config/.env"); err != nil {
			log.Println("config/.env not found, trying .env in root")
			_ = godotenv.Load(".env")
		} else {
			log.Println("Loaded config from config/.env")
		}
	}

	cfg := Config{
		DatabaseURL: getenv("DATABASE_URL", ""),
		JWTSecret:   getenv("JWT_SECRET", "your-secret-key-change-in-production"),
		Port:        getenv("PORT", "3011"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required in config/.env file")
	}
	// PORT defaults to 3011 if not set

	return cfg
}

func getenv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
