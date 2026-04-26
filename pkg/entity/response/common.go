package response

import (
	"time"
)

type Commons struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Commons []Common `json:"commons,omitempty"`
}

type Common struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// ValidateToken represents token validation response
type ValidateToken struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// User is the API representation of a local DB user.
type User struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Users wraps a list of local DB users.
type Users struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Users   []User `json:"users,omitempty"`
}

// Group is the API representation of a local DB group.
type Group struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// Groups wraps a list of local DB groups.
type Groups struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Groups  []Group `json:"groups,omitempty"`
}

// Member is the API representation of a local DB group membership.
type Member struct {
	UUID      string `json:"uuid"`
	UserUUID  string `json:"user_uuid"`
	GroupUUID string `json:"group_uuid"`
}

// Members wraps a list of local DB group memberships.
type Members struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Members []Member `json:"members,omitempty"`
}

// Login represents a login response.
type Login struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

