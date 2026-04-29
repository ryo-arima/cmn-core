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

// UserInternal is the data-access interface for user operations via /v1/internal.
type UserInternal interface {
	GetMyUser() response.RrSingleIdPUser
	UpdateMyUser(req request.RrUpdateUser) response.RrCommons
	GetUser(id string) response.RrSingleIdPUser
	ListGroupUsers() response.RrIdPUsers
}

// UserPrivate extends UserInternal with admin CRUD via /v1/private.
type UserPrivate interface {
	UserInternal
	ListUsers() response.RrIdPUsers
	CreateUser(req request.RrCreateUser) response.RrSingleIdPUser
	UpdateUser(id string, req request.RrUpdateUser) response.RrCommons
	DeleteUser(id string) response.RrCommons
}

// UserPublic is the data-access interface for unauthenticated user creation.
type UserPublic interface {
	RegisterUser(req request.RrCreateUser) response.RrSingleIdPUser
}

// ---- internal repo ---------------------------------------------------------

type userInternalRepo struct {
	base   string
	client *http.Client
}

// NewUserInternal creates a UserInternal repository backed by /v1/internal.
func NewUserInternal(conf config.BaseConfig, manager *clientauth.Manager) UserInternal {
	return &userInternalRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal",
		client: manager.HTTPClient(),
	}
}

func (r *userInternalRepo) doJSON(method, url string, body interface{}, out interface{}) error {
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

func (r *userInternalRepo) GetMyUser() response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	if err := r.doJSON("GET", r.base+"/user", nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *userInternalRepo) UpdateMyUser(req request.RrUpdateUser) response.RrCommons {
	var out response.RrCommons
	if err := r.doJSON("PUT", r.base+"/user", req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_USER_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (r *userInternalRepo) GetUser(id string) response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	url := fmt.Sprintf("%s/user?id=%s", r.base, id)
	if err := r.doJSON("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *userInternalRepo) ListGroupUsers() response.RrIdPUsers {
	var out response.RrIdPUsers
	if err := r.doJSON("GET", r.base+"/users", nil, &out); err != nil {
		out.Code = "CLIENT_USER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

// ---- private repo ----------------------------------------------------------

type userPrivateRepo struct {
	base   string
	client *http.Client
}

// NewUserPrivate creates a UserPrivate repository backed by /v1/private.
func NewUserPrivate(conf config.BaseConfig, manager *clientauth.Manager) UserPrivate {
	return &userPrivateRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/private",
		client: manager.HTTPClient(),
	}
}

func (r *userPrivateRepo) doJSON(method, url string, body interface{}, out interface{}) error {
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

func (r *userPrivateRepo) GetMyUser() response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	if err := r.doJSON("GET", r.base+"/user", nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *userPrivateRepo) UpdateMyUser(req request.RrUpdateUser) response.RrCommons {
	var out response.RrCommons
	if err := r.doJSON("PUT", r.base+"/user", req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_USER_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (r *userPrivateRepo) GetUser(id string) response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	url := fmt.Sprintf("%s/user?id=%s", r.base, id)
	if err := r.doJSON("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *userPrivateRepo) ListGroupUsers() response.RrIdPUsers {
	var out response.RrIdPUsers
	if err := r.doJSON("GET", r.base+"/users", nil, &out); err != nil {
		out.Code = "CLIENT_USER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *userPrivateRepo) ListUsers() response.RrIdPUsers {
	var out response.RrIdPUsers
	if err := r.doJSON("GET", r.base+"/users", nil, &out); err != nil {
		out.Code = "CLIENT_USER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *userPrivateRepo) CreateUser(req request.RrCreateUser) response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	if err := r.doJSON("POST", r.base+"/users", req, &out); err != nil {
		out.Code = "CLIENT_USER_CREATE_001"
		out.Message = err.Error()
	}
	return out
}

func (r *userPrivateRepo) UpdateUser(id string, req request.RrUpdateUser) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/users/%s", r.base, id)
	if err := r.doJSON("PUT", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_USER_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (r *userPrivateRepo) DeleteUser(id string) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/users/%s", r.base, id)
	if err := r.doJSON("DELETE", url, nil, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_USER_DELETE_001", Message: err.Error()}
	}
	return out
}

// ---- public repo -----------------------------------------------------------

type userPublicRepo struct {
	base   string
	client *http.Client
}

// NewUserPublic creates a UserPublic repository backed by /v1/public (no auth required).
func NewUserPublic(conf config.BaseConfig) UserPublic {
	return &userPublicRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/public",
		client: &http.Client{},
	}
}

func (r *userPublicRepo) RegisterUser(req request.RrCreateUser) response.RrSingleIdPUser {
	var result response.RrSingleIdPUser
	b, err := json.Marshal(req)
	if err != nil {
		result.Code = "ANON_REGISTER_001"
		result.Message = "failed to encode request"
		return result
	}
	httpReq, err := http.NewRequest("POST", r.base+"/user", bytes.NewReader(b))
	if err != nil {
		result.Code = "ANON_REGISTER_002"
		result.Message = "failed to create request"
		return result
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(httpReq)
	if err != nil {
		result.Code = "ANON_REGISTER_003"
		result.Message = "request failed"
		return result
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.Code = "ANON_REGISTER_004"
		result.Message = "failed to decode response"
	}
	return result
}
