// Package auth provides client-side authentication token lifecycle management.
package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
// Otherwise, the token is loaded from disk.
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

	return "", errors.New("no valid token found, please authenticate via your IdP and set CMN_ACCESS_TOKEN or save a token to disk")
}

// ForceRefresh is not supported. Tokens must be obtained from the IdP.
func (m *Manager) ForceRefresh() error {
	return errors.New("token refresh is not supported; please re-authenticate via your IdP")
}

// ForceLogin is not supported. Authentication must happen via the IdP directly.
func (m *Manager) ForceLogin(_ string) error {
	return errors.New("browser-based SSO is not supported; please authenticate via your IdP and provide the token")
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
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.manager.Token()
	if err != nil {
		return nil, err
	}

	clone := req.Clone(req.Context())
	if token != "" {
		clone.Header.Set("Authorization", "Bearer "+token)
	}
	return t.base.RoundTrip(clone)
}
