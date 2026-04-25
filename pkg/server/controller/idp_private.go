package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// IdPPrivate handles user/group/member admin endpoints that proxy to the
// configured external identity provider (Keycloak or Casdoor).
// All routes require the admin role (ForPrivate middleware).
type IdPPrivate interface {
	// User CRUD
	ListUsers(c *gin.Context)
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)

	// Group CRUD
	ListGroups(c *gin.Context)
	GetGroup(c *gin.Context)
	CreateGroup(c *gin.Context)
	UpdateGroup(c *gin.Context)
	DeleteGroup(c *gin.Context)

	// Group membership
	ListGroupMembers(c *gin.Context)
	AddGroupMember(c *gin.Context)
	RemoveGroupMember(c *gin.Context)
}

type idpPrivate struct {
	idpUsecase usecase.IdP
}

// NewIdPPrivate creates a new IdPPrivate controller.
func NewIdPPrivate(iu usecase.IdP) IdPPrivate {
	return &idpPrivate{idpUsecase: iu}
}

// ---- Users -----------------------------------------------------------------

// ListUsers lists all users from the IdP.
// GET /v1/private/users
func (ic *idpPrivate) ListUsers(c *gin.Context) {
	users, err := ic.idpUsecase.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_LIST_001", "message": err.Error()})
		return
	}
	resp := make([]response.IdPUser, 0, len(users))
	for _, u := range users {
		resp = append(resp, response.IdPUser{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Enabled:   u.Enabled,
			CreatedAt: u.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, response.IdPUsers{Code: "SUCCESS", Message: "ok", Users: resp})
}

// GetUser returns a single user from the IdP.
// GET /v1/private/user?id=...
func (ic *idpPrivate) GetUser(c *gin.Context) {
	u, err := ic.idpUsecase.GetUser(c.Request.Context(), c.Query("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "IDP_USER_GET_404", "message": "User not found"})
		return
	}
	resp := &response.IdPUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Enabled:   u.Enabled,
		CreatedAt: u.CreatedAt,
	}
	c.JSON(http.StatusOK, response.SingleIdPUser{Code: "SUCCESS", Message: "ok", User: resp})
}

// CreateUser creates a new user in the IdP.
// POST /v1/private/users
func (ic *idpPrivate) CreateUser(c *gin.Context) {
	var req request.CreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_USER_CREATE_001", "message": "Invalid request body"})
		return
	}
	u, err := ic.idpUsecase.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_CREATE_002", "message": err.Error()})
		return
	}
	resp := &response.IdPUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Enabled:   u.Enabled,
		CreatedAt: u.CreatedAt,
	}
	c.JSON(http.StatusCreated, response.SingleIdPUser{Code: "SUCCESS", Message: "created", User: resp})
}

// UpdateUser updates an existing user in the IdP.
// PUT /v1/private/users/:id
func (ic *idpPrivate) UpdateUser(c *gin.Context) {
	var req request.UpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_USER_UPDATE_001", "message": "Invalid request body"})
		return
	}
	if err := ic.idpUsecase.UpdateUser(c.Request.Context(), c.Param("id"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_UPDATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "updated"})
}

// DeleteUser deletes a user from the IdP.
// DELETE /v1/private/users/:id
func (ic *idpPrivate) DeleteUser(c *gin.Context) {
	if err := ic.idpUsecase.DeleteUser(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_DELETE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "deleted"})
}

// ---- Groups ----------------------------------------------------------------

// ListGroups lists all groups from the IdP.
// GET /v1/private/groups
func (ic *idpPrivate) ListGroups(c *gin.Context) {
	groups, err := ic.idpUsecase.ListGroups(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_LIST_001", "message": err.Error()})
		return
	}
	resp := make([]response.IdPGroup, 0, len(groups))
	for _, g := range groups {
		resp = append(resp, response.IdPGroup{ID: g.ID, Name: g.Name, Path: g.Path})
	}
	c.JSON(http.StatusOK, response.IdPGroups{Code: "SUCCESS", Message: "ok", Groups: resp})
}

// GetGroup returns a single group from the IdP.
// GET /v1/private/group?id=...
func (ic *idpPrivate) GetGroup(c *gin.Context) {
	g, err := ic.idpUsecase.GetGroup(c.Request.Context(), c.Query("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "IDP_GROUP_GET_404", "message": "Group not found"})
		return
	}
	c.JSON(http.StatusOK, response.SingleIdPGroup{
		Code:    "SUCCESS",
		Message: "ok",
		Group:   &response.IdPGroup{ID: g.ID, Name: g.Name, Path: g.Path},
	})
}

// CreateGroup creates a new group in the IdP.
// POST /v1/private/groups
func (ic *idpPrivate) CreateGroup(c *gin.Context) {
	var req request.CreateGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_GROUP_CREATE_001", "message": "Invalid request body"})
		return
	}
	g, err := ic.idpUsecase.CreateGroup(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_CREATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.SingleIdPGroup{
		Code:    "SUCCESS",
		Message: "created",
		Group:   &response.IdPGroup{ID: g.ID, Name: g.Name, Path: g.Path},
	})
}

// UpdateGroup updates an existing group in the IdP.
// PUT /v1/private/groups/:id
func (ic *idpPrivate) UpdateGroup(c *gin.Context) {
	var req request.UpdateGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_GROUP_UPDATE_001", "message": "Invalid request body"})
		return
	}
	if err := ic.idpUsecase.UpdateGroup(c.Request.Context(), c.Param("id"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_UPDATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "updated"})
}

// DeleteGroup deletes a group from the IdP.
// DELETE /v1/private/groups/:id
func (ic *idpPrivate) DeleteGroup(c *gin.Context) {
	if err := ic.idpUsecase.DeleteGroup(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_DELETE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "deleted"})
}

// ---- Members ---------------------------------------------------------------

// ListGroupMembers lists members of a group.
// GET /v1/private/members?group_id=...
func (ic *idpPrivate) ListGroupMembers(c *gin.Context) {
	users, err := ic.idpUsecase.ListGroupMembers(c.Request.Context(), c.Query("group_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_LIST_001", "message": err.Error()})
		return
	}
	resp := make([]response.IdPUser, 0, len(users))
	for _, u := range users {
		resp = append(resp, response.IdPUser{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Enabled:   u.Enabled,
			CreatedAt: u.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, response.IdPUsers{Code: "SUCCESS", Message: "ok", Users: resp})
}

// AddGroupMember adds a user to a group.
// POST /v1/private/member/:group_id  (user_id specified in request body)
func (ic *idpPrivate) AddGroupMember(c *gin.Context) {
	var req request.AddGroupMember
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_MEMBER_ADD_001", "message": "Invalid request body"})
		return
	}
	if err := ic.idpUsecase.AddUserToGroup(c.Request.Context(), req.UserID, c.Param("group_id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_ADD_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "member added"})
}

// RemoveGroupMember removes a user from a group.
// DELETE /v1/private/member/:group_id  (user_id specified in request body)
func (ic *idpPrivate) RemoveGroupMember(c *gin.Context) {
	var req request.AddGroupMember
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_MEMBER_REMOVE_400", "message": "Invalid request body"})
		return
	}
	if err := ic.idpUsecase.RemoveUserFromGroup(c.Request.Context(), req.UserID, c.Param("group_id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_REMOVE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "member removed"})
}
