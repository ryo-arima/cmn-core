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

// Resource is the data-access interface for resource operations.
type Resource interface {
	ListResources() response.Resources
	GetResource(uuid string) response.SingleResource
	CreateResource(req request.CreateResource) response.SingleResource
	UpdateResource(uuid string, req request.UpdateResource) response.Commons
	DeleteResource(uuid string) response.Commons
	GetResourceGroupRoles(uuid string) response.ResourceGroupRoles
	SetResourceGroupRole(uuid string, req request.SetResourceGroupRole) response.Commons
	DeleteResourceGroupRole(uuid, groupUUID string) response.Commons
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

func (r *resourceRepo) do(method, url string, body interface{}, out interface{}) error {
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

func (r *resourceRepo) ListResources() response.Resources {
	var out response.Resources
	if err := r.do("GET", r.base+"/resources", nil, &out); err != nil {
		out.Code = "CLIENT_RESOURCE_LIST_001"
		out.Message = err.Error()
	}
	return out
}

func (r *resourceRepo) GetResource(uuid string) response.SingleResource {
	var out response.SingleResource
	url := fmt.Sprintf("%s/resource?uuid=%s", r.base, uuid)
	if err := r.do("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_RESOURCE_GET_001"
		out.Message = err.Error()
	}
	return out
}

func (r *resourceRepo) CreateResource(req request.CreateResource) response.SingleResource {
	var out response.SingleResource
	if err := r.do("POST", r.base+"/resources", req, &out); err != nil {
		out.Code = "CLIENT_RESOURCE_CREATE_001"
		out.Message = err.Error()
	}
	return out
}

func (r *resourceRepo) UpdateResource(uuid string, req request.UpdateResource) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/resources/%s", r.base, uuid)
	if err := r.do("PUT", url, req, &out); err != nil {
		return response.Commons{Code: "CLIENT_RESOURCE_UPDATE_001", Message: err.Error()}
	}
	return out
}

func (r *resourceRepo) DeleteResource(uuid string) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/resources/%s", r.base, uuid)
	if err := r.do("DELETE", url, nil, &out); err != nil {
		return response.Commons{Code: "CLIENT_RESOURCE_DELETE_001", Message: err.Error()}
	}
	return out
}

func (r *resourceRepo) GetResourceGroupRoles(uuid string) response.ResourceGroupRoles {
	var out response.ResourceGroupRoles
	url := fmt.Sprintf("%s/resource/groups?uuid=%s", r.base, uuid)
	if err := r.do("GET", url, nil, &out); err != nil {
		out.Code = "CLIENT_RESOURCE_GROUPS_001"
		out.Message = err.Error()
	}
	return out
}

func (r *resourceRepo) SetResourceGroupRole(uuid string, req request.SetResourceGroupRole) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/resources/%s/groups", r.base, uuid)
	if err := r.do("PUT", url, req, &out); err != nil {
		return response.Commons{Code: "CLIENT_RESOURCE_SETGROUP_001", Message: err.Error()}
	}
	return out
}

func (r *resourceRepo) DeleteResourceGroupRole(uuid, groupUUID string) response.Commons {
	var out response.Commons
	url := fmt.Sprintf("%s/resources/%s/groups/%s", r.base, uuid, groupUUID)
	if err := r.do("DELETE", url, nil, &out); err != nil {
		return response.Commons{Code: "CLIENT_RESOURCE_DELGROUP_001", Message: err.Error()}
	}
	return out
}
