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
