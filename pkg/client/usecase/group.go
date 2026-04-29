package usecase

import (
	clientauth "github.com/ryo-arima/cmn-core/pkg/client/share"
	"github.com/ryo-arima/cmn-core/pkg/client/repository"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// Group is the business-logic interface for app group operations (/v1/internal).
type Group interface {
	ListMyGroups() response.RrIdPGroups
	GetGroup(id string) response.RrSingleIdPGroup
	CreateGroup(req request.RrCreateGroup) response.RrSingleIdPGroup
	UpdateGroup(id string, req request.RrUpdateGroup) response.RrCommons
	DeleteGroup(id string) response.RrCommons
}

// GroupAdmin extends Group with admin-scope group listing (/v1/private).
type GroupAdmin interface {
	Group
	ListGroups() response.RrIdPGroups
}

// ---- internal usecase -------------------------------------------------------

type groupUsecase struct {
	repo repository.GroupInternal
}

// NewGroup creates a Group usecase backed by /v1/internal.
func NewGroup(conf config.BaseConfig, manager *clientauth.Manager) Group {
	return &groupUsecase{repo: repository.NewGroupInternal(conf, manager)}
}

func (u *groupUsecase) ListMyGroups() response.RrIdPGroups { return u.repo.ListMyGroups() }
func (u *groupUsecase) GetGroup(id string) response.RrSingleIdPGroup { return u.repo.GetGroup(id) }
func (u *groupUsecase) CreateGroup(r request.RrCreateGroup) response.RrSingleIdPGroup {
	return u.repo.CreateGroup(r)
}
func (u *groupUsecase) UpdateGroup(id string, r request.RrUpdateGroup) response.RrCommons {
	return u.repo.UpdateGroup(id, r)
}
func (u *groupUsecase) DeleteGroup(id string) response.RrCommons { return u.repo.DeleteGroup(id) }

// ---- admin usecase ----------------------------------------------------------

type groupAdminUsecase struct {
	repo repository.GroupPrivate
}

// NewGroupAdmin creates a GroupAdmin usecase backed by /v1/private.
func NewGroupAdmin(conf config.BaseConfig, manager *clientauth.Manager) GroupAdmin {
	return &groupAdminUsecase{repo: repository.NewGroupPrivate(conf, manager)}
}

func (u *groupAdminUsecase) ListMyGroups() response.RrIdPGroups { return u.repo.ListMyGroups() }
func (u *groupAdminUsecase) ListGroups() response.RrIdPGroups   { return u.repo.ListGroups() }
func (u *groupAdminUsecase) GetGroup(id string) response.RrSingleIdPGroup {
	return u.repo.GetGroup(id)
}
func (u *groupAdminUsecase) CreateGroup(r request.RrCreateGroup) response.RrSingleIdPGroup {
	return u.repo.CreateGroup(r)
}
func (u *groupAdminUsecase) UpdateGroup(id string, r request.RrUpdateGroup) response.RrCommons {
	return u.repo.UpdateGroup(id, r)
}
func (u *groupAdminUsecase) DeleteGroup(id string) response.RrCommons { return u.repo.DeleteGroup(id) }
