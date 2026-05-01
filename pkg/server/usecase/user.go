package usecase

import (
	"context"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
)

// User is the usecase interface for user management via the external IdP.
type User interface {
	GetUser(ctx context.Context, id string) (*model.LoUser, error)
	ListUsers(ctx context.Context) ([]model.LoUser, error)
	CreateUser(ctx context.Context, input request.RrCreateUser) (*model.LoUser, error)
	UpdateUser(ctx context.Context, id string, input request.RrUpdateUser) error
	DeleteUser(ctx context.Context, id string) error
	Login(ctx context.Context, username, password string) (string, error)
}

type userUsecase struct {
	manager repository.IdPManager
}

// NewUser creates a new User usecase backed by the given IdPManager.
func NewUser(manager repository.IdPManager) User {
	return &userUsecase{manager: manager}
}

func (rcvr *userUsecase) GetUser(ctx context.Context, id string) (*model.LoUser, error) {
	return rcvr.manager.GetUser(ctx, id)
}

func (rcvr *userUsecase) ListUsers(ctx context.Context) ([]model.LoUser, error) {
	return rcvr.manager.ListUsers(ctx)
}

func (rcvr *userUsecase) CreateUser(ctx context.Context, input request.RrCreateUser) (*model.LoUser, error) {
	return rcvr.manager.CreateUser(ctx, input)
}

func (rcvr *userUsecase) UpdateUser(ctx context.Context, id string, input request.RrUpdateUser) error {
	return rcvr.manager.UpdateUser(ctx, id, input)
}

func (rcvr *userUsecase) DeleteUser(ctx context.Context, id string) error {
	return rcvr.manager.DeleteUser(ctx, id)
}

func (rcvr *userUsecase) Login(ctx context.Context, username, password string) (string, error) {
	return rcvr.manager.Login(ctx, username, password)
}
