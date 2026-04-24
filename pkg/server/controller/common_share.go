package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ryo-arima/cmn-core/pkg/auth"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	share "github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// CommonShare provides authentication endpoints integrated with external IdPs (Keycloak / Authentik)
// via OIDC and SAML 2.0.
type CommonShare interface {
	// OIDC flow
	OIDCLogin(c *gin.Context)
	OIDCCallback(c *gin.Context)
	// SAML flow
	SAMLLogin(c *gin.Context)
	SAMLCallback(c *gin.Context)
	// SSO polling (for CLI / non-browser clients)
	SSOStart(c *gin.Context)
	SSOPoll(c *gin.Context)
	// Token management
	ValidateToken(c *gin.Context)
	GetUserInfo(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

type commonShare struct {
	CommonUsecase usecase.Common
	OIDCProvider  auth.Provider // nil when OIDC is not configured
	SAMLProvider  auth.Provider // nil when SAML is not configured
}

// NewCommonShare creates a new CommonShare controller.
func NewCommonShare(commonUsecase usecase.Common, oidcProvider, samlProvider auth.Provider) CommonShare {
	return &commonShare{
		CommonUsecase: commonUsecase,
		OIDCProvider:  oidcProvider,
		SAMLProvider:  samlProvider,
	}
}

// OIDCLogin initiates the OIDC Authorization Code flow.
// A random state is generated, stored in Redis, and the browser is redirected to the IdP.
//
// Route: GET /v1/share/auth/oidc/login
func (rcvr commonShare) OIDCLogin(c *gin.Context) {
	if rcvr.OIDCProvider == nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code":    "OIDC_NOT_CONFIGURED",
			"message": "OIDC provider is not configured",
		})
		return
	}

	state := uuid.New().String()
	if err := rcvr.CommonUsecase.StoreAuthState(c.Request.Context(), state); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "OIDC_LOGIN_001",
			"message": "Failed to generate auth state",
		})
		return
	}

	loginURL, err := rcvr.OIDCProvider.LoginURL(state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "OIDC_LOGIN_002",
			"message": "Failed to build OIDC login URL",
		})
		return
	}

	c.Redirect(http.StatusFound, loginURL)
}

// OIDCCallback handles the authorization code callback from the OIDC provider.
//
// Route: GET /v1/share/auth/oidc/callback
func (rcvr commonShare) OIDCCallback(c *gin.Context) {
	if rcvr.OIDCProvider == nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": "OIDC_NOT_CONFIGURED",
		})
		return
	}

	var cb request.OIDCCallback
	if err := c.ShouldBindQuery(&cb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "OIDC_CALLBACK_001",
			"message": "Invalid callback parameters",
		})
		return
	}

	// Validate and consume the state (CSRF guard)
	if err := rcvr.CommonUsecase.ConsumeAuthState(c.Request.Context(), cb.State); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "OIDC_CALLBACK_002",
			"message": "Invalid or expired state parameter",
		})
		return
	}

	claims, err := rcvr.OIDCProvider.HandleCallback(c.Request.Context(), c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "OIDC_CALLBACK_003",
			"message": "OIDC token exchange failed",
		})
		return
	}

	role := rcvr.CommonUsecase.ResolveRole(claims.Email)
	tokenPair, err := rcvr.CommonUsecase.GenerateTokenPair(0, uuid.New().String(), claims.Email, claims.Name, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "OIDC_CALLBACK_004",
			"message": "Failed to issue session token",
		})
		return
	}

	// Store for CLI polling (if state was a poll session ID)
	_ = rcvr.CommonUsecase.StoreTokenForState(c.Request.Context(), cb.State, tokenPair)

	c.JSON(http.StatusOK, &response.AuthCallback{
		Code:      "SUCCESS",
		Message:   "Authentication successful",
		TokenPair: tokenPair,
	})
}

// SAMLLogin initiates SAML SP-initiated SSO.
//
// Route: GET /v1/share/auth/saml/login
func (rcvr commonShare) SAMLLogin(c *gin.Context) {
	if rcvr.SAMLProvider == nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code":    "SAML_NOT_CONFIGURED",
			"message": "SAML provider is not configured",
		})
		return
	}

	state := uuid.New().String()
	if err := rcvr.CommonUsecase.StoreAuthState(c.Request.Context(), state); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "SAML_LOGIN_001",
			"message": "Failed to generate auth state",
		})
		return
	}

	loginURL, err := rcvr.SAMLProvider.LoginURL(state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "SAML_LOGIN_002",
			"message": "Failed to build SAML login URL",
		})
		return
	}

	c.Redirect(http.StatusFound, loginURL)
}

// SAMLCallback handles the ACS POST from the SAML IdP.
//
// Route: POST /v1/share/auth/saml/callback
func (rcvr commonShare) SAMLCallback(c *gin.Context) {
	if rcvr.SAMLProvider == nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": "SAML_NOT_CONFIGURED",
		})
		return
	}

	var cb request.SAMLCallback
	if err := c.ShouldBind(&cb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "SAML_CALLBACK_001",
			"message": "Invalid SAML response",
		})
		return
	}

	relayState := cb.RelayState
	if relayState != "" {
		if err := rcvr.CommonUsecase.ConsumeAuthState(c.Request.Context(), relayState); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "SAML_CALLBACK_002",
				"message": "Invalid or expired RelayState",
			})
			return
		}
	}

	claims, err := rcvr.SAMLProvider.HandleCallback(c.Request.Context(), c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "SAML_CALLBACK_003",
			"message": "SAML assertion validation failed",
		})
		return
	}

	role := rcvr.CommonUsecase.ResolveRole(claims.Email)
	tokenPair, err := rcvr.CommonUsecase.GenerateTokenPair(0, uuid.New().String(), claims.Email, claims.Name, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "SAML_CALLBACK_004",
			"message": "Failed to issue session token",
		})
		return
	}

	if relayState != "" {
		_ = rcvr.CommonUsecase.StoreTokenForState(c.Request.Context(), relayState, tokenPair)
	}

	c.JSON(http.StatusOK, &response.AuthCallback{
		Code:      "SUCCESS",
		Message:   "Authentication successful",
		TokenPair: tokenPair,
	})
}

// SSOStart returns a login URL and session ID for CLI/polling-based SSO.
// The caller opens the URL in a browser, then polls SSOPoll to retrieve the token.
//
// Route: GET /v1/share/auth/sso/start?provider=oidc|saml
func (rcvr commonShare) SSOStart(c *gin.Context) {
	provider := c.DefaultQuery("provider", "oidc")

	var p auth.Provider
	switch provider {
	case "oidc":
		p = rcvr.OIDCProvider
	case "saml":
		p = rcvr.SAMLProvider
	}

	if p == nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code":    "SSO_NOT_CONFIGURED",
			"message": "Requested SSO provider is not configured: " + provider,
		})
		return
	}

	sessionID := uuid.New().String()
	if err := rcvr.CommonUsecase.StoreAuthState(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "SSO_START_001",
			"message": "Failed to create SSO session",
		})
		return
	}

	loginURL, err := p.LoginURL(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "SSO_START_002",
			"message": "Failed to build SSO login URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       "SUCCESS",
		"login_url":  loginURL,
		"session_id": sessionID,
		"provider":   provider,
	})
}

// SSOPoll returns the token pair once the SSO callback has completed.
// Returns 202 Accepted while authentication is still pending.
//
// Route: GET /v1/share/auth/sso/poll?session_id=<id>
func (rcvr commonShare) SSOPoll(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "SSO_POLL_001",
			"message": "session_id query parameter is required",
		})
		return
	}

	tokenPair, found, err := rcvr.CommonUsecase.GetTokenForState(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "SSO_POLL_002",
			"message": "Failed to check SSO session",
		})
		return
	}
	if !found {
		c.JSON(http.StatusAccepted, gin.H{
			"code":    "SSO_PENDING",
			"message": "Authentication is still pending",
		})
		return
	}

	c.JSON(http.StatusOK, &response.AuthCallback{
		Code:      "SUCCESS",
		Message:   "Authentication successful",
		TokenPair: tokenPair,
	})
}

// ValidateToken validates a JWT token issued after OIDC/SAML authentication.
//
// Route: GET /v1/share/token/validate
func (rcvr commonShare) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "TOKEN_VALIDATE_001",
			"message": "Authorization header required",
		})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := rcvr.CommonUsecase.ValidateJWTToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "TOKEN_VALIDATE_002",
			"message": "Invalid token",
		})
		return
	}

	isInvalidated, err := rcvr.CommonUsecase.IsTokenInvalidated(c.Request.Context(), claims.Jti)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "TOKEN_VALIDATE_003",
			"message": "Failed to check token status",
		})
		return
	}
	if isInvalidated {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "TOKEN_VALIDATE_004",
			"message": "Token has been revoked",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "Token is valid",
		"data": gin.H{
			"uuid":       claims.UUID,
			"email":      claims.Email,
			"name":       claims.Name,
			"role":       claims.Role,
			"expires_at": claims.ExpiresAt,
		},
	})
}

// GetUserInfo returns the authenticated user's information from the token context.
//
// Route: GET /v1/share/token/userinfo
func (rcvr commonShare) GetUserInfo(c *gin.Context) {
	userClaims, exists := share.GetUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "TOKEN_USERINFO_001",
			"message": "User not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, &response.Commons{
		Code:    "SUCCESS",
		Message: "User information retrieved successfully",
		Commons: []response.Common{
			{UUID: userClaims.UUID},
		},
	})
}

// RefreshToken exchanges a valid refresh token for a new token pair.
//
// Route: POST /v1/share/token/refresh
func (rcvr commonShare) RefreshToken(c *gin.Context) {
	var req request.RefreshToken
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &response.RefreshToken{
			Code:    "TOKEN_REFRESH_001",
			Message: "Invalid request format",
		})
		return
	}

	claims, err := rcvr.CommonUsecase.ValidateJWTToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.RefreshToken{
			Code:    "TOKEN_REFRESH_002",
			Message: "Invalid refresh token",
		})
		return
	}

	tokenPair, err := rcvr.CommonUsecase.GenerateTokenPair(
		claims.UserID, claims.UUID, claims.Email, claims.Name, claims.Role,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.RefreshToken{
			Code:    "TOKEN_REFRESH_003",
			Message: "Failed to generate new tokens",
		})
		return
	}

	c.JSON(http.StatusOK, &response.RefreshToken{
		Code:      "SUCCESS",
		Message:   "Token refreshed successfully",
		TokenPair: tokenPair,
	})
}

// Logout invalidates the current access token.
//
// Route: DELETE /v1/share/token
func (rcvr commonShare) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "TOKEN_LOGOUT_001",
			"message": "Authorization header required",
		})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "TOKEN_LOGOUT_002",
			"message": "Token not found in Authorization header",
		})
		return
	}

	if err := rcvr.CommonUsecase.InvalidateToken(c.Request.Context(), tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "TOKEN_LOGOUT_003",
			"message": "Failed to invalidate token",
		})
		return
	}

	rcvr.CommonUsecase.DeleteTokenCache(tokenString)

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "Logout successful",
	})
}
