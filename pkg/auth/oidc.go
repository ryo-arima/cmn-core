package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OIDCConfig holds the configuration required to connect to an OIDC provider
// such as Keycloak or Casdoor.
type OIDCConfig struct {
	// ProviderName is a human-readable label, e.g. "keycloak" or "casdoor".
	ProviderName string `yaml:"provider_name"`
	// IssuerURL is the OIDC issuer URL (used for discovery).
	// Keycloak: https://<host>/realms/<realm>
	// Casdoor:  http://<host>
	IssuerURL string `yaml:"issuer_url"`
	// ClientID is the OIDC client ID registered in the IdP.
	ClientID string `yaml:"client_id"`
	// ClientSecret is the OIDC client secret.
	ClientSecret string `yaml:"client_secret"`
	// RedirectURL is the callback URL registered in the IdP.
	// Example: https://example.com/v1/share/auth/oidc/callback
	RedirectURL string `yaml:"redirect_url"`
	// Scopes is the list of OAuth2 scopes to request.
	// Defaults to ["openid", "email", "profile"].
	Scopes []string `yaml:"scopes"`
}

// oidcProvider implements Provider for OIDC-based identity providers.
type oidcProvider struct {
	cfg      OIDCConfig
	provider *gooidc.Provider
	oauth2   oauth2.Config
	verifier *gooidc.IDTokenVerifier
}

// NewOIDCProvider creates a new OIDC provider, performing provider discovery.
// Call this once during server startup.
func NewOIDCProvider(ctx context.Context, cfg OIDCConfig) (Provider, error) {
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

	provider, err := gooidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("auth: OIDC discovery failed for %s: %w", cfg.IssuerURL, err)
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

func (p *oidcProvider) Name() string { return p.cfg.ProviderName }

// LoginURL builds the OAuth2 authorization URL with PKCE-compatible state.
func (p *oidcProvider) LoginURL(state string) (string, error) {
	url := p.oauth2.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return url, nil
}

// HandleCallback exchanges the authorization code for tokens and returns normalized Claims.
// The request must contain "code" and "state" query parameters.
// The caller is responsible for verifying state before calling this method.
func (p *oidcProvider) HandleCallback(ctx context.Context, r *http.Request) (*Claims, error) {
	code := r.FormValue("code")
	if code == "" {
		return nil, errors.New("auth: OIDC callback missing code parameter")
	}

	oauth2Token, err := p.oauth2.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("auth: OIDC code exchange failed: %w", err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("auth: OIDC token response missing id_token")
	}

	idToken, err := p.verifier.Verify(ctx, rawIDToken)
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

	return &Claims{
		Subject: idClaims.Sub,
		Email:   idClaims.Email,
		Name:    idClaims.Name,
		Groups:  idClaims.Groups,
	}, nil
}

