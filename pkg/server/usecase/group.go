package usecase

import (
	"context"

	"github.com/google/uuid"
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
	// groupStore is non-nil for Casdoor, which has no native group UUIDs.
	// In that mode the IdP group Name is a generated UUID, and groupStore stores the display name.
	groupStore repository.Group
}

// NewGroup creates a new Group usecase.
// Pass a non-nil groupStore when using Casdoor so that display names are persisted in PostgreSQL.
func NewGroup(manager repository.IdPManager, groupStore repository.Group) Group {
	return &groupUsecase{manager: manager, groupStore: groupStore}
}

func (rcvr *groupUsecase) GetGroup(ctx context.Context, id string) (*model.LoGroup, error) {
	g, err := rcvr.manager.GetGroup(ctx, id)
	if err != nil {
		return nil, err
	}
	if rcvr.groupStore != nil {
		g.Name = rcvr.groupStore.LookupName(ctx, g.ID)
	}
	return g, nil
}

func (rcvr *groupUsecase) ListGroups(ctx context.Context) ([]model.LoGroup, error) {
	groups, err := rcvr.manager.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	if rcvr.groupStore != nil {
		ids := make([]string, len(groups))
		for i, g := range groups {
			ids[i] = g.ID
		}
		names := rcvr.groupStore.LookupNames(ctx, ids)
		for i := range groups {
			if n, ok := names[groups[i].ID]; ok {
				groups[i].Name = n
			}
		}
	}
	return groups, nil
}

func (rcvr *groupUsecase) CreateGroup(ctx context.Context, input request.RrCreateGroup) (*model.LoGroup, error) {
	if rcvr.groupStore != nil {
		// Casdoor: use a generated UUID as the IdP-internal name to ensure uniqueness.
		id := uuid.New().String()
		g, err := rcvr.manager.CreateGroup(ctx, request.RrCreateGroup{Name: id})
		if err != nil {
			return nil, err
		}
		_ = rcvr.groupStore.Upsert(ctx, id, input.Name)
		g.Name = input.Name
		return g, nil
	}
	return rcvr.manager.CreateGroup(ctx, input)
}

func (rcvr *groupUsecase) UpdateGroup(ctx context.Context, id string, input request.RrUpdateGroup) error {
	if rcvr.groupStore != nil {
		// Casdoor: the IdP Name (UUID) is immutable; only update the display name in psql.
		return rcvr.groupStore.Upsert(ctx, id, input.Name)
	}
	return rcvr.manager.UpdateGroup(ctx, id, input)
}

func (rcvr *groupUsecase) DeleteGroup(ctx context.Context, id string) error {
	if err := rcvr.manager.DeleteGroup(ctx, id); err != nil {
		return err
	}
	if rcvr.groupStore != nil {
		_ = rcvr.groupStore.SoftDelete(ctx, id)
	}
	return nil
}
