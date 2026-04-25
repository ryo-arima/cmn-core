package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	gooidc "github.com/coreos/go-oidc/v3/oidc"
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

	// Initialize OIDC verifier for JWT validation
	var verifier *gooidc.IDTokenVerifier
	oidcCfg := conf.YamlConfig.Application.Server.Auth.OIDC
	if oidcCfg.IssuerURL != "" {
		provider, err := gooidc.NewProvider(context.Background(), oidcCfg.IssuerURL)
		if err != nil {
			log.Printf("OIDC init failed, JWT validation disabled: %v", err)
		} else {
			verifier = provider.Verifier(&gooidc.Config{SkipClientIDCheck: true})
		}
	}

	commonRepository := repository.NewCommon(conf, verifier)
	commonUsecase := usecase.NewCommon(commonRepository)
	commonShareController := controller.NewCommonShare(commonUsecase)

	// Initialize IdP manager (required)
	idpManager, err := repository.NewIdPManager(conf)
	if err != nil {
		panic("IdP init failed: " + err.Error())
	}
	idpUsecase := usecase.NewIdP(idpManager)
	idpInternalCtrl := controller.NewIdPInternal(idpUsecase)
	idpPrivateCtrl := controller.NewIdPPrivate(idpUsecase)

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

	// Token management (JWT validation via IdP JWKS)
	shareAPI.GET("/token/validate", share.ForShare(commonRepository), commonShareController.ValidateToken)
	shareAPI.GET("/token/userinfo", share.ForShare(commonRepository), commonShareController.GetUserInfo)

	// ============ PUBLIC API (anonymous — no auth) ============
	publicAPI := v1.Group("/public")
	publicAPI.Use(loggerMW)
	publicAPI.Use(share.ForPublic())
	// Business logic public routes go here

	// ============ INTERNAL API (app — authenticated) ============
	internalAPI := v1.Group("/internal")
	internalAPI.Use(loggerMW)
	internalAPI.Use(share.ForInternal(commonRepository))

	// Resource routes (authenticated users)
	resourceRepo := repository.NewResource(conf)
	resourceUsecase := usecase.NewResource(resourceRepo)
	resourceInternalCtrl := controller.NewResourceInternal(resourceUsecase)

	// Own user
	internalAPI.GET("/user", idpInternalCtrl.GetMyUser)
	internalAPI.PUT("/user", idpInternalCtrl.UpdateMyUser)

	// Groups the caller belongs to (JWT is issued per-request; claims.Groups is always current)
	internalAPI.GET("/groups", idpInternalCtrl.ListMyGroups)
	internalAPI.POST("/groups", idpInternalCtrl.CreateGroup)
	internalAPI.GET("/group", idpInternalCtrl.GetGroup)
	internalAPI.PUT("/groups/:id", idpInternalCtrl.UpdateGroup)
	internalAPI.DELETE("/groups/:id", idpInternalCtrl.DeleteGroup)

	// Member management
	internalAPI.GET("/members", idpInternalCtrl.ListGroupMembers)
	internalAPI.POST("/member/:group_id", idpInternalCtrl.AddGroupMember)
	internalAPI.DELETE("/member/:group_id", idpInternalCtrl.RemoveGroupMember)

	// Resources
	internalAPI.GET("/resources", resourceInternalCtrl.ListResources)
	internalAPI.POST("/resources", resourceInternalCtrl.CreateResource)
	internalAPI.GET("/resource", resourceInternalCtrl.GetResource)
	internalAPI.PUT("/resources/:uuid", resourceInternalCtrl.UpdateResource)
	internalAPI.DELETE("/resources/:uuid", resourceInternalCtrl.DeleteResource)
	internalAPI.GET("/resource/groups", resourceInternalCtrl.GetResourceGroupRoles)
	internalAPI.PUT("/resources/:uuid/groups", resourceInternalCtrl.SetResourceGroupRole)
	internalAPI.DELETE("/resources/:uuid/groups/:group_uuid", resourceInternalCtrl.DeleteResourceGroupRole)

	// ============ PRIVATE API (admin — admin role required) ============
	privateAPI := v1.Group("/private")
	privateAPI.Use(loggerMW)
	privateAPI.Use(share.ForPrivate(commonRepository))

	// Users
	privateAPI.GET("/users", idpPrivateCtrl.ListUsers)
	privateAPI.POST("/users", idpPrivateCtrl.CreateUser)
	privateAPI.GET("/user", idpPrivateCtrl.GetUser)
	privateAPI.PUT("/users/:id", idpPrivateCtrl.UpdateUser)
	privateAPI.DELETE("/users/:id", idpPrivateCtrl.DeleteUser)

	// Groups
	privateAPI.GET("/groups", idpPrivateCtrl.ListGroups)
	privateAPI.POST("/groups", idpPrivateCtrl.CreateGroup)
	privateAPI.GET("/group", idpPrivateCtrl.GetGroup)
	privateAPI.PUT("/groups/:id", idpPrivateCtrl.UpdateGroup)
	privateAPI.DELETE("/groups/:id", idpPrivateCtrl.DeleteGroup)

	// Members (group membership)
	privateAPI.GET("/members", idpPrivateCtrl.ListGroupMembers)
	privateAPI.POST("/member/:group_id", idpPrivateCtrl.AddGroupMember)
	privateAPI.DELETE("/member/:group_id", idpPrivateCtrl.RemoveGroupMember)

	// Resources
	resourcePrivateCtrl := controller.NewResourcePrivate(resourceUsecase)
	privateAPI.GET("/resources", resourcePrivateCtrl.ListAllResources)
	privateAPI.POST("/resources", resourcePrivateCtrl.CreateResource)
	privateAPI.GET("/resource", resourcePrivateCtrl.GetResource)
	privateAPI.PUT("/resources/:uuid", resourcePrivateCtrl.UpdateResource)
	privateAPI.DELETE("/resources/:uuid", resourcePrivateCtrl.DeleteResource)
	privateAPI.GET("/resource/groups", resourcePrivateCtrl.GetResourceGroupRoles)
	privateAPI.PUT("/resources/:uuid/groups", resourcePrivateCtrl.SetResourceGroupRole)
	privateAPI.DELETE("/resources/:uuid/groups", resourcePrivateCtrl.DeleteResourceGroupRole)

	_ = publicAPI

	return router
}

