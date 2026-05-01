package usecase

import (
	"context"

	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
)

// Member is the usecase interface for group membership management via the external IdP.
type Member interface {
	ListGroupMembers(ctx context.Context, groupID string) ([]model.LoUser, error)
	AddUserToGroup(ctx context.Context, userID, groupID, role string) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID string) error
}

type memberUsecase struct {
	manager repository.IdPManager
}

// NewMember creates a new Member usecase backed by the given IdPManager.
func NewMember(manager repository.IdPManager) Member {
	return &memberUsecase{manager: manager}
}

func (rcvr *memberUsecase) ListGroupMembers(ctx context.Context, groupID string) ([]model.LoUser, error) {
	return rcvr.manager.ListGroupMembers(ctx, groupID)
}

func (rcvr *memberUsecase) AddUserToGroup(ctx context.Context, userID, groupID, role string) error {
	return rcvr.manager.AddUserToGroup(ctx, userID, groupID, role)
}

func (rcvr *memberUsecase) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	return rcvr.manager.RemoveUserFromGroup(ctx, userID, groupID)
}
