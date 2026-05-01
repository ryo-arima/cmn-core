package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// ResourceInternal handles resource endpoints for authenticated (non-admin) users.
type ResourceInternal interface {
	ListResources(c *gin.Context)
	GetResource(c *gin.Context)
	CreateResource(c *gin.Context)
	UpdateResource(c *gin.Context)
	DeleteResource(c *gin.Context)
	// Group-role management
	GetResourceGroupRoles(c *gin.Context)
	SetResourceGroupRole(c *gin.Context)
	DeleteResourceGroupRole(c *gin.Context)
}

type resourceInternal struct {
	resourceUsecase usecase.Resource
}

// NewResourceInternal creates a new ResourceInternal controller.
func NewResourceInternal(ru usecase.Resource) ResourceInternal {
	return &resourceInternal{resourceUsecase: ru}
}

// ListResources lists resources accessible to the caller.
// GET /v1/internal/resources
func (rcvr *resourceInternal) ListResources(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	resources, err := rcvr.resourceUsecase.ListResources(c.Request.Context(), claims.UUID, claims.Groups)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_LIST_001", "message": err.Error()})
		return
	}
	resp := make([]response.RrResource, 0, len(resources))
	for _, r := range resources {
		resp = append(resp, toResponseResource(r))
	}
	c.JSON(http.StatusOK, response.RrResources{Code: "SUCCESS", Message: "ok", Resources: resp})
}

// GetResource returns a single resource by UUID.
// GET /v1/internal/resource?uuid=...
func (rcvr *resourceInternal) GetResource(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	res, err := rcvr.resourceUsecase.GetResource(c.Request.Context(), c.Query("uuid"), claims.UUID, claims.Groups, false)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"code": "RESOURCE_GET_403", "message": "Access denied"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"code": "RESOURCE_GET_404", "message": "Resource not found"})
		return
	}
	c.JSON(http.StatusOK, response.RrSingleResource{Code: "SUCCESS", Message: "ok", Resource: ptr(toResponseResource(*res))})
}

// CreateResource creates a new resource.
// POST /v1/internal/resources
func (rcvr *resourceInternal) CreateResource(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	var req request.RrCreateResource
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "RESOURCE_CREATE_001", "message": "Invalid request body"})
		return
	}
	res, err := rcvr.resourceUsecase.CreateResource(c.Request.Context(), req.Name, req.Description, claims.UUID, req.OwnerGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_CREATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.RrSingleResource{Code: "SUCCESS", Message: "created", Resource: ptr(toResponseResource(*res))})
}

// UpdateResource updates an existing resource.
// PUT /v1/internal/resources/:uuid
func (rcvr *resourceInternal) UpdateResource(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	var req request.RrUpdateResource
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "RESOURCE_UPDATE_001", "message": "Invalid request body"})
		return
	}
	res, err := rcvr.resourceUsecase.UpdateResource(c.Request.Context(), c.Param("uuid"), req.Name, req.Description, claims.UUID, claims.Groups, false)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"code": "RESOURCE_UPDATE_403", "message": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_UPDATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.RrSingleResource{Code: "SUCCESS", Message: "updated", Resource: ptr(toResponseResource(*res))})
}

// DeleteResource soft-deletes a resource.
// DELETE /v1/internal/resources/:uuid
func (rcvr *resourceInternal) DeleteResource(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	if err := rcvr.resourceUsecase.DeleteResource(c.Request.Context(), c.Param("uuid"), claims.UUID, claims.Groups, false); err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"code": "RESOURCE_DELETE_403", "message": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_DELETE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "deleted"})
}

// GetResourceGroupRoles lists the group-role entries for a resource.
// GET /v1/internal/resource/groups?uuid=...
func (rcvr *resourceInternal) GetResourceGroupRoles(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	roles, err := rcvr.resourceUsecase.GetGroupRoles(c.Request.Context(), c.Query("uuid"), claims.UUID, claims.Groups, false)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"code": "RESOURCE_GROUP_403", "message": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_GROUP_001", "message": err.Error()})
		return
	}
	resp := make([]response.RrResourceGroupRole, 0, len(roles))
	for _, r := range roles {
		resp = append(resp, response.RrResourceGroupRole{ResourceUUID: r.ResourceUUID, GroupID: r.GroupID, Role: r.Role})
	}
	c.JSON(http.StatusOK, response.RrResourceGroupRoles{Code: "SUCCESS", Message: "ok", Groups: resp})
}

// SetResourceGroupRole adds or updates a group-role entry.
// PUT /v1/internal/resources/:uuid/groups
func (rcvr *resourceInternal) SetResourceGroupRole(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	var req request.RrSetResourceGroupRole
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "RESOURCE_GROUP_SET_001", "message": "Invalid request body"})
		return
	}
	if err := rcvr.resourceUsecase.SetGroupRole(c.Request.Context(), c.Param("uuid"), req.GroupID, req.Role, claims.UUID, claims.Groups, false); err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"code": "RESOURCE_GROUP_SET_403", "message": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_GROUP_SET_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "group role set"})
}

// DeleteResourceGroupRole removes a group-role entry.
// DELETE /v1/internal/resources/:uuid/groups/:group_id
func (rcvr *resourceInternal) DeleteResourceGroupRole(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	if err := rcvr.resourceUsecase.DeleteGroupRole(c.Request.Context(), c.Param("uuid"), c.Param("group_id"), claims.UUID, claims.Groups, false); err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"code": "RESOURCE_GROUP_DEL_403", "message": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_GROUP_DEL_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "group role removed"})
}

// ---- helpers ---------------------------------------------------------------
// (toResponseResource and ptr are defined in resource_helper.go)
