package request_test

import (
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/stretchr/testify/assert"
)

func TestUserRequest_Structure(t *testing.T) {
	req := request.User{
		UUID:     "user-uuid-123",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	assert.Equal(t, "user-uuid-123", req.UUID)
	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
	assert.Equal(t, "Test User", req.Name)
}

func TestGroupRequest_Structure(t *testing.T) {
	req := request.Group{
		UUID: "group-uuid-123",
		Name: "Test Group",
	}

	assert.Equal(t, "group-uuid-123", req.UUID)
	assert.Equal(t, "Test Group", req.Name)
}

func TestMemberRequest_Structure(t *testing.T) {
	req := request.Member{
		UUID:      "member-uuid-123",
		UserUUID:  "user-uuid-456",
		GroupUUID: "group-uuid-789",
	}

	assert.Equal(t, "member-uuid-123", req.UUID)
	assert.Equal(t, "user-uuid-456", req.UserUUID)
	assert.Equal(t, "group-uuid-789", req.GroupUUID)
}

func TestLoginRequest_Structure(t *testing.T) {
	req := request.Login{
		Email:    "user@example.com",
		Password: "securepassword",
	}

	assert.Equal(t, "user@example.com", req.Email)
	assert.Equal(t, "securepassword", req.Password)
}
