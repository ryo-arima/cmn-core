package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
)

// --- User usecase ---

type UserUC interface {
	GetUsers(c *gin.Context) []model.Users
	CreateUser(c *gin.Context, user model.Users) model.Users
	UpdateUser(c *gin.Context, user model.Users) model.Users
	DeleteUser(c *gin.Context, user model.Users) model.Users
	ListUsers(c *gin.Context, filter repository.UserQueryFilter) ([]model.Users, error)
}

type userUsecase struct {
	repo repository.User
}

// NewUser creates a User usecase backed by the given repository.
func NewUser(repo repository.User) UserUC {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) GetUsers(c *gin.Context) []model.Users {
	return u.repo.GetUsers(c)
}

func (u *userUsecase) CreateUser(c *gin.Context, user model.Users) model.Users {
	return u.repo.CreateUser(c, user)
}

func (u *userUsecase) UpdateUser(c *gin.Context, user model.Users) model.Users {
	return u.repo.UpdateUser(c, user)
}

func (u *userUsecase) DeleteUser(c *gin.Context, user model.Users) model.Users {
	return u.repo.DeleteUser(c, user)
}

func (u *userUsecase) ListUsers(c *gin.Context, filter repository.UserQueryFilter) ([]model.Users, error) {
	return u.repo.ListUsers(c, filter)
}

// --- Group usecase ---

// GroupRoleEnforcer is an optional interface for group-based role enforcement.
// Pass nil if not needed.
type GroupRoleEnforcer interface {
	Enforce(sub, obj, act string) (bool, error)
}

type GroupUC interface {
	GetGroups(c *gin.Context) []model.Groups
	GetGroupByUUID(c *gin.Context, uuid string) (model.Groups, error)
	CreateGroup(c *gin.Context, group *model.Groups) error
	UpdateGroup(c *gin.Context, group *model.Groups) error
	DeleteGroup(c *gin.Context, uuid string) error
	ListGroups(c *gin.Context, filter repository.GroupQueryFilter) ([]model.Groups, error)
}

type groupUsecase struct {
	groupRepo  repository.Group
	memberRepo repository.Member
	enforcer   GroupRoleEnforcer
}

// NewGroup creates a Group usecase. enforcer may be nil to disable role-based checks.
func NewGroup(groupRepo repository.Group, memberRepo repository.Member, enforcer GroupRoleEnforcer) GroupUC {
	return &groupUsecase{groupRepo: groupRepo, memberRepo: memberRepo, enforcer: enforcer}
}

func (u *groupUsecase) GetGroups(c *gin.Context) []model.Groups {
	return u.groupRepo.GetGroups(c)
}

func (u *groupUsecase) GetGroupByUUID(c *gin.Context, uuid string) (model.Groups, error) {
	return u.groupRepo.GetGroupByUUID(c, uuid)
}

func (u *groupUsecase) CreateGroup(c *gin.Context, group *model.Groups) error {
	return u.groupRepo.CreateGroup(c, group).Error
}

func (u *groupUsecase) UpdateGroup(c *gin.Context, group *model.Groups) error {
	return u.groupRepo.UpdateGroup(c, group).Error
}

func (u *groupUsecase) DeleteGroup(c *gin.Context, uuid string) error {
	return u.groupRepo.DeleteGroup(c, uuid).Error
}

func (u *groupUsecase) ListGroups(c *gin.Context, filter repository.GroupQueryFilter) ([]model.Groups, error) {
	return u.groupRepo.ListGroups(c, filter)
}

// --- Member usecase ---

type MemberUC interface {
	GetMembers(c *gin.Context) []model.Members
	GetMemberByUUID(c *gin.Context, uuid string) (model.Members, error)
	CreateMember(c *gin.Context, member *model.Members) interface{}
	DeleteMember(c *gin.Context, uuid string) interface{}
	ListMembers(c *gin.Context, filter repository.MemberQueryFilter) ([]model.Members, error)
}

type memberUsecase struct {
	repo repository.Member
}

// NewMember creates a Member usecase backed by the given repository.
func NewMember(repo repository.Member) MemberUC {
	return &memberUsecase{repo: repo}
}

func (u *memberUsecase) GetMembers(c *gin.Context) []model.Members {
	return u.repo.GetMembers(c)
}

func (u *memberUsecase) GetMemberByUUID(c *gin.Context, uuid string) (model.Members, error) {
	return u.repo.GetMemberByUUID(c, uuid)
}

func (u *memberUsecase) CreateMember(c *gin.Context, member *model.Members) interface{} {
	return u.repo.CreateMember(c, member)
}

func (u *memberUsecase) DeleteMember(c *gin.Context, uuid string) interface{} {
	return u.repo.DeleteMember(c, uuid)
}

func (u *memberUsecase) ListMembers(c *gin.Context, filter repository.MemberQueryFilter) ([]model.Members, error) {
	return u.repo.ListMembers(c, filter)
}
