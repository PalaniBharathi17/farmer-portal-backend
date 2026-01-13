package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

// Load loads configuration from environment variables and optional .env file.
func Load() Config {
	// Explicitly load .env file (not .env.example)
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("WARNING: .env file not found - using environment variables or defaults: %v", err)
		log.Println("INFO: Create a .env file based on .env.example template")
	} else {
		log.Println("INFO: Successfully loaded .env file")
	}

	cfg := Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "farmer_buyer"),
		JWTSecret:  getEnv("JWT_SECRET", "changeme"),
	}

	// Log confirmation of loaded DB config (never print password)
	passwordSet := "yes"
	if cfg.DBPassword == "" {
		passwordSet = "no"
		log.Println("WARNING: DB_PASSWORD is empty - ensure your MySQL user allows passwordless access or set DB_PASSWORD in .env")
	}
	log.Printf("INFO: Database configuration loaded - DB_HOST: %s, DB_PORT: %s, DB_USER: %s, DB_NAME: %s, Password set: %s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, passwordSet)

	return cfg
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
