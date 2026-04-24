package response_test

import (
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/entity/response"
	"github.com/stretchr/testify/assert"
)

func TestUserResponse_Structure(t *testing.T) {
	resp := response.Users{
		Code:    "SUCCESS",
		Message: "Operation successful",
		Users: []response.User{
			{
				UUID:  "user-uuid-123",
				Email: "test@example.com",
				Name:  "Test User",
			},
		},
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Operation successful", resp.Message)
	assert.Len(t, resp.Users, 1)
	assert.Equal(t, "user-uuid-123", resp.Users[0].UUID)
}

func TestGroupResponse_Structure(t *testing.T) {
	resp := response.Groups{
		Code:    "SUCCESS",
		Message: "Group operation successful",
		Groups: []response.Group{
			{
				UUID: "group-uuid-123",
				Name: "Test Group",
			},
		},
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Group operation successful", resp.Message)
	assert.Len(t, resp.Groups, 1)
	assert.Equal(t, "group-uuid-123", resp.Groups[0].UUID)
}

func TestMemberResponse_Structure(t *testing.T) {
	resp := response.Members{
		Code:    "SUCCESS",
		Message: "Member operation successful",
		Members: []response.Member{
			{
				UUID:      "member-uuid-123",
				UserUUID:  "user-uuid-456",
				GroupUUID: "group-uuid-789",
			},
		},
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Member operation successful", resp.Message)
	assert.Len(t, resp.Members, 1)
	assert.Equal(t, "member-uuid-123", resp.Members[0].UUID)
}

func TestLoginResponse_Structure(t *testing.T) {
	resp := response.Login{
		Code:    "SUCCESS",
		Message: "Login successful",
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Login successful", resp.Message)
}

func TestCountResponse_Structure(t *testing.T) {
	// Commons structure doesn't have Count field, using basic structure test
	resp := response.Commons{
		Code:    "SUCCESS",
		Message: "Count retrieved",
	}

	assert.Equal(t, "SUCCESS", resp.Code)
	assert.Equal(t, "Count retrieved", resp.Message)
}
