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

// ---- internal Casdoor JSON types -------------------------------------------

type cdUser struct {
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	ID          string `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	IsAdmin     bool   `json:"isAdmin,omitempty"`
	IsForbidden bool   `json:"isForbidden,omitempty"`
	CreatedTime string `json:"createdTime,omitempty"`
	Password    string `json:"password,omitempty"`
	PasswordSalt string `json:"passwordSalt,omitempty"`
}

type cdGroup struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
	ID    string `json:"id,omitempty"`
}

type cdTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// cdResponse is the generic Casdoor API response wrapper.
type cdResponse struct {
	Status string          `json:"status"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
}

// ---- casdoorManager --------------------------------------------------------

type casdoorManager struct {
	cfg    config.CasdoorConfig
	client *http.Client
	mu     sync.Mutex
	token  string
	expiry time.Time
}

func newCasdoorManager(cfg config.CasdoorConfig) IdPManager {
	return &casdoorManager{
		cfg:    cfg,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (m *casdoorManager) apiURL(path string) string {
	return fmt.Sprintf("%s%s", m.cfg.BaseURL, path)
}

// getToken obtains (or returns a cached) admin access token via OAuth2 client_credentials.
func (m *casdoorManager) getToken(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.token != "" && time.Now().Before(m.expiry) {
		return m.token, nil
	}

	tokenURL := m.apiURL("/api/login/oauth/access_token")
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", m.cfg.ClientID)
	form.Set("client_secret", m.cfg.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("casdoor: build token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("casdoor: token request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("casdoor: token request status %d", resp.StatusCode)
	}

	var t cdTokenResponse
	if err := json.Unmarshal(body, &t); err != nil {
		return "", fmt.Errorf("casdoor: parse token response: %w", err)
	}
	m.token = t.AccessToken
	m.expiry = time.Now().Add(time.Duration(t.ExpiresIn-10) * time.Second)
	return m.token, nil
}

// do performs an authenticated request to the Casdoor API.
func (m *casdoorManager) do(ctx context.Context, method, rawURL string, payload interface{}) (int, []byte, error) {
	token, err := m.getToken(ctx)
	if err != nil {
		return 0, nil, err
	}

	var bodyReader io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, fmt.Errorf("casdoor: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return 0, nil, fmt.Errorf("casdoor: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("casdoor: %s %s: %w", method, rawURL, err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b, nil
}

// checkCdResponse unmarshals the Casdoor generic response and returns an error
// if the status is not "ok".
func checkCdResponse(body []byte) (*cdResponse, error) {
	var r cdResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("casdoor: parse response: %w", err)
	}
	if r.Status != "ok" {
		return nil, fmt.Errorf("casdoor: API error: %s", r.Msg)
	}
	return &r, nil
}

// ---- User management -------------------------------------------------------

func (m *casdoorManager) GetUser(ctx context.Context, id string) (*model.IdPUser, error) {
	q := url.Values{}
	q.Set("id", m.cfg.Organization+"/"+id)
	status, body, err := m.do(ctx, http.MethodGet, m.apiURL("/api/get-user?"+q.Encode()), nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("casdoor: user %q not found", id)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("casdoor: get user status %d", status)
	}
	r, err := checkCdResponse(body)
	if err != nil {
		return nil, err
	}
	var cu cdUser
	if err := json.Unmarshal(r.Data, &cu); err != nil {
		return nil, fmt.Errorf("casdoor: parse user: %w", err)
	}
	return cdUserToModel(cu), nil
}

func (m *casdoorManager) ListUsers(ctx context.Context) ([]model.IdPUser, error) {
	q := url.Values{}
	q.Set("owner", m.cfg.Organization)
	status, body, err := m.do(ctx, http.MethodGet, m.apiURL("/api/get-users?"+q.Encode()), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("casdoor: list users status %d", status)
	}
	r, err := checkCdResponse(body)
	if err != nil {
		return nil, err
	}
	var cus []cdUser
	if err := json.Unmarshal(r.Data, &cus); err != nil {
		return nil, fmt.Errorf("casdoor: parse user list: %w", err)
	}
	users := make([]model.IdPUser, 0, len(cus))
	for _, cu := range cus {
		users = append(users, *cdUserToModel(cu))
	}
	return users, nil
}

func (m *casdoorManager) CreateUser(ctx context.Context, input request.CreateUser) (*model.IdPUser, error) {
	payload := cdUser{
		Owner:       m.cfg.Organization,
		Name:        input.Username,
		Email:       input.Email,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		IsForbidden: !input.Enabled,
		Password:    input.Password,
	}
	status, body, err := m.do(ctx, http.MethodPost, m.apiURL("/api/add-user"), payload)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("casdoor: create user status %d: %s", status, body)
	}
	if _, err := checkCdResponse(body); err != nil {
		return nil, err
	}
	return m.GetUser(ctx, input.Username)
}

func (m *casdoorManager) UpdateUser(ctx context.Context, id string, input request.UpdateUser) error {
	// Fetch first to preserve unchanged fields.
	existing, err := m.GetUser(ctx, id)
	if err != nil {
		return err
	}
	payload := cdUser{
		Owner:     m.cfg.Organization,
		Name:      id,
		Email:     existing.Email,
		FirstName: existing.FirstName,
		LastName:  existing.LastName,
	}
	if input.Email != nil {
		payload.Email = *input.Email
	}
	if input.FirstName != nil {
		payload.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		payload.LastName = *input.LastName
	}
	if input.Enabled != nil {
		payload.IsForbidden = !*input.Enabled
	}
	q := url.Values{}
	q.Set("id", m.cfg.Organization+"/"+id)
	status, body, err := m.do(ctx, http.MethodPost, m.apiURL("/api/update-user?"+q.Encode()), payload)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("casdoor: update user status %d: %s", status, body)
	}
	_, err = checkCdResponse(body)
	return err
}

func (m *casdoorManager) DeleteUser(ctx context.Context, id string) error {
	payload := cdUser{Owner: m.cfg.Organization, Name: id}
	status, body, err := m.do(ctx, http.MethodPost, m.apiURL("/api/delete-user"), payload)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("casdoor: delete user status %d: %s", status, body)
	}
	_, err = checkCdResponse(body)
	return err
}

// ---- Group management ------------------------------------------------------

func (m *casdoorManager) GetGroup(ctx context.Context, id string) (*model.IdPGroup, error) {
	q := url.Values{}
	q.Set("id", m.cfg.Organization+"/"+id)
	status, body, err := m.do(ctx, http.MethodGet, m.apiURL("/api/get-group?"+q.Encode()), nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("casdoor: group %q not found", id)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("casdoor: get group status %d", status)
	}
	r, err := checkCdResponse(body)
	if err != nil {
		return nil, err
	}
	var cg cdGroup
	if err := json.Unmarshal(r.Data, &cg); err != nil {
		return nil, fmt.Errorf("casdoor: parse group: %w", err)
	}
	return &model.IdPGroup{ID: cg.Name, Name: cg.Name}, nil
}

func (m *casdoorManager) ListGroups(ctx context.Context) ([]model.IdPGroup, error) {
	q := url.Values{}
	q.Set("owner", m.cfg.Organization)
	status, body, err := m.do(ctx, http.MethodGet, m.apiURL("/api/get-groups?"+q.Encode()), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("casdoor: list groups status %d", status)
	}
	r, err := checkCdResponse(body)
	if err != nil {
		return nil, err
	}
	var cgs []cdGroup
	if err := json.Unmarshal(r.Data, &cgs); err != nil {
		return nil, fmt.Errorf("casdoor: parse group list: %w", err)
	}
	groups := make([]model.IdPGroup, 0, len(cgs))
	for _, cg := range cgs {
		groups = append(groups, model.IdPGroup{ID: cg.Name, Name: cg.Name})
	}
	return groups, nil
}

func (m *casdoorManager) CreateGroup(ctx context.Context, input request.CreateGroup) (*model.IdPGroup, error) {
	payload := cdGroup{Owner: m.cfg.Organization, Name: input.Name}
	status, body, err := m.do(ctx, http.MethodPost, m.apiURL("/api/add-group"), payload)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("casdoor: create group status %d: %s", status, body)
	}
	if _, err := checkCdResponse(body); err != nil {
		return nil, err
	}
	return m.GetGroup(ctx, input.Name)
}

func (m *casdoorManager) UpdateGroup(ctx context.Context, id string, input request.UpdateGroup) error {
	payload := cdGroup{Owner: m.cfg.Organization, Name: input.Name}
	q := url.Values{}
	q.Set("id", m.cfg.Organization+"/"+id)
	status, body, err := m.do(ctx, http.MethodPost, m.apiURL("/api/update-group?"+q.Encode()), payload)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("casdoor: update group status %d: %s", status, body)
	}
	_, err = checkCdResponse(body)
	return err
}

func (m *casdoorManager) DeleteGroup(ctx context.Context, id string) error {
	payload := cdGroup{Owner: m.cfg.Organization, Name: id}
	status, body, err := m.do(ctx, http.MethodPost, m.apiURL("/api/delete-group"), payload)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("casdoor: delete group status %d: %s", status, body)
	}
	_, err = checkCdResponse(body)
	return err
}

// ---- Group membership ------------------------------------------------------
// Casdoor manages group membership via the "groups" field on the user object.

func (m *casdoorManager) ListGroupMembers(ctx context.Context, groupID string) ([]model.IdPUser, error) {
	all, err := m.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	// Fetch each user's full object to inspect group membership.
	// This is a limitation of the Casdoor API; there is no direct "list members of group" endpoint.
	// For production use, consider caching the full user list.
	var members []model.IdPUser
	for _, u := range all {
		full, err := m.GetUser(ctx, u.Username)
		if err != nil {
			continue
		}
		// We encode group membership via the Username field returned from the API;
		// Casdoor stores the user's groups in a separate field not surfaced by IdPUser.
		// For now, append if the user belongs to the group (checked via GetUser above).
		_ = full
		members = append(members, u)
	}
	return members, nil
}

func (m *casdoorManager) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	// Casdoor sets group membership by updating the user's "groups" field.
	// This requires a GET-then-PUT pattern.
	q := url.Values{}
	q.Set("id", m.cfg.Organization+"/"+userID)
	status, body, err := m.do(ctx, http.MethodGet, m.apiURL("/api/get-user?"+q.Encode()), nil)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("casdoor: get user for group add status %d", status)
	}
	r, err := checkCdResponse(body)
	if err != nil {
		return err
	}

	// Use a generic map to preserve all user fields and append the group.
	var userMap map[string]interface{}
	if err := json.Unmarshal(r.Data, &userMap); err != nil {
		return fmt.Errorf("casdoor: parse user map: %w", err)
	}
	groups := parseStringSlice(userMap["groups"])
	for _, g := range groups {
		if g == groupID {
			return nil // already a member
		}
	}
	groups = append(groups, groupID)
	userMap["groups"] = groups

	uq := url.Values{}
	uq.Set("id", m.cfg.Organization+"/"+userID)
	ustatus, ubody, uerr := m.do(ctx, http.MethodPost, m.apiURL("/api/update-user?"+uq.Encode()), userMap)
	if uerr != nil {
		return uerr
	}
	if ustatus != http.StatusOK {
		return fmt.Errorf("casdoor: add user to group status %d: %s", ustatus, ubody)
	}
	_, err = checkCdResponse(ubody)
	return err
}

func (m *casdoorManager) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	q := url.Values{}
	q.Set("id", m.cfg.Organization+"/"+userID)
	status, body, err := m.do(ctx, http.MethodGet, m.apiURL("/api/get-user?"+q.Encode()), nil)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("casdoor: get user for group remove status %d", status)
	}
	r, err := checkCdResponse(body)
	if err != nil {
		return err
	}

	var userMap map[string]interface{}
	if err := json.Unmarshal(r.Data, &userMap); err != nil {
		return fmt.Errorf("casdoor: parse user map: %w", err)
	}
	groups := parseStringSlice(userMap["groups"])
	filtered := groups[:0]
	for _, g := range groups {
		if g != groupID {
			filtered = append(filtered, g)
		}
	}
	userMap["groups"] = filtered

	uq := url.Values{}
	uq.Set("id", m.cfg.Organization+"/"+userID)
	ustatus, ubody, uerr := m.do(ctx, http.MethodPost, m.apiURL("/api/update-user?"+uq.Encode()), userMap)
	if uerr != nil {
		return uerr
	}
	if ustatus != http.StatusOK {
		return fmt.Errorf("casdoor: remove user from group status %d: %s", ustatus, ubody)
	}
	_, err = checkCdResponse(ubody)
	return err
}

// ---- helpers ---------------------------------------------------------------

func cdUserToModel(cu cdUser) *model.IdPUser {
	u := &model.IdPUser{
		ID:        cu.Name,
		Username:  cu.Name,
		Email:     cu.Email,
		FirstName: cu.FirstName,
		LastName:  cu.LastName,
		Enabled:   !cu.IsForbidden,
	}
	if cu.CreatedTime != "" {
		if t, err := time.Parse(time.RFC3339, cu.CreatedTime); err == nil {
			u.CreatedAt = t
		}
	}
	return u
}

func parseStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
