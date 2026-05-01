package repository_test

import (
	"context"
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewCommonRepository(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)
	assert.NotNil(t, repo)
	assert.Equal(t, cfg, repo.GetBaseConfig())
}

func TestCommonRepository_ResolveRole(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					Admin: config.Admin{
						Emails: []string{"admin@cmn.local"},
					},
				},
			},
		},
	}
	repo := repository.NewCommon(cfg, nil)

	assert.Equal(t, "admin", repo.ResolveRole("admin@cmn.local"))
	assert.Equal(t, "user", repo.ResolveRole("other@example.com"))
}

func TestCommonRepository_ValidateToken_NilVerifier(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)

	_, err := repo.ValidateToken(context.Background(), "invalid.token.here")
	assert.Error(t, err)
}

func TestNewGroup(t *testing.T) {
	repo := repository.NewGroup(nil)
	assert.NotNil(t, repo)
}

func TestNewMember(t *testing.T) {
	t.Skip("MemberRepository requires keycloak/casdoor backend")
}

func TestNewRole(t *testing.T) {
	t.Skip("RoleRepository requires casbin enforcers")
}

func TestNewUser(t *testing.T) {
	t.Skip("UserRepository requires keycloak/casdoor backend")
}

