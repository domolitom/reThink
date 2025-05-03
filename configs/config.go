package configs

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port        string
	Mode        string
	CORSEnabled bool
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret string
	TTL    int // time to live in hours
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8080"),
			Mode:        getEnv("GIN_MODE", "debug"),
			CORSEnabled: getEnvAsBool("CORS_ENABLED", true),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "prediction_social"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-default-secret-key"),
			TTL:    getEnvAsInt("JWT_TTL", 168), // Default: 7 days (168 hours)
		},
	}

	return config
}

// Helper function to get an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as a boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

// Helper function to get an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
