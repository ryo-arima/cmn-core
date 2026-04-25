package repository

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"gorm.io/gorm"
)

// Query filter types

type UserQueryFilter struct {
	UUID  string
	Email string
	Name  string
}

type GroupQueryFilter struct {
	UUID string
	Name string
}

type MemberQueryFilter struct {
	UUID      string
	GroupUUID string
	UserUUID  string
}

// User repository interface

type User interface {
	GetUsers(c *gin.Context) []model.Users
	CreateUser(c *gin.Context, user model.Users) model.Users
	UpdateUser(c *gin.Context, user model.Users) model.Users
	DeleteUser(c *gin.Context, user model.Users) model.Users
	ListUsers(c *gin.Context, filter UserQueryFilter) ([]model.Users, error)
	CountUsers(c *gin.Context, filter UserQueryFilter) (int64, error)
}

// Group repository interface

type Group interface {
	GetGroups(c *gin.Context) []model.Groups
	GetGroupByUUID(c *gin.Context, uuid string) (model.Groups, error)
	GetGroupByID(c *gin.Context, id uint) (model.Groups, error)
	CreateGroup(c *gin.Context, group *model.Groups) *gorm.DB
	UpdateGroup(c *gin.Context, group *model.Groups) *gorm.DB
	DeleteGroup(c *gin.Context, uuid string) *gorm.DB
	ListGroups(c *gin.Context, filter GroupQueryFilter) ([]model.Groups, error)
	CountGroups(c *gin.Context, filter GroupQueryFilter) (int64, error)
}

// Member repository interface

type Member interface {
	GetMembers(c *gin.Context) []model.Members
	GetMemberByUUID(c *gin.Context, uuid string) (model.Members, error)
	CreateMember(c *gin.Context, member *model.Members) interface{}
	UpdateMember(c *gin.Context, member *model.Members) interface{}
	DeleteMember(c *gin.Context, uuid string) interface{}
	ListMembers(c *gin.Context, filter MemberQueryFilter) ([]model.Members, error)
	CountMembers(c *gin.Context, filter MemberQueryFilter) (int64, error)
}

// ResourceQueryFilter holds optional filter conditions for resource queries.
type ResourceQueryFilter struct {
	CreatedBy  string
	GroupUUIDs []string // user's group UUIDs for membership-based access
}

// Resource repository interface.
type Resource interface {
	GetResourceByUUID(ctx context.Context, uuid string) (*model.Resource, error)
	ListResources(ctx context.Context, filter ResourceQueryFilter) ([]model.Resource, error)
	ListAllResources(ctx context.Context) ([]model.Resource, error)
	CreateResource(ctx context.Context, resource *model.Resource) error
	UpdateResource(ctx context.Context, resource *model.Resource) error
	SoftDeleteResource(ctx context.Context, resource *model.Resource, deletedBy string) error
	// Group-role management
	GetGroupRoles(ctx context.Context, resourceUUID string) ([]model.ResourceGroupRole, error)
	SetGroupRole(ctx context.Context, rgr *model.ResourceGroupRole) error
	DeleteGroupRole(ctx context.Context, resourceUUID, groupUUID string) error
}
