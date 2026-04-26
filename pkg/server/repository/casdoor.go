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
	"time"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
)

// ---- internal Casdoor JSON types -------------------------------------------

type cdUser struct {
	Owner        string            `json:"owner"`
	Name         string            `json:"name"`
	ID           string            `json:"id,omitempty"`
	Email        string            `json:"email,omitempty"`
	DisplayName  string            `json:"displayName,omitempty"`
	FirstName    string            `json:"firstName,omitempty"`
	LastName     string            `json:"lastName,omitempty"`
	IsAdmin      bool              `json:"isAdmin,omitempty"`
	IsForbidden  bool              `json:"isForbidden,omitempty"`
	CreatedTime  string            `json:"createdTime,omitempty"`
	Password     string            `json:"password,omitempty"`
	PasswordSalt string            `json:"passwordSalt,omitempty"`
	Groups       []string          `json:"groups,omitempty"`
	Properties   map[string]string `json:"properties,omitempty"`
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

// authParams appends clientId and clientSecret to query params.
// Casdoor admin API authenticates via these query parameters.
func (m *casdoorManager) authParams(q url.Values) url.Values {
	if q == nil {
		q = url.Values{}
	}
	q.Set("clientId", m.cfg.ClientID)
	q.Set("clientSecret", m.cfg.ClientSecret)
	return q
}

// do performs an authenticated request to the Casdoor admin API.
// Authentication is via clientId/clientSecret query parameters.
func (m *casdoorManager) do(ctx context.Context, method, rawURL string, payload interface{}) (int, []byte, error) {
	// Append auth params to URL.
	u, err := url.Parse(rawURL)
	if err != nil {
		return 0, nil, fmt.Errorf("casdoor: parse URL: %w", err)
	}
	u.RawQuery = m.authParams(u.Query()).Encode()

	var bodyReader io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, fmt.Errorf("casdoor: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return 0, nil, fmt.Errorf("casdoor: build request: %w", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("casdoor: %s %s: %w", method, u.Path, err)
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
	// Casdoor API expects "org/name" format.
	// Accept both plain name ("group001") and already-qualified ("cmn/group001").
	fullID := id
	if !strings.Contains(id, "/") {
		fullID = m.cfg.Organization + "/" + id
	}
	q.Set("id", fullID)
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

// qualifyGroup ensures the group ID has the org prefix ("cmn/group001").
// Accepts both plain name ("group001") and already-qualified ("cmn/group001").
func (m *casdoorManager) qualifyGroup(id string) string {
	if strings.Contains(id, "/") {
		return id
	}
	return m.cfg.Organization + "/" + id
}
func (m *casdoorManager) ListGroupMembers(ctx context.Context, groupID string) ([]model.IdPUser, error) {
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
	var members []model.IdPUser
	for _, cu := range cus {
		for _, g := range cu.Groups {
			// Normalize to plain name: "cmn/group001" → "group001"
			normG := g
			if i := strings.LastIndex(g, "/"); i >= 0 {
				normG = g[i+1:]
			}
			normGID := groupID
			if i := strings.LastIndex(groupID, "/"); i >= 0 {
				normGID = groupID[i+1:]
			}
			if normG == normGID {
				u := *cdUserToModel(cu)
				if cu.Properties != nil {
					attrKey := "cmn_group_" + normGID + "_role"
					u.Role = cu.Properties[attrKey]
				}
				members = append(members, u)
				break
			}
		}
	}
	return members, nil
}

func (m *casdoorManager) AddUserToGroup(ctx context.Context, userID, groupID, role string) error {
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
	qualified := m.qualifyGroup(groupID)
	for _, g := range groups {
		if g == qualified {
			return nil // already a member
		}
	}
	groups = append(groups, qualified)
	userMap["groups"] = groups

	// Casdoor requires displayName to be non-empty on update.
	// Fall back to the user's name if displayName is absent or blank.
	if dn, _ := userMap["displayName"].(string); dn == "" {
		if name, _ := userMap["name"].(string); name != "" {
			userMap["displayName"] = name
		}
	}

	// Store the group-specific role as a user property.
	normGID := groupID
	if i := strings.LastIndex(groupID, "/"); i >= 0 {
		normGID = groupID[i+1:]
	}
	attrKey := "cmn_group_" + normGID + "_role"
	propsMap, _ := userMap["properties"].(map[string]interface{})
	if propsMap == nil {
		propsMap = make(map[string]interface{})
	}
	propsMap[attrKey] = role
	userMap["properties"] = propsMap

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
	qualified := m.qualifyGroup(groupID)
	filtered := make([]string, 0, len(groups))
	for _, g := range groups {
		if g != qualified {
			filtered = append(filtered, g)
		}
	}
	userMap["groups"] = filtered

	// Casdoor requires displayName to be non-empty on update.
	if dn, _ := userMap["displayName"].(string); dn == "" {
		if name, _ := userMap["name"].(string); name != "" {
			userMap["displayName"] = name
		}
	}

	// Clean up the group-specific role property.
	normGID := groupID
	if i := strings.LastIndex(groupID, "/"); i >= 0 {
		normGID = groupID[i+1:]
	}
	attrKey := "cmn_group_" + normGID + "_role"
	if propsMap, ok := userMap["properties"].(map[string]interface{}); ok {
		delete(propsMap, attrKey)
		userMap["properties"] = propsMap
	}

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

// Login performs an ROPC grant on behalf of a user and returns the issued access token.
func (m *casdoorManager) Login(ctx context.Context, username, password string) (string, error) {
	tokenURL := m.apiURL("/api/login/oauth/access_token")
	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", m.cfg.ClientID)
	form.Set("client_secret", m.cfg.ClientSecret)
	form.Set("username", username)
	form.Set("password", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("casdoor: build login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("casdoor: login request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("casdoor: login status %d", resp.StatusCode)
	}

	var t cdTokenResponse
	if err := json.Unmarshal(body, &t); err != nil {
		return "", fmt.Errorf("casdoor: parse login response: %w", err)
	}
	return t.AccessToken, nil
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
