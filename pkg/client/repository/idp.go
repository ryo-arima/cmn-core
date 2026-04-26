package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	clientauth "github.com/ryo-arima/cmn-core/pkg/client/auth"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// IdP is the data-access interface for identity-provider operations.
// basePath is either "/v1/internal" (app client) or "/v1/private" (admin client).
type IdP interface {
	// Own user (internal only; admin uses GetUser)
	GetMyUser() response.SingleIdPUser
	UpdateMyUser(req request.UpdateUser) response.Commons

	// Any user by ID; app client fetches from /v1/internal/user?id=
	GetUser(id string) response.SingleIdPUser
	// Users in caller's groups (app: /v1/internal/users)
	ListGroupUsers() response.IdPUsers

	// Groups
	ListMyGroups() response.IdPGroups
	GetGroup(id string) response.SingleIdPGroup
	CreateGroup(req request.CreateGroup) response.SingleIdPGroup
	UpdateGroup(id string, req request.UpdateGroup) response.Commons
	DeleteGroup(id string) response.Commons

	// Members
	ListGroupMembers(groupID string) response.IdPUsers
	AddGroupMember(groupID string, req request.AddGroupMember) response.Commons
	RemoveGroupMember(groupID string, req request.AddGroupMember) response.Commons
}

// IdPAdmin extends IdP with admin-only user management operations.
type IdPAdmin interface {
	IdP
	ListUsers() response.IdPUsers
	CreateUser(req request.CreateUser) response.SingleIdPUser
	UpdateUser(id string, req request.UpdateUser) response.Commons
	DeleteUser(id string) response.Commons
	ListGroups() response.IdPGroups
}

type idpRepo struct {
	base   string // e.g. "http://localhost:8000/v1/internal"
	client *http.Client
}

// NewIdP creates an IdP repository for the given API prefix.
// prefix should be "/v1/internal" for the app client.
func NewIdP(conf config.BaseConfig, manager *clientauth.Manager) IdP {
	return &idpRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal",
		client: manager.HTTPClient(),
	}
}

// NewIdPAdmin creates an IdPAdmin repository backed by /v1/private.
func NewIdPAdmin(conf config.BaseConfig, manager *clientauth.Manager) IdPAdmin {
	return &idpRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/private",
		client: manager.HTTPClient(),
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func (r *idpRepo) do(method, url string, body interface{}, out interface{}) error {
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

func errCommons(code, msg string) response.Commons {
	return response.Commons{Code: code, Message: msg}
}

// ── Own user ─────────────────────────────────────────────────────────────────

func (r *idpRepo) GetMyUser() response.SingleIdPUser {
	var out response.SingleIdPUser
	if err := r.do("GET", r.base+"/user", nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) UpdateMyUser(req request.UpdateUser) response.Commons {
	var out response.Commons
	if err := r.do("PUT", r.base+"/user", req, &out); err != nil {
		return errCommons("CLIENT_USER_UPDATE_001", err.Error())
	}
	return out
}

// ── Groups ───────────────────────────────────────────────────────────────────

func (r *idpRepo) ListMyGroups() response.IdPGroups {
	var out response.IdPGroups
	if err := r.do("GET", r.base+"/groups", nil, &out); err != nil {
		out.Code = "CLIENT_GROUP_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) ListGroups() response.IdPGroups {
	var out response.IdPGroups
	if err := r.do("GET", r.base+"/groups", nil, &out); err != nil {
		out.Code = "CLIENT_GROUP_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) GetGroup(id string) response.SingleIdPGroup {
	var out response.SingleIdPGroup
	url := fmt.Sprintf("%s/group?id=%s", r.base, id)
	if err := r.do("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_GROUP_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) CreateGroup(req request.CreateGroup) response.SingleIdPGroup {
	var out response.SingleIdPGroup
	if err := r.do("POST", r.base+"/groups", req, &out); err != nil {
		out.Code = "CLIENT_GROUP_CREATE_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) UpdateGroup(id string, req request.UpdateGroup) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/groups/%s", r.base, id)
	if err := r.do("PUT", url, req, &out); err != nil {
		return errCommons("CLIENT_GROUP_UPDATE_001", err.Error())
	}
	return out
}

func (r *idpRepo) DeleteGroup(id string) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/groups/%s", r.base, id)
	if err := r.do("DELETE", url, nil, &out); err != nil {
		return errCommons("CLIENT_GROUP_DELETE_001", err.Error())
	}
	return out
}

// ── Members ──────────────────────────────────────────────────────────────────

func (r *idpRepo) ListGroupMembers(groupID string) response.IdPUsers {
	var out response.IdPUsers
	url := fmt.Sprintf("%s/members?group_id=%s", r.base, groupID)
	if err := r.do("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_MEMBER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) AddGroupMember(groupID string, req request.AddGroupMember) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/member/%s", r.base, groupID)
	if err := r.do("POST", url, req, &out); err != nil {
		return errCommons("CLIENT_MEMBER_ADD_001", err.Error())
	}
	return out
}

func (r *idpRepo) RemoveGroupMember(groupID string, req request.AddGroupMember) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/member/%s", r.base, groupID)
	if err := r.do("DELETE", url, req, &out); err != nil {
		return errCommons("CLIENT_MEMBER_REMOVE_001", err.Error())
	}
	return out
}

// ── Admin-only: users ────────────────────────────────────────────────────────

func (r *idpRepo) ListGroupUsers() response.IdPUsers {
	var out response.IdPUsers
	if err := r.do("GET", r.base+"/users", nil, &out); err != nil {
		out.Code = "CLIENT_USER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) ListUsers() response.IdPUsers {
	var out response.IdPUsers
	if err := r.do("GET", r.base+"/users", nil, &out); err != nil {
		out.Code = "CLIENT_USER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) GetUser(id string) response.SingleIdPUser {
	var out response.SingleIdPUser
	url := fmt.Sprintf("%s/user?id=%s", r.base, id)
	if err := r.do("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) CreateUser(req request.CreateUser) response.SingleIdPUser {
	var out response.SingleIdPUser
	if err := r.do("POST", r.base+"/users", req, &out); err != nil {
		out.Code = "CLIENT_USER_CREATE_001"
		out.Message = err.Error()
	}
	return out
}

func (r *idpRepo) UpdateUser(id string, req request.UpdateUser) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/users/%s", r.base, id)
	if err := r.do("PUT", url, req, &out); err != nil {
		return errCommons("CLIENT_USER_UPDATE_001", err.Error())
	}
	return out
}

func (r *idpRepo) DeleteUser(id string) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/users/%s", r.base, id)
	if err := r.do("DELETE", url, nil, &out); err != nil {
		return errCommons("CLIENT_USER_DELETE_001", err.Error())
	}
	return out
}
