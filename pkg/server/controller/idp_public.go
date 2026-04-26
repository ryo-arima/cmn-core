package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// IdPPublic handles unauthenticated endpoints for self-service user registration and login.
// Routes are mounted under /v1/public (no auth middleware required).
type IdPPublic interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
}

type idpPublic struct {
	idpUsecase usecase.IdP
}

// NewIdPPublic creates a new IdPPublic controller.
func NewIdPPublic(iu usecase.IdP) IdPPublic {
	return &idpPublic{idpUsecase: iu}
}

// RegisterUser creates a new user in the IdP without requiring authentication.
//
// Route: POST /v1/public/user
func (ic *idpPublic) RegisterUser(c *gin.Context) {
	var req request.CreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "PUB_REGISTER_400", "message": err.Error()})
		return
	}
	u, err := ic.idpUsecase.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "PUB_REGISTER_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.SingleIdPUser{
		Code:    "SUCCESS",
		Message: "user created",
		User: &response.IdPUser{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Enabled:   u.Enabled,
			CreatedAt: u.CreatedAt,
		},
	})
}

// Login authenticates a user via the IdP using the Resource Owner Password Credentials grant
// and returns the issued access token.
//
// Route: POST /v1/public/login
func (ic *idpPublic) Login(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "PUB_LOGIN_400", "message": err.Error()})
		return
	}
	token, err := ic.idpUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "PUB_LOGIN_401", "message": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":         "SUCCESS",
		"message":      "login successful",
		"access_token": token,
	})
}
