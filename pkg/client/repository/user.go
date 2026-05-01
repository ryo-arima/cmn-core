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

func (rcvr *userInternalRepo) doJSON(method, url string, body interface{}, out interface{}) error {
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
	resp, err := rcvr.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

func (rcvr *userInternalRepo) GetMyUser() response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	if err := rcvr.doJSON("GET", rcvr.base+"/user", nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *userInternalRepo) UpdateMyUser(req request.RrUpdateUser) response.RrCommons {
	var out response.RrCommons
	if err := rcvr.doJSON("PUT", rcvr.base+"/user", req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_USER_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (rcvr *userInternalRepo) GetUser(id string) response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	url := fmt.Sprintf("%s/user?id=%s", rcvr.base, id)
	if err := rcvr.doJSON("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *userInternalRepo) ListGroupUsers() response.RrIdPUsers {
	var out response.RrIdPUsers
	if err := rcvr.doJSON("GET", rcvr.base+"/users", nil, &out); err != nil {
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

func (rcvr *userPrivateRepo) doJSON(method, url string, body interface{}, out interface{}) error {
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
	resp, err := rcvr.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

func (rcvr *userPrivateRepo) GetMyUser() response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	if err := rcvr.doJSON("GET", rcvr.base+"/user", nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *userPrivateRepo) UpdateMyUser(req request.RrUpdateUser) response.RrCommons {
	var out response.RrCommons
	if err := rcvr.doJSON("PUT", rcvr.base+"/user", req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_USER_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (rcvr *userPrivateRepo) GetUser(id string) response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	url := fmt.Sprintf("%s/user?id=%s", rcvr.base, id)
	if err := rcvr.doJSON("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_USER_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *userPrivateRepo) ListGroupUsers() response.RrIdPUsers {
	var out response.RrIdPUsers
	if err := rcvr.doJSON("GET", rcvr.base+"/users", nil, &out); err != nil {
		out.Code = "CLIENT_USER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *userPrivateRepo) ListUsers() response.RrIdPUsers {
	var out response.RrIdPUsers
	if err := rcvr.doJSON("GET", rcvr.base+"/users", nil, &out); err != nil {
		out.Code = "CLIENT_USER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *userPrivateRepo) CreateUser(req request.RrCreateUser) response.RrSingleIdPUser {
	var out response.RrSingleIdPUser
	if err := rcvr.doJSON("POST", rcvr.base+"/users", req, &out); err != nil {
		out.Code = "CLIENT_USER_CREATE_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *userPrivateRepo) UpdateUser(id string, req request.RrUpdateUser) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/users/%s", rcvr.base, id)
	if err := rcvr.doJSON("PUT", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_USER_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (rcvr *userPrivateRepo) DeleteUser(id string) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/users/%s", rcvr.base, id)
	if err := rcvr.doJSON("DELETE", url, nil, &out); err != nil {
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

func (rcvr *userPublicRepo) RegisterUser(req request.RrCreateUser) response.RrSingleIdPUser {
	var result response.RrSingleIdPUser
	b, err := json.Marshal(req)
	if err != nil {
		result.Code = "ANON_REGISTER_001"
		result.Message = "failed to encode request"
		return result
	}
	httpReq, err := http.NewRequest("POST", rcvr.base+"/user", bytes.NewReader(b))
	if err != nil {
		result.Code = "ANON_REGISTER_002"
		result.Message = "failed to create request"
		return result
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := rcvr.client.Do(httpReq)
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
