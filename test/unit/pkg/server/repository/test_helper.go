package repository

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestHelper provides utilities for testing
type TestHelper struct {
	DB         *gorm.DB
	BaseConfig config.BaseConfig
	MockDB     sqlmock.Sqlmock
	SqlDB      *sql.DB
}

// CreateTestConfig creates a minimal BaseConfig for testing
// This replaces the individual createXxxTestConfig functions in each test file
func CreateTestConfig() config.BaseConfig {
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
						Email: "test@test.com", Password: "testpass",
					},
				},
			},
			PostgreSQL: config.PostgreSQL{
				Host: "127.0.0.1",
				User: "user",
				Pass: "password",
				Port: "5432",
				Db:   "cmn_core_test",
			},
		},
	}
}

// NewTestHelper creates a new test helper with mock database
func NewTestHelper() *TestHelper {
	// Create a mock database connection
	db, mockDB, sqlDB := setupTestDB()

	baseConfig := CreateTestConfig()
	baseConfig.DBConnection = db

	return &TestHelper{
		DB:         db,
		BaseConfig: baseConfig,
		MockDB:     mockDB,
		SqlDB:      sqlDB,
	}
}

// setupTestDB creates a test database connection with mock
func setupTestDB() (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	// Create a mock database for testing
	return createMockDB()
}

// createMockDB creates a mock database connection for testing using sqlmock
func createMockDB() (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	// Create a mock SQL database
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create mock database: %v", err))
	}

	// Create GORM DB with the mock (postgres driver, no init queries with PreferSimpleProtocol)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create GORM database: %v", err))
	}

	return gormDB, mock, sqlDB
}

// CleanupDB cleans up the test database
func (th *TestHelper) CleanupDB() {
	if th.SqlDB != nil {
		th.SqlDB.Close()
	}
}
