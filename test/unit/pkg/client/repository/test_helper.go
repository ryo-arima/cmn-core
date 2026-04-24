package client

import (
	"github.com/ryo-arima/cmn-core/pkg/config"
	"gorm.io/gorm"
)

// TestHelper provides utilities for testing client components
type TestHelper struct {
	BaseConfig config.BaseConfig
	DB         *gorm.DB
}

// NewTestHelper creates a new test helper with test configuration
func NewTestHelper() *TestHelper {
	// Create test configuration (without real database connection for client tests)
	yamlConfig := config.YamlConfig{
		Application: config.Application{
			Client: config.Client{
				ServerEndpoint: "http://localhost:8080",
			},
			Server: config.Server{
				Admin: config.Admin{
					Emails: []string{"admin@test.com"},
				},
			},
		},
	}

	baseConfig := config.BaseConfig{
		DBConnection: nil, // Client tests don't need real DB connection
		YamlConfig:   yamlConfig,
	}

	return &TestHelper{
		BaseConfig: baseConfig,
		DB:         nil,
	}
}

// CleanupDB cleans up the test database (no-op for client tests)
func (th *TestHelper) CleanupDB() {
	// Client tests don't use real database, so nothing to cleanup
}
