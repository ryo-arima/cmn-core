// Package auth provides the interface and implementations for integrating with
// external identity providers (IdPs) such as Keycloak and Casdoor via
// OIDC (OpenID Connect) and SAML 2.0.
package auth

import (
	"context"
	"net/http"
)

// Claims represents the normalized identity attributes extracted from either
// an OIDC id_token or a SAML Assertion.
type Claims struct {
	// Subject is the unique identifier for the user within the IdP.
	Subject string
	// Email is the user's email address.
	Email string
	// Name is the display name of the user.
	Name string
	// Groups contains the groups / roles assigned to the user in the IdP.
	Groups []string
	// RawAttributes holds additional IdP-specific attributes.
	RawAttributes map[string][]string
}

// Provider is the common interface implemented by both the OIDC and SAML
// authentication providers.
type Provider interface {
	// Name returns a human-readable identifier for the provider (e.g. "keycloak-oidc").
	Name() string

	// LoginURL returns the URL to which the user's browser should be redirected
	// to begin the authentication flow.
	// state is an opaque value used to prevent CSRF; it will be echoed back in
	// the callback.
	LoginURL(state string) (string, error)

	// HandleCallback processes the IdP callback request (OIDC code exchange or
	// SAML ACS POST) and returns the normalized claims.
	HandleCallback(ctx context.Context, r *http.Request) (*Claims, error)
}
