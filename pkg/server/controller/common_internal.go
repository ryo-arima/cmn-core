package controller

import (
	"github.com/ryo-arima/cmn-core/pkg/server/share"
)

type CommonInternal interface {
}

type commonInternal struct {
	CommonRepository share.Common
}

func NewCommonInternal(commonRepository share.Common) CommonInternal {
	return &commonInternal{CommonRepository: commonRepository}
}
