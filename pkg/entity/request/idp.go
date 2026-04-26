package request

// CreateUser is the request body for creating a new user in the IdP.
type CreateUser struct {
	Username  string `json:"username"  binding:"required"`
	Email     string `json:"email"     binding:"required,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// Password is set as a temporary credential; the user must change it on first login.
	Password string `json:"password" binding:"required"`
	Enabled  bool   `json:"enabled"`
}

// UpdateUser is the request body for updating an existing user in the IdP.
// All fields are optional; nil / zero-value fields are ignored.
type UpdateUser struct {
	Email     *string `json:"email"      binding:"omitempty,email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Enabled   *bool   `json:"enabled"`
}

// CreateGroup is the request body for creating a new group in the IdP.
type CreateGroup struct {
	Name string `json:"name" binding:"required"`
}

// UpdateGroup is the request body for updating an existing group in the IdP.
type UpdateGroup struct {
	Name string `json:"name" binding:"required"`
}

// AddGroupMember is the request body for adding a user to a group.
type AddGroupMember struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role"    binding:"required,oneof=owner editor viewer"`
}

// RemoveGroupMember is the request body for removing a user from a group.
type RemoveGroupMember struct {
	UserID string `json:"user_id" binding:"required"`
}
