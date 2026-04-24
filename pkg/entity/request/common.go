package request

import "time"

type Common struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// RefreshToken represents a request to refresh an access token.
type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// OIDCCallback represents the authorization code callback from an OIDC provider.
type OIDCCallback struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state"`
}

// SAMLCallback represents the ACS POST callback from a SAML identity provider.
type SAMLCallback struct {
	SAMLResponse string `form:"SAMLResponse" binding:"required"`
	RelayState   string `form:"RelayState"`
}
