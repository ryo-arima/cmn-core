package usecase

import (
	"context"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
)

// IdP is the usecase interface for managing users and groups via an external
// identity provider (Keycloak or Casdoor).
type IdP interface {
	// User management
	GetUser(ctx context.Context, id string) (*model.IdPUser, error)
	ListUsers(ctx context.Context) ([]model.IdPUser, error)
	CreateUser(ctx context.Context, input request.CreateUser) (*model.IdPUser, error)
	UpdateUser(ctx context.Context, id string, input request.UpdateUser) error
	DeleteUser(ctx context.Context, id string) error

	// Group management
	GetGroup(ctx context.Context, id string) (*model.IdPGroup, error)
	ListGroups(ctx context.Context) ([]model.IdPGroup, error)
	CreateGroup(ctx context.Context, input request.CreateGroup) (*model.IdPGroup, error)
	UpdateGroup(ctx context.Context, id string, input request.UpdateGroup) error
	DeleteGroup(ctx context.Context, id string) error

	// Group membership
	ListGroupMembers(ctx context.Context, groupID string) ([]model.IdPUser, error)
	AddUserToGroup(ctx context.Context, userID, groupID string) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID string) error

	// Authentication
	Login(ctx context.Context, username, password string) (string, error)
}

type idpUsecase struct {
	manager repository.IdPManager
}

// NewIdP creates a new IdP usecase backed by the given IdPManager.
func NewIdP(manager repository.IdPManager) IdP {
	return &idpUsecase{manager: manager}
}

func (u *idpUsecase) GetUser(ctx context.Context, id string) (*model.IdPUser, error) {
	return u.manager.GetUser(ctx, id)
}

func (u *idpUsecase) ListUsers(ctx context.Context) ([]model.IdPUser, error) {
	return u.manager.ListUsers(ctx)
}

func (u *idpUsecase) CreateUser(ctx context.Context, input request.CreateUser) (*model.IdPUser, error) {
	return u.manager.CreateUser(ctx, input)
}

func (u *idpUsecase) UpdateUser(ctx context.Context, id string, input request.UpdateUser) error {
	return u.manager.UpdateUser(ctx, id, input)
}

func (u *idpUsecase) DeleteUser(ctx context.Context, id string) error {
	return u.manager.DeleteUser(ctx, id)
}

func (u *idpUsecase) GetGroup(ctx context.Context, id string) (*model.IdPGroup, error) {
	return u.manager.GetGroup(ctx, id)
}

func (u *idpUsecase) ListGroups(ctx context.Context) ([]model.IdPGroup, error) {
	return u.manager.ListGroups(ctx)
}

func (u *idpUsecase) CreateGroup(ctx context.Context, input request.CreateGroup) (*model.IdPGroup, error) {
	return u.manager.CreateGroup(ctx, input)
}

func (u *idpUsecase) UpdateGroup(ctx context.Context, id string, input request.UpdateGroup) error {
	return u.manager.UpdateGroup(ctx, id, input)
}

func (u *idpUsecase) DeleteGroup(ctx context.Context, id string) error {
	return u.manager.DeleteGroup(ctx, id)
}

func (u *idpUsecase) ListGroupMembers(ctx context.Context, groupID string) ([]model.IdPUser, error) {
	return u.manager.ListGroupMembers(ctx, groupID)
}

func (u *idpUsecase) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	return u.manager.AddUserToGroup(ctx, userID, groupID)
}

func (u *idpUsecase) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	return u.manager.RemoveUserFromGroup(ctx, userID, groupID)
}

func (u *idpUsecase) Login(ctx context.Context, username, password string) (string, error) {
	return u.manager.Login(ctx, username, password)
}
