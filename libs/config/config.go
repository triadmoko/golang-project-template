package config

import (
	"os"
)

// Config holds all configuration for our application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// Load loads configuration from environment variables
func Load() Config {
	config := Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:    getEnv("DB_HOST", "localhost"),
			Port:    getEnv("DB_PORT", "5432"),
			User:    getEnv("DB_USER", "postgres"),
			Pass:    getEnv("DB_PASS", "password"),
			Name:    getEnv("DB_NAME", "app"),
			SSLMode: getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
		},
	}

	return config
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
