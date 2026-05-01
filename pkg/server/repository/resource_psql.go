package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/entity/request"
	"gorm.io/gorm"
)

// Resource repository interface.
type Resource interface {
	GetResourceByUUID(ctx context.Context, uuid string) (*model.PgResource, error)
	ListResources(ctx context.Context, filter request.LoResourceQueryFilter) ([]model.PgResource, error)
	ListAllResources(ctx context.Context) ([]model.PgResource, error)
	CreateResource(ctx context.Context, resource *model.PgResource) error
	UpdateResource(ctx context.Context, resource *model.PgResource) error
	SoftDeleteResource(ctx context.Context, resource *model.PgResource, deletedBy string) error
	// Group-role management
	GetGroupRoles(ctx context.Context, resourceUUID string) ([]model.PgResourceGroupRole, error)
	SetGroupRole(ctx context.Context, rgr *model.PgResourceGroupRole) error
	DeleteGroupRole(ctx context.Context, resourceUUID, groupID string) error
}

type resourceRepository struct {
	db *gorm.DB
}

// NewResource creates a new Resource repository backed by the given database connection.
func NewResource(conf config.BaseConfig) Resource {
	return &resourceRepository{db: conf.DBConnection}
}

func (rcvr *resourceRepository) GetResourceByUUID(ctx context.Context, uuid string) (*model.PgResource, error) {
	var res model.PgResource
	if err := rcvr.db.WithContext(ctx).Where("uuid = ? AND deleted_at IS NULL", uuid).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("resource not found: %s", uuid)
		}
		return nil, err
	}
	return &res, nil
}

// ListResources returns resources accessible to the caller:
//   - all resources whose created_by matches filter.CreatedBy, or
//   - all resources that have a ResourceGroupRole entry for any of filter.GroupUUIDs.
func (rcvr *resourceRepository) ListResources(ctx context.Context, filter request.LoResourceQueryFilter) ([]model.PgResource, error) {
	var resources []model.PgResource

	query := rcvr.db.WithContext(ctx).Where("deleted_at IS NULL")

	if len(filter.GroupIDs) > 0 {
		// created_by matches OR resource has a group-role entry for one of the user's groups
		query = query.Where(
			"created_by = ? OR uuid IN (SELECT resource_uuid FROM pg_resource_group_roles WHERE group_id IN ?)",
			filter.CreatedBy, filter.GroupIDs,
		)
	} else {
		query = query.Where("created_by = ?", filter.CreatedBy)
	}

	if err := query.Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

// ListAllResources returns every non-deleted resource (admin only).
func (rcvr *resourceRepository) ListAllResources(ctx context.Context) ([]model.PgResource, error) {
	var resources []model.PgResource
	if err := rcvr.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func (rcvr *resourceRepository) CreateResource(ctx context.Context, resource *model.PgResource) error {
	return rcvr.db.WithContext(ctx).Create(resource).Error
}

func (rcvr *resourceRepository) UpdateResource(ctx context.Context, resource *model.PgResource) error {
	return rcvr.db.WithContext(ctx).Save(resource).Error
}

func (rcvr *resourceRepository) SoftDeleteResource(ctx context.Context, resource *model.PgResource, deletedBy string) error {
	now := time.Now()
	resource.DeletedBy = deletedBy
	resource.DeletedAt = &now
	return rcvr.db.WithContext(ctx).Save(resource).Error
}

func (rcvr *resourceRepository) GetGroupRoles(ctx context.Context, resourceUUID string) ([]model.PgResourceGroupRole, error) {
	var roles []model.PgResourceGroupRole
	if err := rcvr.db.WithContext(ctx).Where("resource_uuid = ?", resourceUUID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// SetGroupRole inserts or updates (upsert) a group-role entry for a resource.
func (rcvr *resourceRepository) SetGroupRole(ctx context.Context, rgr *model.PgResourceGroupRole) error {
	var existing model.PgResourceGroupRole
	err := rcvr.db.WithContext(ctx).
		Where("resource_uuid = ? AND group_id = ?", rgr.ResourceUUID, rgr.GroupID).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return rcvr.db.WithContext(ctx).Create(rgr).Error
	}
	if err != nil {
		return err
	}
	existing.Role = rgr.Role
	return rcvr.db.WithContext(ctx).Save(&existing).Error
}

func (rcvr *resourceRepository) DeleteGroupRole(ctx context.Context, resourceUUID, groupID string) error {
	return rcvr.db.WithContext(ctx).
		Where("resource_uuid = ? AND group_id = ?", resourceUUID, groupID).
		Delete(&model.PgResourceGroupRole{}).Error
}
