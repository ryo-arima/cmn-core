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

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
)

// ---- User management -------------------------------------------------------

func (m *keycloakManager) GetUser(ctx context.Context, id string) (*model.LoUser, error) {
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
	var ku model.KcUser
	if err := json.Unmarshal(body, &ku); err != nil {
		return nil, fmt.Errorf("keycloak: parse user: %w", err)
	}
	return kcUserToModel(ku), nil
}

func (m *keycloakManager) ListUsers(ctx context.Context) ([]model.LoUser, error) {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/users"), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("keycloak: list users status %d", status)
	}
	var kus []model.KcUser
	if err := json.Unmarshal(body, &kus); err != nil {
		return nil, fmt.Errorf("keycloak: parse user list: %w", err)
	}
	users := make([]model.LoUser, 0, len(kus))
	for _, ku := range kus {
		users = append(users, *kcUserToModel(ku))
	}
	return users, nil
}

func (m *keycloakManager) CreateUser(ctx context.Context, input request.RrCreateUser) (*model.LoUser, error) {
	payload := model.KcUser{
		Username:  input.Username,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Enabled:   input.Enabled,
	}
	if input.Password != "" {
		payload.Credentials = []model.KcCredential{
			{Type: "password", Value: input.Password, Temporary: true},
		}
	}

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

func (m *keycloakManager) UpdateUser(ctx context.Context, id string, input request.RrUpdateUser) error {
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

	var t model.KcTokenResponse
	if err := json.Unmarshal(body, &t); err != nil {
		return "", fmt.Errorf("keycloak: parse login response: %w", err)
	}
	return t.AccessToken, nil
}

// kcUserToModel converts a Keycloak internal user struct to the domain model.
func kcUserToModel(ku model.KcUser) *model.LoUser {
	u := &model.LoUser{
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
