package testcase

import (
	"testing"

	"github.com/ryo-arima/cmn-core/test/e2e/server"
)

func TestMain(m *testing.M) {
	// Start test server
	if err := server.StartTestServer(); err != nil {
		panic("Failed to start test server: " + err.Error())
	}
	defer server.StopTestServer()

	// Initialize database schema
	if err := server.InitializeDatabase(); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	// Run tests
	m.Run()
}
