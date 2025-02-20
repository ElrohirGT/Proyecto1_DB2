package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	APIPort string
	DatabaseConfig
	CorsConfig
}

type DatabaseConfig struct {
	DBUri      string
	DBUser     string
	DBPassword string
}

type CorsConfig struct {
	AllowedOrigins      []string
	AllowedContentTypes []string
	AllowedMethods      []string
	AllowedHeaders      []string
}

// LoadConfig loads configuration from environment variables.
// It exits with a fatal error if any required variable is missing.
func LoadConfig() Config {
	return Config{
		APIPort: mustGetEnv("API_PORT"),

		// Database
		DatabaseConfig: DatabaseConfig{
			DBUri:      mustGetEnv("DB_HOST"),
			DBUser:     mustGetEnv("DB_USER"),
			DBPassword: mustGetEnv("DB_USER_PASSWORD"),
		},

		CorsConfig: CorsConfig{
			AllowedOrigins:      mustGetEnvAsStringSlice("ALLOWED_ORIGINS"),
			AllowedContentTypes: mustGetEnvAsStringSlice("ALLOWED_CONTENT_TYPES"),
			AllowedMethods:      mustGetEnvAsStringSlice("ALLOWED_METHODS"),
			AllowedHeaders:      mustGetEnvAsStringSlice("ALLOWED_HEADERS"),
		},
	}
}

// mustGetEnv retrieves the value of the given environment variable
// or exits with a fatal error if the variable is not set.
func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

// mustGetEnvAsInt retrieves the value of the given environment variable,
// converts it to an integer, or exits with a fatal error if the variable
// is not set or cannot be converted to an integer.
func mustGetEnvAsInt(key string) int {
	valueStr := mustGetEnv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Environment variable %s must be a valid integer: %v", key, err)
	}
	return value
}

func mustGetEnvAsStringSlice(key string) []string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is required but not set", key)
	}

	return strings.Split(value, ",")
}
