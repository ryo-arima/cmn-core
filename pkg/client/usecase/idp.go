package usecase

import (
	clientauth "github.com/ryo-arima/cmn-core/pkg/client/auth"
	"github.com/ryo-arima/cmn-core/pkg/client/repository"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// IdP is the business-logic interface for identity-provider operations (app user).
type IdP interface {
	GetMyUser() response.SingleIdPUser
	UpdateMyUser(req request.UpdateUser) response.Commons

	// Any user by ID
	GetUser(id string) response.SingleIdPUser
	// Users in caller's groups
	ListGroupUsers() response.IdPUsers

	ListMyGroups() response.IdPGroups
	GetGroup(id string) response.SingleIdPGroup
	CreateGroup(req request.CreateGroup) response.SingleIdPGroup
	UpdateGroup(id string, req request.UpdateGroup) response.Commons
	DeleteGroup(id string) response.Commons

	ListGroupMembers(groupID string) response.IdPUsers
	AddGroupMember(groupID string, req request.AddGroupMember) response.Commons
	RemoveGroupMember(groupID string, req request.RemoveGroupMember) response.Commons
}

// IdPAdmin extends IdP with admin-only user management.
type IdPAdmin interface {
	IdP
	ListUsers() response.IdPUsers
	CreateUser(req request.CreateUser) response.SingleIdPUser
	UpdateUser(id string, req request.UpdateUser) response.Commons
	DeleteUser(id string) response.Commons
	ListGroups() response.IdPGroups
}

type idpUsecase struct {
	repo repository.IdP
}

type idpAdminUsecase struct {
	repo repository.IdPAdmin
}

// NewIdP creates an IdP usecase for the app client (/v1/internal).
func NewIdP(conf config.BaseConfig, manager *clientauth.Manager) IdP {
	return &idpUsecase{repo: repository.NewIdP(conf, manager)}
}

// NewIdPAdmin creates an IdPAdmin usecase for the admin client (/v1/private).
func NewIdPAdmin(conf config.BaseConfig, manager *clientauth.Manager) IdPAdmin {
	return &idpAdminUsecase{repo: repository.NewIdPAdmin(conf, manager)}
}

// ── IdP (app) ─────────────────────────────────────────────────────────────

func (u *idpUsecase) GetMyUser() response.SingleIdPUser           { return u.repo.GetMyUser() }
func (u *idpUsecase) UpdateMyUser(r request.UpdateUser) response.Commons {
	return u.repo.UpdateMyUser(r)
}
func (u *idpUsecase) GetUser(id string) response.SingleIdPUser { return u.repo.GetUser(id) }
func (u *idpUsecase) ListGroupUsers() response.IdPUsers        { return u.repo.ListGroupUsers() }
func (u *idpUsecase) ListMyGroups() response.IdPGroups { return u.repo.ListMyGroups() }
func (u *idpUsecase) GetGroup(id string) response.SingleIdPGroup { return u.repo.GetGroup(id) }
func (u *idpUsecase) CreateGroup(r request.CreateGroup) response.SingleIdPGroup {
	return u.repo.CreateGroup(r)
}
func (u *idpUsecase) UpdateGroup(id string, r request.UpdateGroup) response.Commons {
	return u.repo.UpdateGroup(id, r)
}
func (u *idpUsecase) DeleteGroup(id string) response.Commons { return u.repo.DeleteGroup(id) }
func (u *idpUsecase) ListGroupMembers(gid string) response.IdPUsers {
	return u.repo.ListGroupMembers(gid)
}
func (u *idpUsecase) AddGroupMember(gid string, r request.AddGroupMember) response.Commons {
	return u.repo.AddGroupMember(gid, r)
}
func (u *idpUsecase) RemoveGroupMember(gid string, r request.RemoveGroupMember) response.Commons {
	return u.repo.RemoveGroupMember(gid, r)
}

// ── IdPAdmin ──────────────────────────────────────────────────────────────

func (u *idpAdminUsecase) GetMyUser() response.SingleIdPUser { return u.repo.GetMyUser() }
func (u *idpAdminUsecase) UpdateMyUser(r request.UpdateUser) response.Commons {
	return u.repo.UpdateMyUser(r)
}
func (u *idpAdminUsecase) GetUser(id string) response.SingleIdPUser { return u.repo.GetUser(id) }
func (u *idpAdminUsecase) ListGroupUsers() response.IdPUsers        { return u.repo.ListGroupUsers() }
func (u *idpAdminUsecase) ListMyGroups() response.IdPGroups { return u.repo.ListMyGroups() }
func (u *idpAdminUsecase) GetGroup(id string) response.SingleIdPGroup { return u.repo.GetGroup(id) }
func (u *idpAdminUsecase) CreateGroup(r request.CreateGroup) response.SingleIdPGroup {
	return u.repo.CreateGroup(r)
}
func (u *idpAdminUsecase) UpdateGroup(id string, r request.UpdateGroup) response.Commons {
	return u.repo.UpdateGroup(id, r)
}
func (u *idpAdminUsecase) DeleteGroup(id string) response.Commons { return u.repo.DeleteGroup(id) }
func (u *idpAdminUsecase) ListGroupMembers(gid string) response.IdPUsers {
	return u.repo.ListGroupMembers(gid)
}
func (u *idpAdminUsecase) AddGroupMember(gid string, r request.AddGroupMember) response.Commons {
	return u.repo.AddGroupMember(gid, r)
}
func (u *idpAdminUsecase) RemoveGroupMember(gid string, r request.RemoveGroupMember) response.Commons {
	return u.repo.RemoveGroupMember(gid, r)
}
func (u *idpAdminUsecase) ListUsers() response.IdPUsers       { return u.repo.ListUsers() }
func (u *idpAdminUsecase) CreateUser(r request.CreateUser) response.SingleIdPUser {
	return u.repo.CreateUser(r)
}
func (u *idpAdminUsecase) UpdateUser(id string, r request.UpdateUser) response.Commons {
	return u.repo.UpdateUser(id, r)
}
func (u *idpAdminUsecase) DeleteUser(id string) response.Commons { return u.repo.DeleteUser(id) }
func (u *idpAdminUsecase) ListGroups() response.IdPGroups        { return u.repo.ListGroups() }
