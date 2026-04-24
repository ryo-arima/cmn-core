package repository

import (
	"encoding/json"
	"net/http"

	clientauth "github.com/ryo-arima/cmn-core/pkg/client/auth"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// Common is the data-access interface for token-related server operations.
// Authentication is handled transparently via the auth.Manager.
type Common interface {
	Logout() response.Commons
	ValidateToken() response.ValidateToken
	GetUserInfo() response.Commons
}

type common struct {
	serverBase string
	client     *http.Client
}

// NewCommon creates a Common repository.
// manager is used to obtain and inject auth tokens automatically.
func NewCommon(conf config.BaseConfig, manager *clientauth.Manager) Common {
	return &common{
		serverBase: conf.YamlConfig.Application.Client.ServerEndpoint,
		client:     manager.HTTPClient(),
	}
}

func (r *common) Logout() response.Commons {
	var result response.Commons
	req, err := http.NewRequest("DELETE", r.serverBase+"/v1/share/token", nil)
	if err != nil {
		result.Code = "CLIENT_LOGOUT_001"
		result.Message = "failed to create request"
		return result
	}
	resp, err := r.client.Do(req)
	if err != nil {
		result.Code = "CLIENT_LOGOUT_002"
		result.Message = "request failed"
		return result
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.Code = "CLIENT_LOGOUT_003"
		result.Message = "failed to decode response"
	}
	return result
}

func (r *common) ValidateToken() response.ValidateToken {
	var result response.ValidateToken
	req, err := http.NewRequest("GET", r.serverBase+"/v1/share/token/validate", nil)
	if err != nil {
		result.Code = "CLIENT_VALIDATE_001"
		result.Message = "failed to create request"
		return result
	}
	resp, err := r.client.Do(req)
	if err != nil {
		result.Code = "CLIENT_VALIDATE_002"
		result.Message = "request failed"
		return result
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.Code = "CLIENT_VALIDATE_003"
		result.Message = "failed to decode response"
	}
	return result
}

func (r *common) GetUserInfo() response.Commons {
	var result response.Commons
	req, err := http.NewRequest("GET", r.serverBase+"/v1/share/token/userinfo", nil)
	if err != nil {
		result.Code = "CLIENT_USERINFO_001"
		result.Message = "failed to create request"
		return result
	}
	resp, err := r.client.Do(req)
	if err != nil {
		result.Code = "CLIENT_USERINFO_002"
		result.Message = "request failed"
		return result
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.Code = "CLIENT_USERINFO_003"
		result.Message = "failed to decode response"
	}
	return result
}

