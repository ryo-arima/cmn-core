package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
)

// ---- Group membership -------------------------------------------------------

func (m *keycloakManager) ListGroupMembers(ctx context.Context, groupID string) ([]model.LoUser, error) {
	status, body, err := m.do(ctx, http.MethodGet, m.adminURL("/groups/"+groupID+"/members"), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("keycloak: list group members status %d", status)
	}
	var kus []model.KcUser
	if err := json.Unmarshal(body, &kus); err != nil {
		return nil, fmt.Errorf("keycloak: parse member list: %w", err)
	}
	attrKey := "cmn_group_" + groupID + "_role"
	users := make([]model.LoUser, 0, len(kus))
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
	var ku model.KcUser
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
	var ku model.KcUser
	if err := json.Unmarshal(body, &ku); err != nil || ku.Attributes == nil {
		return nil
	}
	delete(ku.Attributes, "cmn_group_"+groupID+"_role")
	payload, _ := json.Marshal(map[string]interface{}{"attributes": ku.Attributes})
	m.do(ctx, http.MethodPut, m.adminURL("/users/"+userID), payload) //nolint:errcheck
	return nil
}
