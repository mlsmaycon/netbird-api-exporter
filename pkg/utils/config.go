package utils

import "os"

// GetEnvWithDefault returns environment variable value or default
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
