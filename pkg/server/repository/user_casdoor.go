package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
)

// ---- User management -------------------------------------------------------

func (m *casdoorManager) GetUser(ctx context.Context, id string) (*model.LoUser, error) {
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
	var cu model.CdUser
	if err := json.Unmarshal(r.Data, &cu); err != nil {
		return nil, fmt.Errorf("casdoor: parse user: %w", err)
	}
	return cdUserToModel(cu), nil
}

func (m *casdoorManager) ListUsers(ctx context.Context) ([]model.LoUser, error) {
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
	var cus []model.CdUser
	if err := json.Unmarshal(r.Data, &cus); err != nil {
		return nil, fmt.Errorf("casdoor: parse user list: %w", err)
	}
	users := make([]model.LoUser, 0, len(cus))
	for _, cu := range cus {
		users = append(users, *cdUserToModel(cu))
	}
	return users, nil
}

func (m *casdoorManager) CreateUser(ctx context.Context, input request.RrCreateUser) (*model.LoUser, error) {
	payload := model.CdUser{
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

func (m *casdoorManager) UpdateUser(ctx context.Context, id string, input request.RrUpdateUser) error {
	// Fetch first to preserve unchanged fields.
	existing, err := m.GetUser(ctx, id)
	if err != nil {
		return err
	}
	payload := model.CdUser{
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
	payload := model.CdUser{Owner: m.cfg.Organization, Name: id}
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

	var t model.CdTokenResponse
	if err := json.Unmarshal(body, &t); err != nil {
		return "", fmt.Errorf("casdoor: parse login response: %w", err)
	}
	return t.AccessToken, nil
}

// cdUserToModel converts a Casdoor internal user struct to the domain model.
func cdUserToModel(cu model.CdUser) *model.LoUser {
	u := &model.LoUser{
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
