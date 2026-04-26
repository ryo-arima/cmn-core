package response

import "time"

// IdPUser is the API representation of a user managed by the identity provider.
type IdPUser struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Enabled   bool      `json:"enabled"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// IdPUsers wraps a list of IdP users.
type IdPUsers struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Users   []IdPUser `json:"users,omitempty"`
}

// SingleIdPUser wraps a single IdP user.
type SingleIdPUser struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	User    *IdPUser `json:"user,omitempty"`
}

// IdPGroup is the API representation of a group managed by the identity provider.
type IdPGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path,omitempty"`
}

// IdPGroups wraps a list of IdP groups.
type IdPGroups struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Groups  []IdPGroup `json:"groups,omitempty"`
}

// SingleIdPGroup wraps a single IdP group.
type SingleIdPGroup struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Group   *IdPGroup `json:"group,omitempty"`
}
