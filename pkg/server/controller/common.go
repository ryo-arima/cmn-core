package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/ryo-arima/cmn-core/pkg/global"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"github.com/ryo-arima/cmn-core/pkg/server/usecase"
)

// ---- Logging helpers -------------------------------------------------------

// Local aliases for cleaner logging code - use functions to get logger dynamically
func INFO(requestID string, mcode global.MCode, message string) {
	if logger := share.GetServerLogger(); logger != nil {
		logger.INFO(requestID, mcode, message)
	}
}

func DEBUG(requestID string, mcode global.MCode, message string, fields ...map[string]interface{}) {
	if logger := share.GetServerLogger(); logger != nil {
		logger.DEBUG(requestID, mcode, message, fields...)
	}
}

func WARN(requestID string, mcode global.MCode, message string) {
	if logger := share.GetServerLogger(); logger != nil {
		logger.WARN(requestID, mcode, message)
	}
}

func ERROR(requestID string, mcode global.MCode, message string) {
	if logger := share.GetServerLogger(); logger != nil {
		logger.ERROR(requestID, mcode, message)
	}
}

// Local MCode definitions
var (
	SRNRSR1 = global.SRNRSR1
	SRNRSR2 = global.SRNRSR2
	Mcode   = global.Mcode
)

// ---- Resource helpers ------------------------------------------------------

func toResponseResource(r model.PgResource) response.RrResource {
	return response.RrResource{
		ID:          r.ID,
		UUID:        r.UUID,
		Name:        r.Name,
		Description: r.Description,
		CreatedBy:   r.CreatedBy,
		UpdatedBy:   r.UpdatedBy,
		DeletedBy:   r.DeletedBy,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}

func ptr[T any](v T) *T { return &v }

// ---- IdP access-control helpers --------------------------------------------

// isMemberOf returns true if groupID is present in the caller's groups claim.
// Handles the Casdoor JWT format where groups are "org/name" (e.g. "cmn/group001")
// and URL-param group IDs are just the name part (e.g. "group001").
func isMemberOf(groups []string, groupID string) bool {
	// Normalize target: strip org prefix if present
	normTarget := groupID
	if i := strings.LastIndex(groupID, "/"); i >= 0 {
		normTarget = groupID[i+1:]
	}
	for _, g := range groups {
		norm := g
		if i := strings.LastIndex(g, "/"); i >= 0 {
			norm = g[i+1:]
		}
		if norm == normTarget {
			return true
		}
	}
	return false
}

// groupName strips the org prefix from a Casdoor group ID.
// "cmn/group001" → "group001"; "group001" → "group001".
func groupName(gid string) string {
	if i := strings.LastIndex(gid, "/"); i >= 0 {
		return gid[i+1:]
	}
	return gid
}

// callerGroupRole returns the caller's role in the given group ("owner",
// "editor", "viewer") by fetching the member list from the IdP.
// Returns "" if the caller is not found in the group.
func callerGroupRole(c *gin.Context, mu usecase.Member, groupID, callerEmail string) (string, error) {
	members, err := mu.ListGroupMembers(c.Request.Context(), groupID)
	if err != nil {
		return "", err
	}
	for _, m := range members {
		if m.Email == callerEmail {
			return m.Role, nil
		}
	}
	return "", nil
}
