package model

import "time"

// IdPUser represents a user record returned from an external identity provider.
// This is NOT a GORM model; it is never persisted to the local database.
type IdPUser struct {
	// ID is the IdP-internal unique identifier.
	ID        string
	Username  string
	Email     string
	FirstName string
	LastName  string
	Enabled   bool
	Role      string // resolved locally from admin emails list
	CreatedAt time.Time
}

// IdPGroup represents a group record returned from an external identity provider.
// This is NOT a GORM model; it is never persisted to the local database.
type IdPGroup struct {
	// ID is the IdP-internal unique identifier.
	ID   string
	Name string
	// Path is the hierarchical path (Keycloak only; empty for Casdoor).
	Path string
}
