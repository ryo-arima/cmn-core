package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	share "github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// CommonShare exposes token-related endpoints.
type CommonShare interface {
	ValidateToken(c *gin.Context)
	GetUserInfo(c *gin.Context)
}

type commonShare struct {
	CommonUsecase usecase.Common
}

// NewCommonShare creates a new CommonShare controller.
func NewCommonShare(commonUsecase usecase.Common) CommonShare {
	return &commonShare{CommonUsecase: commonUsecase}
}

// ValidateToken returns 200 OK when the middleware has already validated the token.
//
// Route: GET /v1/share/token/validate
func (rcvr commonShare) ValidateToken(c *gin.Context) {
	userClaims, exists := share.GetUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "TOKEN_VALIDATE_001",
			"message": "User not authenticated",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "Token is valid",
		"data": gin.H{
			"uuid":       userClaims.UUID,
			"email":      userClaims.Email,
			"name":       userClaims.Name,
			"role":       userClaims.Role,
			"expires_at": userClaims.ExpiresAt,
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
	c.JSON(http.StatusOK, &response.RrCommons{
		Code:    "SUCCESS",
		Message: "User information retrieved successfully",
		Commons: []response.RrCommon{
			{UUID: userClaims.UUID},
		},
	})
}


