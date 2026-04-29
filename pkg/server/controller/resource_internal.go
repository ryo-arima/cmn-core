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
func (rc *resourceInternal) ListResources(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	resources, err := rc.resourceUsecase.ListResources(c.Request.Context(), claims.UUID, claims.Groups)
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
func (rc *resourceInternal) GetResource(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	res, err := rc.resourceUsecase.GetResource(c.Request.Context(), c.Query("uuid"), claims.UUID, claims.Groups, false)
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
func (rc *resourceInternal) CreateResource(c *gin.Context) {
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
	res, err := rc.resourceUsecase.CreateResource(c.Request.Context(), req.Name, req.Description, claims.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_CREATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.RrSingleResource{Code: "SUCCESS", Message: "created", Resource: ptr(toResponseResource(*res))})
}

// UpdateResource updates an existing resource.
// PUT /v1/internal/resources/:uuid
func (rc *resourceInternal) UpdateResource(c *gin.Context) {
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
	res, err := rc.resourceUsecase.UpdateResource(c.Request.Context(), c.Param("uuid"), req.Name, req.Description, claims.UUID, claims.Groups, false)
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
func (rc *resourceInternal) DeleteResource(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	if err := rc.resourceUsecase.DeleteResource(c.Request.Context(), c.Param("uuid"), claims.UUID, claims.Groups, false); err != nil {
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
func (rc *resourceInternal) GetResourceGroupRoles(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	roles, err := rc.resourceUsecase.GetGroupRoles(c.Request.Context(), c.Query("uuid"), claims.UUID, claims.Groups, false)
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
		resp = append(resp, response.RrResourceGroupRole{ResourceUUID: r.ResourceUUID, GroupUUID: r.GroupUUID, Role: r.Role})
	}
	c.JSON(http.StatusOK, response.RrResourceGroupRoles{Code: "SUCCESS", Message: "ok", Groups: resp})
}

// SetResourceGroupRole adds or updates a group-role entry.
// PUT /v1/internal/resources/:uuid/groups
func (rc *resourceInternal) SetResourceGroupRole(c *gin.Context) {
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
	if err := rc.resourceUsecase.SetGroupRole(c.Request.Context(), c.Param("uuid"), req.GroupUUID, req.Role, claims.UUID, claims.Groups, false); err != nil {
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
// DELETE /v1/internal/resources/:uuid/groups/:group_uuid
func (rc *resourceInternal) DeleteResourceGroupRole(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	if err := rc.resourceUsecase.DeleteGroupRole(c.Request.Context(), c.Param("uuid"), c.Param("group_uuid"), claims.UUID, claims.Groups, false); err != nil {
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
