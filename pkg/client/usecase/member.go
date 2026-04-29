package usecase

import (
	clientauth "github.com/ryo-arima/cmn-core/pkg/client/share"
	"github.com/ryo-arima/cmn-core/pkg/client/repository"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// Member is the business-logic interface for group membership operations.
type Member interface {
	ListGroupMembers(groupID string) response.RrIdPUsers
	AddGroupMember(groupID string, req request.RrAddGroupMember) response.RrCommons
	RemoveGroupMember(groupID string, req request.RrRemoveGroupMember) response.RrCommons
}

// MemberAdmin is the same interface backed by /v1/private.
type MemberAdmin interface {
	Member
}

// ---- internal usecase -------------------------------------------------------

type memberUsecase struct {
	repo repository.MemberInternal
}

// NewMember creates a Member usecase backed by /v1/internal.
func NewMember(conf config.BaseConfig, manager *clientauth.Manager) Member {
	return &memberUsecase{repo: repository.NewMemberInternal(conf, manager)}
}

func (u *memberUsecase) ListGroupMembers(gid string) response.RrIdPUsers {
	return u.repo.ListGroupMembers(gid)
}
func (u *memberUsecase) AddGroupMember(gid string, r request.RrAddGroupMember) response.RrCommons {
	return u.repo.AddGroupMember(gid, r)
}
func (u *memberUsecase) RemoveGroupMember(gid string, r request.RrRemoveGroupMember) response.RrCommons {
	return u.repo.RemoveGroupMember(gid, r)
}

// ---- admin usecase ----------------------------------------------------------

type memberAdminUsecase struct {
	repo repository.MemberPrivate
}

// NewMemberAdmin creates a MemberAdmin usecase backed by /v1/private.
func NewMemberAdmin(conf config.BaseConfig, manager *clientauth.Manager) MemberAdmin {
	return &memberAdminUsecase{repo: repository.NewMemberPrivate(conf, manager)}
}

func (u *memberAdminUsecase) ListGroupMembers(gid string) response.RrIdPUsers {
	return u.repo.ListGroupMembers(gid)
}
func (u *memberAdminUsecase) AddGroupMember(gid string, r request.RrAddGroupMember) response.RrCommons {
	return u.repo.AddGroupMember(gid, r)
}
func (u *memberAdminUsecase) RemoveGroupMember(gid string, r request.RrRemoveGroupMember) response.RrCommons {
	return u.repo.RemoveGroupMember(gid, r)
}
