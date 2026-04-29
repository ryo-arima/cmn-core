package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/global"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"gorm.io/gorm"
)

// Local aliases for cleaner logging code - use functions to get logger dynamically
func INFO(requestID string, mcode global.MCode, message string) {
	if logger := share.GetServerLogger(); logger != nil {
		logger.INFO(requestID, mcode, message)
	}
}

func DEBUG(requestID string, mcode global.MCode, message string, fields ...map[string]interface{}) {
	if logger := share.GetServerLogger(); logger != nil {
		logger.DEBUG(requestID, mcode, message, fields...)
	}
}

func WARN(requestID string, mcode global.MCode, message string) {
	if logger := share.GetServerLogger(); logger != nil {
		logger.WARN(requestID, mcode, message)
	}
}

func ERROR(requestID string, mcode global.MCode, message string) {
	if logger := share.GetServerLogger(); logger != nil {
		logger.ERROR(requestID, mcode, message)
	}
}

// Local MCode definitions
var (
	SRNRSR1 = global.SRNRSR1
	SRNRSR2 = global.SRNRSR2
	Mcode   = global.Mcode
)

// RunInTx executes fn within a single database transaction.
// The transaction is committed when fn returns nil, and rolled back on error.
// Use this in the usecase layer to wrap multiple repository calls atomically.
func RunInTx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}

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
func (rcvr *common) ValidateToken(ctx context.Context, tokenString string) (*model.LoJWTClaims, error) {
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
	claims := &model.LoJWTClaims{
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

// UserQueryFilter holds optional filter conditions for local-DB user queries.
type UserQueryFilter struct {
	UUID  string
	Email string
}

// GroupQueryFilter holds optional filter conditions for local-DB group queries.
type GroupQueryFilter struct {
	UUID string
	Name string
}

// MemberQueryFilter holds optional filter conditions for local-DB member queries.
type MemberQueryFilter struct {
	GroupUUID string
	UserUUID  string
}
