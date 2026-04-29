package request

// RrCreateUser is the request body for creating a new user in the IdP.
type RrCreateUser struct {
	Username  string `json:"username"  binding:"required"`
	Email     string `json:"email"     binding:"required,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// Password is set as a temporary credential; the user must change it on first login.
	Password string `json:"password" binding:"required"`
	Enabled  bool   `json:"enabled"`
}

// RrUpdateUser is the request body for updating an existing user in the IdP.
// All fields are optional; nil / zero-value fields are ignored.
type RrUpdateUser struct {
	Email     *string `json:"email"      binding:"omitempty,email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Enabled   *bool   `json:"enabled"`
}
