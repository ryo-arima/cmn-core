package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// MemberInternal handles group membership endpoints for authenticated (non-admin) users.
type MemberInternal interface {
	ListGroupMembers(c *gin.Context)
	AddGroupMember(c *gin.Context)
	RemoveGroupMember(c *gin.Context)
}

type memberInternal struct {
	memberUsecase usecase.Member
	commonUsecase usecase.Common
}

// NewMemberInternal creates a new MemberInternal controller.
func NewMemberInternal(mu usecase.Member, cu usecase.Common) MemberInternal {
	return &memberInternal{memberUsecase: mu, commonUsecase: cu}
}

// ListGroupMembers lists members of a group — only accessible if the caller is a member.
// GET /v1/internal/members?group_id=...
func (ic *memberInternal) ListGroupMembers(c *gin.Context) {
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
	members, err := ic.memberUsecase.ListGroupMembers(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_LIST_001", "message": err.Error()})
		return
	}
	resp := make([]response.RrIdPUser, 0, len(members))
	for _, u := range members {
		resp = append(resp, response.RrIdPUser{
			ID:        u.ID,
			UUID:      u.UUID,
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

// AddGroupMember adds a user to a group — only group owners may do this.
// POST /v1/internal/member/:group_id
func (ic *memberInternal) AddGroupMember(c *gin.Context) {
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
	role, err := callerGroupRole(c, ic.memberUsecase, groupID, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_ADD_001", "message": err.Error()})
		return
	}
	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_MEMBER_ADD_403", "message": "Only group owners can add members"})
		return
	}
	var req request.RrAddGroupMember
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_MEMBER_ADD_400", "message": "Invalid request body"})
		return
	}
	if err := ic.memberUsecase.AddUserToGroup(c.Request.Context(), req.UserID, groupID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_ADD_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "member added"})
}

// RemoveGroupMember removes a user from a group — only group owners may do this.
// DELETE /v1/internal/member/:group_id
func (ic *memberInternal) RemoveGroupMember(c *gin.Context) {
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
	role, err := callerGroupRole(c, ic.memberUsecase, groupID, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_REMOVE_001", "message": err.Error()})
		return
	}
	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"code": "IDP_MEMBER_REMOVE_403", "message": "Only group owners can remove members"})
		return
	}
	var req request.RrRemoveGroupMember
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_MEMBER_REMOVE_400", "message": "Invalid request body"})
		return
	}
	if err := ic.memberUsecase.RemoveUserFromGroup(c.Request.Context(), req.UserID, groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_MEMBER_REMOVE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "member removed"})
}
