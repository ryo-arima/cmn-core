package usecase_test

import (
	"context"
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewCommonUsecase(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)
	uc := usecase.NewCommon(repo)

	assert.NotNil(t, uc)
	assert.Equal(t, cfg, uc.GetBaseConfig())
}

func TestCommonUsecase_ResolveRole(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					Admin: config.Admin{
						Emails: []string{"admin@example.com"},
					},
				},
			},
		},
	}
	repo := repository.NewCommon(cfg, nil)
	uc := usecase.NewCommon(repo)

	assert.Equal(t, "admin", uc.ResolveRole("admin@example.com"))
	assert.Equal(t, "user", uc.ResolveRole("other@example.com"))
}

func TestCommonUsecase_ValidateToken_NilVerifier(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)
	uc := usecase.NewCommon(repo)

	_, err := uc.ValidateToken(context.Background(), "invalid.token.here")
	assert.Error(t, err)
}

func TestNewGroup(t *testing.T) {
	cfg := config.BaseConfig{}
	groupRepo := repository.NewGroup(cfg)
	memberRepo := repository.NewMember(cfg)

	uc := usecase.NewGroup(groupRepo, memberRepo, nil)
	assert.NotNil(t, uc)
}

func TestNewMember(t *testing.T) {
	cfg := config.BaseConfig{}
	memberRepo := repository.NewMember(cfg)

	uc := usecase.NewMember(memberRepo)
	assert.NotNil(t, uc)
}

func TestNewRole(t *testing.T) {
	t.Skip("Skipping role usecase test - requires casbin enforcer setup")
}

func TestNewUser(t *testing.T) {
	cfg := config.BaseConfig{}
	userRepo := repository.NewUser(cfg)

	uc := usecase.NewUser(userRepo)
	assert.NotNil(t, uc)
}

func TestCommonUsecase_Authorize(t *testing.T) {
	t.Skip("Skipping authorize test - requires casbin config files")
}

