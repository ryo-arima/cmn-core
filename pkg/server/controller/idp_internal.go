package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// IdPInternal handles user/group/member endpoints for authenticated (non-admin) users.
// All operations are proxied to the external IdP; cmn-core stores no auth data locally.
//
// Access control (enforced in this layer using JWT claims):
//   - User operations: self only (claims.UUID)
//   - Group create: any authenticated user
//   - Group read / update / delete / member operations: only for groups present in claims.Groups
//     JWT is issued on every request, so claims.Groups is always current.
type IdPInternal interface {
	// Own user
	GetMyUser(c *gin.Context)
	UpdateMyUser(c *gin.Context)

	// Groups the caller belongs to
	ListMyGroups(c *gin.Context)
	GetGroup(c *gin.Context)
	CreateGroup(c *gin.Context)
	UpdateGroup(c *gin.Context)
	DeleteGroup(c *gin.Context)

	// Member management (groups the caller belongs to)
	ListGroupMembers(c *gin.Context)
	AddGroupMember(c *gin.Context)
	RemoveGroupMember(c *gin.Context)
}

type idpInternal struct {
	idpUsecase usecase.IdP
}

// NewIdPInternal creates a new IdPInternal controller.
func NewIdPInternal(iu usecase.IdP) IdPInternal {
	return &idpInternal{idpUsecase: iu}
}

// isMemberOf returns true if groupID is present in the caller's groups claim.
func isMemberOf(groups []string, groupID string) bool {
	for _, g := range groups {
		if g == groupID {
			return true
		}
	}
	return false
}

// ---- Own user --------------------------------------------------------------

// GetMyUser returns the authenticated user's own profile from the IdP.
// GET /v1/internal/user
func (ic *idpInternal) GetMyUser(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	u, err := ic.idpUsecase.GetUser(c.Request.Context(), claims.UUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "IDP_USER_GET_404", "message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, response.SingleIdPUser{
		Code:    "SUCCESS",
		Message: "ok",
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

// UpdateMyUser updates the authenticated user's own profile in the IdP.
// PUT /v1/internal/user
func (ic *idpInternal) UpdateMyUser(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	var req request.UpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_USER_UPDATE_400", "message": "Invalid request body"})
		return
	}
	if err := ic.idpUsecase.UpdateUser(c.Request.Context(), claims.UUID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_UPDATE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "updated"})
}

// ---- Groups ----------------------------------------------------------------

// ListMyGroups lists the groups the authenticated user belongs to (from JWT claims),
// fetching each group's details from the IdP.
// GET /v1/internal/groups
func (ic *idpInternal) ListMyGroups(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	resp := make([]response.IdPGroup, 0, len(claims.Groups))
	for _, gid := range claims.Groups {
		g, err := ic.idpUsecase.GetGroup(c.Request.Context(), gid)
		if err != nil {
			continue // skip groups that can't be fetched
		}
		resp = append(resp, response.IdPGroup{ID: g.ID, Name: g.Name, Path: g.Path})
	}
	c.JSON(http.StatusOK, response.IdPGroups{Code: "SUCCESS", Message: "ok", Groups: resp})
}

// CreateGroup creates a new group in the IdP. Any authenticated user may create a group.
// POST /v1/internal/groups
func (ic *idpInternal) CreateGroup(c *gin.Context) {
	_, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	var req request.CreateGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_GROUP_CREATE_400", "message": "Invalid request body"})
		return
	}
	g, err := ic.idpUsecase.CreateGroup(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_CREATE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.SingleIdPGroup{
		Code:    "SUCCESS",
		Message: "created",
		Group:   &response.IdPGroup{ID: g.ID, Name: g.Name, Path: g.Path},
	})
}

// UpdateGroup updates a group — only accessible if the caller is a member.
// PUT /v1/internal/groups/:id
func (ic *idpInternal) UpdateGroup(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	groupID := c.Param("id")
	if !isMemberOf(claims.Groups, groupID) {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_GROUP_UPDATE_403", "message": "Access denied"})
		return
	}
	var req request.UpdateGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_GROUP_UPDATE_400", "message": "Invalid request body"})
		return
	}
	if err := ic.idpUsecase.UpdateGroup(c.Request.Context(), groupID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_UPDATE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "updated"})
}

// DeleteGroup deletes a group — only accessible if the caller is a member.
// DELETE /v1/internal/groups/:id
func (ic *idpInternal) DeleteGroup(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	groupID := c.Param("id")
	if !isMemberOf(claims.Groups, groupID) {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_GROUP_DELETE_403", "message": "Access denied"})
		return
	}
	if err := ic.idpUsecase.DeleteGroup(c.Request.Context(), groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_DELETE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "deleted"})
}

// GetGroup returns a single group — only accessible if the caller is a member.
// GET /v1/internal/group?id=...
func (ic *idpInternal) GetGroup(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	groupID := c.Query("id")
	if !isMemberOf(claims.Groups, groupID) {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_GROUP_GET_403", "message": "Access denied"})
		return
	}
	g, err := ic.idpUsecase.GetGroup(c.Request.Context(), groupID)
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

// ---- Members ---------------------------------------------------------------

// ListGroupMembers lists members of a group the caller belongs to.
// GET /v1/internal/members?group_id=...
func (ic *idpInternal) ListGroupMembers(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	groupID := c.Query("group_id")
	if !isMemberOf(claims.Groups, groupID) {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_MEMBER_LIST_403", "message": "Access denied"})
		return
	}
	users, err := ic.idpUsecase.ListGroupMembers(c.Request.Context(), groupID)
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

// AddGroupMember adds a user to a group the caller belongs to.
// POST /v1/internal/member/:group_id
func (ic *idpInternal) AddGroupMember(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	groupID := c.Param("group_id")
	if !isMemberOf(claims.Groups, groupID) {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_MEMBER_ADD_403", "message": "Access denied"})
		return
	}
	var req request.AddGroupMember
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_MEMBER_ADD_400", "message": "Invalid request body"})
		return
	}
	if err := ic.idpUsecase.AddUserToGroup(c.Request.Context(), req.UserID, groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_ADD_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "member added"})
}

// RemoveGroupMember removes a user from a group the caller belongs to.
// DELETE /v1/internal/member/:group_id
func (ic *idpInternal) RemoveGroupMember(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	groupID := c.Param("group_id")
	if !isMemberOf(claims.Groups, groupID) {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_MEMBER_REMOVE_403", "message": "Access denied"})
		return
	}
	var req request.AddGroupMember
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_MEMBER_REMOVE_400", "message": "Invalid request body"})
		return
	}
	if err := ic.idpUsecase.RemoveUserFromGroup(c.Request.Context(), req.UserID, groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_REMOVE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "member removed"})
}
