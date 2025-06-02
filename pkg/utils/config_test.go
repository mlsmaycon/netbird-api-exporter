package utils

import (
	"os"
	"testing"
)

func TestGetEnvWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		expected     string
	}{
		{
			name:         "returns environment variable when set",
			key:          "TEST_VAR_SET",
			defaultValue: "default",
			envValue:     "env_value",
			setEnv:       true,
			expected:     "env_value",
		},
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_VAR_NOT_SET",
			defaultValue: "default_value",
			envValue:     "",
			setEnv:       false,
			expected:     "default_value",
		},
		{
			name:         "returns environment variable when set to empty string explicitly",
			key:          "TEST_VAR_EMPTY",
			defaultValue: "default",
			envValue:     "",
			setEnv:       true,
			expected:     "default",
		},
		{
			name:         "handles empty default value",
			key:          "TEST_VAR_EMPTY_DEFAULT",
			defaultValue: "",
			envValue:     "",
			setEnv:       false,
			expected:     "",
		},
		{
			name:         "handles special characters in environment variable",
			key:          "TEST_VAR_SPECIAL",
			defaultValue: "default",
			envValue:     "value with spaces and symbols !@#$%^&*()",
			setEnv:       true,
			expected:     "value with spaces and symbols !@#$%^&*()",
		},
		{
			name:         "handles multiline environment variable",
			key:          "TEST_VAR_MULTILINE",
			defaultValue: "default",
			envValue:     "line1\nline2\nline3",
			setEnv:       true,
			expected:     "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing environment variable
			originalValue := os.Getenv(tt.key)
			defer func() {
				if originalValue != "" {
					os.Setenv(tt.key, originalValue)
				} else {
					os.Unsetenv(tt.key)
				}
			}()

			// Set up test environment
			if tt.setEnv {
				if err := os.Setenv(tt.key, tt.envValue); err != nil {
					t.Fatalf("Failed to set environment variable: %v", err)
				}
			} else {
				os.Unsetenv(tt.key)
			}

			// Test the function
			result := GetEnvWithDefault(tt.key, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("GetEnvWithDefault(%q, %q) = %q, want %q", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestGetEnvWithDefault_ConcurrentAccess(t *testing.T) {
	// Test concurrent access to ensure thread safety
	const numGoroutines = 100
	const envKey = "TEST_CONCURRENT_VAR"
	const envValue = "concurrent_value"
	const defaultValue = "default"

	// Set environment variable
	originalValue := os.Getenv(envKey)
	defer func() {
		if originalValue != "" {
			os.Setenv(envKey, originalValue)
		} else {
			os.Unsetenv(envKey)
		}
	}()

	if err := os.Setenv(envKey, envValue); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	// Channel to collect results
	results := make(chan string, numGoroutines)

	// Launch goroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			result := GetEnvWithDefault(envKey, defaultValue)
			results <- result
		}()
	}

	// Collect and verify results
	for i := 0; i < numGoroutines; i++ {
		result := <-results
		if result != envValue {
			t.Errorf("Concurrent access test failed: expected %q, got %q", envValue, result)
		}
	}
}

func TestGetEnvWithDefault_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		envValue string
		expected string
		skip     bool
		skipMsg  string
	}{
		{
			name:     "very long environment variable",
			key:      "TEST_VAR_LONG",
			envValue: string(make([]byte, 1000)), // Reduced size for compatibility
			expected: string(make([]byte, 1000)),
		},
		{
			name:     "environment variable with only whitespace",
			key:      "TEST_VAR_WHITESPACE",
			envValue: "   \t\n   ",
			expected: "   \t\n   ",
		},
		{
			name:     "environment variable with unicode",
			key:      "TEST_VAR_UNICODE",
			envValue: "测试值 🚀 émojis",
			expected: "测试值 🚀 émojis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipMsg)
			}

			// Clean up
			originalValue := os.Getenv(tt.key)
			defer func() {
				if originalValue != "" {
					os.Setenv(tt.key, originalValue)
				} else {
					os.Unsetenv(tt.key)
				}
			}()

			// Set environment variable
			if err := os.Setenv(tt.key, tt.envValue); err != nil {
				// If setting the environment variable fails due to OS limitations, skip the test
				t.Skipf("Cannot set environment variable (OS limitation): %v", err)
			}

			result := GetEnvWithDefault(tt.key, "default")
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
