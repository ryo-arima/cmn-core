package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/server"
)

var testServer *http.Server

// StartTestServer starts the test server in a goroutine
func StartTestServer() error {
	// Set test config file path relative to project root
	// Get current working directory and find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Navigate to project root (3 levels up from test/e2e/testcase)
	projectRoot := filepath.Join(cwd, "..", "..", "..")

	// Change to project root directory for relative paths (casbin, etc.)
	if err := os.Chdir(projectRoot); err != nil {
		return fmt.Errorf("failed to change to project root: %w", err)
	}

	testConfigPath := filepath.Join(projectRoot, "test", ".etc", "app.yaml")

	// Verify config file exists
	if _, err := os.Stat(testConfigPath); err != nil {
		return fmt.Errorf("config file not found at %s: %w", testConfigPath, err)
	}

	os.Setenv("CONFIG_FILE", testConfigPath)

	conf := config.NewBaseConfig()

	if err := conf.ConnectDB(); err != nil {
		return fmt.Errorf("failed to connect DB: %w", err)
	}

	router := server.InitRouter(*conf)

	testServer = &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		if err := testServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("test server failed: %v", err))
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Health check
	healthURL := "http://localhost:8000/health"
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(healthURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("test server did not start successfully")
}

// StopTestServer gracefully stops the test server
func StopTestServer() error {
	if testServer == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := testServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown test server: %w", err)
	}

	return nil
}

// InitializeDatabase creates all necessary tables for testing
func InitializeDatabase() error {
	conf := config.NewBaseConfig()

	if err := conf.ConnectDB(); err != nil {
		return fmt.Errorf("failed to connect DB: %w", err)
	}

	// Auto-migrate all tables
	if err := conf.DBConnection.AutoMigrate(
		&model.Users{},
		&model.Groups{},
		&model.Members{},
		&model.Resource{},
		&model.ResourceGroupRole{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
