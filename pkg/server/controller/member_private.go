package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// MemberPrivate handles group membership admin endpoints. All routes require admin role.
type MemberPrivate interface {
	ListGroupMembers(c *gin.Context)
	AddGroupMember(c *gin.Context)
	RemoveGroupMember(c *gin.Context)
}

type memberPrivate struct {
	memberUsecase usecase.Member
}

// NewMemberPrivate creates a new MemberPrivate controller.
func NewMemberPrivate(mu usecase.Member) MemberPrivate {
	return &memberPrivate{memberUsecase: mu}
}

// ListGroupMembers lists members of a group.
// GET /v1/private/members?group_id=...
func (ic *memberPrivate) ListGroupMembers(c *gin.Context) {
	members, err := ic.memberUsecase.ListGroupMembers(c.Request.Context(), c.Query("group_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_LIST_001", "message": err.Error()})
		return
	}
	resp := make([]response.RrIdPUser, 0, len(members))
	for _, u := range members {
		resp = append(resp, response.RrIdPUser{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Enabled:   u.Enabled,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, response.RrIdPUsers{Code: "SUCCESS", Message: "ok", Users: resp})
}

// AddGroupMember adds a user to a group.
// POST /v1/private/member/:group_id
func (ic *memberPrivate) AddGroupMember(c *gin.Context) {
	var req request.RrAddGroupMember
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_MEMBER_ADD_001", "message": "Invalid request body"})
		return
	}
	if err := ic.memberUsecase.AddUserToGroup(c.Request.Context(), req.UserID, c.Param("group_id"), req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_ADD_002", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "member added"})
}

// RemoveGroupMember removes a user from a group.
// DELETE /v1/private/member/:group_id
func (ic *memberPrivate) RemoveGroupMember(c *gin.Context) {
	var req request.RrRemoveGroupMember
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_MEMBER_REMOVE_400", "message": "Invalid request body"})
		return
	}
	if err := ic.memberUsecase.RemoveUserFromGroup(c.Request.Context(), req.UserID, c.Param("group_id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_REMOVE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "member removed"})
}
