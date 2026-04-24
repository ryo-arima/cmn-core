package usecase

import (
	"github.com/ryo-arima/cmn-core/pkg/global"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
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
