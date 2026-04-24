package usecase

import (
	"context"
	"strings"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
)

type Common interface {
	GetBaseConfig() config.BaseConfig
	GenerateJWTToken(claims model.JWTClaims) (string, error)
	ValidateJWTToken(tokenString string) (*model.JWTClaims, error)
	ParseTokenUnverified(tokenString string) (*model.JWTClaims, error)
	IsTokenInvalidated(ctx context.Context, jti string) (bool, error)
	InvalidateToken(ctx context.Context, tokenString string) error
	GenerateTokenPair(userID uint, userUUID, email, name, role string) (*model.TokenPair, error)
	GenerateJWTSecret() (string, error)
	ValidateJWTSecretStrength(secret string) error
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
	ValidatePasswordStrength(password string) error
	DeleteTokenCache(token string)
	SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error
	SendWelcomeEmail(ctx context.Context, to, name string) error
	SendPasswordResetEmail(ctx context.Context, to, name, resetURL string) error
	// SSO state management
	StoreAuthState(ctx context.Context, state string) error
	ConsumeAuthState(ctx context.Context, state string) error
	StoreTokenForState(ctx context.Context, state string, tokenPair *model.TokenPair) error
	GetTokenForState(ctx context.Context, state string) (*model.TokenPair, bool, error)
	// Helper: determine role from admin email list
	ResolveRole(email string) string
}

type common struct {
	commonRepo repository.Common
}

func NewCommon(commonRepo repository.Common) Common {
	return &common{
		commonRepo: commonRepo,
	}
}

func (uc *common) GetBaseConfig() config.BaseConfig {
	return uc.commonRepo.GetBaseConfig()
}

func (uc *common) GenerateJWTToken(claims model.JWTClaims) (string, error) {
	return uc.commonRepo.GenerateJWTToken(claims)
}

func (uc *common) ValidateJWTToken(tokenString string) (*model.JWTClaims, error) {
	return uc.commonRepo.ValidateJWTToken(tokenString)
}

func (uc *common) ParseTokenUnverified(tokenString string) (*model.JWTClaims, error) {
	return uc.commonRepo.ParseTokenUnverified(tokenString)
}

func (uc *common) IsTokenInvalidated(ctx context.Context, jti string) (bool, error) {
	return uc.commonRepo.IsTokenInvalidated(ctx, jti)
}

func (uc *common) InvalidateToken(ctx context.Context, tokenString string) error {
	return uc.commonRepo.InvalidateToken(ctx, tokenString)
}

func (uc *common) GenerateTokenPair(userID uint, userUUID, email, name, role string) (*model.TokenPair, error) {
	return uc.commonRepo.GenerateTokenPair(userID, userUUID, email, name, role)
}

func (uc *common) GenerateJWTSecret() (string, error) {
	return uc.commonRepo.GenerateJWTSecret()
}

func (uc *common) ValidateJWTSecretStrength(secret string) error {
	return uc.commonRepo.ValidateJWTSecretStrength(secret)
}

func (uc *common) HashPassword(password string) (string, error) {
	return uc.commonRepo.HashPassword(password)
}

func (uc *common) VerifyPassword(hashedPassword, password string) error {
	return uc.commonRepo.VerifyPassword(hashedPassword, password)
}

func (uc *common) ValidatePasswordStrength(password string) error {
	return uc.commonRepo.ValidatePasswordStrength(password)
}

func (uc *common) DeleteTokenCache(token string) {
	uc.commonRepo.DeleteTokenCache(token)
}

func (uc *common) SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error {
	return uc.commonRepo.SendEmail(ctx, to, subject, body, isHTML)
}

func (uc *common) SendWelcomeEmail(ctx context.Context, to, name string) error {
	return uc.commonRepo.SendWelcomeEmail(ctx, to, name)
}

func (uc *common) SendPasswordResetEmail(ctx context.Context, to, name, resetURL string) error {
	return uc.commonRepo.SendPasswordResetEmail(ctx, to, name, resetURL)
}

func (uc *common) StoreAuthState(ctx context.Context, state string) error {
	return uc.commonRepo.StoreAuthState(ctx, state)
}

func (uc *common) ConsumeAuthState(ctx context.Context, state string) error {
	return uc.commonRepo.ConsumeAuthState(ctx, state)
}

func (uc *common) StoreTokenForState(ctx context.Context, state string, tokenPair *model.TokenPair) error {
	return uc.commonRepo.StoreTokenForState(ctx, state, tokenPair)
}

func (uc *common) GetTokenForState(ctx context.Context, state string) (*model.TokenPair, bool, error) {
	return uc.commonRepo.GetTokenForState(ctx, state)
}

// ResolveRole returns "admin" if the email is in the admin list, otherwise "app".
func (uc *common) ResolveRole(email string) string {
	cfg := uc.commonRepo.GetBaseConfig()
	for _, a := range cfg.YamlConfig.Application.Server.Admin.Emails {
		if strings.EqualFold(a, email) {
			return "admin"
		}
	}
	return "app"
}
