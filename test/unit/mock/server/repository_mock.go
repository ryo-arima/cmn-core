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

func (m *MockUserRepository) GetUsers(c *gin.Context) []model.PgUsers {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(c)
	}
	return m.Users
}

func (m *MockUserRepository) CreateUser(c *gin.Context, user model.PgUsers) model.PgUsers {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(c, user)
	}
	user.ID = uint(len(m.Users) + 1)
	m.Users = append(m.Users, user)
	return user
}

func (m *MockUserRepository) UpdateUser(c *gin.Context, user model.PgUsers) model.PgUsers {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(c, user)
	}
	for i, u := range m.Users {
		if u.ID == user.ID {
			m.Users[i] = user
			return user
		}
	}
	return model.PgUsers{}
}

func (m *MockUserRepository) DeleteUser(c *gin.Context, user model.PgUsers) model.PgUsers {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(c, user)
	}
	for i, u := range m.Users {
		if u.ID == user.ID {
			m.Users = append(m.Users[:i], m.Users[i+1:]...)
			return user
		}
	}
	return model.PgUsers{}
}

func (m *MockUserRepository) ListUsers(c *gin.Context, filter repository.UserQueryFilter) ([]model.PgUsers, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(c, filter)
	}
	return m.Users, nil
}

func (m *MockUserRepository) CountUsers(c *gin.Context, filter repository.UserQueryFilter) (int64, error) {
	if m.CountUsersFunc != nil {
		return m.CountUsersFunc(c, filter)
	}
	return int64(len(m.Users)), nil
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

func (m *MockGroupRepository) GetGroups(c *gin.Context) []model.PgGroups {
	if m.GetGroupsFunc != nil {
		return m.GetGroupsFunc(c)
	}
	return m.Groups
}

func (m *MockGroupRepository) GetGroupByUUID(c *gin.Context, uuid string) (model.PgGroups, error) {
	if m.GetGroupByUUIDFunc != nil {
		return m.GetGroupByUUIDFunc(c, uuid)
	}
	for _, g := range m.Groups {
		if g.UUID == uuid {
			return g, nil
		}
	}
	return model.PgGroups{}, fmt.Errorf("group not found")
}

func (m *MockGroupRepository) GetGroupByID(c *gin.Context, id uint) (model.PgGroups, error) {
	if m.GetGroupByIDFunc != nil {
		return m.GetGroupByIDFunc(c, id)
	}
	for _, g := range m.Groups {
		if g.ID == id {
			return g, nil
		}
	}
	return model.PgGroups{}, fmt.Errorf("group not found")
}

func (m *MockGroupRepository) CreateGroup(c *gin.Context, group *model.PgGroups) *gorm.DB {
	if m.CreateGroupFunc != nil {
		m.CreateGroupFunc(c, group)
	} else {
		group.ID = uint(len(m.Groups) + 1)
		m.Groups = append(m.Groups, *group)
	}
	return &gorm.DB{}
}

func (m *MockGroupRepository) UpdateGroup(c *gin.Context, group *model.PgGroups) *gorm.DB {
	if m.UpdateGroupFunc != nil {
		m.UpdateGroupFunc(c, group)
	} else {
		for i, g := range m.Groups {
			if g.ID == group.ID {
				m.Groups[i] = *group
				break
			}
		}
	}
	return &gorm.DB{}
}

func (m *MockGroupRepository) DeleteGroup(c *gin.Context, uuid string) *gorm.DB {
	if m.DeleteGroupFunc != nil {
		m.DeleteGroupFunc(c, uuid)
	} else {
		for i, g := range m.Groups {
			if g.UUID == uuid {
				m.Groups = append(m.Groups[:i], m.Groups[i+1:]...)
				break
			}
		}
	}
	return &gorm.DB{}
}

func (m *MockGroupRepository) ListGroups(c *gin.Context, filter repository.GroupQueryFilter) ([]model.PgGroups, error) {
	if m.ListGroupsFunc != nil {
		return m.ListGroupsFunc(c, filter)
	}
	return m.Groups, nil
}

func (m *MockGroupRepository) CountGroups(c *gin.Context, filter repository.GroupQueryFilter) (int64, error) {
	if m.CountGroupsFunc != nil {
		return m.CountGroupsFunc(c, filter)
	}
	return int64(len(m.Groups)), nil
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

func (m *MockMemberRepository) GetMembers(c *gin.Context) []model.PgMembers {
	if m.GetMembersFunc != nil {
		return m.GetMembersFunc(c)
	}
	return m.Members
}

func (m *MockMemberRepository) GetMemberByUUID(c *gin.Context, uuid string) (model.PgMembers, error) {
	if m.GetMemberByUUIDFunc != nil {
		return m.GetMemberByUUIDFunc(c, uuid)
	}
	for _, mem := range m.Members {
		if mem.UUID == uuid {
			return mem, nil
		}
	}
	return model.PgMembers{}, fmt.Errorf("member not found")
}

func (m *MockMemberRepository) CreateMember(c *gin.Context, member *model.PgMembers) interface{} {
	if m.CreateMemberFunc != nil {
		return m.CreateMemberFunc(c, member)
	}
	member.ID = uint(len(m.Members) + 1)
	m.Members = append(m.Members, *member)
	return nil
}

func (m *MockMemberRepository) UpdateMember(c *gin.Context, member *model.PgMembers) interface{} {
	if m.UpdateMemberFunc != nil {
		return m.UpdateMemberFunc(c, member)
	}
	for i, mem := range m.Members {
		if mem.ID == member.ID {
			m.Members[i] = *member
			return nil
		}
	}
	return nil
}

func (m *MockMemberRepository) DeleteMember(c *gin.Context, uuid string) interface{} {
	if m.DeleteMemberFunc != nil {
		return m.DeleteMemberFunc(c, uuid)
	}
	for i, mem := range m.Members {
		if mem.UUID == uuid {
			m.Members = append(m.Members[:i], m.Members[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockMemberRepository) ListMembers(c *gin.Context, filter repository.MemberQueryFilter) ([]model.PgMembers, error) {
	if m.ListMembersFunc != nil {
		return m.ListMembersFunc(c, filter)
	}
	return m.Members, nil
}

func (m *MockMemberRepository) CountMembers(c *gin.Context, filter repository.MemberQueryFilter) (int64, error) {
	if m.CountMembersFunc != nil {
		return m.CountMembersFunc(c, filter)
	}
	return int64(len(m.Members)), nil
}

// MockCommonRepository implements repository.Common for testing
type MockCommonRepository struct {
	JWTSecret         string
	ValidateTokenFunc func(ctx context.Context, token string) (*model.LoJWTClaims, error)
}

func (m *MockCommonRepository) ValidateToken(ctx context.Context, tokenString string) (*model.LoJWTClaims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(ctx, tokenString)
	}
	return &model.LoJWTClaims{
		Email: "test@example.com",
		Role:  "user",
		UUID:  "test-uuid",
	}, nil
}

func (m *MockCommonRepository) GetBaseConfig() config.BaseConfig {
	return config.BaseConfig{}
}

func (m *MockCommonRepository) ResolveRole(email string) string {
	return "user"
}
