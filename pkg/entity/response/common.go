package response

import (
	"time"
)

type RrCommons struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Commons []RrCommon `json:"commons,omitempty"`
}

type RrCommon struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// RrValidateToken represents token validation response
type RrValidateToken struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// RrUser is the API representation of a local DB user.
type RrUser struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// RrUsers wraps a list of local DB users.
type RrUsers struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Users   []RrUser `json:"users,omitempty"`
}

// RrGroup is the API representation of a local DB group.
type RrGroup struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// RrGroups wraps a list of local DB groups.
type RrGroups struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Groups  []RrGroup `json:"groups,omitempty"`
}

// RrMember is the API representation of a local DB group membership.
type RrMember struct {
	UUID      string `json:"uuid"`
	UserUUID  string `json:"user_uuid"`
	GroupUUID string `json:"group_uuid"`
}

// RrMembers wraps a list of local DB group memberships.
type RrMembers struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Members []RrMember `json:"members,omitempty"`
}

// RrLogin represents a login response.
type RrLogin struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

