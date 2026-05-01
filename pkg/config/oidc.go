package config

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OIDCClaims represents the normalized identity attributes extracted from an OIDC id_token.
type OIDCClaims struct {
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

// OIDCProvider is the interface for OIDC-based identity provider integrations.
type OIDCProvider interface {
	// Name returns a human-readable identifier for the provider (e.g. "keycloak").
	Name() string

	// LoginURL returns the URL to which the user's browser should be redirected
	// to begin the authentication flow.
	// state is an opaque value used to prevent CSRF.
	LoginURL(state string) (string, error)

	// HandleCallback processes the IdP callback request (OIDC code exchange)
	// and returns the normalized claims.
	// The caller is responsible for verifying state before calling this method.
	HandleCallback(ctx context.Context, r *http.Request) (*OIDCClaims, error)
}

// oidcProvider implements OIDCProvider.
type oidcProvider struct {
	cfg      OIDCConfig
	provider *gooidc.Provider
	oauth2   oauth2.Config
	verifier *gooidc.IDTokenVerifier
}

// NewOIDCProvider creates a new OIDC provider, performing provider discovery.
// Call this once during server startup.
func NewOIDCProvider(ctx context.Context, cfg OIDCConfig) (OIDCProvider, error) {
	if cfg.IssuerURL == "" {
		return nil, errors.New("auth: OIDC issuer URL must not be empty")
	}
	if cfg.ClientID == "" {
		return nil, errors.New("auth: OIDC client ID must not be empty")
	}

	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{gooidc.ScopeOpenID, "email", "profile"}
	}

	// ProviderURL allows overriding the discovery endpoint (e.g. internal Docker URL)
	// while keeping IssuerURL as the token issuer claim.
	discoveryURL := cfg.IssuerURL
	if cfg.ProviderURL != "" {
		discoveryURL = cfg.ProviderURL
	}

	provider, err := gooidc.NewProvider(ctx, discoveryURL)
	if err != nil {
		return nil, fmt.Errorf("auth: OIDC discovery failed for %s: %w", discoveryURL, err)
	}

	oa2cfg := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	verifier := provider.Verifier(&gooidc.Config{ClientID: cfg.ClientID})

	return &oidcProvider{
		cfg:      cfg,
		provider: provider,
		oauth2:   oa2cfg,
		verifier: verifier,
	}, nil
}

func (rcvr *oidcProvider) Name() string { return rcvr.cfg.ProviderName }

// LoginURL builds the OAuth2 authorization URL with PKCE-compatible state.
func (rcvr *oidcProvider) LoginURL(state string) (string, error) {
	url := rcvr.oauth2.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return url, nil
}

// HandleCallback exchanges the authorization code for tokens and returns normalized OIDCClaims.
func (rcvr *oidcProvider) HandleCallback(ctx context.Context, r *http.Request) (*OIDCClaims, error) {
	code := r.FormValue("code")
	if code == "" {
		return nil, errors.New("auth: OIDC callback missing code parameter")
	}

	oauth2Token, err := rcvr.oauth2.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("auth: OIDC code exchange failed: %w", err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("auth: OIDC token response missing id_token")
	}

	idToken, err := rcvr.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("auth: OIDC id_token verification failed: %w", err)
	}

	var idClaims struct {
		Sub    string   `json:"sub"`
		Email  string   `json:"email"`
		Name   string   `json:"name"`
		Groups []string `json:"groups"`
	}
	if err := idToken.Claims(&idClaims); err != nil {
		return nil, fmt.Errorf("auth: failed to extract OIDC claims: %w", err)
	}

	return &OIDCClaims{
		Subject: idClaims.Sub,
		Email:   idClaims.Email,
		Name:    idClaims.Name,
		Groups:  idClaims.Groups,
	}, nil
}
