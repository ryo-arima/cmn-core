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

// MemberInternal is the data-access interface for group membership operations via /v1/internal.
type MemberInternal interface {
	ListGroupMembers(groupID string) response.RrIdPUsers
	AddGroupMember(groupID string, req request.RrAddGroupMember) response.RrCommons
	RemoveGroupMember(groupID string, req request.RrRemoveGroupMember) response.RrCommons
}

// MemberPrivate is the same interface backed by /v1/private.
type MemberPrivate interface {
	MemberInternal
}

// ---- internal repo ---------------------------------------------------------

type memberInternalRepo struct {
	base   string
	client *http.Client
}

// NewMemberInternal creates a MemberInternal repository backed by /v1/internal.
func NewMemberInternal(conf config.BaseConfig, manager *clientauth.Manager) MemberInternal {
	return &memberInternalRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal",
		client: manager.HTTPClient(),
	}
}

func (rcvr *memberInternalRepo) doJSON(method, url string, body interface{}, out interface{}) error {
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

func (rcvr *memberInternalRepo) ListGroupMembers(groupID string) response.RrIdPUsers {
	var out response.RrIdPUsers
	url := fmt.Sprintf("%s/members?group_id=%s", rcvr.base, groupID)
	if err := rcvr.doJSON("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_MEMBER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *memberInternalRepo) AddGroupMember(groupID string, req request.RrAddGroupMember) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/member/%s", rcvr.base, groupID)
	if err := rcvr.doJSON("POST", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_MEMBER_ADD_001", Message: err.Error()}
	}
	return out
}

func (rcvr *memberInternalRepo) RemoveGroupMember(groupID string, req request.RrRemoveGroupMember) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/member/%s", rcvr.base, groupID)
	if err := rcvr.doJSON("DELETE", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_MEMBER_REMOVE_001", Message: err.Error()}
	}
	return out
}

// ---- private repo ----------------------------------------------------------

type memberPrivateRepo struct {
	base   string
	client *http.Client
}

// NewMemberPrivate creates a MemberPrivate repository backed by /v1/private.
func NewMemberPrivate(conf config.BaseConfig, manager *clientauth.Manager) MemberPrivate {
	return &memberPrivateRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/private",
		client: manager.HTTPClient(),
	}
}

func (rcvr *memberPrivateRepo) doJSON(method, url string, body interface{}, out interface{}) error {
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

func (rcvr *memberPrivateRepo) ListGroupMembers(groupID string) response.RrIdPUsers {
	var out response.RrIdPUsers
	url := fmt.Sprintf("%s/members?group_id=%s", rcvr.base, groupID)
	if err := rcvr.doJSON("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_MEMBER_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *memberPrivateRepo) AddGroupMember(groupID string, req request.RrAddGroupMember) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/member/%s", rcvr.base, groupID)
	if err := rcvr.doJSON("POST", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_MEMBER_ADD_001", Message: err.Error()}
	}
	return out
}

func (rcvr *memberPrivateRepo) RemoveGroupMember(groupID string, req request.RrRemoveGroupMember) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/member/%s", rcvr.base, groupID)
	if err := rcvr.doJSON("DELETE", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_MEMBER_REMOVE_001", Message: err.Error()}
	}
	return out
}
