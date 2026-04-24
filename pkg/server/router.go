package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/auth"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/server/controller"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

func InitRouter(conf config.BaseConfig) *gin.Engine {
	// Initialize global server logger
	if logger, ok := conf.Logger.(*share.ServerLogger); ok {
		share.SetServerLogger(logger)
	}

	if logger, ok := conf.Logger.(share.LoggerInterface); ok {
		gin.DefaultWriter = share.NewGinLoggerWriter(logger)
		gin.DefaultErrorWriter = share.NewGinLoggerWriter(logger)
	}

	redisClient, err := repository.NewRedisClient(conf.YamlConfig.Redis)
	if err != nil {
		panic(err)
	}

	commonRepository := repository.NewCommon(conf, redisClient)
	commonUsecase := usecase.NewCommon(commonRepository)

	// Initialize OIDC provider (optional — skip if not configured)
	var oidcProvider auth.Provider
	oidcCfg := conf.YamlConfig.Application.Server.Auth.OIDC
	if oidcCfg.IssuerURL != "" && oidcCfg.ClientID != "" {
		p, err := auth.NewOIDCProvider(context.Background(), auth.OIDCConfig{
			ProviderName: oidcCfg.ProviderName,
			IssuerURL:    oidcCfg.IssuerURL,
			ClientID:     oidcCfg.ClientID,
			ClientSecret: oidcCfg.ClientSecret,
			RedirectURL:  oidcCfg.RedirectURL,
			Scopes:       oidcCfg.Scopes,
		})
		if err != nil {
			log.Printf("OIDC provider init failed (OIDC disabled): %v", err)
		} else {
			oidcProvider = p
		}
	}

	// Initialize SAML provider (optional — skip if not configured)
	var samlProvider auth.Provider
	samlCfg := conf.YamlConfig.Application.Server.Auth.SAML
	if samlCfg.SPACSURL != "" && samlCfg.SPEntityID != "" && (samlCfg.IDPMetadataURL != "" || samlCfg.IDPCertificatePEM != "") {
		p, err := auth.NewSAMLProvider(context.Background(), auth.SAMLConfig{
			ProviderName:      samlCfg.ProviderName,
			IDPMetadataURL:    samlCfg.IDPMetadataURL,
			IDPCertificatePEM: samlCfg.IDPCertificatePEM,
			SPEntityID:        samlCfg.SPEntityID,
			SPACSURL:          samlCfg.SPACSURL,
			SPKeyPEM:          samlCfg.SPKeyPEM,
			SPCertPEM:         samlCfg.SPCertPEM,
		})
		if err != nil {
			log.Printf("SAML provider init failed (SAML disabled): %v", err)
		} else {
			samlProvider = p
		}
	}

	commonShareController := controller.NewCommonShare(commonUsecase, oidcProvider, samlProvider)

	router := gin.Default()

	// Health check (no authentication required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	var loggerMW gin.HandlerFunc
	if logger, ok := conf.Logger.(share.LoggerInterface); ok {
		loggerMW = share.LoggerWithConfig(logger)
	} else {
		loggerMW = func(c *gin.Context) { c.Next() }
	}
	requestIDMW := share.RequestID()

	v1 := router.Group("/v1")
	v1.Use(requestIDMW)

	// ============ AUTH ENDPOINTS (share) ============
	shareAPI := v1.Group("/share")
	shareAPI.Use(loggerMW)

	// OIDC flow
	shareAPI.GET("/auth/oidc/login", commonShareController.OIDCLogin)
	shareAPI.GET("/auth/oidc/callback", commonShareController.OIDCCallback)

	// SAML flow
	shareAPI.GET("/auth/saml/login", commonShareController.SAMLLogin)
	shareAPI.POST("/auth/saml/callback", commonShareController.SAMLCallback)

	// SSO polling for CLI clients
	shareAPI.GET("/auth/sso/start", commonShareController.SSOStart)
	shareAPI.GET("/auth/sso/poll", commonShareController.SSOPoll)

	// Token management
	shareAPI.POST("/token/refresh", share.ForShare(commonRepository), commonShareController.RefreshToken)
	shareAPI.DELETE("/token", share.ForShare(commonRepository), commonShareController.Logout)
	shareAPI.GET("/token/validate", share.ForShare(commonRepository), commonShareController.ValidateToken)
	shareAPI.GET("/token/userinfo", share.ForShare(commonRepository), commonShareController.GetUserInfo)

	// ============ PUBLIC API (anonymous — no auth) ============
	publicAPI := v1.Group("/public")
	publicAPI.Use(loggerMW)
	publicAPI.Use(share.ForPublic(conf))
	// Business logic public routes go here

	// ============ INTERNAL API (app — authenticated) ============
	internalAPI := v1.Group("/internal")
	internalAPI.Use(loggerMW)
	internalAPI.Use(share.ForInternal(commonRepository))
	// Business logic internal routes go here

	// ============ PRIVATE API (admin — admin role required) ============
	privateAPI := v1.Group("/private")
	privateAPI.Use(loggerMW)
	privateAPI.Use(share.ForPrivate(commonRepository))
	// Business logic private routes go here

	_ = publicAPI
	_ = internalAPI
	_ = privateAPI

	return router
}

