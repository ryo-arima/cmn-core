package model

import "time"

// PgGroups is the GORM/PostgreSQL model for the groups table.
type PgGroups struct {
	ID        uint       `gorm:"primaryKey,autoIncrement"`
	UUID      string     `gorm:"uniqueIndex"`
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// LoGroup represents a group record returned from an external identity provider.
// This is NOT a GORM model; it is never persisted to the local database.
type LoGroup struct {
	// ID is the IdP-internal unique identifier.
	ID   string
	// UUID is the IdP-internal UUID (e.g. Keycloak's "id" field, Casdoor's "id" field).
	UUID string
	Name string
	// Path is the hierarchical path (Keycloak only; empty for Casdoor).
	Path string
}

// KcGroup represents a Keycloak group JSON payload.
type KcGroup struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

// CdGroup represents a Casdoor group JSON payload.
type CdGroup struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
	ID    string `json:"id,omitempty"`
}
