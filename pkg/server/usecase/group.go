package usecase

import (
	"context"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
)

// Group is the usecase interface for group management via the external IdP.
type Group interface {
	GetGroup(ctx context.Context, id string) (*model.LoGroup, error)
	ListGroups(ctx context.Context) ([]model.LoGroup, error)
	CreateGroup(ctx context.Context, input request.RrCreateGroup) (*model.LoGroup, error)
	UpdateGroup(ctx context.Context, id string, input request.RrUpdateGroup) error
	DeleteGroup(ctx context.Context, id string) error
}

type groupUsecase struct {
	manager repository.IdPManager
}

// NewGroup creates a new Group usecase backed by the given IdPManager.
func NewGroup(manager repository.IdPManager) Group {
	return &groupUsecase{manager: manager}
}

func (u *groupUsecase) GetGroup(ctx context.Context, id string) (*model.LoGroup, error) {
	return u.manager.GetGroup(ctx, id)
}

func (u *groupUsecase) ListGroups(ctx context.Context) ([]model.LoGroup, error) {
	return u.manager.ListGroups(ctx)
}

func (u *groupUsecase) CreateGroup(ctx context.Context, input request.RrCreateGroup) (*model.LoGroup, error) {
	return u.manager.CreateGroup(ctx, input)
}

func (u *groupUsecase) UpdateGroup(ctx context.Context, id string, input request.RrUpdateGroup) error {
	return u.manager.UpdateGroup(ctx, id, input)
}

func (u *groupUsecase) DeleteGroup(ctx context.Context, id string) error {
	return u.manager.DeleteGroup(ctx, id)
}
