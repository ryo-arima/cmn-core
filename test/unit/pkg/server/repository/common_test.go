package repository

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
)

func TestCommonRepository_NewCommonRepository(t *testing.T) {
	baseConfig := CreateTestConfig()
	commonRepo := repository.NewCommon(baseConfig, nil)

	// Test that the repository is created successfully
	if commonRepo == nil {
		t.Error("NewCommonRepository() returned nil")
	}

	// Test that the base config is accessible
	config := commonRepo.GetBaseConfig()
	if config.YamlConfig.Application.Client.ServerEndpoint != "http://localhost:8080" {
		t.Errorf("Expected ServerEndpoint to be 'http://localhost:8080', got %s", config.YamlConfig.Application.Client.ServerEndpoint)
	}
}

// createTestConfig creates a minimal BaseConfig for testing
func createTestConfig() config.BaseConfig {
	return config.BaseConfig{
		DBConnection: nil, // No database connection needed for these tests
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Common: config.Common{},
				Server: config.Server{
					Admin: config.Admin{
						Emails: []string{"admin@test.com"},
					},
				},
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
					Credentials: config.ClientCredentials{
						Email:    "test@test.com",
						Password: "testpass",
					},
				},
			},
			PostgreSQL: config.PostgreSQL{
				Host: "127.0.0.1",
				User: "root",
				Pass: "postgres",
				Port: "5432",
				Db:   "cmn_core_test",
			},
		},
	}
}

func TestCommonRepository_GetBaseConfig(t *testing.T) {
	// Setup
	baseConfig := CreateTestConfig()
	commonRepo := repository.NewCommon(baseConfig, nil)

	// Test
	config := commonRepo.GetBaseConfig()

	// Assert using go-cmp for better comparison
	expectedEndpoint := "http://localhost:8080"
	if diff := cmp.Diff(expectedEndpoint, config.YamlConfig.Application.Client.ServerEndpoint); diff != "" {
		t.Errorf("ServerEndpoint mismatch (-want +got):\n%s", diff)
	}

	if len(config.YamlConfig.Application.Server.Admin.Emails) == 0 {
		t.Error("Expected admin emails to be set, got empty slice")
	} else {
		expectedEmail := "admin@test.com"
		if diff := cmp.Diff(expectedEmail, config.YamlConfig.Application.Server.Admin.Emails[0]); diff != "" {
			t.Errorf("Admin email mismatch (-want +got):\n%s", diff)
		}
	}

	// Verify that the config is not nil and has expected structure
	if config.DBConnection != baseConfig.DBConnection {
		t.Error("Expected DBConnection to match the original config")
	}
}

func TestNewCommonRepository(t *testing.T) {
	// Setup
	baseConfig := CreateTestConfig()

	// Test
	commonRepo := repository.NewCommon(baseConfig, nil)

	// Assert using go-cmp
	if diff := cmp.Diff(false, commonRepo == nil); diff != "" {
		t.Errorf("Repository creation check mismatch (-want +got):\n%s", diff)
	}

	// Test that we can call methods on the repository
	config := commonRepo.GetBaseConfig()
	expectedEndpoint := "http://localhost:8080"
	if diff := cmp.Diff(expectedEndpoint, config.YamlConfig.Application.Client.ServerEndpoint); diff != "" {
		t.Errorf("ServerEndpoint configuration check failed (-want +got):\n%s", diff)
	}
}

// TestCommonRepository_GetBaseConfig_Multiple tests multiple calls to GetBaseConfig
func TestCommonRepository_GetBaseConfig_Multiple(t *testing.T) {
	// Setup
	baseConfig := CreateTestConfig()
	commonRepo := repository.NewCommon(baseConfig, nil)

	// Test - call GetBaseConfig multiple times
	config1 := commonRepo.GetBaseConfig()
	config2 := commonRepo.GetBaseConfig()

	// Assert - both should return valid configurations
	expectedEndpoint := "http://localhost:8080"
	if diff := cmp.Diff(expectedEndpoint, config1.YamlConfig.Application.Client.ServerEndpoint); diff != "" {
		t.Errorf("First config ServerEndpoint mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(expectedEndpoint, config2.YamlConfig.Application.Client.ServerEndpoint); diff != "" {
		t.Errorf("Second config ServerEndpoint mismatch (-want +got):\n%s", diff)
	}

	// Verify consistency using go-cmp
	if diff := cmp.Diff(config1.YamlConfig.Application.Client.ServerEndpoint,
		config2.YamlConfig.Application.Client.ServerEndpoint); diff != "" {
		t.Errorf("Server endpoint mismatch between calls (-first +second):\n%s", diff)
	}
}

// TestCommonRepository_GetBaseConfig_Validation tests configuration validation
func TestCommonRepository_GetBaseConfig_Validation(t *testing.T) {
	// Setup
	baseConfig := CreateTestConfig()
	commonRepo := repository.NewCommon(baseConfig, nil)

	// Test
	config := commonRepo.GetBaseConfig()

	// Assert - validate all expected configuration values
	expectedEndpoint := "http://localhost:8080"
	if diff := cmp.Diff(expectedEndpoint, config.YamlConfig.Application.Client.ServerEndpoint); diff != "" {
		t.Errorf("Server endpoint mismatch (-want +got):\n%s", diff)
	}

	if len(config.YamlConfig.Application.Server.Admin.Emails) == 0 {
		t.Error("Expected admin emails to be set")
	} else {
		expectedEmail := "admin@test.com"
		if diff := cmp.Diff(expectedEmail, config.YamlConfig.Application.Server.Admin.Emails[0]); diff != "" {
			t.Errorf("Admin email mismatch (-want +got):\n%s", diff)
		}
	}

	// Test PostgreSQL configuration
	expectedHost := "127.0.0.1"
	if diff := cmp.Diff(expectedHost, config.YamlConfig.PostgreSQL.Host); diff != "" {
		t.Errorf("PostgreSQL host mismatch (-want +got):\n%s", diff)
	}

	expectedDB := "cmn_core_test"
	if diff := cmp.Diff(expectedDB, config.YamlConfig.PostgreSQL.Db); diff != "" {
		t.Errorf("PostgreSQL database name mismatch (-want +got):\n%s", diff)
	}
}

// TestCommonRepository_TableDriven demonstrates table-driven tests for all operations
func TestCommonRepository_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		operation     string
		shouldPass    bool
		expectedValue interface{}
	}{
		{
			name:          "get base config - server endpoint",
			operation:     "get_config_endpoint",
			shouldPass:    true,
			expectedValue: "http://localhost:8080",
		},
		{
			name:          "get base config - admin email",
			operation:     "get_config_email",
			shouldPass:    true,
			expectedValue: "admin@test.com",
		},
		{
			name:          "get base config - user email",
			operation:     "get_config_user",
			shouldPass:    true,
			expectedValue: "test@test.com",
		},
		{
			name:          "get base config - psql host",
			operation:     "get_config_psql_host",
			shouldPass:    true,
			expectedValue: "127.0.0.1",
		},
		{
			name:          "get base config - psql db",
			operation:     "get_config_psql_db",
			shouldPass:    true,
			expectedValue: "cmn_core_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			baseConfig := CreateTestConfig()
			commonRepo := repository.NewCommon(baseConfig, nil)

			// Test based on operation
			var actualValue interface{}
			var success bool

			config := commonRepo.GetBaseConfig()

			switch tt.operation {
			case "get_config_endpoint":
				actualValue = config.YamlConfig.Application.Client.ServerEndpoint
				success = actualValue == tt.expectedValue.(string)
			case "get_config_email":
				if len(config.YamlConfig.Application.Server.Admin.Emails) > 0 {
					actualValue = config.YamlConfig.Application.Server.Admin.Emails[0]
					success = actualValue == tt.expectedValue.(string)
				}
			case "get_config_user":
				actualValue = config.YamlConfig.Application.Client.Credentials.Email
				success = actualValue == tt.expectedValue.(string)
			case "get_config_psql_host":
				actualValue = config.YamlConfig.PostgreSQL.Host
				success = actualValue == tt.expectedValue.(string)
			case "get_config_psql_db":
				actualValue = config.YamlConfig.PostgreSQL.Db
				success = actualValue == tt.expectedValue.(string)
			}

			// Assert using go-cmp
			if diff := cmp.Diff(tt.shouldPass, success); diff != "" {
				t.Errorf("Operation success mismatch for %s (-want +got):\n%s", tt.operation, diff)
			}

			if success && actualValue != nil {
				if diff := cmp.Diff(tt.expectedValue, actualValue); diff != "" {
					t.Errorf("Value mismatch for %s (-want +got):\n%s", tt.operation, diff)
				}
			}
		})
	}
}

// TestCommonRepository_EdgeCases tests edge cases and error scenarios
func TestCommonRepository_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		baseConfig config.BaseConfig
		expectNil  bool
	}{
		{
			name: "minimal config",
			baseConfig: config.BaseConfig{
				YamlConfig: config.YamlConfig{},
			},
			expectNil: false,
		},
		{
			name: "empty admin emails",
			baseConfig: config.BaseConfig{
				YamlConfig: config.YamlConfig{
					Application: config.Application{
						Server: config.Server{
							Admin: config.Admin{
								Emails: []string{},
							},
						},
					},
				},
			},
			expectNil: false,
		},
		{
			name: "nil yaml config",
			baseConfig: config.BaseConfig{
				DBConnection: nil,
			},
			expectNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			commonRepo := repository.NewCommon(tt.baseConfig, nil)

			// Assert
			isNil := commonRepo == nil
			if diff := cmp.Diff(tt.expectNil, isNil); diff != "" {
				t.Errorf("Repository nil check mismatch (-want +got):\n%s", diff)
			}

			if !tt.expectNil && commonRepo != nil {
				// Test that basic operations don't panic
				config := commonRepo.GetBaseConfig()
				_ = config // Just verify it doesn't panic

				// Verify interface compliance
				var _ repository.Common = commonRepo
			}
		})
	}
}

// TestCommonRepository_InterfaceCompliance verifies interface implementation
func TestCommonRepository_InterfaceCompliance(t *testing.T) {
	// Setup
	baseConfig := createTestConfig()

	// Test
	commonRepo := repository.NewCommon(baseConfig, nil)

	// Assert - verify interface compliance at compile time
	var _ repository.Common = commonRepo

	// Runtime verification that all methods are available
	config := commonRepo.GetBaseConfig()
	if diff := cmp.Diff(false, config.YamlConfig.Application.Client.ServerEndpoint == ""); diff != "" {
		t.Errorf("Interface method GetBaseConfig() failed (-want +got):\n%s", diff)
	}
}
