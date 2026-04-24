package server

import (
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/global"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
)

func Main(conf config.BaseConfig) {
	if logger, ok := conf.Logger.(share.LoggerInterface); ok {
		logger.INFO("server-init", global.SSM1, "Starting cmn-core server on port 8000")
	}
	router := InitRouter(conf)
	if logger, ok := conf.Logger.(share.LoggerInterface); ok {
		logger.INFO("server-init", global.SSM3, "Server is ready")
	}
	router.Run(":8000")
}
