// Package auth provides client-side authentication token lifecycle management.
package auth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/config"
)

// Manager handles client-side authentication token lifecycle.
// It is created once per process and shared across repository/usecase/controller layers.
type Manager struct {
	conf          config.BaseConfig
	profile       string // "admin", "app", or "anonymous"
	explicitToken string // when set, Token() always returns this value verbatim
}

// NewManager creates a Manager for the given profile.
// profile must be "admin", "app", or "anonymous".
func NewManager(conf config.BaseConfig, profile string) *Manager {
	return &Manager{conf: conf, profile: profile}
}

// WithToken returns a new Manager that always returns the given token without
// any file I/O, refresh, or SSO logic.  Useful for the anonymous client when
// validating an externally-provided token.
func (m *Manager) WithToken(token string) *Manager {
	return &Manager{conf: m.conf, profile: m.profile, explicitToken: token}
}

// Conf returns the BaseConfig embedded in this manager.
func (m *Manager) Conf() config.BaseConfig { return m.conf }

// IsAnonymous reports whether this manager is for unauthenticated access.
func (m *Manager) IsAnonymous() bool { return m.profile == "anonymous" }

// tokenDir returns the directory where token files are stored for this profile.
func (m *Manager) tokenDir() string {
	return filepath.Join("etc", ".cmn", "client", m.profile)
}

func (m *Manager) readFile(name string) string {
	b, err := os.ReadFile(filepath.Join(m.tokenDir(), name))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// Token returns a valid access token.
//
// If an explicit token was set via WithToken, it is returned verbatim.
// For anonymous profiles, "" is returned.
// Otherwise the token is loaded from the on-disk cache; if missing or expired
// it is obtained automatically via password-based login (POST /v1/public/login).
func (m *Manager) Token() (string, error) {
	if m.explicitToken != "" {
		return m.explicitToken, nil
	}
	if m.IsAnonymous() {
		return "", nil
	}

	token := m.readFile("access_token")
	if token != "" && !isTokenExpired(token) {
		return token, nil
	}

	// Token missing or expired — re-authenticate transparently.
	return m.loginWithPassword(context.Background())
}

// loginWithPassword authenticates using the configured credentials by calling
// POST /v1/public/login on the server.  The server is responsible for all IdP
// communication; the client never contacts Casdoor or Keycloak directly.
func (m *Manager) loginWithPassword(ctx context.Context) (string, error) {
	creds := m.conf.YamlConfig.Application.Client.Credentials
	email, password := creds.Email, creds.Password
	if email == "" || password == "" {
		return "", fmt.Errorf(
			"no credentials configured — set Application.Client.credentials.email/password in the config file",
		)
	}
	return m.loginViaServer(ctx, email, password)
}

// loginViaServer posts credentials to POST /v1/public/login on the app server
// and returns the access token contained in the response.
func (m *Manager) loginViaServer(ctx context.Context, email, password string) (string, error) {
	body, err := json.Marshal(map[string]string{"email": email, "password": password})
	if err != nil {
		return "", fmt.Errorf("build login body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		m.serverBase()+"/v1/public/login", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil || result.AccessToken == "" {
		return "", fmt.Errorf("parse login response: %w", err)
	}
	m.SaveTokenPair(result.AccessToken, "")
	return result.AccessToken, nil
}

// ForceRefresh re-authenticates using the configured credentials.
func (m *Manager) ForceRefresh() error {
	m.ClearTokens()
	_, err := m.loginWithPassword(context.Background())
	return err
}

// ForceLogin re-authenticates using the configured credentials (provider argument is ignored).
func (m *Manager) ForceLogin(_ string) error {
	return m.ForceRefresh()
}

// SaveTokenPair persists the access token to disk.
func (m *Manager) SaveTokenPair(access, _ string) {
	dir := m.tokenDir()
	_ = os.MkdirAll(dir, 0o755)
	if access != "" {
		_ = os.WriteFile(filepath.Join(dir, "access_token"), []byte(access), 0o600)
	}
}

// ClearTokens removes the stored access token file.
func (m *Manager) ClearTokens() {
	_ = os.Remove(filepath.Join(m.tokenDir(), "access_token"))
	_ = os.Unsetenv("CMN_ACCESS_TOKEN")
}

// HTTPClient returns an *http.Client that transparently injects the Bearer token
// into every outgoing request.  On HTTP 401 responses, it refreshes the token
// once and retries the request (for GET/HEAD/DELETE only, since those have no body).
func (m *Manager) HTTPClient() *http.Client {
	return &http.Client{
		Transport: &authTransport{manager: m, base: http.DefaultTransport},
	}
}

// serverBase returns the configured server endpoint URL.
func (m *Manager) serverBase() string {
	return m.conf.YamlConfig.Application.Client.ServerEndpoint
}

// isTokenExpired returns true if the JWT is expired or structurally invalid.
// It decodes the payload without verifying the signature (expiry check only).
func isTokenExpired(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return true
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return true
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return true
	}
	// Treat tokens expiring within the next 10 seconds as already expired.
	return time.Now().Unix() >= claims.Exp-10
}

// authTransport is an http.RoundTripper that injects the Bearer token into every request.
type authTransport struct {
	manager *Manager
	base    http.RoundTripper
}

// RoundTrip obtains a valid token from the Manager and sets the Authorization header.
// On HTTP 401, it forces a re-login and retries once (safe for bodyless methods only).
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.manager.Token()
	if err != nil {
		return nil, err
	}

	// Buffer the request body so that POST/PUT/PATCH can also be retried on 401.
	var bodyBuf []byte
	if req.Body != nil && req.Body != http.NoBody {
		bodyBuf, err = io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("auth transport: buffer request body: %w", err)
		}
	}

	doRequest := func(tok string) (*http.Response, error) {
		clone := req.Clone(req.Context())
		if tok != "" {
			clone.Header.Set("Authorization", "Bearer "+tok)
		}
		if bodyBuf != nil {
			clone.Body = io.NopCloser(bytes.NewReader(bodyBuf))
			clone.ContentLength = int64(len(bodyBuf))
		}
		return t.base.RoundTrip(clone)
	}

	resp, err := doRequest(token)
	if err != nil {
		return nil, err
	}

	// On 401, force re-login and retry once for any HTTP method.
	if resp.StatusCode == http.StatusUnauthorized && !t.manager.IsAnonymous() {
		resp.Body.Close()
		t.manager.ClearTokens()
		newToken, loginErr := t.manager.loginWithPassword(req.Context())
		if loginErr != nil {
			return nil, loginErr
		}
		return doRequest(newToken)
	}

	return resp, nil
}
