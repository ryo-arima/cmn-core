package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	clientauth "github.com/ryo-arima/cmn-core/pkg/client/share"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// GroupInternal is the data-access interface for group operations via /v1/internal.
type GroupInternal interface {
	ListMyGroups() response.RrIdPGroups
	GetGroup(id string) response.RrSingleIdPGroup
	CreateGroup(req request.RrCreateGroup) response.RrSingleIdPGroup
	UpdateGroup(id string, req request.RrUpdateGroup) response.RrCommons
	DeleteGroup(id string) response.RrCommons
}

// GroupPrivate extends GroupInternal with admin-scope group listing.
type GroupPrivate interface {
	GroupInternal
	ListGroups() response.RrIdPGroups
}

// ---- internal repo ---------------------------------------------------------

type groupInternalRepo struct {
	base   string
	client *http.Client
}

// NewGroupInternal creates a GroupInternal repository backed by /v1/internal.
func NewGroupInternal(conf config.BaseConfig, manager *clientauth.Manager) GroupInternal {
	return &groupInternalRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal",
		client: manager.HTTPClient(),
	}
}

func (r *groupInternalRepo) doJSON(method, url string, body interface{}, out interface{}) error {
	var req *http.Request
	var err error
	if body != nil {
		b, merr := json.Marshal(body)
		if merr != nil {
			return fmt.Errorf("marshal: %w", merr)
		}
		req, err = http.NewRequest(method, url, bytes.NewReader(b))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

func (r *groupInternalRepo) ListMyGroups() response.RrIdPGroups {
	var out response.RrIdPGroups
	if err := r.doJSON("GET", r.base+"/groups", nil, &out); err != nil {
		out.Code = "CLIENT_GROUP_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *groupInternalRepo) GetGroup(id string) response.RrSingleIdPGroup {
	var out response.RrSingleIdPGroup
	url := fmt.Sprintf("%s/group?id=%s", r.base, id)
	if err := r.doJSON("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_GROUP_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *groupInternalRepo) CreateGroup(req request.RrCreateGroup) response.RrSingleIdPGroup {
	var out response.RrSingleIdPGroup
	if err := r.doJSON("POST", r.base+"/groups", req, &out); err != nil {
		out.Code = "CLIENT_GROUP_CREATE_001"
		out.Message = err.Error()
	}
	return out
}

func (r *groupInternalRepo) UpdateGroup(id string, req request.RrUpdateGroup) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/groups/%s", r.base, id)
	if err := r.doJSON("PUT", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_GROUP_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (r *groupInternalRepo) DeleteGroup(id string) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/groups/%s", r.base, id)
	if err := r.doJSON("DELETE", url, nil, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_GROUP_DELETE_001", Message: err.Error()}
	}
	return out
}

// ---- private repo ----------------------------------------------------------

type groupPrivateRepo struct {
	base   string
	client *http.Client
}

// NewGroupPrivate creates a GroupPrivate repository backed by /v1/private.
func NewGroupPrivate(conf config.BaseConfig, manager *clientauth.Manager) GroupPrivate {
	return &groupPrivateRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/private",
		client: manager.HTTPClient(),
	}
}

func (r *groupPrivateRepo) doJSON(method, url string, body interface{}, out interface{}) error {
	var req *http.Request
	var err error
	if body != nil {
		b, merr := json.Marshal(body)
		if merr != nil {
			return fmt.Errorf("marshal: %w", merr)
		}
		req, err = http.NewRequest(method, url, bytes.NewReader(b))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

func (r *groupPrivateRepo) ListMyGroups() response.RrIdPGroups {
	var out response.RrIdPGroups
	if err := r.doJSON("GET", r.base+"/groups", nil, &out); err != nil {
		out.Code = "CLIENT_GROUP_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *groupPrivateRepo) ListGroups() response.RrIdPGroups {
	var out response.RrIdPGroups
	if err := r.doJSON("GET", r.base+"/groups", nil, &out); err != nil {
		out.Code = "CLIENT_GROUP_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *groupPrivateRepo) GetGroup(id string) response.RrSingleIdPGroup {
	var out response.RrSingleIdPGroup
	url := fmt.Sprintf("%s/group?id=%s", r.base, id)
	if err := r.doJSON("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_GROUP_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *groupPrivateRepo) CreateGroup(req request.RrCreateGroup) response.RrSingleIdPGroup {
	var out response.RrSingleIdPGroup
	if err := r.doJSON("POST", r.base+"/groups", req, &out); err != nil {
		out.Code = "CLIENT_GROUP_CREATE_001"
		out.Message = err.Error()
	}
	return out
}

func (r *groupPrivateRepo) UpdateGroup(id string, req request.RrUpdateGroup) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/groups/%s", r.base, id)
	if err := r.doJSON("PUT", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_GROUP_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (r *groupPrivateRepo) DeleteGroup(id string) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/groups/%s", r.base, id)
	if err := r.doJSON("DELETE", url, nil, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_GROUP_DELETE_001", Message: err.Error()}
	}
	return out
}
