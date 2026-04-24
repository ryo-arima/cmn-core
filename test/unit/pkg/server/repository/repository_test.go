package repository_test

import (
	"context"
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommonRepository(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			MySQL: config.MySQL{
				Host: "localhost",
				User: "test",
				Pass: "test",
				Port: "3306",
				Db:   "testdb",
			},
		},
	}

	repo := repository.NewCommon(cfg, nil)
	assert.NotNil(t, repo)
	assert.Equal(t, cfg, repo.GetBaseConfig())
}

func TestNewGroup(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			MySQL: config.MySQL{
				Host: "localhost",
			},
		},
	}

	repo := repository.NewGroup(cfg)
	assert.NotNil(t, repo)
}

func TestNewMember(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewMember(cfg)
	assert.NotNil(t, repo)
}

func TestNewRole(t *testing.T) {
	// RoleRepository requires casbin enforcers, skip basic initialization test
	// Test covered in E2E tests
	t.Skip("RoleRepository requires casbin enforcers")
}

func TestNewUser(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewUser(cfg)
	assert.NotNil(t, repo)
}

func TestHashPassword(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)

	password := "SecurePassword123!"
	hashed, err := repo.HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestVerifyPassword(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)

	password := "SecurePassword123!"
	hashed, err := repo.HashPassword(password)
	require.NoError(t, err)

	// Correct password
	err = repo.VerifyPassword(hashed, password)
	assert.NoError(t, err)

	// Incorrect password
	err = repo.VerifyPassword(hashed, "WrongPassword")
	assert.Error(t, err)
}

func TestValidatePasswordStrength(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid strong password",
			password: "StrongPass123!",
			wantErr:  false,
		},
		{
			name:     "Too short",
			password: "Short1!",
			wantErr:  true,
		},
		{
			name:     "No uppercase",
			password: "weakpass123!",
			wantErr:  true,
		},
		{
			name:     "No lowercase",
			password: "WEAKPASS123!",
			wantErr:  true,
		},
		{
			name:     "No digit",
			password: "WeakPassword!",
			wantErr:  true,
		},
		{
			name:     "No special char",
			password: "WeakPassword123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.ValidatePasswordStrength(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateJWTSecret(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)

	secret, err := repo.GenerateJWTSecret()
	require.NoError(t, err)
	assert.NotEmpty(t, secret)
	assert.GreaterOrEqual(t, len(secret), 32)
}

func TestValidateJWTSecretStrength(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)

	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{
			name:    "Valid secret",
			secret:  "this-is-a-very-long-and-secure-secret-key-for-jwt",
			wantErr: false,
		},
		{
			name:    "Too short",
			secret:  "short",
			wantErr: true,
		},
		{
			name:    "Empty",
			secret:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.ValidateJWTSecretStrength(tt.secret)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateTokenPair(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key-for-jwt-that-is-long-enough",
				},
			},
		},
	}
	repo := repository.NewCommon(cfg, nil)

	tokens, err := repo.GenerateTokenPair(1, "user-uuid-123", "test@example.com", "Test User", "user")

	require.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Greater(t, tokens.ExpiresIn, int64(0))
}

func TestValidateJWTToken(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key-for-jwt-that-is-long-enough",
				},
			},
		},
	}
	repo := repository.NewCommon(cfg, nil)

	// Generate a token
	tokens, err := repo.GenerateTokenPair(1, "user-uuid-123", "test@example.com", "Test User", "user")
	require.NoError(t, err)

	// Validate the token
	claims, err := repo.ValidateJWTToken(tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
}

func TestParseTokenUnverified(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key-for-jwt-that-is-long-enough",
				},
			},
		},
	}
	repo := repository.NewCommon(cfg, nil)

	// Generate a token
	tokens, err := repo.GenerateTokenPair(1, "user-uuid-123", "test@example.com", "Test User", "user")
	require.NoError(t, err)

	// Parse without verification
	claims, err := repo.ParseTokenUnverified(tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestIsTokenInvalidated(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommon(cfg, nil)

	// Without Redis, this should return false (token is valid)
	invalidated, err := repo.IsTokenInvalidated(context.Background(), "test-jti")
	assert.NoError(t, err)
	assert.False(t, invalidated)
}

func TestJWTClaims_Structure(t *testing.T) {
	claims := model.JWTClaims{
		UserID: 1,
		UUID:   "user-uuid-123",
		Email:  "test@example.com",
		Name:   "Test User",
		Role:   "admin",
		Jti:    "jti-123",
	}

	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "user-uuid-123", claims.UUID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, "jti-123", claims.Jti)
}
