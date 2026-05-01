package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// UserPublic handles unauthenticated user registration and login endpoints.
// Routes are mounted under /v1/public (no auth middleware required).
type UserPublic interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
}

type userPublic struct {
	userUsecase usecase.User
}

// NewUserPublic creates a new UserPublic controller.
func NewUserPublic(uu usecase.User) UserPublic {
	return &userPublic{userUsecase: uu}
}

// RegisterUser creates a new user in the IdP without requiring authentication.
// Route: POST /v1/public/user
func (rcvr *userPublic) RegisterUser(c *gin.Context) {
	var req request.RrCreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "PUB_REGISTER_400", "message": err.Error()})
		return
	}
	u, err := rcvr.userUsecase.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "PUB_REGISTER_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.RrSingleIdPUser{
		Code:    "SUCCESS",
		Message: "user created",
		User: &response.RrIdPUser{
			ID:        u.ID,
			UUID:      u.UUID,
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
// Route: POST /v1/public/login
func (rcvr *userPublic) Login(c *gin.Context) {
	var req request.RrLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "PUB_LOGIN_400", "message": err.Error()})
		return
	}
	token, err := rcvr.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
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
