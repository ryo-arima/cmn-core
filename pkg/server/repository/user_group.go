package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"gorm.io/gorm"
)

// --- User repository ---

type userRepository struct {
	db *gorm.DB
}

// NewUser creates a User repository backed by the database in conf.
func NewUser(conf config.BaseConfig) User {
	return &userRepository{db: conf.DBConnection}
}

func (r *userRepository) GetUsers(c *gin.Context) []model.Users {
	var users []model.Users
	if r.db != nil {
		r.db.Find(&users)
	}
	return users
}

func (r *userRepository) CreateUser(c *gin.Context, user model.Users) model.Users {
	if r.db != nil {
		r.db.Create(&user)
	}
	return user
}

func (r *userRepository) UpdateUser(c *gin.Context, user model.Users) model.Users {
	if r.db != nil {
		r.db.Save(&user)
	}
	return user
}

func (r *userRepository) DeleteUser(c *gin.Context, user model.Users) model.Users {
	if r.db != nil {
		r.db.Delete(&user)
	}
	return user
}

func (r *userRepository) ListUsers(c *gin.Context, filter UserQueryFilter) ([]model.Users, error) {
	var users []model.Users
	if r.db == nil {
		return users, nil
	}
	q := r.db
	if filter.UUID != "" {
		q = q.Where("uuid = ?", filter.UUID)
	}
	if filter.Email != "" {
		q = q.Where("email = ?", filter.Email)
	}
	if filter.Name != "" {
		q = q.Where("name = ?", filter.Name)
	}
	return users, q.Find(&users).Error
}

func (r *userRepository) CountUsers(c *gin.Context, filter UserQueryFilter) (int64, error) {
	if r.db == nil {
		return 0, nil
	}
	var count int64
	q := r.db.Model(&model.Users{})
	if filter.UUID != "" {
		q = q.Where("uuid = ?", filter.UUID)
	}
	if filter.Email != "" {
		q = q.Where("email = ?", filter.Email)
	}
	if filter.Name != "" {
		q = q.Where("name = ?", filter.Name)
	}
	return count, q.Count(&count).Error
}

// --- Group repository ---

type groupRepository struct {
	db *gorm.DB
}

// NewGroup creates a Group repository backed by the database in conf.
func NewGroup(conf config.BaseConfig) Group {
	return &groupRepository{db: conf.DBConnection}
}

func (r *groupRepository) GetGroups(c *gin.Context) []model.Groups {
	var groups []model.Groups
	if r.db != nil {
		r.db.Find(&groups)
	}
	return groups
}

func (r *groupRepository) GetGroupByUUID(c *gin.Context, uuid string) (model.Groups, error) {
	var g model.Groups
	if r.db == nil {
		return g, nil
	}
	return g, r.db.Where("uuid = ?", uuid).First(&g).Error
}

func (r *groupRepository) GetGroupByID(c *gin.Context, id uint) (model.Groups, error) {
	var g model.Groups
	if r.db == nil {
		return g, nil
	}
	return g, r.db.First(&g, id).Error
}

func (r *groupRepository) CreateGroup(c *gin.Context, group *model.Groups) *gorm.DB {
	if r.db == nil {
		return &gorm.DB{}
	}
	return r.db.Create(group)
}

func (r *groupRepository) UpdateGroup(c *gin.Context, group *model.Groups) *gorm.DB {
	if r.db == nil {
		return &gorm.DB{}
	}
	return r.db.Save(group)
}

func (r *groupRepository) DeleteGroup(c *gin.Context, uuid string) *gorm.DB {
	if r.db == nil {
		return &gorm.DB{}
	}
	return r.db.Where("uuid = ?", uuid).Delete(&model.Groups{})
}

func (r *groupRepository) ListGroups(c *gin.Context, filter GroupQueryFilter) ([]model.Groups, error) {
	var groups []model.Groups
	if r.db == nil {
		return groups, nil
	}
	q := r.db
	if filter.UUID != "" {
		q = q.Where("uuid = ?", filter.UUID)
	}
	if filter.Name != "" {
		q = q.Where("name = ?", filter.Name)
	}
	return groups, q.Find(&groups).Error
}

func (r *groupRepository) CountGroups(c *gin.Context, filter GroupQueryFilter) (int64, error) {
	if r.db == nil {
		return 0, nil
	}
	var count int64
	q := r.db.Model(&model.Groups{})
	if filter.UUID != "" {
		q = q.Where("uuid = ?", filter.UUID)
	}
	if filter.Name != "" {
		q = q.Where("name = ?", filter.Name)
	}
	return count, q.Count(&count).Error
}

// --- Member repository ---

type memberRepository struct {
	db *gorm.DB
}

// NewMember creates a Member repository backed by the database in conf.
func NewMember(conf config.BaseConfig) Member {
	return &memberRepository{db: conf.DBConnection}
}

func (r *memberRepository) GetMembers(c *gin.Context) []model.Members {
	var members []model.Members
	if r.db != nil {
		r.db.Find(&members)
	}
	return members
}

func (r *memberRepository) GetMemberByUUID(c *gin.Context, uuid string) (model.Members, error) {
	var m model.Members
	if r.db == nil {
		return m, nil
	}
	return m, r.db.Where("uuid = ?", uuid).First(&m).Error
}

func (r *memberRepository) CreateMember(c *gin.Context, member *model.Members) interface{} {
	if r.db == nil {
		return &gorm.DB{}
	}
	return r.db.Create(member)
}

func (r *memberRepository) UpdateMember(c *gin.Context, member *model.Members) interface{} {
	if r.db == nil {
		return &gorm.DB{}
	}
	return r.db.Save(member)
}

func (r *memberRepository) DeleteMember(c *gin.Context, uuid string) interface{} {
	if r.db == nil {
		return &gorm.DB{}
	}
	return r.db.Where("uuid = ?", uuid).Delete(&model.Members{})
}

func (r *memberRepository) ListMembers(c *gin.Context, filter MemberQueryFilter) ([]model.Members, error) {
	var members []model.Members
	if r.db == nil {
		return members, nil
	}
	q := r.db
	if filter.UUID != "" {
		q = q.Where("uuid = ?", filter.UUID)
	}
	if filter.GroupUUID != "" {
		q = q.Where("group_uuid = ?", filter.GroupUUID)
	}
	if filter.UserUUID != "" {
		q = q.Where("user_uuid = ?", filter.UserUUID)
	}
	return members, q.Find(&members).Error
}

func (r *memberRepository) CountMembers(c *gin.Context, filter MemberQueryFilter) (int64, error) {
	if r.db == nil {
		return 0, nil
	}
	var count int64
	q := r.db.Model(&model.Members{})
	if filter.UUID != "" {
		q = q.Where("uuid = ?", filter.UUID)
	}
	if filter.GroupUUID != "" {
		q = q.Where("group_uuid = ?", filter.GroupUUID)
	}
	if filter.UserUUID != "" {
		q = q.Where("user_uuid = ?", filter.UserUUID)
	}
	return count, q.Count(&count).Error
}
