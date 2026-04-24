package config_test

import (
	"os"
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigFromEnv(t *testing.T) {
	tests := []struct {
		name             string
		secretID         string
		useLocalStack    string
		expectedSecretID string
		expectedUseLocal bool
	}{
		{
			name:             "Both env vars set to use LocalStack",
			secretID:         "test-secret-123",
			useLocalStack:    "true",
			expectedSecretID: "test-secret-123",
			expectedUseLocal: true,
		},
		{
			name:             "Use production AWS",
			secretID:         "prod-secret-456",
			useLocalStack:    "false",
			expectedSecretID: "prod-secret-456",
			expectedUseLocal: false,
		},
		{
			name:             "USE_LOCALSTACK not set",
			secretID:         "secret-789",
			useLocalStack:    "",
			expectedSecretID: "secret-789",
			expectedUseLocal: false,
		},
		{
			name:             "Empty SECRET_ID",
			secretID:         "",
			useLocalStack:    "true",
			expectedSecretID: "",
			expectedUseLocal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env vars
			origSecretID := os.Getenv("SECRET_ID")
			origUseLocalStack := os.Getenv("USE_LOCALSTACK")

			defer func() {
				// Restore original env vars
				os.Setenv("SECRET_ID", origSecretID)
				os.Setenv("USE_LOCALSTACK", origUseLocalStack)
			}()

			// Set test env vars
			os.Setenv("SECRET_ID", tt.secretID)
			os.Setenv("USE_LOCALSTACK", tt.useLocalStack)

			// Test GetConfigFromEnv
			secretID, useLocal := config.GetConfigFromEnv()

			assert.Equal(t, tt.expectedSecretID, secretID)
			assert.Equal(t, tt.expectedUseLocal, useLocal)
		})
	}
}

func TestGetConfigFromEnv_DefaultValues(t *testing.T) {
	// Save original env vars
	origSecretID := os.Getenv("SECRET_ID")
	origUseLocalStack := os.Getenv("USE_LOCALSTACK")

	defer func() {
		os.Setenv("SECRET_ID", origSecretID)
		os.Setenv("USE_LOCALSTACK", origUseLocalStack)
	}()

	// Unset env vars
	os.Unsetenv("SECRET_ID")
	os.Unsetenv("USE_LOCALSTACK")

	secretID, useLocal := config.GetConfigFromEnv()

	assert.Empty(t, secretID)
	assert.False(t, useLocal)
}
