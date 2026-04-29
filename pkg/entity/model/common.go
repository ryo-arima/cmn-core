package model

import (
	"encoding/json"
	"time"
)

// LoCommons is an internal base struct for local DB entities.
type LoCommons struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// LoJWTClaims represents claims extracted from an IdP-issued JWT token.
type LoJWTClaims struct {
	UUID      string   `json:"sub"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	Groups    []string `json:"groups,omitempty"`
	Role      string   `json:"role,omitempty"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
}

// KcTokenResponse is the Keycloak token endpoint response.
type KcTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// CdTokenResponse is the Casdoor token endpoint response.
type CdTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// CdResponse is the generic Casdoor API response wrapper.
type CdResponse struct {
	Status string          `json:"status"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
}
