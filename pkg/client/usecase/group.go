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

func (rcvr *groupUsecase) ListMyGroups() response.RrIdPGroups { return rcvr.repo.ListMyGroups() }
func (rcvr *groupUsecase) GetGroup(id string) response.RrSingleIdPGroup { return rcvr.repo.GetGroup(id) }
func (rcvr *groupUsecase) CreateGroup(r request.RrCreateGroup) response.RrSingleIdPGroup {
	return rcvr.repo.CreateGroup(r)
}
func (rcvr *groupUsecase) UpdateGroup(id string, r request.RrUpdateGroup) response.RrCommons {
	return rcvr.repo.UpdateGroup(id, r)
}
func (rcvr *groupUsecase) DeleteGroup(id string) response.RrCommons { return rcvr.repo.DeleteGroup(id) }

// ---- admin usecase ----------------------------------------------------------

type groupAdminUsecase struct {
	repo repository.GroupPrivate
}

// NewGroupAdmin creates a GroupAdmin usecase backed by /v1/private.
func NewGroupAdmin(conf config.BaseConfig, manager *clientauth.Manager) GroupAdmin {
	return &groupAdminUsecase{repo: repository.NewGroupPrivate(conf, manager)}
}

func (rcvr *groupAdminUsecase) ListMyGroups() response.RrIdPGroups { return rcvr.repo.ListMyGroups() }
func (rcvr *groupAdminUsecase) ListGroups() response.RrIdPGroups   { return rcvr.repo.ListGroups() }
func (rcvr *groupAdminUsecase) GetGroup(id string) response.RrSingleIdPGroup {
	return rcvr.repo.GetGroup(id)
}
func (rcvr *groupAdminUsecase) CreateGroup(r request.RrCreateGroup) response.RrSingleIdPGroup {
	return rcvr.repo.CreateGroup(r)
}
func (rcvr *groupAdminUsecase) UpdateGroup(id string, r request.RrUpdateGroup) response.RrCommons {
	return rcvr.repo.UpdateGroup(id, r)
}
func (rcvr *groupAdminUsecase) DeleteGroup(id string) response.RrCommons { return rcvr.repo.DeleteGroup(id) }
