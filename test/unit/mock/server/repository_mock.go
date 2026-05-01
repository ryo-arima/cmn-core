package mock

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/server/repository"
	"gorm.io/gorm"
)

// MockUserRepository implements repository.User for testing
type MockUserRepository struct {
	Users          []model.PgUsers
	GetUsersFunc   func(c *gin.Context) []model.PgUsers
	CreateUserFunc func(c *gin.Context, user model.PgUsers) model.PgUsers
	UpdateUserFunc func(c *gin.Context, user model.PgUsers) model.PgUsers
	DeleteUserFunc func(c *gin.Context, user model.PgUsers) model.PgUsers
	ListUsersFunc  func(c *gin.Context, filter repository.UserQueryFilter) ([]model.PgUsers, error)
	CountUsersFunc func(c *gin.Context, filter repository.UserQueryFilter) (int64, error)
}

func (rcvr *MockUserRepository) GetUsers(c *gin.Context) []model.PgUsers {
	if rcvr.GetUsersFunc != nil {
		return rcvr.GetUsersFunc(c)
	}
	return rcvr.Users
}

func (rcvr *MockUserRepository) CreateUser(c *gin.Context, user model.PgUsers) model.PgUsers {
	if rcvr.CreateUserFunc != nil {
		return rcvr.CreateUserFunc(c, user)
	}
	user.ID = uint(len(rcvr.Users) + 1)
	rcvr.Users = append(rcvr.Users, user)
	return user
}

func (rcvr *MockUserRepository) UpdateUser(c *gin.Context, user model.PgUsers) model.PgUsers {
	if rcvr.UpdateUserFunc != nil {
		return rcvr.UpdateUserFunc(c, user)
	}
	for i, u := range rcvr.Users {
		if u.ID == user.ID {
			rcvr.Users[i] = user
			return user
		}
	}
	return model.PgUsers{}
}

func (rcvr *MockUserRepository) DeleteUser(c *gin.Context, user model.PgUsers) model.PgUsers {
	if rcvr.DeleteUserFunc != nil {
		return rcvr.DeleteUserFunc(c, user)
	}
	for i, u := range rcvr.Users {
		if u.ID == user.ID {
			rcvr.Users = append(rcvr.Users[:i], rcvr.Users[i+1:]...)
			return user
		}
	}
	return model.PgUsers{}
}

func (rcvr *MockUserRepository) ListUsers(c *gin.Context, filter repository.UserQueryFilter) ([]model.PgUsers, error) {
	if rcvr.ListUsersFunc != nil {
		return rcvr.ListUsersFunc(c, filter)
	}
	return rcvr.Users, nil
}

func (rcvr *MockUserRepository) CountUsers(c *gin.Context, filter repository.UserQueryFilter) (int64, error) {
	if rcvr.CountUsersFunc != nil {
		return rcvr.CountUsersFunc(c, filter)
	}
	return int64(len(rcvr.Users)), nil
}

// MockGroupRepository implements repository.Group for testing
type MockGroupRepository struct {
	Groups             []model.PgGroups
	GetGroupsFunc      func(c *gin.Context) []model.PgGroups
	GetGroupByUUIDFunc func(c *gin.Context, uuid string) (model.PgGroups, error)
	GetGroupByIDFunc   func(c *gin.Context, id uint) (model.PgGroups, error)
	CreateGroupFunc    func(c *gin.Context, group *model.PgGroups) error
	UpdateGroupFunc    func(c *gin.Context, group *model.PgGroups) error
	DeleteGroupFunc    func(c *gin.Context, uuid string) error
	ListGroupsFunc     func(c *gin.Context, filter repository.GroupQueryFilter) ([]model.PgGroups, error)
	CountGroupsFunc    func(c *gin.Context, filter repository.GroupQueryFilter) (int64, error)
}

func (rcvr *MockGroupRepository) GetGroups(c *gin.Context) []model.PgGroups {
	if rcvr.GetGroupsFunc != nil {
		return rcvr.GetGroupsFunc(c)
	}
	return rcvr.Groups
}

func (rcvr *MockGroupRepository) GetGroupByUUID(c *gin.Context, uuid string) (model.PgGroups, error) {
	if rcvr.GetGroupByUUIDFunc != nil {
		return rcvr.GetGroupByUUIDFunc(c, uuid)
	}
	for _, g := range rcvr.Groups {
		if g.UUID == uuid {
			return g, nil
		}
	}
	return model.PgGroups{}, fmt.Errorf("group not found")
}

func (rcvr *MockGroupRepository) GetGroupByID(c *gin.Context, id uint) (model.PgGroups, error) {
	if rcvr.GetGroupByIDFunc != nil {
		return rcvr.GetGroupByIDFunc(c, id)
	}
	for _, g := range rcvr.Groups {
		if g.ID == id {
			return g, nil
		}
	}
	return model.PgGroups{}, fmt.Errorf("group not found")
}

func (rcvr *MockGroupRepository) CreateGroup(c *gin.Context, group *model.PgGroups) *gorm.DB {
	if rcvr.CreateGroupFunc != nil {
		rcvr.CreateGroupFunc(c, group)
	} else {
		group.ID = uint(len(rcvr.Groups) + 1)
		rcvr.Groups = append(rcvr.Groups, *group)
	}
	return &gorm.DB{}
}

func (rcvr *MockGroupRepository) UpdateGroup(c *gin.Context, group *model.PgGroups) *gorm.DB {
	if rcvr.UpdateGroupFunc != nil {
		rcvr.UpdateGroupFunc(c, group)
	} else {
		for i, g := range rcvr.Groups {
			if g.ID == group.ID {
				rcvr.Groups[i] = *group
				break
			}
		}
	}
	return &gorm.DB{}
}

func (rcvr *MockGroupRepository) DeleteGroup(c *gin.Context, uuid string) *gorm.DB {
	if rcvr.DeleteGroupFunc != nil {
		rcvr.DeleteGroupFunc(c, uuid)
	} else {
		for i, g := range rcvr.Groups {
			if g.UUID == uuid {
				rcvr.Groups = append(rcvr.Groups[:i], rcvr.Groups[i+1:]...)
				break
			}
		}
	}
	return &gorm.DB{}
}

func (rcvr *MockGroupRepository) ListGroups(c *gin.Context, filter repository.GroupQueryFilter) ([]model.PgGroups, error) {
	if rcvr.ListGroupsFunc != nil {
		return rcvr.ListGroupsFunc(c, filter)
	}
	return rcvr.Groups, nil
}

func (rcvr *MockGroupRepository) CountGroups(c *gin.Context, filter repository.GroupQueryFilter) (int64, error) {
	if rcvr.CountGroupsFunc != nil {
		return rcvr.CountGroupsFunc(c, filter)
	}
	return int64(len(rcvr.Groups)), nil
}

// MockMemberRepository implements repository.Member for testing
type MockMemberRepository struct {
	Members             []model.PgMembers
	GetMembersFunc      func(c *gin.Context) []model.PgMembers
	CreateMemberFunc    func(c *gin.Context, member *model.PgMembers) error
	UpdateMemberFunc    func(c *gin.Context, member *model.PgMembers) error
	DeleteMemberFunc    func(c *gin.Context, uuid string) error
	GetMemberByUUIDFunc func(c *gin.Context, uuid string) (model.PgMembers, error)
	ListMembersFunc     func(c *gin.Context, filter repository.MemberQueryFilter) ([]model.PgMembers, error)
	CountMembersFunc    func(c *gin.Context, filter repository.MemberQueryFilter) (int64, error)
}

func (rcvr *MockMemberRepository) GetMembers(c *gin.Context) []model.PgMembers {
	if rcvr.GetMembersFunc != nil {
		return rcvr.GetMembersFunc(c)
	}
	return rcvr.Members
}

func (rcvr *MockMemberRepository) GetMemberByUUID(c *gin.Context, uuid string) (model.PgMembers, error) {
	if rcvr.GetMemberByUUIDFunc != nil {
		return rcvr.GetMemberByUUIDFunc(c, uuid)
	}
	for _, mem := range rcvr.Members {
		if mem.UUID == uuid {
			return mem, nil
		}
	}
	return model.PgMembers{}, fmt.Errorf("member not found")
}

func (rcvr *MockMemberRepository) CreateMember(c *gin.Context, member *model.PgMembers) interface{} {
	if rcvr.CreateMemberFunc != nil {
		return rcvr.CreateMemberFunc(c, member)
	}
	member.ID = uint(len(rcvr.Members) + 1)
	rcvr.Members = append(rcvr.Members, *member)
	return nil
}

func (rcvr *MockMemberRepository) UpdateMember(c *gin.Context, member *model.PgMembers) interface{} {
	if rcvr.UpdateMemberFunc != nil {
		return rcvr.UpdateMemberFunc(c, member)
	}
	for i, mem := range rcvr.Members {
		if mem.ID == member.ID {
			rcvr.Members[i] = *member
			return nil
		}
	}
	return nil
}

func (rcvr *MockMemberRepository) DeleteMember(c *gin.Context, uuid string) interface{} {
	if rcvr.DeleteMemberFunc != nil {
		return rcvr.DeleteMemberFunc(c, uuid)
	}
	for i, mem := range rcvr.Members {
		if mem.UUID == uuid {
			rcvr.Members = append(rcvr.Members[:i], rcvr.Members[i+1:]...)
			return nil
		}
	}
	return nil
}

func (rcvr *MockMemberRepository) ListMembers(c *gin.Context, filter repository.MemberQueryFilter) ([]model.PgMembers, error) {
	if rcvr.ListMembersFunc != nil {
		return rcvr.ListMembersFunc(c, filter)
	}
	return rcvr.Members, nil
}

func (rcvr *MockMemberRepository) CountMembers(c *gin.Context, filter repository.MemberQueryFilter) (int64, error) {
	if rcvr.CountMembersFunc != nil {
		return rcvr.CountMembersFunc(c, filter)
	}
	return int64(len(rcvr.Members)), nil
}

// MockCommonRepository implements repository.Common for testing
type MockCommonRepository struct {
	JWTSecret         string
	ValidateTokenFunc func(ctx context.Context, token string) (*model.LoJWTClaims, error)
}

func (rcvr *MockCommonRepository) ValidateToken(ctx context.Context, tokenString string) (*model.LoJWTClaims, error) {
	if rcvr.ValidateTokenFunc != nil {
		return rcvr.ValidateTokenFunc(ctx, tokenString)
	}
	return &model.LoJWTClaims{
		Email: "test@example.com",
		Role:  "user",
		UUID:  "test-uuid",
	}, nil
}

func (rcvr *MockCommonRepository) GetBaseConfig() config.BaseConfig {
	return config.BaseConfig{}
}

func (rcvr *MockCommonRepository) ResolveRole(email string) string {
	return "user"
}
