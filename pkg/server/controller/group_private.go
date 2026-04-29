package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// GroupPrivate handles group admin endpoints. All routes require admin role.
type GroupPrivate interface {
	ListGroups(c *gin.Context)
	GetGroup(c *gin.Context)
	CreateGroup(c *gin.Context)
	UpdateGroup(c *gin.Context)
	DeleteGroup(c *gin.Context)
}

type groupPrivate struct {
	groupUsecase  usecase.Group
	commonUsecase usecase.Common
}

// NewGroupPrivate creates a new GroupPrivate controller.
func NewGroupPrivate(gu usecase.Group, cu usecase.Common) GroupPrivate {
	return &groupPrivate{groupUsecase: gu, commonUsecase: cu}
}

// ListGroups lists all groups from the IdP.
// GET /v1/private/groups
func (ic *groupPrivate) ListGroups(c *gin.Context) {
	groups, err := ic.groupUsecase.ListGroups(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_LIST_001", "message": err.Error()})
		return
	}
	resp := make([]response.RrIdPGroup, 0, len(groups))
	for _, g := range groups {
		resp = append(resp, response.RrIdPGroup{ID: g.ID, UUID: g.UUID, Name: g.Name, Path: g.Path})
	}
	c.JSON(http.StatusOK, response.RrIdPGroups{Code: "SUCCESS", Message: "ok", Groups: resp})
}

// GetGroup returns a single group from the IdP.
// GET /v1/private/group?id=...
func (ic *groupPrivate) GetGroup(c *gin.Context) {
	g, err := ic.groupUsecase.GetGroup(c.Request.Context(), c.Query("id"))
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

// CreateGroup creates a new group in the IdP.
// POST /v1/private/groups
func (ic *groupPrivate) CreateGroup(c *gin.Context) {
	var req request.RrCreateGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_GROUP_CREATE_001", "message": "Invalid request body"})
		return
	}
	g, err := ic.groupUsecase.CreateGroup(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_CREATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.RrSingleIdPGroup{
		Code:    "SUCCESS",
		Message: "created",
		Group:   &response.RrIdPGroup{ID: g.ID, UUID: g.UUID, Name: g.Name, Path: g.Path},
	})
}

// UpdateGroup updates an existing group in the IdP.
// PUT /v1/private/groups/:id
func (ic *groupPrivate) UpdateGroup(c *gin.Context) {
	var req request.RrUpdateGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_GROUP_UPDATE_001", "message": "Invalid request body"})
		return
	}
	if err := ic.groupUsecase.UpdateGroup(c.Request.Context(), c.Param("id"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_UPDATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "updated"})
}

// DeleteGroup deletes a group from the IdP.
// DELETE /v1/private/groups/:id
func (ic *groupPrivate) DeleteGroup(c *gin.Context) {
	if err := ic.groupUsecase.DeleteGroup(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_GROUP_DELETE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "deleted"})
}
