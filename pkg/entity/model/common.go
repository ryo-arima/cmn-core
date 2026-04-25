package model

import "time"

type Commons struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// JWTClaims represents claims extracted from an IdP-issued JWT token.
type JWTClaims struct {
	UUID      string   `json:"sub"`            // IdP subject identifier
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	Groups    []string `json:"groups,omitempty"`
	Role      string   `json:"role,omitempty"` // resolved locally from admin emails list
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
}
