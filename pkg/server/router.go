package server

import (
	"context"
	"log"
	"net/http"
	"strings"

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
		// providerURL is where we fetch the OIDC discovery document (internal Docker URL).
		// IssuerURL is the public issuer that appears in the JWT "iss" claim.
		// They differ in Docker environments (e.g. Casdoor: issuer=http://localhost:9000,
		// provider=http://casdoor:8000).
		providerURL := oidcCfg.ProviderURL
		if providerURL == "" {
			providerURL = oidcCfg.IssuerURL
		}
		// Build a custom HTTP client that rewrites the public issuer host to the
		// internal provider host for all OIDC-related fetches (discovery + JWKS).
		// This is necessary because the JWKS URI embedded in the discovery document
		// uses the public hostname (e.g. localhost:9000), which is unreachable from
		// inside the Docker network.
		oidcHTTPClient := &http.Client{
			Transport: newHostRewriteTransport(oidcCfg.IssuerURL, providerURL),
		}
		// InsecureIssuerURLContext tells go-oidc to accept the issuer from the
		// discovery document even though we fetched it from a different URL.
		oidcCtx := gooidc.InsecureIssuerURLContext(context.Background(), oidcCfg.IssuerURL)
		oidcCtx = gooidc.ClientContext(oidcCtx, oidcHTTPClient)
		provider, err := gooidc.NewProvider(oidcCtx, providerURL)
		if err != nil {
			log.Printf("OIDC init failed, JWT validation disabled: %v", err)
		} else {
			verifier = provider.Verifier(&gooidc.Config{SkipClientIDCheck: true})
		}
	}

	// ============================================================
	// Repository / Usecase / Controller initialization
	// ============================================================

	// -- common --
	commonRepository := repository.NewCommon(conf, verifier)
	commonUsecase := usecase.NewCommon(commonRepository)
	commonShareCtrl := controller.NewCommonShare(commonUsecase)

	// -- IdP --
	idpManager, err := repository.NewIdPManager(conf)
	if err != nil {
		log.Fatalf("IdP init failed: %v", err)
	}
	idpUsecase := usecase.NewIdP(idpManager)
	idpPublicCtrl := controller.NewIdPPublic(idpUsecase)
	idpInternalCtrl := controller.NewIdPInternal(idpUsecase, commonUsecase)
	idpPrivateCtrl := controller.NewIdPPrivate(idpUsecase, commonUsecase)

	// -- resource --
	resourceRepo := repository.NewResource(conf)
	resourceUsecase := usecase.NewResource(resourceRepo)
	resourceInternalCtrl := controller.NewResourceInternal(resourceUsecase)
	resourcePrivateCtrl := controller.NewResourcePrivate(resourceUsecase)

	// ============================================================
	// Router setup
	// ============================================================

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

	// ============ PUBLIC API (anonymous — no auth) ============
	publicAPI := v1.Group("/public")
	publicAPI.Use(loggerMW)
	publicAPI.Use(share.ForPublic())

	publicAPI.POST("/user", idpPublicCtrl.RegisterUser)
	publicAPI.POST("/login", idpPublicCtrl.Login)

	// ============ SHARE API (any authenticated role) ============
	shareAPI := v1.Group("/share")
	shareAPI.Use(loggerMW)
	shareAPI.Use(share.ForShare(commonRepository))

	shareAPI.GET("/token/validate", commonShareCtrl.ValidateToken)
	shareAPI.GET("/token/userinfo", commonShareCtrl.GetUserInfo)
	shareAPI.GET("/user", idpInternalCtrl.GetMyUser)
	shareAPI.PUT("/user", idpInternalCtrl.UpdateMyUser)

	// ============ INTERNAL API (app — authenticated) ============
	internalAPI := v1.Group("/internal")
	internalAPI.Use(loggerMW)
	internalAPI.Use(share.ForInternal(commonRepository))

	// Own user (or any user by ?id=) and group users
	internalAPI.GET("/user", idpInternalCtrl.GetMyUser)
	internalAPI.PUT("/user", idpInternalCtrl.UpdateMyUser)
	internalAPI.GET("/users", idpInternalCtrl.ListGroupUsers)

	// Groups the caller belongs to
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

	// Members
	privateAPI.GET("/members", idpPrivateCtrl.ListGroupMembers)
	privateAPI.POST("/member/:group_id", idpPrivateCtrl.AddGroupMember)
	privateAPI.DELETE("/member/:group_id", idpPrivateCtrl.RemoveGroupMember)

	// Resources
	privateAPI.GET("/resources", resourcePrivateCtrl.ListAllResources)
	privateAPI.POST("/resources", resourcePrivateCtrl.CreateResource)
	privateAPI.GET("/resource", resourcePrivateCtrl.GetResource)
	privateAPI.PUT("/resources/:uuid", resourcePrivateCtrl.UpdateResource)
	privateAPI.DELETE("/resources/:uuid", resourcePrivateCtrl.DeleteResource)
	privateAPI.GET("/resource/groups", resourcePrivateCtrl.GetResourceGroupRoles)
	privateAPI.PUT("/resources/:uuid/groups", resourcePrivateCtrl.SetResourceGroupRole)
	privateAPI.DELETE("/resources/:uuid/groups", resourcePrivateCtrl.DeleteResourceGroupRole)

	return router
}

// hostRewriteTransport rewrites the host of outgoing HTTP requests.
// This is used to redirect OIDC-related requests (discovery + JWKS) from the
// public issuer URL (e.g. http://localhost:9000) to the internal provider URL
// (e.g. http://casdoor:8000) when running inside a Docker network.
type hostRewriteTransport struct {
	fromHost string
	toHost   string
	base     http.RoundTripper
}

func newHostRewriteTransport(fromURL, toURL string) *hostRewriteTransport {
	stripScheme := func(u string) string {
		u = strings.TrimPrefix(u, "https://")
		u = strings.TrimPrefix(u, "http://")
		return u
	}
	return &hostRewriteTransport{
		fromHost: stripScheme(fromURL),
		toHost:   stripScheme(toURL),
		base:     http.DefaultTransport,
	}
}

func (t *hostRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host == t.fromHost {
		req = req.Clone(req.Context())
		req.URL.Host = t.toHost
	}
	return t.base.RoundTrip(req)
}
