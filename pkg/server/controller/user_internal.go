package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// UserInternal handles user endpoints for authenticated (non-admin) users.
type UserInternal interface {
	GetMyUser(c *gin.Context)
	UpdateMyUser(c *gin.Context)
	ListGroupUsers(c *gin.Context)
}

type userInternal struct {
	userUsecase   usecase.User
	memberUsecase usecase.Member
	commonUsecase usecase.Common
}

// NewUserInternal creates a new UserInternal controller.
func NewUserInternal(uu usecase.User, mu usecase.Member, cu usecase.Common) UserInternal {
	return &userInternal{userUsecase: uu, memberUsecase: mu, commonUsecase: cu}
}

// GetMyUser returns a user's profile.
// If the query param ?id= is provided, returns that user.
// Otherwise returns the authenticated user's own profile.
// GET /v1/internal/user
func (ic *userInternal) GetMyUser(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	userID := c.Query("id")
	if userID == "" {
		userID = claims.UUID
	}
	u, err := ic.userUsecase.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "IDP_USER_GET_404", "message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, response.RrSingleIdPUser{
		Code:    "SUCCESS",
		Message: "ok",
		User: &response.RrIdPUser{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Enabled:   u.Enabled,
			Role:      ic.commonUsecase.ResolveRole(u.Email),
			CreatedAt: u.CreatedAt,
		},
	})
}

// ListGroupUsers returns all users who are members of any group the caller belongs to.
// Results are deduplicated. Group membership is read from the caller's JWT claims.
// GET /v1/internal/users
func (ic *userInternal) ListGroupUsers(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	seen := make(map[string]struct{})
	var users []response.RrIdPUser
	for _, gid := range claims.Groups {
		members, err := ic.memberUsecase.ListGroupMembers(c.Request.Context(), groupName(gid))
		if err != nil {
			continue
		}
		for _, u := range members {
			if _, dup := seen[u.ID]; dup {
				continue
			}
			seen[u.ID] = struct{}{}
			users = append(users, response.RrIdPUser{
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
	}
	c.JSON(http.StatusOK, response.RrIdPUsers{Code: "SUCCESS", Message: "ok", Users: users})
}

// UpdateMyUser updates the authenticated user's own profile in the IdP.
// PUT /v1/internal/user
func (ic *userInternal) UpdateMyUser(c *gin.Context) {
	claims, ok := share.GetUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "IDP_AUTH_001", "message": "Unauthorized"})
		return
	}
	var req request.RrUpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "IDP_USER_UPDATE_400", "message": "Invalid request body"})
		return
	}
	if err := ic.userUsecase.UpdateUser(c.Request.Context(), claims.UUID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "IDP_USER_UPDATE_001", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "updated"})
}
