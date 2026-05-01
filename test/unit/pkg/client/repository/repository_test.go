package client

import (
	"testing"

	clientauth "github.com/ryo-arima/cmn-core/pkg/client/share"
	"github.com/ryo-arima/cmn-core/pkg/client/repository"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/stretchr/testify/assert"
)

func newTestCfg() config.BaseConfig {
	return config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
					Credentials: config.ClientCredentials{
						Email:    "test@example.com",
						Password: "password",
					},
				},
			},
		},
	}
}

// TestNewCommonRepository verifies that a Common repository is constructed without error.
func TestNewCommonRepository(t *testing.T) {
	cfg := newTestCfg()
	manager := clientauth.NewManager(cfg, "app")
	repo := repository.NewCommon(cfg, manager)
	assert.NotNil(t, repo)
}

// TestRepositoryInterface verifies that NewCommon satisfies the Common interface.
func TestRepositoryInterface(t *testing.T) {
	cfg := newTestCfg()
	manager := clientauth.NewManager(cfg, "app")
	var _ repository.Common = repository.NewCommon(cfg, manager)
}



