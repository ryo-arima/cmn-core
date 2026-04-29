package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
)

// ---- keycloakManager -------------------------------------------------------

type keycloakManager struct {
	cfg    config.KeycloakConfig
	client *http.Client
	mu     sync.Mutex
	token  string
	expiry time.Time
}

func newKeycloakManager(cfg config.KeycloakConfig) IdPManager {
	return &keycloakManager{
		cfg:    cfg,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (m *keycloakManager) adminURL(path string) string {
	return fmt.Sprintf("%s/admin/realms/%s%s", m.cfg.BaseURL, m.cfg.Realm, path)
}

// getToken obtains (or returns a cached) admin access token via client_credentials.
func (m *keycloakManager) getToken(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.token != "" && time.Now().Before(m.expiry) {
		return m.token, nil
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", m.cfg.BaseURL, m.cfg.Realm)
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", m.cfg.AdminClientID)
	form.Set("client_secret", m.cfg.AdminClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("keycloak: build token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("keycloak: token request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("keycloak: token request status %d", resp.StatusCode)
	}

	var t model.KcTokenResponse
	if err := json.Unmarshal(body, &t); err != nil {
		return "", fmt.Errorf("keycloak: parse token response: %w", err)
	}
	m.token = t.AccessToken
	// Subtract 10 s to avoid using an almost-expired token.
	m.expiry = time.Now().Add(time.Duration(t.ExpiresIn-10) * time.Second)
	return m.token, nil
}

// do performs an authenticated request and returns status + response body.
func (m *keycloakManager) do(ctx context.Context, method, rawURL string, payload interface{}) (int, []byte, error) {
	token, err := m.getToken(ctx)
	if err != nil {
		return 0, nil, err
	}

	var bodyReader io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, fmt.Errorf("keycloak: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return 0, nil, fmt.Errorf("keycloak: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("keycloak: %s %s: %w", method, rawURL, err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b, nil
}

// kcExtractCreatedID retrieves the new resource ID from the Location header
// that Keycloak returns after a 201 Created response.
func kcExtractCreatedID(location string) string {
	parts := strings.Split(strings.TrimRight(location, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
