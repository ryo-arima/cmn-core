package repository

import (
	"context"
	"fmt"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
)

// IdPManager is the unified interface for managing users and groups in an external
// identity provider (Keycloak or Casdoor).
// The active implementation is selected at startup based on
// Application.Server.IdP.Provider in app.yaml.
type IdPManager interface {
	// --- User management ---
	GetUser(ctx context.Context, id string) (*model.LoUser, error)
	ListUsers(ctx context.Context) ([]model.LoUser, error)
	CreateUser(ctx context.Context, input request.RrCreateUser) (*model.LoUser, error)
	UpdateUser(ctx context.Context, id string, input request.RrUpdateUser) error
	DeleteUser(ctx context.Context, id string) error

	// --- Group management ---
	GetGroup(ctx context.Context, id string) (*model.LoGroup, error)
	ListGroups(ctx context.Context) ([]model.LoGroup, error)
	CreateGroup(ctx context.Context, input request.RrCreateGroup) (*model.LoGroup, error)
	UpdateGroup(ctx context.Context, id string, input request.RrUpdateGroup) error
	DeleteGroup(ctx context.Context, id string) error

	// --- Group membership ---
	ListGroupMembers(ctx context.Context, groupID string) ([]model.LoUser, error)
	AddUserToGroup(ctx context.Context, userID, groupID, role string) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID string) error

	// --- Authentication ---
	// Login performs a Resource Owner Password Credentials (ROPC) grant and returns
	// the access token issued by the IdP.
	Login(ctx context.Context, username, password string) (string, error)
}

// NewIdPManager creates the IdPManager implementation selected by
// Application.Server.IdP.Provider in app.yaml.
// Returns an error if the provider name is unknown or required fields are missing.
func NewIdPManager(conf config.BaseConfig) (IdPManager, error) {
	idpCfg := conf.YamlConfig.Application.Server.IdP
	switch idpCfg.Provider {
	case "keycloak":
		if idpCfg.Keycloak.BaseURL == "" {
			return nil, fmt.Errorf("idp: keycloak.base_url must not be empty")
		}
		return newKeycloakManager(idpCfg.Keycloak), nil
	case "casdoor":
		if idpCfg.Casdoor.BaseURL == "" {
			return nil, fmt.Errorf("idp: casdoor.base_url must not be empty")
		}
		return newCasdoorManager(idpCfg.Casdoor), nil
	default:
		return nil, fmt.Errorf("idp: unknown provider %q (must be \"keycloak\" or \"casdoor\")", idpCfg.Provider)
	}
}
