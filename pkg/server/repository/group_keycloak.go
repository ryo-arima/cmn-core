package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
)

// ---- Group management -------------------------------------------------------

func (m *keycloakManager) GetGroup(ctx context.Context, id string) (*model.LoGroup, error) {
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
	var kg model.KcGroup
	if err := json.Unmarshal(body, &kg); err != nil {
		return nil, fmt.Errorf("keycloak: parse group: %w", err)
	}
	return &model.LoGroup{ID: kg.ID, UUID: kg.ID, Name: kg.Name, Path: kg.Path}, nil
}

func (m *keycloakManager) ListGroups(ctx context.Context) ([]model.LoGroup, error) {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/groups"), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("keycloak: list groups status %d", status)
	}
	var kgs []model.KcGroup
	if err := json.Unmarshal(body, &kgs); err != nil {
		return nil, fmt.Errorf("keycloak: parse group list: %w", err)
	}
	groups := make([]model.LoGroup, 0, len(kgs))
	for _, kg := range kgs {
		groups = append(groups, model.LoGroup{ID: kg.ID, UUID: kg.ID, Name: kg.Name, Path: kg.Path})
	}
	return groups, nil
}

func (m *keycloakManager) CreateGroup(ctx context.Context, input request.RrCreateGroup) (*model.LoGroup, error) {
	token, err := m.getToken(ctx)
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(model.KcGroup{Name: input.Name})
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

func (m *keycloakManager) UpdateGroup(ctx context.Context, id string, input request.RrUpdateGroup) error {
	status, body, err := m.do(ctx, http.MethodPut, m.adminURL("/groups/"+id), model.KcGroup{Name: input.Name})
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
