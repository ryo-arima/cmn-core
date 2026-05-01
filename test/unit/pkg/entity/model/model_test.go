package model_test

import (
	"testing"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/stretchr/testify/assert"
)

func TestUsers_Structure(t *testing.T) {
	now := time.Now()
	user := model.PgUsers{
		ID:        1,
		UUID:      "test-uuid-123",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Name:      "Test User",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "test-uuid-123", user.UUID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, "Test User", user.Name)
	assert.NotNil(t, user.CreatedAt)
	assert.NotNil(t, user.UpdatedAt)
	assert.Nil(t, user.DeletedAt)
}

func TestGroups_Structure(t *testing.T) {
	now := time.Now()
	group := model.PgGroups{
		ID:        1,
		UUID:      "group-uuid-123",
		Name:      "Test Group",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	assert.Equal(t, uint(1), group.ID)
	assert.Equal(t, "group-uuid-123", group.UUID)
	assert.Equal(t, "Test Group", group.Name)
	assert.NotNil(t, group.CreatedAt)
	assert.NotNil(t, group.UpdatedAt)
	assert.Nil(t, group.DeletedAt)
}

func TestMembers_Structure(t *testing.T) {
	now := time.Now()
	member := model.PgMembers{
		ID:        1,
		UUID:      "member-uuid-123",
		UserUUID:  "user-uuid-456",
		GroupUUID: "group-uuid-789",
		Role:      "member",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	assert.Equal(t, uint(1), member.ID)
	assert.Equal(t, "member-uuid-123", member.UUID)
	assert.Equal(t, "user-uuid-456", member.UserUUID)
	assert.Equal(t, "group-uuid-789", member.GroupUUID)
	assert.Equal(t, "member", member.Role)
	assert.NotNil(t, member.CreatedAt)
	assert.NotNil(t, member.UpdatedAt)
	assert.Nil(t, member.DeletedAt)
}
