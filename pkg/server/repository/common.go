package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
)

// Common interface for repository layer
type Common interface {
	share.Common // embed ValidateToken for middleware compatibility
	GetBaseConfig() config.BaseConfig
	ResolveRole(email string) string
}

// common implements Common interface
type common struct {
	BaseConfig   config.BaseConfig
	oidcVerifier *gooidc.IDTokenVerifier
	adminEmails  map[string]struct{}
}

func NewCommon(conf config.BaseConfig, verifier *gooidc.IDTokenVerifier) Common {
	adminEmails := make(map[string]struct{})
	for _, e := range conf.YamlConfig.Application.Server.Admin.Emails {
		adminEmails[e] = struct{}{}
	}
	if verifier == nil {
		log.Println("repository.NewCommon: OIDC verifier is nil – ValidateToken will always return an error")
	}
	return &common{
		BaseConfig:   conf,
		oidcVerifier: verifier,
		adminEmails:  adminEmails,
	}
}

func (rcvr *common) GetBaseConfig() config.BaseConfig {
	return rcvr.BaseConfig
}

func (rcvr *common) ResolveRole(email string) string {
	if _, ok := rcvr.adminEmails[email]; ok {
		return "admin"
	}
	return "user"
}

// ValidateToken validates an IdP-issued JWT token via the OIDC JWKS endpoint.
func (rcvr *common) ValidateToken(ctx context.Context, tokenString string) (*model.JWTClaims, error) {
	if rcvr.oidcVerifier == nil {
		return nil, errors.New("OIDC verifier not configured: cannot validate JWT")
	}
	idToken, err := rcvr.oidcVerifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}
	var raw struct {
		Email  string   `json:"email"`
		Name   string   `json:"name"`
		Groups []string `json:"groups"`
	}
	if err := idToken.Claims(&raw); err != nil {
		return nil, fmt.Errorf("failed to extract token claims: %w", err)
	}
	claims := &model.JWTClaims{
		UUID:      idToken.Subject,
		Email:     raw.Email,
		Name:      raw.Name,
		Groups:    raw.Groups,
		Role:      rcvr.ResolveRole(raw.Email),
		IssuedAt:  idToken.IssuedAt.Unix(),
		ExpiresAt: idToken.Expiry.Unix(),
	}
	return claims, nil
}





