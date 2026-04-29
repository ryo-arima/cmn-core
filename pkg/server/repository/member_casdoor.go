package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
)

// ---- Group membership -------------------------------------------------------
// Casdoor manages group membership via the "groups" field on the user object.

// qualifyGroup ensures the group ID has the org prefix ("cmn/group001").
// Accepts both plain name ("group001") and already-qualified ("cmn/group001").
func (m *casdoorManager) qualifyGroup(id string) string {
	if strings.Contains(id, "/") {
		return id
	}
	return m.cfg.Organization + "/" + id
}

func (m *casdoorManager) ListGroupMembers(ctx context.Context, groupID string) ([]model.LoUser, error) {
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
	var members []model.LoUser
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

// parseStringSlice safely converts an interface{} holding []interface{} of strings.
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
