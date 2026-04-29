package response

import "time"

// RrIdPUser is the API representation of a user managed by the identity provider.
type RrIdPUser struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Enabled   bool      `json:"enabled"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// RrIdPUsers wraps a list of IdP users.
type RrIdPUsers struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Users   []RrIdPUser `json:"users,omitempty"`
}

// RrSingleIdPUser wraps a single IdP user.
type RrSingleIdPUser struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	User    *RrIdPUser `json:"user,omitempty"`
}
