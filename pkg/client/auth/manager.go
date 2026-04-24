// Package auth provides client-side authentication token lifecycle management.
// The Manager transparently loads, refreshes, and renews tokens so that callers
// do not need to handle tokens explicitly.
package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
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
//   - If an explicit token was set via WithToken, it is returned verbatim.
//   - For anonymous profiles, "" is returned.
//   - Otherwise, the token is loaded from disk.  If it is missing or expired,
//     a refresh is attempted; if that also fails, an SSO login flow is triggered.
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

	// Token absent or expired – try to refresh silently.
	if err := m.ForceRefresh(); err == nil {
		token = m.readFile("access_token")
		if token != "" && !isTokenExpired(token) {
			return token, nil
		}
	}

	// Refresh unavailable – trigger SSO login.
	if err := m.ForceLogin(""); err != nil {
		return "", fmt.Errorf("authentication required: %w", err)
	}
	token = m.readFile("access_token")
	if token == "" {
		return "", fmt.Errorf("no token available after SSO login")
	}
	return token, nil
}

// ForceRefresh refreshes the access token using the stored refresh token.
// Returns an error if no refresh token is available or if the server rejects it.
func (m *Manager) ForceRefresh() error {
	rt := m.readFile("refresh_token")
	if rt == "" {
		return fmt.Errorf("no refresh token available")
	}
	pair, err := m.callRefresh(rt)
	if err != nil {
		return err
	}
	m.SaveTokenPair(pair.AccessToken, pair.RefreshToken)
	return nil
}

// ForceLogin triggers a fresh SSO login flow regardless of the current token state.
// provider overrides the configured provider; pass "" to use the value from app.yaml
// (defaults to "oidc" when not configured).
func (m *Manager) ForceLogin(provider string) error {
	if provider == "" {
		provider = m.conf.YamlConfig.Application.Client.Auth.Provider
	}
	if provider == "" {
		provider = "oidc"
	}

	loginURL, sessionID, err := m.callSSOStart(provider)
	if err != nil {
		return fmt.Errorf("SSO start: %w", err)
	}

	fmt.Printf("\nOpen the following URL in your browser to authenticate:\n\n  %s\n\n", loginURL)
	if !m.conf.YamlConfig.Application.Client.Auth.NoBrowser {
		openBrowser(loginURL)
	}

	fmt.Print("Waiting for authentication")
	pair, err := m.callSSOPoll(sessionID)
	if err != nil {
		fmt.Println()
		return fmt.Errorf("SSO poll: %w", err)
	}
	fmt.Println(" done.")

	m.SaveTokenPair(pair.AccessToken, pair.RefreshToken)
	return nil
}

// SaveTokenPair persists the access and refresh tokens to disk.
func (m *Manager) SaveTokenPair(access, refresh string) {
	dir := m.tokenDir()
	_ = os.MkdirAll(dir, 0o755)
	if access != "" {
		_ = os.WriteFile(filepath.Join(dir, "access_token"), []byte(access), 0o600)
	}
	if refresh != "" {
		_ = os.WriteFile(filepath.Join(dir, "refresh_token"), []byte(refresh), 0o600)
	}
}

// ClearTokens removes stored token files and clears environment variables.
func (m *Manager) ClearTokens() {
	_ = os.Remove(filepath.Join(m.tokenDir(), "access_token"))
	_ = os.Remove(filepath.Join(m.tokenDir(), "refresh_token"))
	_ = os.Unsetenv("CMN_ACCESS_TOKEN")
	_ = os.Unsetenv("CMN_REFRESH_TOKEN")
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

// callRefresh calls the server's token refresh endpoint.
// It uses a plain http.DefaultClient to avoid recursive token resolution.
func (m *Manager) callRefresh(refreshToken string) (*model.TokenPair, error) {
	endpoint := m.serverBase() + "/v1/share/token/refresh"
	body, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("create refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Code      string           `json:"code"`
		Message   string           `json:"message"`
		TokenPair *model.TokenPair `json:"token_pair"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode refresh response: %w", err)
	}
	if resp.StatusCode != http.StatusOK || result.TokenPair == nil {
		return nil, fmt.Errorf("token refresh failed: %s", result.Message)
	}
	return result.TokenPair, nil
}

// callSSOStart requests a login URL and session ID from the server.
func (m *Manager) callSSOStart(provider string) (loginURL, sessionID string, err error) {
	endpoint := fmt.Sprintf("%s/v1/share/auth/sso/start?provider=%s", m.serverBase(), provider)
	req, e := http.NewRequest("GET", endpoint, nil)
	if e != nil {
		return "", "", fmt.Errorf("create SSO start request: %w", e)
	}
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return "", "", fmt.Errorf("SSO start request failed: %w", e)
	}
	defer resp.Body.Close()

	var body struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		LoginURL  string `json:"login_url"`
		SessionID string `json:"session_id"`
	}
	if e := json.NewDecoder(resp.Body).Decode(&body); e != nil {
		return "", "", fmt.Errorf("decode SSO start response: %w", e)
	}
	if body.LoginURL == "" {
		return "", "", fmt.Errorf("server error: %s – %s", body.Code, body.Message)
	}
	return body.LoginURL, body.SessionID, nil
}

// callSSOPoll polls the server until the token is ready, up to 5 minutes.
func (m *Manager) callSSOPoll(sessionID string) (*model.TokenPair, error) {
	endpoint := fmt.Sprintf("%s/v1/share/auth/sso/poll?session_id=%s", m.serverBase(), sessionID)
	deadline := time.Now().Add(5 * time.Minute)

	for time.Now().Before(deadline) {
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("create SSO poll request: %w", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("SSO poll request failed: %w", err)
		}

		var body struct {
			Code      string           `json:"code"`
			Message   string           `json:"message"`
			TokenPair *model.TokenPair `json:"token_pair"`
		}
		decodeErr := json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()
		if decodeErr != nil {
			return nil, fmt.Errorf("decode SSO poll response: %w", decodeErr)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			if body.TokenPair == nil {
				return nil, fmt.Errorf("no token pair in SSO poll response")
			}
			return body.TokenPair, nil
		case http.StatusAccepted:
			fmt.Print(".")
			time.Sleep(2 * time.Second)
		default:
			return nil, fmt.Errorf("SSO poll error %d: %s", resp.StatusCode, body.Message)
		}
	}
	return nil, fmt.Errorf("SSO login timed out after 5 minutes")
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

// openBrowser attempts to open the URL in the system default browser.
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	_ = cmd.Start()
}

// authTransport is an http.RoundTripper that injects the Bearer token into every request.
type authTransport struct {
	manager *Manager
	base    http.RoundTripper
}

// RoundTrip obtains a valid token from the Manager and sets the Authorization header.
// On HTTP 401, it refreshes the token and retries once for body-less request methods.
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.manager.Token()
	if err != nil {
		return nil, fmt.Errorf("auth transport: %w", err)
	}

	clone := req.Clone(req.Context())
	if token != "" {
		clone.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := t.base.RoundTrip(clone)
	if err != nil {
		return nil, err
	}

	// On 401, attempt a one-time refresh+retry for requests without a body.
	if resp.StatusCode == http.StatusUnauthorized &&
		(req.Method == http.MethodGet || req.Method == http.MethodHead || req.Method == http.MethodDelete) {
		resp.Body.Close()
		if refreshErr := t.manager.ForceRefresh(); refreshErr == nil {
			token, _ = t.manager.Token()
			clone2 := req.Clone(req.Context())
			if token != "" {
				clone2.Header.Set("Authorization", "Bearer "+token)
			}
			return t.base.RoundTrip(clone2)
		}
	}

	return resp, nil
}
