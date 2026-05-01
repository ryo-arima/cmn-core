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

// Resource is the data-access interface for resource operations.
type Resource interface {
	ListResources() response.RrResources
	GetResource(uuid string) response.RrSingleResource
	CreateResource(req request.RrCreateResource) response.RrSingleResource
	UpdateResource(uuid string, req request.RrUpdateResource) response.RrCommons
	DeleteResource(uuid string) response.RrCommons
	GetResourceGroupRoles(uuid string) response.RrResourceGroupRoles
	SetResourceGroupRole(uuid string, req request.RrSetResourceGroupRole) response.RrCommons
	DeleteResourceGroupRole(uuid, groupID string) response.RrCommons
}

type resourceRepo struct {
	base   string
	client *http.Client
}

// NewResource creates a Resource repository backed by /v1/internal.
func NewResource(conf config.BaseConfig, manager *clientauth.Manager) Resource {
	return &resourceRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal",
		client: manager.HTTPClient(),
	}
}

// NewResourceAdmin creates a Resource repository backed by /v1/private.
func NewResourceAdmin(conf config.BaseConfig, manager *clientauth.Manager) Resource {
	return &resourceRepo{
		base:   conf.YamlConfig.Application.Client.ServerEndpoint + "/v1/private",
		client: manager.HTTPClient(),
	}
}

func (rcvr *resourceRepo) do(method, url string, body interface{}, out interface{}) error {
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

func (rcvr *resourceRepo) ListResources() response.RrResources {
	var out response.RrResources
	if err := rcvr.do("GET", rcvr.base+"/resources", nil, &out); err != nil {
		out.Code = "CLIENT_RESOURCE_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *resourceRepo) GetResource(uuid string) response.RrSingleResource {
	var out response.RrSingleResource
	url := fmt.Sprintf("%s/resource?uuid=%s", rcvr.base, uuid)
	if err := rcvr.do("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_RESOURCE_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *resourceRepo) CreateResource(req request.RrCreateResource) response.RrSingleResource {
	var out response.RrSingleResource
	if err := rcvr.do("POST", rcvr.base+"/resources", req, &out); err != nil {
		out.Code = "CLIENT_RESOURCE_CREATE_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *resourceRepo) UpdateResource(uuid string, req request.RrUpdateResource) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/resources/%s", rcvr.base, uuid)
	if err := rcvr.do("PUT", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_RESOURCE_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (rcvr *resourceRepo) DeleteResource(uuid string) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/resources/%s", rcvr.base, uuid)
	if err := rcvr.do("DELETE", url, nil, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_RESOURCE_DELETE_001", Message: err.Error()}
	}
	return out
}

func (rcvr *resourceRepo) GetResourceGroupRoles(uuid string) response.RrResourceGroupRoles {
	var out response.RrResourceGroupRoles
	url := fmt.Sprintf("%s/resource/groups?uuid=%s", rcvr.base, uuid)
	if err := rcvr.do("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_RESOURCE_GROUPS_001"
		out.Message = err.Error()
	}
	return out
}

func (rcvr *resourceRepo) SetResourceGroupRole(uuid string, req request.RrSetResourceGroupRole) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/resources/%s/groups", rcvr.base, uuid)
	if err := rcvr.do("PUT", url, req, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_RESOURCE_SETGROUP_001", Message: err.Error()}
	}
	return out
}

func (rcvr *resourceRepo) DeleteResourceGroupRole(uuid, groupID string) response.RrCommons {
	var out response.RrCommons
	url := fmt.Sprintf("%s/resources/%s/groups/%s", rcvr.base, uuid, groupID)
	if err := rcvr.do("DELETE", url, nil, &out); err != nil {
		return response.RrCommons{Code: "CLIENT_RESOURCE_DELGROUP_001", Message: err.Error()}
	}
	return out
}
