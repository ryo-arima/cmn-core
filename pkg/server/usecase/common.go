package usecase

import (
	"context"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
)

type Common interface {
	GetBaseConfig() config.BaseConfig
	ValidateToken(ctx context.Context, tokenString string) (*model.JWTClaims, error)
	ResolveRole(email string) string
}

type common struct {
	commonRepo repository.Common
}

func NewCommon(commonRepo repository.Common) Common {
	return &common{commonRepo: commonRepo}
}

func (uc *common) GetBaseConfig() config.BaseConfig {
	return uc.commonRepo.GetBaseConfig()
}

func (uc *common) ValidateToken(ctx context.Context, tokenString string) (*model.JWTClaims, error) {
	return uc.commonRepo.ValidateToken(ctx, tokenString)
}

func (uc *common) ResolveRole(email string) string {
	return uc.commonRepo.ResolveRole(email)
}



