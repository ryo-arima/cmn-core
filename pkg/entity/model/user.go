package model

import "time"

// PgUsers is the GORM/PostgreSQL model for the users table.
type PgUsers struct {
	ID        uint       `gorm:"primaryKey,autoIncrement"`
	UUID      string
	Email     string
	Password  string
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// LoUser represents a user record returned from an external identity provider.
// This is NOT a GORM model; it is never persisted to the local database.
type LoUser struct {
	// ID is the IdP-internal unique identifier.
	ID        string
	// UUID is the IdP-internal UUID (e.g. Keycloak's "id" field, Casdoor's "id" field).
	UUID      string
	Username  string
	Email     string
	FirstName string
	LastName  string
	Enabled   bool
	Role      string // resolved locally from admin emails list
	CreatedAt time.Time
}

// KcUser represents a Keycloak user JSON payload.
type KcUser struct {
	ID               string              `json:"id,omitempty"`
	Username         string              `json:"username,omitempty"`
	Email            string              `json:"email,omitempty"`
	FirstName        string              `json:"firstName,omitempty"`
	LastName         string              `json:"lastName,omitempty"`
	Enabled          bool                `json:"enabled"`
	EmailVerified    bool                `json:"emailVerified,omitempty"`
	CreatedTimestamp int64               `json:"createdTimestamp,omitempty"`
	Credentials      []KcCredential      `json:"credentials,omitempty"`
	Attributes       map[string][]string `json:"attributes,omitempty"`
}

// KcCredential represents a Keycloak credential JSON payload.
type KcCredential struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

// CdUser represents a Casdoor user JSON payload.
type CdUser struct {
	Owner        string            `json:"owner"`
	Name         string            `json:"name"`
	ID           string            `json:"id,omitempty"`
	Email        string            `json:"email,omitempty"`
	DisplayName  string            `json:"displayName,omitempty"`
	FirstName    string            `json:"firstName,omitempty"`
	LastName     string            `json:"lastName,omitempty"`
	IsAdmin      bool              `json:"isAdmin,omitempty"`
	IsForbidden  bool              `json:"isForbidden,omitempty"`
	CreatedTime  string            `json:"createdTime,omitempty"`
	Password     string            `json:"password,omitempty"`
	PasswordSalt string            `json:"passwordSalt,omitempty"`
	Groups       []string          `json:"groups,omitempty"`
	Properties   map[string]string `json:"properties,omitempty"`
}
