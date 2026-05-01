package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// GroupInternal handles group endpoints for authenticated (non-admin) users.
type GroupInternal interface {
	ListMyGroups(c *gin.Context)
	GetGroup(c *gin.Context)
	CreateGroup(c *gin.Context)
	UpdateGroup(c *gin.Context)
	DeleteGroup(c *gin.Context)
}

type groupInternal struct {
	groupUsecase  usecase.Group
	memberUsecase usecase.Member
	commonUsecase usecase.Common
}

// NewGroupInternal creates a new GroupInternal controller.
func NewGroupInternal(gu usecase.Group, mu usecase.Member, cu usecase.Common) GroupInternal {
	return &groupInternal{groupUsecase: gu, memberUsecase: mu, commonUsecase: cu}
}

// ListMyGroups lists the groups the authenticated user belongs to (from JWT claims),
// fetching each group's details from the IdP.
// GET /v1/internal/groups
func (rcvr *groupInternal) ListMyGroups(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	resp := make([]response.RrIdPGroup, 0, len(claims.Groups))
	for _, gid := range claims.Groups {
		// Strip org prefix: JWT uses "cmn/group001", usecase expects "group001"
		g, err := rcvr.groupUsecase.GetGroup(c.Request.Context(), groupName(gid))
		if err != nil {
			continue // skip groups that can't be fetched
		}
		resp = append(resp, response.RrIdPGroup{ID: g.ID, UUID: g.UUID, Name: g.Name, Path: g.Path})
	}
	c.JSON(http.StatusOK, response.RrIdPGroups{Code: "SUCCESS", Message: "ok", Groups: resp})
}

// GetGroup returns a single group — only accessible if the caller is a member.
// GET /v1/internal/group?id=...
func (rcvr *groupInternal) GetGroup(c *gin.Context) {
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
	g, err := rcvr.groupUsecase.GetGroup(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "IDP_GROUP_GET_404", "message": "Group not found"})
		return
	}
	c.JSON(http.StatusOK, response.RrSingleIdPGroup{
		Code:    "SUCCESS",
		Message: "ok",
		Group:   &response.RrIdPGroup{ID: g.ID, UUID: g.UUID, Name: g.Name, Path: g.Path},
	})
}

// CreateGroup creates a new group in the IdP. Any authenticated user may create a group.
// POST /v1/internal/groups
func (rcvr *groupInternal) CreateGroup(c *gin.Context) {
	_, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	var req request.RrCreateGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_GROUP_CREATE_400", "message": "Invalid request body"})
		return
	}
	g, err := rcvr.groupUsecase.CreateGroup(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_CREATE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.RrSingleIdPGroup{
		Code:    "SUCCESS",
		Message: "created",
		Group:   &response.RrIdPGroup{ID: g.ID, UUID: g.UUID, Name: g.Name, Path: g.Path},
	})
}

// UpdateGroup updates a group — only accessible if the caller is a member.
// PUT /v1/internal/groups/:id
func (rcvr *groupInternal) UpdateGroup(c *gin.Context) {
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
	role, err := callerGroupRole(c, rcvr.memberUsecase, groupID, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_UPDATE_001", "message": err.Error()})
		return
	}
	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_GROUP_UPDATE_403", "message": "Only group owners can update a group"})
		return
	}
	var req request.RrUpdateGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_GROUP_UPDATE_400", "message": "Invalid request body"})
		return
	}
	if err := rcvr.groupUsecase.UpdateGroup(c.Request.Context(), groupID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_UPDATE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "updated"})
}

// DeleteGroup deletes a group — only accessible if the caller is a member.
// DELETE /v1/internal/groups/:id
func (rcvr *groupInternal) DeleteGroup(c *gin.Context) {
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
	delRole, err := callerGroupRole(c, rcvr.memberUsecase, groupID, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_DELETE_001", "message": err.Error()})
		return
	}
	if delRole != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_GROUP_DELETE_403", "message": "Only group owners can delete a group"})
		return
	}
	if err := rcvr.groupUsecase.DeleteGroup(c.Request.Context(), groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_DELETE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "deleted"})
}
