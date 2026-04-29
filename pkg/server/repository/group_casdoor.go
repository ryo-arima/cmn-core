package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
)

// ---- Group management -------------------------------------------------------

func (m *casdoorManager) GetGroup(ctx context.Context, id string) (*model.LoGroup, error) {
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
	var cg model.CdGroup
	if err := json.Unmarshal(r.Data, &cg); err != nil {
		return nil, fmt.Errorf("casdoor: parse group: %w", err)
	}
	return &model.LoGroup{ID: cg.Name, UUID: cg.Name, Name: cg.Name}, nil
}

func (m *casdoorManager) ListGroups(ctx context.Context) ([]model.LoGroup, error) {
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
	var cgs []model.CdGroup
	if err := json.Unmarshal(r.Data, &cgs); err != nil {
		return nil, fmt.Errorf("casdoor: parse group list: %w", err)
	}
	groups := make([]model.LoGroup, 0, len(cgs))
	for _, cg := range cgs {
		groups = append(groups, model.LoGroup{ID: cg.Name, UUID: cg.Name, Name: cg.Name})
	}
	return groups, nil
}

func (m *casdoorManager) CreateGroup(ctx context.Context, input request.RrCreateGroup) (*model.LoGroup, error) {
	payload := model.CdGroup{Owner: m.cfg.Organization, Name: input.Name}
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

func (m *casdoorManager) UpdateGroup(ctx context.Context, id string, input request.RrUpdateGroup) error {
	payload := model.CdGroup{Owner: m.cfg.Organization, Name: input.Name}
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
	payload := model.CdGroup{Owner: m.cfg.Organization, Name: id}
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
