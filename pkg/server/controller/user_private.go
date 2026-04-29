package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// UserPrivate handles user admin endpoints that proxy to the external IdP.
// All routes require the admin role (ForPrivate middleware).
type UserPrivate interface {
	ListUsers(c *gin.Context)
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userPrivate struct {
	userUsecase   usecase.User
	commonUsecase usecase.Common
}

// NewUserPrivate creates a new UserPrivate controller.
func NewUserPrivate(uu usecase.User, cu usecase.Common) UserPrivate {
	return &userPrivate{userUsecase: uu, commonUsecase: cu}
}

// ListUsers lists all users from the IdP.
// GET /v1/private/users
func (ic *userPrivate) ListUsers(c *gin.Context) {
	users, err := ic.userUsecase.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_LIST_001", "message": err.Error()})
		return
	}
	resp := make([]response.RrIdPUser, 0, len(users))
	for _, u := range users {
		resp = append(resp, response.RrIdPUser{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Enabled:   u.Enabled,
			Role:      ic.commonUsecase.ResolveRole(u.Email),
			CreatedAt: u.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, response.RrIdPUsers{Code: "SUCCESS", Message: "ok", Users: resp})
}

// GetUser returns a single user from the IdP.
// GET /v1/private/user?id=...
func (ic *userPrivate) GetUser(c *gin.Context) {
	u, err := ic.userUsecase.GetUser(c.Request.Context(), c.Query("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "IDP_USER_GET_404", "message": "User not found"})
		return
	}
	resp := &response.RrIdPUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Enabled:   u.Enabled,
		Role:      ic.commonUsecase.ResolveRole(u.Email),
		CreatedAt: u.CreatedAt,
	}
	c.JSON(http.StatusOK, response.RrSingleIdPUser{Code: "SUCCESS", Message: "ok", User: resp})
}

// CreateUser creates a new user in the IdP.
// POST /v1/private/users
func (ic *userPrivate) CreateUser(c *gin.Context) {
	var req request.RrCreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_USER_CREATE_001", "message": "Invalid request body"})
		return
	}
	u, err := ic.userUsecase.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_CREATE_002", "message": err.Error()})
		return
	}
	resp := &response.RrIdPUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Enabled:   u.Enabled,
		Role:      ic.commonUsecase.ResolveRole(u.Email),
		CreatedAt: u.CreatedAt,
	}
	c.JSON(http.StatusCreated, response.RrSingleIdPUser{Code: "SUCCESS", Message: "created", User: resp})
}

// UpdateUser updates an existing user in the IdP.
// PUT /v1/private/users/:id
func (ic *userPrivate) UpdateUser(c *gin.Context) {
	var req request.RrUpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_USER_UPDATE_001", "message": "Invalid request body"})
		return
	}
	if err := ic.userUsecase.UpdateUser(c.Request.Context(), c.Param("id"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_UPDATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "updated"})
}

// DeleteUser deletes a user from the IdP.
// DELETE /v1/private/users/:id
func (ic *userPrivate) DeleteUser(c *gin.Context) {
	if err := ic.userUsecase.DeleteUser(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_DELETE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "deleted"})
}
