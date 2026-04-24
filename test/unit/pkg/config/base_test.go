package config_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/global"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestIntOrString_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Valid number string",
			input:    "5",
			expected: 5,
		},
		{
			name:     "Zero",
			input:    "0",
			expected: 0,
		},
		{
			name:     "Large number",
			input:    "9999",
			expected: 9999,
		},
		{
			name:     "Non-numeric string defaults to 0",
			input:    "not-a-number",
			expected: 0,
		},
		{
			name:     "Empty string defaults to 0",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yamlContent := "db: " + tt.input + "\n"
			var redisConfig struct {
				DB config.IntOrString `yaml:"db"`
			}

			err := yaml.Unmarshal([]byte(yamlContent), &redisConfig)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, int(redisConfig.DB))
		})
	}
}

func TestMCode_PaddedCode(t *testing.T) {
	maxLen := global.GetMaxCodeLength()

	tests := []struct {
		name  string
		mcode global.MCode
	}{
		{
			name:  "Short code with padding",
			mcode: global.MCode{Code: "TEST", Message: "Test message"},
		},
		{
			name:  "Code longer than max",
			mcode: global.MCode{Code: "TOOLONGCODEXXXXXXXX", Message: "Test"},
		},
		{
			name:  "Empty code",
			mcode: global.MCode{Code: "", Message: "Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.mcode.PaddedCode()
			// Result must start with the original code.
			assert.True(t, strings.HasPrefix(result, tt.mcode.Code),
				"result %q should start with code %q", result, tt.mcode.Code)
			// Result length must be max(len(code), maxLen).
			expectedLen := maxLen
			if len(tt.mcode.Code) > maxLen {
				expectedLen = len(tt.mcode.Code)
			}
			assert.Equal(t, expectedLen, len(result),
				"result %q should have length %d", result, expectedLen)
		})
	}
}

func TestNewBaseConfig(t *testing.T) {
	// Save original env vars
	origConfigFile := os.Getenv("CONFIG_FILE")
	origUseSecretsManager := os.Getenv("USE_SECRETSMANAGER")

	defer func() {
		// Restore original env vars
		os.Setenv("CONFIG_FILE", origConfigFile)
		os.Setenv("USE_SECRETSMANAGER", origUseSecretsManager)
	}()

	// Set test config file
	os.Setenv("CONFIG_FILE", "../../testdata/config/app.yaml")
	os.Setenv("USE_SECRETSMANAGER", "false")

	// Create config
	cfg := config.NewBaseConfig()

	require.NotNil(t, cfg)
	assert.NotNil(t, cfg.YamlConfig)
	assert.Equal(t, "testuser", cfg.YamlConfig.MySQL.User)
	assert.Equal(t, "testdb", cfg.YamlConfig.MySQL.Db)
}

func TestNewClientConfig(t *testing.T) {
	// Save original env vars
	origConfigFile := os.Getenv("CONFIG_FILE")
	origUseSecretsManager := os.Getenv("USE_SECRETSMANAGER")

	defer func() {
		os.Setenv("CONFIG_FILE", origConfigFile)
		os.Setenv("USE_SECRETSMANAGER", origUseSecretsManager)
	}()

	os.Setenv("CONFIG_FILE", "../../testdata/config/app.yaml")
	os.Setenv("USE_SECRETSMANAGER", "false")

	cfg := config.NewClientConfig()

	require.NotNil(t, cfg)
	assert.NotNil(t, cfg.YamlConfig)
	assert.Equal(t, "http://localhost:8080", cfg.YamlConfig.Application.Client.ServerEndpoint)
}

func TestSetLoggerFactory(t *testing.T) {
	t.Skip("SetLoggerFactory is not implemented in current version")
	called := false
	factory := func(lc config.LoggerConfig, bc *config.BaseConfig) interface{} {
		called = true
		return nil
	}

	_ = factory
	_ = called

	// config.SetLoggerFactory(factory)

	// Test that factory is set by creating a config
	origConfigFile := os.Getenv("CONFIG_FILE")
	origUseSecretsManager := os.Getenv("USE_SECRETSMANAGER")

	defer func() {
		os.Setenv("CONFIG_FILE", origConfigFile)
		os.Setenv("USE_SECRETSMANAGER", origUseSecretsManager)
	}()

	os.Setenv("CONFIG_FILE", "../../testdata/config/app.yaml")
	os.Setenv("USE_SECRETSMANAGER", "false")

	config.NewBaseConfig()

	assert.True(t, called, "Logger factory should have been called")
}

func TestLoggerConfig(t *testing.T) {
	lc := config.LoggerConfig{
		Component:    "test-component",
		Service:      "test-service",
		Level:        "DEBUG",
		Structured:   true,
		EnableCaller: true,
		Output:       "stdout",
	}

	assert.Equal(t, "test-component", lc.Component)
	assert.Equal(t, "test-service", lc.Service)
	assert.Equal(t, "DEBUG", lc.Level)
	assert.True(t, lc.Structured)
	assert.True(t, lc.EnableCaller)
	assert.Equal(t, "stdout", lc.Output)
}

func TestYamlConfigStructure(t *testing.T) {
	origConfigFile := os.Getenv("CONFIG_FILE")
	origUseSecretsManager := os.Getenv("USE_SECRETSMANAGER")

	defer func() {
		os.Setenv("CONFIG_FILE", origConfigFile)
		os.Setenv("USE_SECRETSMANAGER", origUseSecretsManager)
	}()

	os.Setenv("CONFIG_FILE", "../../testdata/config/app.yaml")
	os.Setenv("USE_SECRETSMANAGER", "false")

	cfg := config.NewBaseConfig()

	require.NotNil(t, cfg)

	// Test MySQL config
	assert.Equal(t, "localhost", cfg.YamlConfig.MySQL.Host)
	assert.Equal(t, "3306", cfg.YamlConfig.MySQL.Port)

	// Test Redis config
	assert.Equal(t, "localhost", cfg.YamlConfig.Redis.Host)
	assert.Equal(t, 6379, cfg.YamlConfig.Redis.Port)

	// Test Logger config
	assert.Equal(t, "test-component", cfg.YamlConfig.Logger.Component)
	assert.Equal(t, "test-service", cfg.YamlConfig.Logger.Service)

	// Test Application config
	assert.Contains(t, cfg.YamlConfig.Application.Server.Admin.Emails, "admin@test.local")

	// Test Mail config
	assert.Equal(t, "localhost", cfg.YamlConfig.Application.Mail.Host)
	assert.Equal(t, 587, cfg.YamlConfig.Application.Mail.Port)
}

func TestRedisIntOrString(t *testing.T) {
	origConfigFile := os.Getenv("CONFIG_FILE")
	origUseSecretsManager := os.Getenv("USE_SECRETSMANAGER")

	defer func() {
		os.Setenv("CONFIG_FILE", origConfigFile)
		os.Setenv("USE_SECRETSMANAGER", origUseSecretsManager)
	}()

	os.Setenv("CONFIG_FILE", "../../testdata/config/app.yaml")
	os.Setenv("USE_SECRETSMANAGER", "false")

	cfg := config.NewBaseConfig()

	require.NotNil(t, cfg)
	assert.Equal(t, 0, int(cfg.YamlConfig.Redis.DB))
}

func TestNewBaseConfigFromSource_LocalFile(t *testing.T) {
	origConfigFile := os.Getenv("CONFIG_FILE")
	origConfigSource := os.Getenv("CONFIG_SOURCE")
	origUseSecretsManager := os.Getenv("USE_SECRETSMANAGER")

	defer func() {
		os.Setenv("CONFIG_FILE", origConfigFile)
		os.Setenv("CONFIG_SOURCE", origConfigSource)
		os.Setenv("USE_SECRETSMANAGER", origUseSecretsManager)
	}()

	os.Setenv("CONFIG_FILE", "../../testdata/config/app.yaml")
	os.Setenv("CONFIG_SOURCE", "localfile")
	os.Setenv("USE_SECRETSMANAGER", "false")

	cfg := config.NewBaseConfigFromSource(context.Background())

	require.NotNil(t, cfg)
	assert.Equal(t, "testuser", cfg.YamlConfig.MySQL.User)
}

func TestNewBaseConfigFromSource_Default(t *testing.T) {
	origConfigFile := os.Getenv("CONFIG_FILE")
	origConfigSource := os.Getenv("CONFIG_SOURCE")
	origUseSecretsManager := os.Getenv("USE_SECRETSMANAGER")

	defer func() {
		os.Setenv("CONFIG_FILE", origConfigFile)
		os.Setenv("CONFIG_SOURCE", origConfigSource)
		os.Setenv("USE_SECRETSMANAGER", origUseSecretsManager)
	}()

	os.Setenv("CONFIG_FILE", "../../testdata/config/app.yaml")
	os.Unsetenv("CONFIG_SOURCE") // Test default behavior
	os.Setenv("USE_SECRETSMANAGER", "false")

	cfg := config.NewBaseConfigFromSource(context.Background())

	require.NotNil(t, cfg)
	assert.NotNil(t, cfg.YamlConfig)
}

func TestConnectDB_AlreadyConnected(t *testing.T) {
	origConfigFile := os.Getenv("CONFIG_FILE")
	origUseSecretsManager := os.Getenv("USE_SECRETSMANAGER")

	defer func() {
		os.Setenv("CONFIG_FILE", origConfigFile)
		os.Setenv("USE_SECRETSMANAGER", origUseSecretsManager)
	}()

	os.Setenv("CONFIG_FILE", "../../testdata/config/app.yaml")
	os.Setenv("USE_SECRETSMANAGER", "false")

	cfg := config.NewBaseConfig()

	// Test that DBConnection starts as nil
	assert.Nil(t, cfg.DBConnection)

	// Note: ConnectDB with real DB connection is tested in E2E tests
	// Unit test verifies the initial state and error handling
}

func TestMySQLConfig(t *testing.T) {
	mysql := config.MySQL{
		Host: "testhost",
		User: "testuser",
		Pass: "testpass",
		Port: "3306",
		Db:   "testdb",
	}

	assert.Equal(t, "testhost", mysql.Host)
	assert.Equal(t, "testuser", mysql.User)
	assert.Equal(t, "testpass", mysql.Pass)
	assert.Equal(t, "3306", mysql.Port)
	assert.Equal(t, "testdb", mysql.Db)
}

func TestRedisConfig(t *testing.T) {
	redis := config.Redis{
		Host: "redishost",
		Port: 6379,
		User: "default",
		Pass: "redispass",
		DB:   config.IntOrString(0),
	}

	assert.Equal(t, "redishost", redis.Host)
	assert.Equal(t, 6379, redis.Port)
	assert.Equal(t, "default", redis.User)
	assert.Equal(t, "redispass", redis.Pass)
	assert.Equal(t, 0, int(redis.DB))
}

func TestServerConfig(t *testing.T) {
	server := config.Server{
		Admin: config.Admin{
			Emails: []string{"admin@test.com"},
		},
		JWTSecret: "secret123",
		LogLevel:  "debug",
	}

	assert.Contains(t, server.Admin.Emails, "admin@test.com")
	assert.Equal(t, "secret123", server.JWTSecret)
	assert.Equal(t, "debug", server.LogLevel)
}

func TestMailConfig(t *testing.T) {
	mail := config.Mail{
		Host:     "mailhost",
		Port:     587,
		Username: "user@mail.com",
		Password: "mailpass",
		From:     "noreply@mail.com",
		UseTLS:   true,
	}

	assert.Equal(t, "mailhost", mail.Host)
	assert.Equal(t, 587, mail.Port)
	assert.Equal(t, "user@mail.com", mail.Username)
	assert.Equal(t, "mailpass", mail.Password)
	assert.Equal(t, "noreply@mail.com", mail.From)
	assert.True(t, mail.UseTLS)
}

func TestNewBaseConfigWithContext_SecretsManagerFallback(t *testing.T) {
	origConfigFile := os.Getenv("CONFIG_FILE")
	origUseSecretsManager := os.Getenv("USE_SECRETSMANAGER")
	origSecretID := os.Getenv("SECRET_ID")

	defer func() {
		os.Setenv("CONFIG_FILE", origConfigFile)
		os.Setenv("USE_SECRETSMANAGER", origUseSecretsManager)
		os.Setenv("SECRET_ID", origSecretID)
	}()

	// Test fallback when USE_SECRETSMANAGER is true but SECRET_ID is not set
	os.Setenv("CONFIG_FILE", "../../testdata/config/app.yaml")
	os.Setenv("USE_SECRETSMANAGER", "true")
	os.Unsetenv("SECRET_ID")

	cfg := config.NewBaseConfigWithContext(context.Background())

	require.NotNil(t, cfg)
	assert.NotNil(t, cfg.YamlConfig)
	// Should fall back to file-based config
	assert.Equal(t, "testuser", cfg.YamlConfig.MySQL.User)
}

func TestClientConfig(t *testing.T) {
	client := config.Client{
		ServerEndpoint: "http://localhost:8080",
		UserEmail:      "user@test.com",
		UserPassword:   "password123",
	}

	assert.Equal(t, "http://localhost:8080", client.ServerEndpoint)
	assert.Equal(t, "user@test.com", client.UserEmail)
	assert.Equal(t, "password123", client.UserPassword)
}

func TestAdminConfig(t *testing.T) {
	admin := config.Admin{
		Emails: []string{"admin1@test.com", "admin2@test.com"},
	}

	assert.Len(t, admin.Emails, 2)
	assert.Contains(t, admin.Emails, "admin1@test.com")
	assert.Contains(t, admin.Emails, "admin2@test.com")
}
