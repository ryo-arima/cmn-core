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

func (rcvr *casdoorManager) GetGroup(ctx context.Context, id string) (*model.LoGroup, error) {
	q := url.Values{}
	// Casdoor API expects "org/name" format.
	// Accept both plain name ("group001") and already-qualified ("cmn/group001").
	fullID := id
	if !strings.Contains(id, "/") {
		fullID = rcvr.cfg.Organization + "/" + id
	}
	q.Set("id", fullID)
	status, body, err := rcvr.do(ctx, http.MethodGet, rcvr.apiURL("/api/get-group?"+q.Encode()), nil)
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

func (rcvr *casdoorManager) ListGroups(ctx context.Context) ([]model.LoGroup, error) {
	q := url.Values{}
	q.Set("owner", rcvr.cfg.Organization)
	status, body, err := rcvr.do(ctx, http.MethodGet, rcvr.apiURL("/api/get-groups?"+q.Encode()), nil)
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

func (rcvr *casdoorManager) CreateGroup(ctx context.Context, input request.RrCreateGroup) (*model.LoGroup, error) {
	payload := model.CdGroup{Owner: rcvr.cfg.Organization, Name: input.Name}
	status, body, err := rcvr.do(ctx, http.MethodPost, rcvr.apiURL("/api/add-group"), payload)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("casdoor: create group status %d: %s", status, body)
	}
	if _, err := checkCdResponse(body); err != nil {
		return nil, err
	}
	return rcvr.GetGroup(ctx, input.Name)
}

func (rcvr *casdoorManager) UpdateGroup(ctx context.Context, id string, input request.RrUpdateGroup) error {
	payload := model.CdGroup{Owner: rcvr.cfg.Organization, Name: input.Name}
	q := url.Values{}
	q.Set("id", rcvr.cfg.Organization+"/"+id)
	status, body, err := rcvr.do(ctx, http.MethodPost, rcvr.apiURL("/api/update-group?"+q.Encode()), payload)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("casdoor: update group status %d: %s", status, body)
	}
	_, err = checkCdResponse(body)
	return err
}

func (rcvr *casdoorManager) DeleteGroup(ctx context.Context, id string) error {
	payload := model.CdGroup{Owner: rcvr.cfg.Organization, Name: id}
	status, body, err := rcvr.do(ctx, http.MethodPost, rcvr.apiURL("/api/delete-group"), payload)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("casdoor: delete group status %d: %s", status, body)
	}
	_, err = checkCdResponse(body)
	return err
}
