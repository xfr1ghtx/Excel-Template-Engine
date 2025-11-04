package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	ServerPort string
	ServerHost string
	BaseURL    string

	// MongoDB configuration
	MongoDBURI        string
	MongoDBDatabase   string
	MongoDBCollection string
	MongoDBTimeout    time.Duration

	// File paths
	TemplatePath  string
	GeneratedPath string

	// Logging
	LogLevel  string
	LogFormat string

	// File cleanup
	CleanupEnabled    bool
	CleanupInterval   time.Duration
	FileRetentionDays int
}

// Load loads configuration from environment variables
func Load() *Config {
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Printf("Could not load .env file: %v", err)
	}

	config := &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		ServerHost:        getEnv("SERVER_HOST", "0.0.0.0"),
		BaseURL:           getEnv("BASE_URL", "http://localhost:8080"),
		MongoDBURI:        getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBDatabase:   getEnv("MONGODB_DATABASE", "acts_db"),
		MongoDBCollection: getEnv("MONGODB_COLLECTION", "acts"),
		MongoDBTimeout:    parseDuration(getEnv("MONGODB_TIMEOUT", "10s"), 10*time.Second),
		TemplatePath:      getEnv("TEMPLATE_PATH", "./templates/act_template.xlsx"),
		GeneratedPath:     getEnv("GENERATED_PATH", "./generated"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		LogFormat:         getEnv("LOG_FORMAT", "json"),
		CleanupEnabled:    getEnv("CLEANUP_ENABLED", "false") == "true",
		CleanupInterval:   parseDuration(getEnv("CLEANUP_INTERVAL", "24h"), 24*time.Hour),
		FileRetentionDays: parseInt(getEnv("FILE_RETENTION_DAYS", "7"), 7),
	}

	log.Printf("Configuration loaded successfully")
	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// parseDuration parses a duration string or returns a default value
func parseDuration(value string, defaultValue time.Duration) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}

// parseInt parses an integer string or returns a default value
func parseInt(value string, defaultValue int) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
