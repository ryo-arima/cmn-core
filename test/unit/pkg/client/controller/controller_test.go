package controller_test

import (
	"strings"
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/client/auth"
	"github.com/ryo-arima/cmn-core/pkg/client/controller"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig() config.BaseConfig {
	return config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
					UserEmail:      "test@example.com",
					UserPassword:   "password123",
				},
			},
		},
	}
}

func setupManager(conf config.BaseConfig) *auth.Manager {
	return auth.NewManager(conf, "app")
}

func TestGetOutputFormat(t *testing.T) {
	format := controller.GetOutputFormat()
	assert.Contains(t, []string{"json", "table", "yaml"}, format)
}

func TestSetOutputFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"JSON format", "json"},
		{"Table format", "table"},
		{"YAML format", "yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller.SetOutputFormat(tt.format)
			result := controller.GetOutputFormat()
			assert.Equal(t, tt.format, result)
		})
	}
}

func TestCommonController_Logout(t *testing.T) {
	conf := setupTestConfig()
	manager := setupManager(conf)
	cmd := controller.InitCommonLogoutCmd(manager)

	assert.NotNil(t, cmd)
	assert.Equal(t, "logout", cmd.Use)
	assert.Contains(t, strings.ToLower(cmd.Short), "logout")
}

func TestCommonController_RefreshToken(t *testing.T) {
	conf := setupTestConfig()
	manager := setupManager(conf)
	cmd := controller.InitCommonRefreshTokenCmd(manager)

	assert.NotNil(t, cmd)
	assert.Equal(t, "refresh", cmd.Use)
}

func TestCommonController_ValidateToken(t *testing.T) {
	conf := setupTestConfig()
	manager := setupManager(conf)
	cmd := controller.InitCommonValidateTokenCmd(manager)

	assert.NotNil(t, cmd)
	assert.Equal(t, "validate", cmd.Use)
}

func TestCommonController_UserInfo(t *testing.T) {
	conf := setupTestConfig()
	manager := setupManager(conf)
	cmd := controller.InitCommonUserInfoCmd(manager)

	assert.NotNil(t, cmd)
	assert.Equal(t, "userinfo", cmd.Use)
}
