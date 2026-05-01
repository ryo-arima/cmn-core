package usecase

import (
	clientauth "github.com/ryo-arima/cmn-core/pkg/client/share"
	"github.com/ryo-arima/cmn-core/pkg/client/repository"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// User is the business-logic interface for app user operations (/v1/internal).
type User interface {
	GetMyUser() response.RrSingleIdPUser
	UpdateMyUser(req request.RrUpdateUser) response.RrCommons
	GetUser(id string) response.RrSingleIdPUser
	ListGroupUsers() response.RrIdPUsers
}

// UserAdmin extends User with admin-only CRUD (/v1/private).
type UserAdmin interface {
	User
	ListUsers() response.RrIdPUsers
	CreateUser(req request.RrCreateUser) response.RrSingleIdPUser
	UpdateUser(id string, req request.RrUpdateUser) response.RrCommons
	DeleteUser(id string) response.RrCommons
}

// Anonymous is the business-logic interface for unauthenticated user registration.
type Anonymous interface {
	RegisterUser(req request.RrCreateUser) response.RrSingleIdPUser
}

// ---- internal usecase -------------------------------------------------------

type userUsecase struct {
	repo repository.UserInternal
}

// NewUser creates a User usecase backed by /v1/internal.
func NewUser(conf config.BaseConfig, manager *clientauth.Manager) User {
	return &userUsecase{repo: repository.NewUserInternal(conf, manager)}
}

func (rcvr *userUsecase) GetMyUser() response.RrSingleIdPUser           { return rcvr.repo.GetMyUser() }
func (rcvr *userUsecase) UpdateMyUser(r request.RrUpdateUser) response.RrCommons {
	return rcvr.repo.UpdateMyUser(r)
}
func (rcvr *userUsecase) GetUser(id string) response.RrSingleIdPUser { return rcvr.repo.GetUser(id) }
func (rcvr *userUsecase) ListGroupUsers() response.RrIdPUsers        { return rcvr.repo.ListGroupUsers() }

// ---- admin usecase ----------------------------------------------------------

type userAdminUsecase struct {
	repo repository.UserPrivate
}

// NewUserAdmin creates a UserAdmin usecase backed by /v1/private.
func NewUserAdmin(conf config.BaseConfig, manager *clientauth.Manager) UserAdmin {
	return &userAdminUsecase{repo: repository.NewUserPrivate(conf, manager)}
}

func (rcvr *userAdminUsecase) GetMyUser() response.RrSingleIdPUser { return rcvr.repo.GetMyUser() }
func (rcvr *userAdminUsecase) UpdateMyUser(r request.RrUpdateUser) response.RrCommons {
	return rcvr.repo.UpdateMyUser(r)
}
func (rcvr *userAdminUsecase) GetUser(id string) response.RrSingleIdPUser { return rcvr.repo.GetUser(id) }
func (rcvr *userAdminUsecase) ListGroupUsers() response.RrIdPUsers        { return rcvr.repo.ListGroupUsers() }
func (rcvr *userAdminUsecase) ListUsers() response.RrIdPUsers             { return rcvr.repo.ListUsers() }
func (rcvr *userAdminUsecase) CreateUser(r request.RrCreateUser) response.RrSingleIdPUser {
	return rcvr.repo.CreateUser(r)
}
func (rcvr *userAdminUsecase) UpdateUser(id string, r request.RrUpdateUser) response.RrCommons {
	return rcvr.repo.UpdateUser(id, r)
}
func (rcvr *userAdminUsecase) DeleteUser(id string) response.RrCommons { return rcvr.repo.DeleteUser(id) }

// ---- anonymous usecase ------------------------------------------------------

type anonymousUsecase struct {
	repo repository.UserPublic
}

// NewAnonymous creates an Anonymous usecase backed by /v1/public.
func NewAnonymous(conf config.BaseConfig) Anonymous {
	return &anonymousUsecase{repo: repository.NewUserPublic(conf)}
}

func (rcvr *anonymousUsecase) RegisterUser(req request.RrCreateUser) response.RrSingleIdPUser {
	return rcvr.repo.RegisterUser(req)
}
