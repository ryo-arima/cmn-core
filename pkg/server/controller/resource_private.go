package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// ResourcePrivate handles resource endpoints that require the admin role.
type ResourcePrivate interface {
	ListAllResources(c *gin.Context)
	GetResource(c *gin.Context)
	CreateResource(c *gin.Context)
	UpdateResource(c *gin.Context)
	DeleteResource(c *gin.Context)
	GetResourceGroupRoles(c *gin.Context)
	SetResourceGroupRole(c *gin.Context)
	DeleteResourceGroupRole(c *gin.Context)
}

type resourcePrivate struct {
	resourceUsecase usecase.Resource
}

// NewResourcePrivate creates a new ResourcePrivate controller.
func NewResourcePrivate(ru usecase.Resource) ResourcePrivate {
	return &resourcePrivate{resourceUsecase: ru}
}

// ListAllResources returns every non-deleted resource.
// GET /v1/private/resources
func (rcvr *resourcePrivate) ListAllResources(c *gin.Context) {
	resources, err := rcvr.resourceUsecase.ListAllResources(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_ADMIN_LIST_001", "message": err.Error()})
		return
	}
	resp := make([]response.RrResource, 0, len(resources))
	for _, r := range resources {
		resp = append(resp, toResponseResource(r))
	}
	c.JSON(http.StatusOK, response.RrResources{Code: "SUCCESS", Message: "ok", Resources: resp})
}

// DeleteResource soft-deletes any resource (admin override).
// DELETE /v1/private/resources/:uuid
func (rcvr *resourcePrivate) DeleteResource(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	if err := rcvr.resourceUsecase.AdminDeleteResource(c.Request.Context(), c.Param("uuid"), claims.UUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_ADMIN_DELETE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "deleted"})
}

// GetResource returns a single resource (admin access bypasses ownership checks).
// GET /v1/private/resource?uuid=...
func (rcvr *resourcePrivate) GetResource(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	res, err := rcvr.resourceUsecase.GetResource(c.Request.Context(), c.Query("uuid"), claims.UUID, claims.Groups, true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "RESOURCE_GET_404", "message": "Resource not found"})
		return
	}
	c.JSON(http.StatusOK, response.RrSingleResource{Code: "SUCCESS", Message: "ok", Resource: ptr(toResponseResource(*res))})
}

// CreateResource creates a new resource owned by the admin.
// POST /v1/private/resources
func (rcvr *resourcePrivate) CreateResource(c *gin.Context) {
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

// UpdateResource updates any resource (admin override).
// PUT /v1/private/resources/:uuid
func (rcvr *resourcePrivate) UpdateResource(c *gin.Context) {
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
	res, err := rcvr.resourceUsecase.UpdateResource(c.Request.Context(), c.Param("uuid"), req.Name, req.Description, claims.UUID, claims.Groups, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_UPDATE_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.RrSingleResource{Code: "SUCCESS", Message: "updated", Resource: ptr(toResponseResource(*res))})
}

// GetResourceGroupRoles lists group-role entries (admin access).
// GET /v1/private/resource/groups?uuid=...
func (rcvr *resourcePrivate) GetResourceGroupRoles(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	roles, err := rcvr.resourceUsecase.GetGroupRoles(c.Request.Context(), c.Param("uuid"), claims.UUID, claims.Groups, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_GROUP_001", "message": err.Error()})
		return
	}
	resp := make([]response.RrResourceGroupRole, 0, len(roles))
	for _, r := range roles {
		resp = append(resp, response.RrResourceGroupRole{ResourceUUID: r.ResourceUUID, GroupID: r.GroupID, Role: r.Role})
	}
	c.JSON(http.StatusOK, response.RrResourceGroupRoles{Code: "SUCCESS", Message: "ok", Groups: resp})
}

// SetResourceGroupRole adds or updates a group-role entry (admin override).
// PUT /v1/private/resources/:uuid/groups
func (rcvr *resourcePrivate) SetResourceGroupRole(c *gin.Context) {
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
	if err := rcvr.resourceUsecase.SetGroupRole(c.Request.Context(), c.Param("uuid"), req.GroupID, req.Role, claims.UUID, claims.Groups, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_GROUP_SET_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "group role set"})
}

// DeleteResourceGroupRole removes a group-role entry (admin override).
// DELETE /v1/private/resources/:uuid/groups  (group_uuid specified in request body)
func (rcvr *resourcePrivate) DeleteResourceGroupRole(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "RESOURCE_AUTH_001", "message": "Unauthorized"})
		return
	}
	var req request.RrDeleteResourceGroupRole
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "RESOURCE_GROUP_DEL_400", "message": "Invalid request body"})
		return
	}
	if err := rcvr.resourceUsecase.DeleteGroupRole(c.Request.Context(), c.Param("uuid"), req.GroupID, claims.UUID, claims.Groups, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_GROUP_DEL_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "group role deleted"})
}
