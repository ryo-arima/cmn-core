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

func (u *userUsecase) GetUser(ctx context.Context, id string) (*model.LoUser, error) {
	return u.manager.GetUser(ctx, id)
}

func (u *userUsecase) ListUsers(ctx context.Context) ([]model.LoUser, error) {
	return u.manager.ListUsers(ctx)
}

func (u *userUsecase) CreateUser(ctx context.Context, input request.RrCreateUser) (*model.LoUser, error) {
	return u.manager.CreateUser(ctx, input)
}

func (u *userUsecase) UpdateUser(ctx context.Context, id string, input request.RrUpdateUser) error {
	return u.manager.UpdateUser(ctx, id, input)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	return u.manager.DeleteUser(ctx, id)
}

func (u *userUsecase) Login(ctx context.Context, username, password string) (string, error) {
	return u.manager.Login(ctx, username, password)
}
