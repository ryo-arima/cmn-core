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
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
)

// ---- internal Keycloak JSON types ------------------------------------------

type kcUser struct {
	ID               string              `json:"id,omitempty"`
	Username         string              `json:"username,omitempty"`
	Email            string              `json:"email,omitempty"`
	FirstName        string              `json:"firstName,omitempty"`
	LastName         string              `json:"lastName,omitempty"`
	Enabled          bool                `json:"enabled"`
	EmailVerified    bool                `json:"emailVerified,omitempty"`
	CreatedTimestamp int64               `json:"createdTimestamp,omitempty"`
	Credentials      []kcCred            `json:"credentials,omitempty"`
	Attributes       map[string][]string `json:"attributes,omitempty"`
}

type kcCred struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type kcGroup struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type kcTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

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

	var t kcTokenResponse
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

// ---- User management -------------------------------------------------------

func (m *keycloakManager) GetUser(ctx context.Context, id string) (*model.IdPUser, error) {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/users/"+id), nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("keycloak: user %q not found", id)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("keycloak: get user status %d", status)
	}
	var ku kcUser
	if err := json.Unmarshal(body, &ku); err != nil {
		return nil, fmt.Errorf("keycloak: parse user: %w", err)
	}
	return kcUserToModel(ku), nil
}

func (m *keycloakManager) ListUsers(ctx context.Context) ([]model.IdPUser, error) {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/users"), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("keycloak: list users status %d", status)
	}
	var kus []kcUser
	if err := json.Unmarshal(body, &kus); err != nil {
		return nil, fmt.Errorf("keycloak: parse user list: %w", err)
	}
	users := make([]model.IdPUser, 0, len(kus))
	for _, ku := range kus {
		users = append(users, *kcUserToModel(ku))
	}
	return users, nil
}

func (m *keycloakManager) CreateUser(ctx context.Context, input request.CreateUser) (*model.IdPUser, error) {
	payload := kcUser{
		Username:  input.Username,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Enabled:   input.Enabled,
	}
	if input.Password != "" {
		payload.Credentials = []kcCred{
			{Type: "password", Value: input.Password, Temporary: true},
		}
	}

	// Keycloak responds with 201 + Location; no body.
	token, err := m.getToken(ctx)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.adminURL("/users"), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("keycloak: create user: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("keycloak: create user status %d: %s", resp.StatusCode, respBody)
	}

	newID := kcExtractCreatedID(resp.Header.Get("Location"))
	if newID == "" {
		return nil, fmt.Errorf("keycloak: could not determine new user ID from Location header")
	}
	return m.GetUser(ctx, newID)
}

func (m *keycloakManager) UpdateUser(ctx context.Context, id string, input request.UpdateUser) error {
	payload := make(map[string]interface{})
	if input.Email != nil {
		payload["email"] = *input.Email
	}
	if input.FirstName != nil {
		payload["firstName"] = *input.FirstName
	}
	if input.LastName != nil {
		payload["lastName"] = *input.LastName
	}
	if input.Enabled != nil {
		payload["enabled"] = *input.Enabled
	}
	status, body, err := m.do(ctx, http.MethodPut, m.adminURL("/users/"+id), payload)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("keycloak: update user status %d: %s", status, body)
	}
	return nil
}

func (m *keycloakManager) DeleteUser(ctx context.Context, id string) error {
	status, body, err := m.do(ctx, http.MethodDelete, m.adminURL("/users/"+id), nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("keycloak: delete user status %d: %s", status, body)
	}
	return nil
}

// ---- Group management ------------------------------------------------------

func (m *keycloakManager) GetGroup(ctx context.Context, id string) (*model.IdPGroup, error) {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/groups/"+id), nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("keycloak: group %q not found", id)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("keycloak: get group status %d", status)
	}
	var kg kcGroup
	if err := json.Unmarshal(body, &kg); err != nil {
		return nil, fmt.Errorf("keycloak: parse group: %w", err)
	}
	return &model.IdPGroup{ID: kg.ID, Name: kg.Name, Path: kg.Path}, nil
}

func (m *keycloakManager) ListGroups(ctx context.Context) ([]model.IdPGroup, error) {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/groups"), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("keycloak: list groups status %d", status)
	}
	var kgs []kcGroup
	if err := json.Unmarshal(body, &kgs); err != nil {
		return nil, fmt.Errorf("keycloak: parse group list: %w", err)
	}
	groups := make([]model.IdPGroup, 0, len(kgs))
	for _, kg := range kgs {
		groups = append(groups, model.IdPGroup{ID: kg.ID, Name: kg.Name, Path: kg.Path})
	}
	return groups, nil
}

func (m *keycloakManager) CreateGroup(ctx context.Context, input request.CreateGroup) (*model.IdPGroup, error) {
	token, err := m.getToken(ctx)
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(kcGroup{Name: input.Name})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, m.adminURL("/groups"), bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("keycloak: create group: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("keycloak: create group status %d: %s", resp.StatusCode, respBody)
	}
	newID := kcExtractCreatedID(resp.Header.Get("Location"))
	if newID == "" {
		return nil, fmt.Errorf("keycloak: could not determine new group ID from Location header")
	}
	return m.GetGroup(ctx, newID)
}

func (m *keycloakManager) UpdateGroup(ctx context.Context, id string, input request.UpdateGroup) error {
	status, body, err := m.do(ctx, http.MethodPut, m.adminURL("/groups/"+id), kcGroup{Name: input.Name})
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("keycloak: update group status %d: %s", status, body)
	}
	return nil
}

func (m *keycloakManager) DeleteGroup(ctx context.Context, id string) error {
	status, body, err := m.do(ctx, http.MethodDelete, m.adminURL("/groups/"+id), nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("keycloak: delete group status %d: %s", status, body)
	}
	return nil
}

// ---- Group membership ------------------------------------------------------

func (m *keycloakManager) ListGroupMembers(ctx context.Context, groupID string) ([]model.IdPUser, error) {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/groups/"+groupID+"/members"), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("keycloak: list group members status %d", status)
	}
	var kus []kcUser
	if err := json.Unmarshal(body, &kus); err != nil {
		return nil, fmt.Errorf("keycloak: parse member list: %w", err)
	}
	attrKey := "cmn_group_" + groupID + "_role"
	users := make([]model.IdPUser, 0, len(kus))
	for _, ku := range kus {
		u := *kcUserToModel(ku)
		if vals, ok := ku.Attributes[attrKey]; ok && len(vals) > 0 {
			u.Role = vals[0]
		}
		users = append(users, u)
	}
	return users, nil
}

func (m *keycloakManager) AddUserToGroup(ctx context.Context, userID, groupID, role string) error {
	path := fmt.Sprintf("/users/%s/groups/%s", userID, groupID)
	status, body, err := m.do(ctx, http.MethodPut, m.adminURL(path), nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("keycloak: add user to group status %d: %s", status, body)
	}
	return m.setUserGroupRole(ctx, userID, groupID, role)
}

// setUserGroupRole stores the group-specific role as a Keycloak user attribute.
func (m *keycloakManager) setUserGroupRole(ctx context.Context, userID, groupID, role string) error {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/users/"+userID), nil)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("keycloak: get user for attr update status %d", status)
	}
	var ku kcUser
	if err := json.Unmarshal(body, &ku); err != nil {
		return fmt.Errorf("keycloak: parse user for attr update: %w", err)
	}
	if ku.Attributes == nil {
		ku.Attributes = make(map[string][]string)
	}
	ku.Attributes["cmn_group_"+groupID+"_role"] = []string{role}
	payload, _ := json.Marshal(map[string]interface{}{"attributes": ku.Attributes})
	status, body, err = m.do(ctx, http.MethodPut, m.adminURL("/users/"+userID), payload)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("keycloak: set user group role attr status %d: %s", status, body)
	}
	return nil
}

func (m *keycloakManager) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	path := fmt.Sprintf("/users/%s/groups/%s", userID, groupID)
	status, body, err := m.do(ctx, http.MethodDelete, m.adminURL(path), nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("keycloak: remove user from group status %d: %s", status, body)
	}
	// Best-effort: clean up the role attribute.
	_ = m.clearUserGroupRole(ctx, userID, groupID)
	return nil
}

// clearUserGroupRole removes the group-specific role attribute from a Keycloak user.
func (m *keycloakManager) clearUserGroupRole(ctx context.Context, userID, groupID string) error {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/users/"+userID), nil)
	if err != nil || status != http.StatusOK {
		return nil // best-effort
	}
	var ku kcUser
	if err := json.Unmarshal(body, &ku); err != nil || ku.Attributes == nil {
		return nil
	}
	delete(ku.Attributes, "cmn_group_"+groupID+"_role")
	payload, _ := json.Marshal(map[string]interface{}{"attributes": ku.Attributes})
	m.do(ctx, http.MethodPut, m.adminURL("/users/"+userID), payload) //nolint:errcheck
	return nil
}

// Login performs an ROPC grant on behalf of a user and returns the issued access token.
func (m *keycloakManager) Login(ctx context.Context, username, password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", m.cfg.BaseURL, m.cfg.Realm)
	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", m.cfg.AdminClientID)
	form.Set("client_secret", m.cfg.AdminClientSecret)
	form.Set("username", username)
	form.Set("password", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("keycloak: build login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("keycloak: login request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("keycloak: login status %d", resp.StatusCode)
	}

	var t kcTokenResponse
	if err := json.Unmarshal(body, &t); err != nil {
		return "", fmt.Errorf("keycloak: parse login response: %w", err)
	}
	return t.AccessToken, nil
}

// ---- helpers ---------------------------------------------------------------

func kcUserToModel(ku kcUser) *model.IdPUser {
	u := &model.IdPUser{
		ID:        ku.ID,
		Username:  ku.Username,
		Email:     ku.Email,
		FirstName: ku.FirstName,
		LastName:  ku.LastName,
		Enabled:   ku.Enabled,
	}
	if ku.CreatedTimestamp > 0 {
		u.CreatedAt = time.UnixMilli(ku.CreatedTimestamp)
	}
	return u
}
